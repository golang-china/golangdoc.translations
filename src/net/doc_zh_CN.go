// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package net provides a portable interface for network I/O, including TCP/IP,
// UDP, domain name resolution, and Unix domain sockets.
//
// Although the package provides access to low-level networking primitives, most
// clients will need only the basic interface provided by the Dial, Listen, and
// Accept functions and the associated Conn and Listener interfaces. The crypto/tls
// package uses the same interfaces and similar Dial and Listen functions.
//
// The Dial function connects to a server:
//
//	conn, err := net.Dial("tcp", "google.com:80")
//	if err != nil {
//		// handle error
//	}
//	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
//	status, err := bufio.NewReader(conn).ReadString('\n')
//	// ...
//
// The Listen function creates servers:
//
//	ln, err := net.Listen("tcp", ":8080")
//	if err != nil {
//		// handle error
//	}
//	for {
//		conn, err := ln.Accept()
//		if err != nil {
//			// handle error
//		}
//		go handleConnection(conn)
//	}
package net

// IP address lengths (bytes).
const (
	IPv4len = 4
	IPv6len = 16
)

// Well-known IPv4 addresses
var (
	IPv4bcast     = IPv4(255, 255, 255, 255) // broadcast
	IPv4allsys    = IPv4(224, 0, 0, 1)       // all systems
	IPv4allrouter = IPv4(224, 0, 0, 2)       // all routers
	IPv4zero      = IPv4(0, 0, 0, 0)         // all zeros
)

// Well-known IPv6 addresses
var (
	IPv6zero                   = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	IPv6unspecified            = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	IPv6loopback               = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	IPv6interfacelocalallnodes = IP{0xff, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
	IPv6linklocalallnodes      = IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
	IPv6linklocalallrouters    = IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02}
)

// Various errors contained in OpError.
var (
	ErrWriteToConnected = errors.New("use of WriteTo with pre-connected connection")
)

// InterfaceAddrs returns a list of the system's network interface addresses.
func InterfaceAddrs() ([]Addr, error)

// Interfaces returns a list of the system's network interfaces.
func Interfaces() ([]Interface, error)

// JoinHostPort combines host and port into a network address of the form
// "host:port" or, if host contains a colon or a percent sign, "[host]:port".
func JoinHostPort(host, port string) string

// LookupAddr performs a reverse lookup for the given address, returning a list of
// names mapping to that address.
func LookupAddr(addr string) (name []string, err error)

// LookupCNAME returns the canonical DNS host for the given name. Callers that do
// not care about the canonical name can call LookupHost or LookupIP directly; both
// take care of resolving the canonical name as part of the lookup.
func LookupCNAME(name string) (cname string, err error)

// LookupHost looks up the given host using the local resolver. It returns an array
// of that host's addresses.
func LookupHost(host string) (addrs []string, err error)

// LookupIP looks up host using the local resolver. It returns an array of that
// host's IPv4 and IPv6 addresses.
func LookupIP(host string) (addrs []IP, err error)

// LookupMX returns the DNS MX records for the given domain name sorted by
// preference.
func LookupMX(name string) (mx []*MX, err error)

// LookupNS returns the DNS NS records for the given domain name.
func LookupNS(name string) (ns []*NS, err error)

// LookupPort looks up the port for the given network and service.
func LookupPort(network, service string) (port int, err error)

// LookupSRV tries to resolve an SRV query of the given service, protocol, and
// domain name. The proto is "tcp" or "udp". The returned records are sorted by
// priority and randomized by weight within a priority.
//
// LookupSRV constructs the DNS name to look up following RFC 2782. That is, it
// looks up _service._proto.name. To accommodate services publishing SRV records
// under non-standard names, if both service and proto are empty strings, LookupSRV
// looks up name directly.
func LookupSRV(service, proto, name string) (cname string, addrs []*SRV, err error)

// LookupTXT returns the DNS TXT records for the given domain name.
func LookupTXT(name string) (txt []string, err error)

// SplitHostPort splits a network address of the form "host:port", "[host]:port" or
// "[ipv6-host%zone]:port" into host or ipv6-host%zone and port. A literal address
// or host name for IPv6 must be enclosed in square brackets, as in "[::1]:80",
// "[ipv6-host]:http" or "[ipv6-host%zone]:80".
func SplitHostPort(hostport string) (host, port string, err error)

// Addr represents a network end point address.
type Addr interface {
	Network() string // name of the network
	String() string  // string form of address
}

type AddrError struct {
	Err  string
	Addr string
}

func (e *AddrError) Error() string

func (e *AddrError) Temporary() bool

func (e *AddrError) Timeout() bool

