// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package math provides basic constants and mathematical functions.

// math 包提供了基本常数和数学函数。
package math

// Mathematical constants. Reference: http://oeis.org/Axxxxxx

// 数学常数。 参考：http://oeis.org/Axxxxxx
const (
	E   = 2.71828182845904523536028747135266249775724709369995957496696763 // A001113
	Pi  = 3.14159265358979323846264338327950288419716939937510582097494459 // A000796
	Phi = 1.61803398874989484820458683436563811772030917980576286213544862 // A001622

	Sqrt2   = 1.41421356237309504880168872420969807856967187537694807317667974 // A002193
	SqrtE   = 1.64872127070012814684865078781416357165377610071014801157507931 // A019774
	SqrtPi  = 1.77245385090551602729816748334114518279754945612238712821380779 // A002161
	SqrtPhi = 1.27201964951406896425242246173749149171560804184009624861664038 // A139339

	Ln2    = 0.693147180559945309417232121458176568075500134360255254120680009 // A002162
	Log2E  = 1 / Ln2
	Ln10   = 2.30258509299404568401799145468436420760110148862877297603332790 // A002392
	Log10E = 1 / Ln10
)

// Floating-point limit values. Max is the largest finite value representable by
// the type. SmallestNonzero is the smallest positive, non-zero value representable
// by the type.

// 浮点数极限值。 Max 为该类型可表示的最大有限值。
// SmallestNonzero 为该类型可表示的最小（非零）正值。
const (
	MaxFloat32             = 3.40282346638528859811704183484516925440e+38  // 2**127 * (2**24 - 1) / 2**23
	SmallestNonzeroFloat32 = 1.401298464324817070923729583289916131280e-45 // 1 / 2**(127 - 1 + 23)

	MaxFloat64             = 1.797693134862315708145274237317043567981e+308 // 2**1023 * (2**53 - 1) / 2**52
	SmallestNonzeroFloat64 = 4.940656458412465441765687928682213723651e-324 // 1 / 2**(1023 - 1 + 52)
)

// Integer limit values.

// 整数极限值。
const (
	MaxInt8   = 1<<7 - 1
	MinInt8   = -1 << 7
	MaxInt16  = 1<<15 - 1
	MinInt16  = -1 << 15
	MaxInt32  = 1<<31 - 1
	MinInt32  = -1 << 31
	MaxInt64  = 1<<63 - 1
	MinInt64  = -1 << 63
	MaxUint8  = 1<<8 - 1
	MaxUint16 = 1<<16 - 1
	MaxUint32 = 1<<32 - 1
	MaxUint64 = 1<<64 - 1
)

// Abs returns the absolute value of x.
//
// Special cases are:
//
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN

// Abs 返回 x 的绝对值。 特殊情况为：
//
//	Abs(±Inf) = +Inf
//	Abs(NaN)  = NaN
func Abs(x float64) float64

// Acos returns the arccosine, in radians, of x.
//
// Special case is:
//
//	Acos(x) = NaN if x < -1 or x > 1

// Acos 返回 x 的反余弦值。
//
// 特殊情况为：
//
//	Acos(x) = NaN（若 x < -1 或 x > 1）
func Acos(x float64) float64

// Acosh returns the inverse hyperbolic cosine of x.
//
// Special cases are:
//
//	Acosh(+Inf) = +Inf
//	Acosh(x) = NaN if x < 1
//	Acosh(NaN) = NaN

// Acosh 返回 x 的反双曲余弦值。
//
// 特殊情况为：
//
//	Acosh(+Inf) = +Inf
//	Acosh(x)    = NaN（若 x < 1）
//	Acosh(NaN)  = NaN
func Acosh(x float64) float64

// Asin returns the arcsine, in radians, of x.
//
// Special cases are:
//
//	Asin(±0) = ±0
//	Asin(x) = NaN if x < -1 or x > 1

// Asin 返回 x 的反正弦值。
//
// 特殊情况为：
//
//	Asin(±0) = ±0
//	Asin(x)  = NaN（若 x < -1 或 x > 1）
func Asin(x float64) float64

// Asinh returns the inverse hyperbolic sine of x.
//
// Special cases are:
//
//	Asinh(±0) = ±0
//	Asinh(±Inf) = ±Inf
//	Asinh(NaN) = NaN

