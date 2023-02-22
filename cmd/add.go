/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var configToAdd = struct {
	Name              string
	Url               string
	Selector          string
	Type              int
	RequiresWebDriver bool
}{}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a scrape config",
	Run: func(cmd *cobra.Command, args []string) {
		configJson, err := json.Marshal(configToAdd)
		if err != nil {
			log.Fatalln(err)
		}
		resp, err := http.Post("http://localhost:3030/elise/config/add", "application/json", bytes.NewBuffer(configJson))
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Successfully created config: %s\n", string(b))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&configToAdd.Name, "name", "n", "", "Name of scrap configuration")
	addCmd.MarkFlagRequired("name")

	addCmd.Flags().StringVarP(&configToAdd.Url, "url", "u", "", "Url to accecss scrape contents")
	addCmd.MarkFlagRequired("url")

	addCmd.Flags().StringVarP(&configToAdd.Selector, "selector", "s", "", "Selector to accecss scrape contents")
	addCmd.MarkFlagRequired("selector")

	addCmd.Flags().IntVarP(&configToAdd.Type, "type", "t", 0, "The type of scrape contents:\n0 - number\n1 - string")
	addCmd.MarkFlagRequired("type")

	addCmd.Flags().BoolVarP(&configToAdd.RequiresWebDriver, "requiresWebDriver", "w", false, "Whether the scraper needs to use a web driver to access the contents")
	addCmd.MarkFlagRequired("requiresWebDriver")
}
