/*
mcserv
Copyright (C) 2017 Joshua Lindsey

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
)

// MinecraftService implements the Service interface for the wrapped
// MC server.
type MinecraftService struct {
	// Path is the path to the MC server script to run
	Path string

	// InputChan is a buffered channel of inputs to send to the running server
	InputChan chan string

	// OuputChan is a buffered channel of server outputs in response to Inputs
	OutputChan chan string

	// OutputLines is a slice of all server output
	OutputLines []string
}

// NewMinecraftService instantiates a new MinecraftService pointer
func NewMinecraftService(path string) *MinecraftService {
	mc := MinecraftService{
		Path:        path,
		InputChan:   make(chan string, 5),
		OutputChan:  make(chan string, 5),
		OutputLines: make([]string, 0),
	}

	return &mc
}

func (mc *MinecraftService) String() string {
	return fmt.Sprintf("MinecraftService{}")
}
