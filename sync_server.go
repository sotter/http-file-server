/* 文件同步服务器 **/
package main

import "net/http"

func sync(local_dir string, remote_dir string) bool {
	http.FileSystem()
}
