package multithread

import (
	"github.com/nuvi/unicycle/slices"
)

// accepts a single input channel and returns a given number of unbuffered child channels that all pull from it
// the child channels close when the parent channel does
func SplitChannel[T any](input chan T, splitCount int) []chan T {
	if splitCount < 1 {
		panic("splitChannel argument splitCount must be > 0")
	}
	output := []chan T{}
	for i := 0; i < splitCount; i++ {
		outChan := make(chan T)
		go func() {
			for value := range input {
				outChan <- value
			}
			close(outChan)
		}()
		output = append(output, outChan)
	}
	return output
}

// accepts any number of channels of the same type and returns a single channel that pulls from all of them at once
// the returned channel closes once all the source channels do
func MergeChannels[T any](input []chan T, capacity int) chan T {
	output := make(chan T, capacity)
	go func() {
		AwaitConcurrent(
			slices.Mapping(
				input,
				func(inputChan chan T) func() {
					return func() {
						for value := range inputChan {
							output <- value
						}
					}
				},
			)...,
		)
		close(output)
	}()
	return output
}