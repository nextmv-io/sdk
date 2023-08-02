package run

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
)

// IOData describes the data that is used in the IOProducer. The input is the
// source of the input data. The option is the source of the option data. The
// writer is the destination of the output data.
type IOData interface {
	Input() any
	Option() any
	Writer() any
}

// NewIOData creates a new IOData.
func NewIOData(input any, option any, writer any) (data IOData, err error) {
	reader, ok := input.(io.Reader)
	if !ok {
		return ioData{
			input:  input,
			option: option,
			writer: writer,
		}, nil
	}

	if closer, ok := reader.(io.Closer); ok {
		defer func() {
			tempErr := closer.Close()
			// the first error is the most important
			if err == nil {
				err = tempErr
			}
		}()
	}

	// Convert to buffered reader and read magic bytes
	bufferedReader := bufio.NewReader(reader)
	testBytes, err := bufferedReader.Peek(2)

	// Test for gzip magic bytes and use corresponding reader, if given
	if err == nil && testBytes[0] == 31 && testBytes[1] == 139 {
		var gzipReader *gzip.Reader
		if gzipReader, err = gzip.NewReader(bufferedReader); err != nil {
			return ioData{}, err
		}
		reader = gzipReader
	} else {
		// Default case: assume text input
		reader = bufferedReader
	}

	// copy input to buffer
	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return ioData{}, err
	}

	return ioData{
		input:  input,
		option: option,
		writer: writer,
		buf:    buf,
	}, nil
}

type ioData struct {
	input  any
	option any
	writer any
	buf    *bytes.Buffer
}

func (d ioData) Input() (input any) {
	// buffer was filled so use that instead of the original reader
	if len(d.buf.Bytes()) > 0 {
		return bytes.NewReader(d.buf.Bytes())
	}
	return d.input
}

func (d ioData) Option() any {
	return d.option
}

func (d ioData) Writer() any {
	return d.writer
}
