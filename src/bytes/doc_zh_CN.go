// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package bytes implements functions for the manipulation of byte slices. It is
// analogous to the facilities of the strings package.

// bytes 包实现了操作 byte 切片的常用函数.
// bytes 包和 strings 包的结构很相似.
package bytes

// MinRead is the minimum slice size passed to a Read call by Buffer.ReadFrom. As
// long as the Buffer has at least MinRead bytes beyond what is required to hold
// the contents of r, ReadFrom will not grow the underlying buffer.

// MinRead 是被 Buffer.ReadFrom 传递给 Read 调用的切片的最小尺寸.
// 只要 Buffer 在保存 r 内容之外有最少 MinRead 字节的余量, 其 ReadFrom 方法就不会增加底层的缓冲.
const MinRead = 512

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a
// buffer.

// 如果不能申请足够保存数据的缓冲, ErrTooLarge 就会被传递给 panic 函数.
var ErrTooLarge = errors.New("bytes.Buffer: too large")

// Compare returns an integer comparing two byte slices lexicographically. The
// result will be 0 if a==b, -1 if a < b, and +1 if a > b. A nil argument is
// equivalent to an empty slice.

// Compare 函数返回一个整数表示两个 byte 切片按字典序的比较结果.
// 如果 a==b 返回0, 如果 a<b 返回 -1, 否则返回+1.
// nil 参数视为空切片.
func Compare(a, b []byte) int

// Contains reports whether subslice is within b.

// Contains 判断切片 b 是否包含子切片 subslice.
func Contains(b, subslice []byte) bool

// Count counts the number of non-overlapping instances of sep in s.

// Count 统计切片 s 包含 切片 b 的次数.
func Count(s, sep []byte) int

// Equal returns a boolean reporting whether a and b are the same length and
// contain the same bytes. A nil argument is equivalent to an empty slice.

// Equal 判断两个切片的内容是否完全相同.
// nil 参数视为空切片.
func Equal(a, b []byte) bool

// EqualFold reports whether s and t, interpreted as UTF-8 strings, are equal under
// Unicode case-folding.

// EqualFold 判断两个 UTF-8 编码切片是否相同, 忽略大小写.
func EqualFold(s, t []byte) bool

// Fields splits the slice s around each instance of one or more consecutive white
// space characters, returning a slice of subslices of s or an empty list if s
// contains only white space.

// Fields 返回将字符串按照空白分割的多个子切片.
// 如果字符串全部是空白或者是空字符串的话, 会返回空切片.
func Fields(s []byte) [][]byte

// FieldsFunc interprets s as a sequence of UTF-8-encoded Unicode code points. It
// splits the slice s at each run of code points c satisfying f(c) and returns a
// slice of subslices of s. If all code points in s satisfy f(c), or len(s) == 0,
// an empty slice is returned. FieldsFunc makes no guarantees about the order in
// which it calls f(c). If f does not return consistent results for a given c,
// FieldsFunc may crash.

// FieldsFunc 类似 Fields, 但使用函数f来确定分割符.
// 如果字符串全部是分隔符或者是空字符串的话, 会返回空切片.
func FieldsFunc(s []byte, f func(rune) bool) [][]byte

// HasPrefix tests whether the byte slice s begins with prefix.

// HasPrefix 判断 s 前缀是否是 prefix.
func HasPrefix(s, prefix []byte) bool

// HasSuffix tests whether the byte slice s ends with suffix.

// HasSuffix 判断 s 后缀是否是 suffix.
func HasSuffix(s, suffix []byte) bool

// Index returns the index of the first instance of sep in s, or -1 if sep is not
// present in s.

// Index 返回子切片 sep 在 s 中第一次出现的位置, 不存在则返回 -1.
func Index(s, sep []byte) int

// IndexAny interprets s as a sequence of UTF-8-encoded Unicode code points. It
// returns the byte index of the first occurrence in s of any of the Unicode code
// points in chars. It returns -1 if chars is empty or if there is no code point in
// common.

