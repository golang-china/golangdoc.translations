// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package quick implements utility functions to help with black box testing.

// Package quick implements utility
// functions to help with black box
// testing.
package quick

// Check looks for an input to f, any function that returns bool, such that f
// returns false. It calls f repeatedly, with arbitrary values for each argument.
// If f returns false on a given input, Check returns that input as a *CheckError.
// For example:
//
//	func TestOddMultipleOfThree(t *testing.T) {
//		f := func(x int) bool {
//			y := OddMultipleOfThree(x)
//			return y%2 == 1 && y%3 == 0
//		}
//		if err := quick.Check(f, nil); err != nil {
//			t.Error(err)
//		}
//	}

// Check looks for an input to f, any
// function that returns bool, such that f
// returns false. It calls f repeatedly,
// with arbitrary values for each argument.
// If f returns false on a given input,
// Check returns that input as a
// *CheckError. For example:
//
//	func TestOddMultipleOfThree(t *testing.T) {
//		f := func(x int) bool {
//			y := OddMultipleOfThree(x)
//			return y%2 == 1 && y%3 == 0
//		}
//		if err := quick.Check(f, nil); err != nil {
//			t.Error(err)
//		}
//	}
func Check(f interface{}, config *Config) (err error)

// CheckEqual looks for an input on which f and g return different results. It
// calls f and g repeatedly with arbitrary values for each argument. If f and g
// return different answers, CheckEqual returns a *CheckEqualError describing the
// input and the outputs.

// CheckEqual looks for an input on which f
// and g return different results. It calls
// f and g repeatedly with arbitrary values
// for each argument. If f and g return
// different answers, CheckEqual returns a
// *CheckEqualError describing the input
// and the outputs.
func CheckEqual(f, g interface{}, config *Config) (err error)

// Value returns an arbitrary value of the given type. If the type implements the
// Generator interface, that will be used. Note: To create arbitrary values for
// structs, all the fields must be exported.

// Value returns an arbitrary value of the
// given type. If the type implements the
// Generator interface, that will be used.
// Note: To create arbitrary values for
// structs, all the fields must be
// exported.
func Value(t reflect.Type, rand *rand.Rand) (value reflect.Value, ok bool)

// A CheckEqualError is the result CheckEqual finding an error.

// A CheckEqualError is the result
// CheckEqual finding an error.
type CheckEqualError struct {
	CheckError
	Out1 []interface{}
	Out2 []interface{}
}

func (s *CheckEqualError) Error() string

// A CheckError is the result of Check finding an error.

// A CheckError is the result of Check
// finding an error.
type CheckError struct {
	Count int
	In    []interface{}
}

func (s *CheckError) Error() string

// A Config structure contains options for running a test.

// A Config structure contains options for
// running a test.
type Config struct {
	// MaxCount sets the maximum number of iterations. If zero,
	// MaxCountScale is used.
	MaxCount int
	// MaxCountScale is a non-negative scale factor applied to the default
	// maximum. If zero, the default is unchanged.
	MaxCountScale float64
	// If non-nil, rand is a source of random numbers. Otherwise a default
	// pseudo-random source will be used.
	Rand *rand.Rand
	// If non-nil, the Values function generates a slice of arbitrary
	// reflect.Values that are congruent with the arguments to the function
	// being tested. Otherwise, the top-level Values function is used
	// to generate them.
	Values func([]reflect.Value, *rand.Rand)
}

// A Generator can generate random values of its own type.

// A Generator can generate random values
// of its own type.
type Generator interface {
	// Generate returns a random instance of the type on which it is a
	// method using the size as a size hint.
	Generate(rand *rand.Rand, size int) reflect.Value
}

// A SetupError is the result of an error in the way that check is being used,
// independent of the functions being tested.

// A SetupError is the result of an error
// in the way that check is being used,
// independent of the functions being
// tested.
type SetupError string

func (s SetupError) Error() string
