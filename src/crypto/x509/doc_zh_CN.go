// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package x509 parses X.509-encoded keys and certificates.

// x509包解析X.509编码的证书和密钥。
package x509

import (
    "C"
    "bytes"
    "crypto"
    "crypto/aes"
    "crypto/cipher"
    "crypto/des"
    "crypto/dsa"
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/md5"
    "crypto/rsa"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "crypto/x509/pkix"
    "encoding/asn1"
    "encoding/hex"
    "encoding/pem"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "math/big"
    "net"
    "os/exec"
    "runtime"
    "strconv"
    "strings"
    "sync"
    "syscall"
    "time"
    "unicode/utf8"
    "unsafe"
)

const (
    ExtKeyUsageAny ExtKeyUsage = iota
    ExtKeyUsageServerAuth
    ExtKeyUsageClientAuth
    ExtKeyUsageCodeSigning
    ExtKeyUsageEmailProtection
    ExtKeyUsageIPSECEndSystem
    ExtKeyUsageIPSECTunnel
    ExtKeyUsageIPSECUser
    ExtKeyUsageTimeStamping
    ExtKeyUsageOCSPSigning
    ExtKeyUsageMicrosoftServerGatedCrypto
    ExtKeyUsageNetscapeServerGatedCrypto
)

const (
    KeyUsageDigitalSignature KeyUsage = 1 << iota
    KeyUsageContentCommitment
    KeyUsageKeyEncipherment
    KeyUsageDataEncipherment
    KeyUsageKeyAgreement
    KeyUsageCertSign
    KeyUsageCRLSign
    KeyUsageEncipherOnly
    KeyUsageDecipherOnly
)

const (
    // NotAuthorizedToSign results when a certificate is signed by another
    // which isn't marked as a CA certificate.
    NotAuthorizedToSign InvalidReason = iota
    // Expired results when a certificate has expired, based on the time
    // given in the VerifyOptions.
    Expired
    // CANotAuthorizedForThisName results when an intermediate or root
    // certificate has a name constraint which doesn't include the name
    // being checked.
    CANotAuthorizedForThisName
    // TooManyIntermediates results when a path length constraint is
    // violated.
    TooManyIntermediates
    // IncompatibleUsage results when the certificate's key usage indicates
    // that it may only be used for a different purpose.
    IncompatibleUsage
)

const (
    UnknownPublicKeyAlgorithm PublicKeyAlgorithm = iota
    RSA
    DSA
    ECDSA
)

const (
    UnknownSignatureAlgorithm SignatureAlgorithm = iota
    MD2WithRSA
    MD5WithRSA
    SHA1WithRSA
    SHA256WithRSA
    SHA384WithRSA
    SHA512WithRSA
    DSAWithSHA1
    DSAWithSHA256
    ECDSAWithSHA1
    ECDSAWithSHA256
    ECDSAWithSHA384
    ECDSAWithSHA512
)

// Possible values for the EncryptPEMBlock encryption algorithm.
const (
    _   PEMCipher = iota
    PEMCipherDES
    PEMCipher3DES
    PEMCipherAES128
    PEMCipherAES192
    PEMCipherAES256
)

// ErrUnsupportedAlgorithm results from attempting to perform an operation that
// involves algorithms that are not currently implemented.

// 当试图执行包含目前未实现的算法的操作时，会返回ErrUnsupportedAlgorithm。
//
//     var IncorrectPasswordError = errors.New("x509: decryption password incorrect")
//
// 当检测到不正确的密码时，会返回IncorrectPasswordError。
var ErrUnsupportedAlgorithm = errors.New("x509: cannot verify signature: algorithm unimplemented")

// IncorrectPasswordError is returned when an incorrect password is detected.
var IncorrectPasswordError = errors.New("x509: decryption password incorrect")

// CertPool is a set of certificates.

// CertPool代表一个证书集合/证书池。
type CertPool struct {
}