// IndexAny 返回字符串 chars 中的任一个 utf-8 编码字符在 s 中第一次出现的位置,
// 如不存在或者 chars 为空字符串则返回 -1.
func IndexAny(s []byte, chars string) int

// IndexByte returns the index of the first instance of c in s, or -1 if c is not
// present in s.
func IndexByte(s []byte, c byte) int

// IndexFunc interprets s as a sequence of UTF-8-encoded Unicode code points. It
// returns the byte index in s of the first Unicode code point satisfying f(c), or
// -1 if none do.

// IndexFunc 返回 s 中第一个满足函数 f 的位置 i, 不存在则返回 -1.
func IndexFunc(s []byte, f func(r rune) bool) int

// IndexRune interprets s as a sequence of UTF-8-encoded Unicode code points. It
// returns the byte index of the first occurrence in s of the given rune. It
// returns -1 if rune is not present in s.

// IndexRune 返回 unicode 字符 r 的 utf-8 编码在 s 中第一次出现的位置, 不存在则返回 -1.
func IndexRune(s []byte, r rune) int

// Join concatenates the elements of s to create a new byte slice. The separator
// sep is placed between elements in the resulting slice.

// Join 将一系列 byte 切片连接为一个 byte 切片, 之间用 sep 来分隔, 返回生成的新切片.
func Join(s [][]byte, sep []byte) []byte

// LastIndex returns the index of the last instance of sep in s, or -1 if sep is
// not present in s.

// LastIndex 返回切片 sep 在字符串 s 中最后一次出现的位置, 不存在则返回-1.
func LastIndex(s, sep []byte) int

// LastIndexAny interprets s as a sequence of UTF-8-encoded Unicode code points. It
// returns the byte index of the last occurrence in s of any of the Unicode code
// points in chars. It returns -1 if chars is empty or if there is no code point in
// common.

// LastIndexAny 返回字符串 chars 中的任一个 utf-8 字符在 s 中最后一次出现的位置,
// 如不存在或者 chars 为空字符串则返回 -1.
func LastIndexAny(s []byte, chars string) int

// LastIndexFunc interprets s as a sequence of UTF-8-encoded Unicode code points.
// It returns the byte index in s of the last Unicode code point satisfying f(c),
// or -1 if none do.

// LastIndexFunc 返回 s 中最后一个满足函数 f 的 unicode 码值的位置i,
// 不存在则返回 -1.
func LastIndexFunc(s []byte, f func(r rune) bool) int

// Map returns a copy of the byte slice s with all its characters modified
// according to the mapping function. If mapping returns a negative value, the
// character is dropped from the string with no replacement. The characters in s
// and the output are interpreted as UTF-8-encoded Unicode code points.

// Map 将 s 的每一个 unicode 码值 r 都替换为 mapping(r), 返回这些新码值组成的切片拷贝.
// 如果 mapping 返回一个负值, 将会丢弃该码值而不会被替换.
func Map(mapping func(r rune) rune, s []byte) []byte

// Repeat returns a new byte slice consisting of count copies of b.

// Repeat 返回 count 个 b 串联形成的新的切片.
func Repeat(b []byte, count int) []byte

// Replace returns a copy of the slice s with the first n non-overlapping instances
// of old replaced by new. If old is empty, it matches at the beginning of the
// slice and after each UTF-8 sequence, yielding up to k+1 replacements for a
// k-rune slice. If n < 0, there is no limit on the number of replacements.

// Replace 返回将 s 中前 n 个不重叠 old 切片序列都替换为 new 的新的切片拷贝.
// 如果 n<0 会替换所有 old 子切片.
func Replace(s, old, new []byte, n int) []byte

// Runes returns a slice of runes (Unicode code points) equivalent to s.

// Runes 函数返回和 s 等价的 rune 切片.
func Runes(s []byte) []rune

// Split slices s into all subslices separated by sep and returns a slice of the
// subslices between those separators. If sep is empty, Split splits after each
// UTF-8 sequence. It is equivalent to SplitN with a count of -1.

