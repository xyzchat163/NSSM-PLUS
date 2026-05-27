package wrapper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const (
	baseDirName    = "NSSM-Plus"
	servicesDir    = "services"
	logsDirName    = "logs"
	defaultLogSize = 10 * 1024 * 1024 // 10 MB default log rotation size
)

// ConfigDir returns the base directory for NSSM-Plus data.
func ConfigDir() string {
	// Use ProgramData for system-wide storage
	base := os.Getenv("ProgramData")
	if base == "" {
		base = filepath.Join(os.Getenv("SystemDrive"), "ProgramData")
	}
	return filepath.Join(base, baseDirName)
}

// ServicesDir returns the directory containing service config files.
func ServicesDir() string {
	return filepath.Join(ConfigDir(), servicesDir)
}

// LogsDir returns the directory for service log files.
func LogsDir() string {
	return filepath.Join(ConfigDir(), logsDirName)
}

// ConfigPath returns the config file path for a given service name.
func ConfigPath(serviceName string) string {
	return filepath.Join(ServicesDir(), serviceName+".json")
}

// LogPath returns the log file path for a given service name.
func LogPath(serviceName string) string {
	return filepath.Join(LogsDir(), serviceName+".log")
}

// WrapperConfig holds the application launch configuration for a service.
type WrapperConfig struct {
	AppPath        string            `json:"appPath"`
	Arguments      string            `json:"arguments"`
	WorkDir        string            `json:"workDir"`
	Env            map[string]string `json:"env,omitempty"`
	RestartDelay   int               `json:"restartDelay"`
	LogStdout      string            `json:"logStdout,omitempty"`
	LogStderr      string            `json:"logStderr,omitempty"`
	RotateLog      bool              `json:"rotateLog,omitempty"`
	RestartTimeout int               `json:"restartTimeout,omitempty"`
}

// SaveConfig writes the wrapper config to disk.
func SaveConfig(serviceName string, cfg *WrapperConfig) error {
	if err := os.MkdirAll(ServicesDir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(serviceName), data, 0644)
}

// LoadConfig reads the wrapper config from disk.
func LoadConfig(serviceName string) (*WrapperConfig, error) {
	data, err := os.ReadFile(ConfigPath(serviceName))
	if err != nil {
		return nil, err
	}
	var cfg WrapperConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ConfigExists checks if a wrapper config file exists for the service.
func ConfigExists(serviceName string) bool {
	_, err := os.Stat(ConfigPath(serviceName))
	return err == nil
}

// DeleteConfig removes the wrapper config file.
func DeleteConfig(serviceName string) error {
	return os.Remove(ConfigPath(serviceName))
}

// splitArgs splits a command-line argument string into individual arguments.
// It properly handles quoted arguments (both single and double quotes).
// Examples:
//
//	"-jar \"C:\\My App\\app.jar\"" -> ["-jar", "C:\\My App\\app.jar"]
//	"--config='my config.json'"     -> ["--config=my config.json"]
//	"say \"hello world\""           -> ["say", "hello world"]
func splitArgs(argStr string) []string {
	argStr = strings.TrimSpace(argStr)
	if argStr == "" {
		return nil
	}

	var args []string
	var current strings.Builder
	inQuote := false
	var quoteChar byte

	for i := 0; i < len(argStr); i++ {
		ch := argStr[i]

		if inQuote {
			if ch == '\\' && i+1 < len(argStr) && (argStr[i+1] == quoteChar || argStr[i+1] == '\\') {
				current.WriteByte(argStr[i+1])
				i++
			} else if ch == quoteChar {
				inQuote = false
			} else {
				current.WriteByte(ch)
			}
		} else {
			switch ch {
			case '"', '\'':
				inQuote = true
				quoteChar = ch
			case ' ', '\t':
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			default:
				current.WriteByte(ch)
			}
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
