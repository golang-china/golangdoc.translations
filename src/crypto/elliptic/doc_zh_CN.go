// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package elliptic implements several standard elliptic curves over prime fields.
package elliptic

// GenerateKey returns a public/private key pair. The private key is generated
// using the given reader, which must return random data.
func GenerateKey(curve Curve, rand io.Reader) (priv []byte, x, y *big.Int, err error)

// Marshal converts a point into the form specified in section 4.3.6 of ANSI X9.62.
func Marshal(curve Curve, x, y *big.Int) []byte

// Unmarshal converts a point, serialized by Marshal, into an x, y pair. On error,
// x = nil.
func Unmarshal(curve Curve, data []byte) (x, y *big.Int)

// A Curve represents a short-form Weierstrass curve with a=-3. See
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw.html
type Curve interface {
	// Params returns the parameters for the curve.
	Params() *CurveParams
	// IsOnCurve returns true if the given (x,y) lies on the curve.
	IsOnCurve(x, y *big.Int) bool
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

// P224 returns a Curve which implements P-224 (see FIPS 186-3, section D.2.2)
func P224() Curve

// P256 returns a Curve which implements P-256 (see FIPS 186-3, section D.2.3)
func P256() Curve

// P384 returns a Curve which implements P-384 (see FIPS 186-3, section D.2.4)
func P384() Curve

// P521 returns a Curve which implements P-521 (see FIPS 186-3, section D.2.5)
func P521() Curve

// CurveParams contains the parameters of an elliptic curve and also provides a
// generic, non-constant time implementation of Curve.
type CurveParams struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	B       *big.Int // the constant of the curve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
}

func (curve *CurveParams) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int)

func (curve *CurveParams) Double(x1, y1 *big.Int) (*big.Int, *big.Int)

func (curve *CurveParams) IsOnCurve(x, y *big.Int) bool

func (curve *CurveParams) Params() *CurveParams

func (curve *CurveParams) ScalarBaseMult(k []byte) (*big.Int, *big.Int)

func (curve *CurveParams) ScalarMult(Bx, By *big.Int, k []byte) (*big.Int, *big.Int)
