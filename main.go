package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

var dbg *goja.Debugger

func main() {
	inspect := false
	filename := ""

	if len(os.Args) == 2 {
		filename = os.Args[1]
	} else if len(os.Args) == 3 {
		inspect = (os.Args[1] == "inspect")
		filename = os.Args[2]
	} else {
		fmt.Printf("Help:\n%s [inspect] <filename>\nUse inspect to enable debugging.\n", os.Args[0])
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	fmt.Println("Welcome to Goja debugger")
	fmt.Println("Type 'help' or 'h' for list of commands.")

	printer := console.PrinterFunc(func(s string) {
		fmt.Printf("< %s\n", s)
	})

	loader := func(requestedPath string) ([]byte, error) {
		if requestedPath != "" {
			fmt.Printf("%sLoaded sourcemap from: %s%s\n", GrayColor, requestedPath, ResetColor)
		}
		return nil, nil
	}

	_, err = parser.ParseFile(nil, filename, string(content), 0, parser.WithSourceMapLoader(loader))
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}

	runtime := goja.New()

	if inspect {
		dbg = runtime.EnableDebugMode()
	}

	registry := new(require.Registry)
	registry.Enable(runtime)
	registry.RegisterNativeModule("console", console.RequireWithPrinter(printer))
	console.Enable(runtime)

	printWhyWeAreDebugging := func(b string) {
		if b == goja.ProgramStartActivation {
			fmt.Printf("Break on start in %s:%d\n", dbg.Filename(), dbg.Line())
		} else if b == goja.BreakpointActivation {
			fmt.Printf("Break on breakpoint in %s:%d\ns", dbg.Filename(), dbg.Line())
		} else {
			fmt.Printf("Break on debugger statement in %s:%d\n", dbg.Filename(), dbg.Line())
		}
	}

	go func() {
		reader := bufio.NewReader(os.Stdin)

		b, c := dbg.WaitToActivate()
		printWhyWeAreDebugging(b)
		for {
			fmt.Printf("debug[%d]> ", dbg.GetPC())
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)
			if !executor(text) {
				c()
				b, c = dbg.WaitToActivate()
				printWhyWeAreDebugging(b)
			}
		}
	}()

	runtime.RunScript(filename, string(content))
	if err != nil {
		log.Fatal(err)
		os.Exit(4)
	}
}
