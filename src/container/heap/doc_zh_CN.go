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
// A heap is a common way to implement a priority queue. To build a priority queue,
// implement the Heap interface with the (negative) priority as the ordering for
// the Less method, so Push adds items while Pop removes the highest-priority item
// from the queue. The Examples include such an implementation; the file
// example_pq_test.go has the complete source.
package heap

// Fix re-establishes the heap ordering after the element at index i has changed
// its value. Changing the value of the element at index i and then calling Fix is
// equivalent to, but less expensive than, calling Remove(h, i) followed by a Push
// of the new value. The complexity is O(log(n)) where n = h.Len().
func Fix(h Interface, i int)

// A heap must be initialized before any of the heap operations can be used. Init
// is idempotent with respect to the heap invariants and may be called whenever the
// heap invariants may have been invalidated. Its complexity is O(n) where n =
// h.Len().
func Init(h Interface)

// Pop removes the minimum element (according to Less) from the heap and returns
// it. The complexity is O(log(n)) where n = h.Len(). It is equivalent to Remove(h,
// 0).
func Pop(h Interface) interface{}

// Push pushes the element x onto the heap. The complexity is O(log(n)) where n =
// h.Len().
func Push(h Interface, x interface{})

// Remove removes the element at index i from the heap. The complexity is O(log(n))
// where n = h.Len().
func Remove(h Interface, i int) interface{}

// Any type that implements heap.Interface may be used as a min-heap with the
// following invariants (established after Init has been called or if the data is
// empty or sorted):
//
//	!h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()
//
// Note that Push and Pop in this interface are for package heap's implementation
// to call. To add and remove things from the heap, use heap.Push and heap.Pop.
type Interface interface {
	sort.Interface
	Push(x interface{}) // add x as element Len()
	Pop() interface{}   // remove and return element Len() - 1.
}
