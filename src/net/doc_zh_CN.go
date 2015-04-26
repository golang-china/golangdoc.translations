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

// net包提供了可移植的网络I/O接口，包括TCP/IP、UDP、域名解析和Unix域socket。
//
// 虽然本包提供了对网络原语的访问，大部分使用者只需要Dial、Listen和Accept函数提供的基本接口；以及相关的Conn和Listener接口。crypto/tls包提供了相同的接口和类似的Dial和Listen函数。
//
// Dial函数和服务端建立连接：
//
//	conn, err := net.Dial("tcp", "google.com:80")
//	if err != nil {
//		// handle error
//	}
//	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
//	status, err := bufio.NewReader(conn).ReadString('\n')
//	// ...
//
// Listen函数创建的服务端：
//
//	ln, err := net.Listen("tcp", ":8080")
//	if err != nil {
//		// handle error
//	}
//	for {
//		conn, err := ln.Accept()
//		if err != nil {
//			// handle error
//			continue
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

// 常用的IPv4地址。
//
//	var (
//	    IPv6zero                   = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
//	    IPv6unspecified            = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
//	    IPv6loopback               = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
//	    IPv6interfacelocalallnodes = IP{0xff, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
//	    IPv6linklocalallnodes      = IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
//	    IPv6linklocalallrouters    = IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02}
//	)
//
// 常用的IPv6地址。
//
//	var (
//	    ErrWriteToConnected = errors.New("use of WriteTo with pre-connected connection")
//	)
//
// 很多OpError类型的错误会包含本错误。
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

// InterfaceAddrs返回该系统的网络接口的地址列表。
func InterfaceAddrs() ([]Addr, error)

// Interfaces returns a list of the system's network interfaces.

// Interfaces返回该系统的网络接口列表。
func Interfaces() ([]Interface, error)

// JoinHostPort combines host and port into a network address of the form
// "host:port" or, if host contains a colon or a percent sign, "[host]:port".

// 函数将host和port合并为一个网络地址。一般格式为"host:port"；如果host含有冒号或百分号，格式为"[host]:port"。
func JoinHostPort(host, port string) string

// LookupAddr performs a reverse lookup for the given address, returning a list of
// names mapping to that address.

// LookupAddr查询某个地址，返回映射到该地址的主机名序列，本函数和LookupHost不互为反函数。
func LookupAddr(addr string) (name []string, err error)

// LookupCNAME returns the canonical DNS host for the given name. Callers that do
// not care about the canonical name can call LookupHost or LookupIP directly; both
// take care of resolving the canonical name as part of the lookup.

// LookupCNAME函数查询name的规范DNS名（但该域名未必可以访问）。如果调用者不关心规范名可以直接调用LookupHost或者LookupIP；这两个函数都会在查询时考虑到规范名。
func LookupCNAME(name string) (cname string, err error)

// LookupHost looks up the given host using the local resolver. It returns an array
// of that host's addresses.

// LookupHost函数查询主机的网络地址序列。
func LookupHost(host string) (addrs []string, err error)

// LookupIP looks up host using the local resolver. It returns an array of that
// host's IPv4 and IPv6 addresses.

// LookupIP函数查询主机的ipv4和ipv6地址序列。
func LookupIP(host string) (addrs []IP, err error)

// LookupMX returns the DNS MX records for the given domain name sorted by
// preference.

// LookupMX函数返回指定主机的按Pref字段排好序的DNS MX记录。
func LookupMX(name string) (mx []*MX, err error)

// LookupNS returns the DNS NS records for the given domain name.

// LookupNS函数返回指定主机的DNS NS记录。
func LookupNS(name string) (ns []*NS, err error)

// LookupPort looks up the port for the given network and service.

// LookupPort函数查询指定网络和服务的（默认）端口。
func LookupPort(network, service string) (port int, err error)

// LookupSRV tries to resolve an SRV query of the given service, protocol, and
// domain name. The proto is "tcp" or "udp". The returned records are sorted by
// priority and randomized by weight within a priority.
//
// LookupSRV constructs the DNS name to look up following RFC 2782. That is, it
// looks up _service._proto.name. To accommodate services publishing SRV records
// under non-standard names, if both service and proto are empty strings, LookupSRV
// looks up name directly.

// LookupSRV函数尝试执行指定服务、协议、主机的SRV查询。协议proto为"tcp"
// 或"udp"。返回的记录按Priority字段排序，同一优先度按Weight字段随机排序。
//
// LookupSRV函数按照RFC
// 2782的规定构建用于查询的DNS名。也就是说，它会查询_service._proto.name。为了适应将服务的SRV记录发布在非规范名下的情况，如果service和proto参数都是空字符串，函数会直接查询name。
func LookupSRV(service, proto, name string) (cname string, addrs []*SRV, err error)

// LookupTXT returns the DNS TXT records for the given domain name.

// LookupTXT函数返回指定主机的DNS TXT记录。
func LookupTXT(name string) (txt []string, err error)

// SplitHostPort splits a network address of the form "host:port", "[host]:port" or
// "[ipv6-host%zone]:port" into host or ipv6-host%zone and port. A literal address
// or host name for IPv6 must be enclosed in square brackets, as in "[::1]:80",
// "[ipv6-host]:http" or "[ipv6-host%zone]:80".

// 函数将格式为"host:port"、"[host]:port"或"[ipv6-host%zone]:port"的网络地址分割为host或ipv6-host%zone和port两个部分。Ipv6的文字地址或者主机名必须用方括号括起来，如"[::1]:80"、"[ipv6-host]:http"、"[ipv6-host%zone]:80"。
func SplitHostPort(hostport string) (host, port string, err error)

// Addr represents a network end point address.

// Addr代表一个网络终端地址。
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

// Conn接口代表通用的面向流的网络连接。多个线程可能会同时调用同一个Conn的方法。
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

// 在网络network上连接地址address，并返回一个Conn接口。可用的网络类型有：
//
// "tcp"、"tcp4"、"tcp6"、"udp"、"udp4"、"udp6"、"ip"、"ip4"、"ip6"、"unix"、"unixgram"、"unixpacket"
//
// 对TCP和UDP网络，地址格式是host:port或[host]:port，参见函数JoinHostPort和SplitHostPort。
//
//	Dial("tcp", "12.34.56.78:80")
//	Dial("tcp", "google.com:http")
//	Dial("tcp", "[2001:db8::1]:http")
//	Dial("tcp", "[fe80::1%lo0]:80")
//
// 对IP网络，network必须是"ip"、"ip4"、"ip6"后跟冒号和协议号或者协议名，地址必须是IP地址字面值。
//
//	Dial("ip4:1", "127.0.0.1")
//	Dial("ip6:ospf", "::1")
//
// 对Unix网络，地址必须是文件系统路径。
func Dial(network, address string) (Conn, error)

