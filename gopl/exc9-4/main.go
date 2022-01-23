package pipe

// type pipe struct {
// 	in  <-chan interface{}
// 	out chan interface{}
// }

// asyncronous pipe
func createPipe(in <-chan interface{}) chan interface{} {
	out := make(chan interface{})
	go func(in <-chan interface{}, out chan interface{}) {
		v := <-in
		out <- v
	}(in, out)

	return out
}

func pipeLine(stages int, in <-chan interface{}) chan interface{} {
	out := createPipe(in)
	for i := 1; i < stages; i++ {
		next := createPipe(out)
		out = next
	}
	return out
}

// TODO: run benchmarks

// Construct a pipeline that connects an arbitrary number of goroutines
// with channels. What is the maximum number of pipeline stages you can create
// without running out of memory? How long does a value take to transit the entire
// pipeline?
