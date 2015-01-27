// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package tls partially implements TLS 1.2, as specified in RFC 5246.
package tls

// A list of the possible cipher suite ids. Taken from
// http://www.iana.org/assignments/tls-parameters/tls-parameters.xml
const (
	TLS_RSA_WITH_RC4_128_SHA                uint16 = 0x0005
	TLS_RSA_WITH_3DES_EDE_CBC_SHA           uint16 = 0x000a
	TLS_RSA_WITH_AES_128_CBC_SHA            uint16 = 0x002f
	TLS_RSA_WITH_AES_256_CBC_SHA            uint16 = 0x0035
	TLS_ECDHE_ECDSA_WITH_RC4_128_SHA        uint16 = 0xc007
	TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA    uint16 = 0xc009
	TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA    uint16 = 0xc00a
	TLS_ECDHE_RSA_WITH_RC4_128_SHA          uint16 = 0xc011
	TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA     uint16 = 0xc012
	TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA      uint16 = 0xc013
	TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA      uint16 = 0xc014
	TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256   uint16 = 0xc02f
	TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 uint16 = 0xc02b

	// TLS_FALLBACK_SCSV isn't a standard cipher suite but an indicator
	// that the client is doing version fallback. See
	// https://tools.ietf.org/html/draft-ietf-tls-downgrade-scsv-00.
	TLS_FALLBACK_SCSV uint16 = 0x5600
)

const (
	VersionSSL30 = 0x0300
	VersionTLS10 = 0x0301
	VersionTLS11 = 0x0302
	VersionTLS12 = 0x0303
)

// Listen creates a TLS listener accepting connections on the given network address
// using net.Listen. The configuration config must be non-nil and must have at
// least one certificate.
func Listen(network, laddr string, config *Config) (net.Listener, error)

// NewListener creates a Listener which accepts connections from an inner Listener
// and wraps each connection with Server. The configuration config must be non-nil
// and must have at least one certificate.
func NewListener(inner net.Listener, config *Config) net.Listener

// A Certificate is a chain of one or more certificates, leaf first.
type Certificate struct {
	Certificate [][]byte
	// PrivateKey contains the private key corresponding to the public key
	// in Leaf. For a server, this must be a *rsa.PrivateKey or
	// *ecdsa.PrivateKey. For a client doing client authentication, this
	// can be any type that implements crypto.Signer (which includes RSA
	// and ECDSA private keys).
	PrivateKey crypto.PrivateKey
	// OCSPStaple contains an optional OCSP response which will be served
	// to clients that request it.
	OCSPStaple []byte
	// Leaf is the parsed form of the leaf certificate, which may be
	// initialized using x509.ParseCertificate to reduce per-handshake
	// processing for TLS clients doing client authentication. If nil, the
	// leaf certificate will be parsed as needed.
	Leaf *x509.Certificate
}

// LoadX509KeyPair reads and parses a public/private key pair from a pair of files.
// The files must contain PEM encoded data.
func LoadX509KeyPair(certFile, keyFile string) (cert Certificate, err error)

// X509KeyPair parses a public/private key pair from a pair of PEM encoded data.
func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (cert Certificate, err error)

// ClientAuthType declares the policy the server will follow for TLS Client
// Authentication.
type ClientAuthType int

const (
	NoClientCert ClientAuthType = iota
	RequestClientCert
	RequireAnyClientCert
	VerifyClientCertIfGiven
	RequireAndVerifyClientCert
)

// ClientHelloInfo contains information from a ClientHello message in order to
// guide certificate selection in the GetCertificate callback.
type ClientHelloInfo struct {
	// CipherSuites lists the CipherSuites supported by the client (e.g.
	// TLS_RSA_WITH_RC4_128_SHA).
	CipherSuites []uint16

	// ServerName indicates the name of the server requested by the client
	// in order to support virtual hosting. ServerName is only set if the
	// client is using SNI (see
	// http://tools.ietf.org/html/rfc4366#section-3.1).
	ServerName string

	// SupportedCurves lists the elliptic curves supported by the client.
	// SupportedCurves is set only if the Supported Elliptic Curves
	// Extension is being used (see
	// http://tools.ietf.org/html/rfc4492#section-5.1.1).
	SupportedCurves []CurveID

	// SupportedPoints lists the point formats supported by the client.
	// SupportedPoints is set only if the Supported Point Formats Extension
	// is being used (see
	// http://tools.ietf.org/html/rfc4492#section-5.1.2).
	SupportedPoints []uint8
}

