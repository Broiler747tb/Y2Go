package main

import (
	"Y2Go/Cli"
	"Y2Go/UniPlayer"
)

func main() {
	path := CLI.GreeterAndSelecter()
	UniPlayer.Play(path)
}
