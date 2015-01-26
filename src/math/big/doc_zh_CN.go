// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package big implements multi-precision arithmetic (big numbers). The following
// numeric types are supported:
//
//	- Int	signed integers
//	- Rat	rational numbers
//
// Methods are typically of the form:
//
//	func (z *Int) Op(x, y *Int) *Int	(similar for *Rat)
//
// and implement operations z = x Op y with the result as receiver; if it is one of
// the operands it may be overwritten (and its memory reused). To enable chaining
// of operations, the result is also returned. Methods returning a result other
// than *Int or *Rat take one of the operands as the receiver.

// big 包实现了（大数的）高精度运算. 它支持以下数值类型：
//
//	- Int	带符号整数
//	- Rat	有理数
//
// 典型的方法形式如下：
//
//	func (z *Int) Op(x, y *Int) *Int	（*Rat 同理）
//
// 它实现了像 z = x Op y
// 这样的操作，并将其结果作为接收者；若接收者为操作数之一，
// 其值可能会被覆盖（而内存则会被重用）。为保留操作，其结果也会被返回。若该方法返回除
// *Int 或 *Rat 之外的结果，其中一个操作数将被作为接收者。
package big

// MaxBase is the largest number base accepted for string conversions.

// MaxBase is the largest number base
// accepted for string conversions.
const MaxBase = 'z' - 'a' + 10 + 1 // = hexValue('z') + 1

// An Int represents a signed multi-precision integer. The zero value for an Int
// represents the value 0.

// Int 表示一个带符号多精度整数。 Int 的零值为值 0。
type Int struct {
	// contains filtered or unexported fields
}

// NewInt allocates and returns a new Int set to x.

// NewInt 为 x 分配并返回一个新的 Int。
func NewInt(x int64) *Int

// Abs sets z to |x| (the absolute value of x) and returns z.

// Abs 将 z 置为 |x|（即 x 的绝对值）并返回 z。
func (z *Int) Abs(x *Int) *Int

// Add sets z to the sum x+y and returns z.

// Add 将 z 置为 x+y 的和并返回 z。
func (z *Int) Add(x, y *Int) *Int

// And sets z = x & y and returns z.

// And 置 z = x & y 并返回 z。
func (z *Int) And(x, y *Int) *Int

// AndNot sets z = x &^ y and returns z.

// AndNot 置 z = x &^ y 并返回 z。
func (z *Int) AndNot(x, y *Int) *Int

// Binomial sets z to the binomial coefficient of (n, k) and returns z.

// Binomial 将 z 置为 (n, k) 的二项式系数并返回 z。
func (z *Int) Binomial(n, k int64) *Int

// Bit returns the value of the i'th bit of x. That is, it returns (x>>i)&1. The
// bit index i must be >= 0.

// Bit 返回 x 第 i 位的值。换言之，它返回 (x>>i)&1。位下标 i
// 必须 >= 0。
func (x *Int) Bit(i int) uint

// BitLen returns the length of the absolute value of x in bits. The bit length of
// 0 is 0.

// BitLen 返回 z 的绝对值的位数长度。0 的位长为 0.
func (x *Int) BitLen() int

// Bits provides raw (unchecked but fast) access to x by returning its absolute
// value as a little-endian Word slice. The result and x share the same underlying
// array. Bits is intended to support implementation of missing low-level Int
// functionality outside this package; it should be avoided otherwise.

// Bits 提供了对 z
// 的原始访问（未经检查但很快）。它通过将其绝对值作为小端序的 Word
// 切片返回来实现。其结果与 x 共享同一底层数组。Bits
// 旨在支持此包外缺失的底层 Int 功能的实现，除此之外应尽量避免。
func (x *Int) Bits() []Word

// Bytes returns the absolute value of x as a big-endian byte slice.

// Bytes 将 x 的绝对值作为大端序的字节切片返回。
func (x *Int) Bytes() []byte

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y

// Cmp 比较 x 和 y 并返回：
//
//	若 x <  y 则为 -1
//	若 x == y 则为  0
//	若 x >  y 则为 +1
func (x *Int) Cmp(y *Int) (r int)

// Div sets z to the quotient x/y for y != 0 and returns z. If y == 0, a
// division-by-zero run-time panic occurs. Div implements Euclidean division
// (unlike Go); see DivMod for more details.

