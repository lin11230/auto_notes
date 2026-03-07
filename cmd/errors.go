package cmd

import (
	"fmt"
	"os"
)

func exitWithError(summary string, err error) {
	if debugMode && err != nil {
		fmt.Fprintf(os.Stderr, "錯誤：%s\n詳細資訊：%v\n", summary, err)
	} else {
		fmt.Fprintf(os.Stderr, "錯誤：%s\n", summary)
	}
	os.Exit(1)
}
