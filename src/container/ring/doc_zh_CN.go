// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package ring implements operations on circular lists.

// Package ring implements operations on
// circular lists.
package ring

// A Ring is an element of a circular list, or ring. Rings do not have a beginning
// or end; a pointer to any ring element serves as reference to the entire ring.
// Empty rings are represented as nil Ring pointers. The zero value for a Ring is a
// one-element ring with a nil Value.

// A Ring is an element of a circular list,
// or ring. Rings do not have a beginning
// or end; a pointer to any ring element
// serves as reference to the entire ring.
// Empty rings are represented as nil Ring
// pointers. The zero value for a Ring is a
// one-element ring with a nil Value.
type Ring struct {
	Value interface{} // for use by client; untouched by this library
	// contains filtered or unexported fields
}

// New creates a ring of n elements.

// New creates a ring of n elements.
func New(n int) *Ring

// Do calls function f on each element of the ring, in forward order. The behavior
// of Do is undefined if f changes *r.

// Do calls function f on each element of
// the ring, in forward order. The behavior
// of Do is undefined if f changes *r.
func (r *Ring) Do(f func(interface{}))

// Len computes the number of elements in ring r. It executes in time proportional
// to the number of elements.

// Len computes the number of elements in
// ring r. It executes in time proportional
// to the number of elements.
func (r *Ring) Len() int

// Link connects ring r with ring s such that r.Next() becomes s and returns the
// original value for r.Next(). r must not be empty.
//
// If r and s point to the same ring, linking them removes the elements between r
// and s from the ring. The removed elements form a subring and the result is a
// reference to that subring (if no elements were removed, the result is still the
// original value for r.Next(), and not nil).
//
// If r and s point to different rings, linking them creates a single ring with the
// elements of s inserted after r. The result points to the element following the
// last element of s after insertion.

// Link connects ring r with ring s such
// that r.Next() becomes s and returns the
// original value for r.Next(). r must not
// be empty.
//
// If r and s point to the same ring,
// linking them removes the elements
// between r and s from the ring. The
// removed elements form a subring and the
// result is a reference to that subring
// (if no elements were removed, the result
// is still the original value for
// r.Next(), and not nil).
//
// If r and s point to different rings,
// linking them creates a single ring with
// the elements of s inserted after r. The
// result points to the element following
// the last element of s after insertion.
func (r *Ring) Link(s *Ring) *Ring

// Move moves n % r.Len() elements backward (n < 0) or forward (n >= 0) in the ring
// and returns that ring element. r must not be empty.

// Move moves n % r.Len() elements backward
// (n < 0) or forward (n >= 0) in the ring
// and returns that ring element. r must
// not be empty.
func (r *Ring) Move(n int) *Ring

// Next returns the next ring element. r must not be empty.

// Next returns the next ring element. r
// must not be empty.
func (r *Ring) Next() *Ring

// Prev returns the previous ring element. r must not be empty.

// Prev returns the previous ring element.
// r must not be empty.
func (r *Ring) Prev() *Ring

// Unlink removes n % r.Len() elements from the ring r, starting at r.Next(). If n
// % r.Len() == 0, r remains unchanged. The result is the removed subring. r must
// not be empty.

// Unlink removes n % r.Len() elements from
// the ring r, starting at r.Next(). If n %
// r.Len() == 0, r remains unchanged. The
// result is the removed subring. r must
// not be empty.
func (r *Ring) Unlink(n int) *Ring
