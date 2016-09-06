// +build ingore

package big

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

// Constants describing the Accuracy of a Float.
const (
	Below Accuracy = -1
	Exact Accuracy = 0
	Above Accuracy = +1
)

// MaxBase is the largest number base accepted for string conversions.
const MaxBase = 'z' - 'a' + 10 + 1

// Exponent and precision limits.
const (
	MaxExp  = math.MaxInt32  // largest supported exponent
	MinExp  = math.MinInt32  // smallest supported exponent
	MaxPrec = math.MaxUint32 // largest (theoretically) supported precision; likely memory-limited
)

// These constants define supported rounding modes.
const (
	ToNearestEven RoundingMode = iota // == IEEE 754-2008 roundTiesToEven
	ToNearestAway                     // == IEEE 754-2008 roundTiesToAway
	ToZero                            // == IEEE 754-2008 roundTowardZero
	AwayFromZero                      // no IEEE 754-2008 equivalent
	ToNegativeInf                     // == IEEE 754-2008 roundTowardNegative
	ToPositiveInf                     // == IEEE 754-2008 roundTowardPositive
)

// Accuracy describes the rounding error produced by the most recent
// operation that generated a Float value, relative to the exact value.
type Accuracy int8

// An ErrNaN panic is raised by a Float operation that would lead to
// a NaN under IEEE-754 rules. An ErrNaN implements the error interface.
type ErrNaN struct {
}

// A nonzero finite Float represents a multi-precision floating point number
//
// 	sign × mantissa × 2**exponent
//
// with 0.5 <= mantissa < 1.0, and MinExp <= exponent <= MaxExp. A Float may
// also be zero (+0, -0) or infinite (+Inf, -Inf). All Floats are ordered, and
// the ordering of two Floats x and y is defined by x.Cmp(y).
//
// Each Float value also has a precision, rounding mode, and accuracy. The
// precision is the maximum number of mantissa bits available to represent the
// value. The rounding mode specifies how a result should be rounded to fit into
// the mantissa bits, and accuracy describes the rounding error with respect to
// the exact result.
//
// Unless specified otherwise, all operations (including setters) that specify a
// *Float variable for the result (usually via the receiver with the exception
// of MantExp), round the numeric result according to the precision and rounding
// mode of the result variable.
//
// If the provided result precision is 0 (see below), it is set to the precision
// of the argument with the largest precision value before any rounding takes
// place, and the rounding mode remains unchanged. Thus, uninitialized Floats
// provided as result arguments will have their precision set to a reasonable
// value determined by the operands and their mode is the zero value for
// RoundingMode (ToNearestEven).
//
// By setting the desired precision to 24 or 53 and using matching rounding mode
// (typically ToNearestEven), Float operations produce the same results as the
// corresponding float32 or float64 IEEE-754 arithmetic for operands that
// correspond to normal (i.e., not denormal) float32 or float64 numbers.
// Exponent underflow and overflow lead to a 0 or an Infinity for different
// values than IEEE-754 because Float exponents have a much larger range.
//
// The zero (uninitialized) value for a Float is ready to use and represents the
// number +0.0 exactly, with precision 0 and rounding mode ToNearestEven.
type Float struct {
}

// An Int represents a signed multi-precision integer.
// The zero value for an Int represents the value 0.
type Int struct {
}

// A Rat represents a quotient a/b of arbitrary precision.
// The zero value for a Rat represents the value 0.
type Rat struct {
}

// RoundingMode determines how a Float value is rounded to the
// desired precision. Rounding may change the Float value; the
// rounding error is described by the Float's Accuracy.
type RoundingMode byte

// A Word represents a single digit of a multi-precision unsigned integer.
type Word uintptr

// Jacobi returns the Jacobi symbol (x/y), either +1, -1, or 0.
// The y argument must be an odd integer.
func Jacobi(x, y *Int) int

// NewFloat allocates and returns a new Float set to x,
// with precision 53 and rounding mode ToNearestEven.
// NewFloat panics with ErrNaN if x is a NaN.
func NewFloat(x float64) *Float

// NewInt allocates and returns a new Int set to x.
func NewInt(x int64) *Int

// NewRat creates a new Rat with numerator a and denominator b.
func NewRat(a, b int64) *Rat

// ParseFloat is like f.Parse(s, base) with f set to the given precision
// and rounding mode.
func ParseFloat(s string, base int, prec uint, mode RoundingMode) (f *Float, b int, err error)

// Abs sets z to the (possibly rounded) value |x| (the absolute value of x)
// and returns z.
func (z *Float) Abs(x *Float) *Float

// Acc returns the accuracy of x produced by the most recent operation.
func (x *Float) Acc() Accuracy

