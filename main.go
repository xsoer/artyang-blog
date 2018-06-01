package main

import (
    _ "github.com/go-sql-driver/mysql"
	"fmt"
	"net/http"
	"log"
    "strings"
    "database/sql"
    "math"
    "io"
    "os"
)

const (
    username = "root"
    userpwd  = "Myblog2018$"
    dbname   = "artyang"
)

type Post struct {
    ID   string
    URL string
    Intro string
	Types string
	CreatedAt string
}


func main() {
    http.HandleFunc("/api/hello", sayhelloName) //设置访问的路由
    http.HandleFunc("/x/hi", sayHi) //设置访问的路由
    http.HandleFunc("/x/img/index", imgIndex)
    http.HandleFunc("/x/img/save", imgSave)
    http.HandleFunc("/x/img/upload", upload)
    err := http.ListenAndServe(":8080", nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func imgIndex(w http.ResponseWriter, r *http.Request) {
    db, err := getDB(username, userpwd, dbname)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    rows, err := db.Query("SELECT id, url,intro,types,created_at FROM img")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    var id, url, intro, types, created_at string
    data := make(map[string]interface{})
    posts := []Post{}
    for rows.Next() {
        err = rows.Scan(&id, &url, &intro, &types, &created_at)
        if err == nil {
            // log.Println(id, name)
            posts = append(posts, Post{id, url,intro, types, created_at})
        }
    }
    data["items"] = posts
    data["total"] = 10
    // b, err := json.Marshal(data)
    // if err != nil {
    //     fmt.Println("json.Marshal failed:", err)
    //     return
    // }
    // w.WriteHeader(http.StatusOK)
    // w.Write([]byte(string(b)))
    fmt.Fprint(w, data)
}

// 处理/upload 逻辑  
func upload(w http.ResponseWriter, r *http.Request) {  
    fmt.Println("method:", r.Method) //获取请求的方法  
    /*if r.Method == "GET" {  
        crutime := time.Now().Unix()  
        h := md5.New()  
        io.WriteString(h, strconv.FormatInt(crutime, 10))  
        token := fmt.Sprintf("%x", h.Sum(nil))  
  
        t, _ := template.ParseFiles("upload.gtpl")  
        t.Execute(w, token)  
    } else */  
    fmt.Println(r)  
    if r.Method == "POST" {  
        r.ParseMultipartForm(32 << 20)  
        file, handler, err := r.FormFile("file")  
        if err != nil {  
            fmt.Println(err)  
            return  
        }  
        defer file.Close()  
        fmt.Fprintf(w, "%v", handler.Header)  
  
        //  
        bool_fileexist := checkFileIsExist("./uploads")  
        fmt.Println("check file 1-----------------")  
        fmt.Println("-------------------------bool_fileexist:", bool_fileexist)  
        if bool_fileexist { //如果文件夹存在  
            //f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件  
            f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  
            if err != nil {  
                fmt.Println(err)  
                return  
            }  
            defer f.Close()  
            io.Copy(f, file)  
            fmt.Println("文件夹存在")  
        } else { //不存在文件夹时 先创建文件夹再上传  
            err1 := os.Mkdir("./uploads", os.ModePerm) //创建文件夹  
            if err1 != nil {  
                fmt.Println(err)  
                return  
            }  
  
            fmt.Println("文件夹不存在")  
            fmt.Println("文件夹创建成功！")  
            f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  
            if err != nil {  
                fmt.Println(err)  
                return  
            }  
            defer f.Close()  
            io.Copy(f, file)  
        }  
  
    }  
}  

//检查目录是否存在  
func checkFileIsExist(filename string) bool {  
    var exist = true  
    if _, err := os.Stat(filename); os.IsNotExist(err) {  
        fmt.Print(filename + " not exist")  
        exist = false  
    }  
    return exist  
}  


func imgSave(w http.ResponseWriter, r *http.Request)  {
    if (r.Method != "POST") {
        fmt.Fprint(w, "method Error!")
    }

    fmt.Fprint(w, "2323");
}

func getDB(username, userpwd, dbname string) (*sql.DB, error) {
    dataSourceName := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8", username, userpwd, dbname)
    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        log.Println(err.Error()) //仅仅是显示异常
        return nil, err
    }
    return db, nil
}

func Paginator(page, prepage int, nums int64) map[string]interface{} {

	var firstpage int //前一页地址
	var lastpage int  //后一页地址
	//根据nums总数，和prepage每页数量 生成分页总数
	totalpages := int(math.Ceil(float64(nums) / float64(prepage))) //page总数
	if page > totalpages {
		page = totalpages
	}
	if page <= 0 {
		page = 1
	}
	var pages []int
	switch {
	case page >= totalpages-5 && totalpages > 5: //最后5页
		start := totalpages - 5 + 1
		firstpage = page - 1
		lastpage = int(math.Min(float64(totalpages), float64(page+1)))
		pages = make([]int, 5)
		for i, _ := range pages {
			pages[i] = start + i
		}
	case page >= 3 && totalpages > 5:
		start := page - 3 + 1
		pages = make([]int, 5)
		firstpage = page - 3
		for i, _ := range pages {
			pages[i] = start + i
		}
		firstpage = page - 1
		lastpage = page + 1
	default:
		pages = make([]int, int(math.Min(5, float64(totalpages))))
		for i, _ := range pages {
			pages[i] = i + 1
		}
		firstpage = int(math.Max(float64(1), float64(page-1)))
		lastpage = page + 1
		//fmt.Println(pages)
	}
	paginatorMap := make(map[string]interface{})
	paginatorMap["pages"] = pages
	paginatorMap["totalpages"] = totalpages
	paginatorMap["firstpage"] = firstpage
	paginatorMap["lastpage"] = lastpage
	paginatorMap["currpage"] = page
	return paginatorMap
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()  //解析参数，默认是不会解析的
    fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func sayHi(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()  //解析参数，默认是不会解析的
    fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hi hhahahahh!") //这个写入到w的是输出到客户端的
}

