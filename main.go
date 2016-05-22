package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/Luzifer/rconfig"
)

var (
	cfg = struct {
		Hook        string `flag:"hook,h" default:"" description:"Slack incoming webhook"`
		Username    string `flag:"username,u" default:"CmdSlack" description:"Username to use in WebHook command"`
		Icon        string `flag:"icon,i" default:":package:" description:"Icon to use for the webhook"`
		Channel     string `flag:"channel,c" default:"" description:"Channel to send the message to"`
		Description string `flag:"description,d" default:"" description:"Add a piece of text to prepent to the output"`
	}{}

	version = "dev"
)

type slack struct {
	Username string `json:"username,omitempty"`
	Icon     string `json:"icon_emoji,omitempty"`
	Channel  string `json:"channel,omitempty"`
	Text     string `json:"text"`
}

func main() {
	rconfig.Parse(&cfg)
	cmdline := rconfig.Args()[1:]

	buf := bytes.NewBuffer([]byte{})
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command error: %s", err)
	}

	if strings.TrimSpace(buf.String()) == "" {
		log.Printf("Command had empty output, ignoring")
		os.Exit(0)
	}

	text := "```\n" + buf.String() + "```"
	if cfg.Description != "" {
		text = cfg.Description + "\n" + text
	}

	slo := slack{
		Username: cfg.Username,
		Icon:     cfg.Icon,
		Channel:  cfg.Channel,
		Text:     text,
	}

	body := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(body).Encode(slo); err != nil {
		log.Fatalf("Encoder error: %s", err)
	}

	if _, err := http.Post(cfg.Hook, "application/json", body); err != nil {
		log.Fatalf("Posting error: %s", err)
	}

	log.Printf("Posted successfully")
}
