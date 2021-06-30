# Goja Debugger

An example Goja application to demonstrate the capabilities of the NEW debugger

## References

- [goja debugger implementation](https://github.com/mostafa/goja/tree/debugger)
- [goja issue for tracking changes and discussion](https://github.com/dop251/goja/issues/294)

## How to build

This repositories contains an example application in Go for testing the debugger functionality I'm developing for [Goja](https://github.com/dop251/goja) and a bunch of scripts in JavaScript, one (`gen-srcmap.js`) for generating sourcemaps for your JavaScript files and two other (`test_*.js`) for testing the debugger app.

To run the debugger application, just run `go build` in the cloned project on your machine, assuming you have Go.

```bash
git clone https://github.com/mostafa/goja_debugger
cd goja_debugger
go build
```

There are two scripts along with their sourcemaps for testing the debugger functionality. If you want to debug your own scripts, you might need to use Node.js to generate your own sourcemap for your JavaScript, which helps the debugger application easily pinpoint the lines in your script.

```bash
npm i
node gen-srcmap.js test.js > test.js.map
```

## How to run

There's just one subcommand in the debugger application, `inspect`. If you use it, it'll tell the compiler to emit `debugger` instruction in your program and also let's the JS VM (Goja VM) to run your script in debug mode. Otherwise it'll just run the script normally. Also, you can pass `line` after `inspect` to have line numbers printed on the `debug[1]>` prompt, instead of program counter (default).

```bash
./goja_debugger inspect test.js
Welcome to Goja debugger
Type 'help' or 'h' for list of commands.
Loaded sourcemap from: test.js.map
debug[0]> n
Break on start in test.js:1:1
> 1     debugger
  2
  3     function test(val) {
  4         debugger
  5         val += 10

debug[0]>
```

## Supported commands

When the debugger starts, it'll pause, so that you instruct it to run a command or set a breakpoint. Most of the following command are supported:

```terminal
setBreakpoint, sb        Set a breakpoint on a given file and line
clearBreakpoint, cb      Clear a breakpoint on a given file and line
breakpoints, b           List all known breakpoints
run, r                   Run program until a breakpoint/debugger statement if program is not started
next, n                  Continue to next line in current file
cont, c                  Resume execution until next debugger line
step, s                  Step into, potentially entering a function
out, o                   Step out, leaving the current function (not implemented yet)
exec, e                  Evaluate the expression and print the value
list, l                  Print the source around the current line where execution is currently paused
print, p                 Print the provided variable's value
help, h                  Print this very help message
quit, q                  Exit debugger and quit (Ctrl+C)
```

## Roadmap

- [x] Expose [a debugger API](https://github.com/dop251/goja/issues/294#issuecomment-869012300) (Thanks to [@nwidger](https://github.com/nwidger) :pray:)
- [x] Implement debugger event-loop using channels for the new API
- [ ] Implementation of step-in (implemented) and step-out commands
- [ ] Possible sourcemap generation on the fly (not sure if it's possible in Go)
- [ ] Fix Goja tests and try to see if changes are needed for TC39 tests
- [x] Revert changes to interfaces in `compiler.go` and others and use a flag for `debugMode`
- [ ] [DAP](https://microsoft.github.io/debug-adapter-protocol/) integration
- [ ] Integration with [k6](https://github.com/k6io/k6)

## License

AGPL v3 or later
