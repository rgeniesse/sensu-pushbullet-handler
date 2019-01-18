package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sensu/sensu-go/types"
	"github.com/spf13/cobra"
	pushbullet "github.com/xconstruct/go-pushbullet"
)

var (
	apiToken string
	foo      string
	stdin    *os.File
	debug    bool
)

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-pushbullet-handler",
		Short: "The Sensu Go handler plugin for pushbullet",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&apiToken,
		"api.token",
		"a",
		os.Getenv("PUSHBULLET_APP_TOKEN"),
		"Pushbullet API app token, use default from PUSHBULLET_APP_TOKEN env var")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	if apiToken == "" {
		_ = cmd.Help()
		return fmt.Errorf("api token is empty")
	}

	eventJSON, err := ioutil.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %s", err)
	}

	event := &types.Event{}
	err = json.Unmarshal(eventJSON, event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stdin data: %s", err)
	}

	if err = event.Validate(); err != nil {
		return fmt.Errorf("failed to validate event: %s", err)
	}

	if !event.HasCheck() {
		return fmt.Errorf("event does not contain check")
	}

	return notifyPushbullet(event)
}

func notifyPushbullet(event *types.Event) error {
	pb := pushbullet.New(apiToken)
	devs, err := pb.Devices()
	if err != nil {
		panic(err)
	}

	title := event.Check.Name
	message := event.Check.Output

	err = pb.PushNote(devs[0].Iden, title, message)
	if err != nil {
		panic(err)
	}

	if debug == true {
		log.Println(err)
	}

	return nil
}
