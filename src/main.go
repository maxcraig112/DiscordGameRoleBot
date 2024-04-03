package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	token          = flag.String("token", "", "Bot access token")
	AppID          = flag.String("app", "", "Application ID")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session
var ctx context.Context

// Parse input parameters
func init() { flag.Parse() }

// Get token from txt file (.gitignore :) ) and authenticate client
func init() {
	ctx = context.Background()

	var err error
	token, err := GetToken("token.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Got Authentication Token")

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Discord Client Created")
}

var (
	err error

	//List of all application commands currently available for giffy bot
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "add",
			Description: "Adds a role that can be mentioned",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "emoji",
					Description: "Emoji that you want associated with this role, leave blank if you want one randomly assigned",
					Required:    false,
				},
			},
		},
	}

	//list of functions that handle the interaction with the commands
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response
			option := optionMap["emoji"].StringValue()
			fmt.Println(option)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: option,
				},
			})
		},
	}

	//list of functions that handle the interaction with buttons
)

// This function takes a txt file of gifs as input and will attempt to process this all
// This has been used initially to migrate all gifs over from the previous giffy bot

func main() {

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Components are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	//delete all previously existing commands, this is to prevent commands no longer in use from being visible
	allCmds, _ := s.ApplicationCommands(*AppID, *GuildID)
	for _, v := range allCmds {
		fmt.Println(v.Name)
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	//add all existing commands
	log.Println("Adding commands...")
	for _, v := range commands {
		fmt.Println(v.Name)
		_, err := s.ApplicationCommandCreate(*AppID, *GuildID, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")

	allCmds, _ = s.ApplicationCommands(*AppID, *GuildID)
	for _, v := range allCmds {
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}
