package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Settings struct {
	Main  MainSection  `toml:"main"`
	Excel ExcelSection `toml:"excel"`
	Smtp  SmtpSection  `toml:"smtp"`
	Debug DebugSection `toml:"debug"`
}

type MainSection struct {
	Template    string `toml:"template"`
	SenderEmail string `toml:"sender-email"`
	SenderName  string `toml:"sender-name"`
	SendEmail   bool   `toml:"send-email"`
}

type SmtpSection struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Auth     bool   `toml:"auth"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type ExcelSection struct {
	File        string `toml:"file"`
	EmailColumn string `toml:"email-column"`
	NameColumn  string `toml:"name-column"`
}

type DebugSection struct {
	SendTo    string `toml:"send-to"`
	Debugging bool   `toml:"debugging"`
}

func New(configFile string) Settings {
	settings := Settings{}
	if _, err := toml.DecodeFile(configFile, &settings); err != nil {
		log.Fatal(err)
	}
	return settings
}
