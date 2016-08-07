// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Binary api computes the exported API of a set of Go packages.
//
// The run program is invoked via the dist tool.
// To invoke manually: go tool dist test -run api --no-rebuild

// Binary api computes the exported API of a set of Go packages.
//
// The run program is invoked via the dist tool.
// To invoke manually: go tool dist test -run api --no-rebuild
package main // go get cmd/api

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    "go/ast"
    "go/build"
    "go/parser"
    "go/token"
    "go/types"
    "io"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "runtime"
    "sort"
    "strings"
    "testing"
)

type Walker struct {
    context  *build.Context
    root     string
    scope    []string
    current  *types.Package
    features map[string]bool           // set
    imported map[string]*types.Package // packages already imported
}

func BenchmarkAll(b *testing.B)

func NewWalker(context *build.Context, root string) *Walker

func TestCompareAPI(t *testing.T)

func TestGolden(t *testing.T)

func TestSkipInternal(t *testing.T)

func (*Walker) Features() (fs []string)

func (*Walker) Import(name string) (*types.Package, error)

