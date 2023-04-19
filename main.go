package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/donuts-are-good/colors"
	"github.com/matrix-org/gomatrix"
)

type Config struct {
	HomeserverURL  string `json:"homeserverURL"`
	UserID         string `json:"userID"`
	AccessToken    string `json:"accessToken"`
	AccessTokenCmd string `json:"accessTokenCmd"`
	RoomID         string `json:"roomID"`
}

const maxMessageSize = 4000
const messageDelay = 2 * time.Second

var (
	configFile = flag.String("config", "$XDG_CONFIG_HOME/neosay/config.json", "path to the configuration file")
	codeFlag   = flag.Bool("code", false, "is this code?")
)

func main() {
	flag.Parse()

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sfailed to get user's home directory. Well, shit.%s\n", colors.BrightRed, colors.Nc)
		os.Exit(1)
	}
	xdgConfigHome := path.Join(home, ".config")
	if env, _ := os.LookupEnv("XDG_CONFIG_HOME"); env != "" {
		xdgConfigHome = env
	}

	if strings.Contains(*configFile, "$XDG_CONFIG_HOME") {
		*configFile = filepath.Join(xdgConfigHome, "neosay", "config.json")
	}

	// open the config
	file, err := os.Open(*configFile)
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
	if config.AccessTokenCmd != "" {
		res := bytes.NewBuffer([]byte{})
		cmd := exec.Command("/bin/sh", "-c", config.AccessTokenCmd)
		cmd.Stdout = res
		cmd.Run()

		config.AccessToken = res.String()
	}
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
	if *codeFlag {
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
