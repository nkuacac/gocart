package main

import (
	"fmt"

	"github.com/vito/gocart/command_runner"
	"github.com/vito/gocart/fetcher"
	"github.com/vito/gocart/set"
)

func install(root string, recursive bool) {
	cartridge, err := set.LoadFrom(root)
	if err != nil {
		fatal(err)
	}

	runner := command_runner.New(false)

	fetcher, err := fetcher.New(runner)
	if err != nil {
		fatal(err)
	}

	err = installDependencies(fetcher, cartridge, recursive, 0)
	if err != nil {
		fatal(err)
	}

	err = cartridge.SaveTo(root)
	if err != nil {
		fatal(err)
	}

	fmt.Println(green("OK"))
}

func installDependencies(fetcher *fetcher.Fetcher, deps *set.Set, recursive bool, depth int) error {
	maxWidth := 0

	for _, dep := range deps.Dependencies {
		if len(dep.Path) > maxWidth {
			maxWidth = len(dep.Path)
		}
	}

	for _, dep := range deps.Dependencies {
		versionDisplay := ""

		if dep.BleedingEdge {
			versionDisplay = "*"
		} else {
			versionDisplay = dep.Version
		}

		fmt.Println(
			indent(
				depth,
				bold(dep.Path)+padding(maxWidth-len(dep.Path)+2)+cyan(versionDisplay),
			),
		)

		lockedDependency, err := fetcher.Fetch(dep)
		if err != nil {
			return err
		}

		deps.Replace(lockedDependency)

		if recursive {
			nextDeps, err := set.LoadFrom(lockedDependency.FullPath(GOPATH))
			if err == set.NoCartridgeError {
				continue
			} else if err != nil {
				return err
			}

			err = installDependencies(fetcher, nextDeps, true, depth+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

