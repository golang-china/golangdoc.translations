// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package http provides HTTP client and server implementations.
//
// Get, Head, Post, and PostForm make HTTP (or HTTPS) requests:
//
//	resp, err := http.Get("http://example.com/")
//	...
//	resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
//	...
//	resp, err := http.PostForm("http://example.com/form",
//		url.Values{"key": {"Value"}, "id": {"123"}})
//
// The client must close the response body when finished with it:
//
//	resp, err := http.Get("http://example.com/")
//	if err != nil {
//		// handle error
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	// ...
//
// For control over HTTP client headers, redirect policy, and other settings,
// create a Client:
//
//	client := &http.Client{
//		CheckRedirect: redirectPolicyFunc,
//	}
//
//	resp, err := client.Get("http://example.com")
//	// ...
//
//	req, err := http.NewRequest("GET", "http://example.com", nil)
//	// ...
//	req.Header.Add("If-None-Match", `W/"wyzzy"`)
//	resp, err := client.Do(req)
//	// ...
//
// For control over proxies, TLS configuration, keep-alives, compression, and other
// settings, create a Transport:
//
//	tr := &http.Transport{
//		TLSClientConfig:    &tls.Config{RootCAs: pool},
//		DisableCompression: true,
//	}
//	client := &http.Client{Transport: tr}
//	resp, err := client.Get("https://example.com")
//
// Clients and Transports are safe for concurrent use by multiple goroutines and
// for efficiency should only be created once and re-used.
//
// ListenAndServe starts an HTTP server with a given address and handler. The
// handler is usually nil, which means to use DefaultServeMux. Handle and
// HandleFunc add handlers to DefaultServeMux:
//
//	http.Handle("/foo", fooHandler)
//
//	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//	})
//
//	log.Fatal(http.ListenAndServe(":8080", nil))
//
// More control over the server's behavior is available by creating a custom
// Server:
//
//	s := &http.Server{
//		Addr:           ":8080",
//		Handler:        myHandler,
//		ReadTimeout:    10 * time.Second,
//		WriteTimeout:   10 * time.Second,
//		MaxHeaderBytes: 1 << 20,
//	}
//	log.Fatal(s.ListenAndServe())
package http

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

// DefaultMaxHeaderBytes is the maximum permitted size of the headers in an HTTP
// request. This can be overridden by setting Server.MaxHeaderBytes.
const DefaultMaxHeaderBytes = 1 << 20 // 1 MB

// DefaultMaxIdleConnsPerHost is the default value of Transport's
// MaxIdleConnsPerHost.
const DefaultMaxIdleConnsPerHost = 2

// TimeFormat is the time format to use with time.Parse and time.Time.Format when
// parsing or generating times in HTTP headers. It is like time.RFC1123 but hard
// codes GMT as the time zone.
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

var (
	ErrHeaderTooLong        = &ProtocolError{"header too long"}
	ErrShortBody            = &ProtocolError{"entity body too short"}
	ErrNotSupported         = &ProtocolError{"feature not supported"}
	ErrUnexpectedTrailer    = &ProtocolError{"trailer header without chunked transfer encoding"}
	ErrMissingContentLength = &ProtocolError{"missing ContentLength in HEAD response"}
	ErrNotMultipart         = &ProtocolError{"request Content-Type isn't multipart/form-data"}
	ErrMissingBoundary      = &ProtocolError{"no multipart boundary param in Content-Type"}
)

// Errors introduced by the HTTP server.
var (
	ErrWriteAfterFlush = errors.New("Conn.Write called after Flush")
	ErrBodyNotAllowed  = errors.New("http: request method or response status code does not allow body")
	ErrHijacked        = errors.New("Conn has been hijacked")
	ErrContentLength   = errors.New("Conn.Write wrote more than the declared Content-Length")
)

// DefaultClient is the default Client and is used by Get, Head, and Post.
var DefaultClient = &Client{}

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = NewServeMux()

