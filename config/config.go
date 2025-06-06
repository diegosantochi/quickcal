/*
QuickCal - A cli CalDAV client
Copyright (C) 2025 tsundoku.dev

This file is part of QuickCal.

QuickCal is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

QuickCal is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with QuickCal. If not, see <https://www.gnu.org/licenses/>.
*/

package config

type Calendar struct {
	Name    string
	Path    string
	Color   string
	Default bool
}
type Server struct {
	Name      string      `mapstructure:"name"`
	URL       string      `mapstructure:"url"`
	User      string      `mapstructure:"user"`
	Password  string      `mapstructure:"password"`
	Calendars []*Calendar `mapstructure:"calendars"`
}

type Config struct {
	Servers  []*Server `mapstructure:"servers"`
	Timezone string    `mapstructure:"timezone"`
}

func GetServerByName(cfg *Config, name string) *Server {
	for _, server := range cfg.Servers {
		if server.Name == name {
			return server
		}
	}
	return nil
}
