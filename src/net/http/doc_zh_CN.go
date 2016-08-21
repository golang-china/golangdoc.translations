// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package http provides HTTP client and server implementations.
//
// Get, Head, Post, and PostForm make HTTP (or HTTPS) requests:
//
//     resp, err := http.Get("http://example.com/")
//     ...
//     resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
//     ...
//     resp, err := http.PostForm("http://example.com/form",
//         url.Values{"key": {"Value"}, "id": {"123"}})
//
// The client must close the response body when finished with it:
//
//     resp, err := http.Get("http://example.com/")
//     if err != nil {
//         // handle error
//     }
//     defer resp.Body.Close()
//     body, err := ioutil.ReadAll(resp.Body)
//     // ...
//
// For control over HTTP client headers, redirect policy, and other
// settings, create a Client:
//
//     client := &http.Client{
//         CheckRedirect: redirectPolicyFunc,
//     }
//
//     resp, err := client.Get("http://example.com")
//     // ...
//
//     req, err := http.NewRequest("GET", "http://example.com", nil)
//     // ...
//     req.Header.Add("If-None-Match", `W/"wyzzy"`)
//     resp, err := client.Do(req)
//     // ...
//
// For control over proxies, TLS configuration, keep-alives,
// compression, and other settings, create a Transport:
//
//     tr := &http.Transport{
//         TLSClientConfig:    &tls.Config{RootCAs: pool},
//         DisableCompression: true,
//     }
//     client := &http.Client{Transport: tr}
//     resp, err := client.Get("https://example.com")
//
// Clients and Transports are safe for concurrent use by multiple
// goroutines and for efficiency should only be created once and re-used.
//
// ListenAndServe starts an HTTP server with a given address and handler.
// The handler is usually nil, which means to use DefaultServeMux.
// Handle and HandleFunc add handlers to DefaultServeMux:
//
//     http.Handle("/foo", fooHandler)
//
//     http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//     })
//
//     log.Fatal(http.ListenAndServe(":8080", nil))
//
// More control over the server's behavior is available by creating a
// custom Server:
//
//     s := &http.Server{
//         Addr:           ":8080",
//         Handler:        myHandler,
//         ReadTimeout:    10 * time.Second,
//         WriteTimeout:   10 * time.Second,
//         MaxHeaderBytes: 1 << 20,
//     }
//     log.Fatal(s.ListenAndServe())
//
// The http package has transparent support for the HTTP/2 protocol when
// using HTTPS. Programs that must disable HTTP/2 can do so by setting
// Transport.TLSNextProto (for clients) or Server.TLSNextProto (for
// servers) to a non-nil, empty map. Alternatively, the following GODEBUG
// environment variables are currently supported:
//
//     GODEBUG=http2client=0  # disable HTTP/2 client support
//     GODEBUG=http2server=0  # disable HTTP/2 server support
//     GODEBUG=http2debug=1   # enable verbose HTTP/2 debug logs
//     GODEBUG=http2debug=2   # ... even more verbose, with frame dumps
//
// The GODEBUG variables are not covered by Go's API compatibility promise.
// HTTP/2 support was added in Go 1.6. Please report any issues instead of
// disabling HTTP/2 support: https://golang.org/s/http2bug

// http包提供了HTTP客户端和服务端的实现。
//
// Get、Head、Post和PostForm函数发出HTTP/ HTTPS请求。
//
//     resp, err := http.Get("http://example.com/")
//     ...
//     resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
//     ...
//     resp, err := http.PostForm("http://example.com/form",
//         url.Values{"key": {"Value"}, "id": {"123"}})
//
// 程序在使用完回复后必须关闭回复的主体。
//
//     resp, err := http.Get("http://example.com/")
//     if err != nil {
//         // handle error
//     }
//     defer resp.Body.Close()
//     body, err := ioutil.ReadAll(resp.Body)
//     // ...
//
// 要管理HTTP客户端的头域、重定向策略和其他设置，创建一个Client：
//
//     client := &http.Client{
//         CheckRedirect: redirectPolicyFunc,
//     }
//     resp, err := client.Get("http://example.com")
//     // ...
//     req, err := http.NewRequest("GET", "http://example.com", nil)
//     // ...
//     req.Header.Add("If-None-Match", `W/"wyzzy"`)
//     resp, err := client.Do(req)
//     // ...
//
// 要管理代理、TLS配置、keep-alive、压缩和其他设置，创建一个Transport：
//
//     tr := &http.Transport{
//         TLSClientConfig:    &tls.Config{RootCAs: pool},
//         DisableCompression: true,
//     }
//     client := &http.Client{Transport: tr}
//     resp, err := client.Get("https://example.com")
//
// Client和Transport类型都可以安全的被多个go程同时使用。出于效率考虑，应该一次建
// 立、尽量重用。
//
// ListenAndServe使用指定的监听地址和处理器启动一个HTTP服务端。处理器参数通常是
// nil，这表示采用包变量DefaultServeMux作为处理器。Handle和HandleFunc函数可以向
// DefaultServeMux添加处理器。
//
//     http.Handle("/foo", fooHandler)
//     http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//     })
//     log.Fatal(http.ListenAndServe(":8080", nil))
//
// 要管理服务端的行为，可以创建一个自定义的Server：
//
//     s := &http.Server{
//         Addr:           ":8080",
//         Handler:        myHandler,
//         ReadTimeout:    10 * time.Second,
//         WriteTimeout:   10 * time.Second,
//         MaxHeaderBytes: 1 << 20,
//     }
//     log.Fatal(s.ListenAndServe())
package http

import (
    "bufio"
    "bytes"
    "compress/gzip"
    "crypto/tls"
    "encoding/base64"
    "encoding/binary"
    "errors"
    "fmt"
    "internal/golang.org/x/net/http2/hpack"
    "io"
    "io/ioutil"
    "log"
    "mime"
    "mime/multipart"
    "net"
    "net/http/internal"
    "net/textproto"
    "net/url"
    "os"
    "path"
    "path/filepath"
    "reflect"
    "runtime"
    "sort"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
    "time"
    "unicode/utf8"
)

// DefaultMaxHeaderBytes is the maximum permitted size of the headers in an HTTP
// request. This can be overridden by setting Server.MaxHeaderBytes.
const DefaultMaxHeaderBytes = 1 << 20 // 1 MB


// DefaultMaxIdleConnsPerHost is the default value of Transport's
// MaxIdleConnsPerHost.
const DefaultMaxIdleConnsPerHost = 2

const (
    // StateNew represents a new connection that is expected to
    // send a request immediately. Connections begin at this
    // state and then transition to either StateActive or
    // StateClosed.
    StateNew ConnState = iota

    // StateActive represents a connection that has read 1 or more
    // bytes of a request. The Server.ConnState hook for
    // StateActive fires before the request has entered a handler
    // and doesn't fire again until the request has been
    // handled. After the request is handled, the state
    // transitions to StateClosed, StateHijacked, or StateIdle.
    StateActive

    // StateIdle represents a connection that has finished
    // handling a request and is in the keep-alive state, waiting
    // for a new request. Connections transition from StateIdle
    // to either StateActive or StateClosed.
    StateIdle

    // StateHijacked represents a hijacked connection.
    // This is a terminal state. It does not transition to StateClosed.
    StateHijacked

    // StateClosed represents a closed connection.
    // This is a terminal state. Hijacked connections do not
    // transition to StateClosed.
    StateClosed
)

// HTTP status codes, defined in RFC 2616.
const (
    StatusContinue           = 100
    StatusSwitchingProtocols = 101

    StatusOK                   = 200
    StatusCreated              = 201
    StatusAccepted             = 202
    StatusNonAuthoritativeInfo = 203
    StatusNoContent            = 204
    StatusResetContent         = 205
    StatusPartialContent       = 206

    StatusMultipleChoices   = 300
    StatusMovedPermanently  = 301
    StatusFound             = 302
    StatusSeeOther          = 303
    StatusNotModified       = 304
    StatusUseProxy          = 305
    StatusTemporaryRedirect = 307

    StatusBadRequest                   = 400
    StatusUnauthorized                 = 401
    StatusPaymentRequired              = 402
    StatusForbidden                    = 403
    StatusNotFound                     = 404
    StatusMethodNotAllowed             = 405
    StatusNotAcceptable                = 406
    StatusProxyAuthRequired            = 407
    StatusRequestTimeout               = 408
    StatusConflict                     = 409
    StatusGone                         = 410
    StatusLengthRequired               = 411
    StatusPreconditionFailed           = 412
    StatusRequestEntityTooLarge        = 413
    StatusRequestURITooLong            = 414
    StatusUnsupportedMediaType         = 415
    StatusRequestedRangeNotSatisfiable = 416
    StatusExpectationFailed            = 417
    StatusTeapot                       = 418

    StatusInternalServerError     = 500
    StatusNotImplemented          = 501
    StatusBadGateway              = 502
    StatusServiceUnavailable      = 503
    StatusGatewayTimeout          = 504
    StatusHTTPVersionNotSupported = 505
)

// TimeFormat is the time format to use when generating times in HTTP
// headers. It is like time.RFC1123 but hard-codes GMT as the time
// zone. The time being formatted must be in UTC for Format to
// generate the correct format.
//
// For parsing this time format, see ParseTime.

// TimeFormat is the time format to use with time.Parse and time.Time.Format
// when parsing or generating times in HTTP headers. It is like time.RFC1123 but
// hard codes GMT as the time zone.
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

// DefaultClient is the default Client and is used by Get, Head, and Post.
var DefaultClient = &Client{}

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = NewServeMux()

// DefaultTransport is the default implementation of Transport and is used by
// DefaultClient. It establishes network connections as needed and caches them
// for reuse by subsequent calls. It uses HTTP proxies as directed by the
// $HTTP_PROXY and $NO_PROXY (or $http_proxy and $no_proxy) environment
// variables.
var DefaultTransport RoundTripper = &Transport{
    Proxy: ProxyFromEnvironment,
    Dial: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
    }).Dial,
    TLSHandshakeTimeout: 10 * time.Second,
}

// ErrBodyReadAfterClose is returned when reading a Request or Response Body
// after the body has been closed. This typically happens when the body is read
// after an HTTP Handler calls WriteHeader or Write on its ResponseWriter.
var ErrBodyReadAfterClose = errors.New("http: invalid Read on closed Body")

