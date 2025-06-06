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
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"tsundoku.dev/quickcal/constants"
	"tsundoku.dev/quickcal/model"
)

var newCmdFlagAlarm []time.Duration
var newCmdFlagCalendar string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [description] [date] [time]",
	Short: "Adds a new event to the default calendar",
	Long: `
Creates a new event in the default calendar. The accepted parameters have some restrictions:

- description: only a single string is accepted, if there are spaces, surround the description in double quotes
- date: the format is either dd/mm (the current year is assumed) or dd/mm/yyyy
- time: the format is hh:mm
`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		eventComponent := ical.NewComponent(ical.CompEvent)

		// UID
		uid := uuid.NewString()
		eventComponent.Props[ical.PropUID] = []ical.Prop{
			{
				Name:   ical.PropUID,
				Params: nil,
				Value:  uid,
			},
		}

		// SUMMARY
		eventComponent.Props[ical.PropSummary] = []ical.Prop{
			{
				Name:   ical.PropSummary,
				Params: nil,
				Value:  args[0],
			},
		}

		// DTSTAMP
		tz, err := time.LoadLocation(cfg.Timezone)
		if err != nil {
			log.Println(err)
			return
		}

		dateParams := ical.Params{}
		dateParams.Set(ical.ParamTimezoneID, tz.String())

		eventComponent.Props[ical.PropDateTimeStamp] = []ical.Prop{
			{
				Name:   ical.PropDateTimeStamp,
				Params: dateParams,
				Value:  time.Now().Format(constants.TimeLayoutICalDateTime),
			},
		}

		// DTSTART
		dateStr := args[1]
		parts := strings.Split(dateStr, "/")
		if len(parts) == 2 {
			dateStr = fmt.Sprintf("%s/%d", dateStr, time.Now().Year())
		}

		var inputTime time.Time
		if len(args) == 2 {
			// full-day event
			inputTime, err = time.Parse(constants.TimeLayoutInputDate, dateStr)
			if err != nil {
				log.Println(err)
				return
			}

			singleDateParam := ical.Params{}
			singleDateParam.Set(ical.ParamValue, string(ical.ValueDate))

			eventComponent.Props[ical.PropDateTimeStart] = []ical.Prop{
				{
					Name:   ical.PropDateTimeStart,
					Params: singleDateParam,
					Value:  inputTime.Format(constants.TimeLayoutICalDate),
				},
			}

			eventComponent.Props[ical.PropDuration] = []ical.Prop{
				{
					Name:  ical.PropDuration,
					Value: "P1D",
				},
			}

		} else {
			// date and time
			timeStr := args[2]

			fullInputStr := fmt.Sprintf("%s %s", dateStr, timeStr)
			inputTime, err = time.Parse(constants.TimeLayoutInputDateTime, fullInputStr)
			if err != nil {
				log.Println(err)
				return
			}

			eventComponent.Props[ical.PropDateTimeStart] = []ical.Prop{
				{
					Name:   ical.PropDateTimeStart,
					Params: dateParams,
					Value:  inputTime.Format(constants.TimeLayoutICalDateTime),
				},
			}
		}

		// DTEND
		if len(args) > 2 {
			eventComponent.Props[ical.PropDateTimeEnd] = []ical.Prop{
				{
					Name:   ical.PropDateTimeEnd,
					Params: dateParams,
					Value:  inputTime.Add(1 * time.Hour).UTC().Format(constants.TimeLayoutICalDateTime),
				},
			}
		}

		if len(newCmdFlagAlarm) != 0 {

			for _, alarmDuration := range newCmdFlagAlarm {

				alarm, err := parseAlarm(alarmDuration, args[0])
				if err != nil {
					log.Println(err)
					return
				}

				eventComponent.Children = append(eventComponent.Children, alarm)
			}
		}

		newCalendar := ical.NewCalendar()
		// TODO: change name
		newCalendar.Props[ical.PropProductID] = []ical.Prop{
			{
				Name:  ical.PropProductID,
				Value: "-//QuickCal//CalDAV Client//EN",
			},
		}
		newCalendar.Props[ical.PropVersion] = []ical.Prop{
			{
				Name:  ical.PropVersion,
				Value: "2.0",
			},
		}
		newCalendar.Children = append(newCalendar.Children, eventComponent)

		var defaultCalendarClient *caldav.Client
		var defaultCalendar *model.Calendar
		for _, server := range caldavServers {
			for _, calendar := range server.Calendars {
				if calendar.Default {
					defaultCalendarClient = server.Client
					defaultCalendar = &calendar
					break
				}
			}
		}

		if defaultCalendarClient == nil || defaultCalendar == nil {
			log.Println("no default calendar found")
			return
		}

		path := fmt.Sprintf("%s%s.ics", defaultCalendar.Path, uid)

		calObject, err := defaultCalendarClient.PutCalendarObject(path, newCalendar)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(calObject.ETag)
	},
}

func init() {
	eventCmd.AddCommand(newCmd)

	newCmd.Flags().DurationSliceVarP(&newCmdFlagAlarm, "alarm", "a", nil, "Add an alarm (it can be used many times). Use h or m (e.g. --alarm 15m --alarm 1h, will create two alarms, one for 15 minutes and one for 1 hour before the event). ")
	newCmd.Flags().StringVarP(&newCmdFlagCalendar, "calendar", "c", "", "Set the calendar to write this event into. Overrides the selected default calendar")
}

func parseAlarm(alarmDuration time.Duration, alarmDescription string) (*ical.Component, error) {

	var ss strings.Builder

	if alarmDuration < 60*time.Minute {

		ss.WriteString("-PT")
		ss.WriteString(fmt.Sprintf("%d", int(alarmDuration.Minutes())))
		ss.WriteString("M")

	} else if alarmDuration < 24*time.Hour {

		ss.WriteString("-PT")
		ss.WriteString(fmt.Sprintf("%d", int(alarmDuration.Hours())))
		ss.WriteString("H")

	} else {

		ss.WriteString("-P")
		ss.WriteString(fmt.Sprintf("%d", int(alarmDuration.Hours()/24)))
		ss.WriteString("D")
	}

	alarm := ical.NewComponent(ical.CompAlarm)

	alarm.Props[ical.PropTrigger] = []ical.Prop{
		{
			Name:  ical.PropTrigger,
			Value: ss.String(),
		},
	}
	alarm.Props[ical.PropAction] = []ical.Prop{
		{
			Name:  ical.PropAction,
			Value: ical.ParamDisplay,
		},
	}
	alarm.Props[ical.PropDescription] = []ical.Prop{
		{
			Name:  ical.PropDescription,
			Value: alarmDescription,
		},
	}

	return alarm, nil
}
