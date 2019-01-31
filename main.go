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
	token, device string
	stdin         *os.File
	allDevices    bool
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

	cmd.Flags().StringVarP(&token,
		"token",
		"t",
		os.Getenv("PUSHBULLET_APP_TOKEN"),
		"Pushbullet API app token, use default from PUSHBULLET_APP_TOKEN env var")

	cmd.Flags().StringVarP(&device,
		"device",
		"d",
		os.Getenv("DEVICE"),
		"A device registered with Pushbullet")

	cmd.Flags().BoolVar(&allDevices,
		"alldevices",
		false,
		"Bool for sending notifications to all devices")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	if token == "" {
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

	if allDevices && device == "" {

		if err = notifyPushbulletAll(event); err != nil {
			return errors.New(err.Error())
		}

	}

	if !allDevices && device != "" {

		if err = notifyPushbulletOne(event); err != nil {
			return errors.New(err.Error())
		}

	}

	if !allDevices && device == "" {
		fmt.Printf("%s\n", "Must chose to send notification to all devices or one")
	}

	if allDevices && device != "" {
		fmt.Printf("%s\n", "Can not send to all device and one device")
	}

	return nil
}

// Send event notification to all devices
func notifyPushbulletAll(event *types.Event) error {
	pushall := pushbullet.New(token)
	devs, err := pushall.Devices()
	if err != nil {
		panic(err)
	}

	title := fmt.Sprintf("%s/%s", event.Entity.Name, event.Check.Name)
	message := event.Check.Output

	for i := 0; i < len(devs); i++ {

		err = pushall.PushNote(devs[i].Iden, title, message)
		// Need to add something better here than a panic.
		// PB seems to hold onto old device Identities.
		// For now doing nothing seems ok
		if err != nil {
			// panic(err)
		}
	}

	return nil
}

func notifyPushbulletOne(event *types.Event) error {
	pushone := pushbullet.New(token)
	dev, err := pushone.Device(device)
	if err != nil {
		panic(err)
	}

	title := fmt.Sprintf("%s/%s", event.Entity.Name, event.Check.Name)
	message := event.Check.Output

	err = dev.PushNote(title, message)
	// Need to add something better here than a panic.
	// PB seems to hold onto old device Identities.
	// For now doing nothing seems ok
	if err != nil {
		// panic(err)
	}

	return nil
}

func validateEvent(event *types.Event) error {
	// Doesn't work for known sample events.
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
