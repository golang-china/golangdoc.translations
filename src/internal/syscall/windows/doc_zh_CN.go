// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package windows // import "internal/syscall/windows"

import (
    "internal/syscall/windows/sysdll"
    "syscall"
    "unsafe"
)

const (
    ComputerNameNetBIOS                   = 0
    ComputerNameDnsHostname               = 1
    ComputerNameDnsDomain                 = 2
    ComputerNameDnsFullyQualified         = 3
    ComputerNamePhysicalNetBIOS           = 4
    ComputerNamePhysicalDnsHostname       = 5
    ComputerNamePhysicalDnsDomain         = 6
    ComputerNamePhysicalDnsFullyQualified = 7
    ComputerNameMax                       = 8

    MOVEFILE_REPLACE_EXISTING      = 0x1
    MOVEFILE_COPY_ALLOWED          = 0x2
    MOVEFILE_DELAY_UNTIL_REBOOT    = 0x4
    MOVEFILE_WRITE_THROUGH         = 0x8
    MOVEFILE_CREATE_HARDLINK       = 0x10
    MOVEFILE_FAIL_IF_NOT_TRACKABLE = 0x20
)

const GAA_FLAG_INCLUDE_PREFIX = 0x00000010

const (
    IF_TYPE_OTHER              = 1
    IF_TYPE_ETHERNET_CSMACD    = 6
    IF_TYPE_ISO88025_TOKENRING = 9
    IF_TYPE_PPP                = 23
    IF_TYPE_SOFTWARE_LOOPBACK  = 24
    IF_TYPE_ATM                = 37
    IF_TYPE_IEEE80211          = 71
    IF_TYPE_TUNNEL             = 131
    IF_TYPE_IEEE1394           = 144
)

const (
    IfOperStatusUp             = 1
    IfOperStatusDown           = 2
    IfOperStatusTesting        = 3
    IfOperStatusUnknown        = 4
    IfOperStatusDormant        = 5
    IfOperStatusNotPresent     = 6
    IfOperStatusLowerLayerDown = 7
)

type IpAdapterAddresses struct {
    Length                uint32
    IfIndex               uint32
    Next                  *IpAdapterAddresses
    AdapterName           *byte
    FirstUnicastAddress   *IpAdapterUnicastAddress
    FirstAnycastAddress   *IpAdapterAnycastAddress
    FirstMulticastAddress *IpAdapterMulticastAddress
    FirstDnsServerAddress *IpAdapterDnsServerAdapter
    DnsSuffix             *uint16
    Description           *uint16
    FriendlyName          *uint16
    PhysicalAddress       [syscall.MAX_ADAPTER_ADDRESS_LENGTH]byte
    PhysicalAddressLength uint32
    Flags                 uint32
    Mtu                   uint32
    IfType                uint32
    OperStatus            uint32
    Ipv6IfIndex           uint32
    ZoneIndices           [16]uint32
    FirstPrefix           *IpAdapterPrefix
}

type IpAdapterAnycastAddress struct {
    Length  uint32
    Flags   uint32
    Next    *IpAdapterAnycastAddress
    Address SocketAddress
}

type IpAdapterDnsServerAdapter struct {
    Length   uint32
    Reserved uint32
    Next     *IpAdapterDnsServerAdapter
    Address  SocketAddress
}

type IpAdapterMulticastAddress struct {
    Length  uint32
    Flags   uint32
    Next    *IpAdapterMulticastAddress
    Address SocketAddress
}

type IpAdapterPrefix struct {
    Length       uint32
    Flags        uint32
    Next         *IpAdapterPrefix
    Address      SocketAddress
    PrefixLength uint32
}

type IpAdapterUnicastAddress struct {
    Length             uint32
    Flags              uint32
    Next               *IpAdapterUnicastAddress
    Address            SocketAddress
    PrefixOrigin       int32
    SuffixOrigin       int32
    DadState           int32
    ValidLifetime      uint32
    PreferredLifetime  uint32
    LeaseLifetime      uint32
    OnLinkPrefixLength uint8
}

type SocketAddress struct {
    Sockaddr       *syscall.RawSockaddrAny
    SockaddrLength int32
}

func GetACP() (acp uint32)

func GetAdaptersAddresses(family uint32, flags uint32, reserved uintptr, adapterAddresses *IpAdapterAddresses, sizePointer *uint32) (errcode error)

func GetComputerNameEx(nameformat uint32, buf *uint16, n *uint32) (err error)

func MoveFileEx(from *uint16, to *uint16, flags uint32) (err error)

func MultiByteToWideChar(codePage uint32, dwFlags uint32, str *byte, nstr int32, wchar *uint16, nwchar int32) (nwrite int32, err error)

func Rename(oldpath, newpath string) error

