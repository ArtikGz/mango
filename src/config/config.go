package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

var gconfig GlobalConfig

type GlobalConfig struct {
	Server     ServerConfig   `json:"server"`
	Logger     LoggerConfig   `json:"logger"`
	Profiler   ProfilerConfig `json:"profiler"`
	ConfigPath string
}

func Motd() string {
	return gconfig.Server.Motd
}

func Host() string {
	return gconfig.Server.Host
}

func Port() int {
	return gconfig.Server.Port
}

func IsOnline() bool {
	return gconfig.Server.Online
}

func Protocol() int {
	return gconfig.Server.Protocol
}

func CompressionThreshold() int {
	return gconfig.Server.CompressionThreshold
}

func LogLevel() LoggerLevel {
	switch strings.ToUpper(gconfig.Logger.Level) {
	case "OFF":
		return OFF
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

type ServerConfig struct {
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	Online               bool   `json:"online"`
	Motd                 string `json:"motd"`
	Protocol             int    `json:"protocol"`
	CompressionThreshold int    `json:"compression_threshold"`
}

type LoggerConfig struct {
	Level string `json:"level"`
}

type ProfilerConfig struct {
	Port int `json:"port"`
}

func ProfilerPort() int {
	return gconfig.Profiler.Port
}

type LoggerLevel int

const (
	DEBUG LoggerLevel = iota
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

func Parse(path string) error {
	text, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(text, &gconfig); err != nil {
		return err
	}

	gconfig.ConfigPath, err = filepath.Abs(path)
	if err != nil {
		gconfig.ConfigPath = path
	}

	return nil
}

func LoadDefaultConfig() {
	gconfig = GlobalConfig{
		Server: ServerConfig{
			Host:                 "127.0.0.1",
			Port:                 25565,
			Online:               false,
			Motd:                 "Powered by man.go",
			Protocol:             762,
			CompressionThreshold: 256,
		},
		Logger: LoggerConfig{
			Level: "INFO",
		},
		Profiler: ProfilerConfig{
			Port: 8080,
		},
		ConfigPath: "DEFAULT_CONFIG",
	}
}

func GetConfigPath() string {
	return gconfig.ConfigPath
}
