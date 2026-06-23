package webconnector

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const Version = "1.0.0"

const (
	DefaultBaseDir       = `C:\Documents\psrt`
	DefaultAllowedOrigin = "https://editor.psrt.app"
	DefaultPort          = 5278
)

type Config struct {
	BaseDir       string
	AllowedOrigin string
	Port          int
}

func DefaultConfig() *Config {
	return &Config{
		BaseDir:       DefaultBaseDir,
		AllowedOrigin: DefaultAllowedOrigin,
		Port:          DefaultPort,
	}
}

// LoadOrCreateConfig loads the INI at path, creating a default config (and its
// base_dir) when the file does not exist yet — e.g. first run of the binary.
func LoadOrCreateConfig(path string) (*Config, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("config path: %w", err)
	}
	if _, err := os.Stat(abs); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		cfg := DefaultConfig()
		if err := os.MkdirAll(cfg.BaseDir, 0o755); err != nil {
			return nil, fmt.Errorf("create base_dir %q: %w", cfg.BaseDir, err)
		}
		if err := cfg.Save(abs); err != nil {
			return nil, fmt.Errorf("write default config %q: %w", abs, err)
		}
	}
	return LoadConfig(abs)
}

func LoadConfig(path string) (*Config, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("config path: %w", err)
	}
	f, err := os.Open(abs)
	if err != nil {
		return nil, fmt.Errorf("open config %q: %w", abs, err)
	}
	defer f.Close()

	vals := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		i := strings.Index(line, "=")
		if i < 0 {
			continue
		}
		key := strings.TrimSpace(line[:i])
		val := strings.TrimSpace(line[i+1:])
		vals[key] = val
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	var missing []string
	if v, ok := vals["base_dir"]; ok && v != "" {
		cfg.BaseDir = v
	} else {
		missing = append(missing, "base_dir")
	}
	if v, ok := vals["allowed_origin"]; ok && v != "" {
		cfg.AllowedOrigin = v
	} else {
		missing = append(missing, "allowed_origin")
	}
	if v, ok := vals["port"]; ok && v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid port %q", v)
		}
		cfg.Port = p
	} else {
		missing = append(missing, "port")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("config missing required keys: %s", strings.Join(missing, ", "))
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if err := cfg.NormalizeBaseDir(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Port < 1024 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1024 and 65535, got %d", c.Port)
	}
	normalized, err := parseAllowedOriginsField(c.AllowedOrigin)
	if err != nil {
		return err
	}
	c.AllowedOrigin = normalized
	info, err := os.Stat(c.BaseDir)
	if err != nil {
		return fmt.Errorf("base_dir %q: %w", c.BaseDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("base_dir %q is not a directory", c.BaseDir)
	}
	return nil
}

func (c *Config) NormalizeBaseDir() error {
	clean := filepath.Clean(c.BaseDir)
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		return fmt.Errorf("base_dir resolve: %w", err)
	}
	c.BaseDir = resolved
	return nil
}

func (c *Config) Save(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("base_dir=")
	b.WriteString(c.BaseDir)
	b.WriteString("\nallowed_origin=")
	b.WriteString(c.AllowedOrigin)
	b.WriteString("\nport=")
	b.WriteString(strconv.Itoa(c.Port))
	b.WriteString("\n")
	return os.WriteFile(abs, []byte(b.String()), 0o600)
}

func (c *Config) Clone() *Config {
	return &Config{
		BaseDir:       c.BaseDir,
		AllowedOrigin: c.AllowedOrigin,
		Port:          c.Port,
	}
}
