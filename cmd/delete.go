/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Scrape Config",
	Run: func(cmd *cobra.Command, args []string) {
		idToDelete, err := cmd.Flags().GetString("id")
		if err != nil {
			fmt.Println("bruh")
			log.Fatalln(err)
		}
		url := fmt.Sprintf("http://localhost:3030/elise/config/delete?id=%s", idToDelete)

		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			fmt.Println("kek")
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("kek3")
			log.Fatalln(err)
		}
		fmt.Println(string("her"))
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("id", "i", "", "ID Of Config To Delete")
	deleteCmd.MarkFlagRequired("id")
}