// A Certificate represents an X.509 certificate.

// Certificate代表一个X.509证书。
type Certificate struct {
    Raw                     []byte // Complete ASN.1 DER content (certificate, signature algorithm and signature).
    RawTBSCertificate       []byte // Certificate part of raw ASN.1 DER content.
    RawSubjectPublicKeyInfo []byte // DER encoded SubjectPublicKeyInfo.
    RawSubject              []byte // DER encoded Subject
    RawIssuer               []byte // DER encoded Issuer

    Signature          []byte
    SignatureAlgorithm SignatureAlgorithm

    PublicKeyAlgorithm PublicKeyAlgorithm
    PublicKey          interface{}

    Version             int
    SerialNumber        *big.Int
    Issuer              pkix.Name
    Subject             pkix.Name
    NotBefore, NotAfter time.Time // Validity bounds.
    KeyUsage            KeyUsage

    // Extensions contains raw X.509 extensions. When parsing certificates,
    // this can be used to extract non-critical extensions that are not
    // parsed by this package. When marshaling certificates, the Extensions
    // field is ignored, see ExtraExtensions.
    Extensions []pkix.Extension

    // ExtraExtensions contains extensions to be copied, raw, into any
    // marshaled certificates. Values override any extensions that would
    // otherwise be produced based on the other fields. The ExtraExtensions
    // field is not populated when parsing certificates, see Extensions.
    ExtraExtensions []pkix.Extension

    ExtKeyUsage        []ExtKeyUsage           // Sequence of extended key usages.
    UnknownExtKeyUsage []asn1.ObjectIdentifier // Encountered extended key usages unknown to this package.

    BasicConstraintsValid bool // if true then the next two fields are valid.
    IsCA                  bool
    MaxPathLen            int
    // MaxPathLenZero indicates that BasicConstraintsValid==true and
    // MaxPathLen==0 should be interpreted as an actual maximum path length
    // of zero. Otherwise, that combination is interpreted as MaxPathLen
    // not being set.
    MaxPathLenZero bool

    SubjectKeyId   []byte
    AuthorityKeyId []byte

    // RFC 5280, 4.2.2.1 (Authority Information Access)
    OCSPServer            []string
    IssuingCertificateURL []string

    // Subject Alternate Name values
    DNSNames       []string
    EmailAddresses []string
    IPAddresses    []net.IP

    // Name constraints
    PermittedDNSDomainsCritical bool // if true then the name constraints are marked critical.
    PermittedDNSDomains         []string

    // CRL Distribution Points
    CRLDistributionPoints []string

    PolicyIdentifiers []asn1.ObjectIdentifier
}

// CertificateInvalidError results when an odd error occurs. Users of this
// library probably want to handle all these errors uniformly.

// 当发生其余的错误时，会返回CertificateInvalidError。本包的使用者可能会想统一处
// 理所有这类错误。
type CertificateInvalidError struct {
    Cert   *Certificate
    Reason InvalidReason
}

// CertificateRequest represents a PKCS #10, certificate signature request.

// CertificateRequest代表一个PKCS #10证书签名请求。
type CertificateRequest struct {
    Raw                      []byte // Complete ASN.1 DER content (CSR, signature algorithm and signature).
    RawTBSCertificateRequest []byte // Certificate request info part of raw ASN.1 DER content.
    RawSubjectPublicKeyInfo  []byte // DER encoded SubjectPublicKeyInfo.
    RawSubject               []byte // DER encoded Subject.

    Version            int
    Signature          []byte
    SignatureAlgorithm SignatureAlgorithm

    PublicKeyAlgorithm PublicKeyAlgorithm
    PublicKey          interface{}

    Subject pkix.Name

    // Attributes is a collection of attributes providing
    // additional information about the subject of the certificate.
    // See RFC 2986 section 4.1.
    Attributes []pkix.AttributeTypeAndValueSET

    // Extensions contains raw X.509 extensions. When parsing CSRs, this
    // can be used to extract extensions that are not parsed by this
    // package.
    Extensions []pkix.Extension

    // ExtraExtensions contains extensions to be copied, raw, into any
    // marshaled CSR. Values override any extensions that would otherwise
    // be produced based on the other fields but are overridden by any
    // extensions specified in Attributes.
    //
    // The ExtraExtensions field is not populated when parsing CSRs, see
    // Extensions.
    ExtraExtensions []pkix.Extension

    // Subject Alternate Name values.
    DNSNames       []string
    EmailAddresses []string
    IPAddresses    []net.IP
}

