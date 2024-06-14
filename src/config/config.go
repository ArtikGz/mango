package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

var gconfig GlobalConfig

type GlobalConfig struct {
	Server   ServerConfig   `json:"server"`
	Logger   LoggerConfig   `json:"logger"`
	Profiler ProfilerConfig `json:"profiler"`
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
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Online   bool   `json:"online"`
	Motd     string `json:"motd"`
	Protocol int    `json:"protocol"`
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

func Parse(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal([]byte(file), &gconfig)
}
