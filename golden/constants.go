package golden

const (
	// StableTime is a stable time to be used in golden file tests.
	StableTime string = "2023-01-01T00:00:00Z"
	// StableDuration is a stable duration to be used in golden file tests.
	StableDuration string = "123ms"
	// StableText is a stable text to be used in golden file tests.
	StableText string = "text"
	// StableFloat is a stable float to be used in golden file tests.
	StableFloat float64 = 0.123
	// StableInt is a stable int to be used in golden file tests.
	StableInt int = 123
	// StableBool is a stable bool to be used in golden file tests.
	StableBool bool = true
	// StableVersion is a stable version to be used in golden file tests.
	StableVersion string = "VERSION"
	// binaryName is the name of the binary used for testing. We use "main.exe"
	// here to avoid running into "cannot exec" problems on Windows platforms.
	binaryName = "main.exe"
)
