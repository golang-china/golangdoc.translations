// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package rpc provides access to the exported methods of an object across a
// network or other I/O connection. A server registers an object, making it
// visible as a service with the name of the type of the object. After
// registration, exported methods of the object will be accessible remotely. A
// server may register multiple objects (services) of different types but it is
// an error to register multiple objects of the same type.
//
// Only methods that satisfy these criteria will be made available for remote
// access; other methods will be ignored:
//
//     - the method's type is exported.
//     - the method is exported.
//     - the method has two arguments, both exported (or builtin) types.
//     - the method's second argument is a pointer.
//     - the method has return type error.
//
// In effect, the method must look schematically like
//
//     func (t *T) MethodName(argType T1, replyType *T2) error
//
// where T, T1 and T2 can be marshaled by encoding/gob. These requirements apply
// even if a different codec is used. (In the future, these requirements may
// soften for custom codecs.)
//
// The method's first argument represents the arguments provided by the caller;
// the second argument represents the result parameters to be returned to the
// caller. The method's return value, if non-nil, is passed back as a string
// that the client sees as if created by errors.New. If an error is returned,
// the reply parameter will not be sent back to the client.
//
// The server may handle requests on a single connection by calling ServeConn.
// More typically it will create a network listener and call Accept or, for an
// HTTP listener, HandleHTTP and http.Serve.
//
// A client wishing to use the service establishes a connection and then invokes
// NewClient on the connection. The convenience function Dial (DialHTTP)
// performs both steps for a raw network connection (an HTTP connection). The
// resulting Client object has two methods, Call and Go, that specify the
// service and method to call, a pointer containing the arguments, and a pointer
// to receive the result parameters.
//
// The Call method waits for the remote call to complete while the Go method
// launches the call asynchronously and signals completion using the Call
// structure's Done channel.
//
// Unless an explicit codec is set up, package encoding/gob is used to transport
// the data.
//
// Here is a simple example. A server wishes to export an object of type Arith:
//
//     package server
//
//     type Args struct {
//         A, B int
//     }
//
//     type Quotient struct {
//         Quo, Rem int
//     }
//
//     type Arith int
//
//     func (t *Arith) Multiply(args *Args, reply *int) error {
//         *reply = args.A * args.B
//         return nil
//     }
//
//     func (t *Arith) Divide(args *Args, quo *Quotient) error {
//         if args.B == 0 {
//             return errors.New("divide by zero")
//         }
//         quo.Quo = args.A / args.B
//         quo.Rem = args.A % args.B
//         return nil
//     }
//
// The server calls (for HTTP service):
//
//     arith := new(Arith)
//     rpc.Register(arith)
//     rpc.HandleHTTP()
//     l, e := net.Listen("tcp", ":1234")
//     if e != nil {
//         log.Fatal("listen error:", e)
//     }
//     go http.Serve(l, nil)
//
// At this point, clients can see a service "Arith" with methods
// "Arith.Multiply" and "Arith.Divide". To invoke one, a client first dials the
// server:
//
//     client, err := rpc.DialHTTP("tcp", serverAddress + ":1234")
//     if err != nil {
//         log.Fatal("dialing:", err)
//     }
//
// Then it can make a remote call:
//
//     // Synchronous call
//     args := &server.Args{7,8}
//     var reply int
//     err = client.Call("Arith.Multiply", args, &reply)
//     if err != nil {
//         log.Fatal("arith error:", err)
//     }
//     fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
//
// or
//
//     // Asynchronous call
//     quotient := new(Quotient)
//     divCall := client.Go("Arith.Divide", args, quotient, nil)
//     replyCall := <-divCall.Done    // will be equal to divCall
//     // check errors, print, etc.
//
// A server implementation will often provide a simple, type-safe wrapper for
// the client.

