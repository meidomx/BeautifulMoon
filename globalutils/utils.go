package globalutils

import "errors"

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicFalse(result bool) {
	if !result {
		panic(errors.New("result is false."))
	}
}
