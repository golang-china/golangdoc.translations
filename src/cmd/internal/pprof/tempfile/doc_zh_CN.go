// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package tempfile provides tools to create and delete temporary files
package tempfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Cleanup removes any temporary files or directories selected for deferred
// cleaning. Similar to defer semantics, the nodes are deleted in LIFO order.

// Cleanup removes any temporary files selected for deferred cleaning.
func Cleanup()

// DeferDelete marks a file or directory to be deleted by next call to Cleanup.

// DeferDelete marks a file to be deleted by next call to Cleanup()
func DeferDelete(path string)

// New returns an unused filename for output files.
func New(dir, prefix, suffix string) (*os.File, error)

