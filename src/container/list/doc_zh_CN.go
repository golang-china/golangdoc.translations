// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package list implements a doubly linked list.
//
// To iterate over a list (where l is a *List):
//     for e := l.Front(); e != nil; e = e.Next() {
//         // do something with e.Value
//     }

// Package list implements a doubly linked list.
//
// To iterate over a list (where l is a *List):
//     for e := l.Front(); e != nil; e = e.Next() {
//         // do something with e.Value
//     }
package list

// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value interface{}
}


// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}


// New returns an initialized list.
func New() *List

// Next returns the next list element or nil.
func (*Element) Next() *Element

// Prev returns the previous list element or nil.
func (*Element) Prev() *Element

// Back returns the last element of list l or nil.
func (*List) Back() *Element

// Front returns the first element of list l or nil.
func (*List) Front() *Element

// Init initializes or clears list l.
func (*List) Init() *List

// InsertAfter inserts a new element e with value v immediately after mark and
// returns e. If mark is not an element of l, the list is not modified.
func (*List) InsertAfter(v interface{}, mark *Element) *Element

// InsertBefore inserts a new element e with value v immediately before mark and
// returns e. If mark is not an element of l, the list is not modified.
func (*List) InsertBefore(v interface{}, mark *Element) *Element

// Len returns the number of elements of list l.
// The complexity is O(1).
func (*List) Len() int

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
func (*List) MoveAfter(e, mark *Element)

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
func (*List) MoveBefore(e, mark *Element)

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
func (*List) MoveToBack(e *Element)

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
func (*List) MoveToFront(e *Element)

// PushBack inserts a new element e with value v at the back of list l and
// returns e.
func (*List) PushBack(v interface{}) *Element

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same.
func (*List) PushBackList(other *List)

// PushFront inserts a new element e with value v at the front of list l and
// returns e.
func (*List) PushFront(v interface{}) *Element

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same.
func (*List) PushFrontList(other *List)

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
func (*List) Remove(e *Element) interface{}

