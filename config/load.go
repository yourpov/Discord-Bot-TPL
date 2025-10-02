package config

import (
	"encoding/json"
	"os"

	"github.com/yourpov/logrite"
)

type BrandConfig struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type cfg struct {
	// Token is the bot token from the Discord Developer Portal
	Token string `json:"token"`
	// Prefix is the command prefix for prefix commands
	Prefix string `json:"prefix"`
	// Brand contains branding information for the bot, such as name and icon URL
	Brand BrandConfig `json:"brand"`
	// GuildID is the ID of the guild (server) to register commands in, leave empty to register globally
	GuildID string `json:"guild_id"`
	// AuthenticatedIds is a list of user IDs that are authorized to use admin-only commands
	AuthenticatedIds []string `json:"authenticated_ids"`
	// PrefixEnabled when true, enables prefix commands
	PrefixEnabled bool `json:"prefix_enabled"`
	// SlashEnabled when true, enables slash commands
	SlashEnabled bool `json:"slash_enabled"`
	// DeRegisterCommandsAfterRestart when true, removes all slash commands from Discord when the bot shuts down
	DeRegisterCommandsAfterRestart bool `json:"deregister_commands_after_restart"`
}

var (
	// Config holds the bot configuration, loaded from config.json so do (config.) when anything from it
	Config *cfg
)

// Load reads the configuration file and unmarshals it into the Config variable
func Load() {
	f, err := os.ReadFile("./config/config.json")
	if err != nil {
		// if we cant read the file, we want to stop
		logrite.Error("Failed to read config file: %v", err)
		os.Exit(1)
	}

	if err = json.Unmarshal(f, &Config); err != nil {
		// if we cant parse the file, we also want to stop
		logrite.Error("Failed to parse config file: %v", err)
		// we exit as we cant run without a valid config you silly silly person you kek
		os.Exit(1)
	}
}
