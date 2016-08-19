// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package net provides a portable interface for network I/O, including TCP/IP,
// UDP, domain name resolution, and Unix domain sockets.
//
// Although the package provides access to low-level networking primitives, most
// clients will need only the basic interface provided by the Dial, Listen, and
// Accept functions and the associated Conn and Listener interfaces. The
// crypto/tls package uses the same interfaces and similar Dial and Listen
// functions.
//
// The Dial function connects to a server:
//
//     conn, err := net.Dial("tcp", "google.com:80")
//     if err != nil {
//     	// handle error
//     }
//     fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
//     status, err := bufio.NewReader(conn).ReadString('\n')
//     // ...
//
// The Listen function creates servers:
//
//     ln, err := net.Listen("tcp", ":8080")
//     if err != nil {
//     	// handle error
//     }
//     for {
//     	conn, err := ln.Accept()
//     	if err != nil {
//     		// handle error
//     	}
//     	go handleConnection(conn)
//     }
//
//
// Name Resolution
//
// The method for resolving domain names, whether indirectly with functions like
// Dial or directly with functions like LookupHost and LookupAddr, varies by
// operating system.
//
// On Unix systems, the resolver has two options for resolving names. It can use
// a pure Go resolver that sends DNS requests directly to the servers listed in
// /etc/resolv.conf, or it can use a cgo-based resolver that calls C library
// routines such as getaddrinfo and getnameinfo.
//
// By default the pure Go resolver is used, because a blocked DNS request
// consumes only a goroutine, while a blocked C call consumes an operating
// system thread. When cgo is available, the cgo-based resolver is used instead
// under a variety of conditions: on systems that do not let programs make
// direct DNS requests (OS X), when the LOCALDOMAIN environment variable is
// present (even if empty), when the RES_OPTIONS or HOSTALIASES environment
// variable is non-empty, when the ASR_CONFIG environment variable is non-empty
// (OpenBSD only), when /etc/resolv.conf or /etc/nsswitch.conf specify the use
// of features that the Go resolver does not implement, and when the name being
// looked up ends in .local or is an mDNS name.
//
// The resolver decision can be overridden by setting the netdns value of the
// GODEBUG environment variable (see package runtime) to go or cgo, as in:
//
//     export GODEBUG=netdns=go    # force pure Go resolver
//     export GODEBUG=netdns=cgo   # force cgo resolver
//
// The decision can also be forced while building the Go source tree by setting
// the netgo or netcgo build tag.
//
// A numeric netdns setting, as in GODEBUG=netdns=1, causes the resolver to
// print debugging information about its decisions. To force a particular
// resolver while also printing debugging information, join the two settings by
// a plus sign, as in GODEBUG=netdns=go+1.
//
// On Plan 9, the resolver always accesses /net/cs and /net/dns.
//
// On Windows, the resolver always uses C library functions, such as GetAddrInfo
// and DnsQuery.

// Package net provides a portable interface for network I/O, including TCP/IP,
// UDP, domain name resolution, and Unix domain sockets.
//
// Although the package provides access to low-level networking primitives, most
// clients will need only the basic interface provided by the Dial, Listen, and
// Accept functions and the associated Conn and Listener interfaces. The
// crypto/tls package uses the same interfaces and similar Dial and Listen
// functions.
//
// The Dial function connects to a server:
//
//     conn, err := net.Dial("tcp", "golang.org:80")
//     if err != nil {
//     	// handle error
//     }
//     fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
//     status, err := bufio.NewReader(conn).ReadString('\n')
//     // ...
//
// The Listen function creates servers:
//
//     ln, err := net.Listen("tcp", ":8080")
//     if err != nil {
//     	// handle error
//     }
//     for {
//     	conn, err := ln.Accept()
//     	if err != nil {
//     		// handle error
//     	}
//     	go handleConnection(conn)
//     }
//
// Name Resolution
//
// The method for resolving domain names, whether indirectly with functions like
// Dial or directly with functions like LookupHost and LookupAddr, varies by
// operating system.
//
// On Unix systems, the resolver has two options for resolving names. It can use
// a pure Go resolver that sends DNS requests directly to the servers listed in
// /etc/resolv.conf, or it can use a cgo-based resolver that calls C library
// routines such as getaddrinfo and getnameinfo.
//
// By default the pure Go resolver is used, because a blocked DNS request
// consumes only a goroutine, while a blocked C call consumes an operating
// system thread. When cgo is available, the cgo-based resolver is used instead
// under a variety of conditions: on systems that do not let programs make
// direct DNS requests (OS X), when the LOCALDOMAIN environment variable is
// present (even if empty), when the RES_OPTIONS or HOSTALIASES environment
// variable is non-empty, when the ASR_CONFIG environment variable is non-empty
// (OpenBSD only), when /etc/resolv.conf or /etc/nsswitch.conf specify the use
// of features that the Go resolver does not implement, and when the name being
// looked up ends in .local or is an mDNS name.
//
// The resolver decision can be overridden by setting the netdns value of the
// GODEBUG environment variable (see package runtime) to go or cgo, as in:
//
//     export GODEBUG=netdns=go    # force pure Go resolver
//     export GODEBUG=netdns=cgo   # force cgo resolver
//
// The decision can also be forced while building the Go source tree by setting
// the netgo or netcgo build tag.
//
// A numeric netdns setting, as in GODEBUG=netdns=1, causes the resolver to
// print debugging information about its decisions. To force a particular
// resolver while also printing debugging information, join the two settings by
// a plus sign, as in GODEBUG=netdns=go+1.
//
// On Plan 9, the resolver always accesses /net/cs and /net/dns.
//
// On Windows, the resolver always uses C library functions, such as GetAddrInfo
// and DnsQuery.
package net

