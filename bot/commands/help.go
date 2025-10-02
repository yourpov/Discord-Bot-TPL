package commands

import (
	"fmt"
	"strings"
	"template/config"
	"template/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

// HelpPagination holds the data we need for the fancy little button pagination
// basically keeps track of what page we're on, what commands to show, and who asked for help
type HelpPagination struct {
	AllCommands   []string // all the regular user commands
	AdminCommands []string // the admin-only commands (for the power users)
	CurrentPage   int      // what page are we currently showing
	MaxPage       int      // how many pages do we have total
	ShowingAdmin  bool     // are we showing admin commands or regular ones
	UserID        string   // who asked for help (so only they can use the buttons)
	MessageID     string   // the message ID so we can track which help menu this is
}

// activePage keeps track of all the open help menus so buttons work properly
var activePage = make(map[string]*HelpPagination)

// handlerRegistered makes sure we only register the button handler once (nobody likes duplicates)
var handlerRegistered = false

/*
Parameters:
  - s (*discordgo.Session): the active Discord session instance
  - m (*discordgo.MessageCreate): the message that triggered the help command
  - args ([]string): command arguments (we dont really use these but they're there)

This creates a fancy paginated help menu with buttons that shows 10 commands per page
and lets users switch between regular and admin commands (yes this is possible with prefix commands too)
*/
func Help(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// we want to register our button handler if it isnt registered yet
	if !handlerRegistered {
		s.AddHandler(HandleHelpButtons)
		handlerRegistered = true
	}

	// we seperatew our command types into admin/regular
	var adminCommands []string
	var regularCommands []string

	// here we collect all our commands and organize them by type
	for _, cmd := range Commands {
		// now we format the command with prefix and aliases
		cmdText := fmt.Sprintf("`%s%s`", config.Config.Prefix, cmd.Name)
		// if there are aliases we add them in parentheses
		if len(cmd.Alias) > 0 {
			aliases := "" // otherwise leave them blank ^^
			for i, alias := range cmd.Alias {
				if i > 0 {
					aliases += ", " // we want to seperate them by a comma
				}
				aliases += fmt.Sprintf("`%s%s`", config.Config.Prefix, alias) // wrap them for (`.alias1`, `.alias2`)
			}
			cmdText += fmt.Sprintf(" (%s)", aliases)
		}
		cmdText += " - " + cmd.Description // add the description after the command

		// now we sort them into admin vs regular
		if cmd.AdminOnly {
			adminCommands = append(adminCommands, cmdText)
		} else {
			regularCommands = append(regularCommands, cmdText)
		}
	}

	// here we create our pagination object with all the info we need
	pagination := &HelpPagination{
		AllCommands:   regularCommands,
		AdminCommands: adminCommands,
		CurrentPage:   0,           // start on page 1 (well, page 0 but whos counting 0 as a page :p)
		ShowingAdmin:  false,       // start with regular commands
		UserID:        m.Author.ID, // remember who asked so only they can use buttons as we dont need other users trolling with it
	}

	// we wanna figure out how many pages we will need (i want 10 commands per page so i''ll do some math (ikr ew actual math is used in coding))
	commandsToShow := regularCommands
	pagination.MaxPage = (len(commandsToShow) - 1) / 10 // integer division to round down (so 0-9 = 1 page, 10-19 = 2 pages, etc)

	// create the initial embed and buttons
	embed, components := createHelpEmbed(pagination)

	// send the message with our fancy buttons
	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})

	if err == nil {
		// store the pagination data so our buttons know what to do
		pagination.MessageID = msg.ID
		activePage[msg.ID] = pagination

		// we gotta clean up after 5 minutes so we dont hoard memory like a damn squirrel stuffing nuts in its cheeks
		// this also prevents people from using old help menus forever
		// after 5 minutes the buttons will stop working and the menu will expire
		go func() {
			time.Sleep(5 * time.Minute)
			delete(activePage, msg.ID)
		}()
	}
}