// ErrHandlerTimeout is returned on ResponseWriter Write calls in handlers which
// have timed out.
var ErrHandlerTimeout = errors.New("http: Handler timeout")

// HTTP请求的解析错误。
//
//     var (
//         ErrWriteAfterFlush = errors.New("Conn.Write called after Flush")
//         ErrBodyNotAllowed  = errors.New("http: request method or response status code does not allow body")
//         ErrHijacked        = errors.New("Conn has been hijacked")
//         ErrContentLength   = errors.New("Conn.Write wrote more than the declared Content-Length")
//     )
//
// 会被HTTP服务端返回的错误。
//
//     var DefaultClient = &Client{}
//
// DefaultClient是用于包函数Get、Head和Post的默认Client。
//
//     var DefaultServeMux = NewServeMux()
//
// DefaultServeMux是用于Serve的默认ServeMux。
//
//     var ErrBodyReadAfterClose = errors.New("http: invalid Read on closed Body")
//
// 在Resquest或Response的Body字段已经关闭后，试图从中读取时，就会返回
// ErrBodyReadAfterClose。这个错误一般发生在：HTTP处理器中调用完ResponseWriter
// 接口的WriteHeader或Write后从请求中读取数据的时候。
//
//     var ErrHandlerTimeout = errors.New("http: Handler timeout")
//
// 在处理器超时以后调用ResponseWriter接口的Write方法，就会返回ErrHandlerTimeout
// 。
//
//     var ErrLineTooLong = errors.New("header line too long")
//
//     var ErrMissingFile = errors.New("http: no such file")
//
// 当请求中没有提供给FormFile函数的文件字段名，或者该字段名不是文件字段时，该函
// 数就会返回ErrMissingFile。
//
//     var ErrNoCookie = errors.New("http: named cookie not present")
//
//     var ErrNoLocation = errors.New("http: no Location header in response")
var (
    ErrHeaderTooLong        = &ProtocolError{"header too long"}
    ErrShortBody            = &ProtocolError{"entity body too short"}
    ErrNotSupported         = &ProtocolError{"feature not supported"}
    ErrUnexpectedTrailer    = &ProtocolError{"trailer header without chunked transfer encoding"}
    ErrMissingContentLength = &ProtocolError{"missing ContentLength in HEAD response"}
    ErrNotMultipart         = &ProtocolError{"request Content-Type isn't multipart/form-data"}
    ErrMissingBoundary      = &ProtocolError{"no multipart boundary param in Content-Type"}
)

// ErrLineTooLong is returned when reading request or response bodies with
// malformed chunked encoding.
var ErrLineTooLong = internal.ErrLineTooLong

// ErrMissingFile is returned by FormFile when the provided file field name is
// either not present in the request or not a file field.
var ErrMissingFile = errors.New("http: no such file")

// ErrNoCookie is returned by Request's Cookie method when a cookie is not
// found.
var ErrNoCookie = errors.New("http: named cookie not present")

// ErrNoLocation is returned by Response's Location method
// when no Location header is present.
var ErrNoLocation = errors.New("http: no Location header in response")

// Errors introduced by the HTTP server.
var (
    ErrWriteAfterFlush = errors.New("Conn.Write called after Flush")
    ErrBodyNotAllowed  = errors.New("http: request method or response status code does not allow body")
    ErrHijacked        = errors.New("Conn has been hijacked")
    ErrContentLength   = errors.New("Conn.Write wrote more than the declared Content-Length")
)

// A Client is an HTTP client. Its zero value (DefaultClient) is a
// usable client that uses DefaultTransport.
//
// The Client's Transport typically has internal state (cached TCP
// connections), so Clients should be reused instead of created as
// needed. Clients are safe for concurrent use by multiple goroutines.
//
// A Client is higher-level than a RoundTripper (such as Transport)
// and additionally handles HTTP details such as cookies and
// redirects.

// Client类型代表HTTP客户端。它的零值（DefaultClient）是一个可用的使用
// DefaultTransport的客户端。
//
// Client的Transport字段一般会含有内部状态（缓存TCP连接），因此Client类型值应尽
// 量被重用而不是每次需要都创建新的。Client类型值可以安全的被多个go程同时使用。
//
// Client类型的层次比RoundTripper接口（如Transport）高，还会管理HTTP的cookie和重
// 定向等细节。
type Client struct {
    // Transport specifies the mechanism by which individual
    // HTTP requests are made.
    // If nil, DefaultTransport is used.
    Transport RoundTripper

    // CheckRedirect specifies the policy for handling redirects.
    // If CheckRedirect is not nil, the client calls it before
    // following an HTTP redirect. The arguments req and via are
    // the upcoming request and the requests made already, oldest
    // first. If CheckRedirect returns an error, the Client's Get
    // method returns both the previous Response and
    // CheckRedirect's error (wrapped in a url.Error) instead of
    // issuing the Request req.
    //
    // If CheckRedirect is nil, the Client uses its default policy,
    // which is to stop after 10 consecutive requests.
    CheckRedirect func(req *Request, via []*Request) error

    // Jar specifies the cookie jar.
    // If Jar is nil, cookies are not sent in requests and ignored
    // in responses.
    Jar CookieJar

    // Timeout specifies a time limit for requests made by this
    // Client. The timeout includes connection time, any
    // redirects, and reading the response body. The timer remains
    // running after Get, Head, Post, or Do return and will
    // interrupt reading of the Response.Body.
    //
    // A Timeout of zero means no timeout.
    //
    // The Client's Transport must support the CancelRequest
    // method or Client will return errors when attempting to make
    // a request with Get, Head, Post, or Do. Client's default
    // Transport (DefaultTransport) supports CancelRequest.
    Timeout time.Duration
}

// The CloseNotifier interface is implemented by ResponseWriters which
// allow detecting when the underlying connection has gone away.
//
// This mechanism can be used to cancel long operations on the server
// if the client has disconnected before the response is ready.

// HTTP处理器ResponseWriter接口参数的下层如果实现了CloseNotifier接口，可以让用户
// 检测下层的连接是否停止。如果客户端在回复准备好之前关闭了连接，该机制可以用于
// 取消服务端耗时较长的操作。
type CloseNotifier interface {
    // CloseNotify returns a channel that receives a single value
    // when the client connection has gone away.
    CloseNotify() <-chan bool
}

// A ConnState represents the state of a client connection to a server.
// It's used by the optional Server.ConnState hook.

// ConnState代表一个客户端到服务端的连接的状态。本类型用于可选的Server.ConnState
// 回调函数。
//
//     const (
//         // StateNew代表一个新的连接，将要立刻发送请求。
//         // 连接从这个状态开始，然后转变为StateAlive或StateClosed。
//         StateNew ConnState = iota
//         // StateActive代表一个已经读取了请求数据1到多个字节的连接。
//         // 用于StateAlive的Server.ConnState回调函数在将连接交付给处理器之前被触发，
//         // 等到请求被处理完后，Server.ConnState回调函数再次被触发。
//         // 在请求被处理后，连接状态改变为StateClosed、StateHijacked或StateIdle。
//         StateActive
//         // StateIdle代表一个已经处理完了请求、处在闲置状态、等待新请求的连接。
//         // 连接状态可以从StateIdle改变为StateActive或StateClosed。
//         StateIdle
//         // 代表一个被劫持的连接。这是一个终止状态，不会转变为StateClosed。
//         StateHijacked
//         // StateClosed代表一个关闭的连接。
//         // 这是一个终止状态。被劫持的连接不会转变为StateClosed。
//         StateClosed
//     )
type ConnState int

// A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an
// HTTP response or the Cookie header of an HTTP request.
//
// See http://tools.ietf.org/html/rfc6265 for details.

// Cookie代表一个出现在HTTP回复的头域中Set-Cookie头的值里或者HTTP请求的头域中
// Cookie头的值里的HTTP cookie。
type Cookie struct {
    Name       string
    Value      string
    Path       string
    Domain     string
    Expires    time.Time
    RawExpires string

    // MaxAge=0 means no 'Max-Age' attribute specified.
    // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
    // MaxAge>0 means Max-Age attribute present and given in seconds
    MaxAge   int
    Secure   bool
    HttpOnly bool
    Raw      string
    Unparsed []string // Raw text of unparsed attribute-value pairs
}

// A CookieJar manages storage and use of cookies in HTTP requests.
//
// Implementations of CookieJar must be safe for concurrent use by multiple
// goroutines.
//
// The net/http/cookiejar package provides a CookieJar implementation.

// CookieJar管理cookie的存储和在HTTP请求中的使用。CookieJar的实现必须能安全的被
// 多个go程同时使用。
//
// net/http/cookiejar包提供了一个CookieJar的实现。
type CookieJar interface {
    // SetCookies handles the receipt of the cookies in a reply for the
    // given URL.  It may or may not choose to save the cookies, depending
    // on the jar's policy and implementation.
    SetCookies(u *url.URL, cookies []*Cookie)

    // Cookies returns the cookies to send in a request for the given URL.
    // It is up to the implementation to honor the standard cookie use
    // restrictions such as in RFC 6265.
    Cookies(u *url.URL) []*Cookie
}

// A Dir implements FileSystem using the native file system restricted to a
// specific directory tree.
//
// While the FileSystem.Open method takes '/'-separated paths, a Dir's string
// value is a filename on the native file system, not a URL, so it is separated
// by filepath.Separator, which isn't necessarily '/'.
//
// An empty Dir is treated as ".".

// Dir使用限制到指定目录树的本地文件系统实现了http.FileSystem接口。空Dir被视为
// "."，即代表当前目录。
type Dir string

// A File is returned by a FileSystem's Open method and can be
// served by the FileServer implementation.
//
// The methods should behave the same as those on an *os.File.

// File是被FileSystem接口的Open方法返回的接口类型，可以被FileServer等函数用于文
// 件访问服务。
//
// 该接口的方法的行为应该和*os.File类型的同名方法相同。
type File interface {
    io.Closer
    io.Reader
    Readdir(count int) ([]os.FileInfo, error)
    Seek(offset int64, whence int) (int64, error)
    Stat() (os.FileInfo, error)
}

