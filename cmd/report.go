/*
Copyright © 2021 Josh Eveleth
This file is part of CLI application standupdate
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"

	"github.com/manifoldco/promptui"
)

var SlackToken = os.Getenv("SLACK_TOKEN")
var ChannelID = os.Getenv("CHANNEL_ID")

type TotalItems struct {
	Yesterday []string
	Today     []string
	Blockers  []string
}

const myRange = "{{range $val := .}}• {{$val}}\n{{end}}"

func ReturnFormattedUpdate(items []string) string {
	var buf bytes.Buffer
	t := template.Must(template.New("tmpl").Parse(myRange))
	err := t.Execute(&buf, items)
	FailErr("error executing", err)
	return buf.String()
}

func AddSection(info string) *slack.SectionBlock {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", info, false, false),
		nil,
		nil,
	)
}

func NotifySlack(items TotalItems) {
	api := slack.New(SlackToken)

	yHeaderSection := AddSection("*Yesterday*")
	yesterday := AddSection(ReturnFormattedUpdate(items.Yesterday))
	if items.Yesterday == nil {
		yesterday = AddSection(ReturnFormattedUpdate([]string{"Left Blank"}))
	}

	tHeaderSection := AddSection("*Today*")
	today := AddSection(ReturnFormattedUpdate(items.Today))
	if items.Today == nil {
		today = AddSection(ReturnFormattedUpdate([]string{"Left Blank"}))
	}
	bHeaderSection := AddSection("*Blockers*")
	blockers := AddSection(ReturnFormattedUpdate(items.Blockers))
	if items.Blockers == nil {
		blockers = AddSection(ReturnFormattedUpdate([]string{"Left Blank"}))
	}

	expectedBlocks := []slack.Block{
		yHeaderSection,
		yesterday,
		tHeaderSection,
		today,
		bHeaderSection,
		blockers,
	}
	_, _, err := api.PostMessage(ChannelID, slack.MsgOptionBlocks(expectedBlocks...))
	if err != nil {
		log.Printf("error is %v", err)
	}
}

var AllItems TotalItems

func SelectSection() (string, error) {
	prompt := promptui.Select{
		Label: "Select Section",
		Items: []string{"Yesterday", "Today", "Blockers", "Done"},
	}

	_, result, err := prompt.Run()
	FailErr("Prompt failed", err)
	return result, nil
}

func WholeShebang(result string) {
	switch result {
	case "Yesterday":
		_, _ = DoWorkFlow(result)
	case "Today":
		_, _ = DoWorkFlow(result)
	case "Blockers":
		_, _ = DoWorkFlow(result)
	case "Done":
		NotifySlack(AllItems)
	}
}

func DoWorkFlow(result string) (string, error) {
	newItem, err := AddItem(result)
	FailErr(fmt.Sprintf("Prompt failed for %v\n", result), err)
	switch result {
	case "Yesterday":
		AllItems.Yesterday = append(AllItems.Yesterday, newItem)
	case "Today":
		AllItems.Today = append(AllItems.Today, newItem)
	case "Blockers":
		AllItems.Blockers = append(AllItems.Blockers, newItem)
	}

	proceed, err := RunAgain(result)
	FailErr("Failed to run again", err)

	switch proceed {
	case "Yes":
		WholeShebang(result)
	case "No":
		result, err := SelectSection()
		FailErr("Error making selection", err)
		WholeShebang(result)
	}
	return result, nil
}

func RunAgain(result string) (string, error) {
	fmt.Printf("Another item for %v?\n", result)
	again := promptui.Select{
		Label: "Select Yes/No",
		Items: []string{"Yes", "No"},
	}
	_, result, err := again.Run()
	FailErr(fmt.Sprintf("Prompt failed for %v\n", result), err)
	return result, nil
}

func AddItem(result string) (string, error) {
	var item string
	fmt.Printf("What did you do %q?\n", result)
	input := func(input string) error {
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Item",
		Validate: input,
		Default:  item,
	}
	result, err := prompt.Run()
	FailErr("Prompt failed", err)
	return result, nil
}

func FailErr(errMsg string, err error) {
	if err != nil {
		log.Fatalf("%s : %v\n", errMsg, err)
	}
}

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := SelectSection()
		FailErr("Prompt failed", err)
		WholeShebang(result)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