// Add sets z to the rounded sum x+y and returns z. If z's precision is 0, it is
// changed to the larger of x's or y's precision before the operation. Rounding
// is performed according to z's precision and rounding mode; and z's accuracy
// reports the result error relative to the exact (not rounded) result. Add
// panics with ErrNaN if x and y are infinities with opposite signs. The value
// of z is undefined in that case.
//
// BUG(gri) When rounding ToNegativeInf, the sign of Float values rounded to 0
// is incorrect.
func (z *Float) Add(x, y *Float) *Float

// Append appends to buf the string form of the floating-point number x,
// as generated by x.Text, and returns the extended buffer.
func (x *Float) Append(buf []byte, fmt byte, prec int) []byte

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y (incl. -0 == 0, -Inf == -Inf, and +Inf == +Inf)
//   +1 if x >  y
func (x *Float) Cmp(y *Float) int

// Copy sets z to x, with the same precision, rounding mode, and
// accuracy as x, and returns z. x is not changed even if z and
// x are the same.
func (z *Float) Copy(x *Float) *Float

// Float32 returns the float32 value nearest to x. If x is too small to be
// represented by a float32 (|x| < math.SmallestNonzeroFloat32), the result
// is (0, Below) or (-0, Above), respectively, depending on the sign of x.
// If x is too large to be represented by a float32 (|x| > math.MaxFloat32),
// the result is (+Inf, Above) or (-Inf, Below), depending on the sign of x.
func (x *Float) Float32() (float32, Accuracy)

// Float64 returns the float64 value nearest to x. If x is too small to be
// represented by a float64 (|x| < math.SmallestNonzeroFloat64), the result
// is (0, Below) or (-0, Above), respectively, depending on the sign of x.
// If x is too large to be represented by a float64 (|x| > math.MaxFloat64),
// the result is (+Inf, Above) or (-Inf, Below), depending on the sign of x.
func (x *Float) Float64() (float64, Accuracy)

// Format implements fmt.Formatter. It accepts all the regular
// formats for floating-point numbers ('b', 'e', 'E', 'f', 'F',
// 'g', 'G') as well as 'p' and 'v'. See (*Float).Text for the
// interpretation of 'p'. The 'v' format is handled like 'g'.
// Format also supports specification of the minimum precision
// in digits, the output field width, as well as the format flags
// '+' and ' ' for sign control, '0' for space or zero padding,
// and '-' for left or right justification. See the fmt package
// for details.

// Format implements fmt.Formatter. It accepts all the regular
// formats for floating-point numbers ('e', 'E', 'f', 'F', 'g',
// 'G') as well as 'b', 'p', and 'v'. See (*Float).Text for the
// interpretation of 'b' and 'p'. The 'v' format is handled like
// 'g'.
// Format also supports specification of the minimum precision
// in digits, the output field width, as well as the format verbs
// '+' and ' ' for sign control, '0' for space or zero padding,
// and '-' for left or right justification. See the fmt package
// for details.
func (x *Float) Format(s fmt.State, format rune)

// Int returns the result of truncating x towards zero;
// or nil if x is an infinity.
// The result is Exact if x.IsInt(); otherwise it is Below
// for x > 0, and Above for x < 0.
// If a non-nil *Int argument z is provided, Int stores
// the result in z instead of allocating a new Int.
func (x *Float) Int(z *Int) (*Int, Accuracy)

// Int64 returns the integer resulting from truncating x towards zero.
// If math.MinInt64 <= x <= math.MaxInt64, the result is Exact if x is
// an integer, and Above (x < 0) or Below (x > 0) otherwise.
// The result is (math.MinInt64, Above) for x < math.MinInt64,
// and (math.MaxInt64, Below) for x > math.MaxInt64.
func (x *Float) Int64() (int64, Accuracy)

// IsInf reports whether x is +Inf or -Inf.
func (x *Float) IsInf() bool

// IsInt reports whether x is an integer.
// ±Inf values are not integers.
func (x *Float) IsInt() bool

// MantExp breaks x into its mantissa and exponent components
// and returns the exponent. If a non-nil mant argument is
// provided its value is set to the mantissa of x, with the
// same precision and rounding mode as x. The components
// satisfy x == mant × 2**exp, with 0.5 <= |mant| < 1.0.
// Calling MantExp with a nil argument is an efficient way to
// get the exponent of the receiver.
//
// Special cases are:
//
// 	(  ±0).MantExp(mant) = 0, with mant set to   ±0
// 	(±Inf).MantExp(mant) = 0, with mant set to ±Inf
//
// x and mant may be the same in which case x is set to its
// mantissa value.

// MantExp breaks x into its mantissa and exponent components
// and returns the exponent. If a non-nil mant argument is
// provided its value is set to the mantissa of x, with the
// same precision and rounding mode as x. The components
// satisfy x == mant × 2**exp, with 0.5 <= |mant| < 1.0.
// Calling MantExp with a nil argument is an efficient way to
// get the exponent of the receiver.
//
// Special cases are:
//
//     (  ±0).MantExp(mant) = 0, with mant set to   ±0
//     (±Inf).MantExp(mant) = 0, with mant set to ±Inf
//
// x and mant may be the same in which case x is set to its
// mantissa value.
func (x *Float) MantExp(mant *Float) (exp int)

