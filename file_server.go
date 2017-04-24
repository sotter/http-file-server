package main

import (
	"net/http"
	"os"
	"log"
	"io/ioutil"
	"fmt"
	"strings"
	"time"
	"strconv"
	"encoding/json"
	"bytes"
)

type HttpFileServer struct {
	UrlPath         string       //挂载的urlpath， 有这个有一定隐私的功能
	LocalRootPath   string       //提供文件服务的root path
	Port            int          // Listen Port
	ReportUrl       string       // 上报自己地址的reportUrl
}

const HttpServeRootPath = "./img"

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

var imgSuffix = []string{"jpg", "jpeg", "png",}
func IsImgDir(files []os.FileInfo) bool {
	for _, f := range files {
		fmt.Println("IsImgDir:", f.Name())
		for _, suffix := range imgSuffix {
			if strings.HasSuffix(f.Name(), suffix) {
				return true
			}
		}
	}
	return false
}

func CollectImgToHtml(r *http.Request, files []os.FileInfo) []byte {

	fmt.Println("http Request:", *r)

	var body string
	body = body + "<html><head></head><body>"
	for _, file := range files {
		body = body + "<img src=\"http://" + r.Host + "/" + r.URL.Path[1:] + "/" + file.Name() + "\"  alt=\"111\"/>\n"
	}
	body = body + "</body></html>"

	fmt.Println("Collect Img To Html : ", string(body))
	return []byte(body)
}

func (this *HttpFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(this.LocalRootPath + r.URL.Path)

	relative_path := r.URL.Path
	log.Println("relative_path - 0:", relative_path)

	relative_path = strings.TrimSuffix(relative_path, this.UrlPath)

	log.Println("relative_path:", relative_path)
	handler := http.FileServer(http.Dir(this.LocalRootPath))
	abs_path := this.LocalRootPath + relative_path
	log.Println("abs_path:", abs_path)

	if exist, _ := IsDirectory(abs_path) ; exist == true {
		files, _ := ioutil.ReadDir(abs_path)
		if IsImgDir(files) { // 如果是图片目录， 直接把目录中的所有内容返回
			//CollectImgToHtml(r, files)
			w.Write(CollectImgToHtml(r, files))
			return
		} else {  //如果是目录，那么FileServer 产生的是一个html，修改一下它的格式
			w.Write([]byte("<style>body{font-size:40px}</style>"))
			handler.ServeHTTP(w, r)
			return
		}
	}
	//如果是一个文件，直接将文件中的内容返回
	handler.ServeHTTP(w, r)
}

func (this *HttpFileServer) ReportHost() {
	ip := GetIp("en0")
	if len(ip) < 0 {
		return
	}

	report_host := &ServerInfo {
		Host : "http://" + ip + ":" + strconv.Itoa(this.Port) + this.UrlPath,
	}

	body, err := json.Marshal(report_host)
	if err != nil {
		log.Println("json marsha fail:", err.Error())
		return
	}

	//忽略响应
	http.Post(this.ReportUrl, "application/json", bytes.NewBuffer(body))
	time.AfterFunc(time.Second * 10, this.ReportHost)
}

func (this *HttpFileServer) Start() {
	go this.ReportHost()

	http.Handle(this.UrlPath, this)
	address := ":" + strconv.Itoa(this.Port)

	fmt.Println("Start Http File Server : ",
	"\nAddress          :   ", address,
	"\nUrlPath          :   ", this.UrlPath,
	"\nLocalRootPath    :   ", this.LocalRootPath,
	"\nReportUrl        :   ", this.ReportUrl)

	log.Fatal(http.ListenAndServe(address, nil))
}
