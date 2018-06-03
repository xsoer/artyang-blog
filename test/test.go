package main

import (
	"time"
	"fmt"
	"strconv"
	"os"
	"strings"
)

func main() {
	path := ("./uploads")
    year := time.Now().Year()
    month := time.Now().Format("01")

	path += "/"+ strconv.Itoa(year) + "/"+ month
	fmt.Println(path)
	// if !checkPathIsExist(path) {
	// 	makePath(path)
	// }

	fileName := "12341231243.jpeg"
	fileSplit := strings.Split(fileName, ".")
	fmt.Println(fileSplit[1])
}

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
    err := os.MkdirAll(path, os.ModePerm) //创建文件夹
		if err != nil {
			fmt.Println(err)
			flag = false
        }
    return flag
}