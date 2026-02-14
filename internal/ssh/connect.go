package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"ok2pus/internal/model"
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

	fmt.Printf("\n--- Connecting to [%s] (%s) ---\n", h.Alias, h.Host)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n[SSH Error] %v\n", err)
	}

	fmt.Printf("\n--- Connection to [%s] closed ---\n", h.Alias)
}
