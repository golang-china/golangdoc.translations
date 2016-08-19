// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package httputil provides HTTP utility functions, complementing the
// more common ones in the net/http package.

// Package httputil provides HTTP utility functions, complementing the
// more common ones in the net/http package.
package httputil

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "net/http/internal"
    "net/textproto"
    "net/url"
    "strings"
    "sync"
    "time"
)

// ErrLineTooLong is returned when reading malformed chunked data
// with lines that are too long.
var ErrLineTooLong = internal.ErrLineTooLong



var (
	ErrPersistEOF = &http.ProtocolError{ErrorString: "persistent connection closed"}
	ErrClosed     = &http.ProtocolError{ErrorString: "connection closed by user"}
	ErrPipeline   = &http.ProtocolError{ErrorString: "pipeline error"}
)


// A BufferPool is an interface for getting and returning temporary
// byte slices for use by io.CopyBuffer.
type BufferPool interface {
	Get() []byte
	Put([]byte)
}


// A ClientConn sends request and receives headers over an underlying
// connection, while respecting the HTTP keepalive logic. ClientConn
// supports hijacking the connection calling Hijack to
// regain control of the underlying net.Conn and deal with it as desired.
//
// ClientConn is low-level and old. Applications should instead use
// Client or Transport in the net/http package.

// ClientConn 是早期人工编写的 HTTP 实现。
// 它是低级而老旧的，Go 当前的 HTTP 栈不用它。我们应该在 Go 1 前删除它。
//
// 反对使用：请使用 net/http 中的 Client 或 Transport 代替。
type ClientConn struct {
	mu              sync.Mutex // read-write protects the following fields
	c               net.Conn
	r               *bufio.Reader
	re, we          error // read/write errors
	lastbody        io.ReadCloser
	nread, nwritten int
	pipereq         map[*http.Request]uint

	pipe     textproto.Pipeline
	writeReq func(*http.Request, io.Writer) error
}


// ReverseProxy is an HTTP Handler that takes an incoming request and
// sends it to another server, proxying the response back to the
// client.

// ReverseProxy是一个HTTP处理器，它接收进来的请求，然后把请求发送给另外一个服
// 务，并把回复返回给客户端。
type ReverseProxy struct {

	// Director是一个回调函数，它能将请求变成一个新的真实传递的请求。
	// 它的响应会原封不动拷贝并传输到最原始的客户端。
	Director func(*http.Request)

	// Transport用来操作代理请求。
	// 如果为空，默认使用http.DefaultTransport。
	Transport http.RoundTripper

	// FlushInterval代表客户端拷贝回复消息体的刷新间隔时间。
	// 如果设置为zero，则不进行定期的刷新。
	FlushInterval time.Duration

	// ErrorLog specifies an optional logger for errors
	// that occur when attempting to proxy the request.
	// If nil, logging goes to os.Stderr via the log package's
	// standard logger.
	ErrorLog *log.Logger

	// BufferPool optionally specifies a buffer pool to
	// get byte slices for use by io.CopyBuffer when
	// copying HTTP response bodies.
	BufferPool BufferPool
}


// A ServerConn reads requests and sends responses over an underlying
// connection, until the HTTP keepalive logic commands an end. ServerConn
// also allows hijacking the underlying connection by calling Hijack
// to regain control over the connection. ServerConn supports pipe-lining,
// i.e. requests can be read out of sync (but in the same order) while the
// respective responses are sent.
//
// ServerConn is low-level and old. Applications should instead use Server
// in the net/http package.

// ServerConn 是早期人工编写的 HTTP 实现。
// 它是低级而老旧的，Go 当前的 HTTP 栈不用它。具体参见 Server。
//
// 反对使用：请使用 net/http 中的 Server 代替。
type ServerConn struct {
	mu              sync.Mutex // read-write protects the following fields
	c               net.Conn
	r               *bufio.Reader
	re, we          error // read/write errors
	lastbody        io.ReadCloser
	nread, nwritten int
	pipereq         map[*http.Request]uint

	pipe textproto.Pipeline
}


