// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package sync provides basic synchronization primitives such as mutual exclusion
// locks. Other than the Once and WaitGroup types, most are intended for use by
// low-level library routines. Higher-level synchronization is better done via
// channels and communication.
//
// Values containing the types defined in this package should not be copied.

// sync 包提供了互斥锁这类的基本的同步原语. 除 Once 和 WaitGroup
// 之外的类型大多用于底层库的例程。
// 更高级的同步操作通过信道与通信进行。
//
// 在此包中定义的类型中包含的值不应当被复制。
package sync

// Cond implements a condition variable, a rendezvous point for goroutines waiting
// for or announcing the occurrence of an event.
//
// Each Cond has an associated Locker L (often a *Mutex or *RWMutex), which must be
// held when changing the condition and when calling the Wait method.
//
// A Cond can be created as part of other structures. A Cond must not be copied
// after first use.

// Cond
// 实现了条件变量，即Go程等待的汇合点或宣布一个事件的发生。
//
// 每个 Cond 都有一个与其相关联的 Locker L（一般是 *Mutex 或 *RWMutex），
// 在改变该条件或调用 Wait 方法时，它必须保持不变。
type Cond struct {
	// L is held while observing or changing the condition
	L Locker
	// contains filtered or unexported fields
}

// NewCond returns a new Cond with Locker l.

// NewCond 用 Locker l 返回一个新的 Cond。
func NewCond(l Locker) *Cond

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L during the call.

// Broadcast 唤醒所有等待 c 的Go程。
//
// during the call.在调用其间可以保存 c.L，但并没有必要。
func (c *Cond) Broadcast()

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L during the call.

// Signal 用于唤醒等待 c 的Go程，如果有的话。
//
// during the call.在调用其间可以保存 c.L，但并没有必要。
func (c *Cond) Signal()

// Wait atomically unlocks c.L and suspends execution of the calling goroutine.
// After later resuming execution, Wait locks c.L before returning. Unlike in other
// systems, Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller typically cannot
// assume that the condition is true when Wait returns. Instead, the caller should
// Wait in a loop:
//
//	c.L.Lock()
//	for !condition() {
//	    c.Wait()
//	}
//	... make use of condition ...
//	c.L.Unlock()

// Wait 原子性地解锁 c.L
// 并挂起调用的Go程的执行。不像其它的系统那样，Wait 不会返回，除非它被 Broadcast 或 Signal 唤醒。
//
// 由于 Wait 第一次恢复时 c.L
// 并未锁定，因此调用者一般不能假定 Wait 返回时条件为真。 取而代之，调用者应当把 Wait 放入循环中：
//
//	c.L.Lock()
//	for !condition() {
//	    c.Wait()
//	}
//	... 使用 condition ...
//	c.L.Unlock()
func (c *Cond) Wait()

// A Locker represents an object that can be locked and unlocked.

// Locker 表示可被锁定并解锁的对象。
type Locker interface {
	Lock()
	Unlock()
}

// A Mutex is a mutual exclusion lock. Mutexes can be created as part of other
// structures; the zero value for a Mutex is an unlocked mutex.

// Mutex 是一个互斥锁。 Mutex
// 可作为其它结构的一部分来创建；Mutex 的零值即为已解锁的互斥体。
type Mutex struct {
	// contains filtered or unexported fields
}

// Lock locks m. If the lock is already in use, the calling goroutine blocks until
// the mutex is available.

// Lock 用于锁定 m。
// 若该锁正在使用，调用的Go程就会阻塞，直到该互斥体可用。
func (m *Mutex) Lock()

// Unlock unlocks m. It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine. It is allowed for
// one goroutine to lock a Mutex and then arrange for another goroutine to unlock
// it.

// Unlock 用于解锁 m。 若 m 在进入 Unlock
// 前并未锁定，就会引发一个运行时错误。
//
// 已锁定的 Mutex
// 并不与特定的Go程相关联，这样便可让一个Go程锁定
// Mutex，然后安排其它Go程来解锁。
func (m *Mutex) Unlock()

