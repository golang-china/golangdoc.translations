// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package fcgi implements the FastCGI protocol. Currently only the responder role
// is supported. The protocol is defined at
// http://www.fastcgi.com/drupal/node/6?q=node/22

// fcgi 包实现了FastCGI协议. 当前只提供了FastCGI的响应服务端。
// 这个协议定义文档是：http://www.fastcgi.com/drupal/node/6?q=node/22
package fcgi

// Serve accepts incoming FastCGI connections on the listener l, creating a new
// goroutine for each. The goroutine reads requests and then calls handler to reply
// to them. If l is nil, Serve accepts connections from os.Stdin. If handler is
// nil, http.DefaultServeMux is used.

// Serve在l监听器中接受传递进来的FastCGI连接，为每个请求创建了一个新的goroutine。
// goroutine读取请求，然后调用handler来回复。
// 如果l是nil，Serve会从os.Stdin接收连接。
// 如果handler是nil，默认使用http.DefaultServeMux。
func Serve(l net.Listener, handler http.Handler) error
