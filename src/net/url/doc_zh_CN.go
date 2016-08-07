// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package url parses URLs and implements query escaping.
// See RFC 3986.

// url包解析URL并实现了查询的逸码，参见RFC 3986。
package url

import (
    "bytes"
    "errors"
    "fmt"
    "sort"
    "strconv"
    "strings"
)

// Error reports an error and the operation and URL that caused it.

// Error会报告一个错误，以及导致该错误发生的URL和操作。
type Error struct {
    Op  string
    URL string
    Err error
}

type EscapeError string

// A URL represents a parsed URL (technically, a URI reference). The general
// form represented is:
//
//     scheme://[userinfo@]host/path[?query][#fragment]
//
// URLs that do not start with a slash after the scheme are interpreted as:
//
//     scheme:opaque[?query][#fragment]
//
// Note that the Path field is stored in decoded form: /%47%6f%2f becomes /Go/.
// A consequence is that it is impossible to tell which slashes in the Path were
// slashes in the raw URL and which were %2f. This distinction is rarely
// important, but when it is, code must not use Path directly.
//
// Go 1.5 introduced the RawPath field to hold the encoded form of Path. The
// Parse function sets both Path and RawPath in the URL it returns, and URL's
// String method uses RawPath if it is a valid encoding of Path, by calling the
// EscapedPath method.
//
// In earlier versions of Go, the more indirect workarounds were that an HTTP
// server could consult req.RequestURI and an HTTP client could construct a URL
// struct directly and set the Opaque field instead of Path. These still work as
// well.

// URL类型代表一个解析后的URL（或者说，一个URL参照）。URL基本格式如下：
//
//     scheme://[userinfo@]host/path[?query][#fragment]
//
// scheme后不是冒号加双斜线的URL被解释为如下格式：
//
//     scheme:opaque[?query][#fragment]
//
// 注意路径字段是以解码后的格式保存的，如/%47%6f%2f会变成/Go/。这导致我们无法确
// 定Path字段中的斜线是来自原始URL还是解码前的%2f。除非一个客户端必须使用其他程
// 序/函数来解析原始URL或者重构原始URL，这个区别并不重要。此时，HTTP服务端可以查
// 询req.RequestURI，而HTTP客户端可以使用URL{Host: "example.com", Opaque:
// "//example.com/Go%2f"}代替{Host: "example.com", Path: "/Go/"}。
type URL struct {
    Scheme   string
    Opaque   string    // encoded opaque data
    User     *Userinfo // username and password information
    Host     string    // host or host:port
    Path     string
    RawQuery string // encoded query values, without '?'
    Fragment string // fragment for references, without '#'
}

// The Userinfo type is an immutable encapsulation of username and
// password details for a URL. An existing Userinfo value is guaranteed
// to have a username set (potentially empty, as allowed by RFC 2396),
// and optionally a password.

// Userinfo类型是一个URL的用户名和密码细节的一个不可修改的封装。一个真实存在的
// Userinfo值必须保证有用户名（但根据 RFC 2396可以是空字符串）以及一个可选的密码
// 。
type Userinfo struct {
}

// Values maps a string key to a list of values.
// It is typically used for query parameters and form values.
// Unlike in the http.Header map, the keys in a Values map
// are case-sensitive.

// Values将建映射到值的列表。它一般用于查询的参数和表单的属性。不同于http.Header
// 这个字典类型，Values的键是大小写敏感的。
type Values map[string][]string

// Parse parses rawurl into a URL structure. The rawurl may be relative or
// absolute.
func Parse(rawurl string) (url *URL, err error)

// ParseQuery parses the URL-encoded query string and returns a map listing the
// values specified for each key. ParseQuery always returns a non-nil map
// containing all the valid query parameters found; err describes the first
// decoding error encountered, if any.
func ParseQuery(query string) (m Values, err error)

// ParseRequestURI parses rawurl into a URL structure.  It assumes that
// rawurl was received in an HTTP request, so the rawurl is interpreted
// only as an absolute URI or an absolute path.
// The string rawurl is assumed not to have a #fragment suffix.
// (Web browsers strip #fragment before sending the URL to a web server.)

// ParseRequestURI parses rawurl into a URL structure. It assumes that rawurl
// was received in an HTTP request, so the rawurl is interpreted only as an
// absolute URI or an absolute path. The string rawurl is assumed not to have a
// #fragment suffix. (Web browsers strip #fragment before sending the URL to a
// web server.)
func ParseRequestURI(rawurl string) (url *URL, err error)

// QueryEscape escapes the string so it can be safely placed
// inside a URL query.

// QueryEscape函数对s进行转码使之可以安全的用在URL查询里。
func QueryEscape(s string) string

// QueryUnescape does the inverse transformation of QueryEscape, converting
// %AB into the byte 0xAB and '+' into ' ' (space). It returns an error if
// any % is not followed by two hexadecimal digits.

// QueryUnescape函数用于将QueryEscape转码的字符串还原。它会把%AB改为字节0xAB，将
// '+'改为' '。如果有某个%后面未跟两个十六进制数字，本函数会返回错误。
func QueryUnescape(s string) (string, error)