// DumpRequest returns the given request in its HTTP/1.x wire
// representation. It should only be used by servers to debug client
// requests. The returned representation is an approximation only;
// some details of the initial request are lost while parsing it into
// an http.Request. In particular, the order and case of header field
// names are lost. The order of values in multi-valued headers is kept
// intact. HTTP/2 requests are dumped in HTTP/1.x form, not in their
// original binary representations.
//
// If body is true, DumpRequest also returns the body. To do so, it
// consumes req.Body and then replaces it with a new io.ReadCloser
// that yields the same bytes. If DumpRequest returns an error,
// the state of req is undefined.
//
// The documentation for http.Request.Write details which fields
// of req are included in the dump.
func DumpRequest(req *http.Request, body bool) ([]byte, error)

// DumpRequestOut is like DumpRequest but for outgoing client requests. It
// includes any headers that the standard http.Transport adds, such as
// User-Agent.

// DumpRequestOut 和 DumpRequest 一样，但是它用于传出客户端请求。它包含任何标准
// http.Transport 添加的 header，例如 User-Agent。
func DumpRequestOut(req *http.Request, body bool) ([]byte, error)

// DumpResponse is like DumpRequest but dumps a response.

// DumpResponse和DumpRequest一样，但是它取出的是一个response。
func DumpResponse(resp *http.Response, body bool) ([]byte, error)

// NewChunkedReader returns a new chunkedReader that translates the data read
// from r out of HTTP "chunked" format before returning it. The chunkedReader
// returns io.EOF when the final 0-length chunk is read.
//
// NewChunkedReader is not needed by normal applications. The http package
// automatically decodes chunking when reading response bodies.
func NewChunkedReader(r io.Reader) io.Reader

// NewChunkedWriter returns a new chunkedWriter that translates writes into HTTP
// "chunked" format before writing them to w. Closing the returned chunkedWriter
// sends the final 0-length chunk that marks the end of the stream.
//
// NewChunkedWriter is not needed by normal applications. The http
// package adds chunking automatically if handlers don't set a
// Content-Length header. Using NewChunkedWriter inside a handler
// would result in double chunking or chunking with a Content-Length
// length, both of which are wrong.
func NewChunkedWriter(w io.Writer) io.WriteCloser

// NewClientConn returns a new ClientConn reading and writing c.  If r is not
// nil, it is the buffer to use when reading c.
//
// ClientConn is low-level and old. Applications should use Client or
// Transport in the net/http package.

// ClientConn 是早期人工编写的 HTTP 实现。
// 它是低级而老旧的，Go 当前的 HTTP 栈不用它。我们应该在 Go 1 前删除它。
//
// 反对使用：请使用 net/http 中的 Client 或 Transport 代替。
func NewClientConn(c net.Conn, r *bufio.Reader) *ClientConn

// NewProxyClientConn works like NewClientConn but writes Requests
// using Request's WriteProxy method.
//
// New code should not use NewProxyClientConn. See Client or
// Transport in the net/http package instead.

// NewProxyClientConn 是早期人工编写的 HTTP 实现。
// 它是低级而老旧的，Go 当前的 HTTP 栈不用它。我们应该在 Go 1 前删除它。
//
// 反对使用：请使用 net/http 中的 Client 或 Transport 代替。
func NewProxyClientConn(c net.Conn, r *bufio.Reader) *ClientConn

// NewServerConn returns a new ServerConn reading and writing c. If r is not
// nil, it is the buffer to use when reading c.
//
// ServerConn is low-level and old. Applications should instead use Server
// in the net/http package.

// NewServerConn 是早期人工编写的 HTTP 实现。
// 它是低级而老旧的，Go 当前的 HTTP 栈不用它。具体参见 Server。
//
// 反对使用：请使用 net/http 中的 Server 代替。
func NewServerConn(c net.Conn, r *bufio.Reader) *ServerConn

// NewSingleHostReverseProxy returns a new ReverseProxy that routes
// URLs to the scheme, host, and base path provided in target. If the
// target's path is "/base" and the incoming request was for "/dir",
// the target request will be for /base/dir.
// NewSingleHostReverseProxy does not rewrite the Host header.
// To rewrite Host headers, use ReverseProxy directly with a custom
// Director policy.

