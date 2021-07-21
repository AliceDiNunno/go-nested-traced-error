package go_nested_traced_error

import (
	"fmt"
	"hash/crc32"
	"os"
	"runtime"
	"strings"
)

var (
	knownPackageRepositories = []string{
		"bitbucket.org/",
		"code.google.com/",
		"github.com/",
		"launchpad.net/",
	}
)

type Frame struct {
	Filename string `json:"filename"`
	Method   string `json:"method"`
	Line     int    `json:"line"`
}

type Stack []Frame

func NewErrorTrace(err error) *ErrorTrace {
	return &ErrorTrace{
		Err:   err,
		Stack: getStack(),
	}
}

type ErrorTrace struct {
	Err   error
	Stack Stack
	Child *ErrorTrace
}

func AppendNewError(trace *ErrorTrace, err error) *ErrorTrace {
	if trace == nil || err == nil {
		return nil
	}

	newErr := NewErrorTrace(err)
	newErr.Child = trace

	return newErr
}

func fingerprint(str string) string {
	hash := crc32.NewIEEE()

	fmt.Fprintf(hash, str)

	println("calculating fingerprint of")
	println(str)

	return fmt.Sprintf("%x", hash.Sum32())
}

func (s Stack) Fingerprint() string {
	println("stack fingerprint")
	hash := ""
	for _, frame := range s {
		hash = fmt.Sprintf("%s%s%s%d", hash, frame.Filename, frame.Method, frame.Line)
	}
	return fingerprint(hash)
}

func (et *ErrorTrace) Fingerprint() string {
	println("global fingerprint")
	hash := ""
	currentTrace := et
	for {
		println(et.Err.Error())
		hash = fmt.Sprintf("%s%s", hash, currentTrace.Stack.Fingerprint())
		currentTrace = currentTrace.Child
		if currentTrace.Child == nil {
			break
		}
	}
	return fingerprint(hash)
}

func getStack() Stack {
	stack := make(Stack, 0)

	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		functionName := functionName(pc)

		if !strings.HasPrefix(functionName, "go-nested-traced-error") {
			stack = append(stack, Frame{file, functionName, line})
		}
	}

	stack = stack[:len(stack)-1]

	return stack
}

func functionName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "<unknown>"
	}
	name := fn.Name()
	end := strings.LastIndex(name, string(os.PathSeparator))
	return name[end+1:]
}
