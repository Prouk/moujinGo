package bin

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

type Moujin struct {
	BotSession *discordgo.Session
	Players    []*Player
	Logger     Logger
	Config     Config
}

func (m *Moujin) SetPlayer(i *discordgo.InteractionCreate, channel *discordgo.VoiceConnection) (*Player, error) {
	var p *Player
	var err error

	for index, player := range m.Players {
		if player.Interaction.GuildID == i.GuildID {
			p = m.Players[index]
		}
		p.VoiceConnection = channel
		return p, err
	}
	if p == nil {
		p = &Player{
			Interaction:     nil,
			VoiceConnection: channel,
			Queue:           nil,
			Initiator: struct {
				Name string
				Icon string
			}{
				Name: i.User.Username,
				Icon: i.User.AvatarURL("24px"),
			},
		}
	}
	return p, err
}

func (m *Moujin) GetPlayer(i *discordgo.InteractionCreate, url string) (*Player, error) {
	var p *Player
	var err error
	for index, player := range m.Players {
		if player.Interaction.GuildID == i.GuildID {
			p = m.Players[index]
		}
	}
	if p == nil {
		err = errors.New("Need to be in a voice channel first")
	}
	return p, err
}