// DialTimeout acts like Dial but takes a timeout. The timeout includes name
// resolution, if required.

// DialTimeout类似Dial但采用了超时。timeout参数如果必要可包含名称解析。
func DialTimeout(network, address string, timeout time.Duration) (Conn, error)

// FileConn returns a copy of the network connection corresponding to the open file
// f. It is the caller's responsibility to close f when finished. Closing c does
// not affect f, and closing f does not affect c.

// FileConn返回一个下层为文件f的网络连接的拷贝。调用者有责任在结束程序前关闭f。关闭c不会影响f，关闭f也不会影响c。本函数与各种实现了Conn接口的类型的File方法是对应的。
func FileConn(f *os.File) (c Conn, err error)

// Pipe creates a synchronous, in-memory, full duplex network connection; both ends
// implement the Conn interface. Reads on one end are matched with writes on the
// other, copying data directly between the two; there is no internal buffering.

// Pipe创建一个内存中的同步、全双工网络连接。连接的两端都实现了Conn接口。一端的读取对应另一端的写入，直接将数据在两端之间作拷贝；没有内部缓冲。
func Pipe() (Conn, Conn)

// DNSConfigError represents an error reading the machine's DNS configuration.

// DNSConfigError代表读取主机DNS配置时出现的错误。
type DNSConfigError struct {
	Err error
}

func (e *DNSConfigError) Error() string

func (e *DNSConfigError) Temporary() bool

func (e *DNSConfigError) Timeout() bool

// DNSError represents a DNS lookup error.

// DNSError代表DNS查询的错误。
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

// Dialer类型包含与某个地址建立连接时的参数。
//
// 每一个字段的零值都等价于没有该字段。因此调用Dialer零值的Dial方法等价于调用Dial函数。
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

// Dial在指定的网络上连接指定的地址。参见Dial函数获取网络和地址参数的描述。
func (d *Dialer) Dial(network, address string) (Conn, error)

// An Error represents a network error.

// Error代表一个网络错误。
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

// HardwareAddr类型代表一个硬件地址（MAC地址）。
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

// ParseMAC函数使用如下格式解析一个IEEE 802 MAC-48、EUI-48或EUI-64硬件地址：
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

// IP类型是代表单个IP地址的[]byte切片。本包的函数都可以接受4字节（IPv4）和16字节（IPv6）的切片作为输入。
//
// 注意，IP地址是IPv4地址还是IPv6地址是语义上的属性，而不取决于切片的长度：16字节的切片也可以是IPv4地址。
type IP []byte

// IPv4 returns the IP address (in 16-byte form) of the IPv4 address a.b.c.d.

// IPv4返回包含一个IPv4地址a.b.c.d的IP地址（16字节格式）。
func IPv4(a, b, c, d byte) IP

// ParseCIDR parses s as a CIDR notation IP address and mask, like
// "192.168.100.1/24" or "2001:DB8::/48", as defined in RFC 4632 and RFC 4291.
//
// It returns the IP address and the network implied by the IP and mask. For
// example, ParseCIDR("192.168.100.1/16") returns the IP address 192.168.100.1 and
// the network 192.168.0.0/16.

// ParseCIDR将s作为一个CIDR（无类型域间路由）的IP地址和掩码字符串，如"192.168.100.1/24"或"2001:DB8::/48"，解析并返回IP地址和IP网络，参见RFC
// 4632和RFC 4291。
//
// 本函数会返回IP地址和该IP所在的网络和掩码。例如，ParseCIDR("192.168.100.1/16")会返回IP地址192.168.100.1和IP网络192.168.0.0/16。
func ParseCIDR(s string) (IP, *IPNet, error)

// ParseIP parses s as an IP address, returning the result. The string s can be in
// dotted decimal ("74.125.19.99") or IPv6 ("2001:4860:0:2001::68") form. If s is
// not a valid textual representation of an IP address, ParseIP returns nil.

// ParseIP将s解析为IP地址，并返回该地址。如果s不是合法的IP地址文本表示，ParseIP会返回nil。
//
// 字符串可以是小数点分隔的IPv4格式（如"74.125.19.99"）或IPv6格式（如"2001:4860:0:2001::68"）格式。
func ParseIP(s string) IP

// DefaultMask returns the default IP mask for the IP address ip. Only IPv4
// addresses have default masks; DefaultMask returns nil if ip is not a valid IPv4
// address.

// 函数返回IP地址ip的默认子网掩码。只有IPv4有默认子网掩码；如果ip不是合法的IPv4地址，会返回nil。
func (ip IP) DefaultMask() IPMask

// Equal returns true if ip and x are the same IP address. An IPv4 address and that
// same address in IPv6 form are considered to be equal.

// 如果ip和x代表同一个IP地址，Equal会返回真。代表同一地址的IPv4地址和IPv6地址也被认为是相等的。
func (ip IP) Equal(x IP) bool

// IsGlobalUnicast returns true if ip is a global unicast address.

// 如果ip是全局单播地址，则返回真。
func (ip IP) IsGlobalUnicast() bool

// IsInterfaceLinkLocalMulticast returns true if ip is an interface-local multicast
// address.

// 如果ip是接口本地组播地址，则返回真。
func (ip IP) IsInterfaceLocalMulticast() bool

// IsLinkLocalMulticast returns true if ip is a link-local multicast address.

// 如果ip是链路本地组播地址，则返回真。
func (ip IP) IsLinkLocalMulticast() bool

// IsLinkLocalUnicast returns true if ip is a link-local unicast address.

// 如果ip是链路本地单播地址，则返回真。
func (ip IP) IsLinkLocalUnicast() bool

// IsLoopback returns true if ip is a loopback address.

// 如果ip是环回地址，则返回真。
func (ip IP) IsLoopback() bool

// IsMulticast returns true if ip is a multicast address.

// 如果ip是组播地址，则返回真。
func (ip IP) IsMulticast() bool

// IsUnspecified returns true if ip is an unspecified address.

// 如果ip是未指定地址，则返回真。
func (ip IP) IsUnspecified() bool

// MarshalText implements the encoding.TextMarshaler interface. The encoding is the
// same as returned by String.

// MarshalText实现了encoding.TextMarshaler接口，返回值和String方法一样。
func (ip IP) MarshalText() ([]byte, error)

// Mask returns the result of masking the IP address ip with mask.

