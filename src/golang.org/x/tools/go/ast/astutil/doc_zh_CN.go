// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package astutil contains common utilities for working with the Go AST.

// astutil 包包含工作于 Go AST 的常见实用工具.
package astutil

// AddImport adds the import path to the file f, if absent.

// AddImport 增加导入路径到文件 f, 如果缺席.
func AddImport(fset *token.FileSet, f *ast.File, ipath string) (added bool)

// AddNamedImport adds the import path to the file f, if absent. If name is not
// empty, it is used to rename the import.
//
// For example, calling
//
//	AddNamedImport(fset, f, "pathpkg", "path")
//
// adds
//
//	import pathpkg "path"

// AddNamedImport  增加导入路径到文件 f, 如果缺席.
// 如果 name 非空, 它用作导入重命名.
//
// 例如, 调用
//
//	AddNamedImport(fset, f, "pathpkg", "path")
//
// 增加
//
// import pathpkg "path"
func AddNamedImport(fset *token.FileSet, f *ast.File, name, ipath string) (added bool)

// DeleteImport deletes the import path from the file f, if present.

// DeleteImport 从文件 f 删除导入路径, 如果出席.
func DeleteImport(fset *token.FileSet, f *ast.File, path string) (deleted bool)

// Imports returns the file imports grouped by paragraph.

// Imports 返回文件按段分组的导入.
func Imports(fset *token.FileSet, f *ast.File) [][]*ast.ImportSpec

// NodeDescription returns a description of the concrete type of n suitable for a
// user interface.
//
// TODO(adonovan): in some cases (e.g. Field, FieldList, Ident, StarExpr) we could
// be much more specific given the path to the AST root. Perhaps we should do that.

// NodeDescription 返回 n 的适合用户界面的具体类型描述 .
//
// TODO(adonovan): 在某些情况下 (e.g. Field, FieldList, Ident, StarExpr)
// 我们可能会更具体地给出到 AST root 的路径. 或许我们应该这样做.
func NodeDescription(n ast.Node) string

// PathEnclosingInterval returns the node that encloses the source interval [start,
// end), and all its ancestors up to the AST root.
//
// The definition of "enclosing" used by this function considers additional
// whitespace abutting a node to be enclosed by it. In this example:
//
//	z := x + y // add them
//	     <-A->
//	    <----B----->
//
// the ast.BinaryExpr(+) node is considered to enclose interval B even though its
// [Pos()..End()) is actually only interval A. This behaviour makes user interfaces
// more tolerant of imperfect input.
//
// This function treats tokens as nodes, though they are not included in the
// result. e.g. PathEnclosingInterval("+") returns the enclosing ast.BinaryExpr("x
// + y").
//
// If start==end, the 1-char interval following start is used instead.
//
// The 'exact' result is true if the interval contains only path[0] and perhaps
// some adjacent whitespace. It is false if the interval overlaps multiple children
// of path[0], or if it contains only interior whitespace of path[0]. In this
// example:
//
//	z := x + y // add them
//	  <--C-->     <---E-->
//	    ^
//	    D
//
// intervals C, D and E are inexact. C is contained by the z-assignment statement,
// because it spans three of its children (:=, x, +). So too is the 1-char interval
// D, because it contains only interior whitespace of the assignment. E is
// considered interior whitespace of the BlockStmt containing the assignment.
//
// Precondition: [start, end) both lie within the same file as root.
// TODO(adonovan): return (nil, false) in this case and remove precond. Requires
// FileSet; see loader.tokenFileContainsPos.
//
// Postcondition: path is never nil; it always contains at least 'root'.

// PathEnclosingInterval 返回源码区间 [start, end) 内围蔽节点,
// 和其所有上至 AST root 的祖先.
//
// 围蔽出于这种功能考虑, 一个节点被邻近的额外空白字符包围. 在这个例子中:
//
//	z := x + y // add them
//	     <-A->
//	    <----B----->
//
// 认为节点 ast.BinaryExpr(+) 被 B 区间包围, 尽管事实上其 [Pos()..End())
// 只是 A 区间. 此行为让不完美的用户输入界面更宽松.
//
// 此功能像 nodes 一样对待 tokens, 尽管结果中不包括它们.
// e.g. PathEnclosingInterval("+") 返回围蔽 ast.BinaryExpr("x + y").
//
// 如果 start==end, 使用 start 后 1-char 的区间替代.
//
// 如果区间只包含 path[0] 和一些相邻的空白字符, 返回值 'exact' 为 true.
// 如果区间重合 path[0] 多个孩子或只包含 path[0] 内的空白字符, 它的值为 false.
// 例子:
//
//	z := x + y // add them
//	  <--C-->     <---E-->
//	    ^
//	    D
//
// C, D 和 E 区间是不确切的.
// C 被 z 赋值语句包含, 因为它跨越三个孩子 (:=, x, +).
// 所以 1-char 区间 D 也是这样, 因为它只包含赋值语句内的空白字符.
// E 被认为含有赋值的 BlockStmt 内的白字符.
//
// 先决条件: [start, end) 都位于相同文件 root 之下.
// TODO(adonovan): return (nil, false) in this case and remove precond. Requires
// FileSet; see loader.tokenFileContainsPos.
//
// 后决条件: path 不为 nil; 它总是至少包含 'root'.
func PathEnclosingInterval(root *ast.File, start, end token.Pos) (path []ast.Node, exact bool)

// RewriteImport rewrites any import of path oldPath to path newPath.

// RewriteImport 重写导入 oldPath 路径到 newPath 路径.
func RewriteImport(fset *token.FileSet, f *ast.File, oldPath, newPath string) (rewrote bool)

// Unparen returns e with any enclosing parentheses stripped.

// Unparen 返回剥离内部圆括号后的 e.
func Unparen(e ast.Expr) ast.Expr

// UsesImport reports whether a given import is used.

// UsesImport 报告给定的导入路径是否被使用.
func UsesImport(f *ast.File, path string) (used bool)
