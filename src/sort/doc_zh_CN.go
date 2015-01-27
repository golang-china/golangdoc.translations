// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package sort provides primitives for sorting slices and user-defined
// collections.

// sort
// 包为切片及用户定义的集合的排序操作提供了原语.
package sort

// Float64s sorts a slice of float64s in increasing order.

// Float64s 以升序排列 float64 切片
func Float64s(a []float64)

// Float64sAreSorted tests whether a slice of float64s is sorted in increasing
// order.

// Float64sAreSorted 判断 float64 切片是否已经按升序排列。
func Float64sAreSorted(a []float64) bool

// Ints sorts a slice of ints in increasing order.

// Ints 以升序排列 int 切片。
func Ints(a []int)

// IntsAreSorted tests whether a slice of ints is sorted in increasing order.

// IntsAreSorted 判断 int 切片是否已经按升序排列。
func IntsAreSorted(a []int) bool

// IsSorted reports whether data is sorted.

// IsSorted 返回数据是否已经排序。
func IsSorted(data Interface) bool

// Search uses binary search to find and return the smallest index i in [0, n) at
// which f(i) is true, assuming that on the range [0, n), f(i) == true implies
// f(i+1) == true. That is, Search requires that f is false for some (possibly
// empty) prefix of the input range [0, n) and then true for the (possibly empty)
// remainder; Search returns the first true index. If there is no such index,
// Search returns n. (Note that the "not found" return value is not -1 as in, for
// instance, strings.Index). Search calls f(i) only for i in the range [0, n).
//
// A common use of Search is to find the index i for a value x in a sorted,
// indexable data structure such as an array or slice. In this case, the argument
// f, typically a closure, captures the value to be searched for, and how the data
// structure is indexed and ordered.
//
// For instance, given a slice data sorted in ascending order, the call
// Search(len(data), func(i int) bool { return data[i] >= 23 }) returns the
// smallest index i such that data[i] >= 23. If the caller wants to find whether 23
// is in the slice, it must test data[i] == 23 separately.
//
// Searching data sorted in descending order would use the <= operator instead of
// the >= operator.
//
// To complete the example above, the following code tries to find the value x in
// an integer slice data sorted in ascending order:
//
//	x := 23
//	i := sort.Search(len(data), func(i int) bool { return data[i] >= x })
//	if i < len(data) && data[i] == x {
//		// x is present at data[i]
//	} else {
//		// x is not present in data,
//		// but i is the index where it would be inserted.
//	}
//
// As a more whimsical example, this program guesses your number:
//
//	func GuessingGame() {
//		var s string
//		fmt.Printf("Pick an integer from 0 to 100.\n")
//		answer := sort.Search(100, func(i int) bool {
//			fmt.Printf("Is your number <= %d? ", i)
//			fmt.Scanf("%s", &s)
//			return s != "" && s[0] == 'y'
//		})
//		fmt.Printf("Your number is %d.\n", answer)
//	}

// Search 使用二分查找法在 [0, n) 中寻找并返回满足 f(i) == true 的最小索引 i，
// 假定该索引在区间 [0, n) 内，则 f(i) == true 就蕴含了 f(i+1) == true。 也就是说，Search 要求 f
// 对于输入区间 [0, n) （可能为空）的前一部分为 false，
// 而对于剩余（可能为空）的部分为 true；Search 返回第一个 f 为 true 时的索引 i。
// 若该索引不存在，Search 就返回
// n。（注意，“未找到”的返回值并不像 strings.Index 这类函数一样返回 -1）。Search 仅当 i 在区间 [0, n)
// 内时才调用 f(i)。
//
// Search
// 常用于在一个已排序的，可索引的数据结构中寻找索引为 i 的值 x，例如数组或切片。 这种情况下，实参
// f，一般是一个闭包，会捕获所要搜索的值，以及索引并排序该数据结构的方式。
//
// 例如，给定一个以升序排列的切片数据，调用
//
//	Search(len(data), func(i int) bool { return data[i] >= 23 })
//
// 会返回满足 data[i] >= 23 的最小索引 i。若调用者想要判断 23 是否在此切片中， 就必须单独测试 data[i] == 23
// 的值。
//
// 搜索降以序排列的数据，需使用 <= 操作符，而非 >= 操作符。
//
// 补全上面的例子,
// 以下代码试图从以升序排列的整数切片中寻找值 x 的索引：
//
//	x := 23
//	i := sort.Search(len(data), func(i int) bool { return data[i] >= x })
//	if i < len(data) && data[i] == x {
//		// x 为 data[i]
//	} else {
//		// x 不在 data 中，但 i 可作为它的索引插入。
//	}
//
// 还有个更有趣的例子，此程序会猜你所想的数字：
//
//	func GuessingGame() {
//		var s string
//		fmt.Printf("Pick an integer from 0 to 100.\n")
//		answer := sort.Search(100, func(i int) bool {
//			fmt.Printf("Is your number <= %d? ", i)
//			fmt.Scanf("%s", &s)
//			return s != "" && s[0] == 'y'
//		})
//		fmt.Printf("Your number is %d.\n", answer)
//	}
func Search(n int, f func(int) bool) int

