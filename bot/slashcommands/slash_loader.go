package slashcommands

import (
	"os"
	"strings"
	"sync"
	"template/config"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/yourpov/logrite"
)

// RegisteredCommands stores the actual registered commands from Discord
var RegisteredCommands []*discordgo.ApplicationCommand

// RegisteredCommandIDs stores command IDS registered with Discord
var RegisteredCommandIDs []string

// Command struct is the structure for a slash command
type Command struct {
	ID          string
	Name        string
	Description string
	Type        discordgo.ApplicationCommandType
	Options     []*discordgo.ApplicationCommandOption
	Admin       bool
	Execute     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Commands map for slash commands
var (
	Commands = make(map[string]*Command)
	lock     sync.Mutex
)

// cmds holds all slash commands to be registered
var cmds = []Command{
	{
		Name:        "test",                           // name of command
		Description: "test command",                   // description of command
		Type:        discordgo.ChatApplicationCommand, // type of command
		Admin:       true,                             // is it admin only?
		Execute:     PingPong,                         // function to execute
	}, {
		Name:        "uptime",
		Description: "Show bot uptime",
		Type:        discordgo.ChatApplicationCommand,
		Admin:       false,
		Execute:     Uptime,
	},
}

// newCommand registeres a command in cmds
func newCommand(c Command) {
	if _, ok := Commands[c.Name]; ok {
		logrite.Error("Conflicting Slash Commands")
		os.Exit(1)
	}
	lock.Lock()
	defer lock.Unlock()
	Commands[c.Name] = &c
}

// Load registers all commands with Discord
func Load(s *discordgo.Session) {
	// here we clear any existing data on startup to prevent duplicates
	// this way if the bot crashes or gets restarted, we dont duplicate commands
	RegisteredCommands = nil
	RegisteredCommandIDs = nil

	for _, v := range cmds {
		newCommand(v)
	}

	for _, cmd := range cmds {
		// here we create a discordgo.ApplicationCommand from our Command struct above
		// this is what we actually register with Discord
		// i did this so we can have our own struct with an Execute function
		discordCmd := &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Type:        cmd.Type,
			Description: cmd.Description,
			Options:     cmd.Options,
		}

		// now we want to register our commands with Discord
		// we store the registered commands and their IDs so we can deregister them later if configured
		registeredCmd, err := s.ApplicationCommandCreate(s.State.User.ID, config.Config.GuildID, discordCmd)
		if err != nil {
			// we shouldnt reach here but just in case (i like my logs clean..)
			logrite.Error("Cannot create '%v' command: %v", cmd.Name, err)
			continue
		}

		// here we are just appending the registered command to our slices
		// so we can deregister them later if needed
		RegisteredCommands = append(RegisteredCommands, registeredCmd)
		RegisteredCommandIDs = append(RegisteredCommandIDs, registeredCmd.ID)
		logrite.Custom("⚙️ ", "COMMAND", "Registered slash command: %s", color.BgGreen, color.FgBlack, registeredCmd.Name)
	}
}

// Unload deregisters all commands from Discord
func Unload(s *discordgo.Session) {
	if !config.Config.DeRegisterCommandsAfterRestart {
		logrite.Info("Command deregistration is disabled")
		return
	}

	for i, cmdID := range RegisteredCommandIDs {
		err := s.ApplicationCommandDelete(s.State.User.ID, config.Config.GuildID, cmdID)
		if err != nil {
			// only way this would fail is if the command ID is invalid or Discord is having issues
			logrite.Error("Failed to delete command ID %s: %v", cmdID, err)
		} else {
			if i < len(RegisteredCommands) {
				logrite.Custom("⚙️ ", "COMMAND", "Deregistered: %s", color.BgGreen, color.FgWhite, RegisteredCommands[i].Name)
			} else {
				logrite.Custom("⚙️ ", "COMMAND", "Deregistered: %s", color.BgGreen, color.BgBlack, cmdID)
			}
		}
	}

	// we clear the slices after deregistering
	// this is to prevent trying to deregister them again if the bot is restarted
	// also helps with memory management
	RegisteredCommands = nil
	RegisteredCommandIDs = nil
}

// Get retrieves a command by name
func Get(cmd string) (*Command, bool) {
	for _, v := range Commands {
		// case insensitive comparison so users can use any case (help, Help, HELP, HeLp) and so on
		if strings.EqualFold(v.Name, cmd) {
			return v, true
		}
	}
	// we cant find the command so we return nil and false
	return nil, false
}

// HasPermission is what we use to check if the user has permission to use the command
func HasPermission(interaction *discordgo.InteractionCreate) bool {
	for _, v := range config.Config.AuthenticatedIds {
		// interaction.Member.User.ID is the discord id of the user who used the command if it wasnt obvious (i definitely didnt add this line so this page was an even 150 lines..)
		if v == interaction.Member.User.ID {
			// well return true since the user is authorized
			return true
		}
	}
	// now well return false because if we reached here they arent authorized

	return false
}