// Div 在 y != 0 时，将 z 置为 x/y 的商并返回 z。 若 y
// == 0，就会产生一个除以零的运行时派错。 Div
// 实现了欧氏除法（与Go不同），更多详情见 DivMod。
func (z *Int) Div(x, y *Int) *Int

// DivMod sets z to the quotient x div y and m to the modulus x mod y and returns
// the pair (z, m) for y != 0. If y == 0, a division-by-zero run-time panic occurs.
//
// DivMod implements Euclidean division and modulus (unlike Go):
//
//	q = x div y  such that
//	m = x - y*q  with 0 <= m < |q|
//
// (See Raymond T. Boute, ``The Euclidean definition of the functions div and
// mod''. ACM Transactions on Programming Languages and Systems (TOPLAS),
// 14(2):127-144, New York, NY, USA, 4/1992. ACM press.) See QuoRem for T-division
// and modulus (like Go).

// DivMod 在 y != 0 时，将 z 置为 x 除以 y 的商，将 m
// 置为 x 取模 y 的模数并返回值对 (z, m)。 若 y ==
// 0，就会产生一个除以零的运行时派错。
//
// DivMod 实现了截断式除法和取模（与Go不同）：
//
//	q = x div y // 使得
//	m = x - y*q // 其中
//	0 <= m < |q|
//
// （详见 Raymond T. Boute，《函数 div 和 mod
// 的欧氏定义》以及《ACM编程语言与系统会议记录》
// （TOPLAS），14(2):127-144, New York, NY,
// USA, 4/1992. ACM 出版社。） 截断式除法和取模（与Go相同）见
// QuoRem。
func (z *Int) DivMod(x, y, m *Int) (*Int, *Int)

// Exp sets z = x**y mod |m| (i.e. the sign of m is ignored), and returns z. If y
// <= 0, the result is 1 mod |m|; if m == nil or m == 0, z = x**y. See Knuth,
// volume 2, section 4.6.3.

// Exp 置 z = x**y mod |m|（换言之，m 的符号被忽略），并返回
// z。 若 y <=0，则其结果为 1，若 m == nil 或 m == 0，则
// z = x**y。 见 Knuth《计算机程序设计艺术》，卷 2，章节
// 4.6.3。
func (z *Int) Exp(x, y, m *Int) *Int

// Format is a support routine for fmt.Formatter. It accepts the formats 'b'
// (binary), 'o' (octal), 'd' (decimal), 'x' (lowercase hexadecimal), and 'X'
// (uppercase hexadecimal). Also supported are the full suite of package fmt's
// format verbs for integral types, including '+', '-', and ' ' for sign control,
// '#' for leading zero in octal and for hexadecimal, a leading "0x" or "0X" for
// "%#x" and "%#X" respectively, specification of minimum digits precision, output
// field width, space or zero padding, and left or right justification.

// Format 是 fmt.Formatter 的一个支持函数。它接受
// 'b'（二进制）、'o'（八进制）、 'd'（十进制）、'x'（小写十六进制）和
// 'X'（大写十六进制）的格式。也同样支持 fmt
// 包的一整套类型的格式占位符，包括用于符号控制的 '+'、'-' 和 '
// '，用于八进制前导零的 '#'，分别用于十六进制前导 "0x" 或 "0X" 的
// "%#x" 和 "%#X"，用于最小数字精度的规范，
// 输出字段的宽度，空格或零的填充，以及左右对齐。
func (x *Int) Format(s fmt.State, ch rune)

// GCD sets z to the greatest common divisor of a and b, which both must be > 0,
// and returns z. If x and y are not nil, GCD sets x and y such that z = a*x + b*y.
// If either a or b is <= 0, GCD sets z = x = y = 0.

// GCD 将 z 置为 a 和 b 的最大公约数，二者必须均 > 0，并返回 z。
// 若 x 或 y 非 nil，GCD 会设置 x 与 y 的值使得 z = a*x
// + b*y。 若 a 或 b <= 0，GCD就会置 z = x = y =
// 0。
func (z *Int) GCD(x, y, a, b *Int) *Int

// GobDecode implements the gob.GobDecoder interface.

// GobDecode 实现了 gob.GobDecoder 接口。
func (z *Int) GobDecode(buf []byte) error

// GobEncode implements the gob.GobEncoder interface.

// GobEncode 实现了 gob.GobEncoder 接口。
func (x *Int) GobEncode() ([]byte, error)

// Int64 returns the int64 representation of x. If x cannot be represented in an
// int64, the result is undefined.

