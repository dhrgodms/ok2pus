package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"ok2pus/internal/model"

	"github.com/fatih/color"
)

func ConnectHost(h model.SSHHost) {
	var cmd *exec.Cmd

	dest := fmt.Sprintf("%s@%s", h.User, h.Host)

	switch h.AuthType {
	case "Password":
		cmd = exec.Command("ssh", dest, "-p", strconv.Itoa(h.Port))
	case "Public Key":
		cmd = exec.Command("ssh", "-i", h.KeyPath, dest, "-p", strconv.Itoa(h.Port))
	default:
		return
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println()
	color.New(color.Bold, color.BgCyan).Printf("--- Connecting to [%s] (%s) ---", h.Alias, h.Host)
	fmt.Println()

	err := cmd.Run()
	if err != nil {
		color.Red("\n[SSH Error] %v\n", err)
	}

	fmt.Println()
	color.New(color.Bold, color.BgCyan).Printf("--- Connection to [%s] closed ---", h.Alias)
	fmt.Println()
}
