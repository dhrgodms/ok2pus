# ok2pus

```
            :::::::::       ::::::::::::
            ::+@@%+::     -::-+%@@@@@#=::
 ::::::::: -:-#@@%=::::::::-*@@@@%%@@@@+:::::::::::::: :::::::::::::::::::::::::::
:::-+*%%#*=-:::=%@@#-:-+*#*=-*@@@=:::-@@@#::+##*-=*%%*=-::-+##+-:::+*#*-::-=+#%##+=:::
::-#@@@@@@@@@=::+@@@*-+@@@%=::::::::::*@@@=::@@@@@@@@@@@#::=%@@#-::-%@@#=:+@@@@@@@@@#-:
:-@@@%-::-@@@%:-#@@@*@@@%=::-    ::-+@@@@-::-@@@@-::-%@@@=:+@@@*:::=@@@*-=@@@#-::=*++::
:%@@@::::-%@@%-=%@@@@@@#-::   :::-*@@@%-::::*@@@=:::-*@@@=-*@@%=:::*@@@+::#@@@@@@@*-:::
:@@@%::::+@@@*:=@@@@@@@@*-:::::=%@@@#-::::::%@@@::::=@@@%:=%@@#-::-#@@%=:::-=+*%@@@@+::
:@@@@+:-*@@@%--*@@@=-%@@@*-:-*@@@@%++++++-:-@@@@#--+@@@@-:=@@@@*=*@@@@#-=%@@*-:-#@@@=::
:-#@@@@@@@%=::=#@@#-:-#@@@*-+@@@@@@@@@@@@-:*@@@*@@@@@@*-::-+@@@@@@%@@@+::*@@@@@@@@%=::
 -::-=====-::::-===-:::-===---===========-::%@@@:-===--::-:::-===---===-:::--=====-:::
    :::::::  ::::::::-:::::::::::::::::::::-@@@#:::::::      ::::::::::::  -:::::::
                                         ::=@@@+::
                                         ::::::::
```

> A terminal-based SSH connection manager.

## Overview

**ok2pus** lets you organize your SSH servers in a single YAML config and connect to any of them through an interactive TUI. It checks server availability before connecting and hands off to native SSH seamlessly.

## Installation

```bash
git clone https://github.com/your-username/ok2pus.git
cd ok2pus
go build -o ok2pus
```

> Requires Go 1.25+

## Configuration

```bash
cp config.template.yaml config.yaml
```

```yaml
servers:
  - name: "Production"
    host: "prod.example.com"
    user: "deploy"
    port: 22
  - name: "Dev Database"
    host: "db.example.com"
    user: "ubuntu"
    port: 2222
```

| Field  | Description          | Default  |
|--------|----------------------|----------|
| `name` | Display name         | Required |
| `host` | Server address       | Required |
| `user` | SSH username         | Required |
| `port` | SSH port             | `22`     |

## Usage

```bash
./ok2pus
```

| Key            | Action  |
|----------------|---------|
| `↑` / `k`     | Up      |
| `↓` / `j`     | Down    |
| `Enter`        | Connect |
| `q` / `Ctrl+C` | Quit   |

## Tech Stack

- **Go** with [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **syscall.Exec** for native SSH process replacement

## Project Structure

```
ok2pus/
├── main.go                # Entry point + SSH execution
├── tui.go                 # Terminal UI
├── config.go              # YAML config parsing
├── config.template.yaml   # Config template
├── go.mod
└── go.sum
```

## License

MIT