// A FileSystem implements access to a collection of named files.
// The elements in a file path are separated by slash ('/', U+002F)
// characters, regardless of host operating system convention.

// FileSystem接口实现了对一系列命名文件的访问。文件路径的分隔符为'/'，不管主机操
// 作系统的惯例如何。
type FileSystem interface {
    Open(name string) (File, error)
}

// The Flusher interface is implemented by ResponseWriters that allow
// an HTTP handler to flush buffered data to the client.
//
// Note that even for ResponseWriters that support Flush,
// if the client is connected through an HTTP proxy,
// the buffered data may not reach the client until the response
// completes.

// HTTP处理器ResponseWriter接口参数的下层如果实现了Flusher接口，可以让HTTP处理器
// 将缓冲中的数据发送到客户端。
//
// 注意：即使ResponseWriter接口的下层支持Flush方法，如果客户端是通过HTTP代理连接
// 的，缓冲中的数据也可能直到回复完毕才被传输到客户端。
type Flusher interface {
    // Flush sends any buffered data to the client.
    Flush()
}

// A Handler responds to an HTTP request.
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
//
// Depending on the HTTP client software, HTTP protocol version, and
// any intermediaries between the client and the Go server, it may not
// be possible to read from the Request.Body after writing to the
// ResponseWriter. Cautious handlers should read the Request.Body
// first, and then reply.
//
// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes
// that the effect of the panic was isolated to the active request.
// It recovers the panic, logs a stack trace to the server error log,
// and hangs up the connection.

// 实现了Handler接口的对象可以注册到HTTP服务端，为特定的路径及其子树提供服务。
//
// ServeHTTP应该将回复的头域和数据写入ResponseWriter接口然后返回。返回标志着该请
// 求已经结束，HTTP服务端可以转移向该连接上的下一个请求。
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.

// HandlerFunc type是一个适配器，通过类型转换让我们可以将普通的函数作为HTTP处理
// 器使用。如果f是一个具有适当签名的函数，HandlerFunc(f)通过调用f实现了Handler接
// 口。
type HandlerFunc func(ResponseWriter, *Request)

// A Header represents the key-value pairs in an HTTP header.

// Header代表HTTP头域的键值对。
type Header map[string][]string

// The Hijacker interface is implemented by ResponseWriters that allow
// an HTTP handler to take over the connection.

// HTTP处理器ResponseWriter接口参数的下层如果实现了Hijacker接口，可以让HTTP处理
// 器接管该连接。
type Hijacker interface {
    // Hijack lets the caller take over the connection.
    // After a call to Hijack(), the HTTP server library
    // will not do anything else with the connection.
    // It becomes the caller's responsibility to manage
    // and close the connection.
    Hijack() (net.Conn, *bufio.ReadWriter, error)
}

// HTTP request parsing errors.

// HTTP请求解析错误。
type ProtocolError struct {
    ErrorString string
}

// A Request represents an HTTP request received by a server
// or to be sent by a client.
//
// The field semantics differ slightly between client and server
// usage. In addition to the notes on the fields below, see the
// documentation for Request.Write and RoundTripper.

// Request类型代表一个服务端接受到的或者客户端发送出去的HTTP请求。
//
// Request各字段的意义和用途在服务端和客户端是不同的。除了字段本身上方文档，还可
// 参见Request.Write方法和RoundTripper接口的文档。
type Request struct {
    // Method specifies the HTTP method (GET, POST, PUT, etc.).
    // For client requests an empty string means GET.
    Method string

    // URL specifies either the URI being requested (for server
    // requests) or the URL to access (for client requests).
    //
    // For server requests the URL is parsed from the URI
    // supplied on the Request-Line as stored in RequestURI.  For
    // most requests, fields other than Path and RawQuery will be
    // empty. (See RFC 2616, Section 5.1.2)
    //
    // For client requests, the URL's Host specifies the server to
    // connect to, while the Request's Host field optionally
    // specifies the Host header value to send in the HTTP
    // request.
    URL *url.URL

    // The protocol version for incoming requests.
    // Client requests always use HTTP/1.1.
    Proto      string // "HTTP/1.0"
    ProtoMajor int    // 1
    ProtoMinor int    // 0

    // A header maps request lines to their values.
    // If the header says
    //
    //	accept-encoding: gzip, deflate
    //	Accept-Language: en-us
    //	Connection: keep-alive
    //
    // then
    //
    //	Header = map[string][]string{
    //		"Accept-Encoding": {"gzip, deflate"},
    //		"Accept-Language": {"en-us"},
    //		"Connection": {"keep-alive"},
    //	}
    //
    // HTTP defines that header names are case-insensitive.
    // The request parser implements this by canonicalizing the
    // name, making the first character and any characters
    // following a hyphen uppercase and the rest lowercase.
    //
    // For client requests certain headers are automatically
    // added and may override values in Header.
    //
    // See the documentation for the Request.Write method.
    Header Header

    // Body is the request's body.
    //
    // For client requests a nil body means the request has no
    // body, such as a GET request. The HTTP Client's Transport
    // is responsible for calling the Close method.
    //
    // For server requests the Request Body is always non-nil
    // but will return EOF immediately when no body is present.
    // The Server will close the request body. The ServeHTTP
    // Handler does not need to.
    Body io.ReadCloser

    // ContentLength records the length of the associated content.
    // The value -1 indicates that the length is unknown.
    // Values >= 0 indicate that the given number of bytes may
    // be read from Body.
    // For client requests, a value of 0 means unknown if Body is not nil.
    ContentLength int64

    // TransferEncoding lists the transfer encodings from outermost to
    // innermost. An empty list denotes the "identity" encoding.
    // TransferEncoding can usually be ignored; chunked encoding is
    // automatically added and removed as necessary when sending and
    // receiving requests.
    TransferEncoding []string

    // Close indicates whether to close the connection after
    // replying to this request (for servers) or after sending
    // the request (for clients).
    Close bool

    // For server requests Host specifies the host on which the
    // URL is sought. Per RFC 2616, this is either the value of
    // the "Host" header or the host name given in the URL itself.
    // It may be of the form "host:port".
    //
    // For client requests Host optionally overrides the Host
    // header to send. If empty, the Request.Write method uses
    // the value of URL.Host.
    Host string

    // Form contains the parsed form data, including both the URL
    // field's query parameters and the POST or PUT form data.
    // This field is only available after ParseForm is called.
    // The HTTP client ignores Form and uses Body instead.
    Form url.Values

    // PostForm contains the parsed form data from POST or PUT
    // body parameters.
    // This field is only available after ParseForm is called.
    // The HTTP client ignores PostForm and uses Body instead.
    PostForm url.Values

    // MultipartForm is the parsed multipart form, including file uploads.
    // This field is only available after ParseMultipartForm is called.
    // The HTTP client ignores MultipartForm and uses Body instead.
    MultipartForm *multipart.Form

    // Trailer specifies additional headers that are sent after the request
    // body.
    //
    // For server requests the Trailer map initially contains only the
    // trailer keys, with nil values. (The client declares which trailers it
    // will later send.)  While the handler is reading from Body, it must
    // not reference Trailer. After reading from Body returns EOF, Trailer
    // can be read again and will contain non-nil values, if they were sent
    // by the client.
    //
    // For client requests Trailer must be initialized to a map containing
    // the trailer keys to later send. The values may be nil or their final
    // values. The ContentLength must be 0 or -1, to send a chunked request.
    // After the HTTP request is sent the map values can be updated while
    // the request body is read. Once the body returns EOF, the caller must
    // not mutate Trailer.
    //
    // Few HTTP clients, servers, or proxies support HTTP trailers.
    Trailer Header

    // RemoteAddr allows HTTP servers and other software to record
    // the network address that sent the request, usually for
    // logging. This field is not filled in by ReadRequest and
    // has no defined format. The HTTP server in this package
    // sets RemoteAddr to an "IP:port" address before invoking a
    // handler.
    // This field is ignored by the HTTP client.
    RemoteAddr string

    // RequestURI is the unmodified Request-URI of the
    // Request-Line (RFC 2616, Section 5.1) as sent by the client
    // to a server. Usually the URL field should be used instead.
    // It is an error to set this field in an HTTP client request.
    RequestURI string

    // TLS allows HTTP servers and other software to record
    // information about the TLS connection on which the request
    // was received. This field is not filled in by ReadRequest.
    // The HTTP server in this package sets the field for
    // TLS-enabled connections before invoking a handler;
    // otherwise it leaves the field nil.
    // This field is ignored by the HTTP client.
    TLS *tls.ConnectionState
}

// Response represents the response from an HTTP request.

// Response代表一个HTTP请求的回复。
type Response struct {
    Status     string // e.g. "200 OK"
    StatusCode int    // e.g. 200
    Proto      string // e.g. "HTTP/1.0"
    ProtoMajor int    // e.g. 1
    ProtoMinor int    // e.g. 0

    // Header maps header keys to values.  If the response had multiple
    // headers with the same key, they may be concatenated, with comma
    // delimiters.  (Section 4.2 of RFC 2616 requires that multiple headers
    // be semantically equivalent to a comma-delimited sequence.) Values
    // duplicated by other fields in this struct (e.g., ContentLength) are
    // omitted from Header.
    //
    // Keys in the map are canonicalized (see CanonicalHeaderKey).
    Header Header

    // Body represents the response body.
    //
    // The http Client and Transport guarantee that Body is always
    // non-nil, even on responses without a body or responses with
    // a zero-length body. It is the caller's responsibility to
    // close Body.
    //
    // The Body is automatically dechunked if the server replied
    // with a "chunked" Transfer-Encoding.
    Body io.ReadCloser

    // ContentLength records the length of the associated content.  The
    // value -1 indicates that the length is unknown.  Unless Request.Method
    // is "HEAD", values >= 0 indicate that the given number of bytes may
    // be read from Body.
    ContentLength int64

    // Contains transfer encodings from outer-most to inner-most. Value is
    // nil, means that "identity" encoding is used.
    TransferEncoding []string

    // Close records whether the header directed that the connection be
    // closed after reading Body.  The value is advice for clients: neither
    // ReadResponse nor Response.Write ever closes a connection.
    Close bool

    // Trailer maps trailer keys to values, in the same
    // format as the header.
    Trailer Header

    // The Request that was sent to obtain this Response.
    // Request's Body is nil (having already been consumed).
    // This is only populated for Client requests.
    Request *Request

    // TLS contains information about the TLS connection on which the
    // response was received. It is nil for unencrypted responses.
    // The pointer is shared between responses and should not be
    // modified.
    TLS *tls.ConnectionState
}

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
//
// A ResponseWriter may not be used after the Handler.ServeHTTP method
// has returned.

