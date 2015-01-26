// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package list implements a doubly linked list.
//
// To iterate over a list (where l is a *List):
//
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}

// Package list implements a doubly linked
// list.
//
// To iterate over a list (where l is a
// *List):
//
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
package list

// Element is an element of a linked list.

// Element is an element of a linked list.
type Element struct {

	// The value stored with this element.
	Value interface{}
	// contains filtered or unexported fields
}

// Next returns the next list element or nil.

// Next returns the next list element or
// nil.
func (e *Element) Next() *Element

// Prev returns the previous list element or nil.

// Prev returns the previous list element
// or nil.
func (e *Element) Prev() *Element

// List represents a doubly linked list. The zero value for List is an empty list
// ready to use.

// List represents a doubly linked list.
// The zero value for List is an empty list
// ready to use.
type List struct {
	// contains filtered or unexported fields
}

// New returns an initialized list.

// New returns an initialized list.
func New() *List

// Back returns the last element of list l or nil.

// Back returns the last element of list l
// or nil.
func (l *List) Back() *Element

// Front returns the first element of list l or nil.

// Front returns the first element of list
// l or nil.
func (l *List) Front() *Element

// Init initializes or clears list l.

// Init initializes or clears list l.
func (l *List) Init() *List

// InsertAfter inserts a new element e with value v immediately after mark and
// returns e. If mark is not an element of l, the list is not modified.

// InsertAfter inserts a new element e with
// value v immediately after mark and
// returns e. If mark is not an element of
// l, the list is not modified.
func (l *List) InsertAfter(v interface{}, mark *Element) *Element

// InsertBefore inserts a new element e with value v immediately before mark and
// returns e. If mark is not an element of l, the list is not modified.

// InsertBefore inserts a new element e
// with value v immediately before mark and
// returns e. If mark is not an element of
// l, the list is not modified.
func (l *List) InsertBefore(v interface{}, mark *Element) *Element

// Len returns the number of elements of list l. The complexity is O(1).

// Len returns the number of elements of
// list l. The complexity is O(1).
func (l *List) Len() int

// MoveAfter moves element e to its new position after mark. If e or mark is not an
// element of l, or e == mark, the list is not modified.

// MoveAfter moves element e to its new
// position after mark. If e or mark is not
// an element of l, or e == mark, the list
// is not modified.
func (l *List) MoveAfter(e, mark *Element)

// MoveBefore moves element e to its new position before mark. If e or mark is not
// an element of l, or e == mark, the list is not modified.

// MoveBefore moves element e to its new
// position before mark. If e or mark is
// not an element of l, or e == mark, the
// list is not modified.
func (l *List) MoveBefore(e, mark *Element)

// MoveToBack moves element e to the back of list l. If e is not an element of l,
// the list is not modified.

// MoveToBack moves element e to the back
// of list l. If e is not an element of l,
// the list is not modified.
func (l *List) MoveToBack(e *Element)

// MoveToFront moves element e to the front of list l. If e is not an element of l,
// the list is not modified.

// MoveToFront moves element e to the front
// of list l. If e is not an element of l,
// the list is not modified.
func (l *List) MoveToFront(e *Element)

// PushBack inserts a new element e with value v at the back of list l and returns
// e.

// PushBack inserts a new element e with
// value v at the back of list l and
// returns e.
func (l *List) PushBack(v interface{}) *Element

// PushBackList inserts a copy of an other list at the back of list l. The lists l
// and other may be the same.

// PushBackList inserts a copy of an other
// list at the back of list l. The lists l
// and other may be the same.
func (l *List) PushBackList(other *List)

// PushFront inserts a new element e with value v at the front of list l and
// returns e.

// PushFront inserts a new element e with
// value v at the front of list l and
// returns e.
func (l *List) PushFront(v interface{}) *Element

// PushFrontList inserts a copy of an other list at the front of list l. The lists
// l and other may be the same.

// PushFrontList inserts a copy of an other
// list at the front of list l. The lists l
// and other may be the same.
func (l *List) PushFrontList(other *List)

// Remove removes e from l if e is an element of list l. It returns the element
// value e.Value.

// Remove removes e from l if e is an
// element of list l. It returns the
// element value e.Value.
func (l *List) Remove(e *Element) interface{}