import (
    "C"
    "context"
    "errors"
    "golang.org/x/net/route"
    "internal/nettrace"
    "internal/singleflight"
    "io"
    "math/rand"
    "os"
    "runtime"
    "sort"
    "sync"
    "sync/atomic"
    "syscall"
    "time"
    "unsafe"
)


const (
	FlagUp           Flags = 1 << iota // interface is up
	FlagBroadcast                      // interface supports broadcast access capability
	FlagLoopback                       // interface is a loopback interface
	FlagPointToPoint                   // interface belongs to a point-to-point link
	FlagMulticast                      // interface supports multicast access capability

)


// IP address lengths (bytes).
const (
	IPv4len = 4
	IPv6len = 16
)


// Various errors contained in OpError.
var (
	ErrWriteToConnected = errors.New("use of WriteTo with pre-connected connection")
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


// Addr represents a network end point address.
type Addr interface {
	Network() string // name of the network
	String() string  // string form of address
}



type AddrError struct {
	Err  string
	Addr string
}


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


// DNSConfigError represents an error reading the machine's DNS configuration.
// (No longer used; kept for compatibility.)
type DNSConfigError struct {
	Err error
}


// DNSError represents a DNS lookup error.
type DNSError struct {
	Err         string // description of the error
	Name        string // name looked for
	Server      string // server used
	IsTimeout   bool   // if true, timed out; not all timeouts set this
	IsTemporary bool   // if true, error is temporary; not all errors set this
}


// A Dialer contains options for connecting to an address.
//
// The zero value for each field is equivalent to dialing
// without that option. Dialing with the zero value of Dialer
// is therefore equivalent to just calling the Dial function.
type Dialer struct {
	// Timeout is the maximum amount of time a dial will wait for
	// a connect to complete. If Deadline is also set, it may fail
	// earlier.
	//
	// The default is no timeout.
	//
	// When dialing a name with multiple IP addresses, the timeout
	// may be divided between them.
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

	// DualStack enables RFC 6555-compliant "Happy Eyeballs" dialing
	// when the network is "tcp" and the destination is a host name
	// with both IPv4 and IPv6 addresses. This allows a client to
	// tolerate networks where one address family is silently broken.
	DualStack bool

	// FallbackDelay specifies the length of time to wait before
	// spawning a fallback connection, when DualStack is enabled.
	// If zero, a default delay of 300ms is used.
	FallbackDelay time.Duration

	// KeepAlive specifies the keep-alive period for an active
	// network connection.
	// If zero, keep-alives are not enabled. Network protocols
	// that do not support keep-alives ignore this field.
	KeepAlive time.Duration

	// Cancel is an optional channel whose closure indicates that
	// the dial should be canceled. Not all types of dials support
	// cancelation.
	//
	// Deprecated: Use DialContext instead.
	Cancel <-chan struct{}
}


// An Error represents a network error.
type Error interface {
	error
	Timeout() bool   // Is the error a timeout?
	Temporary() bool // Is the error temporary?
}



type Flags uint


// A HardwareAddr represents a physical hardware address.
type HardwareAddr []byte


// An IP is a single IP address, a slice of bytes.
// Functions in this package accept either 4-byte (IPv4)
// or 16-byte (IPv6) slices as input.
//
// Note that in this documentation, referring to an
// IP address as an IPv4 address or an IPv6 address
// is a semantic property of the address, not just the
// length of the byte slice: a 16-byte slice can still
// be an IPv4 address.
type IP []byte


// IPAddr represents the address of an IP end point.
type IPAddr struct {
	IP   IP
	Zone string // IPv6 scoped addressing zone
}


// IPConn is the implementation of the Conn and PacketConn interfaces
// for IP network connections.
type IPConn struct {
}


// An IP mask is an IP address.
type IPMask []byte


// An IPNet represents an IP network.
type IPNet struct {
	IP   IP     // network number
	Mask IPMask // network mask
}


// Interface represents a mapping between network interface name
// and index.  It also represents network interface facility
// information.

// Interface represents a mapping between network interface name
// and index. It also represents network interface facility
// information.
type Interface struct {
	Index        int          // positive integer that starts at one, zero is never used
	MTU          int          // maximum transmission unit
	Name         string       // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr HardwareAddr // IEEE MAC-48, EUI-48 and EUI-64 form
	Flags        Flags        // e.g., FlagUp, FlagLoopback, FlagMulticast
}



type InvalidAddrError string


// A Listener is a generic network listener for stream-oriented protocols.
//
// Multiple goroutines may invoke methods on a Listener simultaneously.
type Listener interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (Conn, error)

	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error

	// Addr returns the listener's network address.
	Addr() Addr
}


// An MX represents a single DNS MX record.
type MX struct {
	Host string
	Pref uint16
}


// An NS represents a single DNS NS record.
type NS struct {
	Host string
}


// OpError is the error type usually returned by functions in the net
// package. It describes the operation, network type, and address of
// an error.
type OpError struct {
	// Op is the operation which caused the error, such as
	// "read" or "write".
	Op string

	// Net is the network type on which this error occurred,
	// such as "tcp" or "udp6".
	Net string

	// For operations involving a remote network connection, like
	// Dial, Read, or Write, Source is the corresponding local
	// network address.
	Source Addr

	// Addr is the network address for which this error occurred.
	// For local operations, like Listen or SetDeadline, Addr is
	// the address of the local endpoint being manipulated.
	// For operations involving a remote network connection, like
	// Dial, Read, or Write, Addr is the remote address of that
	// connection.
	Addr Addr

	// Err is the error that occurred during the operation.
	Err error
}


