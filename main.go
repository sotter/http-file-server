package main

import (
	"flag"
	"log"
)

func main() {
	//go StartWebServer()
	root_path := flag.String("path", HttpServeRootPath, "a string")
	port := flag.Int("port", 38001, "an int")
	flag.Parse()

	//首先判断RootPath是否存在，如果不存在，直接退出：
	if exist, _ := IsDirectory(*root_path); exist == false  {
		log.Fatal("RootPath:", *root_path, " is not exist!!!")
		return
	}

	file_server := HttpFileServer{
		UrlPath:       "/",
		LocalRootPath: *root_path,
		Port:          *port,
		ReportUrl:     "http://127.0.0.1:38000/reportserver",
	}

	go file_server.Start()

	StartWebServer()
}