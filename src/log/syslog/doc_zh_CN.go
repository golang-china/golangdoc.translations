// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package syslog provides a simple interface to the system log
// service. It can send messages to the syslog daemon using UNIX
// domain sockets, UDP or TCP.
//
// Only one call to Dial is necessary. On write failures,
// the syslog client will attempt to reconnect to the server
// and write again.
//
// The syslog package is frozen and not accepting new features.
// Some external packages provide more functionality. See:
//
//   https://godoc.org/?q=syslog

// Package syslog provides a simple interface to the system log service. It can
// send messages to the syslog daemon using UNIX domain sockets, UDP or TCP.
//
// Only one call to Dial is necessary. On write failures, the syslog client will
// attempt to reconnect to the server and write again.
//
// Package syslog provides a simple interface to the system log service.
//
// Package syslog provides a simple interface to the system log service.
package syslog