// Mask方法认为mask为ip的子网掩码，返回ip的网络地址部分的ip。（主机地址部分都置0）
func (ip IP) Mask(mask IPMask) IP

// String returns the string form of the IP address ip. If the address is an IPv4
// address, the string representation is dotted decimal ("74.125.19.99"). Otherwise
// the representation is IPv6 ("2001:4860:0:2001::68").

// String返回IP地址ip的字符串表示。如果ip是IPv4地址，返回值的格式为点分隔的，如"74.125.19.99"；否则表示为IPv6格式，如"2001:4860:0:2001::68"。
func (ip IP) String() string

// To16 converts the IP address ip to a 16-byte representation. If ip is not an IP
// address (it is the wrong length), To16 returns nil.

// To16将一个IP地址转换为16字节表示。如果ip不是一个IP地址（长度错误），To16会返回nil。
func (ip IP) To16() IP

// To4 converts the IPv4 address ip to a 4-byte representation. If ip is not an
// IPv4 address, To4 returns nil.

// To4将一个IPv4地址转换为4字节表示。如果ip不是IPv4地址，To4会返回nil。
func (ip IP) To4() IP

// UnmarshalText implements the encoding.TextUnmarshaler interface. The IP address
// is expected in a form accepted by ParseIP.

// UnmarshalText实现了encoding.TextUnmarshaler接口。IP地址字符串应该是ParseIP函数可以接受的格式。
func (ip *IP) UnmarshalText(text []byte) error

// IPAddr represents the address of an IP end point.

// IPAddr代表一个IP终端的地址。
type IPAddr struct {
	IP   IP
	Zone string // IPv6 scoped addressing zone
}

// ResolveIPAddr parses addr as an IP address of the form "host" or
// "ipv6-host%zone" and resolves the domain name on the network net, which must be
// "ip", "ip4" or "ip6".

// ResolveIPAddr将addr作为一个格式为"host"或"ipv6-host%zone"的IP地址来解析。
// 函数会在参数net指定的网络类型上解析，net必须是"ip"、"ip4"或"ip6"。
func ResolveIPAddr(net, addr string) (*IPAddr, error)

// Network returns the address's network name, "ip".

// Network返回地址的网络类型："ip"。
func (a *IPAddr) Network() string

func (a *IPAddr) String() string

// IPConn is the implementation of the Conn and PacketConn interfaces for IP
// network connections.

// IPConn类型代表IP网络连接，实现了Conn和PacketConn接口。
type IPConn struct {
	// contains filtered or unexported fields
}

// DialIP connects to the remote address raddr on the network protocol netProto,
// which must be "ip", "ip4", or "ip6" followed by a colon and a protocol number or
// name.

// DialIP在网络协议netProto上连接本地地址laddr和远端地址raddr，netProto必须是"ip"、"ip4"或"ip6"后跟冒号和协议名或协议号。
func DialIP(netProto string, laddr, raddr *IPAddr) (*IPConn, error)

// ListenIP listens for incoming IP packets addressed to the local address laddr.
// The returned connection's ReadFrom and WriteTo methods can be used to receive
// and send IP packets with per-packet addressing.

// ListenIP创建一个接收目的地是本地地址laddr的IP数据包的网络连接，返回的*IPConn的ReadFrom和WriteTo方法可以用来发送和接收IP数据包。（每个包都可获取来源址或者设置目标地址）
func ListenIP(netProto string, laddr *IPAddr) (*IPConn, error)

// Close closes the connection.

// Close关闭连接
func (c *IPConn) Close() error

// File sets the underlying os.File to blocking mode and returns a copy. It is the
// caller's responsibility to close f when finished. Closing c does not affect f,
// and closing f does not affect c.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.

// File方法设置下层的os.File为阻塞模式并返回其副本。
//
// 使用者有责任在用完后关闭f。关闭c不影响f，关闭f也不影响c。返回的os.File类型文件描述符和原本的网络连接是不同的。试图使用该副本修改本体的属性可能会（也可能不会）得到期望的效果。
func (c *IPConn) File() (f *os.File, err error)

// LocalAddr returns the local network address.

// LocalAddr返回本地网络地址
func (c *IPConn) LocalAddr() Addr

// Read implements the Conn Read method.

// Read实现Conn接口Read方法
func (c *IPConn) Read(b []byte) (int, error)

// ReadFrom implements the PacketConn ReadFrom method.

// ReadFrom实现PacketConn接口ReadFrom方法。注意本方法有bug，应避免使用。
func (c *IPConn) ReadFrom(b []byte) (int, Addr, error)

// ReadFromIP reads an IP packet from c, copying the payload into b. It returns the
// number of bytes copied into b and the return address that was on the packet.
//
// ReadFromIP can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.

// ReadFromIP从c读取一个IP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
//
// ReadFromIP方法会在超过一个固定的时间点之后超时，并返回一个错误。注意本方法有bug，应避免使用。
func (c *IPConn) ReadFromIP(b []byte) (int, *IPAddr, error)

// ReadMsgIP reads a packet from c, copying the payload into b and the associated
// out-of-band data into oob. It returns the number of bytes copied into b, the
// number of bytes copied into oob, the flags that were set on the packet and the
// source address of the packet.

// ReadMsgIP从c读取一个数据包，将有效负载拷贝进b，相关的带外数据拷贝进oob，返回拷贝进b的字节数，拷贝进oob的字节数，数据包的flag，数据包来源地址和可能的错误。
func (c *IPConn) ReadMsgIP(b, oob []byte) (n, oobn, flags int, addr *IPAddr, err error)

// RemoteAddr returns the remote network address.

// RemoteAddr返回远端网络地址
func (c *IPConn) RemoteAddr() Addr

// SetDeadline implements the Conn SetDeadline method.

// SetDeadline设置读写操作绝对期限，实现了Conn接口的SetDeadline方法
func (c *IPConn) SetDeadline(t time.Time) error

// SetReadBuffer sets the size of the operating system's receive buffer associated
// with the connection.

// SetReadBuffer设置该连接的系统接收缓冲
func (c *IPConn) SetReadBuffer(bytes int) error

// SetReadDeadline implements the Conn SetReadDeadline method.

// SetReadDeadline设置读操作绝对期限，实现了Conn接口的SetReadDeadline方法
func (c *IPConn) SetReadDeadline(t time.Time) error

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the connection.

// SetWriteBuffer设置该连接的系统发送缓冲
func (c *IPConn) SetWriteBuffer(bytes int) error

// SetWriteDeadline implements the Conn SetWriteDeadline method.