// PacketConn is a generic packet-oriented network connection.
//
// Multiple goroutines may invoke methods on a PacketConn simultaneously.
type PacketConn interface {
	// ReadFrom reads a packet from the connection,
	// copying the payload into b. It returns the number of
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


// A ParseError is the error type of literal network address parsers.
type ParseError struct {
	// Type is the type of string that was expected, such as
	// "IP address", "CIDR address".
	Type string

	// Text is the malformed text string.
	Text string
}


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


// TCPConn is an implementation of the Conn interface for TCP network
// connections.
type TCPConn struct {
}


// TCPListener is a TCP network listener.  Clients should typically
// use variables of type Listener instead of assuming TCP.

// TCPListener is a TCP network listener. Clients should typically
// use variables of type Listener instead of assuming TCP.
type TCPListener struct {
	fd *netFD
}


// UDPAddr represents the address of a UDP end point.
type UDPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone
}


// UDPConn is the implementation of the Conn and PacketConn interfaces
// for UDP network connections.
type UDPConn struct {
}


// UnixAddr represents the address of a Unix domain socket end point.
type UnixAddr struct {
	Name string
	Net  string
}


// UnixConn is an implementation of the Conn interface for connections
// to Unix domain sockets.
type UnixConn struct {
}


// UnixListener is a Unix domain socket listener.  Clients should
// typically use variables of type Listener instead of assuming Unix
// domain sockets.

// UnixListener is a Unix domain socket listener. Clients should
// typically use variables of type Listener instead of assuming Unix
// domain sockets.
type UnixListener struct {
	fd     *netFD
	path   string
	unlink bool
}



type UnknownNetworkError string


// CIDRMask returns an IPMask consisting of `ones' 1 bits
// followed by 0s up to a total length of `bits' bits.
// For a mask of this form, CIDRMask is the inverse of IPMask.Size.
func CIDRMask(ones, bits int) IPMask

// Dial connects to the address on the named network.
//
// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
// "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4"
// (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and
// "unixpacket".
//
// For TCP and UDP networks, addresses have the form host:port.
// If host is a literal IPv6 address it must be enclosed
// in square brackets as in "[::1]:80" or "[ipv6-host%zone]:80".
// The functions JoinHostPort and SplitHostPort manipulate addresses
// in this form.
// If the host is empty, as in ":80", the local system is assumed.
//
// Examples:
//     Dial("tcp", "12.34.56.78:80")
//     Dial("tcp", "google.com:http")
//     Dial("tcp", "[2001:db8::1]:http")
//     Dial("tcp", "[fe80::1%lo0]:80")
//     Dial("tcp", ":80")
//
// For IP networks, the network must be "ip", "ip4" or "ip6" followed
// by a colon and a protocol number or name and the addr must be a
// literal IP address.
//
// Examples:
//     Dial("ip4:1", "127.0.0.1")
//     Dial("ip6:ospf", "::1")
//
// For Unix networks, the address must be a file system path.

// Dial connects to the address on the named network.
//
// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
// "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4"
// (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and
// "unixpacket".
//
// For TCP and UDP networks, addresses have the form host:port.
// If host is a literal IPv6 address it must be enclosed
// in square brackets as in "[::1]:80" or "[ipv6-host%zone]:80".
// The functions JoinHostPort and SplitHostPort manipulate addresses
// in this form.
// If the host is empty, as in ":80", the local system is assumed.
//
// Examples:
//     Dial("tcp", "192.0.2.1:80")
//     Dial("tcp", "golang.org:http")
//     Dial("tcp", "[2001:db8::1]:http")
//     Dial("tcp", "[fe80::1%lo0]:80")
//     Dial("tcp", ":80")
//
// For IP networks, the network must be "ip", "ip4" or "ip6" followed
// by a colon and a protocol number or name and the addr must be a
// literal IP address.
//
// Examples:
//     Dial("ip4:1", "192.0.2.1")
//     Dial("ip6:ipv6-icmp", "2001:db8::1")
//
// For Unix networks, the address must be a file system path.
func Dial(network, address string) (Conn, error)

// DialIP connects to the remote address raddr on the network protocol
// netProto, which must be "ip", "ip4", or "ip6" followed by a colon
// and a protocol number or name.
func DialIP(netProto string, laddr, raddr *IPAddr) (*IPConn, error)

// DialTCP connects to the remote address raddr on the network net,
// which must be "tcp", "tcp4", or "tcp6".  If laddr is not nil, it is
// used as the local address for the connection.
func DialTCP(net string, laddr, raddr *TCPAddr) (*TCPConn, error)

// DialTimeout acts like Dial but takes a timeout.
// The timeout includes name resolution, if required.
func DialTimeout(network, address string, timeout time.Duration) (Conn, error)

// DialUDP connects to the remote address raddr on the network net,
// which must be "udp", "udp4", or "udp6".  If laddr is not nil, it is
// used as the local address for the connection.
func DialUDP(net string, laddr, raddr *UDPAddr) (*UDPConn, error)

// DialUnix connects to the remote address raddr on the network net,
// which must be "unix", "unixgram" or "unixpacket".  If laddr is not
// nil, it is used as the local address for the connection.
func DialUnix(net string, laddr, raddr *UnixAddr) (*UnixConn, error)

// FileConn returns a copy of the network connection corresponding to
// the open file f.
// It is the caller's responsibility to close f when finished.
// Closing c does not affect f, and closing f does not affect c.
func FileConn(f *os.File) (c Conn, err error)

// FileListener returns a copy of the network listener corresponding
// to the open file f.
// It is the caller's responsibility to close ln when finished.
// Closing ln does not affect f, and closing f does not affect ln.
func FileListener(f *os.File) (ln Listener, err error)

// FilePacketConn returns a copy of the packet network connection
// corresponding to the open file f.
// It is the caller's responsibility to close f when finished.
// Closing c does not affect f, and closing f does not affect c.
func FilePacketConn(f *os.File) (c PacketConn, err error)