// User returns a Userinfo containing the provided username
// and no password set.

// User函数返回一个用户名设置为username的不设置密码的*Userinfo。
func User(username string) *Userinfo

// UserPassword returns a Userinfo containing the provided username
// and password.
// This functionality should only be used with legacy web sites.
// RFC 2396 warns that interpreting Userinfo this way
// ``is NOT RECOMMENDED, because the passing of authentication
// information in clear text (such as URI) has proven to be a
// security risk in almost every case where it has been used.''

// UserPassword函数返回一个用户名设置为username、密码设置为password的*Userinfo。
//
// 这个函数应该只用于老式的站点，因为风险很大，不建议使用，参见RFC 2396。
func UserPassword(username, password string) *Userinfo

func (*Error) Error() string

// IsAbs reports whether the URL is absolute.

// IsAbs returns true if the URL is absolute.
func (*URL) IsAbs() bool

// Parse parses a URL in the context of the receiver.  The provided URL
// may be relative or absolute.  Parse returns nil, err on parse
// failure, otherwise its return value is the same as ResolveReference.

// Parse parses a URL in the context of the receiver. The provided URL may be
// relative or absolute. Parse returns nil, err on parse failure, otherwise its
// return value is the same as ResolveReference.
func (*URL) Parse(ref string) (*URL, error)

// Query parses RawQuery and returns the corresponding values.
func (*URL) Query() Values

// RequestURI returns the encoded path?query or opaque?query string that would
// be used in an HTTP request for u.
func (*URL) RequestURI() string

// ResolveReference resolves a URI reference to an absolute URI from
// an absolute base URI, per RFC 3986 Section 5.2.  The URI reference
// may be relative or absolute.  ResolveReference always returns a new
// URL instance, even if the returned URL is identical to either the
// base or reference. If ref is an absolute URL, then ResolveReference
// ignores base and returns a copy of ref.

// ResolveReference resolves a URI reference to an absolute URI from an absolute
// base URI, per RFC 3986 Section 5.2. The URI reference may be relative or
// absolute. ResolveReference always returns a new URL instance, even if the
// returned URL is identical to either the base or reference. If ref is an
// absolute URL, then ResolveReference ignores base and returns a copy of ref.
func (*URL) ResolveReference(ref *URL) *URL

// String reassembles the URL into a valid URL string.
// The general form of the result is one of:
//
//     scheme:opaque?query#fragment
//     scheme://userinfo@host/path?query#fragment
//
// If u.Opaque is non-empty, String uses the first form;
// otherwise it uses the second form.
// To obtain the path, String uses u.EscapedPath().
//
// In the second form, the following rules apply:
//     - if u.Scheme is empty, scheme: is omitted.
//     - if u.User is nil, userinfo@ is omitted.
//     - if u.Host is empty, host/ is omitted.
//     - if u.Scheme and u.Host are empty and u.User is nil,
//        the entire scheme://userinfo@host/ is omitted.
//     - if u.Host is non-empty and u.Path begins with a /,
//        the form host/path does not add its own /.
//     - if u.RawQuery is empty, ?query is omitted.
//     - if u.Fragment is empty, #fragment is omitted.

// String reassembles the URL into a valid URL string. The general form of the
// result is one of:
//
//     scheme:opaque
//     scheme://userinfo@host/path?query#fragment
//
// If u.Opaque is non-empty, String uses the first form; otherwise it uses the
// second form.
//
// In the second form, the following rules apply:
//
//     - if u.Scheme is empty, scheme: is omitted.
//     - if u.User is nil, userinfo@ is omitted.
//     - if u.Host is empty, host/ is omitted.
//     - if u.Scheme and u.Host are empty and u.User is nil,
//        the entire scheme://userinfo@host/ is omitted.
//     - if u.Host is non-empty and u.Path begins with a /,
//        the form host/path does not add its own /.
//     - if u.RawQuery is empty, ?query is omitted.
//     - if u.Fragment is empty, #fragment is omitted.
func (*URL) String() string

// Password returns the password in case it is set, and whether it is set.

// 如果设置了密码返回密码和真，否则会返回假。
func (*Userinfo) Password() (string, bool)

// String returns the encoded userinfo information in the standard form
// of "username[:password]".

// String方法返回编码后的用户信息，格式为"username[:password]"。
func (*Userinfo) String() string

// Username returns the username.

// Username方法返回用户名。
func (*Userinfo) Username() string

func (EscapeError) Error() string

// Add adds the value to key. It appends to any existing values associated with
// key.
func (Values) Add(key, value string)

// Del deletes the values associated with key.
func (Values) Del(key string)

// Encode encodes the values into ``URL encoded'' form ("bar=baz&foo=quux")
// sorted by key.
func (Values) Encode() string

// Get gets the first value associated with the given key. If there are no
// values associated with the key, Get returns the empty string. To access
// multiple values, use the map directly.
func (Values) Get(key string) string

// Set sets the key to value. It replaces any existing values.
func (Values) Set(key, value string)