// SetWriteDeadline设置写操作绝对期限，实现了Conn接口的SetWriteDeadline方法
func (c *IPConn) SetWriteDeadline(t time.Time) error

// Write implements the Conn Write method.

// Write实现Conn接口Write方法
func (c *IPConn) Write(b []byte) (int, error)

// WriteMsgIP writes a packet to addr via c, copying the payload from b and the
// associated out-of-band data from oob. It returns the number of payload and
// out-of-band bytes written.

// WriteMsgIP通过c向地址addr发送一个数据包，b和oob分别为包有效负载和对应的带外数据，返回写入的字节数（包数据、带外数据）和可能的错误。
func (c *IPConn) WriteMsgIP(b, oob []byte, addr *IPAddr) (n, oobn int, err error)

// WriteTo implements the PacketConn WriteTo method.

// WriteTo实现PacketConn接口WriteTo方法
func (c *IPConn) WriteTo(b []byte, addr Addr) (int, error)

// WriteToIP writes an IP packet to addr via c, copying the payload from b.
//
// WriteToIP can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline. On
// packet-oriented connections, write timeouts are rare.

// WriteToIP通过c向地址addr发送一个数据包，b为包的有效负载，返回写入的字节。
//
// WriteToIP方法会在超过一个固定的时间点之后超时，并返回一个错误。在面向数据包的连接上，写入超时是十分罕见的。
func (c *IPConn) WriteToIP(b []byte, addr *IPAddr) (int, error)

// An IP mask is an IP address.

// IPMask代表一个IP地址的掩码。
type IPMask []byte

// CIDRMask returns an IPMask consisting of `ones' 1 bits followed by 0s up to a
// total length of `bits' bits. For a mask of this form, CIDRMask is the inverse of
// IPMask.Size.

// CIDRMask返回一个IPMask类型值，该返回值总共有bits个字位，其中前ones个字位都是1，其余字位都是0。
func CIDRMask(ones, bits int) IPMask

// IPv4Mask returns the IP mask (in 4-byte form) of the IPv4 mask a.b.c.d.

// IPv4Mask返回一个4字节格式的IPv4掩码a.b.c.d。
func IPv4Mask(a, b, c, d byte) IPMask

// Size returns the number of leading ones and total bits in the mask. If the mask
// is not in the canonical form--ones followed by zeros--then Size returns 0, 0.

// Size返回m的前导的1字位数和总字位数。如果m不是规范的子网掩码（字位：/^1+0+$/），将返会(0, 0)。
func (m IPMask) Size() (ones, bits int)

// String returns the hexadecimal form of m, with no punctuation.

// String返回m的十六进制格式，没有标点。
func (m IPMask) String() string

// An IPNet represents an IP network.

// IPNet表示一个IP网络。
type IPNet struct {
	IP   IP     // network number
	Mask IPMask // network mask
}

// Contains reports whether the network includes ip.

// Contains报告该网络是否包含地址ip。
func (n *IPNet) Contains(ip IP) bool

// Network returns the address's network name, "ip+net".

// Network返回网络类型名："ip+net"，注意该类型名是不合法的。
func (n *IPNet) Network() string

// String returns the CIDR notation of n like "192.168.100.1/24" or "2001:DB8::/48"
// as defined in RFC 4632 and RFC 4291. If the mask is not in the canonical form,
// it returns the string which consists of an IP address, followed by a slash
// character and a mask expressed as hexadecimal form with no punctuation like
// "192.168.100.1/c000ff00".

// String返回n的CIDR表示，如"192.168.100.1/24"或"2001:DB8::/48"，参见RFC 4632和RFC
// 4291。如果n的Mask字段不是规范格式，它会返回一个包含n.IP.String()、斜线、n.Mask.String()（此时表示为无标点十六进制格式）的字符串，如"192.168.100.1/c000ff00"。
func (n *IPNet) String() string

// Interface represents a mapping between network interface name and index. It also
// represents network interface facility information.

// Interface类型代表一个网络接口（系统与网络的一个接点）。包含接口索引到名字的映射，也包含接口的设备信息。
type Interface struct {
	Index        int          // positive integer that starts at one, zero is never used
	MTU          int          // maximum transmission unit
	Name         string       // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr HardwareAddr // IEEE MAC-48, EUI-48 and EUI-64 form
	Flags        Flags        // e.g., FlagUp, FlagLoopback, FlagMulticast
}

// InterfaceByIndex returns the interface specified by index.

// InterfaceByIndex返回指定索引的网络接口。
func InterfaceByIndex(index int) (*Interface, error)

// InterfaceByName returns the interface specified by name.

// InterfaceByName返回指定名字的网络接口。
func InterfaceByName(name string) (*Interface, error)

// Addrs returns interface addresses for a specific interface.

// Addrs方法返回网络接口ifi的一或多个接口地址。
func (ifi *Interface) Addrs() ([]Addr, error)

// MulticastAddrs returns multicast, joined group addresses for a specific
// interface.

// MulticastAddrs返回网络接口ifi加入的多播组地址。
func (ifi *Interface) MulticastAddrs() ([]Addr, error)

type InvalidAddrError string

func (e InvalidAddrError) Error() string

func (e InvalidAddrError) Temporary() bool

func (e InvalidAddrError) Timeout() bool

// A Listener is a generic network listener for stream-oriented protocols.
//
// Multiple goroutines may invoke methods on a Listener simultaneously.

// Listener是一个用于面向流的网络协议的公用的网络监听器接口。多个线程可能会同时调用一个Listener的方法。
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

// FileListener返回一个下层为文件f的网络监听器的拷贝。调用者有责任在使用结束后改变l。关闭l不会影响f，关闭f也不会影响l。本函数与各种实现了Listener接口的类型的File方法是对应的。
func FileListener(f *os.File) (l Listener, err error)

// Listen announces on the local network address laddr. The network net must be a
// stream-oriented network: "tcp", "tcp4", "tcp6", "unix" or "unixpacket". See Dial
// for the syntax of laddr.
func Listen(net, laddr string) (Listener, error)

// An MX represents a single DNS MX record.

// MX代表一条DNS
// MX记录（邮件交换记录），根据收信人的地址后缀来定位邮件服务器。
type MX struct {
	Host string
	Pref uint16
}

// An NS represents a single DNS NS record.

// NS代表一条DNS
// NS记录（域名服务器记录），指定该域名由哪个DNS服务器来进行解析。
type NS struct {
	Host string
}

// OpError is the error type usually returned by functions in the net package. It
// describes the operation, network type, and address of an error.

// OpError是经常被net包的函数返回的错误类型。它描述了该错误的操作、网络类型和网络地址。
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

