// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package tls partially implements TLS 1.2, as specified in RFC 5246.

// tls包实现了TLS 1.2，细节参见RFC 5246。
package tls

import (
	"bytes"
	"container/list"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rc4"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	CurveP256 CurveID = 23
	CurveP384 CurveID = 24
	CurveP521 CurveID = 25
)

const (
	NoClientCert ClientAuthType = iota
	RequestClientCert
	RequireAnyClientCert
	VerifyClientCertIfGiven
	RequireAndVerifyClientCert
)

// A list of the possible cipher suite ids. Taken from
// http://www.iana.org/assignments/tls-parameters/tls-parameters.xml

// 可选的加密组的ID的列表。参见：
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

// A Certificate is a chain of one or more certificates, leaf first.

// Certificate是一个或多个证书的链条，叶证书在最前面。
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

// ClientAuthType declares the policy the server will follow for
// TLS Client Authentication.

// ClientAuthType类型声明服务端将遵循的TLS客户端验证策略。
type ClientAuthType int

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

// ClientSessionCache is a cache of ClientSessionState objects that can be used
// by a client to resume a TLS session with a given server. ClientSessionCache
// implementations should expect to be called concurrently from different
// goroutines.

// ClientSessionCache是ClientSessionState对象的缓存，可以被客户端用于恢复与某个
// 服务端的TLS会话。本类型的实现期望被不同线程并行的调用。
type ClientSessionCache interface {
	// Get searches for a ClientSessionState associated with the given key.
	// On return, ok is true if one was found.
	Get(sessionKey string) (session *ClientSessionState, ok bool)

	// Put adds the ClientSessionState to the cache with the given key.
	Put(sessionKey string, cs *ClientSessionState)
}

// ClientSessionState contains the state needed by clients to resume TLS
// sessions.

// ClientSessionState包含客户端所需的用于恢复TLS会话的状态。
type ClientSessionState struct {
}

// A Config structure is used to configure a TLS client or server.
// After one has been passed to a TLS function it must not be
// modified. A Config may be reused; the tls package will also not
// modify it.

// Config结构类型用于配置TLS客户端或服务端。在本类型的值提供给TLS函数后，就不应
// 再修改该值。Config类型值可能被重用；tls包也不会修改它。
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
}

// A Conn represents a secured connection.
// It implements the net.Conn interface.

// Conn代表一个安全连接。本类型实现了net.Conn接口。
type Conn struct {
}

// ConnectionState records basic TLS details about the connection.

// ConnectionState类型记录连接的基本TLS细节。
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
//
//
// http://www.iana.org/assignments/tls-parameters/tls-parameters.xml#tls-parameters-8

// CurveID是TLS椭圆曲线的标识符的类型。参见：
//
// http://www.iana.org/assignments/tls-parameters/tls-parameters.xml#tls-parameters-8
type CurveID uint16

// Client returns a new TLS client side connection
// using conn as the underlying transport.
// The config cannot be nil: users must set either ServerName or
// InsecureSkipVerify in the config.

// Client使用conn作为下层传输接口返回一个TLS连接的客户端侧。配置参数config必须是
// 非nil的且必须设置了ServerName或者InsecureSkipVerify字段。
func Client(conn net.Conn, config *Config) *Conn

// Dial connects to the given network address using net.Dial
// and then initiates a TLS handshake, returning the resulting
// TLS connection.
// Dial interprets a nil configuration as equivalent to
// the zero configuration; see the documentation of Config
// for the defaults.

// Dial使用net.Dial连接指定的网络和地址，然后发起TLS握手，返回生成的TLS连接。
// Dial会将nil的配置视为零值的配置；参见Config类型的文档获取细节。
func Dial(network, addr string, config *Config) (*Conn, error)

// DialWithDialer connects to the given network address using dialer.Dial and
// then initiates a TLS handshake, returning the resulting TLS connection. Any
// timeout or deadline given in the dialer apply to connection and TLS
// handshake as a whole.
//
// DialWithDialer interprets a nil configuration as equivalent to the zero
// configuration; see the documentation of Config for the defaults.

// DialWithDialer使用dialer.Dial连接指定的网络和地址，然后发起TLS握手，返回生成
// 的TLS连接。dialer中的超时和期限设置会将连接和TLS握手作为一个整体来应用。
//
// DialWithDialer会将nil的配置视为零值的配置；参见Config类型的文档获取细节。
func DialWithDialer(dialer *net.Dialer, network, addr string, config *Config) (*Conn, error)

// Listen creates a TLS listener accepting connections on the
// given network address using net.Listen.
// The configuration config must be non-nil and must include
// at least one certificate or else set GetCertificate.

// 函数创建一个TLS监听器，使用net.Listen函数接收给定地址上的连接。配置参数config
// 必须是非nil的且必须含有至少一个证书。
func Listen(network, laddr string, config *Config) (net.Listener, error)

// LoadX509KeyPair reads and parses a public/private key pair from a pair of
// files. The files must contain PEM encoded data. On successful return,
// Certificate.Leaf will be nil because the parsed form of the certificate is
// not retained.

