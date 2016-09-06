// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package ring implements operations on circular lists.

// ring实现了环形链表的操作。
package ring

// A Ring is an element of a circular list, or ring.
// Rings do not have a beginning or end; a pointer to any ring element
// serves as reference to the entire ring. Empty rings are represented
// as nil Ring pointers. The zero value for a Ring is a one-element
// ring with a nil Value.

// Ring类型代表环形链表的一个元素，同时也代表链表本身。环形链表没有头尾；指向环
// 形链表任一元素的指针都可以作为整个环形链表看待。Ring零值是具有一个（Value字段
// 为nil的）元素的链表。
type Ring struct {
	Value interface{} // for use by client; untouched by this library
}

// New creates a ring of n elements.

// New创建一个具有n个元素的环形链表。
func New(n int) *Ring

// Do calls function f on each element of the ring, in forward order.
// The behavior of Do is undefined if f changes *r.

// 对链表的每一个元素都执行f（正向顺序），注意如果f改变了*r，Do的行为是未定义的
// 。
func (r *Ring) Do(f func(interface{}))

// Len computes the number of elements in ring r.
// It executes in time proportional to the number of elements.

// Len返回环形链表中的元素个数，复杂度O(n)。
func (r *Ring) Len() int

// Link connects ring r with ring s such that r.Next()
// becomes s and returns the original value for r.Next().
// r must not be empty.
//
// If r and s point to the same ring, linking
// them removes the elements between r and s from the ring.
// The removed elements form a subring and the result is a
// reference to that subring (if no elements were removed,
// the result is still the original value for r.Next(),
// and not nil).
//
// If r and s point to different rings, linking
// them creates a single ring with the elements of s inserted
// after r. The result points to the element following the
// last element of s after insertion.

// Link连接r和s，并返回r原本的后继元素r.Next()。r不能为空。
//
// 如果r和s指向同一个环形链表，则会删除掉r和s之间的元素，删掉的元素构成一个子链
// 表，返回指向该子链表的指针（r的原后继元素）；如果没有删除元素，则仍然返回r的
// 原后继元素，而不是nil。如果r和s指向不同的链表，将创建一个单独的链表，将s指向
// 的链表插入r后面，返回s原最后一个元素后面的元素（即r的原后继元素）。
func (r *Ring) Link(s *Ring) *Ring

// Move moves n % r.Len() elements backward (n < 0) or forward (n >= 0)
// in the ring and returns that ring element. r must not be empty.

// 返回移动n个位置（n>=0向前移动，n<0向后移动）后的元素，r不能为空。
func (r *Ring) Move(n int) *Ring

// Next returns the next ring element. r must not be empty.

// 返回后一个元素，r不能为空。
func (r *Ring) Next() *Ring

// Prev returns the previous ring element. r must not be empty.

// 返回前一个元素，r不能为空。
func (r *Ring) Prev() *Ring

// Unlink removes n % r.Len() elements from the ring r, starting
// at r.Next(). If n % r.Len() == 0, r remains unchanged.
// The result is the removed subring. r must not be empty.

// 删除链表中n % r.Len()个元素，从r.Next()开始删除。如果n % r.Len() ==
// 0，不修改r。返回删除的元素构成的链表，r不能为空。
func (r *Ring) Unlink(n int) *Ring

