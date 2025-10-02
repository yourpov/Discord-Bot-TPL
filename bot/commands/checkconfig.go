package commands

import (
	"fmt"
	"template/config"
	"template/util"

	"github.com/bwmarrin/discordgo"
)

// CheckConfig shows the current bot configuration in an embed
func CheckConfig(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	authenticatedIDs := config.Config.AuthenticatedIds
	Admins := ""
	for i, userID := range authenticatedIDs {
		if i > 0 {
			Admins += ", "
		}
		Admins += "<@" + userID + ">"
	}

	// could say on or off but emojis seem classic enough
	prefixStatus := "❌"
	if config.Config.PrefixEnabled {
		prefixStatus = "✅"
	}

	slashStatus := "❌"
	if config.Config.SlashEnabled {
		slashStatus = "✅"
	}

	deregisterStatus := "❌"
	if config.Config.DeRegisterCommandsAfterRestart {
		deregisterStatus = "✅"
	}

	embed := util.NewEmbed().
		SetTitle(fmt.Sprintf("⚙️ %s Configuration", config.Config.Brand.Name)).
		SetDescription("Bot Configuration and setand tings").
		AddField("Basic Settings", fmt.Sprintf("**Command Prefix:** `%s`\n**Brand Name:** `%s`", config.Config.Prefix, config.Config.Brand.Name)).
		AddField("Authenticated Users", Admins).
		AddField("Command Systems", fmt.Sprintf("**Prefix Commands:** %s\n**Slash Commands:** %s", prefixStatus, slashStatus)).
		AddField("Advanced Settings", fmt.Sprintf("**Auto-Deregister:** %s", deregisterStatus)).
		SetColor(255, 255, 255).
		SetFooter(fmt.Sprintf("%s • Configuration checked by %s", config.Config.Brand.Name, m.Author.Username), config.Config.Brand.Icon).
		SetThumbnail(config.Config.Brand.Icon).
		Truncate()
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}
