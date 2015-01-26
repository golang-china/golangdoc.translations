// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package x509 parses X.509-encoded keys and certificates.

// Package x509 parses X.509-encoded keys
// and certificates.
package x509

// ErrUnsupportedAlgorithm results from attempting to perform an operation that
// involves algorithms that are not currently implemented.

// ErrUnsupportedAlgorithm results from
// attempting to perform an operation that
// involves algorithms that are not
// currently implemented.
var ErrUnsupportedAlgorithm = errors.New("x509: cannot verify signature: algorithm unimplemented")

// IncorrectPasswordError is returned when an incorrect password is detected.

// IncorrectPasswordError is returned when
// an incorrect password is detected.
var IncorrectPasswordError = errors.New("x509: decryption password incorrect")

// CreateCertificate creates a new certificate based on a template. The following
// members of template are used: SerialNumber, Subject, NotBefore, NotAfter,
// KeyUsage, ExtKeyUsage, UnknownExtKeyUsage, BasicConstraintsValid, IsCA,
// MaxPathLen, SubjectKeyId, DNSNames, PermittedDNSDomainsCritical,
// PermittedDNSDomains, SignatureAlgorithm.
//
// The certificate is signed by parent. If parent is equal to template then the
// certificate is self-signed. The parameter pub is the public key of the signee
// and priv is the private key of the signer.
//
// The returned slice is the certificate in DER encoding.
//
// The only supported key types are RSA and ECDSA (*rsa.PublicKey or
// *ecdsa.PublicKey for pub, *rsa.PrivateKey or *ecdsa.PrivateKey for priv).

// CreateCertificate creates a new
// certificate based on a template. The
// following members of template are used:
// SerialNumber, Subject, NotBefore,
// NotAfter, KeyUsage, ExtKeyUsage,
// UnknownExtKeyUsage,
// BasicConstraintsValid, IsCA, MaxPathLen,
// SubjectKeyId, DNSNames,
// PermittedDNSDomainsCritical,
// PermittedDNSDomains, SignatureAlgorithm.
//
// The certificate is signed by parent. If
// parent is equal to template then the
// certificate is self-signed. The
// parameter pub is the public key of the
// signee and priv is the private key of
// the signer.
//
// The returned slice is the certificate in
// DER encoding.
//
// The only supported key types are RSA and
// ECDSA (*rsa.PublicKey or
// *ecdsa.PublicKey for pub,
// *rsa.PrivateKey or *ecdsa.PrivateKey for
// priv).
func CreateCertificate(rand io.Reader, template, parent *Certificate, pub interface{}, priv interface{}) (cert []byte, err error)

// CreateCertificateRequest creates a new certificate based on a template. The
// following members of template are used: Subject, Attributes, SignatureAlgorithm,
// Extensions, DNSNames, EmailAddresses, and IPAddresses. The private key is the
// private key of the signer.
//
// The returned slice is the certificate request in DER encoding.
//
// The only supported key types are RSA (*rsa.PrivateKey) and ECDSA
// (*ecdsa.PrivateKey).

// CreateCertificateRequest creates a new
// certificate based on a template. The
// following members of template are used:
// Subject, Attributes, SignatureAlgorithm,
// Extensions, DNSNames, EmailAddresses,
// and IPAddresses. The private key is the
// private key of the signer.
//
// The returned slice is the certificate
// request in DER encoding.
//
// The only supported key types are RSA
// (*rsa.PrivateKey) and ECDSA
// (*ecdsa.PrivateKey).
func CreateCertificateRequest(rand io.Reader, template *CertificateRequest, priv interface{}) (csr []byte, err error)

// DecryptPEMBlock takes a password encrypted PEM block and the password used to
// encrypt it and returns a slice of decrypted DER encoded bytes. It inspects the
// DEK-Info header to determine the algorithm used for decryption. If no DEK-Info
// header is present, an error is returned. If an incorrect password is detected an
// IncorrectPasswordError is returned.