// IPv4 returns the IP address (in 16-byte form) of the
// IPv4 address a.b.c.d.
func IPv4(a, b, c, d byte) IP

// IPv4Mask returns the IP mask (in 4-byte form) of the
// IPv4 mask a.b.c.d.
func IPv4Mask(a, b, c, d byte) IPMask

// InterfaceAddrs returns a list of the system's network interface
// addresses.
func InterfaceAddrs() ([]Addr, error)

// InterfaceByIndex returns the interface specified by index.
func InterfaceByIndex(index int) (*Interface, error)

// InterfaceByName returns the interface specified by name.
func InterfaceByName(name string) (*Interface, error)

// Interfaces returns a list of the system's network interfaces.
func Interfaces() ([]Interface, error)

// JoinHostPort combines host and port into a network address of the
// form "host:port" or, if host contains a colon or a percent sign,
// "[host]:port".
func JoinHostPort(host, port string) string

// Listen announces on the local network address laddr.
// The network net must be a stream-oriented network: "tcp", "tcp4",
// "tcp6", "unix" or "unixpacket".
// For TCP and UDP, the syntax of laddr is "host:port", like "127.0.0.1:8080".
// If host is omitted, as in ":8080", Listen listens on all available interfaces
// instead of just the interface with the given host address.
// See Dial for more details about address syntax.
func Listen(net, laddr string) (Listener, error)

// ListenIP listens for incoming IP packets addressed to the local
// address laddr.  The returned connection's ReadFrom and WriteTo
// methods can be used to receive and send IP packets with per-packet
// addressing.

// ListenIP listens for incoming IP packets addressed to the local
// address laddr. The returned connection's ReadFrom and WriteTo
// methods can be used to receive and send IP packets with per-packet
// addressing.
func ListenIP(netProto string, laddr *IPAddr) (*IPConn, error)

// ListenMulticastUDP listens for incoming multicast UDP packets
// addressed to the group address gaddr on the interface ifi.
// Network must be "udp", "udp4" or "udp6".
// ListenMulticastUDP uses the system-assigned multicast interface
// when ifi is nil, although this is not recommended because the
// assignment depends on platforms and sometimes it might require
// routing configuration.
//
// ListenMulticastUDP is just for convenience of simple, small
// applications. There are golang.org/x/net/ipv4 and
// golang.org/x/net/ipv6 packages for general purpose uses.
func ListenMulticastUDP(network string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error)

// ListenPacket announces on the local network address laddr. The network net
// must be a packet-oriented network: "udp", "udp4", "udp6", "ip", "ip4", "ip6"
// or "unixgram". For TCP and UDP, the syntax of laddr is "host:port", like
// "127.0.0.1:8080". If host is omitted, as in ":8080", ListenPacket listens on
// all available interfaces instead of just the interface with the given host
// address. See Dial for the syntax of laddr.
func ListenPacket(net, laddr string) (PacketConn, error)

// ListenTCP announces on the TCP address laddr and returns a TCP
// listener.  Net must be "tcp", "tcp4", or "tcp6".  If laddr has a
// port of 0, ListenTCP will choose an available port.  The caller can
// use the Addr method of TCPListener to retrieve the chosen address.

// ListenTCP announces on the TCP address laddr and returns a TCP
// listener. Net must be "tcp", "tcp4", or "tcp6".  If laddr has a
// port of 0, ListenTCP will choose an available port. The caller can
// use the Addr method of TCPListener to retrieve the chosen address.
func ListenTCP(net string, laddr *TCPAddr) (*TCPListener, error)

// ListenUDP listens for incoming UDP packets addressed to the local
// address laddr.  Net must be "udp", "udp4", or "udp6".  If laddr has
// a port of 0, ListenUDP will choose an available port.
// The LocalAddr method of the returned UDPConn can be used to
// discover the port.  The returned connection's ReadFrom and WriteTo
// methods can be used to receive and send UDP packets with per-packet
// addressing.

// ListenUDP listens for incoming UDP packets addressed to the local
// address laddr. Net must be "udp", "udp4", or "udp6".  If laddr has
// a port of 0, ListenUDP will choose an available port.
// The LocalAddr method of the returned UDPConn can be used to
// discover the port. The returned connection's ReadFrom and WriteTo
// methods can be used to receive and send UDP packets with per-packet
// addressing.
func ListenUDP(net string, laddr *UDPAddr) (*UDPConn, error)

// ListenUnix announces on the Unix domain socket laddr and returns a
// Unix listener.  The network net must be "unix" or "unixpacket".

// ListenUnix announces on the Unix domain socket laddr and returns a
// Unix listener. The network net must be "unix" or "unixpacket".
func ListenUnix(net string, laddr *UnixAddr) (*UnixListener, error)

// ListenUnixgram listens for incoming Unix datagram packets addressed
// to the local address laddr.  The network net must be "unixgram".
// The returned connection's ReadFrom and WriteTo methods can be used
// to receive and send packets with per-packet addressing.

// ListenUnixgram listens for incoming Unix datagram packets addressed
// to the local address laddr. The network net must be "unixgram".
// The returned connection's ReadFrom and WriteTo methods can be used
// to receive and send packets with per-packet addressing.
func ListenUnixgram(net string, laddr *UnixAddr) (*UnixConn, error)

// LookupAddr performs a reverse lookup for the given address, returning a list
// of names mapping to that address.
func LookupAddr(addr string) (names []string, err error)

// LookupCNAME returns the canonical DNS host for the given name.
// Callers that do not care about the canonical name can call
// LookupHost or LookupIP directly; both take care of resolving
// the canonical name as part of the lookup.
func LookupCNAME(name string) (cname string, err error)

// LookupHost looks up the given host using the local resolver.
// It returns an array of that host's addresses.
func LookupHost(host string) (addrs []string, err error)

