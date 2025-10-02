package bot

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"template/bot/commands"
	"template/bot/slashcommands"
	"template/config"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yourpov/logrite"
)

// Start starts the Discord bot
func Start() {
	discord, err := discordgo.New("Bot " + config.Config.Token)
	if err != nil {
		logrite.Error("Failed to open a conn to Discord: %v", err)
		os.Exit(1)
	}

	// we need these intents to tell if a user is using a prefix command
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	if config.Config.PrefixEnabled {
		discord.AddHandler(messageCreate)
	}

	discord.AddHandler(ready)

	if config.Config.SlashEnabled {
		discord.AddHandler(handler)
	}

	err = discord.Open()
	if err != nil {
		logrite.Error("Failed to open conn: %v", err)
		os.Exit(1)
	}
	defer discord.Close()

	sc := make(chan os.Signal, 1)
	// we want to listen for termination signals to gracefully shutdown
	// this is especially important if you want to deregister commands on shutdown
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// we want to unload slash commands on shutdown
	// this is optional but recommended so users dont try to use commands when the bot is offline
	// here we unload slash commands
	if config.Config.SlashEnabled && config.Config.DeRegisterCommandsAfterRestart {
		logrite.Warn("Shutting down, deregistering commands...")
		slashcommands.Unload(discord)
	}

}

// ready is a handler for when the bot is ready
func ready(session *discordgo.Session, event *discordgo.Ready) {
	logrite.Info("Brand: %s", config.Config.Brand.Name)
	logrite.Info("User: %s (%s)", session.State.User.Username, session.State.User.ID)

	if config.Config.SlashEnabled && config.Config.PrefixEnabled {
		logrite.Info("Mode: Prefix/Slash")
	} else if config.Config.SlashEnabled {
		logrite.Info("Mode: Slash")
	} else if config.Config.PrefixEnabled {
		logrite.Info("Mode: Prefix")
	} else {
		logrite.Warn("Enable at least one command mode in config/config.json")
	}

	if config.Config.PrefixEnabled {
		// load our prefix commands if enabled
		commands.Load()
		//logrite.Success("Prefix Commands Loaded")
	}

	if config.Config.SlashEnabled {
		// load our slash commands if enabled
		slashcommands.Load(session)
		//logrite.Success("Slash Commands Loaded")
	}
}

// messageCreate is a handler for message-based commands
func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == session.State.User.ID {
		// we dont want the bot to read its own messages as commands so we
		return
	}

	args := strings.Split(m.Content, " ")

	if strings.HasPrefix(args[0], config.Config.Prefix) {
		stripped := strings.TrimPrefix(args[0], config.Config.Prefix)
		ok, command := commands.GetCommand(stripped, m)

		if !ok && command != nil {
			// here we send a temp message and delete it after 5s to mimic ephemeral since discordgo dont support ephem messages in normal text channels like clyde :(
			msg, _ := session.ChannelMessageSend(m.ChannelID, "You are not authorized to use this command")

			go func() {
				if msg != nil {
					<-time.After(5 * time.Second)
					session.ChannelMessageDelete(m.ChannelID, msg.ID)
				}
			}()
			return
		} else if ok {
			command.Execute(session, m, args)
		} else if !ok && command == nil {
			// we do the same thing here ^^
			notFound, _ := session.ChannelMessageSend(m.ChannelID, "Command not found")
			<-time.After(5 * time.Second)
			session.ChannelMessageDelete(m.ChannelID, notFound.ID)
		}
	}
}

// handler is a handler for slash commands
func handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// we get the command data from the interaction
	// and check if it exists in our registered commands map
	// if it does, we execute the command
	data := i.ApplicationCommandData()
	command, ok := slashcommands.Get(data.Name)

	if ok && command != nil {
		if command.Admin {
			if slashcommands.HasPermission(i) {
				command.Execute(s, i)
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You are not permitted to use this command",
						Flags:   discordgo.MessageFlagsEphemeral}})
			}
		} else {
			command.Execute(s, i)
		}
	}
}
