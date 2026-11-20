// Command client calls the Orbital Wiki MCP server's tools from the command line.
// Each subcommand starts the server, calls one tool, prints the result and exits.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

const version = "0.1.0"

func main() {
	cmd := &cli.Command{
		Name:    "client",
		Usage:   "call the Orbital Wiki MCP server's tools",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "server",
				Usage:   "path to the server binary",
				Value:   "./bin/server",
				Sources: cli.EnvVars("ORBITALWIKI_MCP_SERVER"),
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "give up after this long",
				Value: 30 * time.Second,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "tools",
				Usage:  "list the server's tools",
				Action: listTools,
			},
			searchCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		if msg := err.Error(); msg != "" {
			fmt.Fprintln(os.Stderr, "client:", msg)
		}
		os.Exit(1)
	}
}