// Asinh 返回 x 的反双曲正弦值。
//
// 特殊情况为：
//
//	Asinh(±0)   = ±0
//	Asinh(±Inf) = ±Inf
//	Asinh(NaN)  = NaN
func Asinh(x float64) float64

// Atan returns the arctangent, in radians, of x.
//
// Special cases are:
//
//	Atan(±0) = ±0
//	Atan(±Inf) = ±Pi/2

// Atan 返回 x 的反正切值。
//
// 特殊情况为：
//
//	Atan(±0)   = ±0
//	Atan(±Inf) = ±Pi/2
func Atan(x float64) float64

// Atan2 returns the arc tangent of y/x, using the signs of the two to determine
// the quadrant of the return value.
//
// Special cases are (in order):
//
//	Atan2(y, NaN) = NaN
//	Atan2(NaN, x) = NaN
//	Atan2(+0, x>=0) = +0
//	Atan2(-0, x>=0) = -0
//	Atan2(+0, x<=-0) = +Pi
//	Atan2(-0, x<=-0) = -Pi
//	Atan2(y>0, 0) = +Pi/2
//	Atan2(y<0, 0) = -Pi/2
//	Atan2(+Inf, +Inf) = +Pi/4
//	Atan2(-Inf, +Inf) = -Pi/4
//	Atan2(+Inf, -Inf) = 3Pi/4
//	Atan2(-Inf, -Inf) = -3Pi/4
//	Atan2(y, +Inf) = 0
//	Atan2(y>0, -Inf) = +Pi
//	Atan2(y<0, -Inf) = -Pi
//	Atan2(+Inf, x) = +Pi/2
//	Atan2(-Inf, x) = -Pi/2

// Atan2 返回 y/x 的反正切值，通过二者的符号决定其返回值的象限。
//
// 特殊情况为（按顺序）：
//
//	Atan2(y, NaN)     = NaN
//	Atan2(NaN, x)     = NaN
//	Atan2(+0, x>=0)   = +0
//	Atan2(-0, x>=0)   = -0
//	Atan2(+0, x<=-0)  = +Pi
//	Atan2(-0, x<=-0)  = -Pi
//	Atan2(y>0, 0)     = +Pi/2
//	Atan2(y<0, 0)     = -Pi/2
//	Atan2(+Inf, +Inf) = +Pi/4
//	Atan2(-Inf, +Inf) = -Pi/4
//	Atan2(+Inf, -Inf) = 3Pi/4
//	Atan2(-Inf, -Inf) = -3Pi/4
//	Atan2(y, +Inf)    = 0
//	Atan2(y>0, -Inf)  = +Pi
//	Atan2(y<0, -Inf)  = -Pi
//	Atan2(+Inf, x)    = +Pi/2
//	Atan2(-Inf, x)    = -Pi/2
func Atan2(y, x float64) float64

// Atanh returns the inverse hyperbolic tangent of x.
//
// Special cases are:
//
//	Atanh(1) = +Inf
//	Atanh(±0) = ±0
//	Atanh(-1) = -Inf
//	Atanh(x) = NaN if x < -1 or x > 1
//	Atanh(NaN) = NaN

// Atanh 返回 x 的反双曲正切值。
//
// 特殊情况为：
//
//	Atanh(1)   = +Inf
//	Atanh(±0)  = ±0
//	Atanh(-1)  = -Inf
//	Atanh(x)   = NaN（若 x < -1 或 x > 1）
//	Atanh(NaN) = NaN
func Atanh(x float64) float64

// Cbrt returns the cube root of x.
//
// Special cases are:
//
//	Cbrt(±0) = ±0
//	Cbrt(±Inf) = ±Inf
//	Cbrt(NaN) = NaN

// Cbrt 返回 x 的立方根。
//
// 特殊情况为：
//
//	Cbrt(±0)   = ±0
//	Cbrt(±Inf) = ±Inf
//	Cbrt(NaN)  = NaN
func Cbrt(x float64) float64

// Ceil returns the least integer value greater than or equal to x.
//
// Special cases are:
//
//	Ceil(±0) = ±0
//	Ceil(±Inf) = ±Inf
//	Ceil(NaN) = NaN

// Ceil 返回大于或等于 x 的最小整数。
//
// 特殊情况为：
//
//	Ceil(±0)   = ±0
//	Ceil(±Inf) = ±Inf
//	Ceil(NaN)  = NaN
func Ceil(x float64) float64

// Copysign returns a value with the magnitude of x and the sign of y.