// PacketConn接口代表通用的面向数据包的网络连接。多个线程可能会同时调用同一个Conn的方法。
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

// FilePacketConn函数返回一个下层为文件f的数据包网络连接的拷贝。调用者有责任在结束程序前关闭f。关闭c不会影响f，关闭f也不会影响c。本函数与各种实现了PacketConn接口的类型的File方法是对应的。
func FilePacketConn(f *os.File) (c PacketConn, err error)

// ListenPacket announces on the local network address laddr. The network net must
// be a packet-oriented network: "udp", "udp4", "udp6", "ip", "ip4", "ip6" or
// "unixgram". See Dial for the syntax of laddr.

// ListenPacket函数监听本地网络地址laddr。网络类型net必须是面向数据包的网络类型：
//
// "ip"、"ip4"、"ip6"、"udp"、"udp4"、"udp6"、或"unixgram"。laddr的格式参见Dial函数。
func ListenPacket(net, laddr string) (PacketConn, error)

// A ParseError represents a malformed text string and the type of string that was
// expected.

// ParseError代表一个格式错误的字符串，Type为期望的格式。
type ParseError struct {
	Type string
	Text string
}

func (e *ParseError) Error() string

// An SRV represents a single DNS SRV record.

// SRV代表一条DNS
// SRV记录（资源记录），记录某个服务由哪台计算机提供。
type SRV struct {
	Target   string
	Port     uint16
	Priority uint16
	Weight   uint16
}

// TCPAddr represents the address of a TCP end point.

// TCPAddr代表一个TCP终端地址。
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

// ResolveTCPAddr将addr作为TCP地址解析并返回。参数addr格式为"host:port"或"[ipv6-host%zone]:port"，解析得到网络名和端口名；net必须是"tcp"、"tcp4"或"tcp6"。
//
// IPv6地址字面值/名称必须用方括号包起来，如"[::1]:80"、"[ipv6-host]:http"或"[ipv6-host%zone]:80"。
func ResolveTCPAddr(net, addr string) (*TCPAddr, error)

// Network returns the address's network name, "tcp".

// 返回地址的网络类型，"tcp"。
func (a *TCPAddr) Network() string

func (a *TCPAddr) String() string

// TCPConn is an implementation of the Conn interface for TCP network connections.

// TCPConn代表一个TCP网络连接，实现了Conn接口。
type TCPConn struct {
	// contains filtered or unexported fields
}

// DialTCP connects to the remote address raddr on the network net, which must be
// "tcp", "tcp4", or "tcp6". If laddr is not nil, it is used as the local address
// for the connection.

// DialTCP在网络协议net上连接本地地址laddr和远端地址raddr。net必须是"tcp"、"tcp4"、"tcp6"；如果laddr不是nil，将使用它作为本地地址，否则自动选择一个本地地址。
func DialTCP(net string, laddr, raddr *TCPAddr) (*TCPConn, error)

// Close closes the connection.

// Close关闭连接
func (c *TCPConn) Close() error

// CloseRead shuts down the reading side of the TCP connection. Most callers should
// just use Close.

// CloseRead关闭TCP连接的读取侧（以后不能读取），应尽量使用Close方法。
func (c *TCPConn) CloseRead() error

// CloseWrite shuts down the writing side of the TCP connection. Most callers
// should just use Close.

// CloseWrite关闭TCP连接的写入侧（以后不能写入），应尽量使用Close方法。
func (c *TCPConn) CloseWrite() error

// File sets the underlying os.File to blocking mode and returns a copy. It is the
// caller's responsibility to close f when finished. Closing c does not affect f,
// and closing f does not affect c.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.

// File方法设置下层的os.File为阻塞模式并返回其副本。
//
// 使用者有责任在用完后关闭f。关闭c不影响f，关闭f也不影响c。返回的os.File类型文件描述符和原本的网络连接是不同的。试图使用该副本修改本体的属性可能会（也可能不会）得到期望的效果。
func (c *TCPConn) File() (f *os.File, err error)

// LocalAddr returns the local network address.

// LocalAddr返回本地网络地址
func (c *TCPConn) LocalAddr() Addr

// Read implements the Conn Read method.

// Read实现了Conn接口Read方法
func (c *TCPConn) Read(b []byte) (int, error)

// ReadFrom implements the io.ReaderFrom ReadFrom method.

// ReadFrom实现了io.ReaderFrom接口的ReadFrom方法
func (c *TCPConn) ReadFrom(r io.Reader) (int64, error)

// RemoteAddr returns the remote network address.

// RemoteAddr返回远端网络地址
func (c *TCPConn) RemoteAddr() Addr

// SetDeadline implements the Conn SetDeadline method.

// SetDeadline设置读写操作期限，实现了Conn接口的SetDeadline方法
func (c *TCPConn) SetDeadline(t time.Time) error

// SetKeepAlive sets whether the operating system should send keepalive messages on
// the connection.

// SetKeepAlive设置操作系统是否应该在该连接中发送keepalive信息
func (c *TCPConn) SetKeepAlive(keepalive bool) error

// SetKeepAlivePeriod sets period between keep alives.

// SetKeepAlivePeriod设置keepalive的周期，超出会断开
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

// SetLinger设定当连接中仍有数据等待发送或接受时的Close方法的行为。
//
// 如果sec <
// 0（默认），Close方法立即返回，操作系统停止后台数据发送；如果 sec ==
// 0，Close立刻返回，操作系统丢弃任何未发送或未接收的数据；如果sec >
// 0，Close方法阻塞最多sec秒，等待数据发送或者接收，在一些操作系统中，在超时后，任何未发送的数据会被丢弃。
func (c *TCPConn) SetLinger(sec int) error

// SetNoDelay controls whether the operating system should delay packet
// transmission in hopes of sending fewer packets (Nagle's algorithm). The default
// is true (no delay), meaning that data is sent as soon as possible after a Write.

// SetNoDelay设定操作系统是否应该延迟数据包传递，以便发送更少的数据包（Nagle's算法）。默认为真，即数据应该在Write方法后立刻发送。
func (c *TCPConn) SetNoDelay(noDelay bool) error

// SetReadBuffer sets the size of the operating system's receive buffer associated
// with the connection.

// SetReadBuffer设置该连接的系统接收缓冲
func (c *TCPConn) SetReadBuffer(bytes int) error

// SetReadDeadline implements the Conn SetReadDeadline method.

// SetReadDeadline设置读操作期限，实现了Conn接口的SetReadDeadline方法
func (c *TCPConn) SetReadDeadline(t time.Time) error

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the connection.

