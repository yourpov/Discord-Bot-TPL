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