// LookupIP looks up host using the local resolver.
// It returns an array of that host's IPv4 and IPv6 addresses.
func LookupIP(host string) (ips []IP, err error)

// LookupMX returns the DNS MX records for the given domain name sorted by
// preference.
func LookupMX(name string) (mxs []*MX, err error)

// LookupNS returns the DNS NS records for the given domain name.
func LookupNS(name string) (nss []*NS, err error)

// LookupPort looks up the port for the given network and service.
func LookupPort(network, service string) (port int, err error)

// LookupSRV tries to resolve an SRV query of the given service,
// protocol, and domain name.  The proto is "tcp" or "udp".
// The returned records are sorted by priority and randomized
// by weight within a priority.
//
// LookupSRV constructs the DNS name to look up following RFC 2782.
// That is, it looks up _service._proto.name.  To accommodate services
// publishing SRV records under non-standard names, if both service
// and proto are empty strings, LookupSRV looks up name directly.

// LookupSRV tries to resolve an SRV query of the given service,
// protocol, and domain name. The proto is "tcp" or "udp".
// The returned records are sorted by priority and randomized
// by weight within a priority.
//
// LookupSRV constructs the DNS name to look up following RFC 2782.
// That is, it looks up _service._proto.name. To accommodate services
// publishing SRV records under non-standard names, if both service
// and proto are empty strings, LookupSRV looks up name directly.
func LookupSRV(service, proto, name string) (cname string, addrs []*SRV, err error)

// LookupTXT returns the DNS TXT records for the given domain name.
func LookupTXT(name string) (txts []string, err error)

// ParseCIDR parses s as a CIDR notation IP address and mask,
// like "192.168.100.1/24" or "2001:DB8::/48", as defined in
// RFC 4632 and RFC 4291.
//
// It returns the IP address and the network implied by the IP
// and mask.  For example, ParseCIDR("192.168.100.1/16") returns
// the IP address 192.168.100.1 and the network 192.168.0.0/16.

// ParseCIDR parses s as a CIDR notation IP address and mask,
// like "192.0.2.0/24" or "2001:db8::/32", as defined in
// RFC 4632 and RFC 4291.
//
// It returns the IP address and the network implied by the IP
// and mask. For example, ParseCIDR("198.51.100.1/24") returns
// the IP address 198.51.100.1 and the network 198.51.100.0/24.
func ParseCIDR(s string) (IP, *IPNet, error)

// ParseIP parses s as an IP address, returning the result.
// The string s can be in dotted decimal ("74.125.19.99")
// or IPv6 ("2001:4860:0:2001::68") form.
// If s is not a valid textual representation of an IP address,
// ParseIP returns nil.

// ParseIP parses s as an IP address, returning the result.
// The string s can be in dotted decimal ("192.0.2.1")
// or IPv6 ("2001:db8::68") form.
// If s is not a valid textual representation of an IP address,
// ParseIP returns nil.
func ParseIP(s string) IP

// ParseMAC parses s as an IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet
// IP over InfiniBand link-layer address using one of the following formats:
//   01:23:45:67:89:ab
//   01:23:45:67:89:ab:cd:ef
//   01:23:45:67:89:ab:cd:ef:00:00:01:23:45:67:89:ab:cd:ef:00:00
//   01-23-45-67-89-ab
//   01-23-45-67-89-ab-cd-ef
//   01-23-45-67-89-ab-cd-ef-00-00-01-23-45-67-89-ab-cd-ef-00-00
//   0123.4567.89ab
//   0123.4567.89ab.cdef
//   0123.4567.89ab.cdef.0000.0123.4567.89ab.cdef.0000
func ParseMAC(s string) (hw HardwareAddr, err error)

// Pipe creates a synchronous, in-memory, full duplex
// network connection; both ends implement the Conn interface.
// Reads on one end are matched with writes on the other,
// copying data directly between the two; there is no internal
// buffering.
func Pipe() (Conn, Conn)

// ResolveIPAddr parses addr as an IP address of the form "host" or
// "ipv6-host%zone" and resolves the domain name on the network net,
// which must be "ip", "ip4" or "ip6".
func ResolveIPAddr(net, addr string) (*IPAddr, error)

// ResolveTCPAddr parses addr as a TCP address of the form "host:port"
// or "[ipv6-host%zone]:port" and resolves a pair of domain name and
// port name on the network net, which must be "tcp", "tcp4" or
// "tcp6".  A literal address or host name for IPv6 must be enclosed
// in square brackets, as in "[::1]:80", "[ipv6-host]:http" or
// "[ipv6-host%zone]:80".
func ResolveTCPAddr(net, addr string) (*TCPAddr, error)

// ResolveUDPAddr parses addr as a UDP address of the form "host:port"
// or "[ipv6-host%zone]:port" and resolves a pair of domain name and
// port name on the network net, which must be "udp", "udp4" or
// "udp6".  A literal address or host name for IPv6 must be enclosed
// in square brackets, as in "[::1]:80", "[ipv6-host]:http" or
// "[ipv6-host%zone]:80".
func ResolveUDPAddr(net, addr string) (*UDPAddr, error)

// ResolveUnixAddr parses addr as a Unix domain socket address.
// The string net gives the network name, "unix", "unixgram" or
// "unixpacket".
func ResolveUnixAddr(net, addr string) (*UnixAddr, error)

// SplitHostPort splits a network address of the form "host:port",
// "[host]:port" or "[ipv6-host%zone]:port" into host or
// ipv6-host%zone and port.  A literal address or host name for IPv6
// must be enclosed in square brackets, as in "[::1]:80",
// "[ipv6-host]:http" or "[ipv6-host%zone]:80".