// Copysign 的返回值由 x 的量和 y 的符号构成。
func Copysign(x, y float64) float64

// Cos returns the cosine of the radian argument x.
//
// Special cases are:
//
//	Cos(±Inf) = NaN
//	Cos(NaN) = NaN

// Cos 返回 x 的余弦值。
//
// 特殊情况为：
//
//	Cos(±Inf) = NaN
//	Cos(NaN)  = NaN
func Cos(x float64) float64

// Cosh returns the hyperbolic cosine of x.
//
// Special cases are:
//
//	Cosh(±0) = 1
//	Cosh(±Inf) = +Inf
//	Cosh(NaN) = NaN

// Cosh 返回 x 的双曲余弦值。
//
// 特殊情况为：
//
//	Cosh(±0)   = 1
//	Cosh(±Inf) = +Inf
//	Cosh(NaN)  = NaN
func Cosh(x float64) float64

// Dim returns the maximum of x-y or 0.
//
// Special cases are:
//
//	Dim(+Inf, +Inf) = NaN
//	Dim(-Inf, -Inf) = NaN
//	Dim(x, NaN) = Dim(NaN, x) = NaN

// Dim 返回 x-y 和 0 中较大的数。
//
// 特殊情况为：
//
//	Dim(+Inf, +Inf)           = NaN
//	Dim(-Inf, -Inf)           = NaN
//	Dim(x, NaN) = Dim(NaN, x) = NaN
func Dim(x, y float64) float64

// Erf returns the error function of x.
//
// Special cases are:
//
//	Erf(+Inf) = 1
//	Erf(-Inf) = -1
//	Erf(NaN) = NaN

// Erf 返回 x 的误差函数。
//
// 特殊情况为：
//
//	Erf(+Inf) = 1
//	Erf(-Inf) = -1
//	Erf(NaN)  = NaN
func Erf(x float64) float64

// Erfc returns the complementary error function of x.
//
// Special cases are:
//
//	Erfc(+Inf) = 0
//	Erfc(-Inf) = 2
//	Erfc(NaN) = NaN

// Erfc 返回 x 的余误差函数。
//
// 特殊情况为：
//
//	Erfc(+Inf) = 0
//	Erfc(-Inf) = 2
//	Erfc(NaN)  = NaN
func Erfc(x float64) float64

// Exp returns e**x, the base-e exponential of x.
//
// Special cases are:
//
//	Exp(+Inf) = +Inf
//	Exp(NaN) = NaN
//
// Very large values overflow to 0 or +Inf. Very small values underflow to 1.

// Exp 返回 e**x，即以 e 为底的 x 次幂。
//
// 特殊情况为：
//
//	Exp(+Inf) = +Inf
//	Exp(NaN)  = NaN
//
// 非常大的数会向上溢出为 0 或 +Inf。 非常小的数会向下溢出为 1。
func Exp(x float64) float64

// Exp2 returns 2**x, the base-2 exponential of x.
//
// Special cases are the same as Exp.

// Exp2 返回 2**x，即以 2 为底的 x 次指数。
//
// 特殊情况与 Exp 相同。
func Exp2(x float64) float64

// Expm1 returns e**x - 1, the base-e exponential of x minus 1. It is more accurate
// than Exp(x) - 1 when x is near zero.
//
// Special cases are:
//
//	Expm1(+Inf) = +Inf
//	Expm1(-Inf) = -1
//	Expm1(NaN) = NaN
//
// Very large values overflow to -1 or +Inf.

// Expm1 返回 e**x - 1，即以 e 为底的 x 次幂减一。 当 x
// 接近 0 时，该函数比 Exp(x) - 1 更精确。
//
// 特殊情况为：
//
//	Expm1(+Inf) = +Inf
//	Expm1(-Inf) = -1
//	Expm1(NaN)  = NaN
//
// 非常大的值会溢出为 -1 或 +Inf。
func Expm1(x float64) float64

// Float32bits returns the IEEE 754 binary representation of f.

// Float32bits 返回 f 的IEEE 754二进制表示。
func Float32bits(f float32) uint32

// Float32frombits returns the floating point number corresponding to the IEEE 754
// binary representation b.

// Float32frombits 返回与IEEE 754二进制表示 b
// 相应的浮点数。
func Float32frombits(b uint32) float32

// Float64bits returns the IEEE 754 binary representation of f.