// ErrBodyReadAfterClose is returned when reading a Request or Response Body after
// the body has been closed. This typically happens when the body is read after an
// HTTP Handler calls WriteHeader or Write on its ResponseWriter.
var ErrBodyReadAfterClose = errors.New("http: invalid Read on closed Body")

// ErrHandlerTimeout is returned on ResponseWriter Write calls in handlers which
// have timed out.
var ErrHandlerTimeout = errors.New("http: Handler timeout")

// ErrLineTooLong is returned when reading request or response bodies with
// malformed chunked encoding.
var ErrLineTooLong = internal.ErrLineTooLong

// ErrMissingFile is returned by FormFile when the provided file field name is
// either not present in the request or not a file field.
var ErrMissingFile = errors.New("http: no such file")

var ErrNoCookie = errors.New("http: named cookie not present")

var ErrNoLocation = errors.New("http: no Location header in response")

// CanonicalHeaderKey returns the canonical format of the header key s. The
// canonicalization converts the first letter and any letter following a hyphen to
// upper case; the rest are converted to lowercase. For example, the canonical key
// for "accept-encoding" is "Accept-Encoding".
func CanonicalHeaderKey(s string) string

// DetectContentType implements the algorithm described at
// http://mimesniff.spec.whatwg.org/ to determine the Content-Type of the given
// data. It considers at most the first 512 bytes of data. DetectContentType always
// returns a valid MIME type: if it cannot determine a more specific one, it
// returns "application/octet-stream".
func DetectContentType(data []byte) string

// Error replies to the request with the specified error message and HTTP code. The
// error message should be plain text.
func Error(w ResponseWriter, error string, code int)

// Handle registers the handler for the given pattern in the DefaultServeMux. The
// documentation for ServeMux explains how patterns are matched.
func Handle(pattern string, handler Handler)

// HandleFunc registers the handler function for the given pattern in the
// DefaultServeMux. The documentation for ServeMux explains how patterns are
// matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))

// ListenAndServe listens on the TCP network address addr and then calls Serve with
// handler to handle requests on incoming connections. Handler is typically nil, in
// which case the DefaultServeMux is used.
//
// A trivial example server is:
//
//	package main
//
//	import (
//		"io"
//		"net/http"
//		"log"
//	)
//
//	// hello world, the web server
//	func HelloServer(w http.ResponseWriter, req *http.Request) {
//		io.WriteString(w, "hello, world!\n")
//	}
//
//	func main() {
//		http.HandleFunc("/hello", HelloServer)
//		err := http.ListenAndServe(":12345", nil)
//		if err != nil {
//			log.Fatal("ListenAndServe: ", err)
//		}
//	}
func ListenAndServe(addr string, handler Handler) error

// ListenAndServeTLS acts identically to ListenAndServe, except that it expects
// HTTPS connections. Additionally, files containing a certificate and matching
// private key for the server must be provided. If the certificate is signed by a
// certificate authority, the certFile should be the concatenation of the server's
// certificate followed by the CA's certificate.
//
// A trivial example server is:
//
//	import (
//		"log"
//		"net/http"
//	)
//
//	func handler(w http.ResponseWriter, req *http.Request) {
//		w.Header().Set("Content-Type", "text/plain")
//		w.Write([]byte("This is an example server.\n"))
//	}
//
//	func main() {
//		http.HandleFunc("/", handler)
//		log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
//		err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler) error

// MaxBytesReader is similar to io.LimitReader but is intended for limiting the
// size of incoming request bodies. In contrast to io.LimitReader, MaxBytesReader's
// result is a ReadCloser, returns a non-EOF error for a Read beyond the limit, and
// Closes the underlying reader when its Close method is called.
//
// MaxBytesReader prevents clients from accidentally or maliciously sending a large
// request and wasting server resources.
func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser

// NotFound replies to the request with an HTTP 404 not found error.
func NotFound(w ResponseWriter, r *Request)

// ParseHTTPVersion parses a HTTP version string. "HTTP/1.0" returns (1, 0, true).
func ParseHTTPVersion(vers string) (major, minor int, ok bool)

