# Tool: `search_satellites`

**Calls:** Orbital Wiki `GET /search` (see `docs/openapi.json`).

**Why:** finds satellites by name, alias, NORAD ID, COSPAR ID, operator or country, optionally filtered by orbit class and operational status. It's the usual starting point before looking up a single satellite in detail.

## Input

| Argument | Type | Notes |
|---|---|---|
| `q` | string | Free-text search. |
| `orbit_class` | string | One of `LEO MEO GEO HEO SSO VLEO UNKNOWN`. Case-insensitive. |
| `status` | string | One of `OPERATIONAL NONOPERATIONAL PARTIALLY_OPERATIONAL STANDBY SPARE EXTENDED_MISSION DECAYED UNKNOWN`. Case-insensitive. |
| `country` | string | Country code. |
| `limit` | integer | 1–100. Defaults to the API's default (20) when omitted. |
| `offset` | integer | ≥ 0. |

At least one of `q`, `orbit_class`, `status` or `country` is required.

The tool's input schema must describe these arguments well enough for an LLM to call it correctly without reading this file.

## Output

The matching satellites and the total match count. The exact format is up to the implementer, but it has to work both for an LLM reading it and for a person running `client search`. At minimum each result includes NORAD ID, name, orbit class, operational status, country code, and what it matched on.

## Errors

- **Invalid input:** a tool error whose message names the bad argument and, for enums, the allowed values. Validate before calling the API.
- **API errors:** a tool error that includes the HTTP status and the API's error `code`.
- The API key never appears in any output.

## Example

    client search starlink --status operational --limit 5