// Split 用去掉 s 中出现的 sep 的方式进行分割, 会分割到结尾,
// 并返回生成的所有 byte 切片组成的切片.
// 如果 sep 为空字符, Split 会将 s 切分成每一个 unicode 码值一个 byte 切片.
func Split(s, sep []byte) [][]byte

// SplitAfter slices s into all subslices after each instance of sep and returns a
// slice of those subslices. If sep is empty, SplitAfter splits after each UTF-8
// sequence. It is equivalent to SplitAfterN with a count of -1.

// SplitAfter 用从 s 中出现的 sep 后面切断的方式进行分割, 会分割到结尾,
// 并返回生成的所有 byte 切片组成的切片.
// 如果 sep 为空字符, Split 会将 s 切分成每一个 unicode 码值一个 byte 切片.
func SplitAfter(s, sep []byte) [][]byte

// SplitAfterN slices s into subslices after each instance of sep and returns a
// slice of those subslices. If sep is empty, SplitAfterN splits after each UTF-8
// sequence. The count determines the number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices

// SplitAfterN 用从 s 中出现的 sep 后面切断的方式进行分割, 会分割到最多 n 个子切片,
// 并返回生成的所有 byte 切片组成的切片.
// 如果 sep 为空字符, Split 会将 s 切分成每一个 unicode 码值一个 byte 切片.
// 参数n决定返回的切片的数目:
//
//	n > 0: 返回的切片最多n个子字符串；最后一个子字符串包含未进行切割的部分。
//	n == 0: 返回nil
//	n < 0: 返回所有的子字符串组成的切片
func SplitAfterN(s, sep []byte, n int) [][]byte

// SplitN slices s into subslices separated by sep and returns a slice of the
// subslices between those separators. If sep is empty, SplitN splits after each
// UTF-8 sequence. The count determines the number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices

// SplitN 用去掉 s 中出现的 sep 的方式进行分割, 会分割到最多n个子切片,
// 并返回生成的所有 byte 切片组成的切片.
// 如果 sep 为空字符, Split 会将 s 切分成每一个 unicode 码值一个 byte 切片.
// 参数n决定返回的切片的数目:
//
//	n > 0: 返回的切片最多n个子字符串, 最后一个子字符串包含未进行切割的部分.
//	n == 0: 返回nil
//	n < 0: 返回所有的子字符串组成的切片
func SplitN(s, sep []byte, n int) [][]byte

// Title returns a copy of s with all Unicode letters that begin words mapped to
// their title case.
//
// BUG: The rule Title uses for word boundaries does not handle Unicode punctuation
// properly.

// Title 返回 s 中每个单词的首字母都改为标题格式的拷贝.
//
// BUG: Title 用于划分单词的规则不能很好的处理 Unicode 标点符号.
func Title(s []byte) []byte

// ToLower returns a copy of the byte slice s with all Unicode letters mapped to
// their lower case.

// ToLower 返回将所有字母都转为对应的小写版本的拷贝.
func ToLower(s []byte) []byte

// ToLowerSpecial returns a copy of the byte slice s with all Unicode letters
// mapped to their lower case, giving priority to the special casing rules.

// ToLowerSpecial 使用 _case 规定的字符映射, 返回将所有字母都转为对应的小写版本的拷贝.
func ToLowerSpecial(_case unicode.SpecialCase, s []byte) []byte

// ToTitle returns a copy of the byte slice s with all Unicode letters mapped to
// their title case.

// ToTitle 返回将所有字母都转为对应的标题版本的拷贝.
func ToTitle(s []byte) []byte

// ToTitleSpecial returns a copy of the byte slice s with all Unicode letters
// mapped to their title case, giving priority to the special casing rules.

// ToTitleSpecial 使用 _case 规定的字符映射, 返回将所有字母都转为对应的标题版本的拷贝.
func ToTitleSpecial(_case unicode.SpecialCase, s []byte) []byte