// Int64 返回 x 的 int64 表示。 若 x 不能被表示为
// int64，则其结果是未定义的。
func (x *Int) Int64() int64

// Lsh sets z = x << n and returns z.

// Lsh 置 z = x << n 并返回 z。
func (z *Int) Lsh(x *Int, n uint) *Int

// MarshalJSON implements the json.Marshaler interface.

// MarshalJSON 实现了 json.Marshaler 接口。
func (z *Int) MarshalJSON() ([]byte, error)

// MarshalText implements the encoding.TextMarshaler interface.

// MarshalText implements the
// encoding.TextMarshaler interface.
func (z *Int) MarshalText() (text []byte, err error)

// Mod sets z to the modulus x%y for y != 0 and returns z. If y == 0, a
// division-by-zero run-time panic occurs. Mod implements Euclidean modulus (unlike
// Go); see DivMod for more details.

// Mod 在 y != 0 时，将 z 置为 x%y 的余数并返回 z。 若 y
// == 0，就会产生一个除以零的运行时派错。 Mod
// 实现了欧氏取模（与Go不同），更多详情见 DivMod。
func (z *Int) Mod(x, y *Int) *Int

// ModInverse sets z to the multiplicative inverse of g in the ring ℤ/nℤ and
// returns z. If g and n are not relatively prime, the result is undefined.

// ModInverse 将 z 置为 g 在环 ℤ/nℤ 中的乘法逆元素并返回
// z。若 g 与 n 并不互质，则结果为未定义。
func (z *Int) ModInverse(g, n *Int) *Int

// Mul sets z to the product x*y and returns z.

// Mul 将 z 置为 x*y 的积并返回 z。
func (z *Int) Mul(x, y *Int) *Int

// MulRange sets z to the product of all integers in the range [a, b] inclusively
// and returns z. If a > b (empty range), the result is 1.

// MulRange 将 z 置为闭区间 [a, b] 内所有整数的积并返回 z。
// 若 a > b（空区间），则其结果为 1。
func (z *Int) MulRange(a, b int64) *Int

// Neg sets z to -x and returns z.

// Neg 将 z 置为 -x 并返回 z。
func (z *Int) Neg(x *Int) *Int

// Not sets z = ^x and returns z.

// Not 置 z = ^x 并返回 z。
func (z *Int) Not(x *Int) *Int

// Or sets z = x | y and returns z.

// Or 置 z = x | y 并返回 z。
func (z *Int) Or(x, y *Int) *Int

// ProbablyPrime performs n Miller-Rabin tests to check whether x is prime. If it
// returns true, x is prime with probability 1 - 1/4^n. If it returns false, x is
// not prime.

// ProbablyPrime 通过执行 n 次 Miller-Rabin
// 测试来检查 x 是否为质数。 若它返回 true，x 有 1 - 1/4^n
// 的可能性为质数。 若它返回 false，则 x 不是质数。x 必须 >0。
func (x *Int) ProbablyPrime(n int) bool

// Quo sets z to the quotient x/y for y != 0 and returns z. If y == 0, a
// division-by-zero run-time panic occurs. Quo implements truncated division (like
// Go); see QuoRem for more details.

// Quo 在 y != 0 时，将 z 置为 x/y 的商并返回 z。 若 y
// == 0，就会产生一个除以零的运行时派错。 Quo
// 实现了截断式除法（与Go相同），更多详情见 QuoRem。
func (z *Int) Quo(x, y *Int) *Int

// QuoRem sets z to the quotient x/y and r to the remainder x%y and returns the
// pair (z, r) for y != 0. If y == 0, a division-by-zero run-time panic occurs.
//
// QuoRem implements T-division and modulus (like Go):
//
//	q = x/y      with the result truncated to zero
//	r = x - y*q
//
// (See Daan Leijen, ``Division and Modulus for Computer Scientists''.) See DivMod
// for Euclidean division and modulus (unlike Go).

// QuoRem 在 y != 0 时，将 z 置为 x/y 的商，将 r 置为
// x%y 的余数并返回值对 (z, r)。 若 y ==
// 0，就会产生一个除以零的运行时派错。
//
// QuoRem 实现了截断式除法和取模（与Go相同）：
//
//	q = x/y      // 其结果向零截断
//	r = x - y*q
//
// （详见 Daan Leijen，《计算机科学家的除法和取模》。）
// 欧氏除法和取模（与Go不同）见 DivMod。
func (z *Int) QuoRem(x, y, r *Int) (*Int, *Int)