// SplitHostPort splits a network address of the form "host:port",
// "[host]:port" or "[ipv6-host%zone]:port" into host or
// ipv6-host%zone and port. A literal address or host name for IPv6
// must be enclosed in square brackets, as in "[::1]:80",
// "[ipv6-host]:http" or "[ipv6-host%zone]:80".
func SplitHostPort(hostport string) (host, port string, err error)

func (*AddrError) Error() string

func (*AddrError) Temporary() bool

func (*AddrError) Timeout() bool

func (*DNSConfigError) Error() string

func (*DNSConfigError) Temporary() bool

func (*DNSConfigError) Timeout() bool

func (*DNSError) Error() string

// Temporary reports whether the DNS error is known to be temporary.
// This is not always known; a DNS lookup may fail due to a temporary
// error and return a DNSError for which Temporary returns false.
func (*DNSError) Temporary() bool

// Timeout reports whether the DNS lookup is known to have timed out.
// This is not always known; a DNS lookup may fail due to a timeout
// and return a DNSError for which Timeout returns false.
func (*DNSError) Timeout() bool

// Dial connects to the address on the named network.
//
// See func Dial for a description of the network and address
// parameters.
func (*Dialer) Dial(network, address string) (Conn, error)

// DialContext connects to the address on the named network using
// the provided context.
//
// The provided Context must be non-nil. If the context expires before
// the connection is complete, an error is returned. Once successfully
// connected, any expiration of the context will not affect the
// connection.
//
// See func Dial for a description of the network and address
// parameters.
func (*Dialer) DialContext(ctx context.Context, network, address string) (Conn, error)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The IP address is expected in a form accepted by ParseIP.
func (*IP) UnmarshalText(text []byte) error

// Network returns the address's network name, "ip".
func (*IPAddr) Network() string

func (*IPAddr) String() string

// ReadFrom implements the PacketConn ReadFrom method.
func (*IPConn) ReadFrom(b []byte) (int, Addr, error)

// ReadFromIP reads an IP packet from c, copying the payload into b.
// It returns the number of bytes copied into b and the return address
// that was on the packet.
//
// ReadFromIP can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetReadDeadline.
func (*IPConn) ReadFromIP(b []byte) (int, *IPAddr, error)

// ReadMsgIP reads a packet from c, copying the payload into b and the
// associated out-of-band data into oob.  It returns the number of
// bytes copied into b, the number of bytes copied into oob, the flags
// that were set on the packet and the source address of the packet.

// ReadMsgIP reads a packet from c, copying the payload into b and the
// associated out-of-band data into oob. It returns the number of
// bytes copied into b, the number of bytes copied into oob, the flags
// that were set on the packet and the source address of the packet.
func (*IPConn) ReadMsgIP(b, oob []byte) (n, oobn, flags int, addr *IPAddr, err error)

// WriteMsgIP writes a packet to addr via c, copying the payload from
// b and the associated out-of-band data from oob.  It returns the
// number of payload and out-of-band bytes written.

// WriteMsgIP writes a packet to addr via c, copying the payload from
// b and the associated out-of-band data from oob. It returns the
// number of payload and out-of-band bytes written.
func (*IPConn) WriteMsgIP(b, oob []byte, addr *IPAddr) (n, oobn int, err error)

// WriteTo implements the PacketConn WriteTo method.
func (*IPConn) WriteTo(b []byte, addr Addr) (int, error)

// WriteToIP writes an IP packet to addr via c, copying the payload
// from b.
//
// WriteToIP can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetWriteDeadline.  On packet-oriented connections, write timeouts
// are rare.

// WriteToIP writes an IP packet to addr via c, copying the payload
// from b.
//
// WriteToIP can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetWriteDeadline. On packet-oriented connections, write timeouts
// are rare.
func (*IPConn) WriteToIP(b []byte, addr *IPAddr) (int, error)

// Contains reports whether the network includes ip.
func (*IPNet) Contains(ip IP) bool

// Network returns the address's network name, "ip+net".
func (*IPNet) Network() string

// String returns the CIDR notation of n like "192.168.100.1/24"
// or "2001:DB8::/48" as defined in RFC 4632 and RFC 4291.
// If the mask is not in the canonical form, it returns the
// string which consists of an IP address, followed by a slash
// character and a mask expressed as hexadecimal form with no
// punctuation like "192.168.100.1/c000ff00".

// String returns the CIDR notation of n like "192.0.2.1/24"
// or "2001:db8::/48" as defined in RFC 4632 and RFC 4291.
// If the mask is not in the canonical form, it returns the
// string which consists of an IP address, followed by a slash
// character and a mask expressed as hexadecimal form with no
// punctuation like "198.51.100.1/c000ff00".
func (*IPNet) String() string

// Addrs returns interface addresses for a specific interface.
func (*Interface) Addrs() ([]Addr, error)

// MulticastAddrs returns multicast, joined group addresses for
// a specific interface.
func (*Interface) MulticastAddrs() ([]Addr, error)

func (*OpError) Error() string

func (*OpError) Temporary() bool

func (*OpError) Timeout() bool

func (*ParseError) Error() string

// Network returns the address's network name, "tcp".
func (*TCPAddr) Network() string

func (*TCPAddr) String() string

// CloseRead shuts down the reading side of the TCP connection.
// Most callers should just use Close.
func (*TCPConn) CloseRead() error

// CloseWrite shuts down the writing side of the TCP connection.
// Most callers should just use Close.
func (*TCPConn) CloseWrite() error

// ReadFrom implements the io.ReaderFrom ReadFrom method.
func (*TCPConn) ReadFrom(r io.Reader) (int64, error)

// SetKeepAlive sets whether the operating system should send
// keepalive messages on the connection.
func (*TCPConn) SetKeepAlive(keepalive bool) error

// SetKeepAlivePeriod sets period between keep alives.
func (*TCPConn) SetKeepAlivePeriod(d time.Duration) error

