package alns

// NewSolveOptions returns a new instance of SolveOptions. The default
// options are:
//   - iterations: -1
//   - solution_buffer_size: 100
func NewSolveOptions() SolveOptions {
	return SolveOptions{
		Iterations:         -1,
		SolutionBufferSize: 100,
	}
}

// SolveOptions contains the options for the solve process.
type SolveOptions struct {
	Iterations         int `json:"iterations" default:"-1"`
	SolutionBufferSize int `json:"solution_buffer_size" default:"100"`
}
