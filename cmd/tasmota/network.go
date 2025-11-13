package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/kradalby/tasmota-go"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func newNetworkCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	return &ffcli.Command{
		Name:       "network",
		ShortUsage: "tasmota network <subcommand>",
		ShortHelp:  "Network configuration and diagnostics",
		LongHelp: `Configure network settings on Tasmota devices.

Network commands allow you to:
  - View current network configuration (IP, hostname, DHCP status)
  - Set device hostname
  - Configure static IP address or enable DHCP
  - Test connectivity by pinging hosts from the device

Note: Network changes may require a device restart to take full effect.
The device may also change IP addresses if switching between DHCP and static.

Examples:
  # View current network configuration
  tasmota --host 192.168.1.100 network get

  # Set hostname
  tasmota --host 192.168.1.100 network set-hostname --hostname tasmota-bedroom

  # Configure static IP
  tasmota --host 192.168.1.100 network set-static-ip \
    --ip 192.168.1.50 \
    --gateway 192.168.1.1 \
    --subnet 255.255.255.0 \
    --dns 8.8.8.8

  # Enable DHCP
  tasmota --host 192.168.1.100 network set-dhcp

  # Test connectivity
  tasmota --host 192.168.1.100 network ping --target 8.8.8.8`,
		Subcommands: []*ffcli.Command{
			newNetworkGetCmd(host, username, password, timeout, debug),
			newNetworkSetHostnameCmd(host, username, password, timeout, debug),
			newNetworkSetStaticIPCmd(host, username, password, timeout, debug),
			newNetworkSetDHCPCmd(host, username, password, timeout, debug),
			newNetworkPingCmd(host, username, password, timeout, debug),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

func newNetworkGetCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota network get", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "tasmota network get",
		ShortHelp:  "Get network configuration",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			config, err := client.GetNetworkConfig(ctx)
			if err != nil {
				return fmt.Errorf("failed to get network config: %w", err)
			}

			fmt.Printf("Network Configuration:\n")
			fmt.Printf("  Hostname: %s\n", config.Hostname)
			fmt.Printf("  IP Address: %s\n", config.IPAddress)
			fmt.Printf("  Gateway: %s\n", config.Gateway)
			fmt.Printf("  Subnet: %s\n", config.Subnet)
			fmt.Printf("  DNS Server: %s\n", config.DNSServer)
			fmt.Printf("  DHCP: %v\n", config.UseDHCP)

			return nil
		},
	}
}

func newNetworkSetHostnameCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota network set-hostname", flag.ExitOnError)
	hostname := fs.String("hostname", "", "New hostname (required)")

	return &ffcli.Command{
		Name:       "set-hostname",
		ShortUsage: "tasmota network set-hostname --hostname <name>",
		ShortHelp:  "Set device hostname",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if *hostname == "" {
				return fmt.Errorf("--hostname is required")
			}

			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.SetHostname(ctx, *hostname); err != nil {
				return fmt.Errorf("failed to set hostname: %w", err)
			}

			fmt.Printf("Hostname set to: %s\n", *hostname)
			return nil
		},
	}
}

func newNetworkSetStaticIPCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota network set-static-ip", flag.ExitOnError)
	ip := fs.String("ip", "", "IP address (required)")
	gateway := fs.String("gateway", "", "Gateway address (required)")
	subnet := fs.String("subnet", "", "Subnet mask (required)")
	dns := fs.String("dns", "", "DNS server (optional)")

	return &ffcli.Command{
		Name:       "set-static-ip",
		ShortUsage: "tasmota network set-static-ip --ip <ip> --gateway <gw> --subnet <mask> [--dns <dns>]",
		ShortHelp:  "Configure static IP address",
		LongHelp: `Configure a static IP address for the Tasmota device.

This disables DHCP and sets a fixed IP address. You must provide:
  - IP address for the device
  - Gateway address (usually your router)
  - Subnet mask (usually 255.255.255.0)

Optionally, you can specify a DNS server. If not provided, the device
will use the DNS server from DHCP or previous configuration.

Warning: If you set an incorrect IP configuration, you may lose network
connectivity to the device. Ensure your settings are correct before applying.

Example:
  tasmota --host 192.168.1.100 network set-static-ip \
    --ip 192.168.1.50 \
    --gateway 192.168.1.1 \
    --subnet 255.255.255.0 \
    --dns 8.8.8.8`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			if *ip == "" || *gateway == "" || *subnet == "" {
				return fmt.Errorf("--ip, --gateway, and --subnet are required")
			}

			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			ipAddr, err := tasmota.NewIPAddr(*ip)
			if err != nil {
				return fmt.Errorf("invalid IP address: %w", err)
			}

			gwAddr, err := tasmota.NewIPAddr(*gateway)
			if err != nil {
				return fmt.Errorf("invalid gateway address: %w", err)
			}

			subnetAddr, err := tasmota.NewIPAddr(*subnet)
			if err != nil {
				return fmt.Errorf("invalid subnet mask: %w", err)
			}

			if err := client.SetStaticIP(ctx, ipAddr, gwAddr, subnetAddr); err != nil {
				return fmt.Errorf("failed to set static IP: %w", err)
			}

			fmt.Printf("Static IP configured:\n")
			fmt.Printf("  IP: %s\n", *ip)
			fmt.Printf("  Gateway: %s\n", *gateway)
			fmt.Printf("  Subnet: %s\n", *subnet)

			if *dns != "" {
				dnsAddr, err := tasmota.NewIPAddr(*dns)
				if err != nil {
					return fmt.Errorf("invalid DNS server: %w", err)
				}
				if err := client.SetDNSServer(ctx, dnsAddr); err != nil {
					return fmt.Errorf("failed to set DNS server: %w", err)
				}
				fmt.Printf("  DNS: %s\n", *dns)
			}

			return nil
		},
	}
}

func newNetworkSetDHCPCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota network set-dhcp", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "set-dhcp",
		ShortUsage: "tasmota network set-dhcp",
		ShortHelp:  "Enable DHCP",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.EnableDHCP(ctx, true); err != nil {
				return fmt.Errorf("failed to enable DHCP: %w", err)
			}

			fmt.Println("DHCP enabled")
			return nil
		},
	}
}

func newNetworkPingCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota network ping", flag.ExitOnError)
	target := fs.String("target", "8.8.8.8", "Target host to ping")

	return &ffcli.Command{
		Name:       "ping",
		ShortUsage: "tasmota network ping [--target <host>]",
		ShortHelp:  "Ping a host from the device",
		LongHelp: `Test network connectivity by pinging a host from the Tasmota device.

This command instructs the device to ping a target host and reports whether
the ping was successful. This is useful for:
  - Verifying the device has internet connectivity
  - Testing if the device can reach specific hosts
  - Diagnosing network issues

The target can be an IP address or hostname. Default is 8.8.8.8 (Google DNS).

Examples:
  # Ping Google DNS (default)
  tasmota --host 192.168.1.100 network ping

  # Ping your router
  tasmota --host 192.168.1.100 network ping --target 192.168.1.1

  # Ping a hostname
  tasmota --host 192.168.1.100 network ping --target mqtt.home`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			success, err := client.Ping(ctx, *target)
			if err != nil {
				return fmt.Errorf("failed to ping: %w", err)
			}

			if success {
				fmt.Printf("Ping to %s: SUCCESS\n", *target)
			} else {
				fmt.Printf("Ping to %s: FAILED\n", *target)
			}

			return nil
		},
	}
}