// SetLinger sets the behavior of Close on a connection which still
// has data waiting to be sent or to be acknowledged.
//
// If sec < 0 (the default), the operating system finishes sending the
// data in the background.
//
// If sec == 0, the operating system discards any unsent or
// unacknowledged data.
//
// If sec > 0, the data is sent in the background as with sec < 0. On
// some operating systems after sec seconds have elapsed any remaining
// unsent data may be discarded.
func (*TCPConn) SetLinger(sec int) error

// SetNoDelay controls whether the operating system should delay
// packet transmission in hopes of sending fewer packets (Nagle's
// algorithm).  The default is true (no delay), meaning that data is
// sent as soon as possible after a Write.
func (*TCPConn) SetNoDelay(noDelay bool) error

// Accept implements the Accept method in the Listener interface; it
// waits for the next call and returns a generic Conn.
func (*TCPListener) Accept() (Conn, error)

// AcceptTCP accepts the next incoming call and returns the new
// connection.
func (*TCPListener) AcceptTCP() (*TCPConn, error)

// Addr returns the listener's network address, a *TCPAddr.
// The Addr returned is shared by all invocations of Addr, so
// do not modify it.
func (*TCPListener) Addr() Addr

// Close stops listening on the TCP address.
// Already Accepted connections are not closed.
func (*TCPListener) Close() error

// File returns a copy of the underlying os.File, set to blocking
// mode.  It is the caller's responsibility to close f when finished.
// Closing l does not affect f, and closing f does not affect l.
//
// The returned os.File's file descriptor is different from the
// connection's.  Attempting to change properties of the original
// using this duplicate may or may not have the desired effect.

// File returns a copy of the underlying os.File, set to blocking
// mode. It is the caller's responsibility to close f when finished.
// Closing l does not affect f, and closing f does not affect l.
//
// The returned os.File's file descriptor is different from the
// connection's. Attempting to change properties of the original
// using this duplicate may or may not have the desired effect.
func (*TCPListener) File() (f *os.File, err error)

// SetDeadline sets the deadline associated with the listener.
// A zero time value disables the deadline.
func (*TCPListener) SetDeadline(t time.Time) error

// Network returns the address's network name, "udp".
func (*UDPAddr) Network() string

func (*UDPAddr) String() string

// ReadFrom implements the PacketConn ReadFrom method.
func (*UDPConn) ReadFrom(b []byte) (int, Addr, error)

// ReadFromUDP reads a UDP packet from c, copying the payload into b.
// It returns the number of bytes copied into b and the return address
// that was on the packet.
//
// ReadFromUDP can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetReadDeadline.
func (*UDPConn) ReadFromUDP(b []byte) (int, *UDPAddr, error)

// ReadMsgUDP reads a packet from c, copying the payload into b and
// the associated out-of-band data into oob.  It returns the number
// of bytes copied into b, the number of bytes copied into oob, the
// flags that were set on the packet and the source address of the
// packet.

// ReadMsgUDP reads a packet from c, copying the payload into b and
// the associated out-of-band data into oob. It returns the number
// of bytes copied into b, the number of bytes copied into oob, the
// flags that were set on the packet and the source address of the
// packet.
func (*UDPConn) ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *UDPAddr, err error)

// WriteMsgUDP writes a packet to addr via c if c isn't connected, or
// to c's remote destination address if c is connected (in which case
// addr must be nil).  The payload is copied from b and the associated
// out-of-band data is copied from oob.  It returns the number of
// payload and out-of-band bytes written.

// WriteMsgUDP writes a packet to addr via c if c isn't connected, or
// to c's remote destination address if c is connected (in which case
// addr must be nil).  The payload is copied from b and the associated
// out-of-band data is copied from oob. It returns the number of
// payload and out-of-band bytes written.
func (*UDPConn) WriteMsgUDP(b, oob []byte, addr *UDPAddr) (n, oobn int, err error)

// WriteTo implements the PacketConn WriteTo method.
func (*UDPConn) WriteTo(b []byte, addr Addr) (int, error)

// WriteToUDP writes a UDP packet to addr via c, copying the payload
// from b.
//
// WriteToUDP can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetWriteDeadline.  On packet-oriented connections, write timeouts
// are rare.

// WriteToUDP writes a UDP packet to addr via c, copying the payload
// from b.
//
// WriteToUDP can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetWriteDeadline. On packet-oriented connections, write timeouts
// are rare.
func (*UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)

// Network returns the address's network name, "unix", "unixgram" or
// "unixpacket".
func (*UnixAddr) Network() string

func (*UnixAddr) String() string

// CloseRead shuts down the reading side of the Unix domain connection.
// Most callers should just use Close.
func (*UnixConn) CloseRead() error

// CloseWrite shuts down the writing side of the Unix domain connection.
// Most callers should just use Close.
func (*UnixConn) CloseWrite() error

// ReadFrom implements the PacketConn ReadFrom method.
func (*UnixConn) ReadFrom(b []byte) (int, Addr, error)

// ReadFromUnix reads a packet from c, copying the payload into b.  It
// returns the number of bytes copied into b and the source address of
// the packet.
//
// ReadFromUnix can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetReadDeadline.

// ReadFromUnix reads a packet from c, copying the payload into b. It
// returns the number of bytes copied into b and the source address of
// the packet.
//
// ReadFromUnix can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetReadDeadline.
func (*UnixConn) ReadFromUnix(b []byte) (int, *UnixAddr, error)

// ReadMsgUnix reads a packet from c, copying the payload into b and
// the associated out-of-band data into oob.  It returns the number of
// bytes copied into b, the number of bytes copied into oob, the flags
// that were set on the packet, and the source address of the packet.

