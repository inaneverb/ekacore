package main

import (
	"fmt"

	"github.com/qioalice/gext/sys"
)

func printStacktrace(s sys.StackTrace) {
	for i, frame := range s {
		fmt.Printf("%4d: PC: %d, Entry: %d\n", i, frame.PC, frame.Entry)
		fmt.Printf("      Func: %s\n", frame.Function)
		fmt.Printf("      File: %s\n", frame.File)
		fmt.Println()
	}
}

//go:noinline
func init() {
	fmt.Println()
	fmt.Println("---------- INIT ----------")
	printStacktrace(sys.GetStackTrace(-3, -1).ExcludeInternal())
}

//go:noinline
func foo() {
	fmt.Println("---------- FOO ----------")
	printStacktrace(sys.GetStackTrace(-3, -1).ExcludeInternal())
}

//go:noinline
func bar() {
	defer foo()
	fmt.Println("---------- BAR ----------")
	printStacktrace(sys.GetStackTrace(-3, -1).ExcludeInternal())
}

//go:noinline
func recv() {
	defer func() {
		recover()
		fmt.Println("---------- RECV ----------")
		printStacktrace(sys.GetStackTrace(-3, -1).ExcludeInternal())
	}()
	panic(1)
}

func main() {
	fmt.Println("---------- MAIN ----------")
	printStacktrace(sys.GetStackTrace(-3, -1).ExcludeInternal())
	bar()
	recv()
}