// MarshalText implements the encoding.TextMarshaler interface.
// Only the Float value is marshaled (in full precision), other
// attributes such as precision or accuracy are ignored.
func (x *Float) MarshalText() (text []byte, err error)

// MinPrec returns the minimum precision required to represent x exactly
// (i.e., the smallest prec before x.SetPrec(prec) would start rounding x).
// The result is 0 for |x| == 0 and |x| == Inf.
func (x *Float) MinPrec() uint

// Mode returns the rounding mode of x.
func (x *Float) Mode() RoundingMode

// Mul sets z to the rounded product x*y and returns z.
// Precision, rounding, and accuracy reporting are as for Add.
// Mul panics with ErrNaN if one operand is zero and the other
// operand an infinity. The value of z is undefined in that case.
func (z *Float) Mul(x, y *Float) *Float

// Neg sets z to the (possibly rounded) value of x with its sign negated,
// and returns z.
func (z *Float) Neg(x *Float) *Float

// Parse parses s which must contain a text representation of a floating-
// point number with a mantissa in the given conversion base (the exponent
// is always a decimal number), or a string representing an infinite value.
//
// It sets z to the (possibly rounded) value of the corresponding floating-
// point value, and returns z, the actual base b, and an error err, if any.
// If z's precision is 0, it is changed to 64 before rounding takes effect.
// The number must be of the form:
//
// 	number   = [ sign ] [ prefix ] mantissa [ exponent ] | infinity .
// 	sign     = "+" | "-" .
//      prefix   = "0" ( "x" | "X" | "b" | "B" ) .
// 	mantissa = digits | digits "." [ digits ] | "." digits .
// 	exponent = ( "E" | "e" | "p" ) [ sign ] digits .
// 	digits   = digit { digit } .
// 	digit    = "0" ... "9" | "a" ... "z" | "A" ... "Z" .
//      infinity = [ sign ] ( "inf" | "Inf" ) .
//
// The base argument must be 0, 2, 10, or 16. Providing an invalid base
// argument will lead to a run-time panic.
//
// For base 0, the number prefix determines the actual base: A prefix of
// "0x" or "0X" selects base 16, and a "0b" or "0B" prefix selects
// base 2; otherwise, the actual base is 10 and no prefix is accepted.
// The octal prefix "0" is not supported (a leading "0" is simply
// considered a "0").
//
// A "p" exponent indicates a binary (rather then decimal) exponent;
// for instance "0x1.fffffffffffffp1023" (using base 0) represents the
// maximum float64 value. For hexadecimal mantissae, the exponent must
// be binary, if present (an "e" or "E" exponent indicator cannot be
// distinguished from a mantissa digit).
//
// The returned *Float f is nil and the value of z is valid but not
// defined if an error is reported.

// Parse parses s which must contain a text representation of a floating-
// point number with a mantissa in the given conversion base (the exponent
// is always a decimal number), or a string representing an infinite value.
//
// It sets z to the (possibly rounded) value of the corresponding floating-
// point value, and returns z, the actual base b, and an error err, if any.
// If z's precision is 0, it is changed to 64 before rounding takes effect.
// The number must be of the form:
//
//     number   = [ sign ] [ prefix ] mantissa [ exponent ] | infinity .
//     sign     = "+" | "-" .
//      prefix   = "0" ( "x" | "X" | "b" | "B" ) .
//     mantissa = digits | digits "." [ digits ] | "." digits .
//     exponent = ( "E" | "e" | "p" ) [ sign ] digits .
//     digits   = digit { digit } .
//     digit    = "0" ... "9" | "a" ... "z" | "A" ... "Z" .
//      infinity = [ sign ] ( "inf" | "Inf" ) .
//
// The base argument must be 0, 2, 10, or 16. Providing an invalid base
// argument will lead to a run-time panic.
//
// For base 0, the number prefix determines the actual base: A prefix of
// "0x" or "0X" selects base 16, and a "0b" or "0B" prefix selects
// base 2; otherwise, the actual base is 10 and no prefix is accepted.
// The octal prefix "0" is not supported (a leading "0" is simply
// considered a "0").
//
// A "p" exponent indicates a binary (rather then decimal) exponent;
// for instance "0x1.fffffffffffffp1023" (using base 0) represents the
// maximum float64 value. For hexadecimal mantissae, the exponent must
// be binary, if present (an "e" or "E" exponent indicator cannot be
// distinguished from a mantissa digit).
//
// The returned *Float f is nil and the value of z is valid but not
// defined if an error is reported.
func (z *Float) Parse(s string, base int) (f *Float, b int, err error)

