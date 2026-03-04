package main

import (
	"github.com/kclin/auto_notes/cmd"
)

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
