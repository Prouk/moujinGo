package bin

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
	"strconv"
	"strings"
)

type Player struct {
	GuildId     int
	Interaction *discordgo.Interaction
	StartedBy   struct {
		Name string
		Icon string
	}
	Queue          []QueueItemInfos
	VoiceConection *discordgo.VoiceConnection
	Music          Music
}

type QueueItemInfos struct {
	Url        string
	MemberName string
	MessageId  string
}

func (p *Player) AddToQueue(url string, i *discordgo.Interaction) *Player {
	toQueue := QueueItemInfos{
		Url:        url,
		MemberName: i.Member.User.Username,
		MessageId:  i.ID,
	}
	p.Queue = append(p.Queue, toQueue)
	return p
}

func (p *Player) GetEmbed(action string) *discordgo.MessageEmbed {
	var desc string
	desc = action
	if len(p.Queue) > 1 {
		return &discordgo.MessageEmbed{
			Title:       p.Music.Title,
			URL:         p.Music.Url,
			Description: p.Music.Channel,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Current Music By",
					Value: p.Queue[0].MemberName,
				},
				{
					Name:  "Next Music By",
					Value: p.Queue[1].MemberName,
				},
				{
					Name:  "Music in queue",
					Value: strconv.Itoa(len(p.Queue)),
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Player Started By : " + p.StartedBy.Name,
				IconURL: p.StartedBy.Icon,
			},
		}
	} else {
		return &discordgo.MessageEmbed{
			Title:       "strconv.Itoa(player.CommandId)",
			URL:         p.Queue[0].Url,
			Description: desc,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Current Music By",
					Value: p.Queue[0].MemberName,
				},
				{
					Name:  "Music in queue",
					Value: strconv.Itoa(len(p.Queue)),
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Player Started By : " + p.StartedBy.Name,
				IconURL: p.StartedBy.Icon,
			},
		}
	}
}

func (p *Player) GetButtons() discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				CustomID: "pause",
				Label:    "Pause",
				Style:    discordgo.SuccessButton,
			},
			discordgo.Button{
				CustomID: "next",
				Label:    "Next",
				Style:    discordgo.PrimaryButton,
			},
			discordgo.Button{
				CustomID: "stop",
				Label:    "Stop",
				Style:    discordgo.DangerButton,
			},
			discordgo.Button{
				CustomID: "send",
				Label:    "Send",
				Style:    discordgo.SecondaryButton,
			},
		},
	}
}

func (p *Player) StreamFirstMusic() error {
	url := p.Queue[0].Url
	idIndex := strings.Index(url, "?v=")
	id := url[idIndex:]
	playlistIdIndex := strings.Index(id, "&list=")
	if playlistIdIndex >= 0 {
		id = id[:playlistIdIndex]
	}
	client := youtube.Client{}
	video, err := client.GetVideo(id)
	if err != nil {
		return err
	}
	formats := video.Formats.WithAudioChannels()
	_, _, err = client.GetStream(video, &formats[0])
	p.Music.Title = video.Title
	p.Music.Channel = video.Author
	p.Music.Thumbnail = video.Thumbnails[0].URL
	p.Music.Url = p.Queue[0].Url
	p.Music.AddedBy = p.Queue[0].MemberName
	return err
}
