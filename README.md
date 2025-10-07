<div align="center" id="top">

# Discord-Bot Template
</div>
<p align="center">
  <img alt="Top language" src="https://img.shields.io/github/languages/top/yourpov/Discord-Bot-TPL?color=56BEB8">
  <img alt="Language count" src="https://img.shields.io/github/languages/count/yourpov/Discord-Bot-TPL?color=56BEB8">
  <img alt="Repository size" src="https://img.shields.io/github/repo-size/yourpov/Discord-Bot-TPL?color=56BEB8">
  <img alt="License" src="https://img.shields.io/github/license/yourpov/Discord-Bot-TPL?color=56BEB8">
</p>

---

## About

**Discord-Bot-TPL** is a Go bot template that supports prefix commands and slash commands. It's modular with customizable embeds, configurable settings and has easy command management.

### ✨ Key Features

- **Hybrid Command System** - Supports both prefix commands (`.help`) and slash commands (`/uptime`)
- **Rich Embeds** - Custom embed utility for easy embed building (`SetTitle("str")`, `SetDescription("str")`, `SetColor(R,G,B)`) etc.
- **Configurable** - JSON-based configuration with customizable prefix, branding, and togglable commands
- **Permission System** - Discord-ID based "Admin" command access for sensitive commands
- **Modular** - Easy to extend with new commands and features

## 🛠️ Tech Stack

