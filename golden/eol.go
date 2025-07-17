package golden

import (
	"bytes"
)

// convertNewLineStyle converts the new line style of the content based on the
// specified NewLineStyle.
func convertNewLineStyle(content []byte, style NewLineStyle) []byte {
	switch style {
	case NewLineStyleUntouched:
		return content // Keep the original line endings
	case NewLineStyleLF:
		return bytes.ReplaceAll(content, []byte{'\r', '\n'}, []byte{'\n'})
	case NewLineStyleCRLF:
		return bytes.ReplaceAll(content, []byte{'\n'}, []byte{'\r', '\n'})
	default:
		return content // Default case, keep as is
	}
}
