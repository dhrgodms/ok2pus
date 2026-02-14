# ok2pus

```
           --===--    ----===---
    ------ --+@@*-------+@@@@@@@=---------  ------------ ------
 ---+*##*=---*@@+-=*#+-*@@=--+@@*-=#*=*##+=--=#*=--=*#=--=*###+--
--+@@%*#@@%-=%@%+#@@*=------=@@%--%@@%#%@@@=-#@@=--#@@+-%@@*+%@@--
-=@@#--=%@@=+@@@@@#--- ---+@@@+---@@%---*@@+=@@%--=%@%==%@@%#*=---
-*@@+--=@@%-*@@%@@@=---=#@@#=----+@@#--=%@@=+@@*--+@@#---++*%@@#-
-=@@@#%@@#-=%@%==@@@+-#@@@@@@@@=-%@@@%%@@@=-+@@@@%@@@+-%@@##@@@=-
 --=*#*+----++=--=++===++++++++--@@%-+**=----=+*+==+=---=+***=--
   -----  ------ ---------------+@@+----      --------  -----
                              --===--
```

> A terminal-based SSH connection manager built with Go.

## Installation

### Homebrew (Recommended)

```bash
brew tap dhrgodms/ok2pus
brew install ok2pus
```

### Build from source

```bash
git clone https://github.com/dhrgodms/ok2pus.git
cd ok2pus
go build -o ok2pus ./cmd/ok2pus
```

> Requires Go 1.25+

## Usage

```bash
ok2pus
```

### Main Menu

```
? Select:
  > [1] List Hosts
    [2] Add New Host
    [3] Options
    [q] Quit
```

### Add New Host

`Add New Host`를 선택하면 인증 방식을 먼저 고른 뒤, 호스트 정보를 순서대로 입력합니다.

```
? Select ssh authentication method:
  > [1] Password
    [2] Public Key

? Alias:       my-server
? Host Address: 192.168.0.10
? Username:     root
? Port:         22
```

Public Key 인증 시 키 파일 경로를 추가로 입력합니다. (기본값: `~/.ssh/id_rsa`)

### List Hosts / Connect

`List Hosts`에서 호스트를 선택하면 아래 액션을 수행할 수 있습니다.

| Action    | Description          |
|-----------|----------------------|
| `connect` | SSH 접속 실행         |
| `edit`    | 호스트 정보 수정 (`$EDITOR` 또는 `vim`) |
| `delete`  | 호스트 삭제 (확인 후) |

### Options

| Option           | Description                        |
|------------------|------------------------------------|
| `Reset Database` | 모든 호스트 데이터 초기화            |
| `Drop Database`  | DB 파일 삭제 후 재생성              |

## Features

- SSH 호스트 접속 정보를 로컬 SQLite DB에 저장 및 관리
- **Password** / **Public Key** 두 가지 인증 방식 지원
- 인터랙티브 TUI 메뉴를 통한 호스트 추가, 목록 조회, 접속, 삭제
- 시스템 에디터를 통한 호스트 설정 편집
- Port 범위 검증 (1–65535)
- DB 초기화(Reset) 및 삭제(Drop) 옵션 제공

## Data Storage

호스트 정보는 `~/.ok2pus/hosts.db` SQLite 파일에 저장됩니다.

| Field       | Description       | Default      |
|-------------|-------------------|--------------|
| `alias`     | 표시 이름 (고유)   | Required     |
| `host`      | 서버 주소          | Required     |
| `user`      | SSH 사용자         | Required     |
| `port`      | SSH 포트           | `22`         |
| `auth_type` | 인증 방식          | `Password`   |
| `key_path`  | 공개키 파일 경로    | `~/.ssh/id_rsa` |

## Project Structure

```
ok2pus/
├── cmd/ok2pus/
│   └── main.go              # Entry point
├── internal/
│   ├── model/
│   │   └── host.go          # SSHHost data model
│   ├── db/
│   │   ├── sqlite.go        # DB initialization & drop
│   │   └── repository.go    # CRUD operations
│   ├── ssh/
│   │   └── connect.go       # SSH connection
│   └── ui/
│       ├── menu.go          # Host list & action menu
│       ├── editor.go        # Config editor & parser
│       ├── host_form.go     # Add new host form
│       ├── options.go       # Reset & drop options
│       └── logo.go          # ASCII logo
├── .goreleaser.yml
├── go.mod
└── go.sum
```

## Tech Stack

- **Go** with [promptui](https://github.com/manifoldco/promptui) — Interactive TUI
- **SQLite** via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) — CGo-free pure Go driver
- **GoReleaser** — Cross-compilation & Homebrew distribution

## License

MIT
