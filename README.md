# Orbital Wiki take-home (Go)

This repo is a small [MCP](https://modelcontextprotocol.io) server and a CLI client, wrapped around the [Orbital Wiki API](https://www.orbitalwiki.com/developers). The server, the client and the CLI plumbing are already built. Your job is the tools themselves.

There are two parts:

- **Part 1:** done **by hand, without AI**.
- **Part 2:** done **with AI assistance**. We want to see how you work with it.

Expect this to take about **2 hours**. If you're going well past that, stop and write down where you got to and what you'd do next. That's a perfectly fine submission.

There's also a TypeScript track (a Next.js app around the same API). Do whichever you prefer, or both if you like.

---

## Setup

1. **Create your repo.** Click "Use this template" and create a **private** repository. Add **jalavosus and ToFut** as a collaborator.
2. **Get an API key.** Create a free account at [orbitalwiki.com](https://www.orbitalwiki.com/developers/quickstart) and generate a key in the dashboard. The key is shown only once, so save it. The free tier allows 1,000 requests a day, which is plenty. Without a key, you share 50 requests a day with everyone on your network.
3. **Configure:**
   ```sh
   cp .env.example .env    # paste your key in
   set -a; source .env; set +a
   ```
   Or just `export ORBITALWIKI_API_KEY=ow_...` in your shell.
4. **Build and run:**
   ```sh
   make build
   ./bin/client tools                # lists the server's tools
   ./bin/client search starlink      # fails until you finish Part 1
   ```

The API reference is `docs/openapi.json` (a copy of https://www.orbitalwiki.com/api/v1/openapi.json). The MCP Go SDK docs are at https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk/mcp.

### How the pieces fit

- `bin/server` is an MCP server that talks over **stdin/stdout**. You don't run it yourself: each `client` command starts it, calls one tool, and exits.
  - Because stdout carries the protocol, **server code must never print to stdout**. Log with `slog`, which is already set up to write to stderr.
- `bin/client` is the tool caller. Each subcommand maps to one tool.
  - Required parameters are positional arguments; optional parameters are flags.
  - `--server` points at a different server binary.
- `internal/orbitalwiki/` is the API client: types and HTTP.
- `internal/tools/` holds the MCP tool definitions.
- `docs/tools/` has one doc per tool.

### Commands

| What | Plain Go | Make |
|---|---|---|
| Build both binaries | `go build -o bin/server ./cmd/server && go build -o bin/client ./cmd/client` | `make build` |
| Run all tests | `go test ./...` | `make test` |
| Tests with the race detector | `go test -race ./...` | `make test-race` |
| One package, verbose | `go test -v ./internal/orbitalwiki/` | |
| One test | `go test -run TestParseStatus ./internal/orbitalwiki/` | |
| Vet | `go vet ./...` | `make vet` |
| Lint | `golangci-lint run` | `make lint` |
| Format | `gofmt -w .` | `make fmt` |
| Rebuild on save | `go tool air` | `make watch` |
| Run the server directly (debugging) | `go run ./cmd/server --log-level debug` | |

**air** ([air-verse/air](https://github.com/air-verse/air)) is installed as a Go tool dependency, so you don't need to install anything. Run `go tool air` in one terminal and it rebuilds `bin/server` and `bin/client` every time you save a `.go` file. Then run `./bin/client …` in another terminal; it always starts the freshly built server. Config is in `.air.toml`.

**Optional: try your server from an MCP host.** For example, in Claude Code (check `claude mcp add --help` if the syntax has changed):

```sh
claude mcp add orbitalwiki -e ORBITALWIKI_API_KEY=$ORBITALWIKI_API_KEY -- "$(pwd)/bin/server"
```

This is handy in Part 2: an agent can call the tool it's helping you build. **Don't use it in Part 1.**

---

## Rules on AI use

**Part 1: AI does not write your code.** That means:

- no chat assistants (ChatGPT, Claude, Copilot Chat, etc.)
- no coding agents
- no multi-line AI completions that fill in function bodies

Ordinary editor autocomplete (gopls, short single-line suggestions) is fine. So are the OpenAPI spec, pkg.go.dev, the Go and MCP docs, Stack Overflow and any other documentation. If your editor has multi-line AI completion (Copilot, Cursor Tab, etc.), turn it off before you start.

**Part 2: use whatever AI tools you like.**

After you submit, we'll have a short call (about 30 minutes). You'll walk us through your Part 1 code, then make a small change to it live, without AI. That call is how we check Part 1. We're not going to try to guess from your commits.

---

## Part 1: the `search_satellites` tool (no AI)

`docs/tools/search.md` is the spec: what the tool takes, what it returns, and how it fails. `client search` is already wired up to call it.

### 1. `internal/orbitalwiki/search.go`

- Define `SearchParams` and `SearchResponse`.
- Implement `Client.Search` with `net/http`. That means:
  - build the URL and query string
  - set the auth header
  - send the request
  - decode the success body
  - turn error responses into useful Go errors (the API's error body is described in `docs/openapi.json`)
- The enums in `enums.go` are done for you. Use them.

### 2. `internal/tools/search.go`

Register the `search_satellites` tool with the MCP Go SDK. That covers:

- the tool's name, description and input schema
- validating input
- calling `Client.Search`
- turning the response, or the error, into a tool result

Working out how to do this from the SDK docs is a big part of the exercise.

### 3. Tests

Write testify table tests for your code. `internal/orbitalwiki/enums_test.go` shows the style. Good things to test:

- query building
- input validation
- response and error decoding

Use `httptest.Server`; tests must not call the real API. They don't need to be exhaustive; we just want to see how you test.

**You're done with Part 1 when:**

- `client search` works as `docs/tools/search.md` describes, including its error cases
- `go vet ./...`, `go test ./...` and `golangci-lint run` all pass

Try a few:

```sh
./bin/client search starlink --limit 5
./bin/client search --status operational --orbit-class geo
./bin/client search --status sideways      # should be a clear tool error
./bin/client search                        # so should this
```

Commit as you go; small commits are welcome. When you're finished:

```sh
git tag part-1 && git push origin part-1
```

---

## Part 2: a tool of your choosing (with AI)

Pick any other endpoint in `docs/openapi.json` and build a tool for it. The API is your oyster. Tip: avoid the `/dataset/*` endpoints. They're bulk downloads limited to 20 a day.

1. **Add an MCP tool** and register it in `internal/tools/register.go`.
2. **Add a `client` subcommand** that calls it. Required parameters are positional arguments, optional ones are flags. Reuse `callTool` in `cmd/client/mcp.go`.
3. **Document it** in `docs/tools/<tool>.md`, in the same shape as `search.md`:
   - which endpoint it calls
   - why the tool is useful
   - inputs, output and errors
   - an example invocation
4. **Write tests** (testify, `httptest`) for the new code.

The tool must work as documented.

### Logging your AI use

Everything goes in `ai-log/`:

- **Raw transcripts or exports** from the tools you used, unedited. For example:
  - Claude Code: the `.jsonl` files under `~/.claude/projects/<this-project>/`
  - Cursor: export the chat
  - ChatGPT / Claude.ai: share links or copied conversations
- **Specs, plans and brainstorms**, whether you wrote them or the AI did.
- **`ai-log/NOTES.md`**, a few paragraphs covering:
  - what you asked for and why
  - what the AI got wrong, or what you rejected
  - how you checked its work

`AGENTS.md` (which `CLAUDE.md` also loads) tells coding agents to save their plans into `ai-log/` and to stay out of the Part 1 code. You're free to change it, but if you do, tell us why in your notes.

When you're finished:

```sh
git tag part-2 && git push origin part-2
```

---

## Submitting

Push both tags, then let **Jamie Alavosus (jamie@tredy.ai)** know. Anything you want us to know goes in `ai-log/NOTES.md` or the message.

## What we look at

- **Part 1:**
  - correctness
  - input validation and error handling
  - how you used the MCP SDK
  - your tests
  - a readable commit history
- **Part 2:**
  - whether the tool works as documented
  - the quality of the doc
  - how you directed the AI (specs, prompts, iteration)
  - what you caught and fixed
  - your tests
