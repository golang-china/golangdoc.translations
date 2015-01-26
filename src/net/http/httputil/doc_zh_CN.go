// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package httputil provides HTTP utility functions, complementing the more common
// ones in the net/http package.

// Package httputil provides HTTP utility
// functions, complementing the more common
// ones in the net/http package.
package httputil

var (
	ErrPersistEOF = &http.ProtocolError{ErrorString: "persistent connection closed"}
	ErrClosed     = &http.ProtocolError{ErrorString: "connection closed by user"}
	ErrPipeline   = &http.ProtocolError{ErrorString: "pipeline error"}
)

// ErrLineTooLong is returned when reading malformed chunked data with lines that
// are too long.

// ErrLineTooLong is returned when reading
// malformed chunked data with lines that
// are too long.
var ErrLineTooLong = internal.ErrLineTooLong

// DumpRequest returns the as-received wire representation of req, optionally
// including the request body, for debugging. DumpRequest is semantically a no-op,
// but in order to dump the body, it reads the body data into memory and changes
// req.Body to refer to the in-memory copy. The documentation for
// http.Request.Write details which fields of req are used.

// DumpRequest返回req的传输结构，可选的包括请求的消息体，调试使用。
// DumpRequest在语义上是非操作性的，但是为了获取出消息体，它会将消息体读取到内存中，
// 并且改变req.Body内存的一个拷贝映射。使用的是req的http.Request.Write属性的文档细节。
func DumpRequest(req *http.Request, body bool) (dump []byte, err error)

// DumpRequestOut is like DumpRequest but includes headers that the standard
// http.Transport adds, such as User-Agent.

// DumpRequestOut和DumpRequest一样，但是包含了header，这个header有标准的http.Transport，
// 比如User-Agent。
func DumpRequestOut(req *http.Request, body bool) ([]byte, error)

// DumpResponse is like DumpRequest but dumps a response.

// DumpResponse和DumpRequest一样，但是它取出的是一个response。
func DumpResponse(resp *http.Response, body bool) (dump []byte, err error)

// NewChunkedReader returns a new chunkedReader that translates the data read from
// r out of HTTP "chunked" format before returning it. The chunkedReader returns
// io.EOF when the final 0-length chunk is read.
//
// NewChunkedReader is not needed by normal applications. The http package
// automatically decodes chunking when reading response bodies.

// NewChunkedReader returns a new
// chunkedReader that translates the data
// read from r out of HTTP "chunked" format
// before returning it. The chunkedReader
// returns io.EOF when the final 0-length
// chunk is read.
//
// NewChunkedReader is not needed by normal
// applications. The http package
// automatically decodes chunking when
// reading response bodies.
func NewChunkedReader(r io.Reader) io.Reader

// NewChunkedWriter returns a new chunkedWriter that translates writes into HTTP
// "chunked" format before writing them to w. Closing the returned chunkedWriter
// sends the final 0-length chunk that marks the end of the stream.
//
// NewChunkedWriter is not needed by normal applications. The http package adds
// chunking automatically if handlers don't set a Content-Length header. Using
// NewChunkedWriter inside a handler would result in double chunking or chunking
// with a Content-Length length, both of which are wrong.

// NewChunkedWriter returns a new
// chunkedWriter that translates writes
// into HTTP "chunked" format before
// writing them to w. Closing the returned
// chunkedWriter sends the final 0-length
// chunk that marks the end of the stream.
//
// NewChunkedWriter is not needed by normal
// applications. The http package adds
// chunking automatically if handlers don't
// set a Content-Length header. Using
// NewChunkedWriter inside a handler would
// result in double chunking or chunking
// with a Content-Length length, both of
// which are wrong.
func NewChunkedWriter(w io.Writer) io.WriteCloser

// A ClientConn sends request and receives headers over an underlying connection,
// while respecting the HTTP keepalive logic. ClientConn supports hijacking the
// connection calling Hijack to regain control of the underlying net.Conn and deal
// with it as desired.
//
// ClientConn is low-level and old. Applications should instead use Client or
// Transport in the net/http package.

