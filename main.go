package main

import (
	"github.com/galaco/vrad/cmd"
	"log"
	"github.com/galaco/vrad/start"
	"github.com/galaco/vrad/loadbsp"
	"github.com/galaco/vrad/computerad"
	"github.com/galaco/vrad/computerotherlighting"
	"github.com/galaco/vrad/finish"
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