// Conn is a generic stream-oriented network connection.
//
// Multiple goroutines may invoke methods on a Conn simultaneously.
type Conn interface {
	// Read reads data from the connection.
	// Read can be made to time out and return a Error with Timeout() == true
	// after a fixed time limit; see SetDeadline and SetReadDeadline.
	Read(b []byte) (n int, err error)

	// Write writes data to the connection.
	// Write can be made to time out and return a Error with Timeout() == true
	// after a fixed time limit; see SetDeadline and SetWriteDeadline.
	Write(b []byte) (n int, err error)

	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close() error

	// LocalAddr returns the local network address.
	LocalAddr() Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr() Addr

	// SetDeadline sets the read and write deadlines associated
	// with the connection. It is equivalent to calling both
	// SetReadDeadline and SetWriteDeadline.
	//
	// A deadline is an absolute time after which I/O operations
	// fail with a timeout (see type Error) instead of
	// blocking. The deadline applies to all future I/O, not just
	// the immediately following call to Read or Write.
	//
	// An idle timeout can be implemented by repeatedly extending
	// the deadline after successful Read or Write calls.
	//
	// A zero value for t means I/O operations will not time out.
	SetDeadline(t time.Time) error

	// SetReadDeadline sets the deadline for future Read calls.
	// A zero value for t means Read will not time out.
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets the deadline for future Write calls.
	// Even if write times out, it may return n > 0, indicating that
	// some of the data was successfully written.
	// A zero value for t means Write will not time out.
	SetWriteDeadline(t time.Time) error
}

// Dial connects to the address on the named network.
//
// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4"
// (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only),
// "unix", "unixgram" and "unixpacket".
//
// For TCP and UDP networks, addresses have the form host:port. If host is a
// literal IPv6 address it must be enclosed in square brackets as in "[::1]:80" or
// "[ipv6-host%zone]:80". The functions JoinHostPort and SplitHostPort manipulate
// addresses in this form.
//
// Examples:
//
//	Dial("tcp", "12.34.56.78:80")
//	Dial("tcp", "google.com:http")
//	Dial("tcp", "[2001:db8::1]:http")
//	Dial("tcp", "[fe80::1%lo0]:80")
//
// For IP networks, the network must be "ip", "ip4" or "ip6" followed by a colon
// and a protocol number or name and the addr must be a literal IP address.
//
// Examples:
//
//	Dial("ip4:1", "127.0.0.1")
//	Dial("ip6:ospf", "::1")
//
// For Unix networks, the address must be a file system path.
func Dial(network, address string) (Conn, error)

// DialTimeout acts like Dial but takes a timeout. The timeout includes name
// resolution, if required.
func DialTimeout(network, address string, timeout time.Duration) (Conn, error)

// FileConn returns a copy of the network connection corresponding to the open file
// f. It is the caller's responsibility to close f when finished. Closing c does
// not affect f, and closing f does not affect c.
func FileConn(f *os.File) (c Conn, err error)

// Pipe creates a synchronous, in-memory, full duplex network connection; both ends
// implement the Conn interface. Reads on one end are matched with writes on the
// other, copying data directly between the two; there is no internal buffering.
func Pipe() (Conn, Conn)

// DNSConfigError represents an error reading the machine's DNS configuration.
type DNSConfigError struct {
	Err error
}

func (e *DNSConfigError) Error() string

func (e *DNSConfigError) Temporary() bool

func (e *DNSConfigError) Timeout() bool

// DNSError represents a DNS lookup error.
type DNSError struct {
	Err       string // description of the error
	Name      string // name looked for
	Server    string // server used
	IsTimeout bool
}

func (e *DNSError) Error() string

func (e *DNSError) Temporary() bool

func (e *DNSError) Timeout() bool

// A Dialer contains options for connecting to an address.
//
// The zero value for each field is equivalent to dialing without that option.
// Dialing with the zero value of Dialer is therefore equivalent to just calling
// the Dial function.
type Dialer struct {
	// Timeout is the maximum amount of time a dial will wait for
	// a connect to complete. If Deadline is also set, it may fail
	// earlier.
	//
	// The default is no timeout.
	//
	// With or without a timeout, the operating system may impose
	// its own earlier timeout. For instance, TCP timeouts are
	// often around 3 minutes.
	Timeout time.Duration

	// Deadline is the absolute point in time after which dials
	// will fail. If Timeout is set, it may fail earlier.
	// Zero means no deadline, or dependent on the operating system
	// as with the Timeout option.
	Deadline time.Time

	// LocalAddr is the local address to use when dialing an
	// address. The address must be of a compatible type for the
	// network being dialed.
	// If nil, a local address is automatically chosen.
	LocalAddr Addr

	// DualStack allows a single dial to attempt to establish
	// multiple IPv4 and IPv6 connections and to return the first
	// established connection when the network is "tcp" and the
	// destination is a host name that has multiple address family
	// DNS records.
	DualStack bool

	// KeepAlive specifies the keep-alive period for an active
	// network connection.
	// If zero, keep-alives are not enabled. Network protocols
	// that do not support keep-alives ignore this field.
	KeepAlive time.Duration
}

// Dial connects to the address on the named network.
//
// See func Dial for a description of the network and address parameters.
func (d *Dialer) Dial(network, address string) (Conn, error)

