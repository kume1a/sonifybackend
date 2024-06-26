package shared

import (
	"sync"
)

type result[T interface{}] struct {
	data T
	err  error
}

func streamInputs[INPUT interface{}](done <-chan struct{}, inputs []INPUT) <-chan INPUT {
	inputCh := make(chan INPUT)
	go func() {
		defer close(inputCh)
		for _, input := range inputs {
			select {
			case inputCh <- input:
			case <-done:
				return
			}
		}
	}()
	return inputCh
}

func ExecuteParallel[DATA interface{}, INPUT interface{}](
	parallelism int,
	inputs []INPUT,
	call func(input INPUT) (DATA, error),
) ([]DATA, error) {
	if len(inputs) == 0 {
		return []DATA{}, nil
	}

	done := make(chan struct{})
	defer close(done)

	inputCh := streamInputs(done, inputs)

	var wg sync.WaitGroup
	wg.Add(parallelism)

	resultCh := make(chan result[DATA])

	for i := 0; i < parallelism; i++ {
		// spawn N worker goroutines, each is consuming a shared input channel.
		go func() {
			for input := range inputCh {
				data, err := call(input)
				resultCh <- result[DATA]{data, err}
			}
			wg.Done()
		}()
	}

	// Wait all worker goroutines to finish. Happens if there's no error (no early return)
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	results := []DATA{}
	for result := range resultCh {
		if result.err != nil {
			// return early. done channel is closed, thus input channel is also closed.
			// all worker goroutines stop working (because input channel is closed)
			return nil, result.err
		}
		results = append(results, result.data)
	}

	return results, nil
}