// DecryptPEMBlock takes a password
// encrypted PEM block and the password
// used to encrypt it and returns a slice
// of decrypted DER encoded bytes. It
// inspects the DEK-Info header to
// determine the algorithm used for
// decryption. If no DEK-Info header is
// present, an error is returned. If an
// incorrect password is detected an
// IncorrectPasswordError is returned.
func DecryptPEMBlock(b *pem.Block, password []byte) ([]byte, error)

// EncryptPEMBlock returns a PEM block of the specified type holding the given
// DER-encoded data encrypted with the specified algorithm and password.

// EncryptPEMBlock returns a PEM block of
// the specified type holding the given
// DER-encoded data encrypted with the
// specified algorithm and password.
func EncryptPEMBlock(rand io.Reader, blockType string, data, password []byte, alg PEMCipher) (*pem.Block, error)

// IsEncryptedPEMBlock returns if the PEM block is password encrypted.

// IsEncryptedPEMBlock returns if the PEM
// block is password encrypted.
func IsEncryptedPEMBlock(b *pem.Block) bool

// MarshalECPrivateKey marshals an EC private key into ASN.1, DER format.

// MarshalECPrivateKey marshals an EC
// private key into ASN.1, DER format.
func MarshalECPrivateKey(key *ecdsa.PrivateKey) ([]byte, error)

// MarshalPKCS1PrivateKey converts a private key to ASN.1 DER encoded form.

// MarshalPKCS1PrivateKey converts a
// private key to ASN.1 DER encoded form.
func MarshalPKCS1PrivateKey(key *rsa.PrivateKey) []byte

// MarshalPKIXPublicKey serialises a public key to DER-encoded PKIX format.

// MarshalPKIXPublicKey serialises a public
// key to DER-encoded PKIX format.
func MarshalPKIXPublicKey(pub interface{}) ([]byte, error)

// ParseCRL parses a CRL from the given bytes. It's often the case that PEM encoded
// CRLs will appear where they should be DER encoded, so this function will
// transparently handle PEM encoding as long as there isn't any leading garbage.

// ParseCRL parses a CRL from the given
// bytes. It's often the case that PEM
// encoded CRLs will appear where they
// should be DER encoded, so this function
// will transparently handle PEM encoding
// as long as there isn't any leading
// garbage.
func ParseCRL(crlBytes []byte) (certList *pkix.CertificateList, err error)

// ParseCertificates parses one or more certificates from the given ASN.1 DER data.
// The certificates must be concatenated with no intermediate padding.

// ParseCertificates parses one or more
// certificates from the given ASN.1 DER
// data. The certificates must be
// concatenated with no intermediate
// padding.
func ParseCertificates(asn1Data []byte) ([]*Certificate, error)

// ParseDERCRL parses a DER encoded CRL from the given bytes.

// ParseDERCRL parses a DER encoded CRL
// from the given bytes.
func ParseDERCRL(derBytes []byte) (certList *pkix.CertificateList, err error)

// ParseECPrivateKey parses an ASN.1 Elliptic Curve Private Key Structure.

// ParseECPrivateKey parses an ASN.1
// Elliptic Curve Private Key Structure.
func ParseECPrivateKey(der []byte) (key *ecdsa.PrivateKey, err error)

// ParsePKCS1PrivateKey returns an RSA private key from its ASN.1 PKCS#1 DER
// encoded form.

// ParsePKCS1PrivateKey returns an RSA
// private key from its ASN.1 PKCS#1 DER
// encoded form.
func ParsePKCS1PrivateKey(der []byte) (key *rsa.PrivateKey, err error)

// ParsePKCS8PrivateKey parses an unencrypted, PKCS#8 private key. See
// http://www.rsa.com/rsalabs/node.asp?id=2130 and RFC5208.

// ParsePKCS8PrivateKey parses an
// unencrypted, PKCS#8 private key. See
// http://www.rsa.com/rsalabs/node.asp?id=2130
// and RFC5208.
func ParsePKCS8PrivateKey(der []byte) (key interface{}, err error)