// Once is an object that will perform exactly one action.

// Once 是只执行一个动作的对象。
type Once struct {
	// contains filtered or unexported fields
}

// Do calls the function f if and only if Do is being called for the first time for
// this instance of Once. In other words, given
//
//	var once Once
//
// if once.Do(f) is called multiple times, only the first call will invoke f, even
// if f has a different value in each invocation. A new instance of Once is
// required for each function to execute.
//
// Do is intended for initialization that must be run exactly once. Since f is
// niladic, it may be necessary to use a function literal to capture the arguments
// to a function to be invoked by Do:
//
//	config.once.Do(func() { config.init(filename) })
//
// Because no call to Do returns until the one call to f returns, if f causes Do to
// be called, it will deadlock.
//
// If f panics, Do considers it to have returned; future calls of Do return without
// calling f.

// Do
// 方法当且仅当连同此接收者第一次被调用是才执行函数 f。
//
//	var once Once
//
// if once.Do(f) is called multiple times, only the first call will invoke f, even
// if f has a different value in each invocation. A new instance of Once is
// required for each function to execute. 若 once.Do(f)
// 被调用多次，即使每一次请求的 f
// 值都不同，也只有第一次调用会请求 f。 Once
// 的新实例需要为每一个函数所执行。
//
// Do 用于必须刚好运行一次的初始化。由于 f
// 是函数，它可能需要使用函数字面来为 Do 所请求的函数捕获实参：
//
//	config.once.Do(func() { config.init(filename) })
//
// 由于 f 的调用返回之前没有 Do 的调用会返回，因此若 f 引起了 Do 的调用，它就会死锁。
//
// 若 f 发生派错（panic），Do 会考虑它是否已返回；将来对 Do
// 的调用会直接返回而不调用 f。
func (o *Once) Do(f func())

// A Pool is a set of temporary objects that may be individually saved and
// retrieved.
//
// Any item stored in the Pool may be removed automatically at any time without
// notification. If the Pool holds the only reference when this happens, the item
// might be deallocated.
//
// A Pool is safe for use by multiple goroutines simultaneously.
//
// Pool's purpose is to cache allocated but unused items for later reuse, relieving
// pressure on the garbage collector. That is, it makes it easy to build efficient,
// thread-safe free lists. However, it is not suitable for all free lists.
//
// An appropriate use of a Pool is to manage a group of temporary items silently
// shared among and potentially reused by concurrent independent clients of a
// package. Pool provides a way to amortize allocation overhead across many
// clients.
//
// An example of good use of a Pool is in the fmt package, which maintains a
// dynamically-sized store of temporary output buffers. The store scales under load
// (when many goroutines are actively printing) and shrinks when quiescent.
//
// On the other hand, a free list maintained as part of a short-lived object is not
// a suitable use for a Pool, since the overhead does not amortize well in that
// scenario. It is more efficient to have such objects implement their own free
// list.
type Pool struct {

	// New optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	New func() interface{}
	// contains filtered or unexported fields
}

// Get selects an arbitrary item from the Pool, removes it from the Pool, and
// returns it to the caller. Get may choose to ignore the pool and treat it as
// empty. Callers should not assume any relation between values passed to Put and
// the values returned by Get.
//
// If Get would otherwise return nil and p.New is non-nil, Get returns the result
// of calling p.New.
func (p *Pool) Get() interface{}

// Put adds x to the pool.
func (p *Pool) Put(x interface{})

// An RWMutex is a reader/writer mutual exclusion lock. The lock can be held by an
// arbitrary number of readers or a single writer. RWMutexes can be created as part
// of other structures; the zero value for a RWMutex is an unlocked mutex.

// RWMutex 是一个读写互斥锁。
// 该说可被任意多个读取器或单个写入器所持有。RWMutex
// 可作为其它结构的一部分来创建； RWMutex 的零值即为已解锁的互斥体。
type RWMutex struct {
	// contains filtered or unexported fields
}

