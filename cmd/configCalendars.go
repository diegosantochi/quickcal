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
	"errors"
	"fmt"

	"github.com/emersion/go-webdav/caldav"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tsundoku.dev/quickcal/config"
)

// calendarCmd represents the calendar command
var configCalendarCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure calendars",
	Long: `
Scans for calendars on the configured servers and adds them to the config file. It also allows you to set the default calendar.
`,
	Run: func(cmd *cobra.Command, args []string) {

		for _, server := range caldavServers {

			cfgServer := config.GetServerByName(&cfg, server.Name)

			principal, err := server.Client.FindCurrentUserPrincipal()
			if err != nil {
				fmt.Println(err)
			}

			homeset, err := server.Client.FindCalendarHomeSet(principal)
			if err != nil {
				fmt.Println(err)
			}

			calendars, err := server.Client.FindCalendars(homeset)
			if err != nil {
				fmt.Println(err)
			}

			for _, calendar := range calendars {

				err := processCalendar(cfgServer, calendar)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		viper.Set("servers", cfg.Servers)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
		}

		return
	},
}

func init() {
	calendarCmd.AddCommand(configCalendarCmd)
}

func processCalendar(server *config.Server, calendar caldav.Calendar) error {

	var cfgCalendar *config.Calendar
	for _, c := range server.Calendars {
		if c.Path == calendar.Path {
			cfgCalendar = c
			break
		}
	}

	fmt.Println()
	fmt.Println("Calendar's server:", server.Name)
	fmt.Println("Calendar's name:", calendar.Name)
	fmt.Println("Calendar's path:", calendar.Path)

	if cfgCalendar != nil {
		fmt.Println("Is currently tracked?: ✔")
	} else {
		fmt.Println("Is currently tracked? ❌")
	}
	fmt.Println()

	// keep or discard
	calendarAction := "Add"
	if cfgCalendar != nil {
		calendarAction = "Keep"
	}

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("%s this calendar?", calendarAction),
		IsConfirm: true,
	}

	result, err := prompt.Run()
	shouldKeepCalendar := err == nil && (result == "y" || result == "Y")

	if !shouldKeepCalendar {

		if cfgCalendar != nil {

			index := -1
			for i, c := range server.Calendars {
				if c.Path == cfgCalendar.Path {
					index = i
					break
				}
			}

			if index == -1 {
				return errors.New("invalid index")
			}

			server.Calendars = append(server.Calendars[:index], server.Calendars[index+1:]...)
		}

		return nil
	}

	if cfgCalendar == nil {
		cfgCalendar = &config.Calendar{
			Name: calendar.Name,
			Path: calendar.Path,
		}

		server.Calendars = append(server.Calendars, cfgCalendar)
	}

	// calendar color
	availableColors := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}
	templates := &promptui.SelectTemplates{
		Label:    `{{.}}`,
		Active:   "* {{ . }}",
		Inactive: "{{ . }}",
	}

	colorPrompt := promptui.Select{
		Label:     "Pick a color for the calendar",
		Items:     availableColors,
		Templates: templates,
		Size:      8,
	}

	if cfgCalendar.Color != "" {
		index := -1
		for i, c := range availableColors {
			if c == cfgCalendar.Color {
				index = i
				break
			}
		}

		if index != -1 {
			colorPrompt.CursorPos = index
		}
	}

	_, selectedColor, err := colorPrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}

	cfgCalendar.Color = selectedColor

	// is default?
	defaultPrompt := promptui.Prompt{
		Label:     "Do you want this calendar to be the default one?",
		IsConfirm: true,
	}

	result, err = defaultPrompt.Run()
	shouldBeDefault := err == nil && (result == "y" || result == "Y")

	if !shouldBeDefault {

		if cfgCalendar.Default {
			cfgCalendar.Default = false
		}

		return nil
	}

	for _, c := range server.Calendars {
		if c.Path == cfgCalendar.Path {
			c.Default = true
		} else {
			c.Default = false
		}
	}

	return nil
}
