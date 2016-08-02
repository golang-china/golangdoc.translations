// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package testenv provides information about what functionality
// is available in different testing environments run by the Go team.
//
// It is an internal package because these details are specific
// to the Go team's test setup (on build.golang.org) and not
// fundamental to tests in general.
package testenv // import "internal/testenv"

import (
    "os"
    "runtime"
    "strings"
    "testing"
)

// Builder reports the name of the builder running this test
// (for example, "linux-amd64" or "windows-386-gce").
// If the test is not running on the build infrastructure,
// Builder returns the empty string.
func Builder() string

// HasExec reports whether the current system can start new processes
// using os.StartProcess or (more commonly) exec.Command.
func HasExec() bool

// HasExternalNetwork reports whether the current system can use
// external (non-localhost) networks.
func HasExternalNetwork() bool

// HasGoBuild reports whether the current system can build programs with ``go build''
// and then run them with os.StartProcess or exec.Command.
func HasGoBuild() bool

// HasGoRun reports whether the current system can run programs with ``go run.''
func HasGoRun() bool

// MustHaveExec checks that the current system can start new processes
// using os.StartProcess or (more commonly) exec.Command.
// If not, MustHaveExec calls t.Skip with an explanation.
func MustHaveExec(t *testing.T)

// MustHaveExternalNetwork checks that the current system can use
// external (non-localhost) networks.
// If not, MustHaveExternalNetwork calls t.Skip with an explanation.
func MustHaveExternalNetwork(t *testing.T)

// MustHaveGoBuild checks that the current system can build programs with ``go build''
// and then run them with os.StartProcess or exec.Command. If not, MustHaveGoBuild
// calls t.Skip with an explanation.
func MustHaveGoBuild(t *testing.T)

// MustHaveGoRun checks that the current system can run programs with ``go run.''
// If not, MustHaveGoRun calls t.Skip with an explanation.
func MustHaveGoRun(t *testing.T)

