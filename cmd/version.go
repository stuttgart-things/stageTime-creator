/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package cmd

import (
	"github.com/fatih/color"
	goVersion "go.hein.dev/go-version"
)

var (
	version = "unset"
	date    = "unknown"
	commit  = "unknown"
	output  = "yaml"
)

// https://fsymbols.com/generators/carty/
const banner = `
█▀ █░█░█ █▀▀ ▄▀█ ▀█▀ █▀ █░█ █▀█ █▀█ ▄▄ █▀▀ █▀█ █▀▀ ▄▀█ ▀█▀ █▀█ █▀█
▄█ ▀▄▀▄▀ ██▄ █▀█ ░█░ ▄█ █▀█ █▄█ █▀▀ ░░ █▄▄ █▀▄ ██▄ █▀█ ░█░ █▄█ █▀▄

`

func PrintBanner() {
	// Output banner + version output
	color.Cyan(banner)
	resp := goVersion.FuncWithOutput(false, version, commit, date, output)
	color.Magenta(resp + "\n")
}
