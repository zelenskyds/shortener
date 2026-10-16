package main

import (
	"flag"
	"os"
	"testing"

	"github.com/zelenskyds/shortener/internal/config"
)

func TestParseFlag(t *testing.T) {
	defaultServer, defaultBase, defaultLogLevel := config.ServerAddress, config.BaseURL, config.LogLevel
	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		wantServer string
		wantBase   string
	}{
		{
			name:       "defaults without flags or environment",
			wantServer: "localhost:8080",
			wantBase:   "http://localhost:8080",
		},
		{
			name:       "flags override defaults",
			args:       []string{"-a", "localhost:9000", "-b", "https://flags.example.com"},
			wantServer: "localhost:9000",
			wantBase:   "https://flags.example.com",
		},
		{
			name: "environment overrides defaults",
			env: map[string]string{
				"SERVER_ADDRESS": "localhost:9001",
				"BASE_URL":       "https://env.example.com",
			},
			wantServer: "localhost:9001",
			wantBase:   "https://env.example.com",
		},
		{
			name: "environment overrides flags",
			args: []string{"-a", "localhost:9000", "-b", "http://flags.example.com"},
			env: map[string]string{
				"SERVER_ADDRESS": "https://localhost:9001",
				"BASE_URL":       "env.example.com",
			},
			wantServer: "localhost:9001",
			wantBase:   "http://env.example.com",
		},
		{
			name:       "server environment and base flag are independent",
			args:       []string{"-a", "localhost:9000", "-b", "flags.example.com"},
			env:        map[string]string{"SERVER_ADDRESS": "http://localhost:9001"},
			wantServer: "localhost:9001",
			wantBase:   "http://flags.example.com",
		},
		{
			name:       "base environment and server flag are independent",
			args:       []string{"-a", "http://localhost:9000", "-b", "flags.example.com"},
			env:        map[string]string{"BASE_URL": "https://env.example.com/links/"},
			wantServer: "localhost:9000",
			wantBase:   "https://env.example.com/links/",
		},
		{
			name:       "empty environment falls back to flags",
			args:       []string{"-a", "localhost:9000", "-b", "flags.example.com"},
			env:        map[string]string{"SERVER_ADDRESS": "", "BASE_URL": ""},
			wantServer: "localhost:9000",
			wantBase:   "http://flags.example.com",
		},
		{
			name:       "empty environment falls back to defaults",
			env:        map[string]string{"SERVER_ADDRESS": "", "BASE_URL": ""},
			wantServer: "localhost:8080",
			wantBase:   "http://localhost:8080",
		},
		{
			name:       "changing server does not change default base URL",
			env:        map[string]string{"SERVER_ADDRESS": "localhost:9001"},
			wantServer: "localhost:9001",
			wantBase:   "http://localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFlags, oldArgs := flag.CommandLine, os.Args
			oldServer, oldBase, oldLogLevel := config.ServerAddress, config.BaseURL, config.LogLevel
			t.Cleanup(func() {
				flag.CommandLine, os.Args = oldFlags, oldArgs
				config.ServerAddress, config.BaseURL = oldServer, oldBase
				config.LogLevel = oldLogLevel
			})
			flag.CommandLine = flag.NewFlagSet("shortener", flag.ContinueOnError)
			os.Args = append([]string{"shortener"}, tt.args...)
			config.ServerAddress, config.BaseURL = defaultServer, defaultBase
			config.LogLevel = defaultLogLevel

			for _, key := range []string{"SERVER_ADDRESS", "BASE_URL", "LOG_LEVEL"} {
				t.Setenv(key, "")
				if err := os.Unsetenv(key); err != nil {
					t.Fatal(err)
				}
			}
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			ParseFlag()

			if config.ServerAddress != tt.wantServer {
				t.Errorf("ServerAddress = %q, want %q", config.ServerAddress, tt.wantServer)
			}
			if config.BaseURL != tt.wantBase {
				t.Errorf("BaseURL = %q, want %q", config.BaseURL, tt.wantBase)
			}
		})
	}
}

func TestParseFlagLogLevel(t *testing.T) {
	defaultLogLevel := config.LogLevel
	tests := []struct {
		name string
		args []string
		env  map[string]string
		want string
	}{
		{
			name: "default without flags or environment",
			want: defaultLogLevel,
		},
		{
			name: "flag overrides default",
			args: []string{"-l", "debug"},
			want: "debug",
		},
		{
			name: "environment overrides default",
			env:  map[string]string{"LOG_LEVEL": "warn"},
			want: "warn",
		},
		{
			name: "environment overrides flag",
			args: []string{"-l", "debug"},
			env:  map[string]string{"LOG_LEVEL": "info"},
			want: "info",
		},
		{
			name: "empty environment falls back to flag",
			args: []string{"-l", "debug"},
			env:  map[string]string{"LOG_LEVEL": ""},
			want: "debug",
		},
		{
			name: "empty environment falls back to default",
			env:  map[string]string{"LOG_LEVEL": ""},
			want: defaultLogLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFlags, oldArgs := flag.CommandLine, os.Args
			oldServer, oldBase, oldLogLevel := config.ServerAddress, config.BaseURL, config.LogLevel
			t.Cleanup(func() {
				flag.CommandLine, os.Args = oldFlags, oldArgs
				config.ServerAddress, config.BaseURL, config.LogLevel = oldServer, oldBase, oldLogLevel
			})

			flag.CommandLine = flag.NewFlagSet("shortener", flag.ContinueOnError)
			os.Args = append([]string{"shortener"}, tt.args...)
			config.LogLevel = defaultLogLevel

			for _, key := range []string{"SERVER_ADDRESS", "BASE_URL", "LOG_LEVEL"} {
				t.Setenv(key, "")
				if err := os.Unsetenv(key); err != nil {
					t.Fatal(err)
				}
			}
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			ParseFlag()

			if config.LogLevel != tt.want {
				t.Errorf("LogLevel = %q, want %q", config.LogLevel, tt.want)
			}
		})
	}
}

func TestSchemeHandling(t *testing.T) {
	tests := []struct {
		address     string
		wantWith    string
		wantTrimmed string
	}{
		{"localhost:8080", "http://localhost:8080", "localhost:8080"},
		{"http://localhost:8080", "http://localhost:8080", "localhost:8080"},
		{"https://localhost:8080", "https://localhost:8080", "localhost:8080"},
		{"https://example.com/links/", "https://example.com/links/", "example.com/links/"},
		{"[::1]:8080", "http://[::1]:8080", "[::1]:8080"},
	}
	for _, tt := range tests {
		t.Run(tt.address, func(t *testing.T) {
			if got := withScheme(tt.address); got != tt.wantWith {
				t.Errorf("withScheme(%q) = %q, want %q", tt.address, got, tt.wantWith)
			}
			if got := trimScheme(tt.address); got != tt.wantTrimmed {
				t.Errorf("trimScheme(%q) = %q, want %q", tt.address, got, tt.wantTrimmed)
			}
		})
	}
}
