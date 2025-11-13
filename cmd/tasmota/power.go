package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/kradalby/tasmota-go"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func newPowerCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	return &ffcli.Command{
		Name:       "power",
		ShortUsage: "tasmota power <subcommand>",
		ShortHelp:  "Control device power (relays)",
		LongHelp: `Control power relays on Tasmota devices.

Tasmota devices can have up to 8 relays. Use --relay to specify which relay
to control (1-8). Relay 1 is the default. Use --relay 0 to control all relays.

Examples:
  # Turn on relay 1 (default)
  tasmota --host 192.168.1.100 power on

  # Turn off relay 2
  tasmota --host 192.168.1.100 power off --relay 2

  # Toggle relay 3
  tasmota --host 192.168.1.100 power toggle --relay 3

  # Get current state of relay 1
  tasmota --host 192.168.1.100 power get

  # Turn on all relays
  tasmota --host 192.168.1.100 power on --relay 0`,
		Subcommands: []*ffcli.Command{
			newPowerOnCmd(host, username, password, timeout, debug),
			newPowerOffCmd(host, username, password, timeout, debug),
			newPowerToggleCmd(host, username, password, timeout, debug),
			newPowerGetCmd(host, username, password, timeout, debug),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

func newPowerOnCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota power on", flag.ExitOnError)
	relay := fs.Int("relay", 1, "Relay number (1-8, 0=all)")

	return &ffcli.Command{
		Name:       "on",
		ShortUsage: "tasmota power on [--relay N]",
		ShortHelp:  "Turn power on",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			var resp *tasmota.PowerResponse
			if *relay == 0 || *relay == 1 {
				resp, err = client.Power(ctx, tasmota.PowerOn)
			} else {
				resp, err = client.PowerN(ctx, *relay, tasmota.PowerOn)
			}
			if err != nil {
				return fmt.Errorf("failed to turn on: %w", err)
			}

			fmt.Printf("Power relay %d turned ON (%s)\n", *relay, resp.Power)
			return nil
		},
	}
}

func newPowerOffCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota power off", flag.ExitOnError)
	relay := fs.Int("relay", 1, "Relay number (1-8, 0=all)")

	return &ffcli.Command{
		Name:       "off",
		ShortUsage: "tasmota power off [--relay N]",
		ShortHelp:  "Turn power off",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			var resp *tasmota.PowerResponse
			if *relay == 0 || *relay == 1 {
				resp, err = client.Power(ctx, tasmota.PowerOff)
			} else {
				resp, err = client.PowerN(ctx, *relay, tasmota.PowerOff)
			}
			if err != nil {
				return fmt.Errorf("failed to turn off: %w", err)
			}

			fmt.Printf("Power relay %d turned OFF (%s)\n", *relay, resp.Power)
			return nil
		},
	}
}

func newPowerToggleCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota power toggle", flag.ExitOnError)
	relay := fs.Int("relay", 1, "Relay number (1-8, 0=all)")

	return &ffcli.Command{
		Name:       "toggle",
		ShortUsage: "tasmota power toggle [--relay N]",
		ShortHelp:  "Toggle power state",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			var resp *tasmota.PowerResponse
			if *relay == 0 || *relay == 1 {
				resp, err = client.Power(ctx, tasmota.PowerToggle)
			} else {
				resp, err = client.PowerN(ctx, *relay, tasmota.PowerToggle)
			}
			if err != nil {
				return fmt.Errorf("failed to toggle: %w", err)
			}

			fmt.Printf("Power relay %d toggled (%s)\n", *relay, resp.Power)
			return nil
		},
	}
}

func newPowerGetCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota power get", flag.ExitOnError)
	relay := fs.Int("relay", 1, "Relay number (1-8, 0=all)")

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "tasmota power get [--relay N]",
		ShortHelp:  "Get current power state",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			var resp *tasmota.PowerResponse
			if *relay == 0 || *relay == 1 {
				resp, err = client.GetPower(ctx)
			} else {
				resp, err = client.GetPowerN(ctx, *relay)
			}
			if err != nil {
				return fmt.Errorf("failed to get power state: %w", err)
			}

			fmt.Printf("Power relay %d: %s\n", *relay, resp.Power)
			return nil
		},
	}
}
