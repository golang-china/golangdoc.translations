// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package cmplx provides basic constants and mathematical functions for complex
// numbers.

// cmplx 包为复数提供了基本的常量和数学函数.
package cmplx

// Abs returns the absolute value (also called the modulus) of x.

// Abs 返回 x 的绝对值（亦称为模）。
func Abs(x complex128) float64

// Acos returns the inverse cosine of x.

// Acos 返回 x 的反余弦值
func Acos(x complex128) complex128

// Acosh returns the inverse hyperbolic cosine of x.

// Acosh 返回 x 的反双曲余弦值。
func Acosh(x complex128) complex128

// Asin returns the inverse sine of x.

// Asin 返回 x 的反正弦值。
func Asin(x complex128) complex128

// Asinh returns the inverse hyperbolic sine of x.

// Asinh 返回 x 的反双曲正弦值。
func Asinh(x complex128) complex128

// Atan returns the inverse tangent of x.

// Atan 返回 x 的反正切。
func Atan(x complex128) complex128

// Atanh returns the inverse hyperbolic tangent of x.

// Atanh 返回 x 的双曲反正切。
func Atanh(x complex128) complex128

// Conj returns the complex conjugate of x.

// Conj 返回 x 的复数共轭。
func Conj(x complex128) complex128

// Cos returns the cosine of x.

// Cos 返回 x 的余弦值。
func Cos(x complex128) complex128

// Cosh returns the hyperbolic cosine of x.

// Cosh 返回 x 的双曲余弦值。
func Cosh(x complex128) complex128

// Cot returns the cotangent of x.

// Cot 返回 x 的反正切值。
func Cot(x complex128) complex128

// Exp returns e**x, the base-e exponential of x.

// Exp 返回 e**x，即以 e 为低的 x 次幂。
func Exp(x complex128) complex128

// Inf returns a complex infinity, complex(+Inf, +Inf).

// Inf 返回一个复数的无限大值，即 complex(+Inf, +Inf)。
func Inf() complex128

// IsInf returns true if either real(x) or imag(x) is an infinity.

// IsInf 在 real(x) 或 imag(x) 为无限大值时返回 true。
func IsInf(x complex128) bool

// IsNaN returns true if either real(x) or imag(x) is NaN and neither is an
// infinity.

// IsNaN 在 real(x) 或 imag(x) 其中之一为 NaN
// 且另一个为无限大值时返回 true。
func IsNaN(x complex128) bool

// Log returns the natural logarithm of x.

// Log 返回 x 的自然对数。
func Log(x complex128) complex128

// Log10 returns the decimal logarithm of x.

// Log10 返回 x 的十进制对数。
func Log10(x complex128) complex128

// NaN returns a complex ``not-a-number'' value.

// NaN 返回一个复数的“非数值”。
func NaN() complex128

// Phase returns the phase (also called the argument) of x. The returned value is
// in the range [-Pi, Pi].

// Phase 返回 x 的 相位（亦称为辐角）。 其返回值在区间 [-Pi,
// Pi] 内。
func Phase(x complex128) float64

// Polar returns the absolute value r and phase θ of x, such that x = r * e**θi.
// The phase is in the range [-Pi, Pi].

// Polar 返回 x 的绝对值 r 和相位 θ，使得 x = r *
// e**θi。 其相位在区间 [-Pi, Pi] 内。
func Polar(x complex128) (r, θ float64)

// Pow returns x**y, the base-x exponential of y. For generalized compatibility
// with math.Pow:
//
//	Pow(0, ±0) returns 1+0i
//	Pow(0, c) for real(c)<0 returns Inf+0i if imag(c) is zero, otherwise Inf+Inf i.

// Pow 返回 x**y，即以 x 为底的 y 次幂。 对于 math.Pow
// 的通用化兼容：
//
//	Pow(0, ±0) 返回 1+0i
//	若 imag(c) 为零，则 Pow(0, c) 在 real(c)<0 时返回 Inf+0i, 否则返回 Inf+Inf i。
func Pow(x, y complex128) complex128

// Rect returns the complex number x with polar coordinates r, θ.

// Rect 返回极坐标形式 (r, θ) 的复数 x。
func Rect(r, θ float64) complex128

// Sin returns the sine of x.

// Sin 返回 x 的正弦值。
func Sin(x complex128) complex128

// Sinh returns the hyperbolic sine of x.

// Sinh 返回 x 的双曲正弦值。
func Sinh(x complex128) complex128

// Sqrt returns the square root of x. The result r is chosen so that real(r) ≥ 0
// and imag(r) has the same sign as imag(x).

// Sqrt returns the square root of x. The
// result r is chosen so that real(r) ≥ 0
// and imag(r) has the same sign as
// imag(x).
func Sqrt(x complex128) complex128

// Tan returns the tangent of x.

// Tan 返回 x 的正切值。
func Tan(x complex128) complex128

// Tanh returns the hyperbolic tangent of x.

// Tanh 返回 x 的双曲正切。
func Tanh(x complex128) complex128