// ParseTime parses a time header (such as the Date: header), trying each of the
// three formats allowed by HTTP/1.1: TimeFormat, time.RFC850, and time.ANSIC.
func ParseTime(text string) (t time.Time, err error)

// ProxyFromEnvironment returns the URL of the proxy to use for a given request, as
// indicated by the environment variables HTTP_PROXY, HTTPS_PROXY and NO_PROXY (or
// the lowercase versions thereof). HTTPS_PROXY takes precedence over HTTP_PROXY
// for https requests.
//
// The environment values may be either a complete URL or a "host[:port]", in which
// case the "http" scheme is assumed. An error is returned if the value is a
// different form.
//
// A nil URL and nil error are returned if no proxy is defined in the environment,
// or a proxy should not be used for the given request, as defined by NO_PROXY.
//
// As a special case, if req.URL.Host is "localhost" (with or without a port
// number), then a nil URL and nil error will be returned.
func ProxyFromEnvironment(req *Request) (*url.URL, error)

// ProxyURL returns a proxy function (for use in a Transport) that always returns
// the same URL.
func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error)

// Redirect replies to the request with a redirect to url, which may be a path
// relative to the request path.
func Redirect(w ResponseWriter, r *Request, urlStr string, code int)

// Serve accepts incoming HTTP connections on the listener l, creating a new
// service goroutine for each. The service goroutines read requests and then call
// handler to reply to them. Handler is typically nil, in which case the
// DefaultServeMux is used.
func Serve(l net.Listener, handler Handler) error

// ServeContent replies to the request using the content in the provided
// ReadSeeker. The main benefit of ServeContent over io.Copy is that it handles
// Range requests properly, sets the MIME type, and handles If-Modified-Since
// requests.
//
// If the response's Content-Type header is not set, ServeContent first tries to
// deduce the type from name's file extension and, if that fails, falls back to
// reading the first block of the content and passing it to DetectContentType. The
// name is otherwise unused; in particular it can be empty and is never sent in the
// response.
//
// If modtime is not the zero time, ServeContent includes it in a Last-Modified
// header in the response. If the request includes an If-Modified-Since header,
// ServeContent uses modtime to decide whether the content needs to be sent at all.
//
// The content's Seek method must work: ServeContent uses a seek to the end of the
// content to determine its size.
//
// If the caller has set w's ETag header, ServeContent uses it to handle requests
// using If-Range and If-None-Match.
//
// Note that *os.File implements the io.ReadSeeker interface.
func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker)

// ServeFile replies to the request with the contents of the named file or
// directory.
func ServeFile(w ResponseWriter, r *Request, name string)

// SetCookie adds a Set-Cookie header to the provided ResponseWriter's headers.
func SetCookie(w ResponseWriter, cookie *Cookie)

// StatusText returns a text for the HTTP status code. It returns the empty string
// if the code is unknown.
func StatusText(code int) string

// A Client is an HTTP client. Its zero value (DefaultClient) is a usable client
// that uses DefaultTransport.
//
// The Client's Transport typically has internal state (cached TCP connections), so
// Clients should be reused instead of created as needed. Clients are safe for
// concurrent use by multiple goroutines.
//
// A Client is higher-level than a RoundTripper (such as Transport) and
// additionally handles HTTP details such as cookies and redirects.
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

// Do sends an HTTP request and returns an HTTP response, following policy (e.g.
// redirects, cookies, auth) as configured on the client.
//
// An error is returned if caused by client policy (such as CheckRedirect), or if
// there was an HTTP protocol error. A non-2xx response doesn't cause an error.
//
// When err is nil, resp always contains a non-nil resp.Body.
//
// Callers should close resp.Body when done reading from it. If resp.Body is not
// closed, the Client's underlying RoundTripper (typically Transport) may not be
// able to re-use a persistent TCP connection to the server for a subsequent
// "keep-alive" request.
//
// The request Body, if non-nil, will be closed by the underlying Transport, even
// on errors.
//
// Generally Get, Post, or PostForm will be used instead of Do.
func (c *Client) Do(req *Request) (resp *Response, err error)

