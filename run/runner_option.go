package run

// RunnerOption configures a Runner.
type RunnerOption[RunnerConfig, Input, Option, Solution any] func(
	Runner[RunnerConfig, Input, Option, Solution],
)

// InputDecode sets the input decoder of a runner.
func InputDecode[
	RunnerConfig, Input, Option, Solution any,
](i Decoder[Input]) func(
	Runner[RunnerConfig, Input, Option, Solution],
) {
	return func(r Runner[RunnerConfig, Input, Option, Solution]) {
		r.SetInputDecoder(i)
	}
}

// InputValidate sets the input validator of a runner.
func InputValidate[
	RunnerConfig, Input, Option, Solution any,
](v Validator[Input]) func(
	Runner[RunnerConfig, Input, Option, Solution],
) {
	return func(r Runner[RunnerConfig, Input, Option, Solution]) {
		r.SetInputValidator(v)
	}
}

// OptionDecode sets the options decoder of a runner.
func OptionDecode[
	RunnerConfig, Input, Option, Solution any,
](o Decoder[Option]) func(
	Runner[RunnerConfig, Input, Option, Solution],
) {
	return func(r Runner[RunnerConfig, Input, Option, Solution]) {
		r.SetOptionDecoder(o)
	}
}

// Encode sets the encoder of a runner.
func Encode[
	RunnerConfig, Input, Option, Solution any,
](e Encoder[Solution, Option]) func(
	Runner[RunnerConfig, Input, Option, Solution],
) {
	return func(r Runner[RunnerConfig, Input, Option, Solution]) {
		r.SetEncoder(e)
	}
}

// IOProduce sets the IOProducer of a runner.
func IOProduce[
	RunnerConfig, Input, Option, Solution any,
](i IOProducer[RunnerConfig]) func(
	Runner[RunnerConfig, Input, Option, Solution],
) {
	return func(r Runner[RunnerConfig, Input, Option, Solution]) {
		r.SetIOProducer(i)
	}
}
