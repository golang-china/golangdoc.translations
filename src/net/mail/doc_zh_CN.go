// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package mail implements parsing of mail messages.
//
// For the most part, this package follows the syntax as specified by RFC 5322.
// Notable divergences:
//
//	* Obsolete address formats are not parsed, including addresses with
//	  embedded route information.
//	* Group addresses are not parsed.
//	* The full range of spacing (the CFWS syntax element) is not supported,
//	  such as breaking addresses across lines.

// mail 包实现了解析邮件消息的功能.
//
// 大多数情况下，这个包跟着RFC 5322定义的格式。 需要注意的：
//
//	* 过时的地址格式将不能被解析, 包括嵌入路由信息的地址格式。
//	* 组地址不能被解析。
//	* 全范围的空格（CFWS样式元素）不支持，比如使用换行分隔地址。
package mail

var ErrHeaderNotPresent = errors.New("mail: header not in message")

// ParseAddressList parses the given string as a list of addresses.

// ParseAddressList解析给的一列地址字符串
func ParseAddressList(list string) ([]*Address, error)

// Address represents a single mail address. An address such as "Barry Gibbs
// <bg@example.com>" is represented as Address{Name: "Barry Gibbs", Address:
// "bg@example.com"}.

// Address代表单个的邮件地址。 一个地址例如"Barry Gibbs
// <bg@example.com>"代表一个地址 {Name: "Barry
// Gibbs", Address: "bg@example.com"}。
type Address struct {
	Name    string // Proper name; may be empty.
	Address string // user@domain
}

// Parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"

// 解析一个单独的RFC 5322地址，例如 “Barry Gibbs
// <bg@example.com>”
func ParseAddress(address string) (*Address, error)

// String formats the address as a valid RFC 5322 address. If the address's name
// contains non-ASCII characters the name will be rendered according to RFC 2047.

// String格式化一个可视的RFC 5322地址。
// 如果地址名字包含非ASCII字符串，名字就会按照RFC 2047来解析。
func (a *Address) String() string

// A Header represents the key-value pairs in a mail message header.

// Header代表邮件header中的key-value值对。
type Header map[string][]string

// AddressList parses the named header field as a list of addresses.

// AddressList将命名后的头部区域作为一列地址列表解析出来。
func (h Header) AddressList(key string) ([]*Address, error)

// Date parses the Date header field.

// Date解析Date头部区域。
func (h Header) Date() (time.Time, error)

// Get gets the first value associated with the given key. If there are no values
// associated with the key, Get returns "".

// Get获取根据key取出的第一个对应的值。 如果key没有对应的值，返回“”。
func (h Header) Get(key string) string

// A Message represents a parsed mail message.

// Message代表解析后的邮件信息。
type Message struct {
	Header Header
	Body   io.Reader
}

// ReadMessage reads a message from r. The headers are parsed, and the body of the
// message will be available for reading from r.

// ReadMessage从r中读取一个邮件。 头部已经被解析了，而邮件体是可见的。
func ReadMessage(r io.Reader) (msg *Message, err error)
