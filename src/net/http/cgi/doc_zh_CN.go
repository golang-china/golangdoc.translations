// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package cgi implements CGI (Common Gateway Interface) as specified
// in RFC 3875.
//
// Note that using CGI means starting a new process to handle each
// request, which is typically less efficient than using a
// long-running server. This package is intended primarily for
// compatibility with existing systems.

// cgi 包实现了RFC3875协议描述的CGI（公共网关接口）.
//
// 使用CGI就意味开启一个新进程来处理每个请求，这种方法当然比持久运行的服务进程的
// 方式低效些。 这个包主要用来和现有的web系统进行交互。
package cgi

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// Handler runs an executable in a subprocess with a CGI environment.

// Handler会在子进程中创建一个CGI环境来运行可执行程序。
type Handler struct {
	Path string // path to the CGI executable  // CGI可执行脚本的地址
	Root string // root URI prefix of handler or empty for "/"  // URI的前缀ROOT部分，如果为空的话就代表“/”

	// Dir specifies the CGI executable's working directory.
	// If Dir is empty, the base directory of Path is used.
	// If Path has no base directory, the current working
	// directory is used.

	// Dir说明CGI可执行脚本运行所在的工作路径。
	// 如果Dir为空，Path参数指的文件所在的路径就会被使用。
	// 如果Path没有路径，那么Dir就会使用当前的执行目录。
	Dir        string
	Env        []string    // extra environment variables to set, if any, as "key=value"  // 需要额外设置的环境变量，如果有的话，形如“key=value”
	InheritEnv []string    // environment variables to inherit from host, as "key" // 需要继承自宿主的环境变量，形如“key”
	Logger     *log.Logger // optional log for errors or nil to use log.Print  // 可选。错误的日志处理器，如果是nil的话就默认使用log.Print
	Args       []string    // optional arguments to pass to child process  // 可选。给子进程传递的附加参数。
	Stderr     io.Writer   // optional stderr for the child process; nil means os.Stderr

	// PathLocationHandler specifies the root http Handler that
	// should handle internal redirects when the CGI process
	// returns a Location header value starting with a "/", as
	// specified in RFC 3875 § 6.3.2. This will likely be
	// http.DefaultServeMux.
	//
	// If nil, a CGI response with a local URI path is instead sent
	// back to the client and not redirected internally.

	// PathLocationHandler是http Handler，它说明的是在CGI进程返回的Location
	// header信息是以“/”开头的时候（location是在RFC 3875 § 6.3.2）， 根目录应
	// 当如何处理内部的重定向规则。这个值可以是http.DefaultServeMux。
	//
	// 如果为空，一个带有本地URI路径的CGI回复会立刻返回给客户端，并且没有进行任
	// 何的内部重定向。
	PathLocationHandler http.Handler
}

// Request returns the HTTP request as represented in the current
// environment. This assumes the current program is being run
// by a web server in a CGI environment.
// The returned Request's Body is populated, if applicable.

// Request()函数返回当前系统环境下的HTTP请求。这个函数假设当前的程序是跑在一个
// CGI环境下的WebServer中。 返回的Request的Body字段是可有可无的，如果有的话才会
// 返回回来，如果Body没有内容的话，这个字段就是空。
func Request() (*http.Request, error)

// RequestFromMap creates an http.Request from CGI variables.
// The returned Request's Body field is not populated.

// RequestFromMap从CGI的变量中提取出http.Request结构。
// 返回的Request的Body字段是不会为空的。
func RequestFromMap(params map[string]string) (*http.Request, error)

// Serve executes the provided Handler on the currently active CGI
// request, if any. If there's no current CGI environment
// an error is returned. The provided handler may be nil to use
// http.DefaultServeMux.

// Serve使用提供的Handler来处理当前的CGI请求。如果CGI环境配置不正确的话，会返回
// 一个error。 如果提供的Hanlder是nil的话，程序就会使用http.DefaultServeMux。
func Serve(handler http.Handler) error

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request)