// Prec returns the mantissa precision of x in bits.
// The result may be 0 for |x| == 0 and |x| == Inf.
func (x *Float) Prec() uint

// Quo sets z to the rounded quotient x/y and returns z.
// Precision, rounding, and accuracy reporting are as for Add.
// Quo panics with ErrNaN if both operands are zero or infinities.
// The value of z is undefined in that case.
func (z *Float) Quo(x, y *Float) *Float

// Rat returns the rational number corresponding to x;
// or nil if x is an infinity.
// The result is Exact if x is not an Inf.
// If a non-nil *Rat argument z is provided, Rat stores
// the result in z instead of allocating a new Rat.
func (x *Float) Rat(z *Rat) (*Rat, Accuracy)

// Set sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to the precision of x
// before setting z (and rounding will have no effect).
// Rounding is performed according to z's precision and rounding
// mode; and z's accuracy reports the result error relative to the
// exact (not rounded) result.
func (z *Float) Set(x *Float) *Float

// SetFloat64 sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to 53 (and rounding will have
// no effect). SetFloat64 panics with ErrNaN if x is a NaN.
func (z *Float) SetFloat64(x float64) *Float

// SetInf sets z to the infinite Float -Inf if signbit is
// set, or +Inf if signbit is not set, and returns z. The
// precision of z is unchanged and the result is always
// Exact.
func (z *Float) SetInf(signbit bool) *Float

// SetInt sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to the larger of x.BitLen()
// or 64 (and rounding will have no effect).
func (z *Float) SetInt(x *Int) *Float

// SetInt64 sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to 64 (and rounding will have
// no effect).
func (z *Float) SetInt64(x int64) *Float

// SetMantExp sets z to mant × 2**exp and and returns z.
// The result z has the same precision and rounding mode
// as mant. SetMantExp is an inverse of MantExp but does
// not require 0.5 <= |mant| < 1.0. Specifically:
//
// 	mant := new(Float)
// 	new(Float).SetMantExp(mant, x.MantExp(mant)).Cmp(x) == 0
//
// Special cases are:
//
// 	z.SetMantExp(  ±0, exp) =   ±0
// 	z.SetMantExp(±Inf, exp) = ±Inf
//
// z and mant may be the same in which case z's exponent
// is set to exp.

// SetMantExp sets z to mant × 2**exp and and returns z.
// The result z has the same precision and rounding mode
// as mant. SetMantExp is an inverse of MantExp but does
// not require 0.5 <= |mant| < 1.0. Specifically:
//
//     mant := new(Float)
//     new(Float).SetMantExp(mant, x.MantExp(mant)).Cmp(x) == 0
//
// Special cases are:
//
//     z.SetMantExp(  ±0, exp) =   ±0
//     z.SetMantExp(±Inf, exp) = ±Inf
//
// z and mant may be the same in which case z's exponent
// is set to exp.
func (z *Float) SetMantExp(mant *Float, exp int) *Float

// SetMode sets z's rounding mode to mode and returns an exact z.
// z remains unchanged otherwise.
// z.SetMode(z.Mode()) is a cheap way to set z's accuracy to Exact.
func (z *Float) SetMode(mode RoundingMode) *Float

// SetPrec sets z's precision to prec and returns the (possibly) rounded value
// of z. Rounding occurs according to z's rounding mode if the mantissa cannot
// be represented in prec bits without loss of precision. SetPrec(0) maps all
// finite values to ±0; infinite values remain unchanged. If prec > MaxPrec, it
// is set to MaxPrec.
func (z *Float) SetPrec(prec uint) *Float

// SetRat sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to the largest of a.BitLen(),
// b.BitLen(), or 64; with x = a/b.
func (z *Float) SetRat(x *Rat) *Float

// SetString sets z to the value of s and returns z and a boolean indicating
// success. s must be a floating-point number of the same format as accepted
// by Parse, with base argument 0.
func (z *Float) SetString(s string) (*Float, bool)

// SetUint64 sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to 64 (and rounding will have
// no effect).
func (z *Float) SetUint64(x uint64) *Float

// Sign returns:
//
// 	-1 if x <   0
// 	 0 if x is ±0
// 	+1 if x >   0

// Sign returns:
//
//     -1 if x <   0
//      0 if x is ±0
//     +1 if x >   0
func (x *Float) Sign() int

// Signbit returns true if x is negative or negative zero.
func (x *Float) Signbit() bool

// String formats x like x.Text('g', 10). (String must be called explicitly,
// Float.Format does not support %s verb.)
func (x *Float) String() string

// Sub sets z to the rounded difference x-y and returns z.
// Precision, rounding, and accuracy reporting are as for Add.
// Sub panics with ErrNaN if x and y are infinities with equal
// signs. The value of z is undefined in that case.
func (z *Float) Sub(x, y *Float) *Float

