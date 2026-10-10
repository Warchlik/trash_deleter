# trash_deleter

[Polski](README_PL.md)

A small CLI tool written in Go that finds and removes junk files (`.DS_Store`, `Thumbs.db`, `*.tmp` and similar) from a directory tree.

> **Warning:** without `-dry-run` files are removed immediately. Run with `-dry-run` first to see what would be removed.

## Build

```sh
go build -o trash_deleter .
```

## Usage

```sh
./trash_deleter -path ./some/dir -dry-run
./trash_deleter -path ./some/dir
```

## Flags

| Flag | Description |
|---|---|
| `-path` | Directory to scan (default `.`) |
| `-dry-run` | Only list what would be removed, change nothing |
| `-patterns` | Extra comma-separated file name patterns, added to the defaults (e.g. `"*.log,backup_*"`) |
| `-ignore` | Extra comma-separated directory names to skip, added to the defaults (e.g. `"vendor,dist"`) |
| `-empty-dirs` | Also remove empty directories, including ones left empty after the cleanup |
| `-trash` | Move items to the system trash instead of deleting them permanently |
| `-max-depth` | How deep to scan: `1` = only files directly in `-path`, `0` = unlimited (default) |
| `-v`, `-verbose` | Print every removed file and directory |

## Defaults

- Patterns: `*:Zone.Identifier`, `.DS_Store`, `Thumbs.db`, `*.tmp`, `*~`
- Ignored directories: `.git`, `node_modules`

Patterns are matched against the file name only and use glob syntax. For a prefix use `backup_*`.