// An Error represents a network error.
type Error interface {
	error
	Timeout() bool   // Is the error a timeout?
	Temporary() bool // Is the error temporary?
}

type Flags uint

const (
	FlagUp           Flags = 1 << iota // interface is up
	FlagBroadcast                      // interface supports broadcast access capability
	FlagLoopback                       // interface is a loopback interface
	FlagPointToPoint                   // interface belongs to a point-to-point link
	FlagMulticast                      // interface supports multicast access capability
)

func (f Flags) String() string

// A HardwareAddr represents a physical hardware address.
type HardwareAddr []byte

// ParseMAC parses s as an IEEE 802 MAC-48, EUI-48, or EUI-64 using one of the
// following formats:
//
//	01:23:45:67:89:ab
//	01:23:45:67:89:ab:cd:ef
//	01-23-45-67-89-ab
//	01-23-45-67-89-ab-cd-ef
//	0123.4567.89ab
//	0123.4567.89ab.cdef
func ParseMAC(s string) (hw HardwareAddr, err error)

func (a HardwareAddr) String() string

// An IP is a single IP address, a slice of bytes. Functions in this package accept
// either 4-byte (IPv4) or 16-byte (IPv6) slices as input.
//
// Note that in this documentation, referring to an IP address as an IPv4 address
// or an IPv6 address is a semantic property of the address, not just the length of
// the byte slice: a 16-byte slice can still be an IPv4 address.
type IP []byte

// IPv4 returns the IP address (in 16-byte form) of the IPv4 address a.b.c.d.
func IPv4(a, b, c, d byte) IP

// ParseCIDR parses s as a CIDR notation IP address and mask, like
// "192.168.100.1/24" or "2001:DB8::/48", as defined in RFC 4632 and RFC 4291.
//
// It returns the IP address and the network implied by the IP and mask. For
// example, ParseCIDR("192.168.100.1/16") returns the IP address 192.168.100.1 and
// the network 192.168.0.0/16.
func ParseCIDR(s string) (IP, *IPNet, error)

// ParseIP parses s as an IP address, returning the result. The string s can be in
// dotted decimal ("74.125.19.99") or IPv6 ("2001:4860:0:2001::68") form. If s is
// not a valid textual representation of an IP address, ParseIP returns nil.
func ParseIP(s string) IP

// DefaultMask returns the default IP mask for the IP address ip. Only IPv4
// addresses have default masks; DefaultMask returns nil if ip is not a valid IPv4
// address.
func (ip IP) DefaultMask() IPMask

// Equal returns true if ip and x are the same IP address. An IPv4 address and that
// same address in IPv6 form are considered to be equal.
func (ip IP) Equal(x IP) bool

// IsGlobalUnicast returns true if ip is a global unicast address.
func (ip IP) IsGlobalUnicast() bool

// IsInterfaceLinkLocalMulticast returns true if ip is an interface-local multicast
// address.
func (ip IP) IsInterfaceLocalMulticast() bool

// IsLinkLocalMulticast returns true if ip is a link-local multicast address.
func (ip IP) IsLinkLocalMulticast() bool

// IsLinkLocalUnicast returns true if ip is a link-local unicast address.
func (ip IP) IsLinkLocalUnicast() bool

// IsLoopback returns true if ip is a loopback address.
func (ip IP) IsLoopback() bool

// IsMulticast returns true if ip is a multicast address.
func (ip IP) IsMulticast() bool

// IsUnspecified returns true if ip is an unspecified address.
func (ip IP) IsUnspecified() bool

// MarshalText implements the encoding.TextMarshaler interface. The encoding is the
// same as returned by String.
func (ip IP) MarshalText() ([]byte, error)

// Mask returns the result of masking the IP address ip with mask.
func (ip IP) Mask(mask IPMask) IP

// String returns the string form of the IP address ip. If the address is an IPv4
// address, the string representation is dotted decimal ("74.125.19.99"). Otherwise
// the representation is IPv6 ("2001:4860:0:2001::68").
func (ip IP) String() string

// To16 converts the IP address ip to a 16-byte representation. If ip is not an IP
// address (it is the wrong length), To16 returns nil.
func (ip IP) To16() IP

// To4 converts the IPv4 address ip to a 4-byte representation. If ip is not an
// IPv4 address, To4 returns nil.
func (ip IP) To4() IP

// UnmarshalText implements the encoding.TextUnmarshaler interface. The IP address
// is expected in a form accepted by ParseIP.
func (ip *IP) UnmarshalText(text []byte) error

// IPAddr represents the address of an IP end point.
type IPAddr struct {
	IP   IP
	Zone string // IPv6 scoped addressing zone
}

// ResolveIPAddr parses addr as an IP address of the form "host" or
// "ipv6-host%zone" and resolves the domain name on the network net, which must be
// "ip", "ip4" or "ip6".
func ResolveIPAddr(net, addr string) (*IPAddr, error)

