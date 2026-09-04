package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	checkForUpdate()

	if _, err := tea.NewProgram(
		initialModel(),
	).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