// ToUpper returns a copy of the byte slice s with all Unicode letters mapped to
// their upper case.

// ToUpper 返回将所有字母都转为对应的大写版本的拷贝.
func ToUpper(s []byte) []byte

// ToUpperSpecial returns a copy of the byte slice s with all Unicode letters
// mapped to their upper case, giving priority to the special casing rules.

// ToUpperSpecial 使用 _case 规定的字符映射, 返回将所有字母都转为对应的大写版本的拷贝.
func ToUpperSpecial(_case unicode.SpecialCase, s []byte) []byte

// Trim returns a subslice of s by slicing off all leading and trailing
// UTF-8-encoded Unicode code points contained in cutset.

// Trim 返回将 s 前后端所有 cutset 包含的 unicode 码值都去掉的子切片.
func Trim(s []byte, cutset string) []byte

// TrimFunc returns a subslice of s by slicing off all leading and trailing
// UTF-8-encoded Unicode code points c that satisfy f(c).
func TrimFunc(s []byte, f func(r rune) bool) []byte

// TrimLeft returns a subslice of s by slicing off all leading UTF-8-encoded
// Unicode code points contained in cutset.

// TrimLeft 返回将 s 前端所有 cutset 包含的 unicode 码值都去掉的子切片.
func TrimLeft(s []byte, cutset string) []byte

// TrimLeftFunc returns a subslice of s by slicing off all leading UTF-8-encoded
// Unicode code points c that satisfy f(c).

// TrimLeftFunc 返回将 s 前端所有满足 f 的 unicode 码值都去掉的子切片.
func TrimLeftFunc(s []byte, f func(r rune) bool) []byte

// TrimPrefix returns s without the provided leading prefix string. If s doesn't
// start with prefix, s is returned unchanged.

// TrimPrefix 返回去除 s 可能的前缀 prefix 的子切片.
func TrimPrefix(s, prefix []byte) []byte

// TrimRight returns a subslice of s by slicing off all trailing UTF-8-encoded
// Unicode code points that are contained in cutset.

// TrimRight 返回将 s 后端所有 cutset 包含的 unicode 码值都去掉的子切片.
func TrimRight(s []byte, cutset string) []byte

// TrimRightFunc returns a subslice of s by slicing off all trailing UTF-8 encoded
// Unicode code points c that satisfy f(c).

// TrimRightFunc 返回将 s 前后端所有满足 f 的 unicode 码值都去掉的子切片.
func TrimRightFunc(s []byte, f func(r rune) bool) []byte

// TrimSpace returns a subslice of s by slicing off all leading and trailing white
// space, as defined by Unicode.

// TrimSpace 返回将 s 前后端所有空白都去掉的子切片.
func TrimSpace(s []byte) []byte

// TrimSuffix returns s without the provided trailing suffix string. If s doesn't
// end with suffix, s is returned unchanged.

// TrimSuffix 返回去除 s 可能的后缀 suffix 的子切片.
func TrimSuffix(s, suffix []byte) []byte

// A Buffer is a variable-sized buffer of bytes with Read and Write methods. The
// zero value for Buffer is an empty buffer ready to use.

// Buffer 是一个实现了读写方法的可变大小的字节缓冲.
// 本类型的零值是一个空的可用于读写的缓冲.
type Buffer struct {
	// contains filtered or unexported fields
}

// NewBuffer creates and initializes a new Buffer using buf as its initial
// contents. It is intended to prepare a Buffer to read existing data. It can also
// be used to size the internal buffer for writing. To do that, buf should have the
// desired capacity but a length of zero.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is sufficient
// to initialize a Buffer.

// NewBuffer 使用 buf 作为初始内容创建并初始化一个 Buffer.
// 本函数用于创建一个用于读取已存在数据的 buffer, 也用于指定用于写入的内部缓冲的大小.
// 此时, buf 应为一个具有指定容量但长度为 0 的切片.
// buf会 被作为返回值的底层缓冲切片.
//
// 大多数情况下, new(Buffer)就足以初始化一个 Buffer 了.
func NewBuffer(buf []byte) *Buffer

