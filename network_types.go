package tasmota

import (
	"encoding/json"
	"net"
	"net/netip"
)

// IPAddr wraps netip.Addr to provide custom JSON marshaling for Tasmota IP addresses.
type IPAddr struct {
	netip.Addr
}

// MarshalJSON implements json.Marshaler for IPAddr.
func (ip IPAddr) MarshalJSON() ([]byte, error) {
	if !ip.IsValid() {
		return []byte(`""`), nil
	}
	return json.Marshal(ip.String())
}

// UnmarshalJSON implements json.Unmarshaler for IPAddr.
func (ip *IPAddr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" || s == "0.0.0.0" {
		ip.Addr = netip.Addr{}
		return nil
	}
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return err
	}
	ip.Addr = addr
	return nil
}

// String returns the string representation of the IP address.
func (ip IPAddr) String() string {
	if !ip.IsValid() {
		return ""
	}
	return ip.Addr.String()
}

// IsZero returns true if the IP address is not set.
func (ip IPAddr) IsZero() bool {
	return !ip.IsValid()
}

// NewIPAddr creates a new IPAddr from a string.
func NewIPAddr(s string) (IPAddr, error) {
	if s == "" || s == "0.0.0.0" {
		return IPAddr{}, nil
	}
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return IPAddr{}, err
	}
	return IPAddr{Addr: addr}, nil
}

// MustParseIPAddr parses an IP address or panics.
func MustParseIPAddr(s string) IPAddr {
	ip, err := NewIPAddr(s)
	if err != nil {
		panic(err)
	}
	return ip
}

// MACAddr wraps net.HardwareAddr to provide better typing for MAC addresses.
type MACAddr struct {
	net.HardwareAddr
}

// MarshalJSON implements json.Marshaler for MACAddr.
func (m MACAddr) MarshalJSON() ([]byte, error) {
	if len(m.HardwareAddr) == 0 {
		return []byte(`""`), nil
	}
	return json.Marshal(m.String())
}

// UnmarshalJSON implements json.Unmarshaler for MACAddr.
func (m *MACAddr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		m.HardwareAddr = nil
		return nil
	}
	mac, err := net.ParseMAC(s)
	if err != nil {
		return err
	}
	m.HardwareAddr = mac
	return nil
}

// String returns the string representation of the MAC address.
func (m MACAddr) String() string {
	if len(m.HardwareAddr) == 0 {
		return ""
	}
	return m.HardwareAddr.String()
}

// IsZero returns true if the MAC address is not set.
func (m MACAddr) IsZero() bool {
	return len(m.HardwareAddr) == 0
}

// NewMACAddr creates a new MACAddr from a string.
func NewMACAddr(s string) (MACAddr, error) {
	if s == "" {
		return MACAddr{}, nil
	}
	mac, err := net.ParseMAC(s)
	if err != nil {
		return MACAddr{}, err
	}
	return MACAddr{HardwareAddr: mac}, nil
}

// MustParseMACAddr parses a MAC address or panics.
func MustParseMACAddr(s string) MACAddr {
	mac, err := NewMACAddr(s)
	if err != nil {
		panic(err)
	}
	return mac
}
