// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package elliptic implements several standard elliptic curves over prime
// fields.

// elliptic包实现了几条覆盖素数有限域的标准椭圆曲线。
package elliptic

import (
	"io"
	"math/big"
	"sync"
)

// A Curve represents a short-form Weierstrass curve with a=-3.
// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw.html

// Curve代表一个短格式的Weierstrass椭圆曲线，其中a=-3。
//
// Weierstrass椭圆曲线的格式：y**2 = x**3 + a*x + b
//
// 参见http://www.hyperelliptic.org/EFD/g1p/auto-shortw.html
type Curve interface {
	// Params returns the parameters for the curve.
	Params()*CurveParams

	// IsOnCurve reports whether the given (x,y) lies on the curve.
	IsOnCurve(x, y *big.Int)bool

	// Add returns the sum of (x1,y1) and (x2,y2)
	Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int)

	// Double returns 2*(x,y)
	Double(x1, y1 *big.Int) (x, y *big.Int)

	// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
	ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int)

	// ScalarBaseMult returns k*G, where G is the base point of the group
	// and k is an integer in big-endian form.
	ScalarBaseMult(k []byte) (x, y *big.Int)
}

// CurveParams contains the parameters of an elliptic curve and also provides
// a generic, non-constant time implementation of Curve.

// CurveParams包含一个椭圆曲线的所有参数，也可提供一般的、非常数时间实现的椭圆曲
// 线。
type CurveParams struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	B       *big.Int // the constant of the curve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
}

// GenerateKey returns a public/private key pair. The private key is
// generated using the given reader, which must return random data.

// GenerateKey返回一个公钥/私钥对。priv是私钥，而(x,y)是公钥。密钥对是通过提供的
// 随机数读取器来生成的，该io.Reader接口必须返回随机数据。
func GenerateKey(curve Curve, rand io.Reader) (priv []byte, x, y *big.Int, err error)

// Marshal converts a point into the form specified in section 4.3.6 of ANSI
// X9.62.

// Marshal将一个点编码为ANSI X9.62指定的格式。
func Marshal(curve Curve, x, y *big.Int) []byte

// P224 returns a Curve which implements P-224 (see FIPS 186-3, section D.2.2)

// 返回一个实现了P-224的曲线。（参见FIPS 186-3, section D.2.2）
func P224() Curve

// P256 returns a Curve which implements P-256 (see FIPS 186-3, section D.2.3)

// 返回一个实现了P-256的曲线。（参见FIPS 186-3, section D.2.3）
func P256() Curve

// P384 returns a Curve which implements P-384 (see FIPS 186-3, section D.2.4)

// 返回一个实现了P-384的曲线。（参见FIPS 186-3, section D.2.4）
func P384() Curve

// P521 returns a Curve which implements P-521 (see FIPS 186-3, section D.2.5)

// 返回一个实现了P-512的曲线。（参见FIPS 186-3, section D.2.5）
func P521() Curve

// Unmarshal converts a point, serialized by Marshal, into an x, y pair.
// It is an error if the point is not on the curve. On error, x = nil.

// 将一个Marshal编码后的点还原；如果出错，x会被设为nil。
func Unmarshal(curve Curve, data []byte) (x, y *big.Int)

func (curve *CurveParams) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int)

func (curve *CurveParams) Double(x1, y1 *big.Int) (*big.Int, *big.Int)

func (curve *CurveParams) IsOnCurve(x, y *big.Int) bool

func (curve *CurveParams) Params() *CurveParams

func (curve *CurveParams) ScalarBaseMult(k []byte) (*big.Int, *big.Int)

func (curve *CurveParams) ScalarMult(Bx, By *big.Int, k []byte) (*big.Int, *big.Int)

