// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package signal implements access to incoming signals.
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
func Notify(c chan<- os.Signal, sig ...os.Signal)

// Stop causes package signal to stop relaying incoming signals to c. It undoes the
// effect of all prior calls to Notify using c. When Stop returns, it is guaranteed
// that c will receive no more signals.
func Stop(c chan<- os.Signal)