// Rand sets z to a pseudo-random number in [0, n) and returns z.

// Rand 将 z 置为区间 [0, n) 中的一个伪随机数并返回 z。
func (z *Int) Rand(rnd *rand.Rand, n *Int) *Int

// Rem sets z to the remainder x%y for y != 0 and returns z. If y == 0, a
// division-by-zero run-time panic occurs. Rem implements truncated modulus (like
// Go); see QuoRem for more details.

// Rem 在 y != 0 时，将 z 置为 x%y 的余数并返回 z。 若 y
// == 0，就会产生一个除以零的运行时派错。 Rem
// 实现了截断式取模（与Go相同），更多详情见 QuoRem。
func (z *Int) Rem(x, y *Int) *Int

// Rsh sets z = x >> n and returns z.

// Rsh 置 z = x >> n 并返回 z。
func (z *Int) Rsh(x *Int, n uint) *Int

// Scan is a support routine for fmt.Scanner; it sets z to the value of the scanned
// number. It accepts the formats 'b' (binary), 'o' (octal), 'd' (decimal), 'x'
// (lowercase hexadecimal), and 'X' (uppercase hexadecimal).

// Scan 是 fmt.Scanner 的一个支持函数；它将 z
// 置为已扫描数字的值。它接受格式'b'（二进制）、
// 'o'（八进制）、'd'（十进制）、'x'（小写十六进制）及'X'（大写十六进制）。
func (z *Int) Scan(s fmt.ScanState, ch rune) error

// Set sets z to x and returns z.

// Set 将 z 置为 x 并返回 z。
func (z *Int) Set(x *Int) *Int

// SetBit sets z to x, with x's i'th bit set to b (0 or 1). That is, if b is 1
// SetBit sets z = x | (1 << i); if b is 0 SetBit sets z = x &^ (1 << i). If b is
// not 0 or 1, SetBit will panic.

// SetBit 将 z 置为 x，将 x 的第 i 位置为 b（0 或 1）。
// 换言之，若 b 为 1，SetBit 会置 z = x | (1 << i)；若
// b 为 0，SetBit 会置 z = x &^ (1 << i)。若 b 非
// 0 或 1，SetBit 就会引发派错。
func (z *Int) SetBit(x *Int, i int, b uint) *Int

// SetBits provides raw (unchecked but fast) access to z by setting its value to
// abs, interpreted as a little-endian Word slice, and returning z. The result and
// abs share the same underlying array. SetBits is intended to support
// implementation of missing low-level Int functionality outside this package; it
// should be avoided otherwise.

// SetBits 提供了对 z 的原始访问（未经检查但很快）。它通过将其值设为
// abs，解释为小端序的 Word 切片，并返回 z 来实现。SetBits
// 旨在支持此包外缺失的底层 Int 功能的实现，除此之外应尽量避免。
func (z *Int) SetBits(abs []Word) *Int

// SetBytes interprets buf as the bytes of a big-endian unsigned integer, sets z to
// that value, and returns z.

// SetBytes 将 buf 解释为大端序的无符号整数字节，置 z 为该值后返回
// z。
func (z *Int) SetBytes(buf []byte) *Int

// SetInt64 sets z to x and returns z.

// SetInt64 将 z 置为 x 并返回 z。
func (z *Int) SetInt64(x int64) *Int

// SetString sets z to the value of s, interpreted in the given base, and returns z
// and a boolean indicating success. If SetString fails, the value of z is
// undefined but the returned value is nil.
//
// The base argument must be 0 or a value from 2 through MaxBase. If the base is 0,
// the string prefix determines the actual conversion base. A prefix of ``0x'' or
// ``0X'' selects base 16; the ``0'' prefix selects base 8, and a ``0b'' or ``0B''
// prefix selects base 2. Otherwise the selected base is 10.

// SetString 将 z 置为 s 的值，按给定的进制 base 解释并返回
// z 和一个指示是否成功的布尔值。 若 SetString 失败，则 z
// 的值是未定义的，其返回值则为 nil。
//
// 进制实参 base 必须为 0 或从 2 到 MaxBase 的值。若 base
// 为 0，则其实际的转换进制由
// 该字符串的前缀决定。前缀“0x”或“0X”会选择16进制，前缀“0”会选择8进制，前缀“0b”或“0B”
// 会选择2进制。其它情况则选择10进制。
func (z *Int) SetString(s string, base int) (*Int, bool)

// SetUint64 sets z to x and returns z.