// ResponseWriter接口被HTTP处理器用于构造HTTP回复。
type ResponseWriter interface {
    // Header returns the header map that will be sent by WriteHeader.
    // Changing the header after a call to WriteHeader (or Write) has
    // no effect.
    Header() Header

    // Write writes the data to the connection as part of an HTTP reply.
    // If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
    // before writing the data.  If the Header does not contain a
    // Content-Type line, Write adds a Content-Type set to the result of passing
    // the initial 512 bytes of written data to DetectContentType.
    Write([]byte) (int, error)

    // WriteHeader sends an HTTP response header with status code.
    // If WriteHeader is not called explicitly, the first call to Write
    // will trigger an implicit WriteHeader(http.StatusOK).
    // Thus explicit calls to WriteHeader are mainly used to
    // send error codes.
    WriteHeader(int)
}

// RoundTripper is an interface representing the ability to execute a
// single HTTP transaction, obtaining the Response for a given Request.
//
// A RoundTripper must be safe for concurrent use by multiple
// goroutines.

// RoundTripper接口是具有执行单次HTTP事务的能力（接收指定请求的回复）的接口。
//
// RoundTripper接口的类型必须可以安全的被多线程同时使用。
type RoundTripper interface {
    // RoundTrip executes a single HTTP transaction, returning
    // the Response for the request req.  RoundTrip should not
    // attempt to interpret the response.  In particular,
    // RoundTrip must return err == nil if it obtained a response,
    // regardless of the response's HTTP status code.  A non-nil
    // err should be reserved for failure to obtain a response.
    // Similarly, RoundTrip should not attempt to handle
    // higher-level protocol details such as redirects,
    // authentication, or cookies.
    //
    // RoundTrip should not modify the request, except for
    // consuming and closing the Body, including on errors. The
    // request's URL and Header fields are guaranteed to be
    // initialized.
    RoundTrip(*Request) (*Response, error)
}

// ServeMux is an HTTP request multiplexer.
// It matches the URL of each incoming request against a list of registered
// patterns and calls the handler for the pattern that
// most closely matches the URL.
//
// Patterns name fixed, rooted paths, like "/favicon.ico",
// or rooted subtrees, like "/images/" (note the trailing slash).
// Longer patterns take precedence over shorter ones, so that
// if there are handlers registered for both "/images/"
// and "/images/thumbnails/", the latter handler will be
// called for paths beginning "/images/thumbnails/" and the
// former will receive requests for any other paths in the
// "/images/" subtree.
//
// Note that since a pattern ending in a slash names a rooted subtree,
// the pattern "/" matches all paths not matched by other registered
// patterns, not just the URL with Path == "/".
//
// If a subtree has been registered and a request is received naming the
// subtree root without its trailing slash, ServeMux redirects that
// request to the subtree root (adding the trailing slash). This behavior can
// be overridden with a separate registration for the path without
// the trailing slash. For example, registering "/images/" causes ServeMux
// to redirect a request for "/images" to "/images/", unless "/images" has
// been registered separately.
//
// Patterns may optionally begin with a host name, restricting matches to
// URLs on that host only.  Host-specific patterns take precedence over
// general patterns, so that a handler might register for the two patterns
// "/codesearch" and "codesearch.google.com/" without also taking over
// requests for "http://www.google.com/".
//
// ServeMux also takes care of sanitizing the URL request path,
// redirecting any request containing . or .. elements or repeated slashes
// to an equivalent, cleaner URL.

// ServeMux类型是HTTP请求的多路转接器。它会将每一个接收的请求的URL与一个注册模式
// 的列表进行匹配，并调用和URL最匹配的模式的处理器。
//
// 模式是固定的、由根开始的路径，如"/favicon.ico"，或由根开始的子树，如
// "/images/"（注意结尾的斜杠）。较长的模式优先于较短的模式，因此如果模式
// "/images/"和"/images/thumbnails/"都注册了处理器，后一个处理器会用于路径以
// "/images/thumbnails/"开始的请求，前一个处理器会接收到其余的路径在"/images/"子
// 树下的请求。
//
// 注意，因为以斜杠结尾的模式代表一个由根开始的子树，模式"/"会匹配所有的未被其他
// 注册的模式匹配的路径，而不仅仅是路径"/"。
//
// 模式也能（可选地）以主机名开始，表示只匹配该主机上的路径。指定主机的模式优先
// 于一般的模式，因此一个注册了两个模式"/codesearch"和"codesearch.google.com/"的
// 处理器不会接管目标为"http://www.google.com/"的请求。
//
// ServeMux还会注意到请求的URL路径的无害化，将任何路径中包含"."或".."元素的请求
// 重定向到等价的没有这两种元素的URL。（参见path.Clean函数）
type ServeMux struct {
}

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.

// Server类型定义了运行HTTP服务端的参数。Server的零值是合法的配置。
type Server struct {
    Addr           string        // TCP address to listen on, ":http" if empty
    Handler        Handler       // handler to invoke, http.DefaultServeMux if nil
    ReadTimeout    time.Duration // maximum duration before timing out read of the request
    WriteTimeout   time.Duration // maximum duration before timing out write of the response
    MaxHeaderBytes int           // maximum size of request headers, DefaultMaxHeaderBytes if 0
    TLSConfig      *tls.Config   // optional TLS config, used by ListenAndServeTLS

    // TLSNextProto optionally specifies a function to take over
    // ownership of the provided TLS connection when an NPN
    // protocol upgrade has occurred.  The map key is the protocol
    // name negotiated. The Handler argument should be used to
    // handle HTTP requests and will initialize the Request's TLS
    // and RemoteAddr if not already set.  The connection is
    // automatically closed when the function returns.
    TLSNextProto map[string]func(*Server, *tls.Conn, Handler)

    // ConnState specifies an optional callback function that is
    // called when a client connection changes state. See the
    // ConnState type and associated constants for details.
    ConnState func(net.Conn, ConnState)

    // ErrorLog specifies an optional logger for errors accepting
    // connections and unexpected behavior from handlers.
    // If nil, logging goes to os.Stderr via the log package's
    // standard logger.
    ErrorLog *log.Logger
}

// Transport is an implementation of RoundTripper that supports HTTP,
// HTTPS, and HTTP proxies (for either HTTP or HTTPS with CONNECT).
//
// By default, Transport caches connections for future re-use.
// This may leave many open connections when accessing many hosts.
// This behavior can be managed using Transport's CloseIdleConnections method
// and the MaxIdleConnsPerHost and DisableKeepAlives fields.
//
// Transports should be reused instead of created as needed.
// Transports are safe for concurrent use by multiple goroutines.
//
// A Transport is a low-level primitive for making HTTP and HTTPS requests.
// For high-level functionality, such as cookies and redirects, see Client.
//
// Transport uses HTTP/1.1 for HTTP URLs and either HTTP/1.1 or HTTP/2
// for HTTPS URLs, depending on whether the server supports HTTP/2.
// See the package docs for more about HTTP/2.

// Transport类型实现了RoundTripper接口，支持http、https和http/https代理。
// Transport类型可以缓存连接以在未来重用。
//
//     var DefaultTransport RoundTripper = &Transport{
//         Proxy: ProxyFromEnvironment,
//         Dial: (&net.Dialer{
//             Timeout:   30 * time.Second,
//             KeepAlive: 30 * time.Second,
//         }).Dial,
//         TLSHandshakeTimeout: 10 * time.Second,
//     }
//
// DefaultTransport是被包变量DefaultClient使用的默认RoundTripper接口。它会根据需
// 要创建网络连接，并缓存以便在之后的请求中重用这些连接。它使用环境变量
// $HTTP_PROXY和$NO_PROXY（或$http_proxy和$no_proxy）指定的HTTP代理。
type Transport struct {

    // Proxy specifies a function to return a proxy for a given
    // Request. If the function returns a non-nil error, the
    // request is aborted with the provided error.
    // If Proxy is nil or returns a nil *URL, no proxy is used.
    Proxy func(*Request) (*url.URL, error)

    // Dial specifies the dial function for creating unencrypted
    // TCP connections.
    // If Dial is nil, net.Dial is used.
    Dial func(network, addr string) (net.Conn, error)

    // DialTLS specifies an optional dial function for creating
    // TLS connections for non-proxied HTTPS requests.
    //
    // If DialTLS is nil, Dial and TLSClientConfig are used.
    //
    // If DialTLS is set, the Dial hook is not used for HTTPS
    // requests and the TLSClientConfig and TLSHandshakeTimeout
    // are ignored. The returned net.Conn is assumed to already be
    // past the TLS handshake.
    DialTLS func(network, addr string) (net.Conn, error)

    // TLSClientConfig specifies the TLS configuration to use with
    // tls.Client. If nil, the default configuration is used.
    TLSClientConfig *tls.Config

    // TLSHandshakeTimeout specifies the maximum amount of time waiting to
    // wait for a TLS handshake. Zero means no timeout.
    TLSHandshakeTimeout time.Duration

    // DisableKeepAlives, if true, prevents re-use of TCP connections
    // between different HTTP requests.
    DisableKeepAlives bool

    // DisableCompression, if true, prevents the Transport from
    // requesting compression with an "Accept-Encoding: gzip"
    // request header when the Request contains no existing
    // Accept-Encoding value. If the Transport requests gzip on
    // its own and gets a gzipped response, it's transparently
    // decoded in the Response.Body. However, if the user
    // explicitly requested gzip it is not automatically
    // uncompressed.
    DisableCompression bool

    // MaxIdleConnsPerHost, if non-zero, controls the maximum idle
    // (keep-alive) to keep per-host.  If zero,
    // DefaultMaxIdleConnsPerHost is used.
    MaxIdleConnsPerHost int

    // ResponseHeaderTimeout, if non-zero, specifies the amount of
    // time to wait for a server's response headers after fully
    // writing the request (including its body, if any). This
    // time does not include the time to read the response body.
    ResponseHeaderTimeout time.Duration
}