// Text converts the floating-point number x to a string according to the given
// format and precision prec. The format is one of:
//
// 	'e'	-d.dddde±dd, decimal exponent, at least two (possibly 0) exponent digits
// 	'E'	-d.ddddE±dd, decimal exponent, at least two (possibly 0) exponent digits
// 	'f'	-ddddd.dddd, no exponent
// 	'g'	like 'e' for large exponents, like 'f' otherwise
// 	'G'	like 'E' for large exponents, like 'f' otherwise
// 	'b'	-ddddddp±dd, binary exponent
// 	'p'	-0x.dddp±dd, binary exponent, hexadecimal mantissa
//
// For the binary exponent formats, the mantissa is printed in normalized form:
//
// 	'b'	decimal integer mantissa using x.Prec() bits, or -0
// 	'p'	hexadecimal fraction with 0.5 <= 0.mantissa < 1.0, or -0
//
// If format is a different character, Text returns a "%" followed by the
// unrecognized format character.
//
// The precision prec controls the number of digits (excluding the exponent)
// printed by the 'e', 'E', 'f', 'g', and 'G' formats. For 'e', 'E', and 'f' it
// is the number of digits after the decimal point. For 'g' and 'G' it is the
// total number of digits. A negative precision selects the smallest number of
// decimal digits necessary to identify the value x uniquely using x.Prec()
// mantissa bits. The prec value is ignored for the 'b' or 'p' format.

// Text converts the floating-point number x to a string according to the given
// format and precision prec. The format is one of:
//
// 	'e'    -d.dddde±dd, decimal exponent, at least two (possibly 0) exponent digits
// 	'E'    -d.ddddE±dd, decimal exponent, at least two (possibly 0) exponent digits
// 	'f'    -ddddd.dddd, no exponent
// 	'g'    like 'e' for large exponents, like 'f' otherwise
// 	'G'    like 'E' for large exponents, like 'f' otherwise
// 	'b'    -ddddddp±dd, binary exponent
// 	'p'    -0x.dddp±dd, binary exponent, hexadecimal mantissa
//
// For the binary exponent formats, the mantissa is printed in normalized form:
//
// 	'b'    decimal integer mantissa using x.Prec() bits, or -0
// 	'p'    hexadecimal fraction with 0.5 <= 0.mantissa < 1.0, or -0
//
// If format is a different character, Text returns a "%" followed by the
// unrecognized format character.
//
// The precision prec controls the number of digits (excluding the exponent)
// printed by the 'e', 'E', 'f', 'g', and 'G' formats. For 'e', 'E', and 'f' it
// is the number of digits after the decimal point. For 'g' and 'G' it is the
// total number of digits. A negative precision selects the smallest number of
// decimal digits necessary to identify the value x uniquely using x.Prec()
// mantissa bits. The prec value is ignored for the 'b' or 'p' format.
func (x *Float) Text(format byte, prec int) string

// Uint64 returns the unsigned integer resulting from truncating x
// towards zero. If 0 <= x <= math.MaxUint64, the result is Exact
// if x is an integer and Below otherwise.
// The result is (0, Above) for x < 0, and (math.MaxUint64, Below)
// for x > math.MaxUint64.
func (x *Float) Uint64() (uint64, Accuracy)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The result is rounded per the precision and rounding mode of z.
// If z's precision is 0, it is changed to 64 before rounding takes
// effect.
func (z *Float) UnmarshalText(text []byte) error

// Abs sets z to |x| (the absolute value of x) and returns z.
func (z *Int) Abs(x *Int) *Int

// Add sets z to the sum x+y and returns z.
func (z *Int) Add(x, y *Int) *Int

// And sets z = x & y and returns z.
func (z *Int) And(x, y *Int) *Int

// AndNot sets z = x &^ y and returns z.
func (z *Int) AndNot(x, y *Int) *Int

// Append appends the string representation of x, as generated by
// x.Text(base), to buf and returns the extended buffer.
func (x *Int) Append(buf []byte, base int) []byte

// Binomial sets z to the binomial coefficient of (n, k) and returns z.
func (z *Int) Binomial(n, k int64) *Int

// Bit returns the value of the i'th bit of x. That is, it
// returns (x>>i)&1. The bit index i must be >= 0.
func (x *Int) Bit(i int) uint

// BitLen returns the length of the absolute value of x in bits.
// The bit length of 0 is 0.
func (x *Int) BitLen() int

// Bits provides raw (unchecked but fast) access to x by returning its
// absolute value as a little-endian Word slice. The result and x share
// the same underlying array.
// Bits is intended to support implementation of missing low-level Int
// functionality outside this package; it should be avoided otherwise.
func (x *Int) Bits() []Word

// Bytes returns the absolute value of x as a big-endian byte slice.
func (x *Int) Bytes() []byte

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
func (x *Int) Cmp(y *Int) (r int)