// SetUint64 将 z 置为 x 并返回 z。
func (z *Int) SetUint64(x uint64) *Int

// Sign returns:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0

// 符号返回：
//
//	若 x <  0 则为 -1
//	若 x == 0 则为  0
//	若 x >  0 则为 +1
func (x *Int) Sign() int

func (x *Int) String() string

// Sub sets z to the difference x-y and returns z.

// Sub 将 z 置为 x-y 的差并返回 z。
func (z *Int) Sub(x, y *Int) *Int

// Uint64 returns the uint64 representation of x. If x cannot be represented in a
// uint64, the result is undefined.

// Uint64 返回 x 的 uint64 表示。 若 x 不能被表示为
// uint64，则其结果是未定义的。
func (x *Int) Uint64() uint64

// UnmarshalJSON implements the json.Unmarshaler interface.

// UnmarshalJSON 实现了 json.Unmarshaler 接口。
func (z *Int) UnmarshalJSON(text []byte) error

// UnmarshalText implements the encoding.TextUnmarshaler interface.

// UnmarshalText implements the
// encoding.TextUnmarshaler interface.
func (z *Int) UnmarshalText(text []byte) error

// Xor sets z = x ^ y and returns z.

// Xor 置 z = x ^ y 并返回 z。
func (z *Int) Xor(x, y *Int) *Int

// A Rat represents a quotient a/b of arbitrary precision. The zero value for a Rat
// represents the value 0.

// A Rat represents a quotient a/b of
// arbitrary precision. The zero value for
// a Rat represents the value 0.
type Rat struct {
	// contains filtered or unexported fields
}

// NewRat creates a new Rat with numerator a and denominator b.

// NewRat creates a new Rat with numerator
// a and denominator b.
func NewRat(a, b int64) *Rat

// Abs sets z to |x| (the absolute value of x) and returns z.

// Abs sets z to |x| (the absolute value of
// x) and returns z.
func (z *Rat) Abs(x *Rat) *Rat

// Add sets z to the sum x+y and returns z.

// Add sets z to the sum x+y and returns z.
func (z *Rat) Add(x, y *Rat) *Rat

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func (x *Rat) Cmp(y *Rat) int

// Denom returns the denominator of x; it is always > 0. The result is a reference
// to x's denominator; it may change if a new value is assigned to x, and vice
// versa.

// Denom returns the denominator of x; it
// is always > 0. The result is a reference
// to x's denominator; it may change if a
// new value is assigned to x, and vice
// versa.
func (x *Rat) Denom() *Int

// Float32 returns the nearest float32 value for x and a bool indicating whether f
// represents x exactly. If the magnitude of x is too large to be represented by a
// float32, f is an infinity and exact is false. The sign of f always matches the
// sign of x, even if f == 0.

// Float32 returns the nearest float32
// value for x and a bool indicating
// whether f represents x exactly. If the
// magnitude of x is too large to be
// represented by a float32, f is an
// infinity and exact is false. The sign of
// f always matches the sign of x, even if
// f == 0.
func (x *Rat) Float32() (f float32, exact bool)

// Float64 returns the nearest float64 value for x and a bool indicating whether f
// represents x exactly. If the magnitude of x is too large to be represented by a
// float64, f is an infinity and exact is false. The sign of f always matches the
// sign of x, even if f == 0.

// Float64 returns the nearest float64
// value for x and a bool indicating
// whether f represents x exactly. If the
// magnitude of x is too large to be
// represented by a float64, f is an
// infinity and exact is false. The sign of
// f always matches the sign of x, even if
// f == 0.
func (x *Rat) Float64() (f float64, exact bool)

// FloatString returns a string representation of x in decimal form with prec
// digits of precision after the decimal point and the last digit rounded.

// FloatString returns a string
// representation of x in decimal form with
// prec digits of precision after the
// decimal point and the last digit
// rounded.
func (x *Rat) FloatString(prec int) string

// GobDecode implements the gob.GobDecoder interface.

// GobDecode implements the gob.GobDecoder
// interface.
func (z *Rat) GobDecode(buf []byte) error

// GobEncode implements the gob.GobEncoder interface.

// GobEncode implements the gob.GobEncoder
// interface.
func (x *Rat) GobEncode() ([]byte, error)

// Inv sets z to 1/x and returns z.

// Inv sets z to 1/x and returns z.
func (z *Rat) Inv(x *Rat) *Rat

