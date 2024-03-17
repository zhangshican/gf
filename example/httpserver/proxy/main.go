// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	PortOfServer1 = 8198
	PortOfServer2 = 8199
	UpStream      = "http://127.0.0.1:8198"
)

// StartServer1 starts Server1: A simple http server for demo.
func StartServer1() {
	s := g.Server(1)
	s.BindHandler("/*", func(r *ghttp.Request) {
		r.Response.Write("response from server 1")
	})
	s.BindHandler("/user/1", func(r *ghttp.Request) {
		r.Response.Write("user info from server 1")
	})
	s.SetPort(PortOfServer1)
	s.Run()
}

// StartServer2 starts Server2:
// All requests to Server2 are directly redirected to Server1.
func StartServer2() {
	s := g.Server(2)
	u, _ := url.Parse(UpStream)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		writer.WriteHeader(http.StatusBadGateway)
	}
	s.BindHandler("/proxy/*url", func(r *ghttp.Request) {
		var (
			originalPath = r.Request.URL.Path
			proxyToPath  = "/" + r.Get("url").String()
		)
		r.Request.URL.Path = proxyToPath
		g.Log().Infof(r.Context(), `server2:"%s" -> server1:"%s"`, originalPath, proxyToPath)
		r.MakeBodyRepeatableRead(false)
		proxy.ServeHTTP(r.Response.Writer.RawWriter(), r.Request)
	})
	s.SetPort(PortOfServer2)
	s.Run()
}

func main() {
	go StartServer1()
	StartServer2()
}