// NewBufferString creates and initializes a new Buffer using string s as its
// initial contents. It is intended to prepare a buffer to read an existing string.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is sufficient
// to initialize a Buffer.

// NewBufferString 使用 s 作为初始内容创建并初始化一个 Buffer.
// 本函数用于创建一个用于读取已存在数据的 buffer.
//
// 大多数情况下, new(Buffer)就足以初始化一个Buffer了.
func NewBufferString(s string) *Buffer

// Bytes returns a slice of the contents of the unread portion of the buffer;
// len(b.Bytes()) == b.Len(). If the caller changes the contents of the returned
// slice, the contents of the buffer will change provided there are no intervening
// method calls on the Buffer.

// 返回未读取部分字节数据的切片, len(b.Bytes()) == b.Len().
// 如果中间没有调用其他方法, 修改返回的切片的内容会直接改变Buffer的内容.
func (b *Buffer) Bytes() []byte

// Grow grows the buffer's capacity, if necessary, to guarantee space for another n
// bytes. After Grow(n), at least n bytes can be written to the buffer without
// another allocation. If n is negative, Grow will panic. If the buffer can't grow
// it will panic with ErrTooLarge.

// Grow 必要时会增加缓冲的容量, 以保证n字节的剩余空间.
// 调用 Grow(n) 后至少可以向缓冲中写入 n 字节数据而无需申请内存.
// 如果 n 小于零或者不能增加容量都会 panic.
func (b *Buffer) Grow(n int)

// Len returns the number of bytes of the unread portion of the buffer; b.Len() ==
// len(b.Bytes()).

// Len 返回缓冲中未读取部分的字节长度: b.Len() == len(b.Bytes()).
func (b *Buffer) Len() int

// Next returns a slice containing the next n bytes from the buffer, advancing the
// buffer as if the bytes had been returned by Read. If there are fewer than n
// bytes in the buffer, Next returns the entire buffer. The slice is only valid
// until the next call to a read or write method.

// Next 返回未读取部分前 n 字节数据的切片, 并且移动读取位置, 就像调用了Read方法一样.
// 如果缓冲内数据不足, 会返回整个数据的切片.
// 切片只在下一次调用 b 的读/写方法前才合法.
func (b *Buffer) Next(n int) []byte

// Read reads the next len(p) bytes from the buffer or until the buffer is drained.
// The return value n is the number of bytes read. If the buffer has no data to
// return, err is io.EOF (unless len(p) is zero); otherwise it is nil.

// Read 方法从缓冲中读取数据直到缓冲中没有数据或者读取了 len(p) 字节数据,
// 将读取的数据写入p.
// 返回值 n 是读取的字节数, 除非缓冲中完全没有数据可以读取并写入 p,
// 此时返回值 err 为 io.EOF, 否则 err 总是 nil.
func (b *Buffer) Read(p []byte) (n int, err error)

// ReadByte reads and returns the next byte from the buffer. If no byte is
// available, it returns error io.EOF.

// ReadByte 读取并返回缓冲中的下一个字节.
// 如果没有数据可用, 返回值 err 为 io.EOF.
func (b *Buffer) ReadByte() (c byte, err error)

// ReadBytes reads until the first occurrence of delim in the input, returning a
// slice containing the data up to and including the delimiter. If ReadBytes
// encounters an error before finding a delimiter, it returns the data read before
// the error and the error itself (often io.EOF). ReadBytes returns err != nil if
// and only if the returned data does not end in delim.

// ReadBytes 读取直到第一次遇到 delim 字节, 返回一个包含已读取的数据和delim字节的切片.
// 如果 ReadBytes 方法在读取到 delim 之前遇到了错误, 它会返回在错误之前读取的数据以及该错误.
// 当且仅当 ReadBytes 方法返回的切片不以 delim 结尾时, 会返回一个非 nil 的错误.
func (b *Buffer) ReadBytes(delim byte) (line []byte, err error)

