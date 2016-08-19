// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package gccgoimporter implements Import for gccgo-generated object files.

// Package gccgoimporter implements Import for gccgo-generated object files.
package gccgoimporter

import (
    "bufio"
    "bytes"
    "debug/elf"
    "errors"
    "fmt"
    "go/constant"
    "go/token"
    "go/types"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
    "text/scanner"
)

// Information about a specific installation of gccgo.
type GccgoInstallation struct {
	// Version of gcc (e.g. 4.8.0).
	GccVersion string

	// Target triple (e.g. x86_64-unknown-linux-gnu).
	TargetTriple string

	// Built-in library paths used by this installation.
	LibPaths []string
}


// An Importer resolves import paths to Packages. The imports map records
// packages already known, indexed by package path.
// An importer must determine the canonical package path and check imports
// to see if it is already present in the map. If so, the Importer can return
// the map entry. Otherwise, the importer must load the package data for the
// given path into a new *Package, record it in imports map, and return the
// package.
type Importer func(imports map[string]*types.Package, path string) (*types.Package, error)


// The gccgo-specific init data for a package.
type InitData struct {
	// Initialization priority of this package relative to other packages.
	// This is based on the maximum depth of the package's dependency graph;
	// it is guaranteed to be greater than that of its dependencies.
	Priority int

	// The list of packages which this package depends on to be initialized,
	// including itself if needed. This is the subset of the transitive closure of
	// the package's dependencies that need initialization.
	Inits []PackageInit
}


// A PackageInit describes an imported package that needs initialization.
type PackageInit struct {
	Name     string // short package name
	InitFunc string // name of init function
	Priority int    // priority of init function, see InitData.Priority
}


func GetImporter(searchpaths []string, initmap map[*types.Package]InitData) Importer

// Return an importer that searches incpaths followed by the gcc installation's
// built-in search paths and the current directory.
func (*GccgoInstallation) GetImporter(incpaths []string, initmap map[*types.Package]InitData) Importer

// Ask the driver at the given path for information for this GccgoInstallation.
func (*GccgoInstallation) InitFromDriver(gccgoPath string) (err error)

// Return the list of export search paths for this GccgoInstallation.
func (*GccgoInstallation) SearchPaths() (paths []string)

