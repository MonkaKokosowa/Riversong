package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (

	// SsrcUserMap is a map of SSRC to user IDs.
	SsrcUserMap = make(map[uint32]string)
	// UserSpeakingMap is a map of user IDs to their speaking state.
	UserSpeakingMap = make(map[string]*UserSpeaking)
	// Mutex to protect access to the maps.
	Mutex            = &sync.Mutex{}
	UserIDToUsername = make(map[string]string)
	ObsPassword      = ""
)

// UserSpeaking contains the speaking state and a timer for a user.
type UserSpeaking struct {
	Speaking bool
	Timer    *time.Timer
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) FetchMembers() map[string]string {
	return UserIDToUsername
}

func (a *App) StartBot(token string, voiceChannelID string, guildID string, obsPassword string) {
	ObsPassword = obsPassword
	// start discord bot client
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(ready)
	dg.AddHandler(guildCreate)

	dg.Identify.Intents = discordgo.IntentsAll

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	joinVoiceChannel(dg, voiceChannelID, guildID)

}

func voiceSpeakingUpdate(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	Mutex.Lock()
	defer Mutex.Unlock()
	SsrcUserMap[uint32(vs.SSRC)] = vs.UserID
}

func joinVoiceChannel(s *discordgo.Session, voiceChannelID, guildID string) {
	vc, err := s.ChannelVoiceJoin(guildID, voiceChannelID, true, false)
	if err != nil {
		log.Fatal(err)
	}

	vc.AddHandler(voiceSpeakingUpdate)

	// Start a goroutine to listen for Opus packets.
	go func() {
		for p := range vc.OpusRecv {
			Mutex.Lock()
			userID, ok := SsrcUserMap[p.SSRC]
			if !ok {
				Mutex.Unlock()
				continue
			}

			user, ok := UserSpeakingMap[userID]
			if !ok {
				log.Printf("User %s started speaking", UserIDToUsername[userID])
				user = &UserSpeaking{Speaking: true}
				UserSpeakingMap[userID] = user
			}

			if user.Timer != nil {
				user.Timer.Stop()
			}
			user.Timer = time.AfterFunc(200*time.Millisecond, func() {
				Mutex.Lock()
				defer Mutex.Unlock()
				log.Printf("User %s stopped speaking", UserIDToUsername[userID])
				delete(UserSpeakingMap, userID)
			})
			Mutex.Unlock()
		}
	}()

}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	log.Println("Bot is ready.")
}
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	for _, member := range event.Guild.Members {
		log.Printf("Guild member added: %s", member.User.Username)
		UserIDToUsername[member.User.ID] = member.User.Username
	}

}
