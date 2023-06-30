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

const banner = `
__  __  _____  ____         _____  ______  _____  ____  _    _ _____   _____ ______       _____ _____  ______       _______ ____  _____
|  \/  |/ ____|/ __ \       |  __ \|  ____|/ ____|/ __ \| |  | |  __ \ / ____|  ____|     / ____|  __ \|  ____|   /\|__   __/ __ \|  __ \
| \  / | (___ | |  | |______| |__) | |__  | (___ | |  | | |  | | |__) | |    | |__ ______| |    | |__) | |__     /  \  | | | |  | | |__) |
| |\/| |\___ \| |  | |______|  _  /|  __|  \___ \| |  | | |  | |  _  /| |    |  __|______| |    |  _  /|  __|   / /\ \ | | | |  | |  _  /
| |  | |____) | |__| |      | | \ \| |____ ____) | |__| | |__| | | \ \| |____| |____     | |____| | \ \| |____ / ____ \| | | |__| | | \ \
|_|  |_|_____/ \____/       |_|  \_|______|_____/ \____/ \____/|_|  \_\\_____|______|     \_____|_|  \_|______/_/    \_|_|  \____/|_|  \_\
`

func PrintBanner() {
	// Output banner + version output
	color.Cyan(banner)
	resp := goVersion.FuncWithOutput(false, version, commit, date, output)
	color.Magenta(resp + "\n")
}
