package common

// A WindowChecker is used to check if a time is within any of a set of windows.
type WindowChecker interface {
	// SetWindows sets all windows that need to be checked. The windows cannot
	// overlap.
	SetWindows(window [][2]float64)
	// InWindow returns true if the given time is in any of the windows in the
	// window checker.
	InWindow(t int64) bool
	// Seal seals the window checker. This needs to be called after all windows
	// have been added to the window checker and before any calls to InWindow.
	Seal() error
}
