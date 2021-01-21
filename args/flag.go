package args

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// Args is command line args
var Args Option

// Parse command line params
func Parse() {
	parser := flags.NewParser(&Args, flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash)
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(parser.Usage)
		os.Exit(-1)
	}
}
