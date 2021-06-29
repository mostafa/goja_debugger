package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

var (
	dbg *goja.Debugger
	wg  sync.WaitGroup
)

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

	printer := console.PrinterFunc(func(s string) {
		fmt.Printf("< %s\n", s)
	})

	loader := func(requestedPath string) ([]byte, error) {
		if requestedPath != "" {
			fmt.Printf("Loaded sourcemap from: %s\n", requestedPath)
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

	wg.Add(2)

	go runtime.RunScript(filename, string(content))
	if err != nil {
		log.Fatal(err)
		os.Exit(4)
	}

	printWhyWeAreDebugging := func(b string) {
		if b == goja.BreakpointActivation {
			fmt.Println("hit breakpoint")
		} else {
			fmt.Println("hit debugger statement")
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

	wg.Wait()
}
