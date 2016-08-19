// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package x509 parses X.509-encoded keys and certificates.

// Package x509 parses X.509-encoded keys and certificates.
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
    "os"
    "runtime"
    "strconv"
    "strings"
    "sync"
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
	_ PEMCipher = iota
	PEMCipherDES
	PEMCipher3DES
	PEMCipherAES128
	PEMCipherAES192
	PEMCipherAES256
)


// ErrUnsupportedAlgorithm results from attempting to perform an operation that
// involves algorithms that are not currently implemented.
var ErrUnsupportedAlgorithm = errors.New("x509: cannot verify signature: algorithm unimplemented")


// IncorrectPasswordError is returned when an incorrect password is detected.
var IncorrectPasswordError = errors.New("x509: decryption password incorrect")


// CertPool is a set of certificates.
type CertPool struct {
	bySubjectKeyId map[string][]int
	byName         map[string][]int
	certs          []*Certificate
}


// A Certificate represents an X.509 certificate.
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

	// UnhandledCriticalExtensions contains a list of extension IDs that
	// were not (fully) processed when parsing. Verify will fail if this
	// slice is non-empty, unless verification is delegated to an OS
	// library which understands all the critical extensions.
	//
	// Users can access these extensions using Extensions and can remove
	// elements from this slice if they believe that they have been
	// handled.
	UnhandledCriticalExtensions []asn1.ObjectIdentifier

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
type CertificateInvalidError struct {
	Cert   *Certificate
	Reason InvalidReason
}


// CertificateRequest represents a PKCS #10, certificate signature request.
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

	// Attributes is the dried husk of a bug and shouldn't be used.
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
type ConstraintViolationError struct{}


// ExtKeyUsage represents an extended set of actions that are valid for a given
// key. Each of the ExtKeyUsage* constants define a unique action.

// ExtKeyUsage represents an extended set of actions that are valid for a given
// key. Each of the ExtKeyUsage* constants define a unique action.
type ExtKeyUsage int


// HostnameError results when the set of authorized names doesn't match the
// requested name.
type HostnameError struct {
	Certificate *Certificate
	Host        string
}


// An InsecureAlgorithmError
type InsecureAlgorithmError SignatureAlgorithm



type InvalidReason int


// KeyUsage represents the set of actions that are valid for a given key. It's
// a bitmap of the KeyUsage* constants.
type KeyUsage int



type PEMCipher int



type PublicKeyAlgorithm int



type SignatureAlgorithm int


// SystemRootsError results when we fail to load the system root certificates.
type SystemRootsError struct {
	Err error
}



type UnhandledCriticalExtension struct{}


// UnknownAuthorityError results when the certificate issuer is unknown
type UnknownAuthorityError struct {
	cert *Certificate
	// hintErr contains an error that may be helpful in determining why an
	// authority wasn't found.
	hintErr error
	// hintCert contains a possible authority certificate that was rejected
	// because of the error in hintErr.
	hintCert *Certificate
}


// VerifyOptions contains parameters for Certificate.Verify. It's a structure
// because other PKIX verification APIs have ended up needing many options.
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
func CreateCertificate(rand io.Reader, template, parent *Certificate, pub, priv interface{}) (cert []byte, err error)

// CreateCertificateRequest creates a new certificate based on a template. The
// following members of template are used: Subject, Attributes,
// SignatureAlgorithm, Extensions, DNSNames, EmailAddresses, and IPAddresses.
// The private key is the private key of the signer.
//
// The returned slice is the certificate request in DER encoding.
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PublicKey and *ecdsa.PublicKey.)

// CreateCertificateRequest creates a new certificate request based on a
// template. The following members of template are used: Subject, Attributes,
// SignatureAlgorithm, Extensions, DNSNames, EmailAddresses, and IPAddresses.
// The private key is the private key of the signer.
//
// The returned slice is the certificate request in DER encoding.
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PublicKey and *ecdsa.PublicKey.)
func CreateCertificateRequest(rand io.Reader, template *CertificateRequest, priv interface{}) (csr []byte, err error)

// DecryptPEMBlock takes a password encrypted PEM block and the password used to
// encrypt it and returns a slice of decrypted DER encoded bytes. It inspects
// the DEK-Info header to determine the algorithm used for decryption. If no
// DEK-Info header is present, an error is returned. If an incorrect password
// is detected an IncorrectPasswordError is returned. Because of deficiencies
// in the encrypted-PEM format, it's not always possible to detect an incorrect
// password. In these cases no error will be returned but the decrypted DER
// bytes will be random noise.
func DecryptPEMBlock(b *pem.Block, password []byte) ([]byte, error)

// EncryptPEMBlock returns a PEM block of the specified type holding the
// given DER-encoded data encrypted with the specified algorithm and
// password.
func EncryptPEMBlock(rand io.Reader, blockType string, data, password []byte, alg PEMCipher) (*pem.Block, error)

// IsEncryptedPEMBlock returns if the PEM block is password encrypted.
func IsEncryptedPEMBlock(b *pem.Block) bool

