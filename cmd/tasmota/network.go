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
		ShortHelp:  "Network configuration",
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
		ShortUsage: "tasmota network set-static-ip --ip <ip> --gateway <gw> --subnet <mask>",
		ShortHelp:  "Configure static IP",
		FlagSet:    fs,
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
		FlagSet:    fs,
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
