package ui

import "runtime"

const (
	ASCIIStylePulse = "pulse"
	ASCIIStyleOS    = "os"
	ASCIIStyleNone  = "none"
)

func ASCIIArt(style string) string {
	switch style {
	case ASCIIStyleNone:
		return ""
	case ASCIIStyleOS:
		return osArt()
	default:
		return pulseArt
	}
}

// pulse art for windows users lol (this is the default fallback)
const pulseArt = `
          ikiru

       /\        /\      /\
 _____/  \__/\__/  \____/  \____
          \/  \/

    your system, alive
`

func osArt() string {
	switch runtime.GOOS {
	case "linux":
		return linuxArt
	case "darwin":
		return macArt
	default:
		return pulseArt
	}
}

const linuxArt = "\n" +
	"        .--.\n" +
	"       |o_o |\n" +
	"       |:_/ |\n" +
	"      //   \\ \\\n" +
	"     (|     | )\n" +
	"    /'\\_   _/`\\\n" +
	"    \\___)=(___/\n"

const macArt = "\n" +
	"           .:'\n" +
	"       __ :'__\n" +
	"    .'`__`-'__``.\n" +
	"   :__________.-'\n" +
	"   :_________:\n" +
	"    :_________`-;\n" +
	"     `.__.-.__.'\n"