// ClientSessionCache is a cache of ClientSessionState objects that can be used by
// a client to resume a TLS session with a given server. ClientSessionCache
// implementations should expect to be called concurrently from different
// goroutines.
type ClientSessionCache interface {
	// Get searches for a ClientSessionState associated with the given key.
	// On return, ok is true if one was found.
	Get(sessionKey string) (session *ClientSessionState, ok bool)

	// Put adds the ClientSessionState to the cache with the given key.
	Put(sessionKey string, cs *ClientSessionState)
}

// NewLRUClientSessionCache returns a ClientSessionCache with the given capacity
// that uses an LRU strategy. If capacity is < 1, a default capacity is used
// instead.
func NewLRUClientSessionCache(capacity int) ClientSessionCache

// ClientSessionState contains the state needed by clients to resume TLS sessions.
type ClientSessionState struct {
	// contains filtered or unexported fields
}

// A Config structure is used to configure a TLS client or server. After one has
// been passed to a TLS function it must not be modified. A Config may be reused;
// the tls package will also not modify it.
type Config struct {
	// Rand provides the source of entropy for nonces and RSA blinding.
	// If Rand is nil, TLS uses the cryptographic random reader in package
	// crypto/rand.
	// The Reader must be safe for use by multiple goroutines.
	Rand io.Reader

	// Time returns the current time as the number of seconds since the epoch.
	// If Time is nil, TLS uses time.Now.
	Time func() time.Time

	// Certificates contains one or more certificate chains
	// to present to the other side of the connection.
	// Server configurations must include at least one certificate.
	Certificates []Certificate

	// NameToCertificate maps from a certificate name to an element of
	// Certificates. Note that a certificate name can be of the form
	// '*.example.com' and so doesn't have to be a domain name as such.
	// See Config.BuildNameToCertificate
	// The nil value causes the first element of Certificates to be used
	// for all connections.
	NameToCertificate map[string]*Certificate

	// GetCertificate returns a Certificate based on the given
	// ClientHelloInfo. If GetCertificate is nil or returns nil, then the
	// certificate is retrieved from NameToCertificate. If
	// NameToCertificate is nil, the first element of Certificates will be
	// used.
	GetCertificate func(clientHello *ClientHelloInfo) (*Certificate, error)

	// RootCAs defines the set of root certificate authorities
	// that clients use when verifying server certificates.
	// If RootCAs is nil, TLS uses the host's root CA set.
	RootCAs *x509.CertPool

	// NextProtos is a list of supported, application level protocols.
	NextProtos []string

	// ServerName is used to verify the hostname on the returned
	// certificates unless InsecureSkipVerify is given. It is also included
	// in the client's handshake to support virtual hosting.
	ServerName string

	// ClientAuth determines the server's policy for
	// TLS Client Authentication. The default is NoClientCert.
	ClientAuth ClientAuthType

	// ClientCAs defines the set of root certificate authorities
	// that servers use if required to verify a client certificate
	// by the policy in ClientAuth.
	ClientCAs *x509.CertPool

	// InsecureSkipVerify controls whether a client verifies the
	// server's certificate chain and host name.
	// If InsecureSkipVerify is true, TLS accepts any certificate
	// presented by the server and any host name in that certificate.
	// In this mode, TLS is susceptible to man-in-the-middle attacks.
	// This should be used only for testing.
	InsecureSkipVerify bool

	// CipherSuites is a list of supported cipher suites. If CipherSuites
	// is nil, TLS uses a list of suites supported by the implementation.
	CipherSuites []uint16

	// PreferServerCipherSuites controls whether the server selects the
	// client's most preferred ciphersuite, or the server's most preferred
	// ciphersuite. If true then the server's preference, as expressed in
	// the order of elements in CipherSuites, is used.
	PreferServerCipherSuites bool

	// SessionTicketsDisabled may be set to true to disable session ticket
	// (resumption) support.
	SessionTicketsDisabled bool

	// SessionTicketKey is used by TLS servers to provide session
	// resumption. See RFC 5077. If zero, it will be filled with
	// random data before the first server handshake.
	//
	// If multiple servers are terminating connections for the same host
	// they should all have the same SessionTicketKey. If the
	// SessionTicketKey leaks, previously recorded and future TLS
	// connections using that key are compromised.
	SessionTicketKey [32]byte

	// SessionCache is a cache of ClientSessionState entries for TLS session
	// resumption.
	ClientSessionCache ClientSessionCache

	// MinVersion contains the minimum SSL/TLS version that is acceptable.
	// If zero, then SSLv3 is taken as the minimum.
	MinVersion uint16

	// MaxVersion contains the maximum SSL/TLS version that is acceptable.
	// If zero, then the maximum version supported by this package is used,
	// which is currently TLS 1.2.
	MaxVersion uint16

	// CurvePreferences contains the elliptic curves that will be used in
	// an ECDHE handshake, in preference order. If empty, the default will
	// be used.
	CurvePreferences []CurveID
	// contains filtered or unexported fields
}