// createHelpEmbed builds the embed and buttons for the current page
// this is where the magic happens for the pagination
func createHelpEmbed(p *HelpPagination) (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	// figure out which commands to show (regular or admin)
	commandsToShow := p.AllCommands
	sectionTitle := "General Commands"

	if p.ShowingAdmin {
		commandsToShow = p.AdminCommands
		sectionTitle = "Admin Commands"
		p.MaxPage = (len(commandsToShow) - 1) / 10 // recalculate pages for admin commands
	} else {
		p.MaxPage = (len(commandsToShow) - 1) / 10 // recalculate pages for regular commands
	}

	// make sure we don't go out of bounds (nobody likes array index errors) | engish translation: nobody likes crashes
	if p.CurrentPage > p.MaxPage {
		// if we go past the max page we set it to the last page
		p.CurrentPage = p.MaxPage
	}
	if p.CurrentPage < 0 {
		p.CurrentPage = 0
	}

	// here we get the commands for the specific page (10 per page)
	start := p.CurrentPage * 10
	end := start + 10
	if end > len(commandsToShow) {
		end = len(commandsToShow) // we dont wanna go don't go past the end of our list because that would be just silly :p
	}

	var pageCommands []string
	if start < len(commandsToShow) {
		pageCommands = commandsToShow[start:end]
	}

	// now we can build our fancy little embed
	embed := util.NewEmbed().
		SetTitle("Available Commands").
		SetDescription(fmt.Sprintf("Here are all the commands you can use with the `%s` prefix:", config.Config.Prefix)).
		SetColor(255, 255, 255)

	// add the commands field with pagination info
	if len(pageCommands) > 0 {
		commandsText := strings.Join(pageCommands, "\n")
		// this just adds the page number to the title so people know where they are (e.g: Page 1/3)
		embed.AddField(fmt.Sprintf("%s (Page %d/%d)", sectionTitle, p.CurrentPage+1, p.MaxPage+1), commandsText)
	}

	// footer with totals so people know how many commands there are
	totalRegular := len(p.AllCommands)
	totalAdmin := len(p.AdminCommands)
	embed.SetFooter(fmt.Sprintf("%s • %d general, %d admin commands", config.Config.Brand.Name, totalRegular, totalAdmin), config.Config.Brand.Icon).
		SetThumbnail(config.Config.Brand.Icon)

	// now we can make our fancy little buttons omgg
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: "help_prev",
					Label:    "◀ Previous", // took me abit to find this specific arrow for some reason
					Style:    discordgo.SecondaryButton,
					Disabled: p.CurrentPage == 0, // can't go back from first page
				},
				discordgo.Button{
					CustomID: "help_toggle",
					Label:    getToggleButtonText(p.ShowingAdmin),
					Style:    discordgo.PrimaryButton, // we want this one to stand out
				},
				discordgo.Button{
					CustomID: "help_next",
					Label:    "Next ▶",
					Style:    discordgo.SecondaryButton,
					Disabled: p.CurrentPage >= p.MaxPage, // this prevents us from going forward from last page
				},
				discordgo.Button{
					CustomID: "help_close",
					Label:    "✖ Close",
					Style:    discordgo.DangerButton, // DangerButton isnt actually dangerous but it its red so it works for our close
				},
			},
		},
	}

	return embed.MessageEmbed, components
}

// getToggleButtonText returns the label for the toggle button
// basically just tells people what they'll see if they click it
func getToggleButtonText(showingAdmin bool) string {
	if showingAdmin {
		return "General" // show this when we're on admin page (clicking goes to general)
	}
	return "Admin" // show this when we're on general page (clicking goes to admin)
}

// HandleHelpButtons handles button interactions for help pagination
func HandleHelpButtons(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// make sure it's one of our help buttons incase some other button interaction comes through (unlikely but still code is weird sometimes)
	customID := i.MessageComponentData().CustomID
	if !strings.HasPrefix(customID, "help_") {
		return
	}

	pagination, exists := activePage[i.Message.ID]
	if !exists {
		fmt.Printf("Pagination not found for message ID: %s\n", i.Message.ID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ This help menu has expired. Use `.help` to get a new one.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else if i.User != nil {
		userID = i.User.ID
	}

	// here we make sure only the person who asked for help can use the buttons
	if userID != pagination.UserID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Only the user who requested help can use these buttons.",
				Flags:   discordgo.MessageFlagsEphemeral, // an ephemeral message is that one hidden message clyde gives you when you try texting that ex that blocked you btw
			},
		})
		return
	}

	// now we handle the button actions
	switch customID {
	case "help_prev":
		pagination.CurrentPage-- // go back a page
	case "help_next":
		pagination.CurrentPage++ // go forward a page
	case "help_toggle":
		pagination.ShowingAdmin = !pagination.ShowingAdmin // switch between general and admin
		pagination.CurrentPage = 0                         // start from the first page when switching
	case "help_close":
		delete(activePage, i.Message.ID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    "Help menu closed",
				Flags:      discordgo.MessageFlagsEphemeral,
				Embeds:     []*discordgo.MessageEmbed{},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	// here we rebuild the embed with the new page
	embed, components := createHelpEmbed(pagination)

	// this updates the message with the new content
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}

// istg if anyone asked for help after i documented 200 lines in this whole source ima fucking lose it ()
// this help page may of taken me some work to finish
// but you have to work for nice things
// except donating  that's not much work at all | @AlisDiscord@icloud.com > donate if you like my work <3
