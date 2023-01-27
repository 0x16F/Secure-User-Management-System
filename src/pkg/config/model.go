package config

type Config struct {
	HTTP     HTTP     `json:"http"`
	Database Database `json:"database"`
	JWT      JWT      `json:"jwt"`
}

type HTTP struct {
	Port uint16 `json:"port"`
}

type Database struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Schema   string `json:"schema"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type JWT struct {
	AccessSecret  string `json:"access_secret"`
	RefreshSecret string `json:"refresh_secret"`
}
