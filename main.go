package main

import (
	"template/bot"
	"template/config"
)

// init loads the config and commands
func init() {
	config.Load()

}

// main starts the bot
func main() {
	bot.Start()
}
