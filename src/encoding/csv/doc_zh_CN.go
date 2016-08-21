// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package csv reads and writes comma-separated values (CSV) files.
//
// A csv file contains zero or more records of one or more fields per record.
// Each record is separated by the newline character. The final record may
// optionally be followed by a newline character.
//
//     field1,field2,field3
//
// White space is considered part of a field.
//
// Carriage returns before newline characters are silently removed.
//
// Blank lines are ignored.  A line with only whitespace characters (excluding
// the ending newline character) is not considered a blank line.
//
// Fields which start and stop with the quote character " are called
// quoted-fields.  The beginning and ending quote are not part of the
// field.
//
// The source:
//
//     normal string,"quoted-field"
//
// results in the fields
//
//     {`normal string`, `quoted-field`}
//
// Within a quoted-field a quote character followed by a second quote
// character is considered a single quote.
//
//     "the ""word"" is true","a ""quoted-field"""
//
// results in
//
//     {`the "word" is true`, `a "quoted-field"`}
//
// Newlines and commas may be included in a quoted-field
//
//     "Multi-line
//     field","comma is ,"
//
// results in
//
//     {`Multi-line
//     field`, `comma is ,`}

// csv读写逗号分隔值（csv）的文件。
//
// 一个csv分拣包含零到多条记录，每条记录一到多个字段。每个记录用换行符分隔。最后
// 一条记录后面可以有换行符，也可以没有。
//
//     field1,field2,field3
//
// 空白视为字段的一部分。
//
// 换行符前面的回车符会悄悄的被删掉。
//
// 忽略空行。只有空白的行（除了末尾换行符之外）不视为空行。
//
// 以双引号"开始和结束的字段成为受引字段，其开始和结束的引号不属于字段。
//
// 资源：
//
//     normal string,"quoted-field"
//
// 产生如下字段：
//
//     {`normal string`, `quoted-field`}
//
// 受引字段内部，如果有两个连续的双引号，则视为一个单纯的双引号字符：
//
//     "the ""word"" is true","a ""quoted-field"""
//
// 生成：
//
//     {`the "word" is true`, `a "quoted-field"`}
//
// 受引字段里可以包含换行和逗号：
//
//     "Multi-line
//     field","comma is ,"
//
// 生成：
//
//     {`Multi-line
//     field`, `comma is ,`}
package csv

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "io"
    "strings"
    "unicode"
    "unicode/utf8"
)

// These are the errors that can be returned in ParseError.Error

// 这些都是PaserError.Err字段可能的值。
var (
    ErrTrailingComma = errors.New("extra delimiter at end of line") // no longer used
    ErrBareQuote     = errors.New("bare \" in non-quoted-field")
    ErrQuote         = errors.New("extraneous \" in field")
    ErrFieldCount    = errors.New("wrong number of fields in line")
)

// A ParseError is returned for parsing errors.
// The first line is 1.  The first column is 0.

// 当解析错误时返回ParseError，第一个行为1，第一列为0。
type ParseError struct {
    Line   int   // Line where the error occurred
    Column int   // Column (rune index) where the error occurred
    Err    error // The actual error
}

// A Reader reads records from a CSV-encoded file.
//
// As returned by NewReader, a Reader expects input conforming to RFC 4180.
// The exported fields can be changed to customize the details before the
// first call to Read or ReadAll.
//
// Comma is the field delimiter.  It defaults to ','.
//
// Comment, if not 0, is the comment character. Lines beginning with the
// Comment character are ignored.
//
// If FieldsPerRecord is positive, Read requires each record to
// have the given number of fields.  If FieldsPerRecord is 0, Read sets it to
// the number of fields in the first record, so that future records must
// have the same field count.  If FieldsPerRecord is negative, no check is
// made and records may have a variable number of fields.
//
// If LazyQuotes is true, a quote may appear in an unquoted field and a
// non-doubled quote may appear in a quoted field.
//
// If TrimLeadingSpace is true, leading white space in a field is ignored.