// Lock locks rw for writing. If the lock is already locked for reading or writing,
// Lock blocks until the lock is available. To ensure that the lock eventually
// becomes available, a blocked Lock call excludes new readers from acquiring the
// lock.

// Lock 为 rw 的写入将其锁定。
// 若该锁已经为读取或写入而锁定，Lock 就会阻塞直到该锁可用。
// 为确保该锁最终可用，已阻塞的 Lock
// 调用会从获得的锁中排除新的读取器。
func (rw *RWMutex) Lock()

// RLock locks rw for reading.

// RLock 为 rw 的读取将其锁定。
func (rw *RWMutex) RLock()

// RLocker returns a Locker interface that implements the Lock and Unlock methods
// by calling rw.RLock and rw.RUnlock.

// RLocker 返回一个 Locker 接口，该接口通过调用 rw.RLock 和 rw.RUnlock 实现了 Lock 和 Unlock
// 方法。
func (rw *RWMutex) RLocker() Locker

// RUnlock undoes a single RLock call; it does not affect other simultaneous
// readers. It is a run-time error if rw is not locked for reading on entry to
// RUnlock.

// RUnlock 撤销单次 RLock
// 调用，它对于其它同时存在的读取器则没有效果。 若 rw 并没有为读取而锁定，调用 RUnlock
// 就会引发一个运行时错误。
func (rw *RWMutex) RUnlock()

// Unlock unlocks rw for writing. It is a run-time error if rw is not locked for
// writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular goroutine.
// One goroutine may RLock (Lock) an RWMutex and then arrange for another goroutine
// to RUnlock (Unlock) it.

// Unlock 为 rw 的写入将其解锁。 若 rw 并没有为写入而锁定，调用 Unlock
// 就会引发一个运行时错误。
//
// As with Mutexes, a locked RWMutex is not associated with a particular goroutine.
// One goroutine may RLock (Lock) an RWMutex and then arrange for another goroutine
// to RUnlock (Unlock) it. 正如 Mutex 一样，已锁定的 RWMutex
// 并不与特定的Go程相关联。一个Go程可 RLock（Lock）一个 RWMutex，然后安排其它Go程来
// RUnlock（Unlock）它。
func (rw *RWMutex) Unlock()

// A WaitGroup waits for a collection of goroutines to finish. The main goroutine
// calls Add to set the number of goroutines to wait for. Then each of the
// goroutines runs and calls Done when finished. At the same time, Wait can be used
// to block until all goroutines have finished.

// WaitGroup 等待一组Go程的结束。 主Go程调用 Add
// 来设置等待的Go程数。然后该组中的每个Go程都会运行，并在结束时调用 Done。同时，Wait
// 可被用于阻塞，直到所有Go程都结束。
type WaitGroup struct {
	// contains filtered or unexported fields
}

// Add adds delta, which may be negative, to the WaitGroup counter. If the counter
// becomes zero, all goroutines blocked on Wait are released. If the counter goes
// negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero must
// happen before a Wait. Calls with a negative delta, or calls with a positive
// delta that start when the counter is greater than zero, may happen at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for. See the WaitGroup
// example.

// Add 添加 delta，对于 WaitGroup 的 counter 来说，它可能为负数。 若 counter 变为零，在 Wait()
// 被释放后所有Go程就会阻塞。 若 counter 变为负数，Add 就会引发Panic。
//
// 注意，当 counter 为零时，用正整数的 delta 调用它必须发生在调用 Wait 之前。 用负整数的 delta
// 调用它，或在 counter 大于零时开始用正整数的 delta 调用它，
// 那么它可以在任何时候发生。 一般来说，这意味着对 Add
// 的调用应当在该语句创建Go程，或等待其它事件之前执行。 具体见 WaitGroup 的示例。
func (wg *WaitGroup) Add(delta int)

// Done decrements the WaitGroup counter.

// Done 递减 WaitGroup 的 counter。
func (wg *WaitGroup) Done()

// Wait blocks until the WaitGroup counter is zero.

// Wait 阻塞 WaitGroup 直到其 counter 为零。
func (wg *WaitGroup) Wait()
