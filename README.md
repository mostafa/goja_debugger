# Goja Debugger

An example Goja application to demonstrate the capabilities of the NEW debugger

## References

- [goja debugger implementation](https://github.com/mostafa/goja/tree/debugger)
- [goja issue for tracking changes and discussion](https://github.com/dop251/goja/issues/294)
- [goja PR](https://github.com/dop251/goja/pull/297)

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
cd scripts
npm i
node gen-srcmap.js test.js > test.js.map
node gen-srcmap.js test_prime.js > test_prime.js.map
```

## How to run

There's just one subcommand in the debugger application, `inspect`. If you use it, it'll tell the compiler to emit `debugger` instruction in your program and also let's the JS VM (Goja VM) to run your script in debug mode. Otherwise it'll just run the script normally. Also, you can pass `line` after `inspect` to have line numbers printed on the `debug[1]>` prompt, instead of program counter (default).

```bash
$ ./goja_debugger inspect scripts/test.js
Welcome to Goja debugger
Type 'help' or 'h' for list of commands.
Loaded sourcemap from: scripts/test.js.map
Break on start in scripts/test.js:1
debug[0]> r
Break on debugger statement in scripts/test.js:11
debug[3]> l
  7         let pest = 2
  8         return val + 2
  9     }
  10
> 11    let i = 1;
  12    console.log(i)
  13    // eval('console.log(i + 2)')
  14
  15    i = 2

debug[3]>
```

## Supported commands

When the debugger starts, it'll pause, so that you instruct it to run a command or set a breakpoint. Most of the following command are supported:

```terminal
setBreakpoint, sb        Set a breakpoint on a given file and line
clearBreakpoint, cb      Clear a breakpoint on a given file and line
breakpoints, b           List all known breakpoints
run, r                   Run program until a breakpoint/debugger statement if program is not started
                         (ProgramStartActivation is disabled, so run doesn't work for now)
next, n                  Continue to next line in current file
cont, c                  Resume execution until next debugger line
step, s                  Step into, potentially entering a function
out, o                   Step out, leaving the current function (not implemented yet)
exec, e                  Evaluate the expression and print the value
list, l                  Print the source around the current line where execution is currently paused
print, p                 Print the provided variable's value
backtrace, bt            Print the current backtrace
help, h                  Print this very help message
quit, q                  Exit debugger and quit (Ctrl+C)
```

## Roadmap

- [x] Expose [a debugger API](https://github.com/dop251/goja/issues/294#issuecomment-869012300) (Thanks to [@nwidger](https://github.com/nwidger) :pray:)
- [x] Implement debugger event-loop using channels for the new API
- [ ] Implementation of step-in (implemented) and step-out commands
- [x] Possible sourcemap generation on the fly (not sure if it's possible in Go) (Thanks to [@nwidger](https://github.com/nwidger) :pray:)
- [ ] Fix Goja tests and try to see if changes are needed for TC39 tests
- [x] Revert changes to interfaces in `compiler.go` and others and use a flag for `debugMode`
- [ ] [DAP](https://microsoft.github.io/debug-adapter-protocol/) integration: [`dap`](https://github.com/mostafa/goja_debugger/tree/dap) branch (boilerplate for now)
- [ ] Integration with [k6](https://github.com/k6io/k6)

## License

AGPL v3 or later
