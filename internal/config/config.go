package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"nssm-plus/internal/dpapi"
	"nssm-plus/internal/service"
)

// ConfigFile represents a multi-service configuration file
type ConfigFile struct {
	Services []service.ServiceConfig `json:"services"`
}

// Manager handles configuration file operations
type Manager struct{}

// NewManager creates a new config manager
func NewManager() *Manager {
	return &Manager{}
}

// SaveToFile saves multiple service configurations to a JSON file.
// Passwords are encrypted using Windows DPAPI before saving.
// The input configs slice is not modified; a copy is used internally.
func (m *Manager) SaveToFile(filePath string, configs []service.ServiceConfig) error {
	log.Printf("[config] SaveToFile: path=%q, serviceCount=%d", filePath, len(configs))

	configsCopy := make([]service.ServiceConfig, len(configs))
	copy(configsCopy, configs)

	for i := range configsCopy {
		if configsCopy[i].Password != "" && !dpapi.IsEncrypted(configsCopy[i].Password) {
			encrypted, err := dpapi.Encrypt(configsCopy[i].Password)
			if err != nil {
				log.Printf("[config] SaveToFile: failed to encrypt password for service %q: %v", configsCopy[i].ServiceName, err)
				return fmt.Errorf("failed to encrypt password for service '%s': %w", configsCopy[i].ServiceName, err)
			}
			configsCopy[i].Password = encrypted
		}
	}

	data, err := json.MarshalIndent(ConfigFile{Services: configsCopy}, "", "  ")
	if err != nil {
		log.Printf("[config] SaveToFile: failed to marshal config: %v", err)
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Printf("[config] SaveToFile: failed to write file %q: %v", filePath, err)
		return fmt.Errorf("failed to write config file: %w", err)
	}

	log.Printf("[config] SaveToFile: saved %d service(s) to %q", len(configsCopy), filePath)
	return nil
}

// LoadFromFile loads multiple service configurations from a JSON file.
// Supports three formats for backward compatibility:
//  1. New format: {"services": [...]}
//  2. Bare array: [...]
//  3. Old single-service format: {...}
//
// For backward compat with old merged appPath format, if appPath contains
// a space and arguments is empty, the first token is kept as appPath
// and the rest is moved to arguments.
func (m *Manager) LoadFromFile(filePath string) ([]service.ServiceConfig, error) {
	log.Printf("[config] LoadFromFile: path=%q", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("[config] LoadFromFile: failed to read file %q: %v", filePath, err)
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configs []service.ServiceConfig

	var cf ConfigFile
	if err := json.Unmarshal(data, &cf); err == nil && len(cf.Services) > 0 {
		configs = cf.Services
		log.Printf("[config] LoadFromFile: parsed as multi-service format, found %d service(s)", len(configs))
	} else {
		var arr []service.ServiceConfig
		if err := json.Unmarshal(data, &arr); err == nil && len(arr) > 0 {
			configs = arr
			log.Printf("[config] LoadFromFile: parsed as bare array format, found %d service(s)", len(configs))
		} else {
			var single service.ServiceConfig
			if err := json.Unmarshal(data, &single); err == nil && single.ServiceName != "" {
				configs = []service.ServiceConfig{single}
				log.Printf("[config] LoadFromFile: parsed as single-service format, name=%q", single.ServiceName)
			} else {
				log.Printf("[config] LoadFromFile: unrecognized format in %q", filePath)
				return nil, fmt.Errorf("failed to parse config file: unrecognized format")
			}
		}
	}

	for i := range configs {
		if configs[i].Arguments == "" && configs[i].AppPath != "" {
			if idx := strings.Index(configs[i].AppPath, " "); idx >= 0 {
				log.Printf("[config] LoadFromFile: splitting merged appPath for service %q (space at index %d)", configs[i].ServiceName, idx)
				configs[i].Arguments = configs[i].AppPath[idx+1:]
				configs[i].AppPath = configs[i].AppPath[:idx]
			}
		}
	}

	for i := range configs {
		if configs[i].Password != "" && dpapi.IsEncrypted(configs[i].Password) {
			decrypted, err := dpapi.Decrypt(configs[i].Password)
			if err != nil {
				log.Printf("[config] LoadFromFile: failed to decrypt password for service %q: %v", configs[i].ServiceName, err)
				fmt.Fprintf(os.Stderr, "Warning: failed to decrypt password for service '%s': %v\n", configs[i].ServiceName, err)
			} else {
				configs[i].Password = decrypted
			}
		}
	}

	log.Printf("[config] LoadFromFile: loaded %d service(s) from %q", len(configs), filePath)
	return configs, nil
}