// rpc 包提供了一个方法来通过网络或者其他的I/O连接进入对象的外部方法. 一个server
// 注册一个对象， 标记它成为可见对象类型名字的服务。注册后，对象的外部方法就可以
// 远程调用了。一个server可以注册多个 不同类型的对象，但是却不可以注册多个相同类
// 型的对象。
//
// 只有满足这些标准的方法才会被远程调用视为可见；其他的方法都会被忽略：
//
//     - 方法是外部可见的。
//     - 方法有两个参数，参数的类型都是外部可见的。
//     - 方法的第二个参数是一个指针。
//     - 方法有返回类型错误
//
// 事实上，方法必须看起来类似这样
//
//     func (t *T) MethodName(argType T1, replyType *T2) error
//
// T，T1和T2可以被encoding/gob序列化。 不管使用什么编解码，这些要求都要满足。 （
// 在未来，这些要求可能对自定义的编解码会放宽）
//
// 方法的第一个参数代表调用者提供的参数；第二个参数代表返回给调用者的参数。方法
// 的返回值，如果是非空的话 就会被作为一个string返回，客户端会error像是被
// errors.New调用返回的一样。如果error返回的话， 返回的参数将会被送回给客户端。
//
// 服务断可以使用ServeConn来处理单个连接上的请求。更通用的方法，服务器可以制造一
// 个网络监听，然后调用 Accept，或者对一个HTTP监听，处理HandleHTTP和http.Serve。
//
// 客户端希望使用服务来建立连接，然后在连接上调用NewClient来建立连接。更方便的方
// 法就是调用Dial(DialHTTP) 来建立一个新的网络连接（一个HTTP连接）。客户端获得到
// 的对象有两个方法，Call和Go，指定的参数有：服务和方法 指向参数的指针，接受返回
// 结果的指针。
//
// call方法等待远程调用完成，但Go方法是异步调用call方法，使用Call通道来标志调用
// 完成。
//
// 除非有明确制定编解码器，否则默认使用encoding/gob来传输数据。
//
// 这是个简单的例子，服务器希望对外服务出Arith对象：
//
//     package server
//
//     type Args struct {
//         A, B int
//     }
//
//     type Quotient struct {
//         Quo, Rem int
//     }
//
//     type Arith int
//
//     func (t *Arith) Multiply(args *Args, reply *int) error {
//         *reply = args.A * args.B
//         return nil
//     }
//
//     func (t *Arith) Divide(args *Args, quo *Quotient) error {
//         if args.B == 0 {
//             return errors.New("divide by zero")
//         }
//         quo.Quo = args.A / args.B
//         quo.Rem = args.A % args.B
//         return nil
//     }
//
// 服务端调用（使用HTTP服务）：
//
//     arith := new(Arith)
//     rpc.Register(arith)
//     rpc.HandleHTTP()
//     l, e := net.Listen("tcp", ":1234")
//     if e != nil {
//         log.Fatal("listen error:", e)
//     }
//     go http.Serve(l, nil)
//
// 在这个时候，客户端可以看见服务“Arith”，并且有“Arith.Multiply”方法和“
// Arith.Divide”方法。 调用其中一个，客户端首先连接服务：
//
//     client, err := rpc.DialHTTP("tcp", serverAddress + ":1234")
//     if err != nil {
//         log.Fatal("dialing:", err)
//     }
//
// 当它要调用远程服务的时候：
//
//     // Synchronous call
//     args := &server.Args{7,8}
//     var reply int
//     err = client.Call("Arith.Multiply", args, &reply)
//     if err != nil {
//         log.Fatal("arith error:", err)
//     }
//     fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
//
// or
//
//     // Asynchronous call
//     quotient := new(Quotient)
//     divCall := client.Go("Arith.Divide", args, quotient, nil)
//     replyCall := <-divCall.Done    // will be equal to divCall
//     // check errors, print, etc.
//
// 服务端的实现需要为客户端提供一个简单的，类型安全服务。
package rpc

import (
    "bufio"
    "encoding/gob"
    "errors"
    "fmt"
    "html/template"
    "io"
    "log"
    "net"
    "net/http"
    "reflect"
    "sort"
    "strings"
    "sync"
    "unicode"
    "unicode/utf8"
)

const (
    // Defaults used by HandleHTTP
    DefaultRPCPath   = "/_goRPC_"
    DefaultDebugPath = "/debug/rpc"
)

// DefaultServer is the default instance of *Server.

// DefaultServer是默认的*Server实例。
var DefaultServer = NewServer()

var ErrShutdown = errors.New("connection is shut down")

// Call represents an active RPC.

// Call 代表一个活跃的RPC
type Call struct {
    ServiceMethod string      // The name of the service and method to call.
    Args          interface{} // The argument to the function (*struct).
    Reply         interface{} // The reply from the function (*struct).
    Error         error       // After completion, the error status.
    Done          chan *Call  // Strobes when call is complete.
}

// Client represents an RPC Client.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.

// Client代表一个RPC客户端。
// 一个客户端可以有多个调用，并且一个客户端可以被多个goroutine同时使用
type Client struct {
}

// A ClientCodec implements writing of RPC requests and
// reading of RPC responses for the client side of an RPC session.
// The client calls WriteRequest to write a request to the connection
// and calls ReadResponseHeader and ReadResponseBody in pairs
// to read responses.  The client calls Close when finished with the
// connection. ReadResponseBody may be called with a nil
// argument to force the body of the response to be read and then
// discarded.