// Network returns the address's network name, "ip".
func (a *IPAddr) Network() string

func (a *IPAddr) String() string

// IPConn is the implementation of the Conn and PacketConn interfaces for IP
// network connections.
type IPConn struct {
	// contains filtered or unexported fields
}

// DialIP connects to the remote address raddr on the network protocol netProto,
// which must be "ip", "ip4", or "ip6" followed by a colon and a protocol number or
// name.
func DialIP(netProto string, laddr, raddr *IPAddr) (*IPConn, error)

// ListenIP listens for incoming IP packets addressed to the local address laddr.
// The returned connection's ReadFrom and WriteTo methods can be used to receive
// and send IP packets with per-packet addressing.
func ListenIP(netProto string, laddr *IPAddr) (*IPConn, error)

// Close closes the connection.
func (c *IPConn) Close() error

// File sets the underlying os.File to blocking mode and returns a copy. It is the
// caller's responsibility to close f when finished. Closing c does not affect f,
// and closing f does not affect c.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.
func (c *IPConn) File() (f *os.File, err error)

// LocalAddr returns the local network address.
func (c *IPConn) LocalAddr() Addr

// Read implements the Conn Read method.
func (c *IPConn) Read(b []byte) (int, error)

// ReadFrom implements the PacketConn ReadFrom method.
func (c *IPConn) ReadFrom(b []byte) (int, Addr, error)

// ReadFromIP reads an IP packet from c, copying the payload into b. It returns the
// number of bytes copied into b and the return address that was on the packet.
//
// ReadFromIP can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *IPConn) ReadFromIP(b []byte) (int, *IPAddr, error)

// ReadMsgIP reads a packet from c, copying the payload into b and the associated
// out-of-band data into oob. It returns the number of bytes copied into b, the
// number of bytes copied into oob, the flags that were set on the packet and the
// source address of the packet.
func (c *IPConn) ReadMsgIP(b, oob []byte) (n, oobn, flags int, addr *IPAddr, err error)

// RemoteAddr returns the remote network address.
func (c *IPConn) RemoteAddr() Addr

// SetDeadline implements the Conn SetDeadline method.
func (c *IPConn) SetDeadline(t time.Time) error

// SetReadBuffer sets the size of the operating system's receive buffer associated
// with the connection.
func (c *IPConn) SetReadBuffer(bytes int) error

// SetReadDeadline implements the Conn SetReadDeadline method.
func (c *IPConn) SetReadDeadline(t time.Time) error

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the connection.
func (c *IPConn) SetWriteBuffer(bytes int) error

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (c *IPConn) SetWriteDeadline(t time.Time) error

// Write implements the Conn Write method.
func (c *IPConn) Write(b []byte) (int, error)

// WriteMsgIP writes a packet to addr via c, copying the payload from b and the
// associated out-of-band data from oob. It returns the number of payload and
// out-of-band bytes written.
func (c *IPConn) WriteMsgIP(b, oob []byte, addr *IPAddr) (n, oobn int, err error)

// WriteTo implements the PacketConn WriteTo method.
func (c *IPConn) WriteTo(b []byte, addr Addr) (int, error)

// WriteToIP writes an IP packet to addr via c, copying the payload from b.
//
// WriteToIP can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline. On
// packet-oriented connections, write timeouts are rare.
func (c *IPConn) WriteToIP(b []byte, addr *IPAddr) (int, error)

// An IP mask is an IP address.
type IPMask []byte

// CIDRMask returns an IPMask consisting of `ones' 1 bits followed by 0s up to a
// total length of `bits' bits. For a mask of this form, CIDRMask is the inverse of
// IPMask.Size.
func CIDRMask(ones, bits int) IPMask

// IPv4Mask returns the IP mask (in 4-byte form) of the IPv4 mask a.b.c.d.
func IPv4Mask(a, b, c, d byte) IPMask

// Size returns the number of leading ones and total bits in the mask. If the mask
// is not in the canonical form--ones followed by zeros--then Size returns 0, 0.
func (m IPMask) Size() (ones, bits int)

// String returns the hexadecimal form of m, with no punctuation.
func (m IPMask) String() string

// An IPNet represents an IP network.
type IPNet struct {
	IP   IP     // network number
	Mask IPMask // network mask
}

// Contains reports whether the network includes ip.
func (n *IPNet) Contains(ip IP) bool

// Network returns the address's network name, "ip+net".
func (n *IPNet) Network() string

// String returns the CIDR notation of n like "192.168.100.1/24" or "2001:DB8::/48"
// as defined in RFC 4632 and RFC 4291. If the mask is not in the canonical form,
// it returns the string which consists of an IP address, followed by a slash
// character and a mask expressed as hexadecimal form with no punctuation like
// "192.168.100.1/c000ff00".
func (n *IPNet) String() string

