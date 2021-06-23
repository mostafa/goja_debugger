package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
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
	if inspect {
		runtime.EnableDebugMode()
	}

	registry := new(require.Registry)
	registry.Enable(runtime)
	registry.RegisterNativeModule("console", console.RequireWithPrinter(printer))
	console.Enable(runtime)

	_, err = runtime.RunScript(filename, string(content))
	if err != nil {
		log.Fatal(err)
		os.Exit(4)
	}
}
