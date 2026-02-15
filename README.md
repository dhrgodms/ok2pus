# ok2pus
## 터미널 기반 SSH 접속 관리 도구

```
            -------    ----------
    ------ --+@@*-------+@@@@@@@=---------  ------------ ------
 ---+*##*=---*@@+-=*#+-*@@=--+@@*-=#*#*##+=--=#*=--=*#=--=*###+--
--+@@%*#@@%-=%@%+#@@*=------=@@%--%@@%#%@@@=-#@@=--#@@+-%@@*+%@@--
-=@@#--=%@@=+@@@@@#--- ---+@@@+---@@%---*@@+=@@%--=%@%==%@@%#*=---
-*@@+--=@@%-*@@%@@@=---=#@@#=----+@@#--=%@@=+@@*--+@@#---++*%@@#-
-=@@@#%@@#-=%@%==@@@+-#@@@@@@@@=-%@@@%%@@@=-+@@@@%@@@+-%@@##@@@=-
 --=*#*+-------------------------@@%-+**=----=+****+=---=+***=--
   -----  ------ ---------------+@@+----      --------  -----
                              -------
```

## Privacy & Security

> [!NOTE]
> ok2pus는 외부 서버와 일체의 네트워크 통신을 하지 않습니다.

- 모든 호스트 정보는 사용자 로컬 디스크(`~/.ok2pus/hosts.db`)에만 저장됩니다.
- 원격 서버 전송, 분석 수집, 텔레메트리 기능이 존재하지 않습니다.
- SSH 접속은 시스템에 설치된 `ssh` 명령어를 직접 실행하며, ok2pus가 비밀번호나 키를 별도로 저장하거나 중계하지 않습니다.
- 소스 코드가 전부 공개되어 있으므로 직접 확인할 수 있습니다.

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

`Add New Host`를 선택하면 별칭을 먼저 입력한 뒤, 인증 방식과 호스트 정보를 순서대로 입력합니다.

```
? Alias:       my-server

? Select ssh authentication method:
  > [1] Password
    [2] Public Key

? Username:     root
? Host Address: 192.168.0.10
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

---

# English

> A terminal-based SSH connection manager built with Go.

## Privacy & Security

> [!NOTE]
> ok2pus does not communicate with any external server.

- All host data is stored only on your local disk (`~/.ok2pus/hosts.db`).
- There is no remote transmission, analytics collection, or telemetry of any kind.
- SSH connections are made by directly invoking the system `ssh` command. ok2pus does not store or relay your passwords or keys separately.
- The source code is fully open — you can verify this yourself.

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

Select `Add New Host`, enter an alias first, then choose an authentication method and fill in the host details.

```
? Alias:       my-server

? Select ssh authentication method:
  > [1] Password
    [2] Public Key

? Username:     root
? Host Address: 192.168.0.10
? Port:         22
```

When using Public Key authentication, you will be prompted for the key file path. (Default: `~/.ssh/id_rsa`)

### List Hosts / Connect

Select a host from `List Hosts` to perform the following actions.

| Action    | Description          |
|-----------|----------------------|
| `connect` | Start SSH connection  |
| `edit`    | Edit host config (`$EDITOR` or `vim`) |
| `delete`  | Delete host (with confirmation) |

### Options

| Option           | Description                        |
|------------------|------------------------------------|
| `Reset Database` | Reset all host data                |
| `Drop Database`  | Delete and recreate the DB file    |

## Features

- Store and manage SSH host connection info in a local SQLite DB
- Support for **Password** and **Public Key** authentication
- Interactive TUI menu for adding, listing, connecting, and deleting hosts
- Edit host configuration via system editor
- Port range validation (1–65535)
- Database reset and drop options

## Data Storage

Host information is stored in the SQLite file at `~/.ok2pus/hosts.db`.

| Field       | Description       | Default      |
|-------------|-------------------|--------------|
| `alias`     | Display name (unique) | Required     |
| `host`      | Server address     | Required     |
| `user`      | SSH username       | Required     |
| `port`      | SSH port           | `22`         |
| `auth_type` | Auth method        | `Password`   |
| `key_path`  | Public key file path | `~/.ssh/id_rsa` |

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