// Interface represents a mapping between network interface name and index. It also
// represents network interface facility information.
type Interface struct {
	Index        int          // positive integer that starts at one, zero is never used
	MTU          int          // maximum transmission unit
	Name         string       // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr HardwareAddr // IEEE MAC-48, EUI-48 and EUI-64 form
	Flags        Flags        // e.g., FlagUp, FlagLoopback, FlagMulticast
}

// InterfaceByIndex returns the interface specified by index.
func InterfaceByIndex(index int) (*Interface, error)

// InterfaceByName returns the interface specified by name.
func InterfaceByName(name string) (*Interface, error)

// Addrs returns interface addresses for a specific interface.
func (ifi *Interface) Addrs() ([]Addr, error)

// MulticastAddrs returns multicast, joined group addresses for a specific
// interface.
func (ifi *Interface) MulticastAddrs() ([]Addr, error)

type InvalidAddrError string

func (e InvalidAddrError) Error() string

func (e InvalidAddrError) Temporary() bool

func (e InvalidAddrError) Timeout() bool

// A Listener is a generic network listener for stream-oriented protocols.
//
// Multiple goroutines may invoke methods on a Listener simultaneously.
type Listener interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (c Conn, err error)

	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error

	// Addr returns the listener's network address.
	Addr() Addr
}

// FileListener returns a copy of the network listener corresponding to the open
// file f. It is the caller's responsibility to close l when finished. Closing l
// does not affect f, and closing f does not affect l.
func FileListener(f *os.File) (l Listener, err error)

// Listen announces on the local network address laddr. The network net must be a
// stream-oriented network: "tcp", "tcp4", "tcp6", "unix" or "unixpacket". See Dial
// for the syntax of laddr.
func Listen(net, laddr string) (Listener, error)

// An MX represents a single DNS MX record.
type MX struct {
	Host string
	Pref uint16
}

// An NS represents a single DNS NS record.
type NS struct {
	Host string
}

// OpError is the error type usually returned by functions in the net package. It
// describes the operation, network type, and address of an error.
type OpError struct {
	// Op is the operation which caused the error, such as
	// "read" or "write".
	Op string

	// Net is the network type on which this error occurred,
	// such as "tcp" or "udp6".
	Net string

	// Addr is the network address on which this error occurred.
	Addr Addr

	// Err is the error that occurred during the operation.
	Err error
}

func (e *OpError) Error() string

func (e *OpError) Temporary() bool

func (e *OpError) Timeout() bool

// PacketConn is a generic packet-oriented network connection.
//
// Multiple goroutines may invoke methods on a PacketConn simultaneously.
type PacketConn interface {
	// ReadFrom reads a packet from the connection,
	// copying the payload into b.  It returns the number of
	// bytes copied into b and the return address that
	// was on the packet.
	// ReadFrom can be made to time out and return
	// an error with Timeout() == true after a fixed time limit;
	// see SetDeadline and SetReadDeadline.
	ReadFrom(b []byte) (n int, addr Addr, err error)

	// WriteTo writes a packet with payload b to addr.
	// WriteTo can be made to time out and return
	// an error with Timeout() == true after a fixed time limit;
	// see SetDeadline and SetWriteDeadline.
	// On packet-oriented connections, write timeouts are rare.
	WriteTo(b []byte, addr Addr) (n int, err error)

	// Close closes the connection.
	// Any blocked ReadFrom or WriteTo operations will be unblocked and return errors.
	Close() error

	// LocalAddr returns the local network address.
	LocalAddr() Addr

	// SetDeadline sets the read and write deadlines associated
	// with the connection.
	SetDeadline(t time.Time) error

	// SetReadDeadline sets the deadline for future Read calls.
	// If the deadline is reached, Read will fail with a timeout
	// (see type Error) instead of blocking.
	// A zero value for t means Read will not time out.
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets the deadline for future Write calls.
	// If the deadline is reached, Write will fail with a timeout
	// (see type Error) instead of blocking.
	// A zero value for t means Write will not time out.
	// Even if write times out, it may return n > 0, indicating that
	// some of the data was successfully written.
	SetWriteDeadline(t time.Time) error
}

// FilePacketConn returns a copy of the packet network connection corresponding to
// the open file f. It is the caller's responsibility to close f when finished.
// Closing c does not affect f, and closing f does not affect c.
func FilePacketConn(f *os.File) (c PacketConn, err error)

// ListenPacket announces on the local network address laddr. The network net must
// be a packet-oriented network: "udp", "udp4", "udp6", "ip", "ip4", "ip6" or
// "unixgram". See Dial for the syntax of laddr.
func ListenPacket(net, laddr string) (PacketConn, error)

// A ParseError represents a malformed text string and the type of string that was
// expected.
type ParseError struct {
	Type string
	Text string
}

func (e *ParseError) Error() string

// An SRV represents a single DNS SRV record.
type SRV struct {
	Target   string
	Port     uint16
	Priority uint16
	Weight   uint16
}

// TCPAddr represents the address of a TCP end point.
type TCPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

