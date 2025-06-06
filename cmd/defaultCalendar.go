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

		/*if len(args) == 0 {

			if defaultCalendar == "" {
				fmt.Println("No default calendar set")
			} else {
				fmt.Printf("Default calendar: %s\n", defaultCalendar)
			}

			return
		}

		principal, err := caldavClient.FindCurrentUserPrincipal()
		if err != nil {
			fmt.Println(err)
		}

		homeset, err := caldavClient.FindCalendarHomeSet(principal)
		if err != nil {
			fmt.Println(err)
		}

		calendars, err := caldavClient.FindCalendars(homeset)
		if err != nil {
			fmt.Println(err)
		}

		for _, calendar := range calendars {
			if calendar.Path == args[0] {

				viper.Set("calendar.default", calendar.Path)

				err = viper.WriteConfig()
				if err != nil {
					fmt.Println(err)
				}

				return
			}
		}

		fmt.Println("No valid calendar found with that path")*/
	},
}

func init() {
	calendarCmd.AddCommand(defaultCalendarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// defaultCalendarCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// defaultCalendarCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
