package main

import (
	"encoding/json"
	"errors"
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
	stdin    *os.File
	debug    bool
)

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
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

	if apiToken == "" {
		_ = cmd.Help()
		return fmt.Errorf("api token is empty")
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	eventJSON, err := ioutil.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %s", err.Error())
	}

	event := &types.Event{}
	err = json.Unmarshal(eventJSON, event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stdin data: %s", err.Error())
	}

	if err = validateEvent(event); err != nil {
		return errors.New(err.Error())
	}

	if err = notifyPushbullet(event); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func notifyPushbullet(event *types.Event) error {
	pb := pushbullet.New(apiToken)
	devs, err := pb.Devices()
	if err != nil {
		panic(err)
	}

	title := fmt.Sprintf("%s/%s", event.Entity.Name, event.Check.Name)
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

func validateEvent(event *types.Event) error {
	// if event.Timestamp <= 0 {
	// 	return errors.New("timestamp is missing or must be greater than zero")
	// }

	if event.Entity == nil {
		return errors.New("entity is missing from event")
	}

	if !event.HasCheck() {
		return errors.New("check is missing from event")
	}

	if err := event.Entity.Validate(); err != nil {
		return err
	}

	if err := event.Check.Validate(); err != nil {
		return errors.New(err.Error())
	}

	return nil
}
