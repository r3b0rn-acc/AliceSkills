package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Env string

const (
	EnvDev  Env = "dev"
	EnvProd Env = "prod"
)

type Config struct {
	Port            int
	Env             Env
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	LogLevel        string
}

func (cfg Config) Addr() string {
	return ":" + strconv.Itoa(cfg.Port)
}

func validate(cfg *Config) error {
	if cfg.Port < minPort || cfg.Port > maxPort {
		return fmt.Errorf("%s invalid: got %d, expected %d..%d", envPort, cfg.Port, minPort, maxPort)
	}
	if !allowedEnvs[cfg.Env] {
		return fmt.Errorf("%s invalid: got %q, expected one of %s",
			envEnv, cfg.Env, strings.Join(allowedEnvsList, "|"))
	}
	if cfg.ReadTimeout < minReadWriteTimeout || cfg.ReadTimeout > maxReadWriteTimeout {
		return fmt.Errorf("%s invalid: got %s, expected %s..%s",
			envReadTimeout, cfg.ReadTimeout, minReadWriteTimeout, maxReadWriteTimeout)
	}
	if cfg.WriteTimeout < minReadWriteTimeout || cfg.WriteTimeout > maxReadWriteTimeout {
		return fmt.Errorf("%s invalid: got %s, expected %s..%s",
			envWriteTimeout, cfg.WriteTimeout, minReadWriteTimeout, maxReadWriteTimeout)
	}
	if cfg.ShutdownTimeout < minShutdownTimeout || cfg.ShutdownTimeout > maxShutdownTimeout {
		return fmt.Errorf("%s invalid: got %s, expected %s..%s",
			envShutdownTimeout, cfg.ShutdownTimeout, minShutdownTimeout, maxShutdownTimeout)
	}
	if !allowedLogLevels[cfg.LogLevel] {
		return fmt.Errorf("%s invalid: got %q, expected one of %s",
			envLogLevel, cfg.LogLevel, strings.Join(allowedLogLevelsList, "|"))
	}
	return nil
}

const (
	envPort            = "SKILLSRV_PORT"
	envEnv             = "SKILLSRV_ENV"
	envReadTimeout     = "SKILLSRV_READ_TIMEOUT"
	envWriteTimeout    = "SKILLSRV_WRITE_TIMEOUT"
	envShutdownTimeout = "SKILLSRV_SHUTDOWN_TIMEOUT"
	envLogLevel        = "SKILLSRV_LOG_LEVEL"
)

const (
	defaultPort            int           = 8080
	defaultEnv             Env           = "dev"
	defaultReadTimeout     time.Duration = 5 * time.Second
	defaultWriteTimeout    time.Duration = 5 * time.Second
	defaultShutdownTimeout time.Duration = 10 * time.Second
	defaultLogLevel        string        = "info"
)

const (
	minPort             int           = 1
	maxPort             int           = 65535
	minReadWriteTimeout time.Duration = 100 * time.Millisecond
	maxReadWriteTimeout time.Duration = 120 * time.Second
	minShutdownTimeout  time.Duration = 1 * time.Second
	maxShutdownTimeout  time.Duration = 30 * time.Second
)

var (
	allowedEnvs = map[Env]bool{
		EnvDev:  true,
		EnvProd: true,
	}
	allowedLogLevels = map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	allowedEnvsList      = []string{"dev", "prod"}
	allowedLogLevelsList = []string{"debug", "info", "warn", "error"}
)

func getEnvInt(name string, defaultValue int) (int, error) {
	val, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue, nil
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return 0, fmt.Errorf("%s is set but empty", name)
	}
	parsedInt, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("%s must be integer, got %q", name, val)
	}
	return parsedInt, nil
}

func getEnvString(name string, defaultValue string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}
	return strings.TrimSpace(val)
}

func getEnvDuration(name string, defaultValue time.Duration) (time.Duration, error) {
	val, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue, nil
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return 0, fmt.Errorf("%s is set but empty", name)
	}
	valDuration, err := time.ParseDuration(val)
	if err != nil {
		valInt, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("%s must be duration (e.g. 200ms, 5s, 2m) or integer seconds, got %q", name, val)
		}
		return time.Duration(valInt) * time.Second, nil
	}
	return valDuration, nil
}

func Load() (Config, error) {
	port, err := getEnvInt(envPort, defaultPort)
	if err != nil {
		return Config{}, err
	}

	env := getEnvString(envEnv, string(defaultEnv))
	env = strings.ToLower(strings.TrimSpace(env))

	readTimeout, err := getEnvDuration(envReadTimeout, defaultReadTimeout)
	if err != nil {
		return Config{}, err
	}
	writeTimeout, err := getEnvDuration(envWriteTimeout, defaultWriteTimeout)
	if err != nil {
		return Config{}, err
	}
	shutdownTimeout, err := getEnvDuration(envShutdownTimeout, defaultShutdownTimeout)
	if err != nil {
		return Config{}, err
	}

	logLevel := getEnvString(envLogLevel, defaultLogLevel)
	logLevel = strings.ToLower(strings.TrimSpace(logLevel))

	cfg := Config{
		Port:            port,
		Env:             Env(env),
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		ShutdownTimeout: shutdownTimeout,
		LogLevel:        logLevel,
	}

	if err := validate(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
