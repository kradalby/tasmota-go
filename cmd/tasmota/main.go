package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kradalby/tasmota-go"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	var (
		rootFlagSet = flag.NewFlagSet("tasmota", flag.ExitOnError)
		host        = rootFlagSet.String("host", "", "Tasmota device host/IP (required)")
		username    = rootFlagSet.String("username", "", "Basic auth username")
		password    = rootFlagSet.String("password", "", "Basic auth password")
		timeout     = rootFlagSet.Duration("timeout", 10*time.Second, "Request timeout")
		debug       = rootFlagSet.Bool("debug", false, "Enable debug logging")
	)

	root := &ffcli.Command{
		Name:       "tasmota",
		ShortUsage: "tasmota [flags] <subcommand>",
		ShortHelp:  "CLI tool for controlling Tasmota devices",
		FlagSet:    rootFlagSet,
		Subcommands: []*ffcli.Command{
			newStatusCmd(host, username, password, timeout, debug),
			newPowerCmd(host, username, password, timeout, debug),
			newInfoCmd(host, username, password, timeout, debug),
			newNetworkCmd(host, username, password, timeout, debug),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func newClient(host, username, password string, timeout time.Duration, debug bool) (*tasmota.Client, error) {
	if host == "" {
		return nil, fmt.Errorf("--host is required")
	}

	opts := []tasmota.ClientOption{
		tasmota.WithTimeout(timeout),
	}

	if username != "" || password != "" {
		opts = append(opts, tasmota.WithAuth(username, password))
	}

	if debug {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		opts = append(opts, tasmota.WithLogger(logger))
	}

	return tasmota.NewClient(host, opts...)
}
