package run

// IOData describes the data that is used in the IOProducer. The input is the
// source of the input data. The option is the source of the option data. The
// writer is the destination of the output data.
type IOData interface {
	Input() any
	Option() any
	Writer() any
}

// NewIOData creates a new IOData.
func NewIOData(input any, option any, writer any) IOData {
	return ioData{
		input:  input,
		option: option,
		writer: writer,
	}
}

type ioData struct {
	input  any
	option any
	writer any
}

func (d ioData) Input() any {
	return d.input
}

func (d ioData) Option() any {
	return d.option
}

func (d ioData) Writer() any {
	return d.writer
}