// Div sets z to the quotient x/y for y != 0 and returns z. If y == 0, a
// division-by-zero run-time panic occurs. Div implements Euclidean division
// (unlike Go); see DivMod for more details.
func (z *Int) Div(x, y *Int) *Int

// DivMod sets z to the quotient x div y and m to the modulus x mod y
// and returns the pair (z, m) for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
//
// DivMod implements Euclidean division and modulus (unlike Go):
//
// 	q = x div y  such that
// 	m = x - y*q  with 0 <= m < |y|
//
// (See Raymond T. Boute, ``The Euclidean definition of the functions
// div and mod''. ACM Transactions on Programming Languages and
// Systems (TOPLAS), 14(2):127-144, New York, NY, USA, 4/1992.
// ACM press.)
// See QuoRem for T-division and modulus (like Go).

// DivMod sets z to the quotient x div y and m to the modulus x mod y
// and returns the pair (z, m) for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
//
// DivMod implements Euclidean division and modulus (unlike Go):
//
//     q = x div y  such that
//     m = x - y*q  with 0 <= m < |y|
//
// (See Raymond T. Boute, ``The Euclidean definition of the functions
// div and mod''. ACM Transactions on Programming Languages and
// Systems (TOPLAS), 14(2):127-144, New York, NY, USA, 4/1992.
// ACM press.)
// See QuoRem for T-division and modulus (like Go).
func (z *Int) DivMod(x, y, m *Int) (*Int, *Int)

// Exp sets z = x**y mod |m| (i.e. the sign of m is ignored), and returns z.
// If y <= 0, the result is 1 mod |m|; if m == nil or m == 0, z = x**y.
// See Knuth, volume 2, section 4.6.3.
func (z *Int) Exp(x, y, m *Int) *Int

// Format implements fmt.Formatter. It accepts the formats
// 'b' (binary), 'o' (octal), 'd' (decimal), 'x' (lowercase
// hexadecimal), and 'X' (uppercase hexadecimal).
// Also supported are the full suite of package fmt's format
// flags for integral types, including '+' and ' ' for sign
// control, '#' for leading zero in octal and for hexadecimal,
// a leading "0x" or "0X" for "%#x" and "%#X" respectively,
// specification of minimum digits precision, output field
// width, space or zero padding, and '-' for left or right
// justification.

// Format is a support routine for fmt.Formatter. It accepts
// the formats 'b' (binary), 'o' (octal), 'd' (decimal), 'x'
// (lowercase hexadecimal), and 'X' (uppercase hexadecimal).
// Also supported are the full suite of package fmt's format
// verbs for integral types, including '+', '-', and ' '
// for sign control, '#' for leading zero in octal and for
// hexadecimal, a leading "0x" or "0X" for "%#x" and "%#X"
// respectively, specification of minimum digits precision,
// output field width, space or zero padding, and left or
// right justification.
func (x *Int) Format(s fmt.State, ch rune)

// GCD sets z to the greatest common divisor of a and b, which both must
// be > 0, and returns z.
// If x and y are not nil, GCD sets x and y such that z = a*x + b*y.
// If either a or b is <= 0, GCD sets z = x = y = 0.
func (z *Int) GCD(x, y, a, b *Int) *Int

// GobDecode implements the gob.GobDecoder interface.
func (z *Int) GobDecode(buf []byte) error

// GobEncode implements the gob.GobEncoder interface.
func (x *Int) GobEncode() ([]byte, error)

// Int64 returns the int64 representation of x.
// If x cannot be represented in an int64, the result is undefined.
func (x *Int) Int64() int64

// Lsh sets z = x << n and returns z.
func (z *Int) Lsh(x *Int, n uint) *Int

// MarshalJSON implements the json.Marshaler interface.
func (x *Int) MarshalJSON() ([]byte, error)

// MarshalText implements the encoding.TextMarshaler interface.
func (x *Int) MarshalText() (text []byte, err error)

// Mod sets z to the modulus x%y for y != 0 and returns z. If y == 0, a
// division-by-zero run-time panic occurs. Mod implements Euclidean modulus
// (unlike Go); see DivMod for more details.
func (z *Int) Mod(x, y *Int) *Int

// ModInverse sets z to the multiplicative inverse of g in the ring ℤ/nℤ and
// returns z. If g and n are not relatively prime, the result is undefined.
func (z *Int) ModInverse(g, n *Int) *Int

// ModSqrt sets z to a square root of x mod p if such a square root exists, and
// returns z. The modulus p must be an odd prime. If x is not a square mod p,
// ModSqrt leaves z unchanged and returns nil. This function panics if p is not
// an odd integer.
func (z *Int) ModSqrt(x, p *Int) *Int

// Mul sets z to the product x*y and returns z.
func (z *Int) Mul(x, y *Int) *Int