// MarshalECPrivateKey marshals an EC private key into ASN.1, DER format.
func MarshalECPrivateKey(key *ecdsa.PrivateKey) ([]byte, error)

// MarshalPKCS1PrivateKey converts a private key to ASN.1 DER encoded form.
func MarshalPKCS1PrivateKey(key *rsa.PrivateKey) []byte

// MarshalPKIXPublicKey serialises a public key to DER-encoded PKIX format.
func MarshalPKIXPublicKey(pub interface{}) ([]byte, error)

// NewCertPool returns a new, empty CertPool.
func NewCertPool() *CertPool

// ParseCRL parses a CRL from the given bytes. It's often the case that PEM
// encoded CRLs will appear where they should be DER encoded, so this function
// will transparently handle PEM encoding as long as there isn't any leading
// garbage.
func ParseCRL(crlBytes []byte) (*pkix.CertificateList, error)

// ParseCertificate parses a single certificate from the given ASN.1 DER data.
func ParseCertificate(asn1Data []byte) (*Certificate, error)

// ParseCertificateRequest parses a single certificate request from the
// given ASN.1 DER data.
func ParseCertificateRequest(asn1Data []byte) (*CertificateRequest, error)

// ParseCertificates parses one or more certificates from the given ASN.1 DER
// data. The certificates must be concatenated with no intermediate padding.
func ParseCertificates(asn1Data []byte) ([]*Certificate, error)

// ParseDERCRL parses a DER encoded CRL from the given bytes.
func ParseDERCRL(derBytes []byte) (*pkix.CertificateList, error)

// ParseECPrivateKey parses an ASN.1 Elliptic Curve Private Key Structure.
func ParseECPrivateKey(der []byte) (*ecdsa.PrivateKey, error)

// ParsePKCS1PrivateKey returns an RSA private key from its ASN.1 PKCS#1 DER
// encoded form.
func ParsePKCS1PrivateKey(der []byte) (*rsa.PrivateKey, error)

// ParsePKCS8PrivateKey parses an unencrypted, PKCS#8 private key. See
// http://www.rsa.com/rsalabs/node.asp?id=2130 and RFC5208.

// ParsePKCS8PrivateKey parses an unencrypted, PKCS#8 private key.
// See RFC 5208.
func ParsePKCS8PrivateKey(der []byte) (key interface{}, err error)

// ParsePKIXPublicKey parses a DER encoded public key. These values are
// typically found in PEM blocks with "BEGIN PUBLIC KEY".

// ParsePKIXPublicKey parses a DER encoded public key. These values are
// typically found in PEM blocks with "BEGIN PUBLIC KEY".
//
// Supported key types include RSA, DSA, and ECDSA. Unknown key
// types result in an error.
//
// On success, pub will be of type *rsa.PublicKey, *dsa.PublicKey,
// or *ecdsa.PublicKey.
func ParsePKIXPublicKey(derBytes []byte) (pub interface{}, err error)

// SystemCertPool returns a copy of the system cert pool.
//
// Any mutations to the returned pool are not written to disk and do
// not affect any other pool.
func SystemCertPool() (*CertPool, error)

// AddCert adds a certificate to a pool.
func (*CertPool) AddCert(cert *Certificate)

// AppendCertsFromPEM attempts to parse a series of PEM encoded certificates.
// It appends any certificates found to s and reports whether any certificates
// were successfully parsed.
//
// On many Linux systems, /etc/ssl/cert.pem will contain the system wide set
// of root CAs in a format suitable for this function.
func (*CertPool) AppendCertsFromPEM(pemCerts []byte) (ok bool)

// Subjects returns a list of the DER-encoded subjects of
// all of the certificates in the pool.
func (*CertPool) Subjects() [][]byte

// CheckCRLSignature checks that the signature in crl is from c.
func (*Certificate) CheckCRLSignature(crl *pkix.CertificateList) error

// CheckSignature verifies that signature is a valid signature over signed from
// c's public key.
func (*Certificate) CheckSignature(algo SignatureAlgorithm, signed, signature []byte) error

// CheckSignatureFrom verifies that the signature on c is a valid signature
// from parent.
func (*Certificate) CheckSignatureFrom(parent *Certificate) error

// CreateCRL returns a DER encoded CRL, signed by this Certificate, that
// contains the given list of revoked certificates.
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
func (*Certificate) Verify(opts VerifyOptions) (chains [][]*Certificate, err error)

// VerifyHostname returns nil if c is a valid certificate for the named host.
// Otherwise it returns an error describing the mismatch.
func (*Certificate) VerifyHostname(h string) error

// CheckSignature verifies that the signature on c is a valid signature

// CheckSignature reports whether the signature on c is valid.
func (*CertificateRequest) CheckSignature() error

func (CertificateInvalidError) Error() string

func (ConstraintViolationError) Error() string

func (HostnameError) Error() string

func (InsecureAlgorithmError) Error() string

func (SignatureAlgorithm) String() string

func (SystemRootsError) Error() string

func (UnhandledCriticalExtension) Error() string

func (UnknownAuthorityError) Error() string

