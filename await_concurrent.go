package unicycle

import (
	"errors"
	"log"
)

// AwaitConcurrent simplifies the common task of waiting until tasks on multiple threads have finished
func AwaitConcurrent(funcs ...func()) {
	pending := make(chan struct{})
	finished := 0
	for _, wrapped := range funcs {
		go awaitSafe(pending, wrapped)
	}
	for _ = range pending {
		finished++
		if finished == len(funcs) {
			return
		}
	}
}

func awaitSafe(pending chan struct{}, wrapped func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicking goroutine in AwaitConcurrent recovered:", r)
			pending <- Empty
		}
	}()
	wrapped()
	pending <- Empty
}

// Like AwaitConcurrent, but accepts functions that return errors, and returns the first error if there is one
func AwaitConcurrentWithErrors(funcs ...func() error) error {
	pending := make(chan error)
	finished := 0
	for _, wrapped := range funcs {
		go awaitUnsafe(pending, wrapped)
	}
	for err := range pending {
		if err != nil {
			return err
		}
		finished++
		if finished == len(funcs) {
			return nil
		}
	}
	panic("pending channel closed!") // this should never actually happen but the line is required to compile
}

func awaitUnsafe(pending chan error, wrapped func() error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			pending <- ErrAwaitConcurrentWithErrorsPanic
		}
	}()
	pending <- wrapped()
}

var ErrAwaitConcurrentWithErrorsPanic = errors.New("panicking goroutine in AwaitConcurrentWithErrors recovered")