// Get issues a GET to the specified URL. If the response is one of the following
// redirect codes, Get follows the redirect after calling the Client's
// CheckRedirect function.
//
//	301 (Moved Permanently)
//	302 (Found)
//	303 (See Other)
//	307 (Temporary Redirect)
//
// An error is returned if the Client's CheckRedirect function fails or if there
// was an HTTP protocol error. A non-2xx response doesn't cause an error.
//
// When err is nil, resp always contains a non-nil resp.Body. Caller should close
// resp.Body when done reading from it.
func (c *Client) Get(url string) (resp *Response, err error)

// Head issues a HEAD to the specified URL. If the response is one of the following
// redirect codes, Head follows the redirect after calling the Client's
// CheckRedirect function.
//
//	301 (Moved Permanently)
//	302 (Found)
//	303 (See Other)
//	307 (Temporary Redirect)
func (c *Client) Head(url string) (resp *Response, err error)

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// If the provided body is also an io.Closer, it is closed after the request.
func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *Response, err error)

// PostForm issues a POST to the specified URL, with data's keys and values
// urlencoded as the request body.
//
// When err is nil, resp always contains a non-nil resp.Body. Caller should close
// resp.Body when done reading from it.
func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)

// The CloseNotifier interface is implemented by ResponseWriters which allow
// detecting when the underlying connection has gone away.
//
// This mechanism can be used to cancel long operations on the server if the client
// has disconnected before the response is ready.
type CloseNotifier interface {
	// CloseNotify returns a channel that receives a single value
	// when the client connection has gone away.
	CloseNotify() <-chan bool
}

// A ConnState represents the state of a client connection to a server. It's used
// by the optional Server.ConnState hook.
type ConnState int

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

func (c ConnState) String() string

// A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an HTTP
// response or the Cookie header of an HTTP request.
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

// String returns the serialization of the cookie for use in a Cookie header (if
// only Name and Value are set) or a Set-Cookie response header (if other fields
// are set).
func (c *Cookie) String() string

// A CookieJar manages storage and use of cookies in HTTP requests.
//
// Implementations of CookieJar must be safe for concurrent use by multiple
// goroutines.
//
// The net/http/cookiejar package provides a CookieJar implementation.
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
// While the FileSystem.Open method takes '/'-separated paths, a Dir's string value
// is a filename on the native file system, not a URL, so it is separated by
// filepath.Separator, which isn't necessarily '/'.
//
// An empty Dir is treated as ".".
type Dir string

func (d Dir) Open(name string) (File, error)

// A File is returned by a FileSystem's Open method and can be served by the
// FileServer implementation.
//
// The methods should behave the same as those on an *os.File.
type File interface {
	io.Closer
	io.Reader
	Readdir(count int) ([]os.FileInfo, error)
	Seek(offset int64, whence int) (int64, error)
	Stat() (os.FileInfo, error)
}

// A FileSystem implements access to a collection of named files. The elements in a
// file path are separated by slash ('/', U+002F) characters, regardless of host
// operating system convention.
type FileSystem interface {
	Open(name string) (File, error)
}

// The Flusher interface is implemented by ResponseWriters that allow an HTTP
// handler to flush buffered data to the client.
//
// Note that even for ResponseWriters that support Flush, if the client is
// connected through an HTTP proxy, the buffered data may not reach the client
// until the response completes.
type Flusher interface {
	// Flush sends any buffered data to the client.
	Flush()
}

// Objects implementing the Handler interface can be registered to serve a
// particular path or subtree in the HTTP server.
//
// ServeHTTP should write reply headers and data to the ResponseWriter and then
// return. Returning signals that the request is finished and that the HTTP server
// can move on to the next request on the connection.
//
// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes that the
// effect of the panic was isolated to the active request. It recovers the panic,
// logs a stack trace to the server error log, and hangs up the connection.
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

