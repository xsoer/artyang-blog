package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
    "time"
    "math/rand"
)

const (
	username = "root"
	userpwd  = "Myblog2018$"
	dbname   = "artyang"
)

type Post struct {
	ID        string
	URL       string
	Intro     string
	Types     string
	CreatedAt string
}
const (
    KC_RAND_KIND_NUM   = 0  // 纯数字
    KC_RAND_KIND_LOWER = 1  // 小写字母
    KC_RAND_KIND_UPPER = 2  // 大写字母
    KC_RAND_KIND_ALL   = 3  // 数字、大小写字母
)
 
// 随机字符串
func Krand(size int, kind int) []byte {
    ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
    isAll := kind > 2 || kind < 0
    rand.Seed(time.Now().UnixNano())
    for i :=0; i < size; i++ {
        if isAll { // random ikind
            ikind = rand.Intn(3)
        }
        scope, base := kinds[ikind][0], kinds[ikind][1]
        result[i] = uint8(base+rand.Intn(scope))
    }
    return result
}

func main() {
	http.HandleFunc("/api/hello", sayhelloName) //设置访问的路由
	http.HandleFunc("/x/hi", sayHi)             //设置访问的路由
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
			posts = append(posts, Post{id, url, intro, types, created_at})
		}
	}
	data["items"] = posts
	data["total"] = 10
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(b)))
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method != "POST" {
		fmt.Fprint(w, "not POST request")
	}
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	path := ("./uploads")
	year := time.Now().Year()
	month := time.Now().Format("01")

	path += "/" + strconv.Itoa(year) + "/" + month
	if !checkPathIsExist(path) {
		makePath(path)
    }
    fileSplit := strings.Split(handler.Filename, ".")
    fileName := Krand(32, KC_RAND_KIND_ALL)
    filePath := path + "/"+string(fileName[:])+"." + fileSplit[1]
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
    io.Copy(f, file)
    data := make(map[string]interface{})
    data["code"] = 0
    data["url"] = filePath
    b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(b)))
}

//检查目录是否存在
func checkPathIsExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Print(path + " not exist")
		exist = false
	}
	return exist
}

func makePath(path string) bool {
	var flag = true
	err := os.MkdirAll(path, os.ModePerm) //递归创建文件夹
	if err != nil {
		fmt.Println(err)
		flag = false
	}
	return flag
}

func imgSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprint(w, "method Error!")
	}

	fmt.Fprint(w, "2323")
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
		for i := range pages {
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
		for i := range pages {
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
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
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
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hi hhahahahh!") //这个写入到w的是输出到客户端的
}
