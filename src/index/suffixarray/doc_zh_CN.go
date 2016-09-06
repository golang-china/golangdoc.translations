// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package suffixarray implements substring search in logarithmic time using an
// in-memory suffix array.
//
// Example use:
//
// 	// create index for some data
// 	index := suffixarray.New(data)
//
// 	// lookup byte slice s
// 	offsets1 := index.Lookup(s, -1) // the list of all indices where s occurs in data
// 	offsets2 := index.Lookup(s, 3)  // the list of at most 3 indices where s occurs in data

// suffixarrayb包通过使用内存中的后缀树实现了对数级时间消耗的子字符串搜索。
//
// 用法举例：
//
//     // 创建数据的索引
//     index := suffixarray.New(data)
//     // 查找切片s
//     offsets1 := index.Lookup(s, -1) // 返回data中所有s出现的位置
//     offsets2 := index.Lookup(s, 3)  // 返回data中最多3个所有s出现的位置
package suffixarray

import (
	"bytes"
	"encoding/binary"
	"io"
	"regexp"
	"sort"
)

// Index implements a suffix array for fast substring search.

// Index类型实现了用于快速子字符串搜索的后缀数组。
type Index struct {
}

// New creates a new Index for data.
// Index creation time is O(N*log(N)) for N = len(data).

// 使用给出的[]byte数据生成一个*Index，时间复杂度O(N*log(N))。
func New(data []byte) *Index

// Bytes returns the data over which the index was created.
// It must not be modified.

// 返回创建x时提供的[]byte数据，注意不能修改返回值。
func (x *Index) Bytes() []byte

// FindAllIndex returns a sorted list of non-overlapping matches of the
// regular expression r, where a match is a pair of indices specifying
// the matched slice of x.Bytes(). If n < 0, all matches are returned
// in successive order. Otherwise, at most n matches are returned and
// they may not be successive. The result is nil if there are no matches,
// or if n == 0.

// 返回一个正则表达式r的不重叠的匹配的经过排序的列表，一个匹配表示为一对指定了匹
// 配结果的切片的索引（相对于x.Bytes())。如果n<0，返回全部匹配；如果n==0或匹配失
// 败，返回nil；否则n为result的最大长度。
func (x *Index) FindAllIndex(r *regexp.Regexp, n int) (result [][]int)

// Lookup returns an unsorted list of at most n indices where the byte string s
// occurs in the indexed data. If n < 0, all occurrences are returned.
// The result is nil if s is empty, s is not found, or n == 0.
// Lookup time is O(log(N)*len(s) + len(result)) where N is the
// size of the indexed data.

// 返回一个未排序的列表，内为s在被索引为index的切片数据中出现的位置。如果n<0，返
// 回全部匹配；如果n==0或s为空，返回nil；否则n为result的最大长度。时间复杂度
// O(log(N)*len(s) + len(result))，其中N是被索引的数据的大小。
func (x *Index) Lookup(s []byte, n int) (result []int)

// Read reads the index from r into x; x must not be nil.

// 从r中读取一个index写入x，x不能为nil。
func (x *Index) Read(r io.Reader) error

// Write writes the index x to w.

// 将x中的index写入w中，x不能为nil。
func (x *Index) Write(w io.Writer) error