// ConstraintViolationError results when a requested usage is not permitted by
// a certificate. For example: checking a signature when the public key isn't a
// certificate signing key.

// 当请求的用途不被证书许可时，会返回ConstraintViolationError。如：当公钥不是证
// 书的签名密钥时用它检查签名。
type ConstraintViolationError struct{}

// ExtKeyUsage represents an extended set of actions that are valid for a given
// key. Each of the ExtKeyUsage* constants define a unique action.

// ExtKeyUsage代表给定密钥的合法操作扩展集。每一个ExtKeyUsage类型常数定义一个特
// 定的操作。
//
//     const (
//         ExtKeyUsageAny ExtKeyUsage = iota
//         ExtKeyUsageServerAuth
//         ExtKeyUsageClientAuth
//         ExtKeyUsageCodeSigning
//         ExtKeyUsageEmailProtection
//         ExtKeyUsageIPSECEndSystem
//         ExtKeyUsageIPSECTunnel
//         ExtKeyUsageIPSECUser
//         ExtKeyUsageTimeStamping
//         ExtKeyUsageOCSPSigning
//         ExtKeyUsageMicrosoftServerGatedCrypto
//         ExtKeyUsageNetscapeServerGatedCrypto
//     )
type ExtKeyUsage int

// HostnameError results when the set of authorized names doesn't match the
// requested name.

// 当认证的名字和请求的名字不匹配时，会返回HostnameError。
type HostnameError struct {
    Certificate *Certificate
    Host        string
}

type InvalidReason int

// KeyUsage represents the set of actions that are valid for a given key. It's
// a bitmap of the KeyUsage* constants.

// KeyUsage代表给定密钥的合法操作集。用KeyUsage类型常数的位图表示。（字位表示有
// 无）
//
//     const (
//         KeyUsageDigitalSignature KeyUsage = 1 << iota
//         KeyUsageContentCommitment
//         KeyUsageKeyEncipherment
//         KeyUsageDataEncipherment
//         KeyUsageKeyAgreement
//         KeyUsageCertSign
//         KeyUsageCRLSign
//         KeyUsageEncipherOnly
//         KeyUsageDecipherOnly
//     )
type KeyUsage int

type PEMCipher int

type PublicKeyAlgorithm int

type SignatureAlgorithm int

// SystemRootsError results when we fail to load the system root certificates.

// 当从系统装载根证书失败时，会返回SystemRootsError。
type SystemRootsError struct{}

type UnhandledCriticalExtension struct{}

// UnknownAuthorityError results when the certificate issuer is unknown

// 当证书的发布者未知时，会返回UnknownAuthorityError。
type UnknownAuthorityError struct {
}

// VerifyOptions contains parameters for Certificate.Verify. It's a structure
// because other PKIX verification APIs have ended up needing many options.

// VerifyOptions包含提供给Certificate.Verify方法的参数。它是结构体类型，因为其他
// PKIX认证API需要很长参数。
type VerifyOptions struct {
    DNSName       string
    Intermediates *CertPool
    Roots         *CertPool // if nil, the system roots are used
    CurrentTime   time.Time // if zero, the current time is used
    // KeyUsage specifies which Extended Key Usage values are acceptable.
    // An empty list means ExtKeyUsageServerAuth. Key usage is considered a
    // constraint down the chain which mirrors Windows CryptoAPI behaviour,
    // but not the spec. To accept any key usage, include ExtKeyUsageAny.
    KeyUsages []ExtKeyUsage
}

