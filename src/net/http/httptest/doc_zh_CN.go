// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package httptest provides utilities for HTTP testing.

// httptest 包提供HTTP测试的单元工具.
package httptest

// DefaultRemoteAddr is the default remote address to return in RemoteAddr if an
// explicit DefaultRemoteAddr isn't set on ResponseRecorder.

// DefaultRemoteAddr是RemoteAddr返回的默认远端地址。如果没有对ResponseRecorder做地址设置的话，
// DefaultRemoteAddr就作为默认值。
const DefaultRemoteAddr = "1.2.3.4"

// ResponseRecorder is an implementation of http.ResponseWriter that records its
// mutations for later inspection in tests.

// ResponseRecorder是http.ResponseWriter的具体实现，它为进一步的观察记录下了任何变化。
type ResponseRecorder struct {
	Code      int           // the HTTP response code from WriteHeader
	HeaderMap http.Header   // the HTTP response headers
	Body      *bytes.Buffer // if non-nil, the bytes.Buffer to append written data to
	Flushed   bool
	// contains filtered or unexported fields
}

// NewRecorder returns an initialized ResponseRecorder.

// NewRecorder返回一个初始化的ResponseRecorder。
func NewRecorder() *ResponseRecorder

// Flush sets rw.Flushed to true.

// Flush将rw.Flushed设置为true。
func (rw *ResponseRecorder) Flush()

// Header returns the response headers.

// Header返回回复的header。
func (rw *ResponseRecorder) Header() http.Header

// Write always succeeds and writes to rw.Body, if not nil.

// Write总是返回成功，并且如果buf非空的话，它会写数据到rw.Body。
func (rw *ResponseRecorder) Write(buf []byte) (int, error)

// WriteHeader sets rw.Code.

// WriteHeader设置rw.Code
func (rw *ResponseRecorder) WriteHeader(code int)

// A Server is an HTTP server listening on a system-chosen port on the local
// loopback interface, for use in end-to-end HTTP tests.

// Server
// 是一个HTTP服务，它在系统选择的端口上监听请求，并且是在本地的接口监听，
// 它完全是为了点到点的HTTP测试而出现。
type Server struct {
	URL      string // base URL of form http://ipaddr:port with no trailing slash
	Listener net.Listener

	// TLS is the optional TLS configuration, populated with a new config
	// after TLS is started. If set on an unstarted server before StartTLS
	// is called, existing fields are copied into the new config.
	TLS *tls.Config

	// Config may be changed after calling NewUnstartedServer and
	// before Start or StartTLS.
	Config *http.Server
	// contains filtered or unexported fields
}

// NewServer starts and returns a new Server. The caller should call Close when
// finished, to shut it down.

// NewServer开启并且返回一个新的Server。
// 调用者应该当结束的时候调用Close来关闭Server。
func NewServer(handler http.Handler) *Server

// NewTLSServer starts and returns a new Server using TLS. The caller should call
// Close when finished, to shut it down.

// NewTLSServer 开启并且返回了一个使用TLS的新的Server。
// 调用者应该在结束的时候调用Close来关闭它。
func NewTLSServer(handler http.Handler) *Server

// NewUnstartedServer returns a new Server but doesn't start it.
//
// After changing its configuration, the caller should call Start or StartTLS.
//
// The caller should call Close when finished, to shut it down.

// NewUnstartedServer返回一个新的Server实例，但是并不启动这个Server。
//
// 在改变了配置之后，调用者应该调用Start或者StartTLS。
//
// 调用者应该在结束之后调用Close来关闭Server。
func NewUnstartedServer(handler http.Handler) *Server

// Close shuts down the server and blocks until all outstanding requests on this
// server have completed.

// Close关闭server并且阻塞server，知道所有的请求完成之后才继续。
func (s *Server) Close()

// CloseClientConnections closes any currently open HTTP connections to the test
// Server.

// CloseClientConnections关闭任何现有打开的HTTP连接到测试服务器上。
func (s *Server) CloseClientConnections()

// Start starts a server from NewUnstartedServer.

// Start开启从NewUnstartedServer获取到的server。
func (s *Server) Start()

// StartTLS starts TLS on a server from NewUnstartedServer.

// StartTLS开启从NewUnstartedServer获取到的server的TLS。
func (s *Server) StartTLS()
