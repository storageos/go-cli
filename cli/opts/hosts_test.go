package opts

import (
	"strings"
	"testing"
)

func TestParseTCP(t *testing.T) {
	var (
		defaultHTTPHost = "tcp://127.0.0.1:5705"
	)
	invalids := map[string]string{
		"tcp:a.b.c.d":          "Invalid bind address format: tcp:a.b.c.d",
		"tcp:a.b.c.d/path":     "Invalid bind address format: tcp:a.b.c.d/path",
		"udp://127.0.0.1":      "Invalid proto, expected tcp: udp://127.0.0.1",
		"udp://127.0.0.1:2375": "Invalid proto, expected tcp: udp://127.0.0.1:2375",
	}
	valids := map[string]string{
		"":                            defaultHTTPHost,
		"tcp://":                      defaultHTTPHost,
		"0.0.0.1:":                    "tcp://0.0.0.1:5705",
		"0.0.0.1:5555":                "tcp://0.0.0.1:5555",
		"0.0.0.1:5555/path":           "tcp://0.0.0.1:5555/path",
		":6666":                       "tcp://127.0.0.1:6666",
		":6666/path":                  "tcp://127.0.0.1:6666/path",
		"tcp://:7777":                 "tcp://127.0.0.1:7777",
		"tcp://:7777/path":            "tcp://127.0.0.1:7777/path",
		"[::1]:":                      "tcp://[::1]:5705",
		"[::1]:5555":                  "tcp://[::1]:5555",
		"[::1]:5555/path":             "tcp://[::1]:5555/path",
		"[0:0:0:0:0:0:0:1]:":          "tcp://[0:0:0:0:0:0:0:1]:5705",
		"[0:0:0:0:0:0:0:1]:5555":      "tcp://[0:0:0:0:0:0:0:1]:5555",
		"[0:0:0:0:0:0:0:1]:5555/path": "tcp://[0:0:0:0:0:0:0:1]:5555/path",
		"localhost:":                  "tcp://localhost:5705",
		"localhost:5555":              "tcp://localhost:5555",
		"localhost:5555/path":         "tcp://localhost:5555/path",
	}
	for invalidAddr, expectedError := range invalids {
		if addr, err := ParseTCPAddr(invalidAddr, defaultHTTPHost); err == nil || err.Error() != expectedError {
			t.Errorf("tcp %v address expected error %v return, got %s and addr %v", invalidAddr, expectedError, err, addr)
		}
	}
	for validAddr, expectedAddr := range valids {
		if addr, err := ParseTCPAddr(validAddr, defaultHTTPHost); err != nil || addr != expectedAddr {
			t.Errorf("%v -> expected %v, got %v and addr %v", validAddr, expectedAddr, err, addr)
		}
	}
}

func TestParseInvalidUnixAddrInvalid(t *testing.T) {
	if _, err := parseSimpleProtoAddr("unix", "tcp://127.0.0.1", "unix:///var/run/docker.sock"); err == nil || err.Error() != "Invalid proto, expected unix: tcp://127.0.0.1" {
		t.Fatalf("Expected an error, got %v", err)
	}
	if _, err := parseSimpleProtoAddr("unix", "unix://tcp://127.0.0.1", "/var/run/docker.sock"); err == nil || err.Error() != "Invalid proto, expected unix: tcp://127.0.0.1" {
		t.Fatalf("Expected an error, got %v", err)
	}
	if v, err := parseSimpleProtoAddr("unix", "", "/var/run/docker.sock"); err != nil || v != "unix:///var/run/docker.sock" {
		t.Fatalf("Expected an %v, got %v", v, "unix:///var/run/docker.sock")
	}
}

func TestValidateExtraHosts(t *testing.T) {
	valid := []string{
		`myhost:192.168.0.1`,
		`thathost:10.0.2.1`,
		`anipv6host:2003:ab34:e::1`,
		`ipv6local:::1`,
	}

	invalid := map[string]string{
		`myhost:192.notanipaddress.1`:  `invalid IP`,
		`thathost-nosemicolon10.0.0.1`: `bad format`,
		`anipv6host:::::1`:             `invalid IP`,
		`ipv6local:::0::`:              `invalid IP`,
	}

	for _, extrahost := range valid {
		if _, err := ValidateExtraHost(extrahost); err != nil {
			t.Fatalf("ValidateExtraHost(`"+extrahost+"`) should succeed: error %v", err)
		}
	}

	for extraHost, expectedError := range invalid {
		if _, err := ValidateExtraHost(extraHost); err == nil {
			t.Fatalf("ValidateExtraHost(`%q`) should have failed validation", extraHost)
		} else {
			if !strings.Contains(err.Error(), expectedError) {
				t.Fatalf("ValidateExtraHost(`%q`) error should contain %q", extraHost, expectedError)
			}
		}
	}
}
