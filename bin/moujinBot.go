package bin

import "github.com/bwmarrin/discordgo"

type Moujin struct {
	BotSession *discordgo.Session
	players    []Player
	Logger     Logger
	Config     Config
}

func (m *Moujin) GetPlayer(guildId int) *Player {
	if len(m.players) == 0 {
		return nil
	}
	for _, player := range m.players {
		if player.GuildId == guildId {
			return &player
		}
	}
	return nil
}

func (m *Moujin) AddGuildPlayer(guildId int, i *discordgo.Interaction) (*Player, error) {
	var err error
	player, err := CreatePlayer(guildId, i)
	m.players = append(m.players, player)
	return &m.players[len(m.players)-1], err
}
