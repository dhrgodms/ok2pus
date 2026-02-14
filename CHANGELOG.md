# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com),
and this project adheres to [Semantic Versioning](https://semver.org).

## [1.0.0] - 2026-02-14

### Added
- Interactive terminal UI powered by [promptui](https://github.com/manifoldco/promptui)
- SSH host management: add, list, edit, delete, connect
- Password and Public Key authentication support
- Persistent host storage with SQLite (`~/.ok2pus/hosts.db`)
- Host configuration editing via system editor (`$EDITOR` or `vim`)
- Database reset and drop options with confirmation prompts
- Port validation (1–65535)
- Alias uniqueness constraint
- Required field validation for alias and host

[1.0.0]: https://github.com/dhrgodms/ok2pus/releases/tag/v1.0.0
