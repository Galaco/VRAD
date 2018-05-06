package main

import (
	"github.com/galaco/vrad/cmd"
	"log"
	"github.com/galaco/vrad/cmd/tasks/start"
	"github.com/galaco/vrad/cmd/tasks/loadbsp"
	"github.com/galaco/vrad/cmd/tasks/computerad"
	"github.com/galaco/vrad/cmd/tasks/computerotherlighting"
	"github.com/galaco/vrad/cmd/tasks/finish"
)

func main() {
	command := cmd.NewCmd()
	command.AddStep(start.Main)
	command.AddStep(loadbsp.Main)
	command.AddStep(computerad.Main)
	command.AddStep(computerotherlighting.Main)
	command.AddStep(finish.Main)

	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
}