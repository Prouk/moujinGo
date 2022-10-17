package main

import (
	"github.com/bwmarrin/discordgo"
	"main/bin"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var moujin bin.Moujin
	var err error
	moujin.Config, err = moujin.Config.InitConf()
	if err != nil {
		moujin.Logger.ErrorLog(err.Error(), 0)
		return
	}
	moujin.Logger, err = moujin.Logger.InitLogger(moujin.Config.ToConsole, moujin.Config.LogPath, moujin.Config.LogLevel)
	if err != nil {
		moujin.Logger.ErrorLog(err.Error(), 0)
		return
	}
	moujin.BotSession, err = discordgo.New("Bot " + moujin.Config.BotToken)
	if err != nil {
		moujin.Logger.ErrorLog(err.Error(), 0)
		os.Exit(1)
	}
	bin.SetHandlers(moujin.BotSession, &moujin)
	moujin.BotSession.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsDirectMessages | discordgo.IntentsGuildWebhooks | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences
	err = moujin.BotSession.Open()
	if err != nil {
		moujin.Logger.ErrorLog(err.Error(), 0)
		return
	}
	moujin.Logger.PassLog("Bot successfully started !", 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	err = moujin.BotSession.Close()
	if err != nil {
		moujin.Logger.ErrorLog(err.Error(), 0)
	}
	moujin.Logger.PassLog("Bot successfully stopped !", 0)

}