// CanonicalHeaderKey returns the canonical format of the
// header key s.  The canonicalization converts the first
// letter and any letter following a hyphen to upper case;
// the rest are converted to lowercase.  For example, the
// canonical key for "accept-encoding" is "Accept-Encoding".
// If s contains a space or invalid header field bytes, it is
// returned without modifications.

// CanonicalHeaderKey函数返回头域（表示为Header类型）的键s的规范化格式。规范化过
// 程中让单词首字母和'-'后的第一个字母大写，其余字母小写。例如，
// "accept-encoding"规范化为"Accept-Encoding"。
func CanonicalHeaderKey(s string) string

// DetectContentType implements the algorithm described
// at http://mimesniff.spec.whatwg.org/ to determine the
// Content-Type of the given data.  It considers at most the
// first 512 bytes of data.  DetectContentType always returns
// a valid MIME type: if it cannot determine a more specific one, it
// returns "application/octet-stream".

// DetectContentType函数实现了
// http://mimesniff.spec.whatwg.org/描述的算法，用于确定数据的Content-Type。函数
// 总是返回一个合法的MIME类型；如果它不能确定数据的类型，将返回
// "application/octet-stream"。它最多检查数据的前512字节。
func DetectContentType(data []byte) string

// Error replies to the request with the specified error message and HTTP code.
// The error message should be plain text.

// Error使用指定的错误信息和状态码回复请求，将数据写入w。错误信息必须是明文。
func Error(w ResponseWriter, error string, code int)

// FileServer returns a handler that serves HTTP requests
// with the contents of the file system rooted at root.
//
// To use the operating system's file system implementation,
// use http.Dir:
//
//     http.Handle("/", http.FileServer(http.Dir("/tmp")))
//
// As a special case, the returned file server redirects any request
// ending in "/index.html" to the same path, without the final
// "index.html".

// FileServer返回一个使用FileSystem接口root提供文件访问服务的HTTP处理器。要使用
// 操作系统的FileSystem接口实现，可使用http.Dir：
//
//     http.Handle("/", http.FileServer(http.Dir("/tmp")))
func FileServer(root FileSystem) Handler

// Get issues a GET to the specified URL. If the response is one of
// the following redirect codes, Get follows the redirect, up to a
// maximum of 10 redirects:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// An error is returned if there were too many redirects or if there
// was an HTTP protocol error. A non-2xx response doesn't cause an
// error.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// Get is a wrapper around DefaultClient.Get.
//
// To make a request with custom headers, use NewRequest and
// DefaultClient.Do.

// Get向指定的URL发出一个GET请求，如果回应的状态码如下，Get会在调用
// c.CheckRedirect后执行重定向：
//
//     301 (Moved Permanently)
//     302 (Found)
//     303 (See Other)
//     307 (Temporary Redirect)
//
// 如果c.CheckRedirect执行失败或存在HTTP协议错误时，本方法将返回该错误；如果回应
// 的状态码不是2xx，本方法并不会返回错误。如果返回值err为nil，resp.Body总是非nil
// 的，调用者应该在读取完resp.Body后关闭它。
//
// Get是对包变量DefaultClient的Get方法的包装。
func Get(url string) (resp *Response, err error)

// Handle registers the handler for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.

// Handle注册HTTP处理器handler和对应的模式pattern（注册到DefaultServeMux）。如果
// 该模式已经注册有一个处理器，Handle会panic。ServeMux的文档解释了模式的匹配机制
// 。
func Handle(pattern string, handler Handler)

// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.

// HandleFunc注册一个处理器函数handler和对应的模式pattern（注册到DefaultServeMux
// ）。ServeMux的文档解释了模式的匹配机制。
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))

// Head issues a HEAD to the specified URL.  If the response is one of
// the following redirect codes, Head follows the redirect, up to a
// maximum of 10 redirects:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// Head is a wrapper around DefaultClient.Head

// Head向指定的URL发出一个HEAD请求，如果回应的状态码如下，Head会在调用
// c.CheckRedirect后执行重定向：
//
//     301 (Moved Permanently)
//     302 (Found)
//     303 (See Other)
//     307 (Temporary Redirect)
//
// Head是对包变量DefaultClient的Head方法的包装。
func Head(url string) (resp *Response, err error)

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// Handler is typically nil, in which case the DefaultServeMux is
// used.
//
// A trivial example server is:
//
//     package main
//
//     import (
//         "io"
//         "net/http"
//         "log"
//     )
//
//     // hello world, the web server
//     func HelloServer(w http.ResponseWriter, req *http.Request) {
//         io.WriteString(w, "hello, world!\n")
//     }
//
//     func main() {
//         http.HandleFunc("/hello", HelloServer)
//         log.Fatal(http.ListenAndServe(":12345", nil))
//     }
//
// ListenAndServe always returns a non-nil error.

// ListenAndServe监听TCP地址addr，并且会使用handler参数调用Serve函数处理接收到的
// 连接。handler参数一般会设为nil，此时会使用DefaultServeMux。
//
// 一个简单的服务端例子：
//
//     package main
//     import (
//         "io"
//         "net/http"
//         "log"
//     )
//     // hello world, the web server
//     func HelloServer(w http.ResponseWriter, req *http.Request) {
//         io.WriteString(w, "hello, world!\n")
//     }
//     func main() {
//         http.HandleFunc("/hello", HelloServer)
//         err := http.ListenAndServe(":12345", nil)
//         if err != nil {
//             log.Fatal("ListenAndServe: ", err)
//         }
//     }
func ListenAndServe(addr string, handler Handler) error

// ListenAndServeTLS acts identically to ListenAndServe, except that it expects
// HTTPS connections. Additionally, files containing a certificate and matching
// private key for the server must be provided. If the certificate is signed by
// a certificate authority, the certFile should be the concatenation of the
// server's certificate, any intermediates, and the CA's certificate.
//
// A trivial example server is:
//
//     import (
//         "log"
//         "net/http"
//     )
//
//     func handler(w http.ResponseWriter, req *http.Request) {
//         w.Header().Set("Content-Type", "text/plain")
//         w.Write([]byte("This is an example server.\n"))
//     }
//
//     func main() {
//         http.HandleFunc("/", handler)
//         log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
//         err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
//         log.Fatal(err)
//     }
//
// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
//
// ListenAndServeTLS always returns a non-nil error.

// ListenAndServeTLS函数和ListenAndServe函数的行为基本一致，除了它期望HTTPS连接
// 之外。此外，必须提供证书文件和对应的私钥文件。如果证书是由权威机构签发的，
// certFile参数必须是顺序串联的服务端证书和CA证书。如果srv.Addr为空字符串，会使
// 用":https"。
//
// 一个简单的服务端例子：
//
//     import (
//         "log"
//         "net/http"
//     )
//     func handler(w http.ResponseWriter, req *http.Request) {
//         w.Header().Set("Content-Type", "text/plain")
//         w.Write([]byte("This is an example server.\n"))
//     }
//     func main() {
//         http.HandleFunc("/", handler)
//         log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
//         err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
//         if err != nil {
//             log.Fatal(err)
//         }
//     }
//
// 程序员可以使用crypto/tls包的generate_cert.go文件来生成cert.pem和key.pem两个文
// 件。
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler) error

// MaxBytesReader is similar to io.LimitReader but is intended for
// limiting the size of incoming request bodies. In contrast to
// io.LimitReader, MaxBytesReader's result is a ReadCloser, returns a
// non-EOF error for a Read beyond the limit, and closes the
// underlying reader when its Close method is called.
//
// MaxBytesReader prevents clients from accidentally or maliciously
// sending a large request and wasting server resources.

// MaxBytesReader类似io.LimitReader，但它是用来限制接收到的请求的Body的大小的。
// 不同于io.LimitReader，本函数返回一个ReadCloser，返回值的Read方法在读取的数据
// 超过大小限制时会返回非EOF错误，其Close方法会关闭下层的io.ReadCloser接口r。
//
// MaxBytesReader预防客户端因为意外或者蓄意发送的“大”请求，以避免尺寸过大的请
// 求浪费服务端资源。
func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser

// NewFileTransport returns a new RoundTripper, serving the provided
// FileSystem. The returned RoundTripper ignores the URL host in its
// incoming requests, as well as most other properties of the
// request.
//
// The typical use case for NewFileTransport is to register the "file"
// protocol with a Transport, as in:
//
//   t := &http.Transport{}
//   t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
//   c := &http.Client{Transport: t}
//   res, err := c.Get("file:///etc/passwd")
//   ...

// NewFileTransport返回一个RoundTripper接口，使用FileSystem接口fs提供文件访问服
// 务。 返回的RoundTripper接口会忽略接收的请求的URL主机及其他绝大多数属性。
//
// NewFileTransport函数的典型使用情况是给Transport类型的值注册"file"协议，如下所
// 示：
//
//     t := &http.Transport{}
//     t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
//     c := &http.Client{Transport: t}
//     res, err := c.Get("file:///etc/passwd")
//     ...
func NewFileTransport(fs FileSystem) RoundTripper

// NewRequest returns a new Request given a method, URL, and optional body.
//
// If the provided body is also an io.Closer, the returned
// Request.Body is set to body and will be closed by the Client
// methods Do, Post, and PostForm, and Transport.RoundTrip.
//
// NewRequest returns a Request suitable for use with Client.Do or
// Transport.RoundTrip.
// To create a request for use with testing a Server Handler use either
// ReadRequest or manually update the Request fields. See the Request
// type's documentation for the difference between inbound and outbound
// request fields.

// NewRequest使用指定的方法、网址和可选的主题创建并返回一个新的*Request。
//
// 如果body参数实现了io.Closer接口，Request返回值的Body 字段会被设置为body，并会
// 被Client类型的Do、Post和PostFOrm方法以及Transport.RoundTrip方法关闭。
func NewRequest(method, urlStr string, body io.Reader) (*Request, error)

// NewServeMux allocates and returns a new ServeMux.

// NewServeMux创建并返回一个新的*ServeMux
func NewServeMux() *ServeMux

// NotFound replies to the request with an HTTP 404 not found error.

// NotFound回复请求404状态码（not found：目标未发现）。
func NotFound(w ResponseWriter, r *Request)