- **[Go](https://golang.org/)** - Language
- **[DiscordGo](https://github.com/bwmarrin/discordgo)** - Discord API wrapper for Go
- **JSON** - Configuration management

## 📋 Commands

### Prefix Commands (configurable, default: `.`)

- `.help` - Display all available commands
- `.ping` - Ping/pong response test

### Slash Commands

- `/uptime` - Show bot uptime and system information  
- `/test` - Test command for debugging

## 🚀 Setup

### Prerequisites

- Go 1.24.4 or higher
- Discord Bot Token
- Discord Server (for testing)

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/YourPOV/Discord-Bot-TPL
cd Discord-Bot-TPL
```

2. **Install dependencies**

```bash
go mod download
```

3. **Configure the bot**

Edit (**`config/config.json`**):

```json
{
    "token": "Put-Bot-Token-Here",
    "guild_id": "Put-Server-ID-Here", 
    "prefix": ".",
    
    "brand": {
        "name": "Your Brand Name",
        "icon": "https://avatars.githubusercontent.com/u/59181303?v=4"
    },
    
    "authenticated_ids": [
        "Put-Discord-ID-Here"
    ],

    "prefix_enabled": true,
    "slash_enabled": true,
    "deregister_commands_after_restart": false
}
```

### Configuration Options

| Option | Type | Description |
|--------|------|-------------|
| `token` | string | Your Discord bot token |
| `guild_id` | string | Discord server ID where slash commands are registered |
| `prefix` | string | Prefix for text commands (default: ".") |
| `brand.name` | string | Bot name displayed in embeds |
| `brand.icon` | string | Icon URL for embeds |
| `authenticated_ids` | array | Discord user IDs with admin command access |
| `prefix_enabled` | boolean | Enable/disable prefix commands |
| `slash_enabled` | boolean | Enable/disable slash commands |
| `deregister_commands_after_restart` | boolean | **Auto-remove slash commands when bot goes offline** |

#### 🔧 Command Deregistration Feature

When `deregister_commands_after_restart` is set to `true`:

- **Slash commands automatically disappear** from Discord when the bot shuts down
- **Prevents users from seeing non-functional commands** when the bot is offline
- Commands are re-registered automatically when the bot starts up again

1. **Run the bot**

```bash
go run main.go
```

## 🚀 Deployment (Linux)

For Deployments on Ubuntu/Debian servers, use the included `manage.sh` script:

### Quick Start

```bash
# Full deployment (build + restart + attach to screen)
bash manage.sh full

# Or step by step:
bash manage.sh build    # Build the bot
bash manage.sh restart  # Restart systemd service
bash manage.sh screen   # Attach to screen session
```

### Management Commands

| Command | Description |
|---------|-------------|
| `bash manage.sh status` | Check bot status and recent activity |
| `bash manage.sh logs` | View live bot logs (Ctrl+C to exit) |
| `bash manage.sh restart` | Restart the bot service |
| `bash manage.sh screen` | Create/attach to bot screen session |
| `bash manage.sh build` | Build bot from source |
| `bash manage.sh full` | Complete deployment |

### Features

- **Auto-Detection**: Automatically finds your Go project and checks if you have all the dependencie
- **Smart Building**: Cleans up old binaries and rebuilds everything with detailed output so you can see what's happening
- **Service Management**: Works with systemd to keep your bot running reliably (because nobody wants a bot that randomly dies)
- **Screen Sessions**: Runs your bot in the background so you can close your terminal and it stays running
- **Status Monitoring**: Shows you exactly what's happening with live status updates and error detection
- **Logging**: Keeps track of everything with timestamps so you can figure out what went wrong

---

## 🎨 Custom Embeds

The bot includes a custom embed utility:

```go
// Example
embed := util.NewEmbed().
    SetTitle("Embed Title").
    SetDescription("Embed description").
    SetColor(255, 100, 50). // RGB color values
    AddField("Field Name", "Field Value").
    SetFooter("Footer text", "footer icon URL").
    SetImage("image URL").
    SetThumbnail("thumbnail URL").
    SetAuthor("Author name", "author icon", "author URL").
    SetURL("https://example.com").
    InlineAllFields(). // Make all fields inline
    Truncate() // Auto-truncate to Discord limits
```

### Embed Methods

- `SetTitle(title)` - Set embed title
- `SetDescription(desc)` - Set embed description  
- `SetColor(r, g, b)` - Set color using RGB values
- `AddField(name, value)` - Add a field to the embed
- `SetFooter(text, iconURL, proxyURL)` - Set footer (variadic args)
- `SetImage(imageURL, proxyURL)` - Set main image (variadic args)
- `SetThumbnail(thumbURL, proxyURL)` - Set thumbnail (variadic args)
- `SetAuthor(name, iconURL, URL, proxyURL)` - Set author (variadic args)
- `SetURL(url)` - Set embed URL
- `InlineAllFields()` - Make all fields display inline
- `Truncate()` - Auto-truncate all content to Discord limits

### Utility Functions

- `NewGenericEmbed(title, message, ...args)` - Generic embed
- `NewErrorEmbed(title, message, ...args)` - Error embed  
- `NewErrorEmbedAdvanced(title, message, hexColor)` - Custom color error embed
---

## 🔧 Adding New Commands

### Prefix Command

```go
// Define in | bot/commands/NewCommand.go
func NewCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
    embed := util.NewEmbed().
        SetTitle("Example Command").
        SetDescription("I am a newly registered civi.. i mean command").
        SetColor(255, 255, 255)
    
    s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

// Register in | bot/commands/prefix_loader.go
newCommand(Command{
    Name:        "test",
    Description: "Command description",
    Execute:     NewCommand,
})
```

### Slash Command

```go
// Define in | bot/slashcommands/NewCommand.go
func NewCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
    embed := util.NewEmbed().
        SetTitle("Example Command").
        SetDescription("I am a newly registered civi.. i mean command").
        SetColor(255, 255, 255)
    
    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed}},
    })
}

// Add to cmds slice in | bot/commands/prefix_loader.go
    Name:        "test",
    Description: "Command description", 
    Type:        discordgo.ChatApplicationCommand,
    Admin:       true, // Only users in (config > authenticated_ids) can run this command
    Execute:     NewCommand,
}
```

### Slash Command with Options

```go
// For commands with parameters
{
    Name:        "say",
    Description: "Make the bot say something",
    Type:        discordgo.ChatApplicationCommand,
    Admin:       false, // anyone can run this command
    Options: []*discordgo.ApplicationCommandOption{
        {
            Type:        discordgo.ApplicationCommandOptionString,
            Name:        "message",
            Description: "Message to say",
            Required:    true,
        },
    },
    Execute:     SayCommand,
}
```

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/newFeature`)
3. Commit your changes (`git commit -m 'Add some newFeature'`)
4. Push to the branch (`git push origin feature/newFeature`)
5. Open a Pull Request

---

<div align="center">Made with ❤️ and Go</div>