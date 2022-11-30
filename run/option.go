package run

// RunnerOption configures a Runner.
type RunnerOption[Input, Option, Solution any] func(
	Runner[Input, Option, Solution],
)

// Decode sets the decoder of a runner using f.
func Decode[Input, Option, Solution any](d InputDecoder[Input]) func(
	Runner[Input, Option, Solution],
) {
	return func(r Runner[Input, Option, Solution]) { r.SetInputDecoder(d) }
}
