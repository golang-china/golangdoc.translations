// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package cookiejar implements an in-memory RFC 6265-compliant http.CookieJar.

// cookiejar包实现了保管在内存中的符合RFC 6265标准的http.CookieJar接口。
package cookiejar

// Jar implements the http.CookieJar interface from the net/http package.

// Jar类型实现了net/http包的http.CookieJar接口。
type Jar struct {
	// contains filtered or unexported fields
}

// New returns a new cookie jar. A nil *Options is equivalent to a zero Options.

// 返回一个新的Jar，nil指针等价于Options零值的指针。
func New(o *Options) (*Jar, error)

// Cookies implements the Cookies method of the http.CookieJar interface.
//
// It returns an empty slice if the URL's scheme is not HTTP or HTTPS.

// 实现CookieJar接口的Cookies方法，如果URL协议不是HTTP/HTTPS会返回空切片。
func (j *Jar) Cookies(u *url.URL) (cookies []*http.Cookie)

// SetCookies implements the SetCookies method of the http.CookieJar interface.
//
// It does nothing if the URL's scheme is not HTTP or HTTPS.

// 实现CookieJar接口的SetCookies方法，如果URL协议不是HTTP/HTTPS则不会有实际操作。
func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie)

// Options are the options for creating a new Jar.

// Options是创建新Jar是的选项。
type Options struct {
	// PublicSuffixList is the public suffix list that determines whether
	// an HTTP server can set a cookie for a domain.
	//
	// A nil value is valid and may be useful for testing but it is not
	// secure: it means that the HTTP server for foo.co.uk can set a cookie
	// for bar.co.uk.
	PublicSuffixList PublicSuffixList
}

// PublicSuffixList provides the public suffix of a domain. For example:
//
//	- the public suffix of "example.com" is "com",
//	- the public suffix of "foo1.foo2.foo3.co.uk" is "co.uk", and
//	- the public suffix of "bar.pvt.k12.ma.us" is "pvt.k12.ma.us".
//
// Implementations of PublicSuffixList must be safe for concurrent use by multiple
// goroutines.
//
// An implementation that always returns "" is valid and may be useful for testing
// but it is not secure: it means that the HTTP server for foo.com can set a cookie
// for bar.com.
//
// A public suffix list implementation is in the package
// golang.org/x/net/publicsuffix.

// PublicSuffixList提供域名的公共后缀。例如：
//
//	- "example.com"的公共后缀是"com"
//	- "foo1.foo2.foo3.co.uk"的公共后缀是"co.uk"
//	- "bar.pvt.k12.ma.us"的公共后缀是"pvt.k12.ma.us"
//
// PublicSuffixList接口的实现必须是并发安全的。一个总是返回""的实现是合法的，也可以通过测试；但却是不安全的：它允许HTTP服务端跨域名设置cookie。推荐实现：code.google.com/p/go.net/publicsuffix
type PublicSuffixList interface {
	// PublicSuffix returns the public suffix of domain.
	//
	// TODO: specify which of the caller and callee is responsible for IP
	// addresses, for leading and trailing dots, for case sensitivity, and
	// for IDN/Punycode.
	PublicSuffix(domain string) string

	// String returns a description of the source of this public suffix
	// list. The description will typically contain something like a time
	// stamp or version number.
	String() string
}