// ReadFrom reads data from r until EOF and appends it to the buffer, growing the
// buffer as needed. The return value n is the number of bytes read. Any error
// except io.EOF encountered during the read is also returned. If the buffer
// becomes too large, ReadFrom will panic with ErrTooLarge.

// ReadFrom 从 r 中读取数据直到结束并将读取的数据写入缓冲中, 如必要会增加缓冲容量.
// 返回值 n 为从 r 读取并写入 b 的字节数, 会返回读取时遇到的除了io.EOF之外的错误.
// 如果缓冲太大, ReadFrom 会采用错误值 ErrTooLarge 引发 panic.
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error)

// ReadRune reads and returns the next UTF-8-encoded Unicode code point from the
// buffer. If no bytes are available, the error returned is io.EOF. If the bytes
// are an erroneous UTF-8 encoding, it consumes one byte and returns U+FFFD, 1.

// ReadRune 读取并返回缓冲中的下一个 utf-8 码值.
// 如果没有数据可用, 返回值 err 为 io.EOF.
// 如果缓冲中的数据是错误的 utf-8 编码, 本方法会吃掉一字节并返回 (U+FFFD, 1, nil).
func (b *Buffer) ReadRune() (r rune, size int, err error)

// ReadString reads until the first occurrence of delim in the input, returning a
// string containing the data up to and including the delimiter. If ReadString
// encounters an error before finding a delimiter, it returns the data read before
// the error and the error itself (often io.EOF). ReadString returns err != nil if
// and only if the returned data does not end in delim.

// ReadString 读取直到第一次遇到 delim 字节, 返回一个包含已读取的数据和delim字节的字符串.
// 如果 ReadString 方法在读取到 delim 之前遇到了错误, 它会返回在错误之前读取的数据以及该错误.
// 当且仅当 ReadString 方法返回的切片不以 delim 结尾时, 会返回一个非 nil 的错误.
func (b *Buffer) ReadString(delim byte) (line string, err error)

// Reset resets the buffer so it has no content. b.Reset() is the same as
// b.Truncate(0).

// Reset 重设缓冲, 因此会丢弃全部内容, 等价于 b.Truncate(0).
func (b *Buffer) Reset()

// String returns the contents of the unread portion of the buffer as a string. If
// the Buffer is a nil pointer, it returns "<nil>".

// String 将未读取部分的字节数据作为字符串返回, 如果 b 是 nil 指针, 会返回 "<nil>".
func (b *Buffer) String() string

// Truncate discards all but the first n unread bytes from the buffer. It panics if
// n is negative or greater than the length of the buffer.

// Truncate 丢弃缓冲中除前 n 字节数据外的其它数据, 如果 n 小于零或者大于缓冲容量将panic.
func (b *Buffer) Truncate(n int)

// UnreadByte unreads the last byte returned by the most recent read operation. If
// write has happened since the last read, UnreadByte returns an error.

// UnreadByte 吐出最近一次读取操作读取的最后一个字节.
// 如果最后一次读取操作之后进行了写入, 本方法会返回错误.
func (b *Buffer) UnreadByte() error

// UnreadRune unreads the last rune returned by ReadRune. If the most recent read
// or write operation on the buffer was not a ReadRune, UnreadRune returns an
// error. (In this regard it is stricter than UnreadByte, which will unread the
// last byte from any read operation.)

// UnreadRune 吐出最近一次调用 ReadRune 方法读取的 unicode 码值.
// 如果最近一次读写操作不是 ReadRune, 本方法会返回错误.
func (b *Buffer) UnreadRune() error

// Write appends the contents of p to the buffer, growing the buffer as needed. The
// return value n is the length of p; err is always nil. If the buffer becomes too
// large, Write will panic with ErrTooLarge.

// Write 将 p 的内容写入缓冲中, 如必要会增加缓冲容量.
// 返回值 n 为 len(p), err 总是nil.
// 如果缓冲变得太大, Write 会采用错误值 ErrTooLarge 引发 panic.
func (b *Buffer) Write(p []byte) (n int, err error)

