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

// defaultCalendarCmd represents the defaultCalendar command
var defaultCalendarCmd = &cobra.Command{
	Use:   "default",
	Short: "Reads or sets the default calendar",
	Long: `
If no arguments are given, this command will output the current default calendar. If an argument is passed-in, the default calendar will be set to that path
`,
	Run: func(cmd *cobra.Command, args []string) {

		for _, cfgServer := range cfg.Servers {

			for _, cgfCalendar := range cfgServer.Calendars {
				if cgfCalendar.Default {
					fmt.Println()
					fmt.Println("Server:", cfgServer.Name)
					fmt.Println("Calendar name:", cgfCalendar.Name)
					fmt.Println("Calendar path:", cgfCalendar.Path)
					fmt.Println("Calendar color:", cgfCalendar.Color)
					fmt.Println()
					return
				}
			}
		}

		fmt.Println("No default calendar found")
	},
}

func init() {
	calendarCmd.AddCommand(defaultCalendarCmd)
}
