package bin

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

type Player struct {
	GuildId   int
	CommandId int
	StartedBy string
	Queue     []QueueItemInfos
}

type QueueItemInfos struct {
	Url    string
	Member string
}

func CreatePlayer(guildId int, i *discordgo.Interaction) (Player, error) {
	var player Player
	var err error
	player.GuildId = guildId
	player.StartedBy = i.Member.Nick
	commandId, err := strconv.Atoi(i.ApplicationCommandData().TargetID)
	if err != nil {
		return player, err
	}
	player.CommandId = commandId
	return player, err
}

func (p *Player) AddToQueue(url string, i *discordgo.Interaction) *Player {
	toQueue := QueueItemInfos{
		url,
		i.Member.Nick,
	}
	p.Queue = append(p.Queue, toQueue)
	return p
}
