package bin

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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
			options := i.ApplicationCommandData().Options
			url := options[0].StringValue()
			player, err := moujin.GetPlayer(i, url)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("`%s`", err.Error()),
					},
				})
				go DeleteInteractionResp(s, i, 5*time.Second)
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("`%s`", err.Error()),
				},
			})
			player.AddToQueue(url, i)
		},
		"join": func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin) {
			var err error
			options := i.ApplicationCommandData().Options
			voiceChannel := options[0].ChannelValue(s)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "`Joining the channel...`",
				},
			})
			moujin.Logger.WarningLog("am there", 0)
			channel, err := s.ChannelVoiceJoin(i.GuildID, voiceChannel.ID, false, true)
			moujin.Logger.WarningLog("am there 2", 0)
			if err != nil {
				moujin.Logger.WarningLog("wtf", 0)
				moujin.Logger.ErrorLog(err.Error(), 0)
				error := "`Error joining the channel...`"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &error,
				})
				go DeleteInteractionResp(s, i, 5*time.Second)
				return
			}
			moujin.Logger.WarningLog("am there 3", 0)
			p, err := moujin.SetPlayer(i, channel)
			moujin.Logger.WarningLog("am there 4", 0)
			if err != nil {
				moujin.Logger.ErrorLog(err.Error(), 0)
				error := "`Could not set the player`"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &error,
				})
				go DeleteInteractionResp(s, i, 5*time.Second)
				return
			}
			moujin.Logger.WarningLog("am there 5", 0)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds:     p.GetWaitingEmbed(),
				Components: p.GetEmbedComponents(),
			})
			moujin.Logger.WarningLog("am there 6", 0)
			if err != nil {
				moujin.Logger.ErrorLog(err.Error(), 0)
			}
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
				moujin.Logger.ErrorLog(err.Error(), 0)
			}
			emptyString := ""
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &emptyString,
				Embeds: &[]*discordgo.MessageEmbed{&discordgo.MessageEmbed{
					Title:       name,
					URL:         ffCharacter.CharUrl,
					Description: ffCharacter.Title,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: ffCharacter.JobImg,
					},
					Image: &discordgo.MessageEmbedImage{
						URL: ffCharacter.ImgUrl,
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Grande compagnie",
							Value: ffCharacter.GrandCompany,
						},
						{
							Name:  "Level: ",
							Value: ffCharacter.Level,
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Developed by Chanours.",
					},
				}},
			})
			if err != nil {
				moujin.Logger.ErrorLog(err.Error(), 0)
			}
		},
	}

	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin){
		"pause": func(s *discordgo.Session, i *discordgo.InteractionCreate, moujin *Moujin) {
		},
	}
)

func DeleteInteractionResp(s *discordgo.Session, i *discordgo.InteractionCreate, t time.Duration) {
	time.Sleep(t)
	s.InteractionResponseDelete(i.Interaction)
}