// CreateCertificate creates a new certificate based on a template. The
// following members of template are used: SerialNumber, Subject, NotBefore,
// NotAfter, KeyUsage, ExtKeyUsage, UnknownExtKeyUsage, BasicConstraintsValid,
// IsCA, MaxPathLen, SubjectKeyId, DNSNames, PermittedDNSDomainsCritical,
// PermittedDNSDomains, SignatureAlgorithm.
//
// The certificate is signed by parent. If parent is equal to template then the
// certificate is self-signed. The parameter pub is the public key of the
// signee and priv is the private key of the signer.
//
// The returned slice is the certificate in DER encoding.
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PublicKey and *ecdsa.PublicKey.)

// CreateCertificate基于模板创建一个新的证书。会用到模板的如下字段：
//
// SerialNumber、Subject、NotBefore、NotAfter、KeyUsage、ExtKeyUsage、
// UnknownExtKeyUsage、
//
// BasicConstraintsValid、IsCA、MaxPathLen、SubjectKeyId、DNSNames、
// PermittedDNSDomainsCritical、
//
// PermittedDNSDomains、SignatureAlgorithm。
//
// 该证书会使用parent签名。如果parent和template相同，则证书是自签名的。Pub参数是
// 被签名者的公钥，而priv是签名者的私钥。
//
// 返回的切片是DER编码的证书。
//
// 只支持RSA和ECDSA类型的密钥。（pub可以是*rsa.PublicKey或*ecdsa.PublicKey，priv
// 可以是*rsa.PrivateKey或*ecdsa.PrivateKey）
func CreateCertificate(rand io.Reader, template, parent *Certificate, pub interface{}, priv interface{}) (cert []byte, err error)

// CreateCertificateRequest creates a new certificate based on a template. The
// following members of template are used: Subject, Attributes,
// SignatureAlgorithm, Extensions, DNSNames, EmailAddresses, and IPAddresses.
// The private key is the private key of the signer.
//
// The returned slice is the certificate request in DER encoding.
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PublicKey and *ecdsa.PublicKey.)

// CreateCertificateRequest基于模板创建一个新的证书请求。会用到模板的如下字段：
//
// Subject、Attributes、Extension、SignatureAlgorithm、DNSNames、EmailAddresses
// 、IPAddresses。
//
// priv是签名者的私钥。返回的切片是DER编码的证书请求。
//
// 只支持RSA（*rsa.PrivateKey）和ECDSA（*ecdsa.PrivateKey）类型的密钥。
func CreateCertificateRequest(rand io.Reader, template *CertificateRequest, priv interface{}) (csr []byte, err error)

// DecryptPEMBlock takes a password encrypted PEM block and the password used to
// encrypt it and returns a slice of decrypted DER encoded bytes. It inspects
// the DEK-Info header to determine the algorithm used for decryption. If no
// DEK-Info header is present, an error is returned. If an incorrect password
// is detected an IncorrectPasswordError is returned. Because of deficiencies
// in the encrypted-PEM format, it's not always possible to detect an incorrect
// password. In these cases no error will be returned but the decrypted DER
// bytes will be random noise.

// DecryptPEMBlock接受一个加密后的PEM块和加密该块的密码password，返回解密后的DER
// 编码字节切片。它会检查DEK信息头域，以确定用于解密的算法。如果b中没有DEK信息头
// 域，会返回错误。如果函数发现密码不正确，会返回IncorrectPasswordError。
func DecryptPEMBlock(b *pem.Block, password []byte) ([]byte, error)

// EncryptPEMBlock returns a PEM block of the specified type holding the
// given DER-encoded data encrypted with the specified algorithm and
// password.

