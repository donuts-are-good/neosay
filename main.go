package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/donuts-are-good/colors"
	"github.com/matrix-org/gomatrix"
)

type Config struct {
	HomeserverURL string `json:"homeserverURL"`
	UserID        string `json:"userID"`
	AccessToken   string `json:"accessToken"`
	RoomID        string `json:"roomID"`
}

const maxMessageSize = 4000
const messageDelay = 2 * time.Second

var codeFlag bool

func main() {

	// name the argset
	args := os.Args[1:]

	// if we get too few, supply the usage
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, colors.BrightRed+"Usage: neosay <config-file> [-c | --code]"+colors.Nc)
		os.Exit(1)
	}

	// name the config file
	configFile := args[0]

	// check for the args, which can be anywhere
	for _, arg := range args[1:] {
		if arg == "--code" || arg == "-c" {
			codeFlag = true
		}
	}

	// open the config
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, colors.BrightRed+"error opening config file:"+colors.Nc, err)
		os.Exit(1)
	}
	defer file.Close()

	// load the config into json
	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		fmt.Fprintln(os.Stderr, colors.BrightRed+"error reading config file:"+colors.Nc, err)
		os.Exit(1)
	}

	// get the homeserver
	homeserverURL := config.HomeserverURL
	if homeserverURL == "" {
		homeserverURL = os.Getenv("MATRIX_HOMESERVER_URL")
	}

	// get the user id
	userID := config.UserID
	if userID == "" {
		userID = os.Getenv("MATRIX_USER_ID")
	}

	// get the access token
	accessToken := config.AccessToken
	if accessToken == "" {
		accessToken = os.Getenv("MATRIX_ACCESS_TOKEN")
	}

	// get the room id
	roomID := config.RoomID
	if roomID == "" {
		roomID = os.Getenv("MATRIX_ROOM_ID")
	}

	// make a new client for matrix
	client, err := gomatrix.NewClient(homeserverURL, userID, accessToken)
	if err != nil {
		fmt.Fprintln(os.Stderr, colors.BrightRed+"error creating Matrix client:"+colors.Nc, err)
		os.Exit(1)
	}

	// make a buffer for msg chunking
	scanner := bufio.NewScanner(os.Stdin)
	buffer := ""
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		newBuffer := buffer + "\n" + text

		// if what we're getting is bigger
		// than max msg size, send what we have
		if len(newBuffer) > maxMessageSize {

			// send it
			sendMessage(client, roomID, buffer)
			buffer = text
		} else {

			// and if it isn't, add it
			buffer = newBuffer
		}
	}

	// keep talking till the buffer's empty
	if len(buffer) > 0 {
		sendMessage(client, roomID, buffer)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, colors.BrightRed+"error reading standard input:"+colors.Nc, err)
		os.Exit(1)
	}
}

func sendMessage(client *gomatrix.Client, roomID, message string) {

	// if we use -c or --code, make it a ```this thing```
	if codeFlag {
		message = "<pre><code>" + message + "</code></pre>"
	}

	// send the message to the room in tne cofing
	_, err := client.SendText(roomID, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, colors.BrightRed+"error sending message to Matrix:"+colors.Nc, err)
		os.Exit(1)
	} else {
		fmt.Println(colors.BrightGreen + "Sent!" + colors.Nc)
	}
	time.Sleep(messageDelay)
}