// IsInt returns true if the denominator of x is 1.

// IsInt returns true if the denominator of
// x is 1.
func (x *Rat) IsInt() bool

// MarshalText implements the encoding.TextMarshaler interface.

// MarshalText implements the
// encoding.TextMarshaler interface.
func (r *Rat) MarshalText() (text []byte, err error)

// Mul sets z to the product x*y and returns z.

// Mul sets z to the product x*y and
// returns z.
func (z *Rat) Mul(x, y *Rat) *Rat

// Neg sets z to -x and returns z.

// Neg sets z to -x and returns z.
func (z *Rat) Neg(x *Rat) *Rat

// Num returns the numerator of x; it may be <= 0. The result is a reference to x's
// numerator; it may change if a new value is assigned to x, and vice versa. The
// sign of the numerator corresponds to the sign of x.

// Num returns the numerator of x; it may
// be <= 0. The result is a reference to
// x's numerator; it may change if a new
// value is assigned to x, and vice versa.
// The sign of the numerator corresponds to
// the sign of x.
func (x *Rat) Num() *Int

// Quo sets z to the quotient x/y and returns z. If y == 0, a division-by-zero
// run-time panic occurs.

// Quo sets z to the quotient x/y and
// returns z. If y == 0, a division-by-zero
// run-time panic occurs.
func (z *Rat) Quo(x, y *Rat) *Rat

// RatString returns a string representation of x in the form "a/b" if b != 1, and
// in the form "a" if b == 1.

// RatString returns a string
// representation of x in the form "a/b" if
// b != 1, and in the form "a" if b == 1.
func (x *Rat) RatString() string

// Scan is a support routine for fmt.Scanner. It accepts the formats 'e', 'E', 'f',
// 'F', 'g', 'G', and 'v'. All formats are equivalent.

// Scan is a support routine for
// fmt.Scanner. It accepts the formats 'e',
// 'E', 'f', 'F', 'g', 'G', and 'v'. All
// formats are equivalent.
func (z *Rat) Scan(s fmt.ScanState, ch rune) error

// Set sets z to x (by making a copy of x) and returns z.

// Set sets z to x (by making a copy of x)
// and returns z.
func (z *Rat) Set(x *Rat) *Rat

// SetFloat64 sets z to exactly f and returns z. If f is not finite, SetFloat
// returns nil.

// SetFloat64 sets z to exactly f and
// returns z. If f is not finite, SetFloat
// returns nil.
func (z *Rat) SetFloat64(f float64) *Rat

// SetFrac sets z to a/b and returns z.

// SetFrac sets z to a/b and returns z.
func (z *Rat) SetFrac(a, b *Int) *Rat

// SetFrac64 sets z to a/b and returns z.

// SetFrac64 sets z to a/b and returns z.
func (z *Rat) SetFrac64(a, b int64) *Rat

// SetInt sets z to x (by making a copy of x) and returns z.

// SetInt sets z to x (by making a copy of
// x) and returns z.
func (z *Rat) SetInt(x *Int) *Rat

// SetInt64 sets z to x and returns z.

// SetInt64 sets z to x and returns z.
func (z *Rat) SetInt64(x int64) *Rat

// SetString sets z to the value of s and returns z and a boolean indicating
// success. s can be given as a fraction "a/b" or as a floating-point number
// optionally followed by an exponent. If the operation failed, the value of z is
// undefined but the returned value is nil.

// SetString sets z to the value of s and
// returns z and a boolean indicating
// success. s can be given as a fraction
// "a/b" or as a floating-point number
// optionally followed by an exponent. If
// the operation failed, the value of z is
// undefined but the returned value is nil.
func (z *Rat) SetString(s string) (*Rat, bool)

// Sign returns:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0

// Sign returns:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0
func (x *Rat) Sign() int

// String returns a string representation of x in the form "a/b" (even if b == 1).

// String returns a string representation
// of x in the form "a/b" (even if b == 1).
func (x *Rat) String() string

// Sub sets z to the difference x-y and returns z.

// Sub sets z to the difference x-y and
// returns z.
func (z *Rat) Sub(x, y *Rat) *Rat

// UnmarshalText implements the encoding.TextUnmarshaler interface.

// UnmarshalText implements the
// encoding.TextUnmarshaler interface.
func (r *Rat) UnmarshalText(text []byte) error

// A Word represents a single digit of a multi-precision unsigned integer.

// Word 表示多精度无符号整数的单个数字。
type Word uintptr
