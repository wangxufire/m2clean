package options

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

// Opt is command line args
var Args Option

func init() {
	_, err := flags.Parse(&Args)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
}
