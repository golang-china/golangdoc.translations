// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package gcimporter implements Import for gc-generated object files.

// Package gcimporter implements Import for gc-generated object files.
package gcimporter

import (
    "bufio"
    "encoding/binary"
    "errors"
    "fmt"
    "go/build"
    "go/constant"
    "go/token"
    "go/types"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "text/scanner"
    "unicode"
    "unicode/utf8"
)

// BImportData imports a package from the serialized package data
// and returns the number of bytes consumed and a reference to the package.
// If data is obviously malformed, an error is returned but in
// general it is not recommended to call BImportData on untrusted data.
func BImportData(imports map[string]*types.Package, data []byte, path string) (int, *types.Package, error)

// FindExportData positions the reader r at the beginning of the
// export data section of an underlying GC-created object/archive
// file by reading from it. The reader must be positioned at the
// start of the file before calling this function. The hdr result
// is the string before the export data, either "$$" or "$$B".
func FindExportData(r *bufio.Reader) (hdr string, err error)

// FindPkg returns the filename and unique package id for an import
// path based on package information provided by build.Import (using
// the build.Default build.Context). A relative srcDir is interpreted
// relative to the current working directory.
// If no file was found, an empty filename is returned.
func FindPkg(path, srcDir string) (filename, id string)

// Import imports a gc-generated package given its import path and srcDir, adds
// the corresponding package object to the packages map, and returns the object.
// The packages map must contain all packages already imported.
func Import(packages map[string]*types.Package, path, srcDir string) (pkg *types.Package, err error)

// ImportData imports a package by reading the gc-generated export data,
// adds the corresponding package object to the packages map indexed by id,
// and returns the object.
//
// The packages map must contains all packages already imported. The data
// reader position must be the beginning of the export data section. The
// filename is only used in error messages.
//
// If packages[id] contains the completely imported package, that package
// can be used directly, and there is no need to call this function (but
// there is also no harm but for extra time used).
func ImportData(packages map[string]*types.Package, filename, id string, data io.Reader) (pkg *types.Package, err error)

