package slashcommands

import (
	"template/config"
	"template/util"

	"github.com/bwmarrin/discordgo"
)

/*
Parameters:
  - s (*discordgo.Session): the active Discord session instance
  - i (*discordgo.InteractionCreate): the slash command interaction that triggered this function
*/

func PingPong(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// make an embed using our util package
	embed := util.NewEmbed().
		// methods are chainable so we can call them one after another
		SetTitle("Pong").                                              // simple title
		SetColor(255, 255, 255).                                       // SetColor takes RGB values
		SetFooter(config.Config.Brand.Name, config.Config.Brand.Icon). // footer with bot branding
		Truncate()                                                     // auto-truncate to Discord message limits

	// slash commands respond differently than regular messages
	// we use InteractionRespond instead of ChannelMessageSend
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource, // this tells Discord were sending a message back
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed}, // slash commands need embeds in a slice format
		},
	})
}