// ReadMsgUnix reads a packet from c, copying the payload into b and
// the associated out-of-band data into oob. It returns the number of
// bytes copied into b, the number of bytes copied into oob, the flags
// that were set on the packet, and the source address of the packet.
func (*UnixConn) ReadMsgUnix(b, oob []byte) (n, oobn, flags int, addr *UnixAddr, err error)

// WriteMsgUnix writes a packet to addr via c, copying the payload
// from b and the associated out-of-band data from oob.  It returns
// the number of payload and out-of-band bytes written.

// WriteMsgUnix writes a packet to addr via c, copying the payload
// from b and the associated out-of-band data from oob. It returns
// the number of payload and out-of-band bytes written.
func (*UnixConn) WriteMsgUnix(b, oob []byte, addr *UnixAddr) (n, oobn int, err error)

// WriteTo implements the PacketConn WriteTo method.
func (*UnixConn) WriteTo(b []byte, addr Addr) (int, error)

// WriteToUnix writes a packet to addr via c, copying the payload from b.
//
// WriteToUnix can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetWriteDeadline.  On packet-oriented connections, write timeouts
// are rare.

// WriteToUnix writes a packet to addr via c, copying the payload from b.
//
// WriteToUnix can be made to time out and return an error with
// Timeout() == true after a fixed time limit; see SetDeadline and
// SetWriteDeadline. On packet-oriented connections, write timeouts
// are rare.
func (*UnixConn) WriteToUnix(b []byte, addr *UnixAddr) (int, error)

// Accept implements the Accept method in the Listener interface; it
// waits for the next call and returns a generic Conn.

// Accept implements the Accept method in the Listener interface.
// Returned connections will be of type *UnixConn.
func (*UnixListener) Accept() (Conn, error)

// AcceptUnix accepts the next incoming call and returns the new
// connection.
func (*UnixListener) AcceptUnix() (*UnixConn, error)

// Addr returns the listener's network address.
// The Addr returned is shared by all invocations of Addr, so
// do not modify it.
func (*UnixListener) Addr() Addr

// Close stops listening on the Unix address.  Already accepted
// connections are not closed.

// Close stops listening on the Unix address. Already accepted
// connections are not closed.
func (*UnixListener) Close() error

// File returns a copy of the underlying os.File, set to blocking
// mode.  It is the caller's responsibility to close f when finished.
// Closing l does not affect f, and closing f does not affect l.
//
// The returned os.File's file descriptor is different from the
// connection's.  Attempting to change properties of the original
// using this duplicate may or may not have the desired effect.

// File returns a copy of the underlying os.File, set to blocking
// mode. It is the caller's responsibility to close f when finished.
// Closing l does not affect f, and closing f does not affect l.
//
// The returned os.File's file descriptor is different from the
// connection's. Attempting to change properties of the original
// using this duplicate may or may not have the desired effect.
func (*UnixListener) File() (f *os.File, err error)

// SetDeadline sets the deadline associated with the listener.
// A zero time value disables the deadline.
func (*UnixListener) SetDeadline(t time.Time) error

func (Flags) String() string

func (HardwareAddr) String() string

// DefaultMask returns the default IP mask for the IP address ip.
// Only IPv4 addresses have default masks; DefaultMask returns
// nil if ip is not a valid IPv4 address.
func (IP) DefaultMask() IPMask

// Equal reports whether ip and x are the same IP address.
// An IPv4 address and that same address in IPv6 form are
// considered to be equal.
func (IP) Equal(x IP) bool

// IsGlobalUnicast reports whether ip is a global unicast
// address.
func (IP) IsGlobalUnicast() bool

// IsInterfaceLocalMulticast reports whether ip is
// an interface-local multicast address.
func (IP) IsInterfaceLocalMulticast() bool

// IsLinkLocalMulticast reports whether ip is a link-local
// multicast address.
func (IP) IsLinkLocalMulticast() bool

// IsLinkLocalUnicast reports whether ip is a link-local
// unicast address.
func (IP) IsLinkLocalUnicast() bool

// IsLoopback reports whether ip is a loopback address.
func (IP) IsLoopback() bool

// IsMulticast reports whether ip is a multicast address.
func (IP) IsMulticast() bool

// IsUnspecified reports whether ip is an unspecified address.
func (IP) IsUnspecified() bool

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (IP) MarshalText() ([]byte, error)

// Mask returns the result of masking the IP address ip with mask.
func (IP) Mask(mask IPMask) IP

// String returns the string form of the IP address ip.
// If the address is an IPv4 address, the string representation
// is dotted decimal ("74.125.19.99").  Otherwise the representation
// is IPv6 ("2001:4860:0:2001::68").

// String returns the string form of the IP address ip.
// It returns one of 4 forms:
//   - "<nil>", if ip has length 0
//   - dotted decimal ("192.0.2.1"), if ip is an IPv4 or IP4-mapped IPv6 address
//   - IPv6 ("2001:db8::1"), if ip is a valid IPv6 address
//   - the hexadecimal form of ip, without punctuation, if no other cases apply
func (IP) String() string

// To16 converts the IP address ip to a 16-byte representation.
// If ip is not an IP address (it is the wrong length), To16 returns nil.
func (IP) To16() IP

// To4 converts the IPv4 address ip to a 4-byte representation.
// If ip is not an IPv4 address, To4 returns nil.
func (IP) To4() IP

// Size returns the number of leading ones and total bits in the mask.
// If the mask is not in the canonical form--ones followed by zeros--then
// Size returns 0, 0.
func (IPMask) Size() (ones, bits int)

// String returns the hexadecimal form of m, with no punctuation.
func (IPMask) String() string

func (InvalidAddrError) Error() string

func (InvalidAddrError) Temporary() bool

func (InvalidAddrError) Timeout() bool

func (UnknownNetworkError) Error() string

func (UnknownNetworkError) Temporary() bool

func (UnknownNetworkError) Timeout() bool

