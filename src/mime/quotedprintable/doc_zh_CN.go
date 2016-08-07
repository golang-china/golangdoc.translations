// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package quotedprintable implements quoted-printable encoding as specified by
// RFC 2045.

// Package quotedprintable implements quoted-printable encoding as specified by
// RFC 2045.
package quotedprintable

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
)

// NewReader returns a quoted-printable reader, decoding from r.
func NewReader(r io.Reader) io.Reader