// EncryptPEMBlock使用指定的密码、加密算法加密data，返回一个具有指定块类型，保管
// 加密后数据的PEM块。
func EncryptPEMBlock(rand io.Reader, blockType string, data, password []byte, alg PEMCipher) (*pem.Block, error)

// IsEncryptedPEMBlock returns if the PEM block is password encrypted.

// IsEncryptedPEMBlock返回PEM块b是否是用密码加密了的。
func IsEncryptedPEMBlock(b *pem.Block) bool

// MarshalECPrivateKey marshals an EC private key into ASN.1, DER format.

// MarshalECPrivateKey将ecdsa私钥序列化为ASN.1 DER编码。
func MarshalECPrivateKey(key *ecdsa.PrivateKey) ([]byte, error)

// MarshalPKCS1PrivateKey converts a private key to ASN.1 DER encoded form.

// MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码。
func MarshalPKCS1PrivateKey(key *rsa.PrivateKey) []byte

// MarshalPKIXPublicKey serialises a public key to DER-encoded PKIX format.

// MarshalPKIXPublicKey将公钥序列化为PKIX格式DER编码。
func MarshalPKIXPublicKey(pub interface{}) ([]byte, error)

// NewCertPool returns a new, empty CertPool.

// NewCertPool创建一个新的、空的CertPool。
func NewCertPool() *CertPool

// ParseCRL parses a CRL from the given bytes. It's often the case that PEM
// encoded CRLs will appear where they should be DER encoded, so this function
// will transparently handle PEM encoding as long as there isn't any leading
// garbage.

// ParseCRL从crlBytes中解析CRL（证书注销列表）。因为经常有PEM编码的CRL出现在应该
// 是DER编码的地方，因此本函数可以透明的处理PEM编码，只要没有前导的垃圾数据。
func ParseCRL(crlBytes []byte) (certList *pkix.CertificateList, err error)

// ParseCertificate parses a single certificate from the given ASN.1 DER data.

// ParseCertificate从ASN.1 DER数据解析单个证书。
func ParseCertificate(asn1Data []byte) (*Certificate, error)

// ParseCertificateRequest parses a single certificate request from the
// given ASN.1 DER data.

// ParseCertificateRequest解析一个ASN.1 DER数据获取单个证书请求。
func ParseCertificateRequest(asn1Data []byte) (*CertificateRequest, error)

// ParseCertificates parses one or more certificates from the given ASN.1 DER
// data. The certificates must be concatenated with no intermediate padding.

// ParseCertificates从ASN.1
// DER编码的asn1Data中解析一到多个证书。这些证书必须是串联的，且中间没有填充。
func ParseCertificates(asn1Data []byte) ([]*Certificate, error)

// ParseDERCRL parses a DER encoded CRL from the given bytes.

// ParseDERCRL从derBytes中解析DER编码的CRL。
func ParseDERCRL(derBytes []byte) (certList *pkix.CertificateList, err error)

// ParseECPrivateKey parses an ASN.1 Elliptic Curve Private Key Structure.

// ParseECPrivateKey解析ASN.1 DER编码的ecdsa私钥。
func ParseECPrivateKey(der []byte) (key *ecdsa.PrivateKey, err error)

// ParsePKCS1PrivateKey returns an RSA private key from its ASN.1 PKCS#1 DER
// encoded form.

// ParsePKCS1PrivateKey解析ASN.1 PKCS#1 DER编码的rsa私钥。
func ParsePKCS1PrivateKey(der []byte) (key *rsa.PrivateKey, err error)

// ParsePKCS8PrivateKey parses an unencrypted, PKCS#8 private key. See
// http://www.rsa.com/rsalabs/node.asp?id=2130 and RFC5208.

// ParsePKCS8PrivateKey解析一个未加密的PKCS#8私钥，参见
// http://www.rsa.com/rsalabs/node.asp?id=2130和RFC5208。
func ParsePKCS8PrivateKey(der []byte) (key interface{}, err error)