// ResolveTCPAddr parses addr as a TCP address of the form "host:port" or
// "[ipv6-host%zone]:port" and resolves a pair of domain name and port name on the
// network net, which must be "tcp", "tcp4" or "tcp6". A literal address or host
// name for IPv6 must be enclosed in square brackets, as in "[::1]:80",
// "[ipv6-host]:http" or "[ipv6-host%zone]:80".
func ResolveTCPAddr(net, addr string) (*TCPAddr, error)

// Network returns the address's network name, "tcp".
func (a *TCPAddr) Network() string

func (a *TCPAddr) String() string

// TCPConn is an implementation of the Conn interface for TCP network connections.
type TCPConn struct {
	// contains filtered or unexported fields
}

// DialTCP connects to the remote address raddr on the network net, which must be
// "tcp", "tcp4", or "tcp6". If laddr is not nil, it is used as the local address
// for the connection.
func DialTCP(net string, laddr, raddr *TCPAddr) (*TCPConn, error)

// Close closes the connection.
func (c *TCPConn) Close() error

// CloseRead shuts down the reading side of the TCP connection. Most callers should
// just use Close.
func (c *TCPConn) CloseRead() error

// CloseWrite shuts down the writing side of the TCP connection. Most callers
// should just use Close.
func (c *TCPConn) CloseWrite() error

// File sets the underlying os.File to blocking mode and returns a copy. It is the
// caller's responsibility to close f when finished. Closing c does not affect f,
// and closing f does not affect c.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.
func (c *TCPConn) File() (f *os.File, err error)

// LocalAddr returns the local network address.
func (c *TCPConn) LocalAddr() Addr

// Read implements the Conn Read method.
func (c *TCPConn) Read(b []byte) (int, error)

// ReadFrom implements the io.ReaderFrom ReadFrom method.
func (c *TCPConn) ReadFrom(r io.Reader) (int64, error)

// RemoteAddr returns the remote network address.
func (c *TCPConn) RemoteAddr() Addr

// SetDeadline implements the Conn SetDeadline method.
func (c *TCPConn) SetDeadline(t time.Time) error

// SetKeepAlive sets whether the operating system should send keepalive messages on
// the connection.
func (c *TCPConn) SetKeepAlive(keepalive bool) error

// SetKeepAlivePeriod sets period between keep alives.
func (c *TCPConn) SetKeepAlivePeriod(d time.Duration) error

// SetLinger sets the behavior of Close on a connection which still has data
// waiting to be sent or to be acknowledged.
//
// If sec < 0 (the default), the operating system finishes sending the data in the
// background.
//
// If sec == 0, the operating system discards any unsent or unacknowledged data.
//
// If sec > 0, the data is sent in the background as with sec < 0. On some
// operating systems after sec seconds have elapsed any remaining unsent data may
// be discarded.
func (c *TCPConn) SetLinger(sec int) error

// SetNoDelay controls whether the operating system should delay packet
// transmission in hopes of sending fewer packets (Nagle's algorithm). The default
// is true (no delay), meaning that data is sent as soon as possible after a Write.
func (c *TCPConn) SetNoDelay(noDelay bool) error

// SetReadBuffer sets the size of the operating system's receive buffer associated
// with the connection.
func (c *TCPConn) SetReadBuffer(bytes int) error

// SetReadDeadline implements the Conn SetReadDeadline method.
func (c *TCPConn) SetReadDeadline(t time.Time) error

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the connection.
func (c *TCPConn) SetWriteBuffer(bytes int) error

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (c *TCPConn) SetWriteDeadline(t time.Time) error

// Write implements the Conn Write method.
func (c *TCPConn) Write(b []byte) (int, error)

// TCPListener is a TCP network listener. Clients should typically use variables of
// type Listener instead of assuming TCP.
type TCPListener struct {
	// contains filtered or unexported fields
}

// ListenTCP announces on the TCP address laddr and returns a TCP listener. Net
// must be "tcp", "tcp4", or "tcp6". If laddr has a port of 0, ListenTCP will
// choose an available port. The caller can use the Addr method of TCPListener to
// retrieve the chosen address.
func ListenTCP(net string, laddr *TCPAddr) (*TCPListener, error)

// Accept implements the Accept method in the Listener interface; it waits for the
// next call and returns a generic Conn.
func (l *TCPListener) Accept() (Conn, error)

// AcceptTCP accepts the next incoming call and returns the new connection.
func (l *TCPListener) AcceptTCP() (*TCPConn, error)

// Addr returns the listener's network address, a *TCPAddr.
func (l *TCPListener) Addr() Addr

// Close stops listening on the TCP address. Already Accepted connections are not
// closed.
func (l *TCPListener) Close() error

// File returns a copy of the underlying os.File, set to blocking mode. It is the
// caller's responsibility to close f when finished. Closing l does not affect f,
// and closing f does not affect l.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.
func (l *TCPListener) File() (f *os.File, err error)

// SetDeadline sets the deadline associated with the listener. A zero time value
// disables the deadline.
func (l *TCPListener) SetDeadline(t time.Time) error

