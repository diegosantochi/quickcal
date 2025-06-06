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
	"fmt"

	"github.com/spf13/cobra"
)

// listCalendarsCmd represents the listCalendars command
var listCalendarsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all calendars",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		for _, server := range caldavServers {

			caldavClient := server.Client
			principal, err := caldavClient.FindCurrentUserPrincipal()
			if err != nil {
				fmt.Printf("Error finding current user principal: %s\n", err)
				return
			}

			homeset, err := caldavClient.FindCalendarHomeSet(principal)
			if err != nil {
				fmt.Printf("Error finding calendar home set: %s\n", err)
				return
			}

			calendars, err := caldavClient.FindCalendars(homeset)
			if err != nil {
				fmt.Printf("Error finding calendars: %s\n", err)
				return
			}

			fmt.Println()
			for _, calendar := range calendars {
				fmt.Printf("Server: %s\nPath: %s\nName: %s\nDescription: %s\n\n", server.Name, calendar.Path, calendar.Name, calendar.Description)
			}
		}
	},
}

func init() {
	calendarCmd.AddCommand(listCalendarsCmd)
}