// Float64bits 返回 f 的IEEE 754二进制表示。
func Float64bits(f float64) uint64

// Float64frombits returns the floating point number corresponding the IEEE 754
// binary representation b.

// Float64frombits 返回与IEEE 754二进制表示 b
// 相应的浮点数。
func Float64frombits(b uint64) float64

// Floor returns the greatest integer value less than or equal to x.
//
// Special cases are:
//
//	Floor(±0) = ±0
//	Floor(±Inf) = ±Inf
//	Floor(NaN) = NaN

// Floor 返回小于或等于 x 的最大整数。
//
// 特殊情况为：
//
//	Floor(±0)   = ±0
//	Floor(±Inf) = ±Inf
//	Floor(NaN)  = NaN
func Floor(x float64) float64

// Frexp breaks f into a normalized fraction and an integral power of two. It
// returns frac and exp satisfying f == frac × 2**exp, with the absolute value of
// frac in the interval [½, 1).
//
// Special cases are:
//
//	Frexp(±0) = ±0, 0
//	Frexp(±Inf) = ±Inf, 0
//	Frexp(NaN) = NaN, 0

// Frexp 将 f 分解为一个规范化的小数和一个 2 的整数次幂。 它返回的
// frac 和 exp 满足 f == frac × 2**exp，且 frac
// 的绝对值在区间 [½, 1) 内。
//
// 特殊情况为：
//
//	Frexp(±0)   = ±0, 0
//	Frexp(±Inf) = ±Inf, 0
//	Frexp(NaN)  = NaN, 0
func Frexp(f float64) (frac float64, exp int)

// Gamma returns the Gamma function of x.
//
// Special cases are:
//
//	Gamma(+Inf) = +Inf
//	Gamma(+0) = +Inf
//	Gamma(-0) = -Inf
//	Gamma(x) = NaN for integer x < 0
//	Gamma(-Inf) = NaN
//	Gamma(NaN) = NaN

// Gamma 返回 x 的伽马函数。
//
// 特殊情况为：
//
//	Gamma(+Inf) = +Inf
//	Gamma(+0)   = +Inf
//	Gamma(-0)   = -Inf
//	Gamma(x)    = NaN（对于整数 x < 0）
//	Gamma(-Inf) = NaN
//	Gamma(NaN)  = NaN
func Gamma(x float64) float64

// Hypot returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and
// underflow.
//
// Special cases are:
//
//	Hypot(±Inf, q) = +Inf
//	Hypot(p, ±Inf) = +Inf
//	Hypot(NaN, q) = NaN
//	Hypot(p, NaN) = NaN

// Hypot 返回 Sqrt(p*p +
// q*q)，小心避免不必要的向上溢出和向下溢出。
//
// 特殊情况为：
//
//	Hypot(±Inf, q) = +Inf
//	Hypot(p, ±Inf) = +Inf
//	Hypot(NaN, q)  = NaN
//	Hypot(p, NaN)  = NaN
func Hypot(p, q float64) float64

// Ilogb returns the binary exponent of x as an integer.
//
// Special cases are:
//
//	Ilogb(±Inf) = MaxInt32
//	Ilogb(0) = MinInt32
//	Ilogb(NaN) = MaxInt32

// Ilogb 将以 2 为底 x 的指数作为整数返回。
//
// 特殊情况为：
//
//	Ilogb(±Inf) = MaxInt32
//	Ilogb(0)    = MinInt32
//	Ilogb(NaN)  = MaxInt32
func Ilogb(x float64) int

// Inf returns positive infinity if sign >= 0, negative infinity if sign < 0.

// Inf 返回无穷大值（infinity）。若 sign >=
// 0，则返回正无穷大（positive infinity）； 若 sign <
// 0，则返回负无穷大（negative infinity）。
func Inf(sign int) float64

// IsInf reports whether f is an infinity, according to sign. If sign > 0, IsInf
// reports whether f is positive infinity. If sign < 0, IsInf reports whether f is
// negative infinity. If sign == 0, IsInf reports whether f is either infinity.

// IsInf 判断 f 是否为无穷大值，视 sign 而定。 若 sign >
// 0，IsInf 就判断 f 是否为正无穷大。 若 sign < 0，IsInf
// 就判断 f 是否为负无穷大。 若 sign == 0，IsInf 就判断 f
// 是否为无穷大。
func IsInf(f float64, sign int) bool