// ClientCodec实现了客户端一方对RPC会话的写RPC请求，和读RPC回复功能。客户端调用
// WriterRequest 往连接中写RPC请求，同时调用ReadResponseHeader和ReadResponseBody
// 来读取RPC返回。当结束连接的时候 客户端调用Close。ReadResponseBody可以使用一个
// nil参数，来读取RPC回复，然后丢弃信息。
type ClientCodec interface {
    // WriteRequest must be safe for concurrent use by multiple goroutines.
    WriteRequest(*Request, interface{}) error
    ReadResponseHeader(*Response) error
    ReadResponseBody(interface{}) error

    Close() error
}

// Request is a header written before every RPC call.  It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.

// Request是在每个RPC调用之前使用的header。它是内部使用的，写在这里是为了调试用
// ，例如分析网络的流量等。
type Request struct {
    ServiceMethod string // format: "Service.Method"
    Seq           uint64 // sequence number chosen by client

}

// Response is a header written before every RPC return.  It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.

// Response是在每个RPC回复之前被写在头里面的。它是内部使用的，写在这里是为了调试
// 用，例如分析网络的流量等。
type Response struct {
    ServiceMethod string // echoes that of the Request
    Seq           uint64 // echoes that of the request
    Error         string // error, if any.

}

// Server represents an RPC Server.

// Server代表一个RPC服务。
type Server struct {
}

// A ServerCodec implements reading of RPC requests and writing of
// RPC responses for the server side of an RPC session.
// The server calls ReadRequestHeader and ReadRequestBody in pairs
// to read requests from the connection, and it calls WriteResponse to
// write a response back.  The server calls Close when finished with the
// connection. ReadRequestBody may be called with a nil
// argument to force the body of the request to be read and discarded.

// ServerCodec实现了为RPC会话提供读RPC请求和写PRC回复的服务端的方法。服务端调用
// ReadRequestHeader和ReadRequestBody来读取连接上的请求，然后调用WriteResponse来
// 写回复。服务端当结束连接的时候调用Close。ReadRequestBody可能会调用一个nil参数
// 来强迫 读取请求内容并忽略。
type ServerCodec interface {
    ReadRequestHeader(*Request) error
    ReadRequestBody(interface{}) error
    // WriteResponse must be safe for concurrent use by multiple goroutines.
    WriteResponse(*Response, interface{}) error

    Close() error
}

// ServerError represents an error that has been returned from
// the remote side of the RPC connection.

// ServerError 代表从远程RPC连接另一端返回的错误
type ServerError string

// Accept accepts connections on the listener and serves requests
// to DefaultServer for each incoming connection.
// Accept blocks; the caller typically invokes it in a go statement.

// Accept在连接上监听和服务请求，为每个连接调用DefaultServer。
// Accept是阻塞的，调用者一般是在go语句中调用。
func Accept(lis net.Listener)

// Dial connects to an RPC server at the specified network address.

// Dial根据指定的网络地址连接到一个RPC服务。
func Dial(network, address string) (*Client, error)

// DialHTTP connects to an HTTP RPC server at the specified network address
// listening on the default HTTP RPC path.

// DialHttp根据制定的网络地址，连接到一个HTTP RPC服务。并且在默认的HTTP RPC路径
// 进行监听。
func DialHTTP(network, address string) (*Client, error)

// DialHTTPPath connects to an HTTP RPC server
// at the specified network address and path.

// DialHTTPPATH根据制定的网络地址和路径连接到一个HTTP RPC服务。
func DialHTTPPath(network, address, path string) (*Client, error)

// HandleHTTP registers an HTTP handler for RPC messages to DefaultServer
// on DefaultRPCPath and a debugging handler on DefaultDebugPath.
// It is still necessary to invoke http.Serve(), typically in a go statement.

// HandleHTTP在DefaultRPCPath上为RPC消息注册了一个HTTP的处理器到DefaultServer上
// ，并且在 DefaultDebugPath上注册了一个debuggin处理器。 它仍然需要调用
// http.Serve()，一般是在go语句中。
func HandleHTTP()

// NewClient returns a new Client to handle requests to the set of services at
// the other end of the connection. It adds a buffer to the write side of the
// connection so the header and payload are sent as a unit.
func NewClient(conn io.ReadWriteCloser) *Client

// NewClientWithCodec is like NewClient but uses the specified codec to encode
// requests and decode responses.
func NewClientWithCodec(codec ClientCodec) *Client

// NewServer returns a new Server.

// NewServer返回一个新的Server
func NewServer() *Server

// Register publishes the receiver's methods in the DefaultServer.

// Register在DefaultServer中发布接收者的方法
func Register(rcvr interface{}) error

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.

// RegisterName就像Register，但是为类型使用自定义的名字而不是接收者定义的名字。
func RegisterName(name string, rcvr interface{}) error

// ServeCodec is like ServeConn but uses the specified codec to
// decode requests and encode responses.

