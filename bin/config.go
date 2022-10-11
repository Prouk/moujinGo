package bin

import (
	"encoding/json"
	"os"
)

type Config struct {
	BotToken    string `json:"bot_token"`
	BotTokenDev string `json:"bot_token_dev"`
	GuildId     string `json:"guild_id"`
	Volume      int    `json:"volume"`
	ToConsole   bool   `json:"to_console"`
	LogPath     string `json:"log_path"`
	LogLevel    int    `json:"log_level"`
}

func (Config) InitConf() (Config, error) {
	var config Config
	var err error
	file, err := os.ReadFile("./config.json")
	err = json.Unmarshal(file, &config)
	return config, err
}