// UDPAddr represents the address of a UDP end point.
type UDPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

// ResolveUDPAddr parses addr as a UDP address of the form "host:port" or
// "[ipv6-host%zone]:port" and resolves a pair of domain name and port name on the
// network net, which must be "udp", "udp4" or "udp6". A literal address or host
// name for IPv6 must be enclosed in square brackets, as in "[::1]:80",
// "[ipv6-host]:http" or "[ipv6-host%zone]:80".
func ResolveUDPAddr(net, addr string) (*UDPAddr, error)

// Network returns the address's network name, "udp".
func (a *UDPAddr) Network() string

func (a *UDPAddr) String() string

// UDPConn is the implementation of the Conn and PacketConn interfaces for UDP
// network connections.
type UDPConn struct {
	// contains filtered or unexported fields
}

// DialUDP connects to the remote address raddr on the network net, which must be
// "udp", "udp4", or "udp6". If laddr is not nil, it is used as the local address
// for the connection.
func DialUDP(net string, laddr, raddr *UDPAddr) (*UDPConn, error)

// ListenMulticastUDP listens for incoming multicast UDP packets addressed to the
// group address gaddr on ifi, which specifies the interface to join.
// ListenMulticastUDP uses default multicast interface if ifi is nil.
func ListenMulticastUDP(net string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error)

// ListenUDP listens for incoming UDP packets addressed to the local address laddr.
// Net must be "udp", "udp4", or "udp6". If laddr has a port of 0, ListenUDP will
// choose an available port. The LocalAddr method of the returned UDPConn can be
// used to discover the port. The returned connection's ReadFrom and WriteTo
// methods can be used to receive and send UDP packets with per-packet addressing.
func ListenUDP(net string, laddr *UDPAddr) (*UDPConn, error)

// Close closes the connection.
func (c *UDPConn) Close() error

// File sets the underlying os.File to blocking mode and returns a copy. It is the
// caller's responsibility to close f when finished. Closing c does not affect f,
// and closing f does not affect c.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.
func (c *UDPConn) File() (f *os.File, err error)

// LocalAddr returns the local network address.
func (c *UDPConn) LocalAddr() Addr

// Read implements the Conn Read method.
func (c *UDPConn) Read(b []byte) (int, error)

// ReadFrom implements the PacketConn ReadFrom method.
func (c *UDPConn) ReadFrom(b []byte) (int, Addr, error)

// ReadFromUDP reads a UDP packet from c, copying the payload into b. It returns
// the number of bytes copied into b and the return address that was on the packet.
//
// ReadFromUDP can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err error)

// ReadMsgUDP reads a packet from c, copying the payload into b and the associated
// out-of-band data into oob. It returns the number of bytes copied into b, the
// number of bytes copied into oob, the flags that were set on the packet and the
// source address of the packet.
func (c *UDPConn) ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *UDPAddr, err error)

// RemoteAddr returns the remote network address.
func (c *UDPConn) RemoteAddr() Addr

// SetDeadline implements the Conn SetDeadline method.
func (c *UDPConn) SetDeadline(t time.Time) error

// SetReadBuffer sets the size of the operating system's receive buffer associated
// with the connection.
func (c *UDPConn) SetReadBuffer(bytes int) error

// SetReadDeadline implements the Conn SetReadDeadline method.
func (c *UDPConn) SetReadDeadline(t time.Time) error

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the connection.
func (c *UDPConn) SetWriteBuffer(bytes int) error

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (c *UDPConn) SetWriteDeadline(t time.Time) error

// Write implements the Conn Write method.
func (c *UDPConn) Write(b []byte) (int, error)

// WriteMsgUDP writes a packet to addr via c, copying the payload from b and the
// associated out-of-band data from oob. It returns the number of payload and
// out-of-band bytes written.
func (c *UDPConn) WriteMsgUDP(b, oob []byte, addr *UDPAddr) (n, oobn int, err error)

// WriteTo implements the PacketConn WriteTo method.
func (c *UDPConn) WriteTo(b []byte, addr Addr) (int, error)

// WriteToUDP writes a UDP packet to addr via c, copying the payload from b.
//
// WriteToUDP can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline. On
// packet-oriented connections, write timeouts are rare.
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)

// UnixAddr represents the address of a Unix domain socket end point.
type UnixAddr struct {
	Name string
	Net  string
}

// ResolveUnixAddr parses addr as a Unix domain socket address. The string net
// gives the network name, "unix", "unixgram" or "unixpacket".
func ResolveUnixAddr(net, addr string) (*UnixAddr, error)

// Network returns the address's network name, "unix", "unixgram" or "unixpacket".
func (a *UnixAddr) Network() string

func (a *UnixAddr) String() string

// UnixConn is an implementation of the Conn interface for connections to Unix
// domain sockets.
type UnixConn struct {
	// contains filtered or unexported fields
}

