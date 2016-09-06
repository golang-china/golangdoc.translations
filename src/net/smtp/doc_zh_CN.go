// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package smtp implements the Simple Mail Transfer Protocol as defined in RFC
// 5321. It also implements the following extensions:
//
// 	8BITMIME  RFC 1652
// 	AUTH      RFC 2554
// 	STARTTLS  RFC 3207
//
// Additional extensions may be handled by clients.
//
// The smtp package is frozen and not accepting new features. Some external
// packages provide more functionality. See:
//
// 	https://godoc.org/?q=smtp

// Package smtp implements the Simple Mail Transfer Protocol as defined in RFC
// 5321. It also implements the following extensions:
//
// 	8BITMIME  RFC 1652
// 	AUTH      RFC 2554
// 	STARTTLS  RFC 3207
//
// Additional extensions may be handled by clients.
package smtp

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strings"
)

// Auth is implemented by an SMTP authentication mechanism.

// Auth接口应被每一个SMTP认证机制实现。
type Auth interface {
	// Start begins an authentication with a server.
	// It returns the name of the authentication protocol
	// and optionally data to include in the initial AUTH message
	// sent to the server. It can return proto == "" to indicate
	// that the authentication should be skipped.
	// If it returns a non-nil error, the SMTP client aborts
	// the authentication attempt and closes the connection.
	Start(server *ServerInfo) (proto string, toServer []byte, err error)

	// Next continues the authentication. The server has just sent
	// the fromServer data. If more is true, the server expects a
	// response, which Next should return as toServer; otherwise
	// Next should return toServer == nil.
	// If Next returns a non-nil error, the SMTP client aborts
	// the authentication attempt and closes the connection.
	Next(fromServer []byte, more bool) (toServer []byte, err error)
}

// A Client represents a client connection to an SMTP server.

// Client代表一个连接到SMTP服务器的客户端。
type Client struct {
	// Text is the textproto.Conn used by the Client. It is exported to allow
	// for clients to add extensions.
	Text *textproto.Conn
}

// ServerInfo records information about an SMTP server.

// ServerInfo类型记录一个SMTP服务器的信息。
type ServerInfo struct {
	Name string   // SMTP server name
	TLS  bool     // using TLS, with valid certificate for Name
	Auth []string // advertised authentication mechanisms
}

// CRAMMD5Auth returns an Auth that implements the CRAM-MD5 authentication
// mechanism as defined in RFC 2195.
// The returned Auth uses the given username and secret to authenticate
// to the server using the challenge-response mechanism.

// 返回一个实现了CRAM-MD5身份认证机制（参见RFC 2195）的Auth接口。返回的接口使用
// 给出的用户名和密码，采用响应——回答机制与服务端进行身份认证。
func CRAMMD5Auth(username, secret string) Auth

// Dial returns a new Client connected to an SMTP server at addr.
// The addr must include a port, as in "mail.example.com:smtp".

// Dial返回一个连接到地址为addr的SMTP服务器的*Client；addr必须包含端口号。
func Dial(addr string) (*Client, error)

// NewClient returns a new Client using an existing connection and host as a
// server name to be used when authenticating.

// NewClient使用已经存在的连接conn和作为服务器名的host（用于身份认证）来创建一个
// *Client。
func NewClient(conn net.Conn, host string) (*Client, error)

// PlainAuth returns an Auth that implements the PLAIN authentication
// mechanism as defined in RFC 4616.
// The returned Auth uses the given username and password to authenticate
// on TLS connections to host and act as identity. Usually identity will be
// left blank to act as username.

// 返回一个实现了PLAIN身份认证机制（参见RFC 4616）的Auth接口。返回的接口使用给出
// 的用户名和密码，通过TLS连接到主机认证，采用identity为身份管理和行动（通常应设
// identity为""，以便使用username为身份）。
func PlainAuth(identity, username, password, host string) Auth

// SendMail connects to the server at addr, switches to TLS if
// possible, authenticates with the optional mechanism a if possible,
// and then sends an email from address from, to addresses to, with
// message msg.
// The addr must include a port, as in "mail.example.com:smtp".
//
// The addresses in the to parameter are the SMTP RCPT addresses.
//
// The msg parameter should be an RFC 822-style email with headers
// first, a blank line, and then the message body. The lines of msg
// should be CRLF terminated. The msg headers should usually include
// fields such as "From", "To", "Subject", and "Cc".  Sending "Bcc"
// messages is accomplished by including an email address in the to
// parameter but not including it in the msg headers.
//
// The SendMail function and the the net/smtp package are low-level
// mechanisms and provide no support for DKIM signing, MIME
// attachments (see the mime/multipart package), or other mail
// functionality. Higher-level packages exist outside of the standard
// library.