// BuildNameToCertificate parses c.Certificates and builds c.NameToCertificate from
// the CommonName and SubjectAlternateName fields of each of the leaf certificates.
func (c *Config) BuildNameToCertificate()

// A Conn represents a secured connection. It implements the net.Conn interface.
type Conn struct {
	// contains filtered or unexported fields
}

// Client returns a new TLS client side connection using conn as the underlying
// transport. The config cannot be nil: users must set either ServerName or
// InsecureSkipVerify in the config.
func Client(conn net.Conn, config *Config) *Conn

// Dial connects to the given network address using net.Dial and then initiates a
// TLS handshake, returning the resulting TLS connection. Dial interprets a nil
// configuration as equivalent to the zero configuration; see the documentation of
// Config for the defaults.
func Dial(network, addr string, config *Config) (*Conn, error)

// DialWithDialer connects to the given network address using dialer.Dial and then
// initiates a TLS handshake, returning the resulting TLS connection. Any timeout
// or deadline given in the dialer apply to connection and TLS handshake as a
// whole.
//
// DialWithDialer interprets a nil configuration as equivalent to the zero
// configuration; see the documentation of Config for the defaults.
func DialWithDialer(dialer *net.Dialer, network, addr string, config *Config) (*Conn, error)

// Server returns a new TLS server side connection using conn as the underlying
// transport. The configuration config must be non-nil and must have at least one
// certificate.
func Server(conn net.Conn, config *Config) *Conn

// Close closes the connection.
func (c *Conn) Close() error

// ConnectionState returns basic TLS details about the connection.
func (c *Conn) ConnectionState() ConnectionState

// Handshake runs the client or server handshake protocol if it has not yet been
// run. Most uses of this package need not call Handshake explicitly: the first
// Read or Write will call it automatically.
func (c *Conn) Handshake() error

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr

// OCSPResponse returns the stapled OCSP response from the TLS server, if any.
// (Only valid for client connections.)
func (c *Conn) OCSPResponse() []byte

// Read can be made to time out and return a net.Error with Timeout() == true after
// a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (n int, err error)

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr

// SetDeadline sets the read and write deadlines associated with the connection. A
// zero value for t means Read and Write will not time out. After a Write has timed
// out, the TLS state is corrupt and all future writes will return the same error.
func (c *Conn) SetDeadline(t time.Time) error

// SetReadDeadline sets the read deadline on the underlying connection. A zero
// value for t means Read will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error

// SetWriteDeadline sets the write deadline on the underlying connection. A zero
// value for t means Write will not time out. After a Write has timed out, the TLS
// state is corrupt and all future writes will return the same error.
func (c *Conn) SetWriteDeadline(t time.Time) error

// VerifyHostname checks that the peer certificate chain is valid for connecting to
// host. If so, it returns nil; if not, it returns an error describing the problem.
func (c *Conn) VerifyHostname(host string) error

// Write writes data to the connection.
func (c *Conn) Write(b []byte) (int, error)

// ConnectionState records basic TLS details about the connection.
type ConnectionState struct {
	Version                    uint16                // TLS version used by the connection (e.g. VersionTLS12)
	HandshakeComplete          bool                  // TLS handshake is complete
	DidResume                  bool                  // connection resumes a previous TLS connection
	CipherSuite                uint16                // cipher suite in use (TLS_RSA_WITH_RC4_128_SHA, ...)
	NegotiatedProtocol         string                // negotiated next protocol (from Config.NextProtos)
	NegotiatedProtocolIsMutual bool                  // negotiated protocol was advertised by server
	ServerName                 string                // server name requested by client, if any (server side only)
	PeerCertificates           []*x509.Certificate   // certificate chain presented by remote peer
	VerifiedChains             [][]*x509.Certificate // verified chains built from PeerCertificates

	// TLSUnique contains the "tls-unique" channel binding value (see RFC
	// 5929, section 3). For resumed sessions this value will be nil
	// because resumption does not include enough context (see
	// https://secure-resumption.com/#channelbindings). This will change in
	// future versions of Go once the TLS master-secret fix has been
	// standardized and implemented.
	TLSUnique []byte
}

// CurveID is the type of a TLS identifier for an elliptic curve. See
// http://www.iana.org/assignments/tls-parameters/tls-parameters.xml#tls-parameters-8
type CurveID uint16

const (
	CurveP256 CurveID = 23
	CurveP384 CurveID = 24
	CurveP521 CurveID = 25
)