// Reader从csv编码的文件中读取记录。
//
// NewReader返回的*Reader期望输入符合RFC 4180。在首次调用Read或者ReadAll之前可以
// 设定导出字段的细节。
//
// Comma是字段分隔符，默认为','。Comment如果不是0，则表示注释标识符，以Comment开
// 始的行会被忽略。如果FieldsPerRecord大于0，Read方法要求每条记录都有给定数目的
// 字段。如果FieldsPerRecord等于0，Read方法会将其设为第一条记录的字段数，因此其
// 余的记录必须有同样数目的字段。如果FieldsPerRecord小于0，不会检查字段数，允许
// 记录有不同数量的字段。如果LazyQuotes为真，引号可以出现在非受引字段里，不连续
// 两个的引号可以出现在受引字段里。如果TrimLeadingSpace为真，字段前导的空白会忽
// 略掉。
type Reader struct {
    Comma            rune // field delimiter (set to ',' by NewReader)
    Comment          rune // comment character for start of line
    FieldsPerRecord  int  // number of expected fields per record
    LazyQuotes       bool // allow lazy quotes
    TrailingComma    bool // ignored; here for backwards compatibility
    TrimLeadingSpace bool // trim leading space

}

// A Writer writes records to a CSV encoded file.
//
// As returned by NewWriter, a Writer writes records terminated by a
// newline and uses ',' as the field delimiter.  The exported fields can be
// changed to customize the details before the first call to Write or WriteAll.
//
// Comma is the field delimiter.
//
// If UseCRLF is true, the Writer ends each record with \r\n instead of \n.

// Writer类型的值将记录写入一个csv编码的文件。
//
// NewWriter返回的*Writer写入记录时，以换行结束记录，用','分隔字段。在第一次调用
// Write或WriteAll之前，可以设置导出字段的细节。
//
// Comma是字段分隔符。如果UseCRLF为真，Writer在每条记录结束时用\r\n代替\n。
type Writer struct {
    Comma   rune // Field delimiter (set to ',' by NewWriter)
    UseCRLF bool // True to use \r\n as the line terminator

}

// NewReader returns a new Reader that reads from r.

// NewReader函数返回一个从r读取的*Reader。
func NewReader(r io.Reader) *Reader

// NewWriter returns a new Writer that writes to w.

// NewWriter返回一个写入w的*Writer。
func NewWriter(w io.Writer) *Writer

func (*ParseError) Error() string

// Read reads one record from r.  The record is a slice of strings with each
// string representing one field.

// Read从r读取一条记录，返回值record是字符串的切片，每个字符串代表一个字段。
func (*Reader) Read() (record []string, err error)

// ReadAll reads all the remaining records from r.
// Each record is a slice of fields.
// A successful call returns err == nil, not err == io.EOF. Because ReadAll is
// defined to read until EOF, it does not treat end of file as an error to be
// reported.

// ReadAll从r中读取所有剩余的记录，每个记录都是字段的切片，成功的调用返回值err为
// nil而不是EOF。因为ReadAll方法定义为读取直到文件结尾，因此它不会将文件结尾视为
// 应该报告的错误。
func (*Reader) ReadAll() (records [][]string, err error)

// Error reports any error that has occurred during a previous Write or Flush.

// Error返回在之前的Write方法和Flush方法执行时出现的任何错误。
func (*Writer) Error() error

// Flush writes any buffered data to the underlying io.Writer.
// To check if an error occurred during the Flush, call Error.

// 将缓存中的数据写入底层的io.Writer。要检查Flush时是否发生错误的话，应调用Error
// 方法。
func (*Writer) Flush()

// Writer writes a single CSV record to w along with any necessary quoting.
// A record is a slice of strings with each string being one field.

// 向w中写入一条记录，会自行添加必需的引号。记录是字符串切片，每个字符串代表一个
// 字段。
func (*Writer) Write(record []string) (err error)

// WriteAll writes multiple CSV records to w using Write and then calls Flush.

// WriteAll方法使用Write方法向w写入多条记录，并在最后调用Flush方法清空缓存。
func (*Writer) WriteAll(records [][]string) (err error)