// SendMail连接到addr指定的服务器；如果支持会开启TLS；如果支持会使用a认证身份；
// 然后以from为邮件源地址发送邮件msg到目标地址to。（可以是多个目标地址：群发）
func SendMail(addr string, a Auth, from string, to []string, msg []byte) error

// Auth authenticates a client using the provided authentication mechanism.
// A failed authentication closes the connection.
// Only servers that advertise the AUTH extension support this function.

// Auth使用提供的认证机制进行认证。失败的认证会关闭该连接。只有服务端支持AUTH时
// ，本方法才有效。（但是不支持时，调用会默默的成功）
func (c *Client) Auth(a Auth) error

// Close closes the connection.

// Close关闭连接。
func (c *Client) Close() error

// Data issues a DATA command to the server and returns a writer that
// can be used to write the mail headers and body. The caller should
// close the writer before calling any more methods on c. A call to
// Data must be preceded by one or more calls to Rcpt.

// Data发送DATA指令到服务器并返回一个io.WriteCloser，用于写入邮件信息。调用者必
// 须在调用c的下一个方法之前关闭这个io.WriteCloser。方法必须在一次或多次Rcpt方法
// 之后调用。
func (c *Client) Data() (io.WriteCloser, error)

// Extension reports whether an extension is support by the server.
// The extension name is case-insensitive. If the extension is supported,
// Extension also returns a string that contains any parameters the
// server specifies for the extension.

// Extension返回服务端是否支持某个扩展，扩展名是大小写不敏感的。如果扩展被支持，
// 方法还会返回一个包含指定给该扩展的各个参数的字符串。
func (c *Client) Extension(ext string) (bool, string)

// Hello sends a HELO or EHLO to the server as the given host name.
// Calling this method is only necessary if the client needs control
// over the host name used. The client will introduce itself as "localhost"
// automatically otherwise. If Hello is called, it must be called before
// any of the other methods.

// Hello发送给服务端一个HELO或EHLO命令。本方法只有使用者需要控制使用的本地主机名
// 时才应使用，否则程序会将本地主机名设为“localhost”，Hello方法只能在最开始调
// 用。
func (c *Client) Hello(localName string) error

// Mail issues a MAIL command to the server using the provided email address. If
// the server supports the 8BITMIME extension, Mail adds the BODY=8BITMIME
// parameter. This initiates a mail transaction and is followed by one or more
// Rcpt calls.

// Mail发送MAIL命令和邮箱地址from到服务器。如果服务端支持8BITMIME扩展，本方法会
// 添加BODY=8BITMIME参数。方法初始化一次邮件传输，后应跟1到多个Rcpt方法的调用。
func (c *Client) Mail(from string) error

// Quit sends the QUIT command and closes the connection to the server.

// Quit发送QUIT命令并关闭到服务端的连接。
func (c *Client) Quit() error

// Rcpt issues a RCPT command to the server using the provided email address.
// A call to Rcpt must be preceded by a call to Mail and may be followed by
// a Data call or another Rcpt call.

// Rcpt发送RCPT命令和邮箱地址to到服务器。调用Rcpt方法之前必须调用了Mail方法，之
// 后可以再一次调用Rcpt方法，也可以调用Data方法。
func (c *Client) Rcpt(to string) error

// Reset sends the RSET command to the server, aborting the current mail
// transaction.

// Reset向服务端发送REST命令，中断当前的邮件传输。
func (c *Client) Reset() error

// StartTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.

// StartTLS方法发送STARTTLS命令，并将之后的所有数据往来加密。只有服务器附加了
// STARTTLS扩展，这个方法才有效。
func (c *Client) StartTLS(config *tls.Config) error

// TLSConnectionState returns the client's TLS connection state.
// The return values are their zero values if StartTLS did
// not succeed.
func (c *Client) TLSConnectionState() (state tls.ConnectionState, ok bool)

// Verify checks the validity of an email address on the server.
// If Verify returns nil, the address is valid. A non-nil return
// does not necessarily indicate an invalid address. Many servers
// will not verify addresses for security reasons.

// Verify检查一个邮箱地址在其服务器是否合法，如果合法会返回nil；但非nil的返回值
// 并不代表不合法，因为许多服务器出于安全原因不支持这种查询。
func (c *Client) Verify(addr string) error

