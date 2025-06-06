/*
QuickCal - A cli CalDAV client
Copyright (C) 2025 tsundoku.dev

This file is part of QuickCal.

QuickCal is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

QuickCal is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with QuickCal. If not, see <https://www.gnu.org/licenses/>.
*/

package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/teambition/rrule-go"
	"tsundoku.dev/quickcal/constants"
)

type CalendarObject struct {
	UID      string
	Summary  string
	Start    *time.Time
	End      *time.Time
	Calendar *Calendar
}

// NewCalendarObjects creates a list of CalendarObject from a caldav.CalendarObject.
// if the calendar object has a recurrence rule, the event will be expanded between the provided from and to times
func NewCalendarObjects(calendarObject caldav.CalendarObject, from time.Time, to time.Time, fromCalendar *Calendar) ([]*CalendarObject, error) {

	objects := make([]*CalendarObject, 0)
	if calendarObject.Data.Component.Name != ical.CompCalendar {
		return objects, errors.New("unexpected calendar object")
	}

	for _, child := range calendarObject.Data.Component.Children {
		if child.Name != ical.CompEvent {
			continue
		}

		// start time
		var startTime *time.Time
		startProp := child.Props.Get(ical.PropDateTimeStart)
		if startProp != nil {
			st, err := parseTime(*startProp)
			if err != nil {
				return objects, err
			}
			startTime = &st
		}
		if startTime == nil {
			continue
		}

		// end time
		var endTime *time.Time
		endProp := child.Props.Get(ical.PropDateTimeEnd)
		if endProp != nil {
			et, err := parseTime(*endProp)
			if err != nil {
				return objects, err
			}
			endTime = &et
		}

		// duration
		var eventDuration time.Duration
		if startTime != nil && endTime != nil {
			eventDuration = endTime.Sub(*startTime)
		}

		// uid and summary
		eventUID := child.Props.Get(ical.PropUID).Value
		eventSummary := child.Props.Get(ical.PropSummary).Value

		// recurrence rule
		rruleOption, err := child.Props.RecurrenceRule()
		if err != nil {
			return objects, err
		}

		// return only the event if it is a regular event within the given time range
		if rruleOption == nil && !startTime.Before(from) && !startTime.After(to) {
			objects = append(objects, &CalendarObject{
				UID:      eventUID,
				Summary:  eventSummary,
				Start:    startTime,
				End:      endTime,
				Calendar: fromCalendar,
			})
			return objects, err
		}

		// The `expand` extension is not implemented in github.com/emersion/go-webdav/caldav, so manually expanding is required
		// skip the event if it is not recurrent, and it is outside the time range
		if rruleOption == nil {
			continue
		}
		rruleOption.Dtstart = *startTime

		rRule, err := rrule.NewRRule(*rruleOption)
		if err != nil {
			return objects, err
		}

		occurrences := rRule.Between(from, to, true)
		for _, occurrence := range occurrences {
			if occurrence.Before(from) || occurrence.After(to) {
				continue
			}

			occuranceEvent := &CalendarObject{
				UID:      eventUID,
				Summary:  eventSummary,
				Start:    &occurrence,
				Calendar: fromCalendar,
			}

			// Calculate end time if original event had one
			if endTime != nil {
				occurranceEndTime := occurrence.Add(eventDuration)
				occuranceEvent.End = &occurranceEndTime
			}

			objects = append(objects, occuranceEvent)
		}

		return objects, nil
	}

	return objects, nil
}

func (z *CalendarObject) String() string {
	return fmt.Sprintf("%s\t%v\t%s", z.Calendar.Name, z.Start, z.Summary)
}

func parseTime(timeProp ical.Prop) (time.Time, error) {

	// The prop param can contain either "DATE" when only the date is specified (seen in recurrent events), or TZID, for full date/time objects
	tzid := timeProp.Params.Get(ical.PropTimezoneID)
	if tzid != "" {
		location, err := time.LoadLocation(tzid)
		if err != nil {
			return time.Time{}, err
		}

		parsedTime, err := time.ParseInLocation(constants.TimeLayoutICalDateTime, timeProp.Value, location)
		if err != nil {
			return time.Time{}, err
		}

		return parsedTime, nil
	}

	date := timeProp.Params.Get(ical.ParamValue)
	if date != "" {
		parsedTime, err := time.Parse(constants.TimeLayoutICalDate, timeProp.Value)
		if err != nil {
			return time.Time{}, err
		}

		return parsedTime, nil
	}

	return time.Time{}, errors.New("unrecognized timeProp type")
}
