package utils

import (
	"fmt"
	"os"

	"my-personal-web/config"
)

func MigrateSeeders() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [up|down|refresh|create <name>|seed create <name>|seed run <name>|seed refresh]")
		return
	}

	switch os.Args[1] {
	// ==== Migration ====
	case "create":
		if len(os.Args) > 2 {
			config.RunMigration("create", os.Args[2])
		} else {
			fmt.Println("Usage: go run main.go create <migration_name>")
		}

	case "up", "down", "refresh":
		config.RunMigration(os.Args[1])

	// ==== Seeder ====
	case "seed":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go seed create <name> | seed run <name> | seed refresh")
			return
		}

		switch os.Args[2] {
		case "create":
			if len(os.Args) > 3 {
				config.RunSeeder("create", os.Args[3])
			} else {
				fmt.Println("Usage: go run main.go seed create <name>")
			}

		case "run":
			name := ""
			if len(os.Args) > 3 {
				name = os.Args[3]
			}
			config.RunSeeder("run", name)

		case "refresh":
			config.RunSeeder("refresh", "")

		default:
			fmt.Println("Usage: go run main.go seed create <name> | seed run <name> | seed refresh")
		}

	default:
		fmt.Println("Unknown command. Use: up | down | refresh | create <name> | seed ...")
	}
}
