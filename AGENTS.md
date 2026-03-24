# nixgram Agent Conventions

## Project Structure & Operational Memory
- `cmd/nixgram/main.go` -> Application entrypoint (loads env vars, starts bot).
- `cmd/nixgram-hey/main.go` -> Example command meant to be run by the bot.
- `nixgram.go` -> Core bot lifecycle, initialization, and message dispatching.
- `runner.go` -> Execution logic, parsing bot commands and finding `$PATH` commands prefixed with `nixgram-`.
- `pkg/errreporter/reporter.go` -> (Planned/Future) Centralized error reporting.

## Documentation Requirements
- Use standard GoDoc-style line comments (`//`) for all exported types and functions.
- DO NOT use JSDoc or block comments (`/** ... */`) in Go code.
- Focus on the "why", constraints, edge cases, and nuances. No redundant comments.

## Error Handling
- All unexpected errors must be explicitly handled and routed through the centralized error-reporting function. Silent failures are strictly prohibited.
- (If no central reporter exists yet, ensure errors are at least properly logged with context, and consider implementing the central reporter in the shared utilities or services layer).

## Build & Tooling
- `mise` is the primary task runner.
- Dependencies and tool versions must be pinned to exact versions, no 'latest' or 'lts'.
- Pre-commit testing and verification is strictly enforced.
