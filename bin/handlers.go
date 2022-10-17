package bin

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func SetHandlers(session *discordgo.Session, moujin *Moujin) {
	session.AddHandler(func(s *discordgo.Session, r discordgo.Ready) {
		moujin.Logger.PassLog("Bot Session initialized !", 0)
	})
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i, moujin)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i, moujin)
			}
		}
	})
}

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin){
		"play": func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin) {
			guildId, err := strconv.Atoi(i.GuildID)
			if err != nil {
				moujin.Logger.ErrorLog(err.Error(), 1)
			}
			options := i.ApplicationCommandData().Options
			url := options[0].StringValue()
			player, playerIndex := moujin.GetPlayer(guildId)
			if player == nil {
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "`Error joining channel`",
						},
					})
					return
				}
				player, err = moujin.AddGuildPlayer(guildId, i.Interaction)
				player.AddToQueue(url, i.Interaction)
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds:     []*discordgo.MessageEmbed{player.GetEmbed("Player Started")},
						Components: []discordgo.MessageComponent{player.GetButtons()},
					},
				})
				if err != nil {
					moujin.Logger.ErrorLog(err.Error(), 2)
				}
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "`Adding Song to Queue ...`",
					},
				})
				moujin.players[playerIndex].Queue = player.AddToQueue(url, i.Interaction).Queue
				_, err = s.InteractionResponseEdit(player.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{player.GetEmbed("Music Added")},
				})
				if err != nil {
					moujin.Logger.ErrorLog(err.Error(), 2)
				}
				s.InteractionResponseDelete(i.Interaction)
			}
		},
		"join": func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin) {
			options := i.ApplicationCommandData().Options
			voiceChannel := options[0].ChannelValue(s)
			channel, err := s.ChannelVoiceJoin(i.GuildID, voiceChannel.ID, false, true)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "`Error joining channel`",
					},
				})
				return
			}
			if len(moujin.Presences) == 0 {
				moujin.Presences = append(moujin.Presences, Presence{
					GuildId:         i.GuildID,
					VoiceConnection: channel,
				})
			} else {
				for index, presence := range moujin.Presences {
					if presence.GuildId == i.GuildID {
						moujin.Presences[index] = Presence{
							GuildId:         i.GuildID,
							VoiceConnection: channel,
						}
					}
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "`Channel " + voiceChannel.Name + " joined`",
				},
			})
			go func() {
				time.Sleep(5 * time.Second)
				s.InteractionResponseDelete(i.Interaction)
			}()
			return
		},
		"character": func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin) {
			options := i.ApplicationCommandData().Options
			dataCenter := options[0].StringValue()
			name := options[1].StringValue()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "`Searching character...`",
				},
			})
			ffCharacter, err := GetFfCharacter(dataCenter, name)
			if err != nil {
				moujin.Logger.ErrorLog(err.Error(),0)
			}
			emptyString := ""
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &emptyString,
				Embeds:     &[]*discordgo.MessageEmbed{&discordgo.MessageEmbed{
					Title: name,
					URL: ffCharacter.CharUrl,
					Description: ffCharacter.Title,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL:      ffCharacter.JobImg,
					},
					Image: &discordgo.MessageEmbedImage{
						URL:      ffCharacter.ImgUrl,
					},
					Fields:[]*discordgo.MessageEmbedField{
						{
							Name: "Grande compagnie",
							Value: ffCharacter.GrandCompany,
						},
						{
							Name: "Level: ",
							Value: ffCharacter.Level,
						},
					},
				}},
			})
			if err != nil {
				moujin.Logger.ErrorLog(err.Error(),0)
			}
		},
	}

	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin){
		"pause": func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin) {
			guildId, err := strconv.Atoi(i.GuildID)
			if err != nil {
				moujin.Logger.ErrorLog(err.Error(), 1)
			}
			player, _ := moujin.GetPlayer(guildId)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "`Pausing Music ...`",
				},
			})
			_, err = s.InteractionResponseEdit(player.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{player.GetEmbed("Music Paused")},
			})
			s.InteractionResponseDelete(i.Interaction)
		},
	}
)
