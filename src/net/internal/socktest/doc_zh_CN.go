// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package socktest provides utilities for socket testing.

// Package socktest provides utilities for socket testing.
package socktest // import "net/internal/socktest"

import (
    "fmt"
    "sync"
    "syscall"
)

const (
    FilterSocket        FilterType = iota // for Socket
    FilterConnect                         // for Connect or ConnectEx
    FilterListen                          // for Listen
    FilterAccept                          // for Accept or Accept4
    FilterGetsockoptInt                   // for GetsockoptInt
    FilterClose                           // for Close or Closesocket
)

// An AfterFilter represents a socket system call filter after an
// execution of a system call.
//
// It will only be executed after a system call for a socket that has
// an entry in internal table.
// If the filter returns a non-nil error, the system call function
// returns the non-nil error.
type AfterFilter func(*Status) error

// A Cookie represents a 3-tuple of a socket; address family, socket
// type and protocol number.
type Cookie uint64

// A Filter represents a socket system call filter.
//
// It will only be executed before a system call for a socket that has
// an entry in internal table.
// If the filter returns a non-nil error, the execution of system call
// will be canceled and the system call function returns the non-nil
// error.
// It can return a non-nil AfterFilter for filtering after the
// execution of the system call.
type Filter func(*Status) (AfterFilter, error)

// A FilterType represents a filter type.
type FilterType int

// Sockets maps a socket descriptor to the status of socket.
type Sockets map[int]Status

// Sockets maps a socket descriptor to the status of socket.
type Sockets map[int]Status

// Sockets maps a socket descriptor to the status of socket.
type Sockets map[syscall.Handle]Status

// A Stat represents a per-cookie socket statistics.
type Stat struct {
    Family   int // address family
    Type     int // socket type
    Protocol int // protocol number

    Opened    uint64 // number of sockets opened
    Connected uint64 // number of sockets connected
    Listened  uint64 // number of sockets listened
    Accepted  uint64 // number of sockets accepted
    Closed    uint64 // number of sockets closed

    OpenFailed    uint64 // number of sockets open failed
    ConnectFailed uint64 // number of sockets connect failed
    ListenFailed  uint64 // number of sockets listen failed
    AcceptFailed  uint64 // number of sockets accept failed
    CloseFailed   uint64 // number of sockets close failed
}

// A Status represents the status of a socket.
type Status struct {
    Cookie    Cookie
    Err       error // error status of socket system call
    SocketErr error // error status of socket by SO_ERROR
}

// A Switch represents a callpath point switch for socket system
// calls.
type Switch struct {
    once sync.Once

    fmu   sync.RWMutex
    fltab map[FilterType]Filter

    smu   sync.RWMutex
    sotab Sockets
    stats stats
}

// Accept wraps syscall.Accept.
func (*Switch) Accept(s int) (ns int, sa syscall.Sockaddr, err error)

// Accept4 wraps syscall.Accept4.
func (*Switch) Accept4(s, flags int) (ns int, sa syscall.Sockaddr, err error)

// Close wraps syscall.Close.
func (*Switch) Close(s int) (err error)

// Closesocket wraps syscall.Closesocket.
func (*Switch) Closesocket(s syscall.Handle) (err error)

// Connect wraps syscall.Connect.
func (*Switch) Connect(s int, sa syscall.Sockaddr) (err error)

// ConnectEx wraps syscall.ConnectEx.
func (*Switch) ConnectEx(s syscall.Handle, sa syscall.Sockaddr, b *byte, n uint32, nwr *uint32, o *syscall.Overlapped) (err error)

// GetsockoptInt wraps syscall.GetsockoptInt.
func (*Switch) GetsockoptInt(s, level, opt int) (soerr int, err error)

// Listen wraps syscall.Listen.
func (*Switch) Listen(s, backlog int) (err error)

// Set deploys the socket system call filter f for the filter type t.
func (*Switch) Set(t FilterType, f Filter)

// Socket wraps syscall.Socket.
func (*Switch) Socket(family, sotype, proto int) (s int, err error)

// Sockets returns mappings of socket descriptor to socket status.
func (*Switch) Sockets() Sockets

// Stats returns a list of per-cookie socket statistics.
func (*Switch) Stats() []Stat

// Family returns an address family.
func (Cookie) Family() int

// Protocol returns a protocol number.
func (Cookie) Protocol() int

// Type returns a socket type.
func (Cookie) Type() int

func (Stat) String() string

func (Status) String() string

