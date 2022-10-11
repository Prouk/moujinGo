package bin

import "github.com/bwmarrin/discordgo"

type Moujin struct {
	BotSession *discordgo.Session
	players    []Player
	Logger     Logger
	Config     Config
}

func (m *Moujin) GetPlayer(guildId int) *Player {
	return &m.players[guildId]
}

func (m *Moujin) AddGuildPlayer(guildId int, i *discordgo.Interaction) (*Player, error) {
	var err error
	m.players[guildId], err = CreatePlayer(guildId, i)
	return &m.players[guildId], err
}
