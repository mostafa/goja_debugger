package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dop251/goja"
)

const (
	GreenColor = "\u001b[32m"
	GrayColor  = "\u001b[38;5;245m"
	ResetColor = "\u001b[0m"
)

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

func executor(in string) bool {
	cmd, err := parseCmd(in)
	if err != nil {
		fmt.Println(err)
		return true
	}

	var result goja.Result
	switch cmd.Name {
	case "setBreakpoint", "sb":
		if len(cmd.Args) < 2 {
			fmt.Println("sb filename linenumber")
			return true
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
			return true
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
		result = dbg.Next()
	case "cont", "continue", "c":
		return false
	case "step", "s":
		result = dbg.StepIn()
	case "out", "o":
		result = dbg.StepOut()
	case "exec", "e":
		result = dbg.Exec(strings.Join(cmd.Args, ";"))
	case "print", "p":
		result = dbg.Print(strings.Join(cmd.Args, ""))
	case "list", "l":
		result = dbg.List()
		if result.Err == nil {
			lines := result.Value.([]string)

			currentLine := dbg.Line()
			lineIndex := currentLine - 1
			var builder strings.Builder
			for idx, lineContents := range lines {
				if inRange(lineIndex, idx-4, idx+4) {
					lineNumber := idx + 1
					totalPadding := 6
					digitCount := countDigits(lineNumber)
					if digitCount >= totalPadding {
						totalPadding = digitCount + 1
					}
					if currentLine == lineNumber {
						padding := strings.Repeat(" ", totalPadding-digitCount)
						builder.Write([]byte(fmt.Sprintf("%s>%s %d%s%s\n", GreenColor, ResetColor, currentLine, padding, lines[lineIndex])))
					} else {
						padding := strings.Repeat(" ", totalPadding-digitCount)
						builder.Write([]byte(fmt.Sprintf("%s  %d%s%s%s\n", GrayColor, lineNumber, padding, lineContents, ResetColor)))
					}
				}
			}
			fmt.Println(builder.String())
		}
		return true
	case "help", "h":
		fmt.Println(help)
		return true
	case "quit", "q":
		os.Exit(0)
	default:
		fmt.Printf("Unknown command, `%s`. You can use `h` to print available commands\n", in)
	}

	prefix := "<" // this should probably be done differently
	if result.Value != nil {
		fmt.Printf("%s%s\n", prefix, result.Value)
	}
	if result.Err != nil {
		fmt.Printf("%sError: %s\n", prefix, result.Err)
	}
	return true
}

func inRange(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}

func countDigits(number int) int {
	if number < 10 {
		return 1
	} else {
		return 1 + countDigits(number/10)
	}
}

var help = `
	setBreakpoint, sb        Set a breakpoint on a given file and line
	clearBreakpoint, cb      Clear a breakpoint on a given file and line
	breakpoints              List all known breakpoints
	next, n                  Continue to next line in current file
	cont, c                  Resume execution until next debugger line
	step, s                  Step into, potentially entering a function (not implemented yet)
	out, o                   Step out, leaving the current function (not implemented yet)
	exec, e                  Evaluate the expression and print the value
	list, l                  Print the source around the current line where execution is currently paused
	print, p                 Print the provided variable's value
	help, h                  Print this very help message
	quit, q                  Exit debugger and quit (Ctrl+C)
`[1:] // this removes the first new line
