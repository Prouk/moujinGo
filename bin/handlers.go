package bin

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

func SetHandlers(session *discordgo.Session, moujin *Moujin) {
	session.AddHandler(func(s *discordgo.Session, r discordgo.Ready) {
		moujin.Logger.PassLog("Bot Session initialized !", 0)
	})
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, moujin)
		}
	})
}

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin){
		"play": func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin) {
			moujin.Logger.CommentLog(`Play command requested`, 3)
			guildId, err := strconv.Atoi(i.GuildID)
			if err != nil {
				moujin.Logger.ErrorLog(err.Error(), 1)
			}
			options := i.ApplicationCommandData().Options
			url := options[0].StringValue()
			player := moujin.GetPlayer(guildId)
			if player == nil {
				player, err = moujin.AddGuildPlayer(guildId, i.Interaction)
				player.AddToQueue(url, i.Interaction)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title: strconv.Itoa(player.CommandId),
								URL:   url,
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:  "State",
										Value: "Launching the player",
									},
									{
										Name:  "Current Music By",
										Value: player.Queue[0].Member,
									},
								},
								Footer: &discordgo.MessageEmbedFooter{
									Text:    "Player Started By : " + player.StartedBy,
									IconURL: i.Interaction.Member.AvatarURL("25"),
								},
							},
						},
					},
				})
			} else {
				player.AddToQueue(url, i.Interaction)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title: strconv.Itoa(player.CommandId),
								URL:   url,
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:  "State",
										Value: "Launching the player",
									},
									{
										Name:  "Current Music By",
										Value: player.Queue[0].Member,
									},
									{
										Name:  "Next Music By",
										Value: player.Queue[1].Member,
									},
								},
								Footer: &discordgo.MessageEmbedFooter{
									Text:    "Player Started By : " + player.StartedBy,
									IconURL: i.Interaction.Member.AvatarURL("25"),
								},
							},
						},
					},
				})
			}
		},
	}
)
