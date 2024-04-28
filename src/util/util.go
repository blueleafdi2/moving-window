package util

import (
	"errors"
	"fmt"
	"runtime/debug"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
func Quietly(call func(), printStack ...bool) (err error) {
	needPrintStack := false
	if len(printStack) > 0 {
		needPrintStack = printStack[0]
	}

	defer func() {
		if e := recover(); e != nil {
			if needPrintStack {
				debug.PrintStack()
			}
			switch e.(type) {
			case error:
				err = e.(error)
			default:
				err = errors.New(fmt.Sprintf("call with error: %v", e))
			}
		}
	}()
	call()
	return err
}
