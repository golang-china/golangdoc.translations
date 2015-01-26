// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Generate Go pakckage doc for translate.
//
// Usage:
//
//	docgen importPath lang... [-GOOS=...] [-GOARCH=...]
//	docgen -h
//
// Example:
//
//	docgen builtin zh_CN
//	docgen unsafe  zh_CN
//	docgen unsafe  zh_CN zh_TW
//	docgen syscall zh_CN zh_TW -GOOS=windows               # for windows
//	docgen syscall zh_CN zh_TW -GOOS=windows -GOARCH=amd64 # for windows/amd64
//	docgen syscall zh_CN zh_TW                             # for non windows
//
// Output:
//
//	translations/src/builtin/doc_zh_CN.go
//	translations/src/unsafe/doc_zh_CN.go unsafe/doc_zh_TW.go
//	translations/src/unsafe/doc_zh_CN.go
//	translations/src/syscall/doc_zh_CN_windows.go          # for windows
//	translations/src/syscall/doc_zh_CN_windows_amd64.go    # for windows/amd64
//	translations/src/syscall/doc_zh_CN.go                  # for non windows
//
// Help:
//
//	docgen -h

// Generate Go pakckage doc for translate.
//
// Usage:
//
//	docgen importPath lang... [-GOOS=...] [-GOARCH=...]
//	docgen -h
//
// Example:
//
//	docgen builtin zh_CN
//	docgen unsafe  zh_CN
//	docgen unsafe  zh_CN zh_TW
//	docgen syscall zh_CN zh_TW -GOOS=windows               # for windows
//	docgen syscall zh_CN zh_TW -GOOS=windows -GOARCH=amd64 # for windows/amd64
//	docgen syscall zh_CN zh_TW                             # for non windows
//
// Output:
//
//	translations/src/builtin/doc_zh_CN.go
//	translations/src/unsafe/doc_zh_CN.go unsafe/doc_zh_TW.go
//	translations/src/unsafe/doc_zh_CN.go
//	translations/src/syscall/doc_zh_CN_windows.go          # for windows
//	translations/src/syscall/doc_zh_CN_windows_amd64.go    # for windows/amd64
//	translations/src/syscall/doc_zh_CN.go                  # for non windows
//
// Help:
//
//	docgen -h
package main

type PackageInfo struct {
	Lang      string
	FSet      *token.FileSet
	PAst      *ast.Package
	PDoc      *doc.Package
	PDocLocal *doc.Package
	PDocMap   map[string]string
}

func ParsePackageInfo(name, lang string) (pkg *PackageInfo, err error)

func (p *PackageInfo) Bytes() []byte
