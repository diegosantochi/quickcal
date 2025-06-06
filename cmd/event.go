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

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "events-related actions",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(eventCmd)
}
