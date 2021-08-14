package GoNestedTracedError

import "fmt"

type Error struct {
	Err   error
	Stack Stack
	Child *Error
}

func Trace(err error) *Error {
	return &Error{
		Err:   err,
		Stack: getStack(),
	}
}

func Wrap(err error) *Error {
	if err == nil {
		return nil
	}
	return Trace(err)
}

func (trace *Error) Append(err error) *Error {
	if trace == nil || err == nil {
		return nil
	}

	newErr := Trace(err)
	newErr.Child = trace

	return newErr
}

func (trace *Error) Fingerprint() string {
	hash := ""
	currentTrace := trace
	for {
		hash = fmt.Sprintf("%s%s", hash, currentTrace.Stack.Fingerprint())
		currentTrace = currentTrace.Child
		if currentTrace == nil {
			break
		}
	}
	return fingerprint(hash)
}