// FileServer returns a handler that serves HTTP requests with the contents of the
// file system rooted at root.
//
// To use the operating system's file system implementation, use http.Dir:
//
//	http.Handle("/", http.FileServer(http.Dir("/tmp")))
func FileServer(root FileSystem) Handler

// NotFoundHandler returns a simple request handler that replies to each request
// with a ``404 page not found'' reply.
func NotFoundHandler() Handler

// RedirectHandler returns a request handler that redirects each request it
// receives to the given url using the given status code.
func RedirectHandler(url string, code int) Handler

// StripPrefix returns a handler that serves HTTP requests by removing the given
// prefix from the request URL's Path and invoking the handler h. StripPrefix
// handles a request for a path that doesn't begin with prefix by replying with an
// HTTP 404 not found error.
func StripPrefix(prefix string, h Handler) Handler

// TimeoutHandler returns a Handler that runs h with the given time limit.
//
// The new Handler calls h.ServeHTTP to handle each request, but if a call runs for
// longer than its time limit, the handler responds with a 503 Service Unavailable
// error and the given message in its body. (If msg is empty, a suitable default
// message will be sent.) After such a timeout, writes by h to its ResponseWriter
// will return ErrHandlerTimeout.
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler

// The HandlerFunc type is an adapter to allow the use of ordinary functions as
// HTTP handlers. If f is a function with the appropriate signature, HandlerFunc(f)
// is a Handler object that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)

// A Header represents the key-value pairs in an HTTP header.
type Header map[string][]string

// Add adds the key, value pair to the header. It appends to any existing values
// associated with key.
func (h Header) Add(key, value string)

// Del deletes the values associated with key.
func (h Header) Del(key string)

// Get gets the first value associated with the given key. If there are no values
// associated with the key, Get returns "". To access multiple values of a key,
// access the map directly with CanonicalHeaderKey.
func (h Header) Get(key string) string

// Set sets the header entries associated with key to the single element value. It
// replaces any existing values associated with key.
func (h Header) Set(key, value string)

// Write writes a header in wire format.
func (h Header) Write(w io.Writer) error

// WriteSubset writes a header in wire format. If exclude is not nil, keys where
// exclude[key] == true are not written.
func (h Header) WriteSubset(w io.Writer, exclude map[string]bool) error

// The Hijacker interface is implemented by ResponseWriters that allow an HTTP
// handler to take over the connection.
type Hijacker interface {
	// Hijack lets the caller take over the connection.
	// After a call to Hijack(), the HTTP server library
	// will not do anything else with the connection.
	// It becomes the caller's responsibility to manage
	// and close the connection.
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}

// HTTP request parsing errors.
type ProtocolError struct {
	ErrorString string
}

func (err *ProtocolError) Error() string

// A Request represents an HTTP request received by a server or to be sent by a
// client.
//
// The field semantics differ slightly between client and server usage. In addition
// to the notes on the fields below, see the documentation for Request.Write and
// RoundTripper.
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

// NewRequest returns a new Request given a method, URL, and optional body.
//
// If the provided body is also an io.Closer, the returned Request.Body is set to
// body and will be closed by the Client methods Do, Post, and PostForm, and
// Transport.RoundTrip.
func NewRequest(method, urlStr string, body io.Reader) (*Request, error)

// ReadRequest reads and parses a request from b.
func ReadRequest(b *bufio.Reader) (req *Request, err error)

// AddCookie adds a cookie to the request. Per RFC 6265 section 5.4, AddCookie does
// not attach more than one Cookie header field. That means all cookies, if any,
// are written into the same line, separated by semicolon.
func (r *Request) AddCookie(c *Cookie)

// BasicAuth returns the username and password provided in the request's
// Authorization header, if the request uses HTTP Basic Authentication. See RFC
// 2617, Section 2.
func (r *Request) BasicAuth() (username, password string, ok bool)

