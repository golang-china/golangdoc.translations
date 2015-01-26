// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package pkix contains shared, low level structures used for ASN.1 parsing and
// serialization of X.509 certificates, CRL and OCSP.

// Package pkix contains shared, low level
// structures used for ASN.1 parsing and
// serialization of X.509 certificates, CRL
// and OCSP.
package pkix

// AlgorithmIdentifier represents the ASN.1 structure of the same name. See RFC
// 5280, section 4.1.1.2.

// AlgorithmIdentifier represents the ASN.1
// structure of the same name. See RFC
// 5280, section 4.1.1.2.
type AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}

// AttributeTypeAndValue mirrors the ASN.1 structure of the same name in
// http://tools.ietf.org/html/rfc5280#section-4.1.2.4

// AttributeTypeAndValue mirrors the ASN.1
// structure of the same name in
// http://tools.ietf.org/html/rfc5280#section-4.1.2.4
type AttributeTypeAndValue struct {
	Type  asn1.ObjectIdentifier
	Value interface{}
}

// AttributeTypeAndValueSET represents a set of ASN.1 sequences of
// AttributeTypeAndValue sequences from RFC 2986 (PKCS #10).

// AttributeTypeAndValueSET represents a
// set of ASN.1 sequences of
// AttributeTypeAndValue sequences from RFC
// 2986 (PKCS #10).
type AttributeTypeAndValueSET struct {
	Type  asn1.ObjectIdentifier
	Value [][]AttributeTypeAndValue `asn1:"set"`
}

// CertificateList represents the ASN.1 structure of the same name. See RFC 5280,
// section 5.1. Use Certificate.CheckCRLSignature to verify the signature.

// CertificateList represents the ASN.1
// structure of the same name. See RFC
// 5280, section 5.1. Use
// Certificate.CheckCRLSignature to verify
// the signature.
type CertificateList struct {
	TBSCertList        TBSCertificateList
	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

// HasExpired reports whether now is past the expiry time of certList.

// HasExpired reports whether now is past
// the expiry time of certList.
func (certList *CertificateList) HasExpired(now time.Time) bool

// Extension represents the ASN.1 structure of the same name. See RFC 5280, section
// 4.2.

// Extension represents the ASN.1 structure
// of the same name. See RFC 5280, section
// 4.2.
type Extension struct {
	Id       asn1.ObjectIdentifier
	Critical bool `asn1:"optional"`
	Value    []byte
}

// Name represents an X.509 distinguished name. This only includes the common
// elements of a DN. Additional elements in the name are ignored.

// Name represents an X.509 distinguished
// name. This only includes the common
// elements of a DN. Additional elements in
// the name are ignored.
type Name struct {
	Country, Organization, OrganizationalUnit []string
	Locality, Province                        []string
	StreetAddress, PostalCode                 []string
	SerialNumber, CommonName                  string

	Names []AttributeTypeAndValue
}

func (n *Name) FillFromRDNSequence(rdns *RDNSequence)

func (n Name) ToRDNSequence() (ret RDNSequence)

type RDNSequence []RelativeDistinguishedNameSET

type RelativeDistinguishedNameSET []AttributeTypeAndValue

// RevokedCertificate represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1.

// RevokedCertificate represents the ASN.1
// structure of the same name. See RFC
// 5280, section 5.1.
type RevokedCertificate struct {
	SerialNumber   *big.Int
	RevocationTime time.Time
	Extensions     []Extension `asn1:"optional"`
}

// TBSCertificateList represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1.

// TBSCertificateList represents the ASN.1
// structure of the same name. See RFC
// 5280, section 5.1.
type TBSCertificateList struct {
	Raw                 asn1.RawContent
	Version             int `asn1:"optional,default:2"`
	Signature           AlgorithmIdentifier
	Issuer              RDNSequence
	ThisUpdate          time.Time
	NextUpdate          time.Time            `asn1:"optional"`
	RevokedCertificates []RevokedCertificate `asn1:"optional"`
	Extensions          []Extension          `asn1:"tag:0,optional,explicit"`
}
