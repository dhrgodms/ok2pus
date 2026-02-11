package main

import (
	"fmt"
	"net"
	"time"

	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	servers  []Server
	cursor   int
	selected int
	loading  bool
	spinner  spinner.Model
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
	s := `                :::::::::       ::::::::::::                                            
                ::+@@%+::     -::-+%@@@@@#=::                                           
     ::::::::: -:-#@@%=::::::::-*@@@@%%@@@@+:::::::::::::: :::::::::::::::::::::::::::  
  :::-+*%%#*=-:::=%@@#-:-+*#*=-*@@@=:::-@@@#::+##*-=*%%*=-::-+##+-:::+*#*-::-=+#%##+=:::
 ::-#@@@@@@@@@=::+@@@*-+@@@%=::::::::::*@@@=::@@@@@@@@@@@#::=%@@#-::-%@@#=:+@@@@@@@@@#-:
::-@@@%-::-@@@%:-#@@@*@@@%=::-    ::-+@@@@-::-@@@@-::-%@@@=:+@@@*:::=@@@*-=@@@#-::=*++::
::%@@@::::-%@@%-=%@@@@@@#-::   :::-*@@@%-::::*@@@=:::-*@@@=-*@@%=:::*@@@+::#@@@@@@@*-:::
::@@@%::::+@@@*:=@@@@@@@@*-:::::=%@@@#-::::::%@@@::::=@@@%:=%@@#-::-#@@%=:::-=+*%@@@@+::
::@@@@+:-*@@@%--*@@@=-%@@@*-:-*@@@@%++++++-:-@@@@#--+@@@@-:=@@@@*=*@@@@#-=%@@*-:-#@@@=::
::-#@@@@@@@%=::=#@@#-:-#@@@*-+@@@@@@@@@@@@-:*@@@*@@@@@@*-::-+@@@@@@%@@@+::*@@@@@@@@%=:: 
 -::-=====-::::-===-:::-===---===========-::%@@@:-===--::-:::-===---===-:::--=====-:::  
    :::::::  ::::::::-:::::::::::::::::::::-@@@#:::::::      ::::::::::::  -:::::::     
                                         ::=@@@+::                                      
                                         ::::::::                                       
	`

	s += "\n\n"
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
