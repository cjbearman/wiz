package wiz

type Waitable interface {
	Wait()
}

// AsyncOperation allows for multiple operations to be performed concurrently
type AsyncOperation[T any] struct {
	done     chan struct{}
	result   T
	error    error
	executor func() (T, error)
	device   *Device
}

// AsyncOperations represents a collection of async operations of the same type
type AsyncOperations[T any] []*AsyncOperation[T]

// DeviceResultAndError represents the result and error of an operation for a specific device
type DeviceResultAndError[T any] struct {
	Device *Device
	Result T
	Err    error
}

// DeviceResultsAndError represents a collection of device results and errors
type DeviceResultsAndError[T any] []DeviceResultAndError[T]

// NewAsyncOperations creates and starts an async operation using the provided excutor function
func NewAsyncOperation[T any](device *Device, executor func() (T, error)) *AsyncOperation[T] {
	op := &AsyncOperation[T]{
		done:     make(chan struct{}),
		executor: executor,
		device:   device,
	}
	go func() {
		op.result, op.error = op.executor()
		close(op.done)
	}()
	return op
}

// AwaitResult waits for the operation to complete and returns the result and any error encountered
func (op *AsyncOperation[T]) AwaitResult() (T, error) {
	<-op.done
	return op.result, op.error
}

// DeliverResultTo will deliver the result of an async operation to the specified channel
func (op *AsyncOperation[T]) DeliverResultTo(ch chan<- DeviceResultAndError[T]) *AsyncOperation[T] {
	res, err := op.AwaitResult()
	ch <- DeviceResultAndError[T]{Device: op.device, Result: res, Err: err}
	return op
}

// DeliverMixedResultTo will deliver the result of an async operation to a channel expecting any type
// of result, not just those of a specific type
func (op *AsyncOperation[T]) DeliverMixedResultTo(ch chan<- DeviceResultAndError[any]) *AsyncOperation[T] {
	res, err := op.AwaitResult()
	ch <- DeviceResultAndError[any]{Device: op.device, Result: res, Err: err}
	return op
}

// IsDone will return true if an async operation has completed
func IsDone[T any](op *AsyncOperation[T]) bool {
	select {
	case <-op.done:
		return true
	default:
		return false
	}
}

// AwaitAll waits for all async operations in the collection to complete and returns their results and errors
func (ops AsyncOperations[T]) AwaitAll() DeviceResultsAndError[T] {
	results := make(DeviceResultsAndError[T], len(ops))

	// It doesn't matter which order we wait for them, they all have to finish
	for i, op := range ops {
		res, err := op.AwaitResult()
		results[i] = DeviceResultAndError[T]{Device: op.device, Result: res, Err: err}
	}

	return results
}

// DeliverResultsTo will deliver the results of all async operations in the collection to the specified channel
func (ops AsyncOperations[T]) DeliverResultsTo(ch chan<- DeviceResultAndError[T]) AsyncOperations[T] {
	for _, op := range ops {
		op.DeliverResultTo(ch)
	}
	return ops
}

// DeliverMixedResultsTo will deliver the results of multiple operations of mixed types
// to the same channel
// This is useful when you have a collection of async operations of different types and want to handle their results in a unified manner.
// but you should examine the type of the embedded result before using it
func (ops AsyncOperations[T]) DeliverMixedResultsTo(ch chan<- DeviceResultAndError[any]) AsyncOperations[T] {
	for _, op := range ops {
		op.DeliverMixedResultTo(ch)
	}
	return ops
}
