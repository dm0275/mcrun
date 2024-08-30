package main

import "github.com/dm0275/mcrun/cmd"

func main() {
	core := cmd.NewCLI()
	core.Execute()
}