// NewSingleHostReverseProxy返回一个新的ReverseProxy，它会重写URL的scheme，host
// 和基本的目标路径。如果目标路径是“/base”并且进入的请求的路径是“/dir”，
// 那么最终请求的目标路径就会变成/base/dir。
// NewSingleHostReverseProxy 不会重写 Host header。要重写 Host header 请直接使用
// ReverseProxy 自定的 Director 政策。
func NewSingleHostReverseProxy(target *url.URL) *ReverseProxy

// Close calls Hijack and then also closes the underlying connection

// Close调用Hijack并且关闭底层的连接。
func (*ClientConn) Close() error

// Do is convenience method that writes a request and reads a response.

// Do是一个写请求和读回复很方便的方法。
func (*ClientConn) Do(req *http.Request) (*http.Response, error)

// Hijack detaches the ClientConn and returns the underlying connection as well
// as the read-side bufio which may have some left over data. Hijack may be
// called before the user or Read have signaled the end of the keep-alive
// logic. The user should not call Hijack while Read or Write is in progress.

// Hijack将ClientConn单独分离出来，并且返回底层的连接，以及可能有一些未读数据的
// 缓存的读取器。 Hijack会在读取获取到keep-alive结束信号之前被调用。在Read或者
// Write进行中不可以调用Hijack。
func (*ClientConn) Hijack() (c net.Conn, r *bufio.Reader)

// Pending returns the number of unanswered requests
// that have been sent on the connection.

// Pending返回已经被发送出去但是却没有获取到应答的请求数。
func (*ClientConn) Pending() int

// Read reads the next response from the wire. A valid response might be
// returned together with an ErrPersistEOF, which means that the remote
// requested that this be the last request serviced. Read can be called
// concurrently with Write, but not with another Read.

// Read读取连接上的下个请求。回复有可能和ErrPersistEOF一起返回，如果返回了这个错
// 误， 则代表远端的请求是最后被服务的请求了。Read可以和Write并发调用，但是却不
// 能和其他Read并发调用。
func (*ClientConn) Read(req *http.Request) (resp *http.Response, err error)

// Write writes a request. An ErrPersistEOF error is returned if the connection
// has been closed in an HTTP keepalive sense. If req.Close equals true, the
// keepalive connection is logically closed after this request and the opposing
// server is informed. An ErrUnexpectedEOF indicates the remote closed the
// underlying TCP connection, which is usually considered as graceful close.

// Write负责写请求。如果HTTP长连接已经被关闭了，ErrPersistEOF错误就会被抛出。如
// 果req.Close设置为true， 在通知请求和对应的服务之后，长连接就会被关闭了。
// ErrUnexpectedEOF则表示TCP连接被远端关闭。 在考虑到关闭连接的时候必须考虑到这
// 种情况。
func (*ClientConn) Write(req *http.Request) error

func (*ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request)

// Close calls Hijack and then also closes the underlying connection

// Close调用Hijack，并且关闭底层的连接。
func (*ServerConn) Close() error

// Hijack detaches the ServerConn and returns the underlying connection as well
// as the read-side bufio which may have some left over data. Hijack may be
// called before Read has signaled the end of the keep-alive logic. The user
// should not call Hijack while Read or Write is in progress.
func (*ServerConn) Hijack() (net.Conn, *bufio.Reader)

// Pending returns the number of unanswered requests
// that have been received on the connection.

// Pending返回已经连接上但未应答的请求数。
func (*ServerConn) Pending() int

// Read returns the next request on the wire. An ErrPersistEOF is returned if
// it is gracefully determined that there are no more requests (e.g. after the
// first request on an HTTP/1.0 connection, or after a Connection:close on a
// HTTP/1.1 connection).

// Read返回连接上的下个请求。如果确认了没有更多请求之后，将会返回ErrPersistEOF
// 。（例如，在HTTP/1.0 的第一个请求之后，或者在HTTP/1.1的Connection:close之后）
func (*ServerConn) Read() (*http.Request, error)

// Write writes resp in response to req. To close the connection gracefully, set
// the Response.Close field to true. Write should be considered operational
// until it returns an error, regardless of any errors returned on the Read
// side.

// Write为请求进行回复。为了要更好的关闭连接，该函数将Response.Close设置为true。
// 直到它返回一个错误之前，Write都可以被调用，并且应该要忽略任何读取端的错误。
func (*ServerConn) Write(req *http.Request, resp *http.Response) error

