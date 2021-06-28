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
	dbg   *goja.Debugger
	dbgCh <-chan *goja.Debugger
	wg    sync.WaitGroup
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

	var lastPrint string
	printer := console.PrinterFunc(func(s string) {
		lastPrint = s
		fmt.Printf("< %s\n", s)
	})

	count := 0
	requestedPath := ""
	loader := func(p string) ([]byte, error) {
		count++
		requestedPath = p
		return nil, nil
	}

	_, err = parser.ParseFile(nil, filename, string(content), 0, parser.WithSourceMapLoader(loader))
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}

	runtime := goja.New()

	// var wg sync.WaitGroup
	// wg.Add(1)
	// defer wg.Done()

	// var debugger *goja.Debugger
	if inspect {
		dbgCh = runtime.EnableDebugMode()
	}

	// go func() {
	// 	var prev *Break
	// 	for b := d.Wait(); b != nil; b, prev = d.Wait(), b {
	// 		if b != prev {
	// 			fmt.Printf("Break at %s:%d", b.Filename(), b.Line())
	// 		}
	// 		fmt.Println("> ")
	// 		cmd := parseCmdFromStdin()
	// 		switch cmd.Name {
	// 		case "cont", "c":
	// 			d.Continue()
	// 		case "next", "n":
	// 			d.Next()
	// 		case "list", "l":
	// 			fmt.Println(b.Source())
	// 		}
	// 	}
	// }()
	// fmt.Println(debugger.Help().Value)
	// fmt.Println(debugger.List().Value)

	registry := new(require.Registry)
	registry.Enable(runtime)
	registry.RegisterNativeModule("console", console.RequireWithPrinter(printer))
	console.Enable(runtime)

	wg.Add(2)

	go runtime.RunScript(filename, string(content))
	// if err != nil {
	// 	log.Fatal(err)
	// 	os.Exit(4)
	// }

	go func() {
		reader := bufio.NewReader(os.Stdin)
		dbg = <-dbgCh // wait for debugger

		for {
			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)
			executor(text)
		}
	}()

	wg.Wait()
}