// SearchFloat64s searches for x in a sorted slice of float64s and returns the
// index as specified by Search. The return value is the index to insert x if x is
// not present (it could be len(a)). The slice must be sorted in ascending order.

// SearchFloat64s 在float64s切片中搜索x并返回索引 如Search函数所述.
// 返回可以插入x值的索引位置，如果x 不存在，返回数组a的长度 切片必须以升序排列
func SearchFloat64s(a []float64, x float64) int

// SearchInts searches for x in a sorted slice of ints and returns the index as
// specified by Search. The return value is the index to insert x if x is not
// present (it could be len(a)). The slice must be sorted in ascending order.

// SearchInts 在ints切片中搜索x并返回索引 如Search函数所述.
// 返回可以插入x值的索引位置，如果x 不存在，返回数组a的长度 切片必须以升序排列
func SearchInts(a []int, x int) int

// SearchStrings searches for x in a sorted slice of strings and returns the index
// as specified by Search. The return value is the index to insert x if x is not
// present (it could be len(a)). The slice must be sorted in ascending order.

// SearchFloat64s 在strings切片中搜索x并返回索引 如Search函数所述.
// 返回可以插入x值的索引位置，如果x 不存在，返回数组a的长度 切片必须以升序排列
func SearchStrings(a []string, x string) int

// Sort sorts data. It makes one call to data.Len to determine n, and O(n*log(n))
// calls to data.Less and data.Swap. The sort is not guaranteed to be stable.

// Sort 对 data 进行排序。 它调用一次 data.Len 来决定排序的长度 n，调用 data.Less 和 data.Swap
// 的开销为 O(n*log(n))。此排序为不稳定排序。
func Sort(data Interface)

// Stable sorts data while keeping the original order of equal elements.
//
// It makes one call to data.Len to determine n, O(n*log(n)) calls to data.Less and
// O(n*log(n)*log(n)) calls to data.Swap.
func Stable(data Interface)

// Strings sorts a slice of strings in increasing order.

// Strings 以升序排列 string 切片。
func Strings(a []string)

// StringsAreSorted tests whether a slice of strings is sorted in increasing order.

// StringsAreSorted 判断 string 切片是否已经按升序排列。
func StringsAreSorted(a []string) bool

// Float64Slice attaches the methods of Interface to []float64, sorting in
// increasing order.

// Float64Slice 针对 []float6
// 实现接口的方法，以升序排列。
type Float64Slice []float64

func (p Float64Slice) Len() int

func (p Float64Slice) Less(i, j int) bool

// Search returns the result of applying SearchFloat64s to the receiver and x.

// Search
// 返回以调用者和x为参数调用SearchFloat64s后的结果
func (p Float64Slice) Search(x float64) int

// Sort is a convenience method.

// Sort 为便捷性方法
func (p Float64Slice) Sort()

func (p Float64Slice) Swap(i, j int)

// IntSlice attaches the methods of Interface to []int, sorting in increasing
// order.

// IntSlice 针对 []int 实现接口的方法，以升序排列。
type IntSlice []int

func (p IntSlice) Len() int

func (p IntSlice) Less(i, j int) bool

// Search returns the result of applying SearchInts to the receiver and x.

// Search
// 返回以调用者和x为参数调用SearchInts后的结果
func (p IntSlice) Search(x int) int

// Sort is a convenience method.

// Sort 为便捷性方法
func (p IntSlice) Sort()

func (p IntSlice) Swap(i, j int)

// A type, typically a collection, that satisfies sort.Interface can be sorted by
// the routines in this package. The methods require that the elements of the
// collection be enumerated by an integer index.

// 任何实现了 sort.Interface
// 的类型（一般为集合），均可使用该包中的方法进行排序。
// 这些方法要求集合内列出元素的索引为整数。
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

// Reverse returns the reverse order for data.
func Reverse(data Interface) Interface

// StringSlice attaches the methods of Interface to []string, sorting in increasing
// order.

// StringSlice 针对 []string 实现接口的方法，以升序排列。
type StringSlice []string

func (p StringSlice) Len() int

func (p StringSlice) Less(i, j int) bool

// Search returns the result of applying SearchStrings to the receiver and x.

// Search
// 返回以调用者和x为参数调用SearchStrings后的结果
func (p StringSlice) Search(x string) int

// Sort is a convenience method.

// Sort 为便捷性方法
func (p StringSlice) Sort()

func (p StringSlice) Swap(i, j int)