// NotFoundHandler returns a simple request handler
// that replies to each request with a ``404 page not found'' reply.

// NotFoundHandler返回一个简单的请求处理器，该处理器会对每个请求都回复"404 page
// not found"。
func NotFoundHandler() Handler

// ParseHTTPVersion parses a HTTP version string.
// "HTTP/1.0" returns (1, 0, true).

// ParseHTTPVersion解析HTTP版本字符串。如"HTTP/1.0"返回(1, 0, true)。
func ParseHTTPVersion(vers string) (major, minor int, ok bool)

// ParseTime parses a time header (such as the Date: header),
// trying each of the three formats allowed by HTTP/1.1:
// TimeFormat, time.RFC850, and time.ANSIC.

// ParseTime用3种格式TimeFormat,
// time.RFC850和time.ANSIC尝试解析一个时间头的值（如Date: header）。
func ParseTime(text string) (t time.Time, err error)

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// If the provided body is an io.Closer, it is closed after the
// request.
//
// Post is a wrapper around DefaultClient.Post.
//
// To set custom headers, use NewRequest and DefaultClient.Do.

// Post向指定的URL发出一个POST请求。bodyType为POST数据的类型， body为POST数据，
// 作为请求的主体。如果参数body实现了io.Closer接口，它会在发送请求后被关闭。调用
// 者有责任在读取完返回值resp的主体后关闭它。
//
// Post是对包变量DefaultClient的Post方法的包装。
func Post(url string, bodyType string, body io.Reader) (resp *Response, err error)

// PostForm issues a POST to the specified URL, with data's keys and
// values URL-encoded as the request body.
//
// The Content-Type header is set to application/x-www-form-urlencoded.
// To set other headers, use NewRequest and DefaultClient.Do.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// PostForm is a wrapper around DefaultClient.PostForm.

// PostForm向指定的URL发出一个POST请求，url.Values类型的data会被编码为请求的主体
// 。如果返回值err为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭
// 它。
//
// PostForm是对包变量DefaultClient的PostForm方法的包装。
func PostForm(url string, data url.Values) (resp *Response, err error)

// ProxyFromEnvironment returns the URL of the proxy to use for a
// given request, as indicated by the environment variables
// HTTP_PROXY, HTTPS_PROXY and NO_PROXY (or the lowercase versions
// thereof). HTTPS_PROXY takes precedence over HTTP_PROXY for https
// requests.
//
// The environment values may be either a complete URL or a
// "host[:port]", in which case the "http" scheme is assumed.
// An error is returned if the value is a different form.
//
// A nil URL and nil error are returned if no proxy is defined in the
// environment, or a proxy should not be used for the given request,
// as defined by NO_PROXY.
//
// As a special case, if req.URL.Host is "localhost" (with or without
// a port number), then a nil URL and nil error will be returned.

// ProxyFromEnvironment使用环境变量$HTTP_PROXY和$NO_PROXY(或$http_proxy和
// $no_proxy)的配置返回用于req的代理。如果代理环境不合法将返回错误；如果环境未设
// 定代理或者给定的request不应使用代理时，将返回(nil, nil)；如果req.URL.Host字段
// 是"localhost"（可以有端口号，也可以没有），也会返回(nil, nil)。
func ProxyFromEnvironment(req *Request) (*url.URL, error)

// ProxyURL returns a proxy function (for use in a Transport)
// that always returns the same URL.

// ProxyURL返回一个代理函数（用于Transport类型），该函数总是返回同一个URL。
func ProxyURL(fixedURL *url.URL) (func(*Request) (*url.URL, error))

// ReadRequest reads and parses an incoming request from b.

// ReadRequest从b读取并解析出一个HTTP请求。（本函数主要用在服务端从下层获取请求
// ）
func ReadRequest(b *bufio.Reader) (req *Request, err error)

// ReadResponse reads and returns an HTTP response from r.
// The req parameter optionally specifies the Request that corresponds
// to this Response. If nil, a GET request is assumed.
// Clients must call resp.Body.Close when finished reading resp.Body.
// After that call, clients can inspect resp.Trailer to find key/value
// pairs included in the response trailer.

// ReadResponse从r读取并返回一个HTTP 回复。req参数是可选的，指定该回复对应的请求
// （即是对该请求的回复）。如果是nil，将假设请求是GET请求。客户端必须在结束
// resp.Body的读取后关闭它。读取完毕并关闭后，客户端可以检查resp.Trailer字段获取
// 回复的trailer的键值对。（本函数主要用在客户端从下层获取回复）
func ReadResponse(r *bufio.Reader, req *Request) (*Response, error)

// Redirect replies to the request with a redirect to url,
// which may be a path relative to the request path.
//
// The provided code should be in the 3xx range and is usually
// StatusMovedPermanently, StatusFound or StatusSeeOther.

// Redirect回复请求一个重定向地址urlStr和状态码code。该重定向地址可以是相对于请
// 求r的相对地址。
func Redirect(w ResponseWriter, r *Request, urlStr string, code int)

// RedirectHandler returns a request handler that redirects
// each request it receives to the given url using the given
// status code.
//
// The provided code should be in the 3xx range and is usually
// StatusMovedPermanently, StatusFound or StatusSeeOther.

// RedirectHandler返回一个请求处理器，该处理器会对每个请求都使用状态码code重定向
// 到网址url。
func RedirectHandler(url string, code int) Handler

// Serve accepts incoming HTTP connections on the listener l,
// creating a new service goroutine for each.  The service goroutines
// read requests and then call handler to reply to them.
// Handler is typically nil, in which case the DefaultServeMux is used.

// Serve会接手监听器l收到的每一个连接，并为每一个连接创建一个新的服务go程。该go
// 程会读取请求，然后调用handler回复请求。handler参数一般会设为nil，此时会使用
// DefaultServeMux。
func Serve(l net.Listener, handler Handler) error

// ServeContent replies to the request using the content in the
// provided ReadSeeker.  The main benefit of ServeContent over io.Copy
// is that it handles Range requests properly, sets the MIME type, and
// handles If-Modified-Since requests.
//
// If the response's Content-Type header is not set, ServeContent
// first tries to deduce the type from name's file extension and,
// if that fails, falls back to reading the first block of the content
// and passing it to DetectContentType.
// The name is otherwise unused; in particular it can be empty and is
// never sent in the response.
//
// If modtime is not the zero time or Unix epoch, ServeContent
// includes it in a Last-Modified header in the response.  If the
// request includes an If-Modified-Since header, ServeContent uses
// modtime to decide whether the content needs to be sent at all.
//
// The content's Seek method must work: ServeContent uses
// a seek to the end of the content to determine its size.
//
// If the caller has set w's ETag header, ServeContent uses it to
// handle requests using If-Range and If-None-Match.
//
// Note that *os.File implements the io.ReadSeeker interface.

// ServeContent使用提供的ReadSeeker的内容回复请求。ServeContent比起io.Copy函数的
// 主要优点，是可以处理范围类请求（只要一部分内容）、设置MIME类型，处理
// If-Modified-Since请求。
//
// 如果未设定回复的Content-Type头，本函数首先会尝试从name的文件扩展名推断数据类
// 型；如果失败，会用读取content的第1块数据并提供给DetectContentType推断类型；之
// 后会设置Content-Type头。参数name不会用于别的地方，甚至于它可以是空字符串，也
// 永远不会发送到回复里。
//
// 如果modtime不是Time零值，函数会在回复的头域里设置Last-Modified头。如果请求的
// 头域包含If-Modified-Since头，本函数会使用modtime参数来确定是否应该发送内容。
// 如果调用者设置了w的ETag头，ServeContent会使用它处理包含If-Range头和
// If-None-Match头的请求。
//
// 参数content的Seek方法必须有效：函数使用Seek来确定它的大小。
//
// 注意：本包File接口和*os.File类型都实现了io.ReadSeeker接口。
func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker)

// ServeFile replies to the request with the contents of the named
// file or directory.
//
// If the provided file or direcory name is a relative path, it is
// interpreted relative to the current directory and may ascend to parent
// directories. If the provided name is constructed from user input, it
// should be sanitized before calling ServeFile. As a precaution, ServeFile
// will reject requests where r.URL.Path contains a ".." path element.
//
// As a special case, ServeFile redirects any request where r.URL.Path
// ends in "/index.html" to the same path, without the final
// "index.html". To avoid such redirects either modify the path or
// use ServeContent.

// ServeFile回复请求name指定的文件或者目录的内容。
func ServeFile(w ResponseWriter, r *Request, name string)

// SetCookie adds a Set-Cookie header to the provided ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.

// SetCookie在w的头域中添加Set-Cookie头，该HTTP头的值为cookie。
func SetCookie(w ResponseWriter, cookie *Cookie)

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.

// StatusText返回HTTP状态码code对应的文本，如220对应"OK"。如果code是未知的状态码
// ，会返回""。
func StatusText(code int) string

// StripPrefix returns a handler that serves HTTP requests
// by removing the given prefix from the request URL's Path
// and invoking the handler h. StripPrefix handles a
// request for a path that doesn't begin with prefix by
// replying with an HTTP 404 not found error.

// StripPrefix返回一个处理器，该处理器会将请求的URL.Path字段中给定前缀prefix去除
// 后再交由h处理。StripPrefix会向URL.Path字段中没有给定前缀的请求回复404 page
// not found。
func StripPrefix(prefix string, h Handler) Handler

// TimeoutHandler returns a Handler that runs h with the given time limit.
//
// The new Handler calls h.ServeHTTP to handle each request, but if a
// call runs for longer than its time limit, the handler responds with
// a 503 Service Unavailable error and the given message in its body.
// (If msg is empty, a suitable default message will be sent.)
// After such a timeout, writes by h to its ResponseWriter will return
// ErrHandlerTimeout.
//
// TimeoutHandler buffers all Handler writes to memory and does not
// support the Hijacker or Flusher interfaces.

// TimeoutHandler返回一个采用指定时间限制的请求处理器。
//
// 返回的Handler会调用h.ServeHTTP去处理每个请求，但如果某一次调用耗时超过了时间
// 限制，该处理器会回复请求状态码503 Service Unavailable，并将msg作为回复的主体
// （如果msg为空字符串，将发送一个合理的默认信息）。在超时后，h对它的
// ResponseWriter接口参数的写入操作会返回ErrHandlerTimeout。
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler

