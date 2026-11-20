package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v3"
)

// connect starts the server binary named by --server as a subprocess and returns an
// initialised MCP session. The server inherits this process's environment, and its
// stderr is wired to ours. The caller must Close the session.
func connect(ctx context.Context, cmd *cli.Command) (*mcp.ClientSession, error) {
	server := exec.CommandContext(ctx, cmd.String("server"))
	server.Stderr = os.Stderr

	client := mcp.NewClient(&mcp.Implementation{Name: "orbitalwiki-client", Version: version}, nil)
	cs, err := client.Connect(ctx, &mcp.CommandTransport{Command: server}, nil)
	if err != nil {
		return nil, fmt.Errorf("starting server %q: %w", cmd.String("server"), err)
	}
	return cs, nil
}

// callTool starts the server, calls one tool with args, and prints the result.
// Text content goes to stdout. If the result is a tool error, its text goes to
// stderr and the command exits with status 1.
func callTool(ctx context.Context, cmd *cli.Command, name string, args map[string]any) error {
	ctx, cancel := context.WithTimeout(ctx, cmd.Duration("timeout"))
	defer cancel()
	cs, err := connect(ctx, cmd)
	if err != nil {
		return err
	}
	defer func() { _ = cs.Close() }()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		return fmt.Errorf("calling %s: %w", name, err)
	}

	out := cmd.Root().Writer
	if res.IsError {
		out = cmd.Root().ErrWriter
	}

	printed := false
	for _, c := range res.Content {
		if t, ok := c.(*mcp.TextContent); ok {
			_, _ = fmt.Fprintln(out, t.Text)
			printed = true
		}
	}
	if !printed && res.StructuredContent != nil {
		b, err := json.MarshalIndent(res.StructuredContent, "", "  ")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(out, string(b))
	}

	if res.IsError {
		return cli.Exit("", 1)
	}
	return nil
}

// listTools prints each of the server's tools with its name, description and input schema.
func listTools(ctx context.Context, cmd *cli.Command) error {
	ctx, cancel := context.WithTimeout(ctx, cmd.Duration("timeout"))
	defer cancel()
	cs, err := connect(ctx, cmd)
	if err != nil {
		return err
	}
	defer func() { _ = cs.Close() }()

	out := cmd.Root().Writer
	for tool, err := range cs.Tools(ctx, nil) {
		if err != nil {
			return fmt.Errorf("listing tools: %w", err)
		}
		b, err := json.MarshalIndent(map[string]any{
			"name":        tool.Name,
			"description": tool.Description,
			"inputSchema": tool.InputSchema,
		}, "", "  ")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(out, string(b))
	}
	return nil
}
