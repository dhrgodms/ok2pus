package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"net"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	spinner "github.com/charmbracelet/bubbles/spinner"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Port int `yaml:"port"`
}

type model struct {
	servers	[]Server
	cursor	int
	selected	int
	loading bool
	spinner spinner.Model
}

type checkResultMsg struct {
	err error
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.servers)-1 {
				m.cursor++
			}
		case "enter":
			if !m.loading {
				m.selected = m.cursor
				m.loading = true
				return m, tea.Batch(m.spinner.Tick, checkServer(m.servers[m.cursor]))
			}
		}
	case checkResultMsg:
		if msg.err != nil {
			m.loading = false
			return m, nil
		}
		return m, tea.Quit

	case spinner.TickMsg: 
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	}
	return m, nil
}

func (m model) View() string {
	s := "🚀 Select a server to connect (q: quit):\n\n"

	for i, server := range m.servers {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s [%d] %s (%s@%s)\n", cursor, i+1, server.Name, server.User, server.Host)
	}

	if m.selected != -1 {
		return fmt.Sprintf("\nCheck for connecting %s server...\n\n", m.servers[m.selected].Name)
	}

	return s + "\n"
}

func checkServer(s Server) tea.Cmd {
	return func() tea.Msg {
		port := s.Port
		if port == 0 {
			port = 22
		}

		address := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err != nil {
			return checkResultMsg{err: err}
		}
		conn.Close()
		return checkResultMsg{err: nil}
	}
}

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Print("Cannot read config.yaml : %v\n", err)
		return
	}

	var config Config
	
	// yaml 형식의 데이터를 go의 구조체, 맵 같은 변수에 담도록 역직렬화
	// 바이트 슬라이스(data)로 분석하여 정의한 변수의 사용자가 정의한 메모리 주소(&config)에 값 하나씩 채워넣음
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Parsing Error: %v\n", err)
		return
	}

	// spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))


	initialModel := model {
		servers: config.Servers,
		cursor: 0,
		selected: -1,
		spinner: s,
	}

	p := tea.NewProgram(initialModel)
	finalModel, _ := p.Run()

	m, ok := finalModel.(model)
	if !ok {
		log.Fatal("Final model is not of type 'model'")
	}

	if m.selected != -1 {
		target := m.servers[m.selected]
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
}