// Do sends an HTTP request and returns an HTTP response, following
// policy (e.g. redirects, cookies, auth) as configured on the client.
//
// An error is returned if caused by client policy (such as
// CheckRedirect), or if there was an HTTP protocol error.
// A non-2xx response doesn't cause an error.
//
// When err is nil, resp always contains a non-nil resp.Body.
//
// Callers should close resp.Body when done reading from it. If
// resp.Body is not closed, the Client's underlying RoundTripper
// (typically Transport) may not be able to re-use a persistent TCP
// connection to the server for a subsequent "keep-alive" request.
//
// The request Body, if non-nil, will be closed by the underlying
// Transport, even on errors.
//
// Generally Get, Post, or PostForm will be used instead of Do.

// Do方法发送请求，返回HTTP回复。它会遵守客户端c设置的策略（如重定向、cookie、认
// 证）。
//
// 如果客户端的策略（如重定向）返回错误或存在HTTP协议错误时，本方法将返回该错误
// ；如果回应的状态码不是2xx，本方法并不会返回错误。
//
// 如果返回值err为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭它
// 。如果返回值resp的主体未关闭，c下层的RoundTripper接口（一般为Transport类型）
// 可能无法重用resp主体下层保持的TCP连接去执行之后的请求。
//
// 请求的主体，如果非nil，会在执行后被c.Transport关闭，即使出现错误。
//
// 一般应使用Get、Post或PostForm方法代替Do方法。
func (*Client) Do(req *Request) (resp *Response, err error)

// Get issues a GET to the specified URL. If the response is one of the
// following redirect codes, Get follows the redirect after calling the
// Client's CheckRedirect function:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// An error is returned if the Client's CheckRedirect function fails
// or if there was an HTTP protocol error. A non-2xx response doesn't
// cause an error.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// To make a request with custom headers, use NewRequest and Client.Do.

// Get向指定的URL发出一个GET请求，如果回应的状态码如下，Get会在调用
// c.CheckRedirect后执行重定向：
//
//     301 (Moved Permanently)
//     302 (Found)
//     303 (See Other)
//     307 (Temporary Redirect)
//
// 如果c.CheckRedirect执行失败或存在HTTP协议错误时，本方法将返回该错误；如果回应
// 的状态码不是2xx，本方法并不会返回错误。如果返回值err为nil，resp.Body总是非nil
// 的，调用者应该在读取完resp.Body后关闭它。
func (*Client) Get(url string) (resp *Response, err error)

// Head issues a HEAD to the specified URL.  If the response is one of the
// following redirect codes, Head follows the redirect after calling the
// Client's CheckRedirect function:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)

// Head向指定的URL发出一个HEAD请求，如果回应的状态码如下，Head会在调用
// c.CheckRedirect后执行重定向：
//
//     301 (Moved Permanently)
//     302 (Found)
//     303 (See Other)
//     307 (Temporary Redirect)
func (*Client) Head(url string) (resp *Response, err error)

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// If the provided body is an io.Closer, it is closed after the
// request.
//
// To set custom headers, use NewRequest and Client.Do.

// Post向指定的URL发出一个POST请求。bodyType为POST数据的类型， body为POST数据，
// 作为请求的主体。如果参数body实现了io.Closer接口，它会在发送请求后被关闭。调用
// 者有责任在读取完返回值resp的主体后关闭它。
func (*Client) Post(url string, bodyType string, body io.Reader) (resp *Response, err error)

// PostForm issues a POST to the specified URL,
// with data's keys and values URL-encoded as the request body.
//
// The Content-Type header is set to application/x-www-form-urlencoded.
// To set other headers, use NewRequest and DefaultClient.Do.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.

// PostForm向指定的URL发出一个POST请求，url.Values类型的data会被编码为请求的主体
// 。POST数据的类型一般会设为"application/x-www-form-urlencoded"。如果返回值err
// 为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭它。
func (*Client) PostForm(url string, data url.Values) (resp *Response, err error)

// String returns the serialization of the cookie for use in a Cookie
// header (if only Name and Value are set) or a Set-Cookie response
// header (if other fields are set).
// If c is nil or c.Name is invalid, the empty string is returned.

// String返回该cookie的序列化结果。如果只设置了Name和Value字段，序列化结果可用于
// HTTP请求的Cookie头或者HTTP回复的Set-Cookie头；如果设置了其他字段，序列化结果
// 只能用于HTTP回复的Set-Cookie头。
func (*Cookie) String() string

func (*ProtocolError) Error() string

// AddCookie adds a cookie to the request.  Per RFC 6265 section 5.4,
// AddCookie does not attach more than one Cookie header field.  That
// means all cookies, if any, are written into the same line,
// separated by semicolon.

// AddCookie向请求中添加一个cookie。按照RFC 6265 section 5.4的跪地，AddCookie不
// 会添加超过一个Cookie头字段。这表示所有的cookie都写在同一行，用分号分隔（
// cookie内部用逗号分隔属性）。
func (*Request) AddCookie(c *Cookie)

// BasicAuth returns the username and password provided in the request's
// Authorization header, if the request uses HTTP Basic Authentication. See RFC
// 2617, Section 2.
func (*Request) BasicAuth() (username, password string, ok bool)

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found.

// Cookie返回请求中名为name的cookie，如果未找到该cookie会返回nil, ErrNoCookie。
func (*Request) Cookie(name string) (*Cookie, error)

// Cookies parses and returns the HTTP cookies sent with the request.

// Cookies解析并返回该请求的Cookie头设置的cookie。
func (*Request) Cookies() []*Cookie

// FormFile returns the first file for the provided form key.
// FormFile calls ParseMultipartForm and ParseForm if necessary.

// FormFile返回以key为键查询r.MultipartForm字段得到结果中的第一个文件和它的信息
// 。如果必要，本函数会隐式调用ParseMultipartForm和ParseForm。查询失败会返回
// ErrMissingFile错误。
func (*Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)

// FormValue returns the first value for the named component of the query.
// POST and PUT body parameters take precedence over URL query string values.
// FormValue calls ParseMultipartForm and ParseForm if necessary and ignores
// any errors returned by these functions.
// If key is not present, FormValue returns the empty string.
// To access multiple values of the same key, call ParseForm and
// then inspect Request.Form directly.

// FormValue返回key为键查询r.Form字段得到结果[]string切片的第一个值。POST和PUT主
// 体中的同名参数优先于URL查询字符串。如果必要，本函数会隐式调用
// ParseMultipartForm和ParseForm。
func (*Request) FormValue(key string) string

// MultipartReader returns a MIME multipart reader if this is a
// multipart/form-data POST request, else returns nil and an error.
// Use this function instead of ParseMultipartForm to
// process the request body as a stream.

// 如果请求是multipart/form-data POST请求，MultipartReader返回一个
// multipart.Reader接口，否则返回nil和一个错误。使用本函数代替ParseMultipartForm
// ，可以将r.Body作为流处理。
func (*Request) MultipartReader() (*multipart.Reader, error)

// ParseForm parses the raw query from the URL and updates r.Form.
//
// For POST or PUT requests, it also parses the request body as a form and
// put the results into both r.PostForm and r.Form.
// POST and PUT body parameters take precedence over URL query string values
// in r.Form.
//
// If the request Body's size has not already been limited by MaxBytesReader,
// the size is capped at 10MB.
//
// ParseMultipartForm calls ParseForm automatically.
// It is idempotent.

// ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
//
// 对于POST或PUT请求，ParseForm还会将body当作表单解析，并将结果既更新到
// r.PostForm也更新到r.Form。解析结果中，POST或PUT请求主体要优先于URL查询字符串
// （同名变量，主体的值在查询字符串的值前面）。
//
// 如果请求的主体的大小没有被MaxBytesReader函数设定限制，其大小默认限制为开头
// 10MB。
//
// ParseMultipartForm会自动调用ParseForm。重复调用本方法是无意义的。
func (*Request) ParseForm() error

// ParseMultipartForm parses a request body as multipart/form-data.
// The whole request body is parsed and up to a total of maxMemory bytes of
// its file parts are stored in memory, with the remainder stored on
// disk in temporary files.
// ParseMultipartForm calls ParseForm if necessary.
// After one call to ParseMultipartForm, subsequent calls have no effect.

// ParseMultipartForm将请求的主体作为multipart/form-data解析。请求的整个主体都会
// 被解析，得到的文件记录最多maxMemery字节保存在内存，其余部分保存在硬盘的temp文
// 件里。如果必要，ParseMultipartForm会自行调用ParseForm。重复调用本方法是无意义
// 的。
func (*Request) ParseMultipartForm(maxMemory int64) error

// PostFormValue returns the first value for the named component of the POST
// or PUT request body. URL query parameters are ignored.
// PostFormValue calls ParseMultipartForm and ParseForm if necessary and ignores
// any errors returned by these functions.
// If key is not present, PostFormValue returns the empty string.

// PostFormValue返回key为键查询r.PostForm字段得到结果[]string切片的第一个值。如
// 果必要，本函数会隐式调用ParseMultipartForm和ParseForm。
func (*Request) PostFormValue(key string) string

// ProtoAtLeast reports whether the HTTP protocol used
// in the request is at least major.minor.

// ProtoAtLeast报告该请求使用的HTTP协议版本至少是major.minor。
func (*Request) ProtoAtLeast(major, minor int) bool

// Referer returns the referring URL, if sent in the request.
//
// Referer is misspelled as in the request itself, a mistake from the
// earliest days of HTTP.  This value can also be fetched from the
// Header map as Header["Referer"]; the benefit of making it available
// as a method is that the compiler can diagnose programs that use the
// alternate (correct English) spelling req.Referrer() but cannot
// diagnose programs that use Header["Referrer"].

// Referer返回请求中的访问来路信息。（请求的Referer头）
//
// Referer在请求中就是拼错了的，这是HTTP早期就有的错误。该值也可以从用
// Header["Referer"]获取； 让获取Referer字段变成方法的好处是，编译器可以诊断使用
// 正确单词拼法的req.Referrer()的程序，但却不能诊断使用Header["Referrer"]的程序
// 。
func (*Request) Referer() string