// SetWriteBuffer设置该连接的系统发送缓冲
func (c *TCPConn) SetWriteBuffer(bytes int) error

// SetWriteDeadline implements the Conn SetWriteDeadline method.

// SetWriteDeadline设置写操作期限，实现了Conn接口的SetWriteDeadline方法
func (c *TCPConn) SetWriteDeadline(t time.Time) error

// Write implements the Conn Write method.

// Write实现了Conn接口Write方法
func (c *TCPConn) Write(b []byte) (int, error)

// TCPListener is a TCP network listener. Clients should typically use variables of
// type Listener instead of assuming TCP.

// TCPListener代表一个TCP网络的监听者。使用者应尽量使用Listener接口而不是假设（网络连接为）TCP。
type TCPListener struct {
	// contains filtered or unexported fields
}

// ListenTCP announces on the TCP address laddr and returns a TCP listener. Net
// must be "tcp", "tcp4", or "tcp6". If laddr has a port of 0, ListenTCP will
// choose an available port. The caller can use the Addr method of TCPListener to
// retrieve the chosen address.

// ListenTCP在本地TCP地址laddr上声明并返回一个*TCPListener，net参数必须是"tcp"、"tcp4"、"tcp6"，如果laddr的端口字段为0，函数将选择一个当前可用的端口，可以用Listener的Addr方法获得该端口。
func ListenTCP(net string, laddr *TCPAddr) (*TCPListener, error)

// Accept implements the Accept method in the Listener interface; it waits for the
// next call and returns a generic Conn.

// Accept用于实现Listener接口的Accept方法；他会等待下一个呼叫，并返回一个该呼叫的Conn接口。
func (l *TCPListener) Accept() (Conn, error)

// AcceptTCP accepts the next incoming call and returns the new connection.

// AcceptTCP接收下一个呼叫，并返回一个新的*TCPConn。
func (l *TCPListener) AcceptTCP() (*TCPConn, error)

// Addr returns the listener's network address, a *TCPAddr.

// Addr返回l监听的的网络地址，一个*TCPAddr。
func (l *TCPListener) Addr() Addr

// Close stops listening on the TCP address. Already Accepted connections are not
// closed.

// Close停止监听TCP地址，已经接收的连接不受影响。
func (l *TCPListener) Close() error

// File returns a copy of the underlying os.File, set to blocking mode. It is the
// caller's responsibility to close f when finished. Closing l does not affect f,
// and closing f does not affect l.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.

// File方法返回下层的os.File的副本，并将该副本设置为阻塞模式。
//
// 使用者有责任在用完后关闭f。关闭c不影响f，关闭f也不影响c。返回的os.File类型文件描述符和原本的网络连接是不同的。试图使用该副本修改本体的属性可能会（也可能不会）得到期望的效果。
func (l *TCPListener) File() (f *os.File, err error)

// SetDeadline sets the deadline associated with the listener. A zero time value
// disables the deadline.

// 设置监听器执行的期限，t为Time零值则会关闭期限限制。
func (l *TCPListener) SetDeadline(t time.Time) error

// UDPAddr represents the address of a UDP end point.

// UDPAddr代表一个UDP终端地址。
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

// ResolveTCPAddr将addr作为TCP地址解析并返回。参数addr格式为"host:port"或"[ipv6-host%zone]:port"，解析得到网络名和端口名；net必须是"udp"、"udp4"或"udp6"。
//
// IPv6地址字面值/名称必须用方括号包起来，如"[::1]:80"、"[ipv6-host]:http"或"[ipv6-host%zone]:80"。
func ResolveUDPAddr(net, addr string) (*UDPAddr, error)

// Network returns the address's network name, "udp".

// 返回地址的网络类型，"udp"。
func (a *UDPAddr) Network() string

func (a *UDPAddr) String() string

// UDPConn is the implementation of the Conn and PacketConn interfaces for UDP
// network connections.

// UDPConn代表一个UDP网络连接，实现了Conn和PacketConn接口。
type UDPConn struct {
	// contains filtered or unexported fields
}

// DialUDP connects to the remote address raddr on the network net, which must be
// "udp", "udp4", or "udp6". If laddr is not nil, it is used as the local address
// for the connection.

// DialTCP在网络协议net上连接本地地址laddr和远端地址raddr。net必须是"udp"、"udp4"、"udp6"；如果laddr不是nil，将使用它作为本地地址，否则自动选择一个本地地址。
func DialUDP(net string, laddr, raddr *UDPAddr) (*UDPConn, error)

// ListenMulticastUDP listens for incoming multicast UDP packets addressed to the
// group address gaddr on ifi, which specifies the interface to join.
// ListenMulticastUDP uses default multicast interface if ifi is nil.

// ListenMulticastUDP接收目的地是ifi接口上的组地址gaddr的UDP数据包。它指定了使用的接口，如果ifi是nil，将使用默认接口。
func ListenMulticastUDP(net string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error)

// ListenUDP listens for incoming UDP packets addressed to the local address laddr.
// Net must be "udp", "udp4", or "udp6". If laddr has a port of 0, ListenUDP will
// choose an available port. The LocalAddr method of the returned UDPConn can be
// used to discover the port. The returned connection's ReadFrom and WriteTo
// methods can be used to receive and send UDP packets with per-packet addressing.

// ListenUDP创建一个接收目的地是本地地址laddr的UDP数据包的网络连接。net必须是"udp"、"udp4"、"udp6"；如果laddr端口为0，函数将选择一个当前可用的端口，可以用Listener的Addr方法获得该端口。返回的*UDPConn的ReadFrom和WriteTo方法可以用来发送和接收UDP数据包（每个包都可获得来源地址或设置目标地址）。
func ListenUDP(net string, laddr *UDPAddr) (*UDPConn, error)

// Close closes the connection.

// Close关闭连接
func (c *UDPConn) Close() error

// File sets the underlying os.File to blocking mode and returns a copy. It is the
// caller's responsibility to close f when finished. Closing c does not affect f,
// and closing f does not affect c.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.

// File方法设置下层的os.File为阻塞模式并返回其副本。
//
// 使用者有责任在用完后关闭f。关闭c不影响f，关闭f也不影响c。返回的os.File类型文件描述符和原本的网络连接是不同的。试图使用该副本修改本体的属性可能会（也可能不会）得到期望的效果。
func (c *UDPConn) File() (f *os.File, err error)

// LocalAddr returns the local network address.

// LocalAddr返回本地网络地址
func (c *UDPConn) LocalAddr() Addr

// Read implements the Conn Read method.

