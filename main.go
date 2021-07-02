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
	"github.com/evanw/esbuild/pkg/api"
	"github.com/urfave/cli/v2"
)

var dbg *goja.Debugger
var runtime *goja.Runtime

func main() {
	app := &cli.App{
		Name:  "goja_debugger",
		Usage: "Runs or inspects a JS script with Goja",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Runs a JS script with Goja",
				Action: func(c *cli.Context) error {
					return debug(false, "", c.Args().First())
				},
			},
			{
				Name:    "inspect",
				Aliases: []string{"i"},
				Usage:   "Debugs a JS script with Goja",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "liveinfo",
						Aliases: []string{"l"},
						Value:   "pc",
						Usage:   "Show program counter (pc) or line number (line) in debug prompt",
					},
				},
				Action: func(c *cli.Context) error {
					return debug(true, c.String("liveinfo"), c.Args().First())
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func debug(inspect bool, liveInfo, filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if inspect {
		// Generate sourceamp on-the-fly
		content = generateSourceMap(filename, string(content))

		fmt.Println("Welcome to Goja debugger")
		fmt.Println("Type 'help' or 'h' for list of commands.")
	}

	printer := console.PrinterFunc(func(s string) {
		prefix := ""
		if inspect {
			prefix = "< "
		}
		fmt.Printf("%s%s\n", prefix, s)
	})

	loader := func(requestedPath string) ([]byte, error) {
		if requestedPath != "" && inspect {
			fmt.Printf("%sLoaded sourcemap from: %s%s\n", GrayColor, requestedPath, ResetColor)
		}
		return nil, nil
	}

	_, err = parser.ParseFile(nil, filename, string(content), 0, parser.WithSourceMapLoader(loader))
	if err != nil {
		log.Fatal(err)
		return err
	}

	runtime = goja.New()

	if inspect {
		dbg = runtime.EnableDebugMode()
	}

	registry := new(require.Registry)
	registry.Enable(runtime)
	registry.RegisterNativeModule("console", console.RequireWithPrinter(printer))
	console.Enable(runtime)

	go func() {
		if inspect {
			reader := bufio.NewReader(os.Stdin)

			reason, resume := dbg.WaitToActivate()
			printDebuggingReason(reason)
			for {
				fmt.Printf("debug%s> ", getInfo(liveInfo))
				userInput, _ := reader.ReadString('\n')
				// convert CRLF to LF
				userInput = strings.Replace(userInput, "\n", "", -1)
				if !repl(userInput) {
					resume()
					reason, resume = dbg.WaitToActivate()
					printDebuggingReason(reason)
				}
			}
		}
	}()

	runtime.RunScript(filename, string(content))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func printDebuggingReason(reason goja.ActivationReason) {
	if reason == goja.ProgramStartActivation {
		fmt.Printf("Break on start in %s:%d\n", dbg.Filename(), dbg.Line())
	} else if reason == goja.BreakpointActivation {
		fmt.Printf("Break on breakpoint in %s:%d\ns", dbg.Filename(), dbg.Line())
	} else {
		fmt.Printf("Break on debugger statement in %s:%d\n", dbg.Filename(), dbg.Line())
	}
}

func getInfo(liveInfo string) string {
	if liveInfo == "line" {
		return fmt.Sprintf("[%d]", dbg.Line())
	}
	return fmt.Sprintf("[%d]", dbg.PC())
}

func generateSourceMap(filename string, src string) []byte {
	result := api.Transform(src, api.TransformOptions{
		Sourcemap:  api.SourceMapInline,
		Sourcefile: filename,
	})

	if len(result.Errors) > 0 {
		fmt.Println(result.Errors)
	}

	return result.Code
}
