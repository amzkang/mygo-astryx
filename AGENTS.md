# AGENTS.md

## What this is

`mygo-astryx` is a native-UI design system, component library, multi-theme token engine, and agent-ready CLI/MCP toolkit built on [MyGo](https://github.com/egoist/mygo) v0.8.3 (inspired by [Meta Astryx](https://github.com/meta/astryx)).

## Verification commands

Run after every change:

```sh
./bin/go build ./... && ./bin/go vet ./... && ./bin/go tool mygo vet .
./bin/go test ./...
```

## MyGo Native UI Rules (`mygo vet`)

- **Elements do not outlive their frame.** Never store `ui.Element`, `*ui.Context`, or `*Parts` in structs, package variables, or closures sent to goroutines.
- **The UI thread never waits.** Run I/O or background work in goroutines and publish back via `Window.Update` / `Window.Invalidate`.
- **Theme scopes stay in-frame.** Use `theme.WithTheme(c, th, fn)` (which wraps `ui.Themed` and pushes onto a frame-local stack via `ui.Local`) and read tokens with `theme.Of(c)`.
- **Identity is stable.** Repeated or stateful controls must use `c.Key(...)` before construction.
- **Swizzle-safe component files.** Exported component files in `ds/` must keep imports limited to `github.com/egoist/mygo/ui`, `github.com/amzkang/mygo-astryx/theme`, `github.com/amzkang/mygo-astryx/icons`, and standard library packages so `astryx swizzle` can rewrite and eject them cleanly.
