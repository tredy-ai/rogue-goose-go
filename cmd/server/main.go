// Command server is the Orbital Wiki MCP server. It speaks MCP over stdio.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload" // loads .env into the environment
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v3"

	"github.com/tredy-ai/rogue-goose-go/internal/orbitalwiki"
	"github.com/tredy-ai/rogue-goose-go/internal/tools"
)

const version = "0.1.0"

func main() {
	cmd := &cli.Command{
		Name:    "server",
		Usage:   "Orbital Wiki MCP server (stdio)",
		Version: version,
		Description: "Speaks the Model Context Protocol over stdin/stdout. You normally don't run this\n" +
			"yourself: each `client` command starts it as a subprocess, calls one tool and exits.\n" +
			"It can also be registered with an MCP host such as Claude Code.\n\n" +
			"The API key is read from ORBITALWIKI_API_KEY. A .env file in the working directory\n" +
			"is loaded at startup; it does not override variables already set.\n\n" +
			"All logging goes to stderr; stdout carries the protocol.",
		// stdout is the protocol channel, so help and errors go to stderr.
		Writer:    os.Stderr,
		ErrWriter: os.Stderr,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "base-url",
				Usage:   "Orbital Wiki API base URL",
				Value:   orbitalwiki.DefaultBaseURL,
				Sources: cli.EnvVars("ORBITALWIKI_BASE_URL"),
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "debug, info, warn or error",
				Value: "info",
			},
		},
		Action: run,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "server:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	apiKey := os.Getenv("ORBITALWIKI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("no API key: set ORBITALWIKI_API_KEY in the environment or in .env")
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(cmd.String("log-level"))); err != nil {
		return fmt.Errorf("--log-level: %w", err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	api := orbitalwiki.NewClient(apiKey, orbitalwiki.WithBaseURL(cmd.String("base-url")))

	// Advertise the tools capability even before any tool is registered.
	opts := &mcp.ServerOptions{Capabilities: &mcp.ServerCapabilities{Tools: &mcp.ToolCapabilities{ListChanged: true}}}
	if level <= slog.LevelDebug {
		// The SDK logs every session's lifecycle at INFO, so only show it when debugging.
		opts.Logger = logger
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "orbitalwiki", Version: version}, opts)
	tools.Register(server, api)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Debug("serving on stdio", "base_url", api.BaseURL)
	return server.Run(ctx, &mcp.StdioTransport{})
}
