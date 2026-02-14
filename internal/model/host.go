package model

type SSHHost struct {
	ID    int
	Alias string
	Host  string
	User  string
	Port  int

	AuthType string
	KeyPath  string
}
