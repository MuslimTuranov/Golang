package modules

import "time"

type PostgreConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	DBName      string
	SSLMode     string
	ExecTimeout time.Duration
}

type AppConfig struct {
	Port      string
	APIKey    string
	JWTSecret string
	JWTTTL    time.Duration
	DB        PostgreConfig
}