// LoadX509KeyPair读取并解析一对文件获取公钥和私钥。这些文件必须是PEM编码的。
func LoadX509KeyPair(certFile, keyFile string) (cert Certificate, err error)

// NewLRUClientSessionCache returns a ClientSessionCache with the given
// capacity that uses an LRU strategy. If capacity is < 1, a default capacity
// is used instead.

// 函数使用给出的容量创建一个采用LRU策略的ClientSessionState，如果capacity<1会采
// 用默认容量。
func NewLRUClientSessionCache(capacity int) ClientSessionCache

// NewListener creates a Listener which accepts connections from an inner
// Listener and wraps each connection with Server.
// The configuration config must be non-nil and must include
// at least one certificate or else set GetCertificate.

// 函数创建一个TLS监听器，该监听器接受inner接收到的每一个连接，并调用Server函数
// 包装这些连接。配置参数config必须是非nil的且必须含有至少一个证书。
func NewListener(inner net.Listener, config *Config) net.Listener

// Server returns a new TLS server side connection
// using conn as the underlying transport.
// The configuration config must be non-nil and must include
// at least one certificate or else set GetCertificate.

// Server使用conn作为下层传输接口返回一个TLS连接的服务端侧。配置参数config必须是
// 非nil的且必须含有至少一个证书。
func Server(conn net.Conn, config *Config) *Conn

// X509KeyPair parses a public/private key pair from a pair of
// PEM encoded data. On successful return, Certificate.Leaf will be nil because
// the parsed form of the certificate is not retained.

// X509KeyPair解析一对PEM编码的数据获取公钥和私钥。
func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (cert Certificate, err error)

// BuildNameToCertificate parses c.Certificates and builds c.NameToCertificate
// from the CommonName and SubjectAlternateName fields of each of the leaf
// certificates.

// BuildNameToCertificate解析c.Certificates并将每一个叶证书的CommonName和
// SubjectAlternateName字段用于创建c.NameToCertificate。
func (*Config) BuildNameToCertificate()

// Close closes the connection.

// Close关闭连接。
func (*Conn) Close() error

// ConnectionState returns basic TLS details about the connection.

// ConnectionState返回该连接的基本TLS细节。
func (*Conn) ConnectionState() ConnectionState

// Handshake runs the client or server handshake
// protocol if it has not yet been run.
// Most uses of this package need not call Handshake
// explicitly: the first Read or Write will call it automatically.

// Handshake执行客户端或服务端的握手协议（如果还没有执行的话）。本包的大多数应用
// 不需要显式的调用Handsake方法：第一次Read或Write方法会自动调用本方法。
func (*Conn) Handshake() error

// LocalAddr returns the local network address.

// LocalAddr返回本地网络地址。
func (*Conn) LocalAddr() net.Addr

// OCSPResponse returns the stapled OCSP response from the TLS server, if
// any. (Only valid for client connections.)

// OCSPResponse返回来自服务端的OCSP
// staple回复（如果有）。只有客户端可以使用本方法。
func (*Conn) OCSPResponse() []byte

// Read can be made to time out and return a net.Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.

// Read从连接读取数据，可设置超时，参见SetDeadline和SetReadDeadline。
func (*Conn) Read(b []byte) (n int, err error)

// RemoteAddr returns the remote network address.

// LocalAddr返回远端网络地址。
func (*Conn) RemoteAddr() net.Addr

// SetDeadline sets the read and write deadlines associated with the connection.
// A zero value for t means Read and Write will not time out. After a Write has
// timed out, the TLS state is corrupt and all future writes will return the
// same error.

// SetDeadline设置该连接的读写操作绝对期限。t为Time零值表示不设置超时。在一次
// Write/Read方法超时后，TLS连接状态会被破坏，之后所有的读写操作都会返回同一错误
// 。
func (*Conn) SetDeadline(t time.Time) error

// SetReadDeadline sets the read deadline on the underlying connection.
// A zero value for t means Read will not time out.

// SetReadDeadline设置该连接的读操作绝对期限。t为Time零值表示不设置超时。
func (*Conn) SetReadDeadline(t time.Time) error

// SetWriteDeadline sets the write deadline on the underlying connection. A zero
// value for t means Write will not time out. After a Write has timed out, the
// TLS state is corrupt and all future writes will return the same error.

// SetReadDeadline设置该连接的写操作绝对期限。t为Time零值表示不设置超时。在一次
// Write方法超时后，TLS连接状态会被破坏，之后所有的写操作都会返回同一错误。
func (*Conn) SetWriteDeadline(t time.Time) error

// VerifyHostname checks that the peer certificate chain is valid for
// connecting to host.  If so, it returns nil; if not, it returns an error
// describing the problem.

// VerifyHostname检查用于连接到host的对等实体证书链是否合法。如果合法，它会返回
// nil；否则，会返回一个描述该问题的错误。
func (*Conn) VerifyHostname(host string) error

// Write writes data to the connection.

// Write将数据写入连接，可设置超时，参见SetDeadline和SetWriteDeadline。
func (*Conn) Write(b []byte) (int, error)
