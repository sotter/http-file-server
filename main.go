package main

import (
	"flag"
	"log"
)

func main() {
	var config struct {
		ServerPort int
		ReportUrl  string
		RootPath   string
	}

	//首先判断RootPath是否存在，如果不存在，直接退出：
	if exist, _ := IsDirectory(config.RootPath); exist == false  {
		log.Fatal("RootPath:", config.RootPath, " is not exist!!!")
		return
	}

	file_server := HttpFileServer{
		UrlPath:       "/",
		LocalRootPath: config.RootPath,
		Port:          config.ServerPort,
		ReportUrl:     "http://127.0.0.1:38000/reportserver",
	}

	go file_server.Start()

	StartWebServer()
}