// Cookie returns the named cookie provided in the request or ErrNoCookie if not
// found.
func (r *Request) Cookie(name string) (*Cookie, error)

// Cookies parses and returns the HTTP cookies sent with the request.
func (r *Request) Cookies() []*Cookie

// FormFile returns the first file for the provided form key. FormFile calls
// ParseMultipartForm and ParseForm if necessary.
func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)

// FormValue returns the first value for the named component of the query. POST and
// PUT body parameters take precedence over URL query string values. FormValue
// calls ParseMultipartForm and ParseForm if necessary and ignores any errors
// returned by these functions. To access multiple values of the same key, call
// ParseForm and then inspect Request.Form directly.
func (r *Request) FormValue(key string) string

// MultipartReader returns a MIME multipart reader if this is a multipart/form-data
// POST request, else returns nil and an error. Use this function instead of
// ParseMultipartForm to process the request body as a stream.
func (r *Request) MultipartReader() (*multipart.Reader, error)

// ParseForm parses the raw query from the URL and updates r.Form.
//
// For POST or PUT requests, it also parses the request body as a form and put the
// results into both r.PostForm and r.Form. POST and PUT body parameters take
// precedence over URL query string values in r.Form.
//
// If the request Body's size has not already been limited by MaxBytesReader, the
// size is capped at 10MB.
//
// ParseMultipartForm calls ParseForm automatically. It is idempotent.
func (r *Request) ParseForm() error

// ParseMultipartForm parses a request body as multipart/form-data. The whole
// request body is parsed and up to a total of maxMemory bytes of its file parts
// are stored in memory, with the remainder stored on disk in temporary files.
// ParseMultipartForm calls ParseForm if necessary. After one call to
// ParseMultipartForm, subsequent calls have no effect.
func (r *Request) ParseMultipartForm(maxMemory int64) error

// PostFormValue returns the first value for the named component of the POST or PUT
// request body. URL query parameters are ignored. PostFormValue calls
// ParseMultipartForm and ParseForm if necessary and ignores any errors returned by
// these functions.
func (r *Request) PostFormValue(key string) string

// ProtoAtLeast reports whether the HTTP protocol used in the request is at least
// major.minor.
func (r *Request) ProtoAtLeast(major, minor int) bool

// Referer returns the referring URL, if sent in the request.
//
// Referer is misspelled as in the request itself, a mistake from the earliest days
// of HTTP. This value can also be fetched from the Header map as
// Header["Referer"]; the benefit of making it available as a method is that the
// compiler can diagnose programs that use the alternate (correct English) spelling
// req.Referrer() but cannot diagnose programs that use Header["Referrer"].
func (r *Request) Referer() string

// SetBasicAuth sets the request's Authorization header to use HTTP Basic
// Authentication with the provided username and password.
//
// With HTTP Basic Authentication the provided username and password are not
// encrypted.
func (r *Request) SetBasicAuth(username, password string)

// UserAgent returns the client's User-Agent, if sent in the request.
func (r *Request) UserAgent() string

// Write writes an HTTP/1.1 request -- header and body -- in wire format. This
// method consults the following fields of the request:
//
//	Host
//	URL
//	Method (defaults to "GET")
//	Header
//	ContentLength
//	TransferEncoding
//	Body
//
// If Body is present, Content-Length is <= 0 and TransferEncoding hasn't been set
// to "identity", Write adds "Transfer-Encoding: chunked" to the header. Body is
// closed after it is sent.
func (r *Request) Write(w io.Writer) error

// WriteProxy is like Write but writes the request in the form expected by an HTTP
// proxy. In particular, WriteProxy writes the initial Request-URI line of the
// request with an absolute URI, per section 5.1.2 of RFC 2616, including the
// scheme and host. In either case, WriteProxy also writes a Host header, using
// either r.Host or r.URL.Host.
func (r *Request) WriteProxy(w io.Writer) error

// Response represents the response from an HTTP request.
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