// ParsePKIXPublicKey parses a DER encoded public key. These values are typically
// found in PEM blocks with "BEGIN PUBLIC KEY".

// ParsePKIXPublicKey parses a DER encoded
// public key. These values are typically
// found in PEM blocks with "BEGIN PUBLIC
// KEY".
func ParsePKIXPublicKey(derBytes []byte) (pub interface{}, err error)

// CertPool is a set of certificates.

// CertPool is a set of certificates.
type CertPool struct {
	// contains filtered or unexported fields
}

// NewCertPool returns a new, empty CertPool.

// NewCertPool returns a new, empty
// CertPool.
func NewCertPool() *CertPool

// AddCert adds a certificate to a pool.

// AddCert adds a certificate to a pool.
func (s *CertPool) AddCert(cert *Certificate)

// AppendCertsFromPEM attempts to parse a series of PEM encoded certificates. It
// appends any certificates found to s and returns true if any certificates were
// successfully parsed.
//
// On many Linux systems, /etc/ssl/cert.pem will contain the system wide set of
// root CAs in a format suitable for this function.

// AppendCertsFromPEM attempts to parse a
// series of PEM encoded certificates. It
// appends any certificates found to s and
// returns true if any certificates were
// successfully parsed.
//
// On many Linux systems, /etc/ssl/cert.pem
// will contain the system wide set of root
// CAs in a format suitable for this
// function.
func (s *CertPool) AppendCertsFromPEM(pemCerts []byte) (ok bool)

// Subjects returns a list of the DER-encoded subjects of all of the certificates
// in the pool.

// Subjects returns a list of the
// DER-encoded subjects of all of the
// certificates in the pool.
func (s *CertPool) Subjects() (res [][]byte)

// A Certificate represents an X.509 certificate.

// A Certificate represents an X.509
// certificate.
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

// ParseCertificate parses a single certificate from the given ASN.1 DER data.

// ParseCertificate parses a single
// certificate from the given ASN.1 DER
// data.
func ParseCertificate(asn1Data []byte) (*Certificate, error)

// CheckCRLSignature checks that the signature in crl is from c.

// CheckCRLSignature checks that the
// signature in crl is from c.
func (c *Certificate) CheckCRLSignature(crl *pkix.CertificateList) (err error)

// CheckSignature verifies that signature is a valid signature over signed from c's
// public key.

// CheckSignature verifies that signature
// is a valid signature over signed from
// c's public key.
func (c *Certificate) CheckSignature(algo SignatureAlgorithm, signed, signature []byte) (err error)

// CheckSignatureFrom verifies that the signature on c is a valid signature from
// parent.

// CheckSignatureFrom verifies that the
// signature on c is a valid signature from
// parent.
func (c *Certificate) CheckSignatureFrom(parent *Certificate) (err error)

// CreateCRL returns a DER encoded CRL, signed by this Certificate, that contains
// the given list of revoked certificates.
//
// The only supported key type is RSA (*rsa.PrivateKey for priv).

// CreateCRL returns a DER encoded CRL,
// signed by this Certificate, that
// contains the given list of revoked
// certificates.
//
// The only supported key type is RSA
// (*rsa.PrivateKey for priv).
func (c *Certificate) CreateCRL(rand io.Reader, priv interface{}, revokedCerts []pkix.RevokedCertificate, now, expiry time.Time) (crlBytes []byte, err error)

func (c *Certificate) Equal(other *Certificate) bool

// Verify attempts to verify c by building one or more chains from c to a
// certificate in opts.Roots, using certificates in opts.Intermediates if needed.
// If successful, it returns one or more chains where the first element of the
// chain is c and the last element is from opts.Roots.
//
// If opts.Roots is nil and system roots are unavailable the returned error will be
// of type SystemRootsError.
//
// WARNING: this doesn't do any revocation checking.

