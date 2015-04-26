// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package signal implements access to incoming signals.

// signal包实现了对输入信号的访问。
package signal

// Notify causes package signal to relay incoming signals to c. If no signals are
// listed, all incoming signals will be relayed to c. Otherwise, just the listed
// signals will.
//
// Package signal will not block sending to c: the caller must ensure that c has
// sufficient buffer space to keep up with the expected signal rate. For a channel
// used for notification of just one signal value, a buffer of size 1 is
// sufficient.
//
// It is allowed to call Notify multiple times with the same channel: each call
// expands the set of signals sent to that channel. The only way to remove signals
// from the set is to call Stop.
//
// It is allowed to call Notify multiple times with different channels and the same
// signals: each channel receives copies of incoming signals independently.

// Notify函数让signal包将输入信号转发到c。如果没有列出要传递的信号，会将所有输入信号传递到c；否则只传递列出的输入信号。
//
// signal包不会为了向c发送信息而阻塞（就是说如果发送时c阻塞了，signal包会直接放弃）：调用者应该保证c有足够的缓存空间可以跟上期望的信号频率。对使用单一信号用于通知的通道，缓存为1就足够了。
//
// 可以使用同一通道多次调用Notify：每一次都会扩展该通道接收的信号集。唯一从信号集去除信号的方法是调用Stop。可以使用同一信号和不同通道多次调用Notify：每一个通道都会独立接收到该信号的一个拷贝。
func Notify(c chan<- os.Signal, sig ...os.Signal)

// Stop causes package signal to stop relaying incoming signals to c. It undoes the
// effect of all prior calls to Notify using c. When Stop returns, it is guaranteed
// that c will receive no more signals.

// Stop函数让signal包停止向c转发信号。它会取消之前使用c调用的所有Notify的效果。当Stop返回后，会保证c不再接收到任何信号。
func Stop(c chan<- os.Signal)
