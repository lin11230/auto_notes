package cmd

import (
	"fmt"
	"os"
)

func exitWithError(summary string, err error) {
	if debugMode && err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\nDetails: %v\n", summary, err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", summary)
	}
	os.Exit(1)
}