// Verify attempts to verify c by building
// one or more chains from c to a
// certificate in opts.Roots, using
// certificates in opts.Intermediates if
// needed. If successful, it returns one or
// more chains where the first element of
// the chain is c and the last element is
// from opts.Roots.
//
// If opts.Roots is nil and system roots
// are unavailable the returned error will
// be of type SystemRootsError.
//
// WARNING: this doesn't do any revocation
// checking.
func (c *Certificate) Verify(opts VerifyOptions) (chains [][]*Certificate, err error)

// VerifyHostname returns nil if c is a valid certificate for the named host.
// Otherwise it returns an error describing the mismatch.

// VerifyHostname returns nil if c is a
// valid certificate for the named host.
// Otherwise it returns an error describing
// the mismatch.
func (c *Certificate) VerifyHostname(h string) error

// CertificateInvalidError results when an odd error occurs. Users of this library
// probably want to handle all these errors uniformly.

// CertificateInvalidError results when an
// odd error occurs. Users of this library
// probably want to handle all these errors
// uniformly.
type CertificateInvalidError struct {
	Cert   *Certificate
	Reason InvalidReason
}

func (e CertificateInvalidError) Error() string

// CertificateRequest represents a PKCS #10, certificate signature request.

// CertificateRequest represents a PKCS
// #10, certificate signature request.
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

// ParseCertificateRequest parses a single certificate request from the given ASN.1
// DER data.

// ParseCertificateRequest parses a single
// certificate request from the given ASN.1
// DER data.
func ParseCertificateRequest(asn1Data []byte) (*CertificateRequest, error)

// ConstraintViolationError results when a requested usage is not permitted by a
// certificate. For example: checking a signature when the public key isn't a
// certificate signing key.

// ConstraintViolationError results when a
// requested usage is not permitted by a
// certificate. For example: checking a
// signature when the public key isn't a
// certificate signing key.
type ConstraintViolationError struct{}

func (ConstraintViolationError) Error() string

// ExtKeyUsage represents an extended set of actions that are valid for a given
// key. Each of the ExtKeyUsage* constants define a unique action.

// ExtKeyUsage represents an extended set
// of actions that are valid for a given
// key. Each of the ExtKeyUsage* constants
// define a unique action.
type ExtKeyUsage int

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

// HostnameError results when the set of authorized names doesn't match the
// requested name.

// HostnameError results when the set of
// authorized names doesn't match the
// requested name.
type HostnameError struct {
	Certificate *Certificate
	Host        string
}

func (h HostnameError) Error() string

type InvalidReason int

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

// KeyUsage represents the set of actions that are valid for a given key. It's a
// bitmap of the KeyUsage* constants.

// KeyUsage represents the set of actions
// that are valid for a given key. It's a
// bitmap of the KeyUsage* constants.
type KeyUsage int

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

type PEMCipher int

// Possible values for the EncryptPEMBlock encryption algorithm.

// Possible values for the EncryptPEMBlock
// encryption algorithm.
const (
	_ PEMCipher = iota
	PEMCipherDES
	PEMCipher3DES
	PEMCipherAES128
	PEMCipherAES192
	PEMCipherAES256
)

type PublicKeyAlgorithm int

const (
	UnknownPublicKeyAlgorithm PublicKeyAlgorithm = iota
	RSA
	DSA
	ECDSA
)

type SignatureAlgorithm int

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

// SystemRootsError results when we fail to load the system root certificates.

// SystemRootsError results when we fail to
// load the system root certificates.
type SystemRootsError struct{}

func (SystemRootsError) Error() string

type UnhandledCriticalExtension struct{}

func (h UnhandledCriticalExtension) Error() string

// UnknownAuthorityError results when the certificate issuer is unknown

// UnknownAuthorityError results when the
// certificate issuer is unknown
type UnknownAuthorityError struct {
	// contains filtered or unexported fields
}

func (e UnknownAuthorityError) Error() string

// VerifyOptions contains parameters for Certificate.Verify. It's a structure
// because other PKIX verification APIs have ended up needing many options.

// VerifyOptions contains parameters for
// Certificate.Verify. It's a structure
// because other PKIX verification APIs
// have ended up needing many options.
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