// MulRange sets z to the product of all integers
// in the range [a, b] inclusively and returns z.
// If a > b (empty range), the result is 1.
func (z *Int) MulRange(a, b int64) *Int

// Neg sets z to -x and returns z.
func (z *Int) Neg(x *Int) *Int

// Not sets z = ^x and returns z.
func (z *Int) Not(x *Int) *Int

// Or sets z = x | y and returns z.
func (z *Int) Or(x, y *Int) *Int

// ProbablyPrime performs n Miller-Rabin tests to check whether x is prime.
// If x is prime, it returns true.
// If x is not prime, it returns false with probability at least 1 - ¼ⁿ.
//
// It is not suitable for judging primes that an adversary may have crafted
// to fool this test.
func (x *Int) ProbablyPrime(n int) bool

// Quo sets z to the quotient x/y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Quo implements truncated division (like Go); see QuoRem for more details.
func (z *Int) Quo(x, y *Int) *Int

// QuoRem sets z to the quotient x/y and r to the remainder x%y
// and returns the pair (z, r) for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
//
// QuoRem implements T-division and modulus (like Go):
//
// 	q = x/y      with the result truncated to zero
// 	r = x - y*q
//
// (See Daan Leijen, ``Division and Modulus for Computer Scientists''.)
// See DivMod for Euclidean division and modulus (unlike Go).

// QuoRem sets z to the quotient x/y and r to the remainder x%y
// and returns the pair (z, r) for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
//
// QuoRem implements T-division and modulus (like Go):
//
//     q = x/y      with the result truncated to zero
//     r = x - y*q
//
// (See Daan Leijen, ``Division and Modulus for Computer Scientists''.)
// See DivMod for Euclidean division and modulus (unlike Go).
func (z *Int) QuoRem(x, y, r *Int) (*Int, *Int)

// Rand sets z to a pseudo-random number in [0, n) and returns z.
func (z *Int) Rand(rnd *rand.Rand, n *Int) *Int

// Rem sets z to the remainder x%y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Rem implements truncated modulus (like Go); see QuoRem for more details.
func (z *Int) Rem(x, y *Int) *Int

// Rsh sets z = x >> n and returns z.
func (z *Int) Rsh(x *Int, n uint) *Int

// Scan is a support routine for fmt.Scanner; it sets z to the value of the
// scanned number. It accepts the formats 'b' (binary), 'o' (octal), 'd'
// (decimal), 'x' (lowercase hexadecimal), and 'X' (uppercase hexadecimal).
func (z *Int) Scan(s fmt.ScanState, ch rune) error

// Set sets z to x and returns z.
func (z *Int) Set(x *Int) *Int

// SetBit sets z to x, with x's i'th bit set to b (0 or 1).
// That is, if b is 1 SetBit sets z = x | (1 << i);
// if b is 0 SetBit sets z = x &^ (1 << i). If b is not 0 or 1,
// SetBit will panic.
func (z *Int) SetBit(x *Int, i int, b uint) *Int

// SetBits provides raw (unchecked but fast) access to z by setting its
// value to abs, interpreted as a little-endian Word slice, and returning
// z. The result and abs share the same underlying array.
// SetBits is intended to support implementation of missing low-level Int
// functionality outside this package; it should be avoided otherwise.
func (z *Int) SetBits(abs []Word) *Int

// SetBytes interprets buf as the bytes of a big-endian unsigned
// integer, sets z to that value, and returns z.
func (z *Int) SetBytes(buf []byte) *Int

// SetInt64 sets z to x and returns z.
func (z *Int) SetInt64(x int64) *Int

// SetString sets z to the value of s, interpreted in the given base, and
// returns z and a boolean indicating success. If SetString fails, the value of
// z is undefined but the returned value is nil.
//
// The base argument must be 0 or a value between 2 and MaxBase. If the base is
// 0, the string prefix determines the actual conversion base. A prefix of
// ``0x'' or ``0X'' selects base 16; the ``0'' prefix selects base 8, and a
// ``0b'' or ``0B'' prefix selects base 2. Otherwise the selected base is 10.
func (z *Int) SetString(s string, base int) (*Int, bool)

// SetUint64 sets z to x and returns z.
func (z *Int) SetUint64(x uint64) *Int

// Sign returns:
//
// 	-1 if x <  0
// 	 0 if x == 0
// 	+1 if x >  0

// Sign returns:
//
//     -1 if x <  0
//      0 if x == 0
//     +1 if x >  0
func (x *Int) Sign() int

func (x *Int) String() string

// Sub sets z to the difference x-y and returns z.
func (z *Int) Sub(x, y *Int) *Int

// Text returns the string representation of x in the given base.
// Base must be between 2 and 36, inclusive. The result uses the
// lower-case letters 'a' to 'z' for digit values >= 10. No base
// prefix (such as "0x") is added to the string.
func (x *Int) Text(base int) string

