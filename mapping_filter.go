package unicycle

// Like Mapping and Filter at the same time
func MappingFilter[INPUT_TYPE any, OUTPUT_TYPE any](input []INPUT_TYPE, mutatingFilter func(INPUT_TYPE) (OUTPUT_TYPE, bool)) []OUTPUT_TYPE {
	keep := make([]OUTPUT_TYPE, 0, len(input))
	for _, value := range input {
		mutated, ok := mutatingFilter(value)
		if ok {
			keep = append(keep, mutated)
		}
	}
	return Trim(keep)
}

// like MappingFilter(), but accepts a mutator that can return an error, and aborts on the first non-nil error returned by a mutator
func MappingFilterWithError[INPUT_TYPE any, OUTPUT_TYPE any](input []INPUT_TYPE, mutatingFilter func(INPUT_TYPE) (OUTPUT_TYPE, bool, error)) ([]OUTPUT_TYPE, error) {
	keep := make([]OUTPUT_TYPE, 0, len(input))
	for _, value := range input {
		mutated, ok, err := mutatingFilter(value)
		if err != nil {
			return []OUTPUT_TYPE{}, err
		}
		if ok {
			keep = append(keep, mutated)
		}
	}
	return Trim(keep), nil
}

type mappingFilterResult[OUTPUT_TYPE any] struct {
	mutated OUTPUT_TYPE
	ok      bool
}

// like MappingFilter(), but all mutating/filter functions run in parallel in their own goroutines
func MappingFilterMultithread[INPUT_TYPE any, OUTPUT_TYPE any](input []INPUT_TYPE, mutatingFilter func(INPUT_TYPE) (OUTPUT_TYPE, bool)) []OUTPUT_TYPE {
	finished := MappingMultithread(input, func(value INPUT_TYPE) mappingFilterResult[OUTPUT_TYPE] {
		mutated, ok := mutatingFilter(value)
		return mappingFilterResult[OUTPUT_TYPE]{
			mutated: mutated,
			ok:      ok,
		}
	})
	return MappingFilter(finished, func(res mappingFilterResult[OUTPUT_TYPE]) (OUTPUT_TYPE, bool) {
		return res.mutated, res.ok
	})
}

// like MappingFilterWithError(), but all mutating/filter functions run in parallel in their own goroutines
func MappingFilterMultithreadWithError[INPUT_TYPE any, OUTPUT_TYPE any](input []INPUT_TYPE, mutatingFilter func(INPUT_TYPE) (OUTPUT_TYPE, bool, error)) ([]OUTPUT_TYPE, error) {
	finished, err := MappingMultithreadWithError(input, func(value INPUT_TYPE) (mappingFilterResult[OUTPUT_TYPE], error) {
		mutated, ok, err := mutatingFilter(value)
		return mappingFilterResult[OUTPUT_TYPE]{
			mutated: mutated,
			ok:      ok,
		}, err
	})
	if err != nil {
		return []OUTPUT_TYPE{}, err
	}
	return MappingFilter(finished, func(res mappingFilterResult[OUTPUT_TYPE]) (OUTPUT_TYPE, bool) {
		return res.mutated, res.ok
	}), nil
}
