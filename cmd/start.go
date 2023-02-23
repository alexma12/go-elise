/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/alexma12/go-elise/config"
	"github.com/alexma12/go-elise/pkg/administrator"
	"github.com/alexma12/go-elise/pkg/db/mysql"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the scraper engine",
	Run: func(cmd *cobra.Command, args []string) {
		errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

		infoLog.Println("Initializing DB Connection...")
		mysqlDB, err := openMySQLDB()
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		defer mysqlDB.Close()
		db := mysql.New(mysqlDB)

		admin := administrator.New(db, errorLog, infoLog)
		admin.Start()
	},
}

func openMySQLDB() (*sql.DB, error) {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		return nil, errors.New("error occured while loading db config")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/elise?parseTime=true", dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBHost, dbConfig.DBPort)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func init() {
	rootCmd.AddCommand(startCmd)
}