// ClientConn从还保持着HTTP
// keepalive的底层连接发送请求，并且接收header。
// ClientConn支持调用Hijack来劫持连接用于获取底层网络连接的控制来处理net.Conn。
//
// ServerConn 是低级而老旧的，应用应当采用 net/http 中的
// Client 或 Transport 来代替。
type ClientConn struct {
	// contains filtered or unexported fields
}

// NewClientConn returns a new ClientConn reading and writing c. If r is not nil,
// it is the buffer to use when reading c.
//
// ClientConn is low-level and old. Applications should use Client or Transport in
// the net/http package.

// NewClientConn返回一个新的ClientConnd对c进行读取和写入。如果r非空，则使用缓存对c进行读取。
//
// ServerConn 是低级而老旧的，应用应当采用 net/http 中的
// Client 或 Transport 来代替。
func NewClientConn(c net.Conn, r *bufio.Reader) *ClientConn

// NewProxyClientConn works like NewClientConn but writes Requests using Request's
// WriteProxy method.
//
// New code should not use NewProxyClientConn. See Client or Transport in the
// net/http package instead.

// NewProxyClientConn像NewClientConn一样，不同的是使用Request的WriteProxy方法对请求进行写操作。
//
// 新代码不应使用 NewProxyClientConn。见 net/http 中的
// Client 或 Transport。
func NewProxyClientConn(c net.Conn, r *bufio.Reader) *ClientConn

// Close calls Hijack and then also closes the underlying connection

// Close调用Hijack并且关闭底层的连接
func (cc *ClientConn) Close() error

// Do is convenience method that writes a request and reads a response.

// Do是一个写请求和读回复很方便的方法。
func (cc *ClientConn) Do(req *http.Request) (resp *http.Response, err error)

// Hijack detaches the ClientConn and returns the underlying connection as well as
// the read-side bufio which may have some left over data. Hijack may be called
// before the user or Read have signaled the end of the keep-alive logic. The user
// should not call Hijack while Read or Write is in progress.

// Hijack将ClientConn单独分离出来，并且返回底层的连接，以及可能有一些未读数据的缓存的读取器。
// Hijack会在读取获取到keep-alive结束信号之前被调用。在Read或者Write进行中不可以调用Hijack。
func (cc *ClientConn) Hijack() (c net.Conn, r *bufio.Reader)

// Pending returns the number of unanswered requests that have been sent on the
// connection.

// Pending返回已经被发送出去但是却没有获取到应答的请求数。
func (cc *ClientConn) Pending() int

// Read reads the next response from the wire. A valid response might be returned
// together with an ErrPersistEOF, which means that the remote requested that this
// be the last request serviced. Read can be called concurrently with Write, but
// not with another Read.

// Read读取连接上的下个请求。回复有可能和ErrPersistEOF一起返回，如果返回了这个错误，
// 则代表远端的请求是最后被服务的请求了。Read可以和Write并发调用，但是却不能和其他Read并发调用。
func (cc *ClientConn) Read(req *http.Request) (resp *http.Response, err error)

// Write writes a request. An ErrPersistEOF error is returned if the connection has
// been closed in an HTTP keepalive sense. If req.Close equals true, the keepalive
// connection is logically closed after this request and the opposing server is
// informed. An ErrUnexpectedEOF indicates the remote closed the underlying TCP
// connection, which is usually considered as graceful close.

// Write负责写请求。如果HTTP长连接已经被关闭了，ErrPersistEOF错误就会被抛出。如果req.Close设置为true，
// 在通知请求和对应的服务之后，长连接就会被关闭了。ErrUnexpectedEOF则表示TCP连接被远端关闭。
// 在考虑到关闭连接的时候必须考虑到这种情况。
func (cc *ClientConn) Write(req *http.Request) (err error)

// ReverseProxy is an HTTP Handler that takes an incoming request and sends it to
// another server, proxying the response back to the client.

// ReverseProxy是一个HTTP处理器，它接收进来的请求，然后把请求发送给另外一个服务，并把回复返回给客户端。
type ReverseProxy struct {
	// Director must be a function which modifies
	// the request into a new request to be sent
	// using Transport. Its response is then copied
	// back to the original client unmodified.
	Director func(*http.Request)

	// The transport used to perform proxy requests.
	// If nil, http.DefaultTransport is used.
	Transport http.RoundTripper

	// FlushInterval specifies the flush interval
	// to flush to the client while copying the
	// response body.
	// If zero, no periodic flushing is done.
	FlushInterval time.Duration

	// ErrorLog specifies an optional logger for errors
	// that occur when attempting to proxy the request.
	// If nil, logging goes to os.Stderr via the log package's
	// standard logger.
	ErrorLog *log.Logger
}