// Read实现Conn接口Read方法
func (c *UDPConn) Read(b []byte) (int, error)

// ReadFrom implements the PacketConn ReadFrom method.

// ReadFrom实现PacketConn接口ReadFrom方法
func (c *UDPConn) ReadFrom(b []byte) (int, Addr, error)

// ReadFromUDP reads a UDP packet from c, copying the payload into b. It returns
// the number of bytes copied into b and the return address that was on the packet.
//
// ReadFromUDP can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.

// ReadFromUDP从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
//
// ReadFromUDP方法会在超过一个固定的时间点之后超时，并返回一个错误。
func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err error)

// ReadMsgUDP reads a packet from c, copying the payload into b and the associated
// out-of-band data into oob. It returns the number of bytes copied into b, the
// number of bytes copied into oob, the flags that were set on the packet and the
// source address of the packet.

// ReadMsgUDP从c读取一个数据包，将有效负载拷贝进b，相关的带外数据拷贝进oob，返回拷贝进b的字节数，拷贝进oob的字节数，数据包的flag，数据包来源地址和可能的错误。
func (c *UDPConn) ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *UDPAddr, err error)

// RemoteAddr returns the remote network address.

// RemoteAddr返回远端网络地址
func (c *UDPConn) RemoteAddr() Addr

// SetDeadline implements the Conn SetDeadline method.

// SetDeadline设置读写操作期限，实现了Conn接口的SetDeadline方法
func (c *UDPConn) SetDeadline(t time.Time) error

// SetReadBuffer sets the size of the operating system's receive buffer associated
// with the connection.

// SetReadBuffer设置该连接的系统接收缓冲
func (c *UDPConn) SetReadBuffer(bytes int) error

// SetReadDeadline implements the Conn SetReadDeadline method.

// SetReadDeadline设置读操作期限，实现了Conn接口的SetReadDeadline方法
func (c *UDPConn) SetReadDeadline(t time.Time) error

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the connection.

// SetWriteBuffer设置该连接的系统发送缓冲
func (c *UDPConn) SetWriteBuffer(bytes int) error

// SetWriteDeadline implements the Conn SetWriteDeadline method.

// SetWriteDeadline设置写操作期限，实现了Conn接口的SetWriteDeadline方法
func (c *UDPConn) SetWriteDeadline(t time.Time) error

// Write implements the Conn Write method.

// Write实现Conn接口Write方法
func (c *UDPConn) Write(b []byte) (int, error)

// WriteMsgUDP writes a packet to addr via c, copying the payload from b and the
// associated out-of-band data from oob. It returns the number of payload and
// out-of-band bytes written.

// WriteMsgUDP通过c向地址addr发送一个数据包，b和oob分别为包有效负载和对应的带外数据，返回写入的字节数（包数据、带外数据）和可能的错误。
func (c *UDPConn) WriteMsgUDP(b, oob []byte, addr *UDPAddr) (n, oobn int, err error)

// WriteTo implements the PacketConn WriteTo method.

// WriteTo实现PacketConn接口WriteTo方法
func (c *UDPConn) WriteTo(b []byte, addr Addr) (int, error)

// WriteToUDP writes a UDP packet to addr via c, copying the payload from b.
//
// WriteToUDP can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline. On
// packet-oriented connections, write timeouts are rare.

// WriteToUDP通过c向地址addr发送一个数据包，b为包的有效负载，返回写入的字节。
//
// WriteToUDP方法会在超过一个固定的时间点之后超时，并返回一个错误。在面向数据包的连接上，写入超时是十分罕见的。
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)

// UnixAddr represents the address of a Unix domain socket end point.

// UnixAddr代表一个Unix域socket终端地址。
type UnixAddr struct {
	Name string
	Net  string
}

// ResolveUnixAddr parses addr as a Unix domain socket address. The string net
// gives the network name, "unix", "unixgram" or "unixpacket".

// ResolveUnixAddr将addr作为Unix域socket地址解析，参数net指定网络类型："unix"、"unixgram"或"unixpacket"。
func ResolveUnixAddr(net, addr string) (*UnixAddr, error)

// Network returns the address's network name, "unix", "unixgram" or "unixpacket".

// 返回地址的网络类型，"unix"，"unixgram"或"unixpacket"。
func (a *UnixAddr) Network() string

func (a *UnixAddr) String() string

// UnixConn is an implementation of the Conn interface for connections to Unix
// domain sockets.

// UnixConn代表Unix域socket连接，实现了Conn和PacketConn接口。
type UnixConn struct {
	// contains filtered or unexported fields
}

// DialUnix connects to the remote address raddr on the network net, which must be
// "unix", "unixgram" or "unixpacket". If laddr is not nil, it is used as the local
// address for the connection.

// DialUnix在网络协议net上连接本地地址laddr和远端地址raddr。net必须是"unix"、"unixgram"、"unixpacket"，如果laddr不是nil将使用它作为本地地址，否则自动选择一个本地地址。
func DialUnix(net string, laddr, raddr *UnixAddr) (*UnixConn, error)

// ListenUnixgram listens for incoming Unix datagram packets addressed to the local
// address laddr. The network net must be "unixgram". The returned connection's
// ReadFrom and WriteTo methods can be used to receive and send packets with
// per-packet addressing.

// ListenUnixgram接收目的地是本地地址laddr的Unix
// datagram网络连接。net必须是"unixgram"，返回的*UnixConn的ReadFrom和WriteTo方法可以用来发送和接收数据包（每个包都可获取来源址或者设置目标地址）。
func ListenUnixgram(net string, laddr *UnixAddr) (*UnixConn, error)

// Close closes the connection.

// Close关闭连接
func (c *UnixConn) Close() error

// CloseRead shuts down the reading side of the Unix domain connection. Most
// callers should just use Close.

// CloseRead关闭TCP连接的读取侧（以后不能读取），应尽量使用Close方法
func (c *UnixConn) CloseRead() error

// CloseWrite shuts down the writing side of the Unix domain connection. Most
// callers should just use Close.

// CloseWrite关闭TCP连接的写入侧（以后不能写入），应尽量使用Close方法
func (c *UnixConn) CloseWrite() error

// File sets the underlying os.File to blocking mode and returns a copy. It is the
// caller's responsibility to close f when finished. Closing c does not affect f,
// and closing f does not affect c.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.

// File方法设置下层的os.File为阻塞模式并返回其副本。
//
// 使用者有责任在用完后关闭f。关闭c不影响f，关闭f也不影响c。返回的os.File类型文件描述符和原本的网络连接是不同的。试图使用该副本修改本体的属性可能会（也可能不会）得到期望的效果。
func (c *UnixConn) File() (f *os.File, err error)