// DialUnix connects to the remote address raddr on the network net, which must be
// "unix", "unixgram" or "unixpacket". If laddr is not nil, it is used as the local
// address for the connection.
func DialUnix(net string, laddr, raddr *UnixAddr) (*UnixConn, error)

// ListenUnixgram listens for incoming Unix datagram packets addressed to the local
// address laddr. The network net must be "unixgram". The returned connection's
// ReadFrom and WriteTo methods can be used to receive and send packets with
// per-packet addressing.
func ListenUnixgram(net string, laddr *UnixAddr) (*UnixConn, error)

// Close closes the connection.
func (c *UnixConn) Close() error

// CloseRead shuts down the reading side of the Unix domain connection. Most
// callers should just use Close.
func (c *UnixConn) CloseRead() error

// CloseWrite shuts down the writing side of the Unix domain connection. Most
// callers should just use Close.
func (c *UnixConn) CloseWrite() error

// File sets the underlying os.File to blocking mode and returns a copy. It is the
// caller's responsibility to close f when finished. Closing c does not affect f,
// and closing f does not affect c.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.
func (c *UnixConn) File() (f *os.File, err error)

// LocalAddr returns the local network address.
func (c *UnixConn) LocalAddr() Addr

// Read implements the Conn Read method.
func (c *UnixConn) Read(b []byte) (int, error)

// ReadFrom implements the PacketConn ReadFrom method.
func (c *UnixConn) ReadFrom(b []byte) (int, Addr, error)

// ReadFromUnix reads a packet from c, copying the payload into b. It returns the
// number of bytes copied into b and the source address of the packet.
//
// ReadFromUnix can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *UnixConn) ReadFromUnix(b []byte) (int, *UnixAddr, error)

// ReadMsgUnix reads a packet from c, copying the payload into b and the associated
// out-of-band data into oob. It returns the number of bytes copied into b, the
// number of bytes copied into oob, the flags that were set on the packet, and the
// source address of the packet.
func (c *UnixConn) ReadMsgUnix(b, oob []byte) (n, oobn, flags int, addr *UnixAddr, err error)

// RemoteAddr returns the remote network address.
func (c *UnixConn) RemoteAddr() Addr

// SetDeadline implements the Conn SetDeadline method.
func (c *UnixConn) SetDeadline(t time.Time) error

// SetReadBuffer sets the size of the operating system's receive buffer associated
// with the connection.
func (c *UnixConn) SetReadBuffer(bytes int) error

// SetReadDeadline implements the Conn SetReadDeadline method.
func (c *UnixConn) SetReadDeadline(t time.Time) error

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the connection.
func (c *UnixConn) SetWriteBuffer(bytes int) error

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (c *UnixConn) SetWriteDeadline(t time.Time) error

// Write implements the Conn Write method.
func (c *UnixConn) Write(b []byte) (int, error)

// WriteMsgUnix writes a packet to addr via c, copying the payload from b and the
// associated out-of-band data from oob. It returns the number of payload and
// out-of-band bytes written.
func (c *UnixConn) WriteMsgUnix(b, oob []byte, addr *UnixAddr) (n, oobn int, err error)

// WriteTo implements the PacketConn WriteTo method.
func (c *UnixConn) WriteTo(b []byte, addr Addr) (int, error)

// WriteToUnix writes a packet to addr via c, copying the payload from b.
//
// WriteToUnix can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline. On
// packet-oriented connections, write timeouts are rare.
func (c *UnixConn) WriteToUnix(b []byte, addr *UnixAddr) (int, error)

// UnixListener is a Unix domain socket listener. Clients should typically use
// variables of type Listener instead of assuming Unix domain sockets.
type UnixListener struct {
	// contains filtered or unexported fields
}

// ListenUnix announces on the Unix domain socket laddr and returns a Unix
// listener. The network net must be "unix" or "unixpacket".
func ListenUnix(net string, laddr *UnixAddr) (*UnixListener, error)

// Accept implements the Accept method in the Listener interface; it waits for the
// next call and returns a generic Conn.
func (l *UnixListener) Accept() (Conn, error)

// AcceptUnix accepts the next incoming call and returns the new connection.
func (l *UnixListener) AcceptUnix() (*UnixConn, error)

// Addr returns the listener's network address.
func (l *UnixListener) Addr() Addr

// Close stops listening on the Unix address. Already accepted connections are not
// closed.
func (l *UnixListener) Close() error

// File returns a copy of the underlying os.File, set to blocking mode. It is the
// caller's responsibility to close f when finished. Closing l does not affect f,
// and closing f does not affect l.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.
func (l *UnixListener) File() (*os.File, error)

// SetDeadline sets the deadline associated with the listener. A zero time value
// disables the deadline.
func (l *UnixListener) SetDeadline(t time.Time) error

type UnknownNetworkError string

func (e UnknownNetworkError) Error() string

func (e UnknownNetworkError) Temporary() bool

func (e UnknownNetworkError) Timeout() bool
