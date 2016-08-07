// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package heap provides heap operations for any type that implements
// heap.Interface. A heap is a tree with the property that each node is the
// minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, implement the Heap interface with the (negative) priority as the
// ordering for the Less method, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete source.

// heap包提供了对任意类型（实现了heap.Interface接口）的堆操作。（最小）堆是具有
// “每个节点都是以其为根的子树中最小值”属性的树。
//
// 树的最小元素为其根元素，索引0的位置。
//
// heap是常用的实现优先队列的方法。要创建一个优先队列，实现一个具有使用（负的）
// 优先级作为比较的依据的Less方法的Heap接口，如此一来可用Push添加项目而用Pop取出
// 队列最高优先级的项目。
package heap

import "sort"

// Any type that implements heap.Interface may be used as a min-heap with the
// following invariants (established after Init has been called or if the data
// is empty or sorted):
//
//     !h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()
//
// Note that Push and Pop in this interface are for package heap's
// implementation to call. To add and remove things from the heap, use heap.Push
// and heap.Pop.

// 任何实现了本接口的类型都可以用于构建最小堆。最小堆可以通过heap.Init建立，数据
// 是递增顺序或者空的话也是最小堆。最小堆的约束条件是：
//
//     !h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()
//
// 注意接口的Push和Pop方法是供heap包调用的，请使用heap.Push和heap.Pop来向一个堆
// 添加或者删除元素。
type Interface interface {
    sort.Interface
    Push(x interface{}) // add x as element Len()
    Pop() interface{}   // remove and return element Len() - 1.
}

// Fix re-establishes the heap ordering after the element at index i has changed
// its value. Changing the value of the element at index i and then calling Fix
// is equivalent to, but less expensive than, calling Remove(h, i) followed by a
// Push of the new value. The complexity is O(log(n)) where n = h.Len().

// 在修改第i个元素后，调用本函数修复堆，比删除第i个元素后插入新元素更有效率。
//
// 复杂度O(log(n))，其中n等于h.Len()。
func Fix(h Interface, i int)

// A heap must be initialized before any of the heap operations
// can be used. Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// Its complexity is O(n) where n = h.Len().

// 一个堆在使用任何堆操作之前应先初始化。Init函数对于堆的约束性是幂等的（多次执
// 行无意义），并可能在任何时候堆的约束性被破坏时被调用。本函数复杂度为O(n)，其
// 中n等于h.Len()。
func Init(h Interface)

// Pop removes the minimum element (according to Less) from the heap
// and returns it. The complexity is O(log(n)) where n = h.Len().
// It is equivalent to Remove(h, 0).

// 删除并返回堆h中的最小元素（不影响约束性）。复杂度O(log(n))，其中n等于h.Len()
// 。等价于Remove(h, 0)。
func Pop(h Interface) interface{}

// Push pushes the element x onto the heap. The complexity is
// O(log(n)) where n = h.Len().

// 向堆h中插入元素x，并保持堆的约束性。复杂度O(log(n))，其中n等于h.Len()。
func Push(h Interface, x interface{})

// Remove removes the element at index i from the heap.
// The complexity is O(log(n)) where n = h.Len().

// 删除堆中的第i个元素，并保持堆的约束性。复杂度O(log(n))，其中n等于h.Len()。
func Remove(h Interface, i int) interface{}

