// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package imageutil contains code shared by image-related packages.
package imageutil // import "image/internal/imageutil"

import "image"

// DrawYCbCr draws the YCbCr source image on the RGBA destination image with
// r.Min in dst aligned with sp in src. It reports whether the draw was
// successful. If it returns false, no dst pixels were changed.
//
// This function assumes that r is entirely within dst's bounds and the
// translation of r from dst coordinate space to src coordinate space is
// entirely within src's bounds.
func DrawYCbCr(dst *image.RGBA, r image.Rectangle, src *image.YCbCr, sp image.Point) (ok bool)

