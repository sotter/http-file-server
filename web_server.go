/*
	定时上报的自己的地址给web-trans
 */

package main

import (
	"net/http"
	"log"
	"fmt"
	"time"
	"sync"
	"io/ioutil"
	"encoding/json"
)

const (
	ReportServerUrl = "/reportserver"
	GetServerUrl    = "/getserver"
)

type ServerManager struct {
	mutex sync.RWMutex
	servers map[string]uint32
}

//接受来自waf-server的心跳请求
type ServerInfo struct {
	Host string `json:"host"`
}

func NewServerManager() *ServerManager{
	return &ServerManager{
		servers : make(map[string]uint32),
	}
}

func (this *ServerManager)MarkServerActive(server string) {
	this.mutex.Lock()
	this.servers[server] = uint32(time.Now().Unix())
	this.mutex.Unlock()
}

//Server的超时监测
var timeout = 60
func (this *ServerManager)MarkServerTimeout() {
	fmt.Println("------- check timeout")
	now := uint32(time.Now().Unix())

	this.mutex.Lock()
	for key, value := range this.servers {
		if now - value > uint32(timeout) {
			log.Println("Server ", key, " timeout, delete it")
			delete(this.servers, key)
		}
	}
	this.mutex.Unlock()

	time.AfterFunc(time.Second * 1, this.MarkServerTimeout)
}

func (this *ServerManager) GetAllServers()  []string {
	var servers []string

	this.mutex.Lock()
	for key, _ := range this.servers {
		servers = append(servers, key)
	}
	this.mutex.Unlock()

	return servers
}

//fileserver上报地址
func ReportServer(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != ReportServerUrl {
		w.WriteHeader(404)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var server_info ServerInfo
	if err := json.Unmarshal(body, &server_info); err != nil {
		w.WriteHeader(405)
		return
	}
	server_manager.MarkServerActive(server_info.Host)
}

var server_manager = NewServerManager()
//获取fileserver的地址
func GetServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Into GetServer ")
	if r.RequestURI != GetServerUrl {
		log.Println("r.RequestURI:", r.RequestURI, " != ", GetServerUrl)
		w.WriteHeader(404)
		return
	}

	var body string
	body = body + "<html><head></head><body><style>body{font-size:40px}</style>"
	servers := server_manager.GetAllServers()
	for _, s := range servers {
		//<a href="http://www.w3school.com.cn">W3School</a>
		body = body + "<a href=\"" + s + "\">" + "MyFileServer" +  "</a>"
	}
	body = body + "</body></html>"

	w.Write([]byte(body))
}

func StartWebServer() {
	mux := http.NewServeMux()
	mux.HandleFunc(GetServerUrl, GetServer)
	mux.HandleFunc(ReportServerUrl, ReportServer)

	log.Println("Start WebServer at:", 38000)
	http.ListenAndServe(":38000", mux)
}