// LocalAddr returns the local network address.

// LocalAddr返回本地网络地址
func (c *UnixConn) LocalAddr() Addr

// Read implements the Conn Read method.

// Read实现了Conn接口Read方法
func (c *UnixConn) Read(b []byte) (int, error)

// ReadFrom implements the PacketConn ReadFrom method.

// ReadFrom实现PacketConn接口ReadFrom方法
func (c *UnixConn) ReadFrom(b []byte) (int, Addr, error)

// ReadFromUnix reads a packet from c, copying the payload into b. It returns the
// number of bytes copied into b and the source address of the packet.
//
// ReadFromUnix can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.

// ReadFromUnix从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
//
// ReadFromUnix方法会在超过一个固定的时间点之后超时，并返回一个错误。
func (c *UnixConn) ReadFromUnix(b []byte) (int, *UnixAddr, error)

// ReadMsgUnix reads a packet from c, copying the payload into b and the associated
// out-of-band data into oob. It returns the number of bytes copied into b, the
// number of bytes copied into oob, the flags that were set on the packet, and the
// source address of the packet.

// ReadMsgUnix从c读取一个数据包，将有效负载拷贝进b，相关的带外数据拷贝进oob，返回拷贝进b的字节数，拷贝进oob的字节数，数据包的flag，数据包来源地址和可能的错误。
func (c *UnixConn) ReadMsgUnix(b, oob []byte) (n, oobn, flags int, addr *UnixAddr, err error)

// RemoteAddr returns the remote network address.

// RemoteAddr返回远端网络地址
func (c *UnixConn) RemoteAddr() Addr

// SetDeadline implements the Conn SetDeadline method.

// SetDeadline设置读写操作期限，实现了Conn接口的SetDeadline方法
func (c *UnixConn) SetDeadline(t time.Time) error

// SetReadBuffer sets the size of the operating system's receive buffer associated
// with the connection.

// SetReadBuffer设置该连接的系统接收缓冲
func (c *UnixConn) SetReadBuffer(bytes int) error

// SetReadDeadline implements the Conn SetReadDeadline method.

// SetReadDeadline设置读操作期限，实现了Conn接口的SetReadDeadline方法
func (c *UnixConn) SetReadDeadline(t time.Time) error

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the connection.

// SetWriteBuffer设置该连接的系统发送缓冲
func (c *UnixConn) SetWriteBuffer(bytes int) error

// SetWriteDeadline implements the Conn SetWriteDeadline method.

// SetWriteDeadline设置写操作期限，实现了Conn接口的SetWriteDeadline方法
func (c *UnixConn) SetWriteDeadline(t time.Time) error

// Write implements the Conn Write method.

// Write实现了Conn接口Write方法
func (c *UnixConn) Write(b []byte) (int, error)

// WriteMsgUnix writes a packet to addr via c, copying the payload from b and the
// associated out-of-band data from oob. It returns the number of payload and
// out-of-band bytes written.

// WriteMsgUnix通过c向地址addr发送一个数据包，b和oob分别为包有效负载和对应的带外数据，返回写入的字节数（包数据、带外数据）和可能的错误。
func (c *UnixConn) WriteMsgUnix(b, oob []byte, addr *UnixAddr) (n, oobn int, err error)

// WriteTo implements the PacketConn WriteTo method.

// WriteTo实现PacketConn接口WriteTo方法
func (c *UnixConn) WriteTo(b []byte, addr Addr) (int, error)

// WriteToUnix writes a packet to addr via c, copying the payload from b.
//
// WriteToUnix can be made to time out and return an error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline. On
// packet-oriented connections, write timeouts are rare.

// WriteToUnix通过c向地址addr发送一个数据包，b为包的有效负载，返回写入的字节。
//
// WriteToUnix方法会在超过一个固定的时间点之后超时，并返回一个错误。在面向数据包的连接上，写入超时是十分罕见的。
func (c *UnixConn) WriteToUnix(b []byte, addr *UnixAddr) (int, error)

// UnixListener is a Unix domain socket listener. Clients should typically use
// variables of type Listener instead of assuming Unix domain sockets.

// UnixListener代表一个Unix域scoket的监听者。使用者应尽量使用Listener接口而不是假设（网络连接为）Unix域scoket。
type UnixListener struct {
	// contains filtered or unexported fields
}

// ListenUnix announces on the Unix domain socket laddr and returns a Unix
// listener. The network net must be "unix" or "unixpacket".

// ListenTCP在Unix域scoket地址laddr上声明并返回一个*UnixListener，net参数必须是"unix"或"unixpacket"。
func ListenUnix(net string, laddr *UnixAddr) (*UnixListener, error)

// Accept implements the Accept method in the Listener interface; it waits for the
// next call and returns a generic Conn.

// Accept用于实现Listener接口的Accept方法；他会等待下一个呼叫，并返回一个该呼叫的Conn接口。
func (l *UnixListener) Accept() (Conn, error)

// AcceptUnix accepts the next incoming call and returns the new connection.

// AcceptUnix接收下一个呼叫，并返回一个新的*UnixConn。
func (l *UnixListener) AcceptUnix() (*UnixConn, error)

// Addr returns the listener's network address.

// Addr返回l的监听的Unix域socket地址
func (l *UnixListener) Addr() Addr

// Close stops listening on the Unix address. Already accepted connections are not
// closed.

// Close停止监听Unix域socket地址，已经接收的连接不受影响。
func (l *UnixListener) Close() error

// File returns a copy of the underlying os.File, set to blocking mode. It is the
// caller's responsibility to close f when finished. Closing l does not affect f,
// and closing f does not affect l.
//
// The returned os.File's file descriptor is different from the connection's.
// Attempting to change properties of the original using this duplicate may or may
// not have the desired effect.

// File方法返回下层的os.File的副本，并将该副本设置为阻塞模式。
//
// 使用者有责任在用完后关闭f。关闭c不影响f，关闭f也不影响c。返回的os.File类型文件描述符和原本的网络连接是不同的。试图使用该副本修改本体的属性可能会（也可能不会）得到期望的效果。
func (l *UnixListener) File() (*os.File, error)

// SetDeadline sets the deadline associated with the listener. A zero time value
// disables the deadline.

// 设置监听器执行的期限，t为Time零值则会关闭期限限制
func (l *UnixListener) SetDeadline(t time.Time) error

type UnknownNetworkError string

func (e UnknownNetworkError) Error() string

func (e UnknownNetworkError) Temporary() bool

func (e UnknownNetworkError) Timeout() bool
