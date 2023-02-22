/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexma12/go-elise/pkg/scrapedb"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scrape configs",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("http://localhost:3030/elise/config/list")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		var configs []scrapedb.ScrapeConfig

		err = decoder.Decode(&configs)
		if err != nil {
			log.Fatalln(err)
		}

		for _, c := range configs {
			fmt.Println(c)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
