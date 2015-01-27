// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package symbolz symbolizes a profile using the output from the symbolz service.
package symbolz

// Symbolize symbolizes profile p by parsing data returned by a symbolz handler.
// syms receives the symbolz query (hex addresses separated by '+') and returns the
// symbolz output in a string. It symbolizes all locations based on their
// addresses, regardless of mapping.
func Symbolize(source string, syms func(string, string) ([]byte, error), p *profile.Profile) error