// SetBasicAuth sets the request's Authorization header to use HTTP
// Basic Authentication with the provided username and password.
//
// With HTTP Basic Authentication the provided username and password
// are not encrypted.

// SetBasicAuth使用提供的用户名和密码，采用HTTP基本认证，设置请求的Authorization
// 头。HTTP基本认证会明码传送用户名和密码。
func (*Request) SetBasicAuth(username, password string)

// UserAgent returns the client's User-Agent, if sent in the request.

// UserAgent返回请求中的客户端用户代理信息（请求的User-Agent头）。
func (*Request) UserAgent() string

// Write writes an HTTP/1.1 request, which is the header and body, in wire
// format. This method consults the following fields of the request:
//
//     Host
//     URL
//     Method (defaults to "GET")
//     Header
//     ContentLength
//     TransferEncoding
//     Body
//
// If Body is present, Content-Length is <= 0 and TransferEncoding hasn't been
// set to "identity", Write adds "Transfer-Encoding: chunked" to the header.
// Body is closed after it is sent.

// Write方法以有线格式将HTTP/1.1请求写入w（用于将请求写入下层TCPConn等）。本方法
// 会考虑请求的如下字段：
//
//     Host
//     URL
//     Method (defaults to "GET")
//     Header
//     ContentLength
//     TransferEncoding
//     Body
//
// 如果存在Body，ContentLength字段<= 0且TransferEncoding字段未显式设置为
// ["identity"]，Write方法会显式添加"Transfer-Encoding: chunked"到请求的头域。
// Body字段会在发送完请求后关闭。
func (*Request) Write(w io.Writer) error

// WriteProxy is like Write but writes the request in the form
// expected by an HTTP proxy.  In particular, WriteProxy writes the
// initial Request-URI line of the request with an absolute URI, per
// section 5.1.2 of RFC 2616, including the scheme and host.
// In either case, WriteProxy also writes a Host header, using
// either r.Host or r.URL.Host.

// WriteProxy类似Write但会将请求以HTTP代理期望的格式发送。
//
// 尤其是，按照RFC 2616 Section 5.1.2，WriteProxy会使用绝对URI（包括协议和主机名
// ）来初始化请求的第1行（Request-URI行）。无论何种情况，WriteProxy都会使用
// r.Host或r.URL.Host设置Host头。
func (*Request) WriteProxy(w io.Writer) error

// Cookies parses and returns the cookies set in the Set-Cookie headers.

// Cookies解析并返回该回复中的Set-Cookie头设置的cookie。
func (*Response) Cookies() []*Cookie

// Location returns the URL of the response's "Location" header,
// if present.  Relative redirects are resolved relative to
// the Response's Request.  ErrNoLocation is returned if no
// Location header is present.

// Location返回该回复的Location头设置的URL。相对地址的重定向会相对于该回复对应的
// 请求来确定绝对地址。如果回复中没有Location头，会返回nil, ErrNoLocation。
func (*Response) Location() (*url.URL, error)

// ProtoAtLeast reports whether the HTTP protocol used
// in the response is at least major.minor.

// ProtoAtLeast报告该回复使用的HTTP协议版本至少是major.minor。
func (*Response) ProtoAtLeast(major, minor int) bool

// Write writes r to w in the HTTP/1.n server response format,
// including the status line, headers, body, and optional trailer.
//
// This method consults the following fields of the response r:
//
//  StatusCode
//  ProtoMajor
//  ProtoMinor
//  Request.Method
//  TransferEncoding
//  Trailer
//  Body
//  ContentLength
//  Header, values for non-canonical keys will have unpredictable behavior
//
// The Response Body is closed after it is sent.

// Write以有线格式将回复写入w（用于将回复写入下层TCPConn等）。本方法会考虑如下字
// 段：
//
//     StatusCode
//     ProtoMajor
//     ProtoMinor
//     Request.Method
//     TransferEncoding
//     Trailer
//     Body
//     ContentLength
//     Header（不规范的键名和它对应的值会导致不可预知的行为）
//
// Body字段在发送完回复后会被关闭。
func (*Response) Write(w io.Writer) error

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.

// Handle注册HTTP处理器handler和对应的模式pattern。如果该模式已经注册有一个处理
// 器，Handle会panic。
func (*ServeMux) Handle(pattern string, handler Handler)

// HandleFunc registers the handler function for the given pattern.
func (*ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))

// Handler returns the handler to use for the given request, consulting
// r.Method, r.Host, and r.URL.Path. It always returns a non-nil handler. If the
// path is not in its canonical form, the handler will be an
// internally-generated handler that redirects to the canonical path.
//
// Handler also returns the registered pattern that matches the request or, in
// the case of internally-generated redirects, the pattern that will match after
// following the redirect.
//
// If there is no registered handler that applies to the request, Handler
// returns a ``page not found'' handler and an empty pattern.
func (*ServeMux) Handler(r *Request) (h Handler, pattern string)

// ServeHTTP dispatches the request to the handler whose pattern most closely
// matches the request URL.
func (*ServeMux) ServeHTTP(w ResponseWriter, r *Request)

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// If srv.Addr is blank, ":http" is used.
// ListenAndServe always returns a non-nil error.

// ListenAndServe监听srv.Addr指定的TCP地址，并且会调用Serve方法接收到的连接。如
// 果srv.Addr为空字符串，会使用":http"。
func (*Server) ListenAndServe() error

// ListenAndServeTLS listens on the TCP network address srv.Addr and
// then calls Serve to handle requests on incoming TLS connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// Filenames containing a certificate and matching private key for the
// server must be provided if neither the Server's TLSConfig.Certificates
// nor TLSConfig.GetCertificate are populated. If the certificate is
// signed by a certificate authority, the certFile should be the
// concatenation of the server's certificate, any intermediates, and
// the CA's certificate.
//
// If srv.Addr is blank, ":https" is used.
//
// ListenAndServeTLS always returns a non-nil error.

// ListenAndServeTLS监听srv.Addr确定的TCP地址，并且会调用Serve方法处理接收到的连
// 接。必须提供证书文件和对应的私钥文件。如果证书是由权威机构签发的，certFile参
// 数必须是顺序串联的服务端证书和CA证书。如果srv.Addr为空字符串，会使用":https"
// 。
func (*Server) ListenAndServeTLS(certFile, keyFile string) error

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each. The service goroutines read requests and
// then call srv.Handler to reply to them.
// Serve always returns a non-nil error.

// Serve会接手监听器l收到的每一个连接，并为每一个连接创建一个新的服务go程。该go
// 程会读取请求，然后调用srv.Handler回复请求。
func (*Server) Serve(l net.Listener) error

// SetKeepAlivesEnabled controls whether HTTP keep-alives are enabled.
// By default, keep-alives are always enabled. Only very
// resource-constrained environments or servers in the process of
// shutting down should disable them.

// SetKeepAlivesEnabled控制是否允许HTTP闲置连接重用（keep-alive）功能。默认该功
// 能总是被启用的。只有资源非常紧张的环境或者服务端在关闭进程中时，才应该关闭该
// 功能。
func (*Server) SetKeepAlivesEnabled(v bool)

// CancelRequest cancels an in-flight request by closing its connection.
// CancelRequest should only be called after RoundTrip has returned.
//
// Deprecated: Use Request.Cancel instead. CancelRequest can not cancel
// HTTP/2 requests.

// CancelRequest通过关闭请求所在的连接取消一个执行中的请求。
func (*Transport) CancelRequest(req *Request)

// CloseIdleConnections closes any connections which were previously
// connected from previous requests but are now sitting idle in
// a "keep-alive" state. It does not interrupt any connections currently
// in use.

// CloseIdleConnections关闭所有之前的请求建立但目前处于闲置状态的连接。本方法不
// 会中断正在使用的连接。
func (*Transport) CloseIdleConnections()

// RegisterProtocol registers a new protocol with scheme.
// The Transport will pass requests using the given scheme to rt.
// It is rt's responsibility to simulate HTTP request semantics.
//
// RegisterProtocol can be used by other packages to provide
// implementations of protocol schemes like "ftp" or "file".
//
// If rt.RoundTrip returns ErrSkipAltProtocol, the Transport will
// handle the RoundTrip itself for that one request, as if the
// protocol were not registered.

// RegisterProtocol注册一个新的名为scheme的协议。t会将使用scheme协议的请求转交给
// rt。rt有责任模拟HTTP请求的语义。
//
// RegisterProtocol可以被其他包用于提供"ftp"或"file"等协议的实现。
func (*Transport) RegisterProtocol(scheme string, rt RoundTripper)

// RoundTrip implements the RoundTripper interface.
//
// For higher-level HTTP client support (such as handling of cookies
// and redirects), see Get, Post, and the Client type.

// RoundTrip方法实现了RoundTripper接口。
//
// 高层次的HTTP客户端支持（如管理cookie和重定向）请参见Get、Post等函数和Client类
// 型。
func (*Transport) RoundTrip(req *Request) (resp *Response, err error)

func (ConnState) String() string

func (Dir) Open(name string) (File, error)

// ServeHTTP calls f(w, r).

// ServeHTTP方法会调用f(w, r)
func (HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.

// Add添加键值对到h，如键已存在则会将新的值附加到旧值切片后面。
func (Header) Add(key, value string)

// Del deletes the values associated with key.

// Del删除键值对。
func (Header) Del(key string)

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
// To access multiple values of a key, access the map directly
// with CanonicalHeaderKey.

// Get返回键对应的第一个值，如果键不存在会返回""。如要获取该键对应的值切片，请直
// 接用规范格式的键访问map。
func (Header) Get(key string) string

// Set sets the header entries associated with key to
// the single element value.  It replaces any existing
// values associated with key.

// Set添加键值对到h，如键已存在则会用只有新值一个元素的切片取代旧值切片。
func (Header) Set(key, value string)

// Write writes a header in wire format.

// Write以有线格式将头域写入w。
func (Header) Write(w io.Writer) error

// WriteSubset writes a header in wire format.
// If exclude is not nil, keys where exclude[key] == true are not written.

// WriteSubset以有线格式将头域写入w。当exclude不为nil时，如果h的键值对的键在
// exclude中存在且其对应值为真，该键值对就不会被写入w。
func (Header) WriteSubset(w io.Writer, exclude map[string]bool) error

