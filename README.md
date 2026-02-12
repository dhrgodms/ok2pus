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

> A terminal-based SSH connection manager.

## Overview

**ok2pus**는 SSH 서버 접속 정보를 SQLite DB에 저장하고, 인터랙티브 TUI 메뉴를 통해 호스트 추가/삭제/접속을 관리하는 CLI 도구입니다.

## Installation

```bash
git clone https://github.com/your-username/ok2pus.git
cd ok2pus
go build -o ok2pus
```

> Requires Go 1.25+

## Usage

```bash
./ok2pus
```

실행하면 메인 메뉴가 표시됩니다:

```
? Select:
  > 1. List Hosts
    2. Add New Host
    q. Quit
```

### Add Host

`Add New Host`를 선택하면 Alias, Host Address, Username, Port를 순서대로 입력합니다.

### List Hosts / Connect

`List Hosts`를 선택하면 등록된 호스트 목록이 표시됩니다. 호스트를 선택하면 다음 액션을 고를 수 있습니다:

| Action      | Description                    |
|-------------|--------------------------------|
| `connect`   | SSH 접속                       |
| `edit`      | 호스트 정보 수정 (미구현)       |
| `delete`    | 호스트 삭제 (확인 후 삭제)      |

## Data Storage

호스트 정보는 `~/.ok2pus/hosts.db` SQLite 파일에 저장됩니다.

| Field   | Description    | Default  |
|---------|----------------|----------|
| `alias` | 표시 이름       | Required |
| `host`  | 서버 주소       | Required |
| `user`  | SSH 사용자      | Required |
| `port`  | SSH 포트        | `22`     |

## Tech Stack

- **Go** with [promptui](https://github.com/manifoldco/promptui) (인터랙티브 TUI)
- **SQLite** via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (CGo-free)
- **os/exec** SSH 프로세스 실행

## Project Structure

```
ok2pus/
├── main.go          # Entry point, 메인 메뉴 및 인터랙티브 흐름
├── connect.go       # SSH 접속 실행
├── dbconnect.go     # SQLite DB 초기화
├── repository.go    # SSHHost CRUD 함수
├── logo.go          # 로고 출력
├── go.mod
└── go.sum
```

## License

MIT