// ServeCodec和ServeConn一样，但是使用特定的codec来解码请求，编码回复。
func ServeCodec(codec ServerCodec)

// ServeConn runs the DefaultServer on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
// ServeConn uses the gob wire format (see package gob) on the
// connection.  To use an alternate codec, use ServeCodec.

// ServeConn在单个连接上调用DefaultServer。 ServeConn阻塞，服务连接，直到客户端
// 关闭。 调用者一般在go语句中调用ServeConn。ServeConn在连接上使用gob格式（参考
// gob包）。 要使用自定义的编解码，使用ServeCodec.
func ServeConn(conn io.ReadWriteCloser)

// ServeRequest is like ServeCodec but synchronously serves a single request.
// It does not close the codec upon completion.

// ServeRequest和ServeCodec相似，但是同步地服务单个请求。
// 它直到完成了才关闭codec。
func ServeRequest(codec ServerCodec) error

// Call invokes the named function, waits for it to complete, and returns its
// error status.

// Call调用方法的名字，等待它完成，然后返回成功或失败的error状态。
func (*Client) Call(serviceMethod string, args interface{}, reply interface{}) error

func (*Client) Close() error

// Go invokes the function asynchronously. It returns the Call structure
// representing the invocation. The done channel will signal when the call is
// complete by returning the same Call object. If done is nil, Go will allocate
// a new channel. If non-nil, done must be buffered or Go will deliberately
// crash.

// Go能异步调用功能。它返回Call结构来代表回调。当调用完成，返回相同的Call对象，
// done channel就会获取到 信息。如果done是空的话，Go就会分配一个新的channel。如
// 果非空的话，done必须缓冲起来，或者Go会立即崩溃。
func (*Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call

// Accept accepts connections on the listener and serves requests
// for each incoming connection. Accept blocks until the listener
// returns a non-nil error. The caller typically invokes Accept in a
// go statement.

// Accept接收连接，为每个连接监听和服务请求。Accept是阻塞的，调用者一般在go语句
// 中使用它。
func (*Server) Accept(lis net.Listener)

// HandleHTTP registers an HTTP handler for RPC messages on rpcPath,
// and a debugging handler on debugPath.
// It is still necessary to invoke http.Serve(), typically in a go statement.

// HandleHTTP在rpcPath上为RPC消息注册一个HTTP处理器，并在debugPath注册一个
// debugging处理器。 它仍然需要调用http.Serve()，一般是在go语句中使用。
func (*Server) HandleHTTP(rpcPath, debugPath string)

// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//     - exported method of exported type
//     - two arguments, both of exported type
//     - the second argument is a pointer
//     - one return value, of type error
// It returns an error if the receiver is not an exported type or has
// no suitable methods. It also logs the error using package log.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.

// Register发布服务器的一系列方法，接受器必须满足这几个条件：
//
//     - 对外可见的方法
//     - 两个参数，都指向对外可见的结构
//     - 一个error类型返回值
//
// 如果接收者不是一个对外可见的类型，或者没有任何方法，或者没有满足条件的方法，
// 都会返回error。 它也会使用log包来记录错误。客户端进入每个方法使用字符串格式形
// 如“Type.Method”, 这里Type是接收者的具体的类型。
func (*Server) Register(rcvr interface{}) error

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.

// RegisterName像Register，但是为type使用提供的名字，而不是使用receivers的具体类
// 型。
func (*Server) RegisterName(name string, rcvr interface{}) error

// ServeCodec is like ServeConn but uses the specified codec to
// decode requests and encode responses.

// ServerCodec和ServeConn相似，但是使用自定义的编解码器来解码请求和编码回复。
func (*Server) ServeCodec(codec ServerCodec)

// ServeConn runs the server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
// ServeConn uses the gob wire format (see package gob) on the
// connection.  To use an alternate codec, use ServeCodec.

// ServeConn在单个连接上跑server。 ServeConn阻塞，知道客户端关闭之后才继续服务其
// 他连接。 调用者一般在go语句中调用ServeConn。 ServeConn在连接传输的时候使用gob
// 格式（参考gob包）。可以使用自定义编码器，ServeCodec。
func (*Server) ServeConn(conn io.ReadWriteCloser)

// ServeHTTP implements an http.Handler that answers RPC requests.

// ServeHTTP实现了http.Handler，并且回复RPC请求。
func (*Server) ServeHTTP(w http.ResponseWriter, req *http.Request)

// ServeRequest is like ServeCodec but synchronously serves a single request.
// It does not close the codec upon completion.

// ServerRequest和ServeCodec相似，但是是同步地服务单个请求。
// 它在结束的时候不会关闭编解码器。
func (*Server) ServeRequest(codec ServerCodec) error

func (ServerError) Error() string