// NewSingleHostReverseProxy returns a new ReverseProxy that rewrites URLs to the
// scheme, host, and base path provided in target. If the target's path is "/base"
// and the incoming request was for "/dir", the target request will be for
// /base/dir.

// NewSingleHostReverseProxy返回一个新的ReverseProxy，它会重写URL的scheme，host
// 和基本的目标路径。如果目标路径是“/base”并且进入的请求的路径是“/dir”，
// 那么最终请求的目标路径就会变成/base/dir。
func NewSingleHostReverseProxy(target *url.URL) *ReverseProxy

func (p *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request)

// A ServerConn reads requests and sends responses over an underlying connection,
// until the HTTP keepalive logic commands an end. ServerConn also allows hijacking
// the underlying connection by calling Hijack to regain control over the
// connection. ServerConn supports pipe-lining, i.e. requests can be read out of
// sync (but in the same order) while the respective responses are sent.
//
// ServerConn is low-level and old. Applications should instead use Server in the
// net/http package.

// ServerConn 在底层连接之上读取请求，发送回复，直到HTTP
// keepalive出现了结束命令。 ServerConn
// 允许靠调用Hijack来对底层连接进行劫持，从而得到连接的控制权。
// ServerConn
// 支持管道连接，例如，当回复发送的时候，请求可以不需要进行同步（但是是在相同的顺序）。
//
// ServerConn 是低级而老旧的，大部分应用都不需要它。具体参见
// Server。
type ServerConn struct {
	// contains filtered or unexported fields
}

// NewServerConn returns a new ServerConn reading and writing c. If r is not nil,
// it is the buffer to use when reading c.
//
// ServerConn is low-level and old. Applications should instead use Server in the
// net/http package.

// NewServerConn返回一个新的ServerConn来读取和写c。如果r非空，则使用缓存对c进行读取。
//
// ServerConn 是低级而老旧的，大部分应用都不需要它。具体参见
// Server。
func NewServerConn(c net.Conn, r *bufio.Reader) *ServerConn

// Close calls Hijack and then also closes the underlying connection

// Close调用Hijack，并且关闭底层的连接。
func (sc *ServerConn) Close() error

// Hijack detaches the ServerConn and returns the underlying connection as well as
// the read-side bufio which may have some left over data. Hijack may be called
// before Read has signaled the end of the keep-alive logic. The user should not
// call Hijack while Read or Write is in progress.

// Hijack将ServerConn单独分离出来，并且返回底层的连接，以及可能有一些未读数据的缓存的读取器。
// Hijack会在读取获取到keep-alive结束信号之前被调用。在Read或者Write进行中不可以调用Hijack。
func (sc *ServerConn) Hijack() (c net.Conn, r *bufio.Reader)

// Pending returns the number of unanswered requests that have been received on the
// connection.

// Pending返回已经连接上但未应答的请求数。
func (sc *ServerConn) Pending() int

// Read returns the next request on the wire. An ErrPersistEOF is returned if it is
// gracefully determined that there are no more requests (e.g. after the first
// request on an HTTP/1.0 connection, or after a Connection:close on a HTTP/1.1
// connection).

// Read返回连接上的下个请求。如果确认了没有更多请求之后，将会返回ErrPersistEOF。（例如，在HTTP/1.0
// 的第一个请求之后，或者在HTTP/1.1的Connection:close之后）
func (sc *ServerConn) Read() (req *http.Request, err error)

// Write writes resp in response to req. To close the connection gracefully, set
// the Response.Close field to true. Write should be considered operational until
// it returns an error, regardless of any errors returned on the Read side.

// Write为请求进行回复。为了要更好的关闭连接，该函数将Response.Close设置为true。
// 直到它返回一个错误之前，Write都可以被调用，并且应该要忽略任何读取端的错误。
func (sc *ServerConn) Write(req *http.Request, resp *http.Response) error