// IsNaN reports whether f is an IEEE 754 ``not-a-number'' value.

// IsNaN 判断 f 是否为IEEE
// 754定义的“非数值”（Not-a-Number）。
func IsNaN(f float64) (is bool)

// J0 returns the order-zero Bessel function of the first kind.
//
// Special cases are:
//
//	J0(±Inf) = 0
//	J0(0) = 1
//	J0(NaN) = NaN

// J0 返回第一类零阶贝塞尔函数。
//
// 特殊情况为：
//
//	J0(±Inf) = 0
//	J0(0)    = 1
//	J0(NaN)  = NaN
func J0(x float64) float64

// J1 returns the order-one Bessel function of the first kind.
//
// Special cases are:
//
//	J1(±Inf) = 0
//	J1(NaN) = NaN

// J1 返回一阶第一类贝塞尔函数。
//
// 特殊情况为
//
//	J1(±Inf) = 0
//	J1(NaN)  = NaN
func J1(x float64) float64

// Jn returns the order-n Bessel function of the first kind.
//
// Special cases are:
//
//	Jn(n, ±Inf) = 0
//	Jn(n, NaN) = NaN

// Jn 返回 n 阶第一类贝塞尔函数。
//
// 特殊情况为：
//
//	Jn(n, ±Inf) = 0
//	Jn(n, NaN)  = NaN
func Jn(n int, x float64) float64

// Ldexp is the inverse of Frexp. It returns frac × 2**exp.
//
// Special cases are:
//
//	Ldexp(±0, exp) = ±0
//	Ldexp(±Inf, exp) = ±Inf
//	Ldexp(NaN, exp) = NaN

// Ldexp 为 Frexp 的反函数。 它返回 frac × 2**exp。
//
// 特殊情况为：
//
//	Ldexp(±0, exp)   = ±0
//	Ldexp(±Inf, exp) = ±Inf
//	Ldexp(NaN, exp)  = NaN
func Ldexp(frac float64, exp int) float64

// Lgamma returns the natural logarithm and sign (-1 or +1) of Gamma(x).
//
// Special cases are:
//
//	Lgamma(+Inf) = +Inf
//	Lgamma(0) = +Inf
//	Lgamma(-integer) = +Inf
//	Lgamma(-Inf) = -Inf
//	Lgamma(NaN) = NaN

// Lgamma 返回 Gamma(x) 的自然对数和符号（-1 或 +1）。
//
// 特殊情况为：
//
//	Lgamma(+Inf)     = +Inf
//	Lgamma(0)        = +Inf
//	Lgamma(-integer) = +Inf
//	Lgamma(-Inf)     = -Inf
//	Lgamma(NaN)      = NaN
func Lgamma(x float64) (lgamma float64, sign int)

// Log returns the natural logarithm of x.
//
// Special cases are:
//
//	Log(+Inf) = +Inf
//	Log(0) = -Inf
//	Log(x < 0) = NaN
//	Log(NaN) = NaN

// Log 返回 x 的自然对数。
//
// 特殊情况为
//
//	Log(+Inf)  = +Inf
//	Log(0)     = -Inf
//	Log(x < 0) = NaN
//	Log(NaN)   = NaN
func Log(x float64) float64

// Log10 returns the decimal logarithm of x. The special cases are the same as for
// Log.

// Log10 返回以 10 为底 x 的对数。 特殊情况与 Log 相同。
func Log10(x float64) float64

// Log1p returns the natural logarithm of 1 plus its argument x. It is more
// accurate than Log(1 + x) when x is near zero.
//
// Special cases are:
//
//	Log1p(+Inf) = +Inf
//	Log1p(±0) = ±0
//	Log1p(-1) = -Inf
//	Log1p(x < -1) = NaN
//	Log1p(NaN) = NaN

// Log1p 返回 1 加其实参 x 的自然对数。 当 x 接近 0 时，该函数比
// Log(1 + x) 精确。
//
// 特殊情况为：
//
//	Log1p(+Inf)   = +Inf
//	Log1p(±0)     = ±0
//	Log1p(-1)     = -Inf
//	Log1p(x < -1) = NaN
//	Log1p(NaN)    = NaN
func Log1p(x float64) float64

// Log2 returns the binary logarithm of x. The special cases are the same as for
// Log.

// Log2 返回以 2 为底 x 的对数。 特殊情况与 Log 相同。
func Log2(x float64) float64

