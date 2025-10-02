package commands

import (
	"template/config"
	"template/util"

	"github.com/bwmarrin/discordgo"
)

/*
Parameters:
  - s (*discordgo.Session): the active Discord session instance
  - m (*discordgo.MessageCreate): the message that triggered the command
  - args ([]string): a slice of arguments passed to the command
*/

func PingPong(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

	// make an embed using our util package
	embed := util.NewEmbed().
		// methods are chainable so we can call them one after another
		SetTitle("Success").
		SetThumbnail("https://i.gifer.com/BH2F.gif").                  // add a thumbnail
		SetColor(255, 255, 255).                                       // SetColor takes RGB values
		SetDescription("Pong!").                                       // short description
		SetFooter(config.Config.Brand.Name, config.Config.Brand.Icon). // footer text + icon

		// extra fields for demonstration
		AddField("Field Name", "Field Value").
		SetImage("image URL").                                                            // large image below description
		SetAuthor("Author name", config.Config.Brand.Icon, "https://github.com/yourpov"). // author block
		SetURL("https://discord.com/developers/applications/1423165125044736021/oauth2"). // clickable embed link
		InlineAllFields().                                                                // make all fields inline
		Truncate()                                                                        // auto-truncate to Discord message limits

	// send the embed back to the same channel the command came from
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}
