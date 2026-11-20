package main

import (
	"context"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/tredy-ai/rogue-goose-go/internal/orbitalwiki"
)

func searchCommand() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "search satellites (search_satellites tool)",
		UsageText: "client search [options] [query]",
		Arguments: []cli.Argument{
			&cli.StringArg{Name: "query", UsageText: "free-text search: name, alias, NORAD ID, COSPAR ID, operator or country"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "orbit-class", Usage: "one of " + strings.Join(orbitalwiki.OrbitClassValues(), ", ")},
			&cli.StringFlag{Name: "status", Usage: "one of " + strings.Join(orbitalwiki.StatusValues(), ", ")},
			&cli.StringFlag{Name: "country", Usage: "country code"},
			&cli.IntFlag{Name: "limit", Usage: "max results, 1-100 (API default 20)", HideDefault: true},
			&cli.IntFlag{Name: "offset", Usage: "results to skip"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Pass everything through as given: the tool validates, not the client.
			args := map[string]any{}
			if q := cmd.StringArg("query"); q != "" {
				args["q"] = q
			}
			for flag, key := range map[string]string{"orbit-class": "orbit_class", "status": "status", "country": "country"} {
				if cmd.IsSet(flag) {
					args[key] = cmd.String(flag)
				}
			}
			for _, flag := range []string{"limit", "offset"} {
				if cmd.IsSet(flag) {
					args[flag] = cmd.Int(flag)
				}
			}
			return callTool(ctx, cmd, "search_satellites", args)
		},
	}
}