// Logb returns the binary exponent of x.
//
// Special cases are:
//
//	Logb(±Inf) = +Inf
//	Logb(0) = -Inf
//	Logb(NaN) = NaN

// Logb 返回以 2 为底 x 的指数。
//
// 特殊情况为：
//
//	Logb(±Inf) = +Inf
//	Logb(0)    = -Inf
//	Logb(NaN)  = NaN
func Logb(x float64) float64

// Max returns the larger of x or y.
//
// Special cases are:
//
//	Max(x, +Inf) = Max(+Inf, x) = +Inf
//	Max(x, NaN) = Max(NaN, x) = NaN
//	Max(+0, ±0) = Max(±0, +0) = +0
//	Max(-0, -0) = -0

// Max 返回 x 和 y 中较大的数。
//
// 特殊情况为：
//
//	Max(x, +Inf) = Max(+Inf, x) = +Inf
//	Max(x, NaN)  = Max(NaN, x)  = NaN
//	Max(+0, ±0)  = Max(±0, +0)  = +0
//	Max(-0, -0)  = -0
func Max(x, y float64) float64

// Min returns the smaller of x or y.
//
// Special cases are:
//
//	Min(x, -Inf) = Min(-Inf, x) = -Inf
//	Min(x, NaN) = Min(NaN, x) = NaN
//	Min(-0, ±0) = Min(±0, -0) = -0

// Min 返回 x 和 y 中较小的数。
//
// 特殊情况为：
//
//	Min(x, -Inf) = Min(-Inf, x) = -Inf
//	Min(x, NaN)  = Min(NaN, x)  = NaN
//	Min(-0, ±0)  = Min(±0, -0)  = -0
func Min(x, y float64) float64

// Mod returns the floating-point remainder of x/y. The magnitude of the result is
// less than y and its sign agrees with that of x.
//
// Special cases are:
//
//	Mod(±Inf, y) = NaN
//	Mod(NaN, y) = NaN
//	Mod(x, 0) = NaN
//	Mod(x, ±Inf) = x
//	Mod(x, NaN) = NaN

// Mod 返回 x/y 的浮点余数。 其结果的大小小于 y 且其符号与 x 一致。
//
// 特殊情况为：
//
//	Mod(±Inf, y) = NaN
//	Mod(NaN, y)  = NaN
//	Mod(x, 0)    = NaN
//	Mod(x, ±Inf) = x
//	Mod(x, NaN)  = NaN
func Mod(x, y float64) float64

// Modf returns integer and fractional floating-point numbers that sum to f. Both
// values have the same sign as f.
//
// Special cases are:
//
//	Modf(±Inf) = ±Inf, NaN
//	Modf(NaN) = NaN, NaN

// Modf 将 f 的整数部分和小数部分分别作为浮点数返回。两值的符号与 f
// 一致。
//
// 特殊情况为：
//
//	Modf(±Inf) = ±Inf, NaN
//	Modf(NaN)  = NaN, NaN
func Modf(f float64) (int float64, frac float64)

// NaN returns an IEEE 754 ``not-a-number'' value.

// NaN 返回IEEE 754定义的“非数值”（Not-a-Number）。
func NaN() float64

// Nextafter returns the next representable float64 value after x towards y.
// Special cases:
//
//		Nextafter64(x, x)   = x
//	     Nextafter64(NaN, y) = NaN
//	     Nextafter64(x, NaN) = NaN

// Nextafter 返回从 x 到 y 的下一个可表示的 float64 值。
//
// 特殊情况为：
//
//	Nextafter64(x, x)   = x
//	Nextafter(NaN, y) = NaN
//	Nextafter(x, NaN) = NaN
func Nextafter(x, y float64) (r float64)

// Nextafter32 returns the next representable float32 value after x towards y.
// Special cases:
//
//		Nextafter32(x, x)   = x
//	     Nextafter32(NaN, y) = NaN
//	     Nextafter32(x, NaN) = NaN

// Nextafter32 返回从 x 到 y 的下一个可表示的 float32
// 值。
//
// 特殊情况为：
//
//	Nextafter32(x, x)   = x
//	Nextafter32(NaN, y) = NaN
//	Nextafter32(x, NaN) = NaN
func Nextafter32(x, y float32) (r float32)

