package slashcommands

import (
	"fmt"
	"template/config"
	"template/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

// StartTime stores when the bot was started
var StartTime = time.Now()

// Uptime calculates and shows how long the bot has been running since startup
func Uptime(s *discordgo.Session, i *discordgo.InteractionCreate) {
	now := time.Now()

	// this gets us years, months, days, hours, minutes, and seconds
	years := now.Year() - StartTime.Year()
	months := int(now.Month()) - int(StartTime.Month())
	days := now.Day() - StartTime.Day()
	hours := now.Hour() - StartTime.Hour()
	minutes := now.Minute() - StartTime.Minute()
	seconds := now.Second() - StartTime.Second()

	// now we handle negative values by borrowing from the next unit
	// this is like doing math on paper where you borrow from the next column and if that dont help explain it go back to math class
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	if days < 0 {
		// now we get the last day of the previous month so we can borrow properly
		prevMonth := now.AddDate(0, -1, 0)
		days += daysIn(prevMonth.Year(), prevMonth.Month())
		months--
	}
	if months < 0 {
		months += 12
		years--
	}

	// now we build the human-readable uptime string
	// we only want to get the units that are greater than 0 because we dont want to show "0 Days"
	var unit []string
	if years > 0 {
		if years == 1 {
			unit = append(unit, "1 Year")
		} else {
			unit = append(unit, fmt.Sprintf("%d Years", years))
		}
	}
	if months > 0 {
		if months == 1 {
			unit = append(unit, "1 Month")
		} else {
			unit = append(unit, fmt.Sprintf("%d Months", months))
		}
	}
	if days > 0 {
		if days == 1 {
			unit = append(unit, "1 Day")
		} else {
			unit = append(unit, fmt.Sprintf("%d Days", days))
		}
	}
	if hours > 0 {
		if hours == 1 {
			unit = append(unit, "1 Hour")
		} else {
			unit = append(unit, fmt.Sprintf("%d Hours", hours))
		}
	}
	if seconds > 0 {
		if seconds == 1 {
			unit = append(unit, "1 second")
		} else {
			unit = append(unit, fmt.Sprintf("%d seconds", seconds))
		}
	}

	// here we join the unit with commas and "and" for the last item
	// this makes it read like "3 Years, 5 Months and 2 Days" and so on
	var details string
	if len(unit) == 0 {
		details = "Just started" // nobody is using the bot this fast :P
	} else if len(unit) == 1 {
		details = unit[0]
	} else {
		details = fmt.Sprintf("%s and %s",
			joinWithCommas(unit[:len(unit)-1]),
			unit[len(unit)-1])
	}

	// here we format the main description with the date and detailed breakdown
	uptime := fmt.Sprintf("%s has been running since **%s** (%s)",
		config.Config.Brand.Name,
		StartTime.Format("January 2nd, 2006"),
		details)

	// now we build our embed
	embed := util.NewEmbed().
		SetTitle("Uptime").                                              // title
		SetDescription(uptime).                                          // uptime we calculated above
		SetColor(255, 255, 255).                                         // i like white better
		AddField("Started", fmt.Sprintf("<t:%d:F>", StartTime.Unix())).  // discord timestamp - shows exact start time
		AddField("Duration", fmt.Sprintf("<t:%d:R>", StartTime.Unix())). // discord relative timestamp - shows "x time ago"
		SetFooter(config.Config.Brand.Name, config.Config.Brand.Icon).   // bot branding in the footer
		Truncate()                                                       // auto-truncate to Discord limits

	// now we can send the response back to Discord
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,                                     // tell Discord were sending a message
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed}}, // now we send our embed

	})
}

// daysIn returns the number of days in a given month/year
// this is needed because months have different numbers of days
func daysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// joinWithCommas joins a slice of strings with commas
// so we can make nice lists like "1 Year, 2 Months, 3 Days"
func joinWithCommas(unit []string) string {
	if len(unit) == 0 {
		return ""
	}
	if len(unit) == 1 {
		return unit[0]
	}

	// start with the first part and add commas + spaces for the rest
	result := unit[0]
	for i := 1; i < len(unit); i++ {
		result += ", " + unit[i]
	}
	return result
}
