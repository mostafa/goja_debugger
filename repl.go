package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
)

// const (
// 	SetBreakpoint   = "sb"
// 	ClearBreakpoint = "cb"
// 	Breakpoints     = "breakpoints"
// 	Next            = "n"
// 	Continue        = "c"
// 	StepIn          = "s"
// 	StepOut         = "o"
// 	Exec            = "e"
// 	Print           = "p"
// 	List            = "l"
// 	Help            = "h"
// 	Quit            = "q"
// 	Empty           = ""
// 	NewLine         = "\n"
// )

type Command struct {
	Name string
	Args []string
}

func parseCmd(in string) (*Command, error) {
	data := strings.Split(in, " ")
	if len(data) == 0 {
		return nil, errors.New("unknown command")
	}
	name := data[0]
	var args []string
	if len(data) > 1 {
		args = append(args, data[1:]...)
	}
	return &Command{Name: name, Args: args}, nil
}

func getResult(prefix string) {
	wg.Add(1)
	result := dbg.GetResult()
	if result.Value != nil {
		fmt.Printf("%s%s\n", prefix, result.Value)
	}
	if result.Err != nil {
		fmt.Printf("%sError: %s\n", prefix, result.Err)
	}
	wg.Done()
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "setBreakpoint, sb", Description: "Set a breakpoint on a given file and line"},
		{Text: "clearBreakpoint, cb", Description: "Clear a breakpoint on a given file and line"},
		{Text: "breakpoints", Description: "List all known breakpoints"},
		{Text: "next, n", Description: "Continue to next line in current file"},
		{Text: "cont, c", Description: "Resume execution until next debugger line"},
		{Text: "step, s", Description: "Step into, potentially entering a function (not implemented yet)"},
		{Text: "out, o", Description: "Step out, leaving the current function (not implemented yet)"},
		{Text: "exec, e", Description: "Evaluate the expression and print the value"},
		{Text: "list, l", Description: "Print the source around the current line where execution is currently paused"},
		{Text: "print, p", Description: "Print the provided variable's value"},
		{Text: "help, h", Description: "Print this very help message"},
		{Text: "quit, q", Description: "Exit debugger and quit (Ctrl+C)"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(in string) {
	cmd, err := parseCmd(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch cmd.Name {
	case "setBreakpoint", "sb":
		if len(cmd.Args) < 2 {
			fmt.Println("sb filename linenumber")
			return
		}
		if line, err := strconv.Atoi(cmd.Args[1]); err != nil {
			fmt.Printf("Cannot convert %s to line number\n", cmd.Args[1])
		} else {
			err := dbg.SetBreakpoint(cmd.Args[0], line)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	case "clearBreakpoint", "cb":
		if len(cmd.Args) < 2 {
			fmt.Println("cb filename linenumber")
			return
		}
		if line, err := strconv.Atoi(cmd.Args[1]); err != nil {
			fmt.Printf("Cannot convert %s to line number\n", cmd.Args[1])
		} else {
			err := dbg.ClearBreakpoint(cmd.Args[0], line)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	case "breakpoints":
		breakpoints, err := dbg.Breakpoints()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			for _, b := range breakpoints {
				fmt.Printf("Breakpoint on %s:%d\n", b.Filename, b.Line)
			}
		}
	case "next", "n":
		dbg.Next()
		go getResult("< ")
	case "cont", "continue", "c":
		dbg.Continue()
		go getResult("< ")
	case "step", "s":
		dbg.StepIn()
		go getResult("< ")
	case "out", "o":
		dbg.StepOut()
		go getResult("< ")
	case "exec", "e":
		dbg.Exec(strings.Join(cmd.Args, ";"))
		go getResult("< ")
	case "print", "p":
		dbg.Print(strings.Join(cmd.Args, ""))
		go getResult("< ")
	case "list", "l":
		dbg.List()
		go getResult("")
	case "help", "h":
		// result := dbg.Help()
		// fmt.Print(result.Value)
	case "quit", "q":
		os.Exit(0)
	default:
		// dbg.Quit(0)
	}

	wg.Done()
}

type Cmd struct {
	Name string
	Args []string
}

// func REPL(dbg *goja.Debugger, intro bool) {
// 	// Refactor this piece of sh!t
// 	debuggerCommands := map[string]string{
// 		"setBreakpoint":   SetBreakpoint,
// 		SetBreakpoint:     SetBreakpoint,
// 		"clearBreakpoint": ClearBreakpoint,
// 		ClearBreakpoint:   ClearBreakpoint,
// 		"breakpoints":     Breakpoints,
// 		"next":            Next,
// 		Next:              Next,
// 		"cont":            Continue,
// 		Continue:          Continue,
// 		"step":            StepIn,
// 		StepIn:            StepIn,
// 		"out":             StepOut,
// 		StepOut:           StepOut,
// 		"exec":            Exec,
// 		Exec:              Exec,
// 		"print":           Print,
// 		Print:             Print,
// 		"list":            List,
// 		List:              List,
// 		"help":            Help,
// 		Help:              Help,
// 		"quit":            Quit,
// 		Quit:              Quit,
// 		NewLine:           "\n",
// 	}
// 	if intro {
// 		fmt.Println("Welcome to Goja debugger")
// 		fmt.Println("Type 'help' or 'h' for list of commands.")
// 	} else {
// 		if dbg.IsBreakOnStart() {
// 			fmt.Printf("Break on start in %s:%d\n", dbg.Filename(), dbg.Line())
// 		} else {
// 			fmt.Printf("Break in %s:%d\n", dbg.Filename(), dbg.Line())
// 		}
// 		result := dbg.List()
// 		fmt.Println(result.Value)
// 		if result.Err != nil {
// 			fmt.Println(result.Err)
// 		}
// 	}

// 	var commandAndArguments []string

// 	// var prev *goja.Breakpoint
// 	for b := dbg.Wait(); b != nil; b, _ = dbg.Wait(), b {
// 		command, err := parseCmd(dbg.GetPC())
// 		commandAndArguments = strings.Split(command[:len(command)-1], " ")
// 		if command == NewLine && len(dbg.LastDebuggerCmdAndArgs) > 0 {
// 			// If enter is pressed and there's a command already executed,
// 			// run the last debugger command
// 			commandAndArguments = make([]string, len(dbg.LastDebuggerCmdAndArgs))
// 			copy(commandAndArguments, dbg.LastDebuggerCmdAndArgs)
// 		}

// 		if v, ok := debuggerCommands[commandAndArguments[0]]; ok {
// 			if command != NewLine {
// 				// FIXME: Exec command acts as Next on the next run
// 				dbg.LastDebuggerCmdAndArgs = make([]string, len(commandAndArguments))
// 				copy(dbg.LastDebuggerCmdAndArgs, commandAndArguments)
// 			}

// 			switch v {
// 			case SetBreakpoint:
// 				if len(commandAndArguments) < 3 {
// 					fmt.Println("sb filename linenumber")
// 					continue
// 				}
// 				if line, err := strconv.Atoi(commandAndArguments[2]); err != nil {
// 					fmt.Printf("Cannot convert %s to line number\n", commandAndArguments[2])
// 				} else {
// 					err := dbg.SetBreakpoint(commandAndArguments[1], line)
// 					if err != nil {
// 						fmt.Println(err.Error())
// 					}
// 				}
// 			case ClearBreakpoint:
// 				if len(commandAndArguments) < 3 {
// 					fmt.Println("cb filename linenumber")
// 					continue
// 				}
// 				if line, err := strconv.Atoi(commandAndArguments[2]); err != nil {
// 					fmt.Printf("Cannot convert %s to line number\n", commandAndArguments[2])
// 				} else {
// 					err := dbg.ClearBreakpoint(commandAndArguments[1], line)
// 					if err != nil {
// 						fmt.Println(err.Error())
// 					}
// 				}
// 			case Breakpoints:
// 				breakpoints, err := dbg.Breakpoints()
// 				if err != nil {
// 					fmt.Println(err.Error())
// 				} else {
// 					for _, b := range breakpoints {
// 						fmt.Printf("Breakpoint on %s:%d\n", b.Filename, b.Line)
// 					}
// 				}
// 			case Next:
// 				return
// 			case Continue:
// 				return
// 			case StepIn:
// 				fmt.Println(dbg.StepIn())
// 			case StepOut:
// 				fmt.Println(dbg.StepOut())
// 			case Exec:
// 				result := dbg.Exec(strings.Join(commandAndArguments[1:], ";"))
// 				if result.Err != nil {
// 					fmt.Println(result.Err)
// 				}
// 			case Print:
// 				result := dbg.Print(strings.Join(commandAndArguments[1:], ""))
// 				fmt.Printf("< %s\n", result.Value)
// 				if err != nil {
// 					fmt.Printf("< Error: %s\n", result.Err)
// 				}
// 			case List:
// 				result := dbg.List()
// 				fmt.Print(result.Value)
// 				if err != nil {
// 					fmt.Println(result.Err)
// 				}
// 			// case Help:
// 			// result := dbg.Help()
// 			// fmt.Print(result.Value)
// 			case Quit:
// 				dbg.Quit(0)
// 			default:
// 				dbg.Quit(0)
// 			}
// 		} else {
// 			fmt.Println("unknown command")
// 		}
// 	}
// }
