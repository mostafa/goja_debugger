package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

var dbg *goja.Debugger
var wg sync.WaitGroup

func getLivePrefix() (string, bool) {
	if dbg == nil {
		return "debug>", false
	} else {
		return fmt.Sprintf("debug[%d]> ", dbg.GetPC()), true
	}
}

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
		dbg = runtime.EnableDebugMode()
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

	p := prompt.New(
		executor,
		completer,
		prompt.OptionLivePrefix(getLivePrefix),
	)
	go p.Run()

	wg.Wait()
}
