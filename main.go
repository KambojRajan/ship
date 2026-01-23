package main

import (
	"fmt"
	"os"

	commands "github.com/KambojRajan/ship/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no command provided")
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("usage: ship add <path>")
			return
		}

		path := os.Args[2]
		err := commands.Add(path)
		if err != nil {
			return
		}
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("usage: ship init <path>")
			return
		}
		path := os.Args[2]
		err := commands.Init(path)
		if err != nil {
			return
		}
	case "cat-file":
		if len(os.Args) < 3 {
			fmt.Println("usage: ship cat-file <hash>")
			return
		}
		hash := os.Args[2]
		err := commands.CateFile(hash)
		if err != nil {
			return
		}
	}
}
