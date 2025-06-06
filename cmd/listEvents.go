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
	"log"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/spf13/cobra"
	"tsundoku.dev/quickcal/constants"
	"tsundoku.dev/quickcal/model"
)

var (
	fromDateStr string
	toDateStr   string
)

// listCmd represents the list command
var eventsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the events in the next 7 days",
	Long: `
Lists all the events, in the next 7 days, on the default calendar. A different calendar can be specified via flags.

The flags "from" and "to"" can be used to override the search time range.
`,
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		from := time.Now()
		if fromDateStr != "" {
			from, err = parseDateString(fromDateStr)
			if err != nil {
				log.Println(err)
				return
			}
		}

		to := from.Add(7 * 24 * time.Hour)
		if toDateStr != "" {
			to, err = parseDateString(toDateStr)
			if err != nil {
				log.Println(err)
				return
			}
		}

		query := caldav.CalendarQuery{
			CompFilter: caldav.CompFilter{
				Name: ical.CompCalendar,
				Comps: []caldav.CompFilter{
					{
						Name:  ical.CompEvent,
						Start: from,
						End:   to,
					},
				},
			},
		}

		var allEvents []*model.CalendarObject
		for _, caldavServer := range caldavServers {

			for _, calendar := range caldavServer.Calendars {

				calendarObjects, err := caldavServer.Client.QueryCalendar(calendar.Path, &query)
				if err != nil {
					fmt.Println(err)
				}

				for _, calendarObject := range calendarObjects {
					zcs, err := model.NewCalendarObjects(calendarObject, from, to, &calendar)
					if err != nil {
						fmt.Println(err)
					}

					allEvents = append(allEvents, zcs...)
				}
			}
		}

		sort.Slice(allEvents, func(i, j int) bool {
			return allEvents[i].Start.Before(*allEvents[j].Start)
		})

		for _, zc := range allEvents {
			_, _ = zc.Calendar.Color.Println(zc)
		}
	},
}

func init() {
	eventCmd.AddCommand(eventsListCmd)

	eventsListCmd.Flags().StringVar(&fromDateStr, "from", "", "List events from this date. Defaults to the current date")
	eventsListCmd.Flags().StringVar(&toDateStr, "to", "", "List events to this date. Defaults to 7 days after the from date")
}

func parseDateString(dateStr string) (time.Time, error) {

	// the date string can be dd/mm or dd/mm/yyyy
	parts := strings.Split(dateStr, "/")
	if len(parts) == 2 {
		dateStr = fmt.Sprintf("%s/%d", dateStr, time.Now().Year())
	}

	inputTime, err := time.Parse(constants.TimeLayoutInputDate, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return inputTime, nil
}
