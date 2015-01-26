// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package smtp implements the Simple Mail Transfer Protocol as defined in RFC
// 5321. It also implements the following extensions:
//
//	8BITMIME  RFC 1652
//	AUTH      RFC 2554
//	STARTTLS  RFC 3207
//
// Additional extensions may be handled by clients.

// Package smtp implements the Simple Mail
// Transfer Protocol as defined in RFC
// 5321. It also implements the following
// extensions:
//
//	8BITMIME  RFC 1652
//	AUTH      RFC 2554
//	STARTTLS  RFC 3207
//
// Additional extensions may be handled by
// clients.
package smtp

// SendMail connects to the server at addr, switches to TLS if possible,
// authenticates with the optional mechanism a if possible, and then sends an email
// from address from, to addresses to, with message msg.

// SendMail connects to the server at addr,
// switches to TLS if possible,
// authenticates with the optional
// mechanism a if possible, and then sends
// an email from address from, to addresses
// to, with message msg.
func SendMail(addr string, a Auth, from string, to []string, msg []byte) error

// Auth is implemented by an SMTP authentication mechanism.

// Auth is implemented by an SMTP
// authentication mechanism.
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

// CRAMMD5Auth returns an Auth that implements the CRAM-MD5 authentication
// mechanism as defined in RFC 2195. The returned Auth uses the given username and
// secret to authenticate to the server using the challenge-response mechanism.

// CRAMMD5Auth returns an Auth that
// implements the CRAM-MD5 authentication
// mechanism as defined in RFC 2195. The
// returned Auth uses the given username
// and secret to authenticate to the server
// using the challenge-response mechanism.
func CRAMMD5Auth(username, secret string) Auth

// PlainAuth returns an Auth that implements the PLAIN authentication mechanism as
// defined in RFC 4616. The returned Auth uses the given username and password to
// authenticate on TLS connections to host and act as identity. Usually identity
// will be left blank to act as username.

// PlainAuth returns an Auth that
// implements the PLAIN authentication
// mechanism as defined in RFC 4616. The
// returned Auth uses the given username
// and password to authenticate on TLS
// connections to host and act as identity.
// Usually identity will be left blank to
// act as username.
func PlainAuth(identity, username, password, host string) Auth

// A Client represents a client connection to an SMTP server.

// A Client represents a client connection
// to an SMTP server.
type Client struct {
	// Text is the textproto.Conn used by the Client. It is exported to allow for
	// clients to add extensions.
	Text *textproto.Conn
	// contains filtered or unexported fields
}

// Dial returns a new Client connected to an SMTP server at addr. The addr must
// include a port number.

// Dial returns a new Client connected to
// an SMTP server at addr. The addr must
// include a port number.
func Dial(addr string) (*Client, error)

// NewClient returns a new Client using an existing connection and host as a server
// name to be used when authenticating.

// NewClient returns a new Client using an
// existing connection and host as a server
// name to be used when authenticating.
func NewClient(conn net.Conn, host string) (*Client, error)

// Auth authenticates a client using the provided authentication mechanism. A
// failed authentication closes the connection. Only servers that advertise the
// AUTH extension support this function.

// Auth authenticates a client using the
// provided authentication mechanism. A
// failed authentication closes the
// connection. Only servers that advertise
// the AUTH extension support this
// function.
func (c *Client) Auth(a Auth) error

// Close closes the connection.

// Close closes the connection.
func (c *Client) Close() error

// Data issues a DATA command to the server and returns a writer that can be used
// to write the data. The caller should close the writer before calling any more
// methods on c. A call to Data must be preceded by one or more calls to Rcpt.

// Data issues a DATA command to the server
// and returns a writer that can be used to
// write the data. The caller should close
// the writer before calling any more
// methods on c. A call to Data must be
// preceded by one or more calls to Rcpt.
func (c *Client) Data() (io.WriteCloser, error)

// Extension reports whether an extension is support by the server. The extension
// name is case-insensitive. If the extension is supported, Extension also returns
// a string that contains any parameters the server specifies for the extension.

// Extension reports whether an extension
// is support by the server. The extension
// name is case-insensitive. If the
// extension is supported, Extension also
// returns a string that contains any
// parameters the server specifies for the
// extension.
func (c *Client) Extension(ext string) (bool, string)

// Hello sends a HELO or EHLO to the server as the given host name. Calling this
// method is only necessary if the client needs control over the host name used.
// The client will introduce itself as "localhost" automatically otherwise. If
// Hello is called, it must be called before any of the other methods.

// Hello sends a HELO or EHLO to the server
// as the given host name. Calling this
// method is only necessary if the client
// needs control over the host name used.
// The client will introduce itself as
// "localhost" automatically otherwise. If
// Hello is called, it must be called
// before any of the other methods.
func (c *Client) Hello(localName string) error

// Mail issues a MAIL command to the server using the provided email address. If
// the server supports the 8BITMIME extension, Mail adds the BODY=8BITMIME
// parameter. This initiates a mail transaction and is followed by one or more Rcpt
// calls.

// Mail issues a MAIL command to the server
// using the provided email address. If the
// server supports the 8BITMIME extension,
// Mail adds the BODY=8BITMIME parameter.
// This initiates a mail transaction and is
// followed by one or more Rcpt calls.
func (c *Client) Mail(from string) error

// Quit sends the QUIT command and closes the connection to the server.

// Quit sends the QUIT command and closes
// the connection to the server.
func (c *Client) Quit() error

// Rcpt issues a RCPT command to the server using the provided email address. A
// call to Rcpt must be preceded by a call to Mail and may be followed by a Data
// call or another Rcpt call.

// Rcpt issues a RCPT command to the server
// using the provided email address. A call
// to Rcpt must be preceded by a call to
// Mail and may be followed by a Data call
// or another Rcpt call.
func (c *Client) Rcpt(to string) error

// Reset sends the RSET command to the server, aborting the current mail
// transaction.

// Reset sends the RSET command to the
// server, aborting the current mail
// transaction.
func (c *Client) Reset() error

// StartTLS sends the STARTTLS command and encrypts all further communication. Only
// servers that advertise the STARTTLS extension support this function.

// StartTLS sends the STARTTLS command and
// encrypts all further communication. Only
// servers that advertise the STARTTLS
// extension support this function.
func (c *Client) StartTLS(config *tls.Config) error

// Verify checks the validity of an email address on the server. If Verify returns
// nil, the address is valid. A non-nil return does not necessarily indicate an
// invalid address. Many servers will not verify addresses for security reasons.

// Verify checks the validity of an email
// address on the server. If Verify returns
// nil, the address is valid. A non-nil
// return does not necessarily indicate an
// invalid address. Many servers will not
// verify addresses for security reasons.
func (c *Client) Verify(addr string) error

// ServerInfo records information about an SMTP server.

// ServerInfo records information about an
// SMTP server.
type ServerInfo struct {
	Name string   // SMTP server name
	TLS  bool     // using TLS, with valid certificate for Name
	Auth []string // advertised authentication mechanisms
}
