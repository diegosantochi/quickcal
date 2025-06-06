/*
QuickCal - A cli CalDAV client
Copyright (C) 2025 tsundoku.dev

This file is part of QuickCal.

QuickCal is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

QuickCal is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with QuickCal. If not, see <https://www.gnu.org/licenses/>.
*/

package cmd

import (
	"log"
	"os"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tsundoku.dev/quickcal/config"
	"tsundoku.dev/quickcal/model"
)

var cfgFile string
var caldavServers map[string]model.CalendarServer
var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "A cli caldav client.",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.calendar.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// init client
	cobra.OnInitialize(func() {
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Panicln("Failed to unmarshal config:", err)
		}

		caldavServers = make(map[string]model.CalendarServer)

		// Create a client for each server
		for _, server := range cfg.Servers {
			if server.URL == "" || server.User == "" || server.Password == "" {
				log.Printf("Skipping server '%s' due to missing configuration", server.Name)
				continue
			}

			httpClient := webdav.HTTPClientWithBasicAuth(nil, server.User, server.Password)

			client, err := caldav.NewClient(httpClient, server.URL)
			if err != nil {
				log.Printf("Failed to create client for server '%s': %v", server.Name, err)
				continue
			}

			calendars := make([]model.Calendar, 0, len(server.Calendars))
			for _, calendar := range server.Calendars {

				var calendarColor *color.Color
				switch calendar.Color {
				case "black":
					calendarColor = color.New(color.FgBlack)
				case "red":
					calendarColor = color.New(color.FgRed)
				case "green":
					calendarColor = color.New(color.FgGreen)
				case "yellow":
					calendarColor = color.New(color.FgYellow)
				case "blue":
					calendarColor = color.New(color.FgBlue)
				case "magenta":
					calendarColor = color.New(color.FgMagenta)
				case "cyan":
					calendarColor = color.New(color.FgCyan)
				case "white":
					calendarColor = color.New(color.FgWhite)
				}

				calendars = append(calendars, model.Calendar{
					Name:    calendar.Name,
					Path:    calendar.Path,
					Color:   calendarColor,
					Default: calendar.Default,
				})
			}

			caldavServers[server.Name] = model.CalendarServer{
				Name:      server.Name,
				Client:    client,
				Calendars: calendars,
			}
		}
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".calendar.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Failed to read config file:", err)
	}
}
