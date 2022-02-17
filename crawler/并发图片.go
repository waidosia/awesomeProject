package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func HandleError3(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}

func DownLoadFile(url string, filename string) (ok bool) {
	resp, err := http.Get(url)
	HandleError3(err, "http.get.url")
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	HandleError3(err, "resp.body")
	filename = "./images" + filename
	err = ioutil.WriteFile(filename, bytes, 0666)
	if err != nil {
		ok = false
	} else {
		ok = true
	}
	return
}

var (
	chanImagesUrls chan string
	waitGroup      sync.WaitGroup
	chanTask       chan string
	reimg          = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

func GetFilenameFromUrl(url string) (filename string) {
	lastIndex := strings.LastIndex(url, "/")
	filename = url[lastIndex+1:]
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	filename = timePrefix + "_" + filename
	return
}

func DownloadImg() {
	for url := range chanImagesUrls {
		filename := GetFilenameFromUrl(url)
		ok := DownLoadFile(url, filename)
		if ok {
			fmt.Printf("%s 下载成功\n", filename)
		} else {
			fmt.Printf("%s 下载失败\n", filename)
		}
	}
	waitGroup.Done()
}
func GetPageStr2(url string) (pageStr string) {
	resp, err := http.Get(url)
	HandleError3(err, "Http.Get url")
	defer resp.Body.Close()
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError3(err, "ioutil,ReadAll")
	pageStr = string(pageBytes)
	return pageStr
}

func GetImgs(url string) (urls []string) {
	pageStr := GetPageStr2(url)
	re := regexp.MustCompile(reimg)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("共找到%d条结果\n", len(results))
	for _, result := range results {
		url := result[0]
		urls = append(urls, url)
	}
	return
}

func GetImgUrls(url string) {
	urls := GetImgs(url)
	for _, url := range urls {
		chanImagesUrls <- url
	}
	chanTask <- url
	waitGroup.Done()
}

func CheckOK() {
	var count int
	for {
		url := <-chanTask
		fmt.Printf("%s 完成了爬取任务\n", url)
		count++
		if count == 26 {
			close(chanImagesUrls)
			break
		}
	}
	waitGroup.Done()
}

func main() {
	chanImagesUrls = make(chan string, 10000000)
	chanTask = make(chan string, 26)
	for i := 0; i < 27; i++ {
		waitGroup.Add(1)
		go GetImgUrls("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/" + strconv.Itoa(i) + ".html")
	}
	waitGroup.Add(1)
	go CheckOK()
	for i := 0; i < 5; i++ {
		waitGroup.Add(1)
		go DownloadImg()
	}
	waitGroup.Wait()
}
