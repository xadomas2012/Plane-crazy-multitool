package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--update-helper" {
		if err := runUpdateHelper(os.Args[2:]); err != nil {
			_ = os.WriteFile(
				filepath.Join(
					os.TempDir(),
					"pc-multitool-update-error",
				),
				[]byte(err.Error()),
				0600,
			)
			os.Exit(1)
		}
		return
	}

	checkForUpdate()

	if _, err := tea.NewProgram(
		initialModel(),
	).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