// WriteByte appends the byte c to the buffer, growing the buffer as needed. The
// returned error is always nil, but is included to match bufio.Writer's WriteByte.
// If the buffer becomes too large, WriteByte will panic with ErrTooLarge.

// WriteByte 将字节 c 写入缓冲中, 如必要会增加缓冲容量.
// 返回值总是 nil, 但仍保留以匹配 bufio.Writer 的 WriteByte 方法.
// 如果缓冲太大, WriteByte 会采用错误值 ErrTooLarge 引发 panic.
func (b *Buffer) WriteByte(c byte) error

// WriteRune appends the UTF-8 encoding of Unicode code point r to the buffer,
// returning its length and an error, which is always nil but is included to match
// bufio.Writer's WriteRune. The buffer is grown as needed; if it becomes too
// large, WriteRune will panic with ErrTooLarge.

// WriteByte 将 unicode 码值 r 的 utf-8 编码写入缓冲中, 如必要会增加缓冲容量.
// 返回值总是 nil, 但仍保留以匹配 bufio.Writer 的 WriteRune 方法.
// 如果缓冲太大, WriteRune 会采用错误值 ErrTooLarge 引发 panic.
func (b *Buffer) WriteRune(r rune) (n int, err error)

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the buffer
// becomes too large, WriteString will panic with ErrTooLarge.

// WriteString 将 s 的内容写入缓冲中, 如必要会增加缓冲容量.
// 返回值 n 为 len(p), err总是nil.
// 如果缓冲变得太大, WriteString 会采用错误值 ErrTooLarge 引发 panic.
func (b *Buffer) WriteString(s string) (n int, err error)

// WriteTo writes data to w until the buffer is drained or an error occurs. The
// return value n is the number of bytes written; it always fits into an int, but
// it is int64 to match the io.WriterTo interface. Any error encountered during the
// write is also returned.

// WriteTo 从缓冲中读取数据直到缓冲内没有数据或遇到错误, 并将这些数据写入 w.
// 返回值 n 为从 b 读取并写入 w 的字节数, 返回值总是可以无溢出的写入int类型,
// 但为了匹配 io.WriterTo 接口设为 int64 类型.
// 从 b 读取是遇到的非 io.EOF 错误及写入 w 时遇到的错误都会终止本方法并返回该错误.
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error)

// A Reader implements the io.Reader, io.ReaderAt, io.WriterTo, io.Seeker,
// io.ByteScanner, and io.RuneScanner interfaces by reading from a byte slice.
// Unlike a Buffer, a Reader is read-only and supports seeking.

// Reader 类型通过从一个 byte 切片读取数据,
// 实现了 io.Reader, io.Seeker, io.ReaderAt, io.WriterTo, io.ByteScanner, io.RuneScanner 接口.
type Reader struct {
	// contains filtered or unexported fields
}

// NewReader returns a new Reader reading from b.

// NewReader 创建一个从 s 读取数据的 Reader.
func NewReader(b []byte) *Reader

// Len returns the number of bytes of the unread portion of the slice.

// Len 返回 r 包含的切片中还没有被读取的部分.
func (r *Reader) Len() int

func (r *Reader) Read(b []byte) (n int, err error)

func (r *Reader) ReadAt(b []byte, off int64) (n int, err error)

func (r *Reader) ReadByte() (b byte, err error)

func (r *Reader) ReadRune() (ch rune, size int, err error)

// Seek implements the io.Seeker interface.

// Seek 实现了 io.Seeker 接口.
func (r *Reader) Seek(offset int64, whence int) (int64, error)

func (r *Reader) UnreadByte() error

func (r *Reader) UnreadRune() error

// WriteTo implements the io.WriterTo interface.

// WriteTo 实现了 io.WriterTo 接口.
func (r *Reader) WriteTo(w io.Writer) (n int64, err error)