// Get issues a GET to the specified URL. If the response is one of the following
// redirect codes, Get follows the redirect, up to a maximum of 10 redirects:
//
//	301 (Moved Permanently)
//	302 (Found)
//	303 (See Other)
//	307 (Temporary Redirect)
//
// An error is returned if there were too many redirects or if there was an HTTP
// protocol error. A non-2xx response doesn't cause an error.
//
// When err is nil, resp always contains a non-nil resp.Body. Caller should close
// resp.Body when done reading from it.
//
// Get is a wrapper around DefaultClient.Get.
func Get(url string) (resp *Response, err error)

// Head issues a HEAD to the specified URL. If the response is one of the following
// redirect codes, Head follows the redirect after calling the Client's
// CheckRedirect function.
//
//	301 (Moved Permanently)
//	302 (Found)
//	303 (See Other)
//	307 (Temporary Redirect)
//
// Head is a wrapper around DefaultClient.Head
func Head(url string) (resp *Response, err error)

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// Post is a wrapper around DefaultClient.Post
func Post(url string, bodyType string, body io.Reader) (resp *Response, err error)

// PostForm issues a POST to the specified URL, with data's keys and values
// URL-encoded as the request body.
//
// When err is nil, resp always contains a non-nil resp.Body. Caller should close
// resp.Body when done reading from it.
//
// PostForm is a wrapper around DefaultClient.PostForm
func PostForm(url string, data url.Values) (resp *Response, err error)

// ReadResponse reads and returns an HTTP response from r. The req parameter
// optionally specifies the Request that corresponds to this Response. If nil, a
// GET request is assumed. Clients must call resp.Body.Close when finished reading
// resp.Body. After that call, clients can inspect resp.Trailer to find key/value
// pairs included in the response trailer.
func ReadResponse(r *bufio.Reader, req *Request) (*Response, error)

// Cookies parses and returns the cookies set in the Set-Cookie headers.
func (r *Response) Cookies() []*Cookie

// Location returns the URL of the response's "Location" header, if present.
// Relative redirects are resolved relative to the Response's Request.
// ErrNoLocation is returned if no Location header is present.
func (r *Response) Location() (*url.URL, error)

// ProtoAtLeast reports whether the HTTP protocol used in the response is at least
// major.minor.
func (r *Response) ProtoAtLeast(major, minor int) bool

// Writes the response (header, body and trailer) in wire format. This method
// consults the following fields of the response:
//
//	StatusCode
//	ProtoMajor
//	ProtoMinor
//	Request.Method
//	TransferEncoding
//	Trailer
//	Body
//	ContentLength
//	Header, values for non-canonical keys will have unpredictable behavior
//
// Body is closed after it is sent.
func (r *Response) Write(w io.Writer) error

// A ResponseWriter interface is used by an HTTP handler to construct an HTTP
// response.
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

// RoundTripper is an interface representing the ability to execute a single HTTP
// transaction, obtaining the Response for a given Request.
//
// A RoundTripper must be safe for concurrent use by multiple goroutines.
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

// DefaultTransport is the default implementation of Transport and is used by
// DefaultClient. It establishes network connections as needed and caches them for
// reuse by subsequent calls. It uses HTTP proxies as directed by the $HTTP_PROXY
// and $NO_PROXY (or $http_proxy and $no_proxy) environment variables.
var DefaultTransport RoundTripper = &Transport{
	Proxy: ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
}

// NewFileTransport returns a new RoundTripper, serving the provided FileSystem.
// The returned RoundTripper ignores the URL host in its incoming requests, as well
// as most other properties of the request.
//
// The typical use case for NewFileTransport is to register the "file" protocol
// with a Transport, as in:
//
//	t := &http.Transport{}
//	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
//	c := &http.Client{Transport: t}
//	res, err := c.Get("file:///etc/passwd")
//	...
func NewFileTransport(fs FileSystem) RoundTripper

