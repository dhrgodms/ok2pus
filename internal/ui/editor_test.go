package ui

import (
	"testing"

	"ok2pus/internal/model"
)

var baseHost = model.SSHHost{
	ID:       1,
	Alias:    "old",
	Host:     "1.1.1.1",
	User:     "root",
	Port:     22,
	AuthType: "Password",
	KeyPath:  "",
}

// --- 정상 케이스 ---

func TestParseHostConfig_ValidFull(t *testing.T) {
	content := "alias=newname\nhost=10.0.0.1\nuser=admin\nport=2222\nauth_type=Public Key\nkey_path=~/.ssh/id_rsa"

	got, err := parseHostConfig(content, baseHost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Alias != "newname" {
		t.Errorf("alias: expected 'newname', got '%s'", got.Alias)
	}
	if got.Host != "10.0.0.1" {
		t.Errorf("host: expected '10.0.0.1', got '%s'", got.Host)
	}
	if got.User != "admin" {
		t.Errorf("user: expected 'admin', got '%s'", got.User)
	}
	if got.Port != 2222 {
		t.Errorf("port: expected 2222, got %d", got.Port)
	}
	if got.AuthType != "Public Key" {
		t.Errorf("auth_type: expected 'Public Key', got '%s'", got.AuthType)
	}
	if got.KeyPath != "~/.ssh/id_rsa" {
		t.Errorf("key_path: expected '~/.ssh/id_rsa', got '%s'", got.KeyPath)
	}
}

func TestParseHostConfig_PreservesID(t *testing.T) {
	content := "alias=test\nhost=10.0.0.1"

	got, err := parseHostConfig(content, baseHost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != baseHost.ID {
		t.Errorf("ID should be preserved: expected %d, got %d", baseHost.ID, got.ID)
	}
}

// --- 빈 콘텐츠 ---

func TestParseHostConfig_EmptyContent(t *testing.T) {
	_, err := parseHostConfig("", baseHost)
	if err == nil {
		t.Fatal("expected error for empty content, got nil")
	}
}

func TestParseHostConfig_WhitespaceOnly(t *testing.T) {
	_, err := parseHostConfig("   \n\n  \t  ", baseHost)
	if err == nil {
		t.Fatal("expected error for whitespace-only content, got nil")
	}
}

// --- 필수 필드 누락 ---

func TestParseHostConfig_MissingAlias(t *testing.T) {
	content := "host=10.0.0.1\nuser=root"

	_, err := parseHostConfig(content, baseHost)
	if err == nil {
		t.Fatal("expected error when alias is missing, got nil")
	}
}

func TestParseHostConfig_MissingHost(t *testing.T) {
	content := "alias=myserver\nuser=root"

	_, err := parseHostConfig(content, baseHost)
	if err == nil {
		t.Fatal("expected error when host is missing, got nil")
	}
}

func TestParseHostConfig_EmptyAliasValue(t *testing.T) {
	content := "alias=\nhost=10.0.0.1"

	_, err := parseHostConfig(content, baseHost)
	if err == nil {
		t.Fatal("expected error when alias value is empty, got nil")
	}
}

func TestParseHostConfig_EmptyHostValue(t *testing.T) {
	content := "alias=myserver\nhost="

	_, err := parseHostConfig(content, baseHost)
	if err == nil {
		t.Fatal("expected error when host value is empty, got nil")
	}
}

// --- 잘못된 port ---

func TestParseHostConfig_InvalidPort_KeepsOriginal(t *testing.T) {
	content := "alias=test\nhost=10.0.0.1\nport=abc"

	got, err := parseHostConfig(content, baseHost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Port != baseHost.Port {
		t.Errorf("port should keep original (%d) on invalid input, got %d", baseHost.Port, got.Port)
	}
}

func TestParseHostConfig_NegativePort_Error(t *testing.T) {
	content := "alias=test\nhost=10.0.0.1\nport=-1"

	_, err := parseHostConfig(content, baseHost)
	if err == nil {
		t.Fatal("expected error for negative port, got nil")
	}
}

func TestParseHostConfig_PortZero_Error(t *testing.T) {
	content := "alias=test\nhost=10.0.0.1\nport=0"

	_, err := parseHostConfig(content, baseHost)
	if err == nil {
		t.Fatal("expected error for port 0, got nil")
	}
}

func TestParseHostConfig_PortTooLarge_Error(t *testing.T) {
	content := "alias=test\nhost=10.0.0.1\nport=70000"

	_, err := parseHostConfig(content, baseHost)
	if err == nil {
		t.Fatal("expected error for port > 65535, got nil")
	}
}

func TestParseHostConfig_PortBoundaryMin(t *testing.T) {
	content := "alias=test\nhost=10.0.0.1\nport=1"

	got, err := parseHostConfig(content, baseHost)
	if err != nil {
		t.Fatalf("expected no error for port 1, got %v", err)
	}
	if got.Port != 1 {
		t.Errorf("expected port 1, got %d", got.Port)
	}
}

func TestParseHostConfig_PortBoundaryMax(t *testing.T) {
	content := "alias=test\nhost=10.0.0.1\nport=65535"

	got, err := parseHostConfig(content, baseHost)
	if err != nil {
		t.Fatalf("expected no error for port 65535, got %v", err)
	}
	if got.Port != 65535 {
		t.Errorf("expected port 65535, got %d", got.Port)
	}
}

// --- 엣지 케이스 ---

func TestParseHostConfig_ExtraWhitespace(t *testing.T) {
	content := "  alias = myserver  \n  host = 10.0.0.1  \n  user = admin  "

	got, err := parseHostConfig(content, baseHost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Alias != "myserver" {
		t.Errorf("alias: expected 'myserver', got '%s'", got.Alias)
	}
}

func TestParseHostConfig_SkipsInvalidLines(t *testing.T) {
	content := "alias=test\nthis line has no equals\nhost=10.0.0.1\n\n"

	got, err := parseHostConfig(content, baseHost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Alias != "test" || got.Host != "10.0.0.1" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestParseHostConfig_UnknownKeysIgnored(t *testing.T) {
	content := "alias=test\nhost=10.0.0.1\nunknown_key=some_value"

	got, err := parseHostConfig(content, baseHost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Alias != "test" {
		t.Errorf("alias: expected 'test', got '%s'", got.Alias)
	}
}

func TestParseHostConfig_ValueContainsEquals(t *testing.T) {
	content := "alias=test\nhost=10.0.0.1\nkey_path=/home/user/.ssh/key=file"

	got, err := parseHostConfig(content, baseHost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.KeyPath != "/home/user/.ssh/key=file" {
		t.Errorf("key_path: expected '/home/user/.ssh/key=file', got '%s'", got.KeyPath)
	}
}

func TestParseHostConfig_OriginalReturnedOnError(t *testing.T) {
	content := ""

	got, err := parseHostConfig(content, baseHost)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got.ID != baseHost.ID || got.Alias != baseHost.Alias {
		t.Errorf("expected original host returned on error, got %+v", got)
	}
}
