// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package jsonrpc implements a JSON-RPC ClientCodec and ServerCodec
// for the rpc package.

// jsonrpc 包使用了rpc的包实现了一个JSON-RPC的客户端解码器和服务端的解码器.
package jsonrpc

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net"
    "net/rpc"
    "sync"
)

// Dial connects to a JSON-RPC server at the specified network address.

// Dial在指定的网络地址上，连接了一个JSON-RPC服务
func Dial(network, address string) (*rpc.Client, error)

// NewClient returns a new rpc.Client to handle requests to the
// set of services at the other end of the connection.

// NewClient返回新的rpc.Client，用于连接的服务器一端来进行rpc服务。
func NewClient(conn io.ReadWriteCloser) *rpc.Client

// NewClientCodec returns a new rpc.ClientCodec using JSON-RPC on conn.

// NewClientCodec在连接中使用JSON-RPC返回一个新的rpc.ClientCodec
func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec

// NewServerCodec returns a new rpc.ServerCodec using JSON-RPC on conn.

// NewServerCodec在连接中使用JSON-RPC返回一个新的rpc.ServerCodec
func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec

// ServeConn runs the JSON-RPC server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.

// ServeConn在一个连接上运行启动一个JSON-RPC。
// ServeConn是阻塞的，直到客户端关闭都服务这个连接。
// 调用者一般是在go语句中调用ServeConn的。
func ServeConn(conn io.ReadWriteCloser)

