package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func Get(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	fmt.Println(url)
	if err != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	buf := make([]byte, 4*1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("讀取完畢")
				break
			} else {
				fmt.Println("resp.body.read err = ", err)
				break
			}
		}
		result += string(buf[:n])
	}
	return
}
func SpiderPage(i int, page chan<- int) {
	url := "https://github.com/search?q=go&type=Repositories&p=" + strconv.Itoa(i)
	fmt.Printf("正在爬第%d頁的資料\n", i)
	result, err := Get(url)
	if err != nil {
		fmt.Println("http.get.url err=", err)
		return
	}
	filename := "page" + strconv.Itoa(i) + ".html"
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("os.creact err=", err)
		return
	}
	f.WriteString(result)
	f.Close()
	page <- i
}
func Run(start, end int) {
	fmt.Printf("正在爬取第%d-%d頁\n", start, end)
	page := make(chan int)
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}

	for i := start; i <= end; i++ {
		fmt.Printf("第%d頁已完成\n", <-page)
	}
}

func main() {
	var start, end int
	fmt.Printf("起始頁數: ")
	fmt.Scan(&start)
	fmt.Printf("結束頁數: ")
	fmt.Scan(&end)
	Run(start, end)
}
