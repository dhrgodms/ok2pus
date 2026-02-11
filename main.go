package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		fmt.Print("Cannot read config.yaml : %v\n", err)
		return
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	initialModel := model{
		servers:  config.Servers,
		cursor:   0,
		selected: -1,
		spinner:  s,
	}

	p := tea.NewProgram(initialModel)
	finalModel, _ := p.Run()

	m, ok := finalModel.(model)
	if !ok {
		log.Fatal("Final model is not of type 'model'")
	}

	if m.selected != -1 {
		runSSH(m.servers[m.selected])
	}
}

func runSSH(target Server) {
	fmt.Printf("You can connect: ssh %s@%s...\n\n", target.User, target.Host)

	// 1. 실행할 ssh 명령어 경로 찾기(binary: ssh실행파일 절대 경로)
	binary, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		fmt.Printf("Command not found for ssh: %v\n", lookErr)
		return
	}

	// 2. ssh 명령어에 넘길 인자 구성
	args := []string{"ssh", fmt.Sprintf("%s@%s", target.User, target.Host)}

	// ssh 포트(22)아닌 경우 -p 옵션 추가
	if target.Port != 0 && target.Port != 22 {
		args = append(args, "-p", fmt.Sprintf("%d", target.Port))
	}

	// 3. 현재 프로세스 환경 변수 가져오기
	env := os.Environ()

	// 4. 현재 go 프로세스를 ssh 프로세스로 대체
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		fmt.Printf("Connection Failed: %v\n", execErr)
	}
}