// Uint64 returns the uint64 representation of x.
// If x cannot be represented in a uint64, the result is undefined.
func (x *Int) Uint64() uint64

// UnmarshalJSON implements the json.Unmarshaler interface.
func (z *Int) UnmarshalJSON(text []byte) error

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (z *Int) UnmarshalText(text []byte) error

// Xor sets z = x ^ y and returns z.
func (z *Int) Xor(x, y *Int) *Int

// Abs sets z to |x| (the absolute value of x) and returns z.
func (z *Rat) Abs(x *Rat) *Rat

// Add sets z to the sum x+y and returns z.
func (z *Rat) Add(x, y *Rat) *Rat

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
func (x *Rat) Cmp(y *Rat) int

// Denom returns the denominator of x; it is always > 0.
// The result is a reference to x's denominator; it
// may change if a new value is assigned to x, and vice versa.
func (x *Rat) Denom() *Int

// Float32 returns the nearest float32 value for x and a bool indicating
// whether f represents x exactly. If the magnitude of x is too large to
// be represented by a float32, f is an infinity and exact is false.
// The sign of f always matches the sign of x, even if f == 0.
func (x *Rat) Float32() (f float32, exact bool)

// Float64 returns the nearest float64 value for x and a bool indicating
// whether f represents x exactly. If the magnitude of x is too large to
// be represented by a float64, f is an infinity and exact is false.
// The sign of f always matches the sign of x, even if f == 0.
func (x *Rat) Float64() (f float64, exact bool)

// FloatString returns a string representation of x in decimal form with prec
// digits of precision after the decimal point. The last digit is rounded to
// nearest, with halves rounded away from zero.
func (x *Rat) FloatString(prec int) string

// GobDecode implements the gob.GobDecoder interface.
func (z *Rat) GobDecode(buf []byte) error

// GobEncode implements the gob.GobEncoder interface.
func (x *Rat) GobEncode() ([]byte, error)

// Inv sets z to 1/x and returns z.
func (z *Rat) Inv(x *Rat) *Rat

// IsInt reports whether the denominator of x is 1.
func (x *Rat) IsInt() bool

// MarshalText implements the encoding.TextMarshaler interface.
func (x *Rat) MarshalText() (text []byte, err error)

// Mul sets z to the product x*y and returns z.
func (z *Rat) Mul(x, y *Rat) *Rat

// Neg sets z to -x and returns z.
func (z *Rat) Neg(x *Rat) *Rat

// Num returns the numerator of x; it may be <= 0.
// The result is a reference to x's numerator; it
// may change if a new value is assigned to x, and vice versa.
// The sign of the numerator corresponds to the sign of x.
func (x *Rat) Num() *Int

// Quo sets z to the quotient x/y and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
func (z *Rat) Quo(x, y *Rat) *Rat

// RatString returns a string representation of x in the form "a/b" if b != 1,
// and in the form "a" if b == 1.
func (x *Rat) RatString() string

// Scan is a support routine for fmt.Scanner. It accepts the formats
// 'e', 'E', 'f', 'F', 'g', 'G', and 'v'. All formats are equivalent.
func (z *Rat) Scan(s fmt.ScanState, ch rune) error

// Set sets z to x (by making a copy of x) and returns z.
func (z *Rat) Set(x *Rat) *Rat

// SetFloat64 sets z to exactly f and returns z.
// If f is not finite, SetFloat returns nil.
func (z *Rat) SetFloat64(f float64) *Rat

// SetFrac sets z to a/b and returns z.
func (z *Rat) SetFrac(a, b *Int) *Rat

// SetFrac64 sets z to a/b and returns z.
func (z *Rat) SetFrac64(a, b int64) *Rat

// SetInt sets z to x (by making a copy of x) and returns z.
func (z *Rat) SetInt(x *Int) *Rat

// SetInt64 sets z to x and returns z.
func (z *Rat) SetInt64(x int64) *Rat

// SetString sets z to the value of s and returns z and a boolean indicating
// success. s can be given as a fraction "a/b" or as a floating-point number
// optionally followed by an exponent. If the operation failed, the value of
// z is undefined but the returned value is nil.
func (z *Rat) SetString(s string) (*Rat, bool)

// Sign returns:
//
// 	-1 if x <  0
// 	 0 if x == 0
// 	+1 if x >  0

// Sign returns:
//
//     -1 if x <  0
//      0 if x == 0
//     +1 if x >  0
func (x *Rat) Sign() int

// String returns a string representation of x in the form "a/b" (even if b ==
// 1).
func (x *Rat) String() string

// Sub sets z to the difference x-y and returns z.
func (z *Rat) Sub(x, y *Rat) *Rat

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (z *Rat) UnmarshalText(text []byte) error

func (i Accuracy) String() string

func (err ErrNaN) Error() string

func (i RoundingMode) String() string

