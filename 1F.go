//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type Command struct {
	textualRepresentation string
	isDefault             bool
}

type Evaluator struct {
	defs  map[string][]Command
	stack []int
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{stack: make([]int, 0), defs: make(map[string][]Command)}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
// The string processing
func (e *Evaluator) Process(row string) ([]int, error) {
	// process into commands and divert handling
	var commandList []Command

	// dictionary of always default commands
	alwaysDefaultCommands := []string{":", ";"}

	for _, word := range strings.Fields(strings.ToLower(row)) {
		isDefault := false
		// check for number
		if _, err := strconv.Atoi(word); err == nil {
			isDefault = true
		}
		// check for default commands
		for _, defaultCommand := range alwaysDefaultCommands {
			if word == defaultCommand {
				isDefault = true
			}
		}
		// create command
		commandList = append(commandList, Command{textualRepresentation: word, isDefault: isDefault})
	}

	// evaluate as commands
	return e.ProcessCommandList(commandList)
}

// The command processing
func (e *Evaluator) ProcessCommandList(words []Command) ([]int, error) {
	// make error placeholder
	var err error

	// dictionary of always default commands
	alwaysDefaultCommands := []string{":", ";"}
	// dictionary of all default commands
	defaultCommands := []string{"dup", "over", "drop", "swap", "+", "-", "*", "/", ":", ";"}

	// evaluate commands
	for idx := 0; idx < len(words) && err == nil; idx++ {
		// handling default command behavior
		if words[idx].isDefault {
			switch words[idx].textualRepresentation {
			case "dup":
				if len(e.stack) == 0 {
					err = errors.New("stack underflow")
				} else {
					e.stack = append(e.stack, e.stack[len(e.stack)-1])
				}
			case "over":
				if len(e.stack) < 2 {
					err = errors.New("stack underflow")
				} else {
					e.stack = append(e.stack, e.stack[len(e.stack)-2])
				}
			case "drop":
				if len(e.stack) == 0 {
					err = errors.New("stack underflow")
				} else {
					e.stack = e.stack[:len(e.stack)-1]
				}
			case "swap":
				if len(e.stack) < 2 {
					err = errors.New("stack underflow")
				} else {
					e.stack[len(e.stack)-2], e.stack[len(e.stack)-1] = e.stack[len(e.stack)-1], e.stack[len(e.stack)-2]
				}
			case "+":
				if len(e.stack) < 2 {
					err = errors.New("stack underflow")
				} else {
					e.stack[len(e.stack)-2] += e.stack[len(e.stack)-1]
					e.stack = e.stack[:len(e.stack)-1]
				}
			case "-":
				if len(e.stack) < 2 {
					err = errors.New("stack underflow")
				} else {
					e.stack[len(e.stack)-2] -= e.stack[len(e.stack)-1]
					e.stack = e.stack[:len(e.stack)-1]
				}
			case "*":
				if len(e.stack) < 2 {
					err = errors.New("stack underflow")
				} else {
					e.stack[len(e.stack)-2] *= e.stack[len(e.stack)-1]
					e.stack = e.stack[:len(e.stack)-1]
				}
			case "/":
				if len(e.stack) < 2 {
					err = errors.New("stack underflow")
				} else {
					// handle division by zero
					if e.stack[len(e.stack)-1] == 0 {
						err = errors.New("division by zero")
					} else {
						e.stack[len(e.stack)-2] /= e.stack[len(e.stack)-1]
						e.stack = e.stack[:len(e.stack)-1]
					}
				}
			case ";":
				err = errors.New("syntax error, unexpected ';'")
			case ":":
				// new definition
				newCommand := words[idx+1].textualRepresentation

				// check for number redefinition
				_, err = strconv.Atoi(newCommand)
				if err == nil {
					err = errors.New("redefinition of number")
					break
				}
				err = nil

				// check for redefinition of always default commands
				isDefault := false
				for _, defaultCommand := range alwaysDefaultCommands {
					if newCommand == defaultCommand {
						isDefault = true
					}
				}
				if isDefault {
					err = errors.New("redefinition of always default command")
				}

				// handle errors
				if err != nil {
					break
				}

				// update and create helper variables
				idx += 2
				def := make([]Command, 0)

				// process definition
				for words[idx].textualRepresentation != ";" && idx < len(words) {
					// if not default command put its definition instead
					if !words[idx].isDefault {

						// handle case of missing definition
						if _, ok := e.defs[words[idx].textualRepresentation]; !ok {
							// handle word with default behavior
							isDefault := false
							for _, defaultCommand := range defaultCommands {
								if words[idx].textualRepresentation == defaultCommand {
									isDefault = true
								}
							}

							if isDefault {
								def = append(def, Command{textualRepresentation: words[idx].textualRepresentation, isDefault: true})
								idx++
								continue
							}

							// handle word with no default behavior
							err = errors.New("undefined word: " + words[idx].textualRepresentation)
							break
						}

						// if definition exists, append it
						def = append(def, e.defs[words[idx].textualRepresentation]...)
					} else {
						def = append(def, words[idx])
					}
					idx++
				}
				if idx >= len(words) {
					err = errors.New("syntax error, definition not terminated")
				}

				// update dictionary
				e.defs[newCommand] = def
			default:
				// handle numbers
				n := 0
				n, err = strconv.Atoi(words[idx].textualRepresentation)
				if err != nil {
					err = errors.New("undefined word: " + words[idx].textualRepresentation)
					break
				}
				e.stack = append(e.stack, n)
			}

			// handle user defined behavior
		} else {
			// check for missing definition
			if _, ok := e.defs[words[idx].textualRepresentation]; !ok {
				// handle word with default behavior
				isDefault := false
				for _, defaultCommand := range defaultCommands {
					if words[idx].textualRepresentation == defaultCommand {
						isDefault = true
					}
				}

				if isDefault {
					e.stack, err = e.ProcessCommandList([]Command{{textualRepresentation: words[idx].textualRepresentation, isDefault: true}})
					continue
				}

				// handle word with no default behavior
				err = errors.New("undefined word: " + words[idx].textualRepresentation)
			} else {
				// evaluate definition
				e.stack, err = e.ProcessCommandList(e.defs[words[idx].textualRepresentation])
			}
		}
	}

	return e.stack, err
}
