package commands

import (
	"os"
	"strings"
	"sync"
	"template/config"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/yourpov/logrite"
)

// Command struct is the structure for a prefix command
type Command struct {
	Name        string
	Alias       []string
	Description string
	AdminOnly   bool
	Execute     func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

// Commands map
var (
	Commands = make(map[string]*Command)
	lock     sync.Mutex
	cmds     = []Command{{
		Name:        "help",               // name of command
		Alias:       []string{"commands"}, // aliases of the command
		Description: "List all commands",  // description of the command
		AdminOnly:   false,                // is it admin only?
		Execute:     Help,                 // function to execute
	}, {
		Name:        "config",
		Alias:       []string{"configuration"},
		Description: "Check bot configuration",
		AdminOnly:   true,
		Execute:     CheckConfig,
	}, {
		Name:        "ping",
		Alias:       []string{"pingpong"},
		Description: "ping pong command",
		AdminOnly:   true,
		Execute:     PingPong,
	},
		// Test commands for pagination
		{Name: "avatar", Alias: []string{"av", "pfp"}, Description: "Get user's avatar", AdminOnly: false, Execute: TestCommand},
		{Name: "serverinfo", Alias: []string{"si", "guildinfo"}, Description: "Display server information", AdminOnly: false, Execute: TestCommand},
		{Name: "userinfo", Alias: []string{"ui", "whois"}, Description: "Display user information", AdminOnly: false, Execute: TestCommand},
		{Name: "ban", Alias: []string{"b"}, Description: "Ban a user from the server", AdminOnly: true, Execute: TestCommand},
		{Name: "kick", Alias: []string{"k"}, Description: "Kick a user from the server", AdminOnly: true, Execute: TestCommand},
		{Name: "mute", Alias: []string{"m", "timeout"}, Description: "Mute a user in the server", AdminOnly: true, Execute: TestCommand},
		{Name: "unmute", Alias: []string{"um"}, Description: "Unmute a user in the server", AdminOnly: true, Execute: TestCommand},
		{Name: "warn", Alias: []string{"w"}, Description: "Warn a user", AdminOnly: true, Execute: TestCommand},
		{Name: "purge", Alias: []string{"clear", "clean"}, Description: "Delete multiple messages", AdminOnly: true, Execute: TestCommand},
		{Name: "lock", Alias: []string{"lockdown"}, Description: "Lock a channel", AdminOnly: true, Execute: TestCommand},
		{Name: "unlock", Alias: []string{}, Description: "Unlock a channel", AdminOnly: true, Execute: TestCommand},
		{Name: "role", Alias: []string{"addrole"}, Description: "Add role to user", AdminOnly: true, Execute: TestCommand},
		{Name: "removerole", Alias: []string{"rr"}, Description: "Remove role from user", AdminOnly: true, Execute: TestCommand},
		{Name: "stats", Alias: []string{"statistics"}, Description: "Show bot statistics", AdminOnly: false, Execute: TestCommand},
		{Name: "weather", Alias: []string{"w"}, Description: "Get weather information", AdminOnly: false, Execute: TestCommand},
		{Name: "translate", Alias: []string{"tr"}, Description: "Translate text to another language", AdminOnly: false, Execute: TestCommand},
		{Name: "joke", Alias: []string{"funny"}, Description: "Get a random joke", AdminOnly: false, Execute: TestCommand},
		{Name: "quote", Alias: []string{"q"}, Description: "Get an inspirational quote", AdminOnly: false, Execute: TestCommand},
		{Name: "fact", Alias: []string{"f"}, Description: "Get a random fact", AdminOnly: false, Execute: TestCommand},
		{Name: "coinflip", Alias: []string{"flip", "coin"}, Description: "Flip a coin", AdminOnly: false, Execute: TestCommand},
		{Name: "dice", Alias: []string{"roll", "d6"}, Description: "Roll a dice", AdminOnly: false, Execute: TestCommand},
		{Name: "8ball", Alias: []string{"eightball"}, Description: "Ask the magic 8-ball", AdminOnly: false, Execute: TestCommand},
		{Name: "calculate", Alias: []string{"calc", "math"}, Description: "Perform calculations", AdminOnly: false, Execute: TestCommand},
		{Name: "base64", Alias: []string{"b64"}, Description: "Encode/decode base64", AdminOnly: false, Execute: TestCommand},
		{Name: "hash", Alias: []string{"md5", "sha1"}, Description: "Generate hash of text", AdminOnly: false, Execute: TestCommand},
		{Name: "qr", Alias: []string{"qrcode"}, Description: "Generate QR code", AdminOnly: false, Execute: TestCommand},
		{Name: "shorten", Alias: []string{"shorturl"}, Description: "Shorten a URL", AdminOnly: false, Execute: TestCommand},
		{Name: "screenshot", Alias: []string{"ss"}, Description: "Take website screenshot", AdminOnly: false, Execute: TestCommand},
		{Name: "color", Alias: []string{"colour", "hex"}, Description: "Show color information", AdminOnly: false, Execute: TestCommand},
		{Name: "reminder", Alias: []string{"remind", "timer"}, Description: "Set a reminder", AdminOnly: false, Execute: TestCommand},
		{Name: "poll", Alias: []string{"vote"}, Description: "Create a poll", AdminOnly: false, Execute: TestCommand},
		{Name: "giveaway", Alias: []string{"ga"}, Description: "Start a giveaway", AdminOnly: true, Execute: TestCommand},
		{Name: "announce", Alias: []string{"announcement"}, Description: "Make an announcement", AdminOnly: true, Execute: TestCommand},
		{Name: "embed", Alias: []string{"em"}, Description: "Create custom embed", AdminOnly: true, Execute: TestCommand},
		{Name: "say", Alias: []string{"echo"}, Description: "Make bot say something", AdminOnly: true, Execute: TestCommand},
		{Name: "react", Alias: []string{"r"}, Description: "Add reaction to message", AdminOnly: true, Execute: TestCommand},
		{Name: "slowmode", Alias: []string{"slow"}, Description: "Set channel slowmode", AdminOnly: true, Execute: TestCommand},
		{Name: "nickname", Alias: []string{"nick"}, Description: "Change user nickname", AdminOnly: true, Execute: TestCommand},
		{Name: "logs", Alias: []string{"log"}, Description: "View server logs", AdminOnly: true, Execute: TestCommand},
		{Name: "backup", Alias: []string{"save"}, Description: "Create server backup", AdminOnly: true, Execute: TestCommand},
		{Name: "restore", Alias: []string{"load"}, Description: "Restore server backup", AdminOnly: true, Execute: TestCommand},
		{Name: "automod", Alias: []string{"am"}, Description: "Configure automod settings", AdminOnly: true, Execute: TestCommand},
		{Name: "filter", Alias: []string{"wordfilter"}, Description: "Manage word filters", AdminOnly: true, Execute: TestCommand},
		{Name: "welcome", Alias: []string{"greeting"}, Description: "Configure welcome messages", AdminOnly: true, Execute: TestCommand},
		{Name: "goodbye", Alias: []string{"farewell"}, Description: "Configure goodbye messages", AdminOnly: true, Execute: TestCommand},
		{Name: "autorole", Alias: []string{"ar"}, Description: "Configure auto roles", AdminOnly: true, Execute: TestCommand},
		{Name: "starboard", Alias: []string{"star"}, Description: "Configure starboard", AdminOnly: true, Execute: TestCommand},
		{Name: "leveling", Alias: []string{"levels", "xp"}, Description: "Configure leveling system", AdminOnly: true, Execute: TestCommand},
		{Name: "economy", Alias: []string{"eco", "money"}, Description: "Economy system commands", AdminOnly: false, Execute: TestCommand},
		{Name: "shop", Alias: []string{"store"}, Description: "View the shop", AdminOnly: false, Execute: TestCommand},
		{Name: "inventory", Alias: []string{"inv", "items"}, Description: "View your inventory", AdminOnly: false, Execute: TestCommand},
		{Name: "balance", Alias: []string{"bal", "coins"}, Description: "Check your balance", AdminOnly: false, Execute: TestCommand},
		{Name: "daily", Alias: []string{"d"}, Description: "Claim daily reward", AdminOnly: false, Execute: TestCommand},
		{Name: "weekly", Alias: []string{"w"}, Description: "Claim weekly reward", AdminOnly: false, Execute: TestCommand},
		{Name: "work", Alias: []string{"job"}, Description: "Work to earn money", AdminOnly: false, Execute: TestCommand},
	}
)

// Load loads all commands into the Commands map
func Load() {
	for _, cmd := range cmds {
		newCommand(cmd)
		//logrite.Success("Registered Prefix Command: %s ", cmd.Name)
		logrite.Custom("⚙️ ", "COMMAND", "Registered prefix command: %s", color.FgWhite, color.BgGreen, cmd.Name)

	}
}

// newCommand adds a new command to the map
func newCommand(c Command) {
	lock.Lock()
	defer lock.Unlock()

	if existing, ok := Commands[c.Name]; ok {
		// here we check for conflicting command names
		logrite.Error("Conflicting command names: '%s' already exists", c.Name)
		logrite.Error("Existing command: %+v", existing)
		logrite.Error("New command: %+v", c)
		os.Exit(1)
	}

	Commands[c.Name] = &c
}

// GetCommand retrieves a command by name or alias, and checks if the user is authorized to use it
func GetCommand(name string, m *discordgo.MessageCreate) (bool, *Command) {
	lock.Lock()
	/* some people get confused with this but i do this to avoid deadlocks
	   i lock the mutex at the start of the function and use defer to unlock it when the function exits
	   this way it doesn't matter how many return statements we have since unlock will always run before the function actually returns */
	defer lock.Unlock()

	cmd, ok := Commands[name]
	// we want to check if the command exists first
	if ok && cmd != nil {
		if cmd.AdminOnly {
			// we also want to check if the user is authorized to use it
			for _, v := range config.Config.AuthenticatedIds {
				if strings.EqualFold(m.Author.ID, v) {
					// the user is in the authorized list for admin-only commands
					return true, cmd
				}
			}
			// we return false when the user's ID is not authorized
			return false, cmd
		}
		// this time we return true because all users are automatically authorized to use it
		return true, cmd
	}

	// we want to check for command aliases as well
	for _, cmd := range Commands {
		for _, alias := range cmd.Alias {
			if strings.EqualFold(name, alias) {
				if cmd.AdminOnly {
					// we also want to check if the user is authorized to use it
					for _, v := range config.Config.AuthenticatedIds {
						if strings.EqualFold(m.Author.ID, v) {
							// user is in the authorized list for admin-only commands
							return true, cmd
						}
					}
					// we return false when the user's ID is not authorized
					return false, cmd
				}
				// this time we return true because all users are automatically authorized to use it
				return true, cmd
			}
		}
	}
	// if we reach here the command was not found
	return false, nil
}

// TestCommand is a placeholder function for test commands
func TestCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, "🧪 This is a test command! It doesn't do anything yet.")
}
