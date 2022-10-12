package bin

import "github.com/bwmarrin/discordgo"

type Presence struct {
	GuildId         string
	VoiceConnection *discordgo.VoiceConnection
}