// Pow returns x**y, the base-x exponential of y.
//
// Special cases are (in order):
//
//	Pow(x, ±0) = 1 for any x
//	Pow(1, y) = 1 for any y
//	Pow(x, 1) = x for any x
//	Pow(NaN, y) = NaN
//	Pow(x, NaN) = NaN
//	Pow(±0, y) = ±Inf for y an odd integer < 0
//	Pow(±0, -Inf) = +Inf
//	Pow(±0, +Inf) = +0
//	Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
//	Pow(±0, y) = ±0 for y an odd integer > 0
//	Pow(±0, y) = +0 for finite y > 0 and not an odd integer
//	Pow(-1, ±Inf) = 1
//	Pow(x, +Inf) = +Inf for |x| > 1
//	Pow(x, -Inf) = +0 for |x| > 1
//	Pow(x, +Inf) = +0 for |x| < 1
//	Pow(x, -Inf) = +Inf for |x| < 1
//	Pow(+Inf, y) = +Inf for y > 0
//	Pow(+Inf, y) = +0 for y < 0
//	Pow(-Inf, y) = Pow(-0, -y)
//	Pow(x, y) = NaN for finite x < 0 and finite non-integer y

// Pow 返回 x**y，即以 x 为底的 y 次幂。
//
// 特殊情况为（按顺序）：
//
//	Pow(x, ±0)    = 1   （对于任何 x）
//	Pow(1, y)     = 1   （对于任何 y）
//	Pow(x, 1)     = x   （对于任何 x）
//	Pow(NaN, y)   = NaN
//	Pow(x, NaN)   = NaN
//	Pow(±0, y)    = ±Inf（对于奇整数 y < 0）
//	Pow(±0, -Inf) = +Inf
//	Pow(±0, +Inf) = +0
//	Pow(±0, y)    = +Inf（对于有限非奇整数 y < 0）
//	Pow(±0, y)    = ±0  （对于奇整数 y > 0）
//	Pow(±0, y)    = +0  （对于有限非奇整数 y >）
//	Pow(-1, ±Inf) = 1
//	Pow(x, +Inf)  = +Inf（对于 |x| > 1）
//	Pow(x, -Inf)  = +0  （对于 |x| > 1）
//	Pow(x, +Inf)  = +0  （对于 |x| < 1）
//	Pow(x, -Inf)  = +Inf（对于 |x| < 1）
//	Pow(+Inf, y)  = +Inf（对于 y > 0）
//	Pow(+Inf, y)  = +0  （对于 y < 0）
//	Pow(-Inf, y)  = Pow(-0, -y)
//	Pow(x, y)     = NaN （对于有限数 x < 0 和有限非整数 y）
func Pow(x, y float64) float64

// Pow10 returns 10**e, the base-10 exponential of e.
//
// Special cases are:
//
//	Pow10(e) = +Inf for e > 309
//	Pow10(e) = 0 for e < -324

// Pow10 返回 10**e，即以 10 为底的 e 次幂。
//
// 特殊情况为：
//
//	对于 e >  309，有 Pow10(e) = +Inf
//	对于 e < -324，有 Pow10(e) = 0
func Pow10(e int) float64

// Remainder returns the IEEE 754 floating-point remainder of x/y.
//
// Special cases are:
//
//	Remainder(±Inf, y) = NaN
//	Remainder(NaN, y) = NaN
//	Remainder(x, 0) = NaN
//	Remainder(x, ±Inf) = x
//	Remainder(x, NaN) = NaN

// Remainder 返回IEEE 754标准 x/y 的余数。
//
// 特殊情况为：
//
//	Remainder(±Inf, y) = NaN
//	Remainder(NaN, y)  = NaN
//	Remainder(x, 0)    = NaN
//	Remainder(x, ±Inf) = x
//	Remainder(x, NaN)  = NaN
func Remainder(x, y float64) float64

// Signbit returns true if x is negative or negative zero.

// Signbit 判断 x 是否为负值或负零。
func Signbit(x float64) bool

// Sin returns the sine of the radian argument x.
//
// Special cases are:
//
//	Sin(±0) = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN) = NaN

// Sin 返回 x 的正弦值。
//
// 特殊情况为：
//
//	Sin(±0)   = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN)  = NaN
func Sin(x float64) float64

// Sincos returns Sin(x), Cos(x).
//
// Special cases are:
//
//	Sincos(±0) = ±0, 1
//	Sincos(±Inf) = NaN, NaN
//	Sincos(NaN) = NaN, NaN

