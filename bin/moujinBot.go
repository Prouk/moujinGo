package bin

import "github.com/bwmarrin/discordgo"

type Moujin struct {
	BotSession *discordgo.Session
	Presences  []Presence
	players    []Player
	Logger     Logger
	Config     Config
}

func (m *Moujin) GetPlayer(guildId int) (*Player, int) {
	if len(m.players) == 0 {
		return nil, -1
	}
	for index, player := range m.players {
		if player.GuildId == guildId {
			return &player, index
		}
	}
	return nil, -1
}

func (m *Moujin) AddGuildPlayer(guildId int, i *discordgo.Interaction) (*Player, error) {
	var err error
	player, err := CreatePlayer(guildId, i)
	m.players = append(m.players, player)
	return &m.players[len(m.players)-1], err
}

func CreatePlayer(guildId int, i *discordgo.Interaction) (Player, error) {
	var player Player
	var err error
	player.GuildId = guildId
	player.Interaction = i
	player.StartedBy.Name = i.Member.User.Username
	player.StartedBy.Icon = i.Member.User.AvatarURL("24px")
	if err != nil {
		return player, err
	}
	return player, err
}