// ParsePKIXPublicKey parses a DER encoded public key. These values are
// typically found in PEM blocks with "BEGIN PUBLIC KEY".

// ParsePKIXPublicKey解析一个DER编码的公钥。这些公钥一般在以"BEGIN PUBLIC
// KEY"出现的PEM块中。
func ParsePKIXPublicKey(derBytes []byte) (pub interface{}, err error)

// AddCert adds a certificate to a pool.

// AddCert向s中添加一个证书。
func (*CertPool) AddCert(cert *Certificate)

// AppendCertsFromPEM attempts to parse a series of PEM encoded certificates.
// It appends any certificates found to s and reports whether any certificates
// were successfully parsed.
//
// On many Linux systems, /etc/ssl/cert.pem will contain the system wide set
// of root CAs in a format suitable for this function.

// AppendCertsFromPEM试图解析一系列PEM编码的证书。它将找到的任何证书都加入s中，
// 如果所有证书都成功被解析，会返回真。
//
// 在许多Linux系统中，/etc/ssl/cert.pem会包含适合本函数的大量系统级根证书。
func (*CertPool) AppendCertsFromPEM(pemCerts []byte) (ok bool)

// Subjects returns a list of the DER-encoded subjects of
// all of the certificates in the pool.

// Subjects返回池中所有证书的DER编码的持有者的列表。
func (*CertPool) Subjects() (res [][]byte)

// CheckCRLSignature checks that the signature in crl is from c.

// CheckCRLSignature检查crl中的签名是否来自c。
func (*Certificate) CheckCRLSignature(crl *pkix.CertificateList) (err error)

// CheckSignature verifies that signature is a valid signature over signed from
// c's public key.

// CheckSignature检查signature是否是c的公钥生成的signed的合法签名。
func (*Certificate) CheckSignature(algo SignatureAlgorithm, signed, signature []byte) (err error)

// CheckSignatureFrom verifies that the signature on c is a valid signature
// from parent.

// CheckSignatureFrom检查c中的签名是否是来自parent的合法签名。
func (*Certificate) CheckSignatureFrom(parent *Certificate) (err error)

// CreateCRL returns a DER encoded CRL, signed by this Certificate, that
// contains the given list of revoked certificates.

// CreateCRL返回一个DER编码的CRL（证书注销列表），使用c签名，并包含给出的已取消
// 签名列表。
//
// 只支持RSA类型的密钥（priv参数必须是*rsa.PrivateKey类型）。
func (*Certificate) CreateCRL(rand io.Reader, priv interface{}, revokedCerts []pkix.RevokedCertificate, now, expiry time.Time) (crlBytes []byte, err error)

func (*Certificate) Equal(other *Certificate) bool

// Verify attempts to verify c by building one or more chains from c to a
// certificate in opts.Roots, using certificates in opts.Intermediates if
// needed. If successful, it returns one or more chains where the first
// element of the chain is c and the last element is from opts.Roots.
//
// If opts.Roots is nil and system roots are unavailable the returned error
// will be of type SystemRootsError.
//
// WARNING: this doesn't do any revocation checking.

// Verify通过创建一到多个从c到opts.Roots中的证书的链条来认证c，如有必要会使用
// opts.Intermediates中的证书。如果成功，它会返回一到多个证书链条，每一条都以c开
// 始，以opts.Roots中的证书结束。
//
// 警告：它不会做任何取消检查。
func (*Certificate) Verify(opts VerifyOptions) (chains [][]*Certificate, err error)

// VerifyHostname returns nil if c is a valid certificate for the named host.
// Otherwise it returns an error describing the mismatch.
func (*Certificate) VerifyHostname(h string) error

func (CertificateInvalidError) Error() string

func (ConstraintViolationError) Error() string

func (HostnameError) Error() string

func (SystemRootsError) Error() string

func (UnhandledCriticalExtension) Error() string

func (UnknownAuthorityError) Error() string