// Sincos 返回 Sin(x)，Cos(x)。
//
// 特殊情况为：
//
//	Sincos(±0)   = ±0, 1
//	Sincos(±Inf) = NaN, NaN
//	Sincos(NaN)  = NaN, NaN
func Sincos(x float64) (sin, cos float64)

// Sinh returns the hyperbolic sine of x.
//
// Special cases are:
//
//	Sinh(±0) = ±0
//	Sinh(±Inf) = ±Inf
//	Sinh(NaN) = NaN

// Sinh 返回 x 的双曲正弦值。
//
// 特殊情况为：
//
//	Sinh(±0)   = ±0
//	Sinh(±Inf) = ±Inf
//	Sinh(NaN)  = NaN
func Sinh(x float64) float64

// Sqrt returns the square root of x.
//
// Special cases are:
//
//	Sqrt(+Inf) = +Inf
//	Sqrt(±0) = ±0
//	Sqrt(x < 0) = NaN
//	Sqrt(NaN) = NaN

// Sqrt 返回 x 的平方根。
//
// 特殊情况为：
//
//	Sqrt(+Inf)  = +Inf
//	Sqrt(±0)    = ±0
//	Sqrt(x < 0) = NaN
//	Sqrt(NaN)   = NaN
func Sqrt(x float64) float64

// Tan returns the tangent of the radian argument x.
//
// Special cases are:
//
//	Tan(±0) = ±0
//	Tan(±Inf) = NaN
//	Tan(NaN) = NaN

// Tan 返回 x 的正切值。
//
// 特殊情况为：
//
//	Tan(±0)   = ±0
//	Tan(±Inf) = NaN
//	Tan(NaN)  = NaN
func Tan(x float64) float64

// Tanh returns the hyperbolic tangent of x.
//
// Special cases are:
//
//	Tanh(±0) = ±0
//	Tanh(±Inf) = ±1
//	Tanh(NaN) = NaN

// Tanh 返回 x 的双曲正切。
//
// 特殊情况为：
//
//	Tanh(±0)   = ±0
//	Tanh(±Inf) = ±1
//	Tanh(NaN)  = NaN
func Tanh(x float64) float64

// Trunc returns the integer value of x.
//
// Special cases are:
//
//	Trunc(±0) = ±0
//	Trunc(±Inf) = ±Inf
//	Trunc(NaN) = NaN

// Trunc 返回 x 的整数部分
//
// 特殊情况为：
//
//	Trunc(±0)   = ±0
//	Trunc(±Inf) = ±Inf
//	Trunc(NaN)  = NaN
func Trunc(x float64) float64

// Y0 returns the order-zero Bessel function of the second kind.
//
// Special cases are:
//
//	Y0(+Inf) = 0
//	Y0(0) = -Inf
//	Y0(x < 0) = NaN
//	Y0(NaN) = NaN

// Y0 返回第二类零阶贝塞尔函数。
//
// 特殊情况为：
//
//	Y0(+Inf) = 0
//	Y0(0)    = -Inf
//	Y0(x<0)  = NaN
//	Y0(NaN)  = NaN
func Y0(x float64) float64

// Y1 returns the order-one Bessel function of the second kind.
//
// Special cases are:
//
//	Y1(+Inf) = 0
//	Y1(0) = -Inf
//	Y1(x < 0) = NaN
//	Y1(NaN) = NaN

// Y1 返回一阶第二类贝塞尔函数。
//
// 特殊情况为：
//
//	Y1(+Inf) = 0
//	Y1(0)    = -Inf
//	Y1(x<0)  = NaN
//	Y1(NaN)  = NaN
func Y1(x float64) float64

// Yn returns the order-n Bessel function of the second kind.
//
// Special cases are:
//
//	Yn(n, +Inf) = 0
//	Yn(n > 0, 0) = -Inf
//	Yn(n < 0, 0) = +Inf if n is odd, -Inf if n is even
//	Y1(n, x < 0) = NaN
//	Y1(n, NaN) = NaN

// Yn 返回 n 阶第二类贝塞尔函数。
//
// 特殊情况为：
//
//	Yn(n, +Inf)  = 0
//	Yn(n > 0, 0) = -Inf
//	Yn(n < 0, 0) = 若 n 为奇数则为 +Inf，若 n 为偶数则为 -Inf
//	Y1(n, x < 0) = NaN
//	Y1(n, NaN)   = NaN
func Yn(n int, x float64) float64