// ServeMux is an HTTP request multiplexer. It matches the URL of each incoming
// request against a list of registered patterns and calls the handler for the
// pattern that most closely matches the URL.
//
// Patterns name fixed, rooted paths, like "/favicon.ico", or rooted subtrees, like
// "/images/" (note the trailing slash). Longer patterns take precedence over
// shorter ones, so that if there are handlers registered for both "/images/" and
// "/images/thumbnails/", the latter handler will be called for paths beginning
// "/images/thumbnails/" and the former will receive requests for any other paths
// in the "/images/" subtree.
//
// Note that since a pattern ending in a slash names a rooted subtree, the pattern
// "/" matches all paths not matched by other registered patterns, not just the URL
// with Path == "/".
//
// Patterns may optionally begin with a host name, restricting matches to URLs on
// that host only. Host-specific patterns take precedence over general patterns, so
// that a handler might register for the two patterns "/codesearch" and
// "codesearch.google.com/" without also taking over requests for
// "http://www.google.com/".
//
// ServeMux also takes care of sanitizing the URL request path, redirecting any
// request containing . or .. elements to an equivalent .- and ..-free URL.
type ServeMux struct {
	// contains filtered or unexported fields
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux

// Handle registers the handler for the given pattern. If a handler already exists
// for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler)

// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))

// Handler returns the handler to use for the given request, consulting r.Method,
// r.Host, and r.URL.Path. It always returns a non-nil handler. If the path is not
// in its canonical form, the handler will be an internally-generated handler that
// redirects to the canonical path.
//
// Handler also returns the registered pattern that matches the request or, in the
// case of internally-generated redirects, the pattern that will match after
// following the redirect.
//
// If there is no registered handler that applies to the request, Handler returns a
// ``page not found'' handler and an empty pattern.
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)

// ServeHTTP dispatches the request to the handler whose pattern most closely
// matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)

// A Server defines parameters for running an HTTP server. The zero value for
// Server is a valid configuration.
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
	// contains filtered or unexported fields
}

// ListenAndServe listens on the TCP network address srv.Addr and then calls Serve
// to handle requests on incoming connections. If srv.Addr is blank, ":http" is
// used.
func (srv *Server) ListenAndServe() error

// ListenAndServeTLS listens on the TCP network address srv.Addr and then calls
// Serve to handle requests on incoming TLS connections.
//
// Filenames containing a certificate and matching private key for the server must
// be provided. If the certificate is signed by a certificate authority, the
// certFile should be the concatenation of the server's certificate followed by the
// CA's certificate.
//
// If srv.Addr is blank, ":https" is used.
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error

// Serve accepts incoming connections on the Listener l, creating a new service
// goroutine for each. The service goroutines read requests and then call
// srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error

// SetKeepAlivesEnabled controls whether HTTP keep-alives are enabled. By default,
// keep-alives are always enabled. Only very resource-constrained environments or
// servers in the process of shutting down should disable them.
func (s *Server) SetKeepAlivesEnabled(v bool)

// Transport is an implementation of RoundTripper that supports HTTP, HTTPS, and
// HTTP proxies (for either HTTP or HTTPS with CONNECT). Transport can also cache
// connections for future re-use.
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
	// contains filtered or unexported fields
}

// CancelRequest cancels an in-flight request by closing its connection.
func (t *Transport) CancelRequest(req *Request)

// CloseIdleConnections closes any connections which were previously connected from
// previous requests but are now sitting idle in a "keep-alive" state. It does not
// interrupt any connections currently in use.
func (t *Transport) CloseIdleConnections()

// RegisterProtocol registers a new protocol with scheme. The Transport will pass
// requests using the given scheme to rt. It is rt's responsibility to simulate
// HTTP request semantics.
//
// RegisterProtocol can be used by other packages to provide implementations of
// protocol schemes like "ftp" or "file".
func (t *Transport) RegisterProtocol(scheme string, rt RoundTripper)

// RoundTrip implements the RoundTripper interface.
//
// For higher-level HTTP client support (such as handling of cookies and
// redirects), see Get, Post, and the Client type.
func (t *Transport) RoundTrip(req *Request) (resp *Response, err error)
