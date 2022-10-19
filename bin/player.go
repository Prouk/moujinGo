package bin

import "github.com/bwmarrin/discordgo"

type Player struct {
	Interaction     *discordgo.Interaction
	VoiceConnection *discordgo.VoiceConnection
	Queue           []Music
	Initiator       struct {
		Name string
		Icon string
	}
}

func (p *Player) AddToQueue(url string, i *discordgo.InteractionCreate) (*Player, error) {
	var music Music
	var err error
	music = Music{
		Url:      url,
		Title:    "",
		Channel:  "",
		Duration: 0,
		AddedBy: struct {
			Name string
			Icon string
		}{
			Name: i.User.Username,
			Icon: i.User.AvatarURL("24px"),
		},
	}

	p.Queue = append(p.Queue, music)

	return p, err
}
func (p *Player) GetWaitingEmbed() *[]*discordgo.MessageEmbed {
	return &[]*discordgo.MessageEmbed{
		{
			Type:        "",
			Title:       "Waiting for a music",
			Description: "Player started by : " + p.Initiator.Name,
			Timestamp:   "",
			Color:       10181046,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Developed by Prouk.",
			},
			Image:     nil,
			Thumbnail: nil,
			Video:     nil,
			Provider:  nil,
			Author:    nil,
			Fields:    nil,
		},
	}
}

func (p *Player) GetEmbedComponents() *[]discordgo.MessageComponent {
	if len(p.Queue) == 0 {
		return &[]discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Quit",
				Style:    discordgo.DangerButton,
				Disabled: false,
				CustomID: "quit",
			},
		}
	} else {
		return &[]discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Pause / Play",
				Style:    discordgo.SuccessButton,
				Disabled: false,
				CustomID: "pause",
			},
			discordgo.Button{
				Label:    "Skip",
				Style:    discordgo.PrimaryButton,
				Disabled: false,
				CustomID: "skip",
			},
			discordgo.Button{
				Label:    "Quit",
				Style:    discordgo.DangerButton,
				Disabled: false,
				CustomID: "quit",
			},
		}
	}
}
