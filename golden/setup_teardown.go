package golden

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	osPath "path"
	"regexp"
	"strings"
)

// SetupWithFileGen gets the main.go from the template before setting up golden
// file testing. Must be called before test execution.
func SetupWithFileGen(templateName string) {
	flag.Parse()
	if *update {
		GenerateTemplate(templateName, "")
		safeDelete("main.go")
		safeDelete("input.json")
		CopyFile(fmt.Sprintf("%s/main.go", templateName), "main.go")
		CopyFile(fmt.Sprintf("%s/input.json", templateName), "input.json")
		if err := os.RemoveAll(templateName); err != nil {
			panic(err)
		}
	}

	Setup()
}

// SetupTemplateTest sets up golden file testing with a template. Must be
// called before test execution.
func SetupTemplateTest(templateName string) {
	GenerateTemplate(templateName, "")

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(templateName)
	if err != nil {
		panic(err)
	}

	Setup()
	CopyFile(binaryName, fmt.Sprintf("../%s", binaryName))

	err = os.Chdir(currentDir)
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(templateName)
	if err != nil {
		panic(err)
	}
}

// Setup sets up golden file testing. Must be called before test execution.
func Setup() {
	safeDelete(binaryName)

	// Compile testing binary
	comp := exec.Command("go", "build", "-trimpath", "-o", binaryName)
	if output, err := comp.CombinedOutput(); err != nil {
		panic(fmt.Errorf(
			"error compiling testing binary: %v\n%s",
			err,
			output,
		))
	}
}

// Reset cleans up golden file testing. Must be called after test execution.
// Keep is a map of files that will not be deleted.
func Reset(keep []string) {
	keepMap := make(map[string]struct{})
	for _, k := range keep {
		keepMap[k] = struct{}{}
	}
	// Remove all files and directories, except for files in keep.
	entries, _ := os.ReadDir(".")
	for _, e := range entries {
		_, keep := keepMap[e.Name()]
		// Files that start with the prefix "__" debug files and automatically
		// removed when debugging finishes. This to avoid trying to remove
		// files that suddenly do not exist anymore.
		if !keep && !strings.HasPrefix(e.Name(), "__") {
			switch e.IsDir() {
			case true:
				if err := os.RemoveAll(e.Name()); err != nil {
					panic(err)
				}
			case false:
				if err := os.Remove(e.Name()); err != nil {
					panic(err)
				}
			}
		}
	}
}

// Teardown cleans up golden file testing. Must be called after test execution.
// Passing multiple directories is optional and it will clean up all
// directories passed to it.
func Teardown(dirs ...string) {
	// Clean up the binary file
	safeDelete(binaryName)

	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(err)
		}
	}
}

// TemplateTeardown cleans up the template directory. Must be called after test
// execution.
func TemplateTeardown(templateName string) {
	entries, _ := os.ReadDir(".")
	for _, e := range entries {
		isTemplateDir, _ := regexp.
			MatchString(fmt.Sprintf("%s\\S*", templateName), e.Name())
		if e.IsDir() && isTemplateDir {
			if err := os.RemoveAll(e.Name()); err != nil {
				panic(err)
			}
		}
	}
}

// GenerateTemplate generates a template with the given name. If destination is
// not empty, it will move the template to the destination directory.
func GenerateTemplate(templateName, destination string) {
	args := []string{
		"nextmv", "community", "clone", "-a", templateName,
	}

	if destination != "" {
		args = append(args, "-d", destination)
	}

	output, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		panic(fmt.Errorf(
			"error generating template: %v\n%s",
			err,
			output,
		))
	}

	// This is a temporary hack that is needed to align the SDK version used in
	// the generated community app with the SDK version that is used by the
	// compiled plugins. When plugins disappear, this can be removed.
	sdkRef := os.Getenv("SDK_REF")
	file := osPath.Join(destination, templateName, "go.mod")
	if _, err := os.Stat(file); sdkRef != "" && err == nil {
		cmd := exec.Command(
			"nextmv", "go", "get",
			fmt.Sprintf("github.com/nextmv-io/sdk@%s", sdkRef),
		)
		cmd.Dir = osPath.Join(destination, templateName)
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}

// CopyFile copies a file from the original path to the new path.
func CopyFile(original, newFile string) {
	bytesRead, err := os.ReadFile(original)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(newFile, bytesRead, 0o755)
	if err != nil {
		log.Fatal(err)
	}
}

func safeDelete(file string) {
	if _, err := os.Stat(file); err == nil {
		if err = os.Remove(file); err != nil {
			panic(err)
		}
	}
}
