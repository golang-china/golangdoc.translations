// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package list implements a doubly linked list.
//
// To iterate over a list (where l is a *List):
// 	for e := l.Front(); e != nil; e = e.Next() {
// 		// do something with e.Value
// 	}

// list包实现了双向链表。要遍历一个链表：
//
//     for e := l.Front(); e != nil; e = e.Next() {
//         // do something with e.Value
//     }
package list

// Element is an element of a linked list.

// Element类型代表是双向链表的一个元素。
type Element struct {
	// The value stored with this element.
	Value interface{}
}

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.

// List代表一个双向链表。List零值为一个空的、可用的链表。
type List struct {
}

// New returns an initialized list.

// New创建一个链表。
func New() *List

// Next returns the next list element or nil.

// Next返回链表的后一个元素或者nil。
func (e *Element) Next() *Element

// Prev returns the previous list element or nil.

// Prev返回链表的前一个元素或者nil。
func (e *Element) Prev() *Element

// Back returns the last element of list l or nil.

// Back返回链表最后一个元素或nil。
func (l *List) Back() *Element

// Front returns the first element of list l or nil.

// Front返回链表第一个元素或nil。
func (l *List) Front() *Element

// Init initializes or clears list l.

// Init清空链表。
func (l *List) Init() *List

// InsertAfter inserts a new element e with value v immediately after mark and
// returns e. If mark is not an element of l, the list is not modified.

// InsertAfter将一个值为v的新元素插入到mark后面，并返回新生成的元素。如果mark不
// 是l的元素，l不会被修改。
func (l *List) InsertAfter(v interface{}, mark *Element) *Element

// InsertBefore inserts a new element e with value v immediately before mark and
// returns e. If mark is not an element of l, the list is not modified.

// InsertDefore将一个值为v的新元素插入到mark前面，并返回生成的新元素。如果mark不
// 是l的元素，l不会被修改。
func (l *List) InsertBefore(v interface{}, mark *Element) *Element

// Len returns the number of elements of list l.
// The complexity is O(1).

// Len返回链表中元素的个数，复杂度O(1)。
func (l *List) Len() int

// MoveAfter moves element e to its new position after mark. If e or mark is not
// an element of l, or e == mark, the list is not modified.

// MoveAfter将元素e移动到mark的后面。如果e或mark不是l的元素，或者e==mark，l不会
// 被修改。
func (l *List) MoveAfter(e, mark *Element)

// MoveBefore moves element e to its new position before mark. If e or mark is
// not an element of l, or e == mark, the list is not modified.

// MoveBefore将元素e移动到mark的前面。如果e或mark不是l的元素，或者e==mark，l不会
// 被修改。
func (l *List) MoveBefore(e, mark *Element)

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.

// MoveToBack将元素e移动到链表的最后一个位置，如果e不是l的元素，l不会被修改。
func (l *List) MoveToBack(e *Element)

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.

// MoveToFront将元素e移动到链表的第一个位置，如果e不是l的元素，l不会被修改。
func (l *List) MoveToFront(e *Element)

// PushBack inserts a new element e with value v at the back of list l and
// returns e.

// PushBack将一个值为v的新元素插入链表的最后一个位置，返回生成的新元素。
func (l *List) PushBack(v interface{}) *Element

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same.

// PushBack创建链表other的拷贝，并将链表l的最后一个位置连接到拷贝的第一个位置。
func (l *List) PushBackList(other *List)

// PushFront inserts a new element e with value v at the front of list l and
// returns e.

// PushBack将一个值为v的新元素插入链表的第一个位置，返回生成的新元素。
func (l *List) PushFront(v interface{}) *Element

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same.

// PushFrontList创建链表other的拷贝，并将拷贝的最后一个位置连接到链表l的第一个位
// 置。
func (l *List) PushFrontList(other *List)

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.

// Remove删除链表中的元素e，并返回e.Value。
func (l *List) Remove(e *Element) interface{}

