/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import "sopsctl/cmd"

// VERSION is set at build time via -ldflags "-X main.VERSION=<version>" (see .goreleaser.yaml).
var VERSION = "dev"

func main() {
	cmd.Execute(VERSION)
}
