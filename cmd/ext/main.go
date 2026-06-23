package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func bailIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func findNpx() (string, error) {
	names := []string{"npx.cmd", "npx.exe", "npx"}
	if runtime.GOOS == "windows" {
		names = []string{"npx.cmd", "npx.exe", "npx"}
	}
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("npx not found in PATH, please install Node.js")
}

func runPrettier(tempFileName string) error {
	tempDir := filepath.Dir(tempFileName)
	args := []string{"--write", "--parser", "babel", tempFileName}

	var cmd *exec.Cmd
	if path, err := exec.LookPath("prettier"); err == nil {
		cmd = exec.Command(path, args...)
	} else if path, err := exec.LookPath("prettier.cmd"); err == nil {
		cmd = exec.Command(path, args...)
	} else {
		npx, err := findNpx()
		if err != nil {
			return err
		}
		// --yes avoids npx blocking on install prompt when run without a TTY (e.g. from go run)
		cmd = exec.Command(npx, append([]string{"--yes", "prettier@3"}, args...)...)
	}

	cmd.Dir = tempDir
	cmd.Env = append(os.Environ(), "NODE_OPTIONS=--max-old-space-size=8192")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	inputPath := flag.String("input", "", "Path to JS file (e.g., extensionHostProcess.js)")
	outputDir := flag.String("output", "", "Output directory for proto files (default: ./cursor_proto)")
	skipFormat := flag.Bool("skip-format", false, "Skip prettier formatting")
	flag.Parse()

	if *inputPath == "" && flag.NArg() > 0 {
		*inputPath = flag.Arg(0)
	}

	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "Usage: ext -input <path-to-js-file> [-output <dir>] [-skip-format]")
		fmt.Fprintln(os.Stderr, "       ext <path-to-js-file>")
		fmt.Fprintln(os.Stderr, "\nExample:")
		fmt.Fprintln(os.Stderr, "  ext -input /path/to/extensionHostProcess.js")
		fmt.Fprintln(os.Stderr, "  ext C:\\Users\\xxx\\scoop\\apps\\cursor\\current\\resources\\app\\out\\vs\\workbench\\api\\node\\extensionHostProcess.js")
		fmt.Fprintln(os.Stderr, "\nTip: copy extensionHostProcess.js out of the Cursor install dir first;")
		fmt.Fprintln(os.Stderr, "     close Cursor if reading directly from its install path.")
		os.Exit(1)
	}

	info, err := os.Stat(*inputPath)
	bailIf(err)

	if info.IsDir() {
		bailIf(fmt.Errorf("expected %s to be file, is dir", *inputPath))
	}

	lowerInput := strings.ToLower(filepath.ToSlash(*inputPath))
	if strings.Contains(lowerInput, "/cursor/") && strings.Contains(lowerInput, "extensionhostprocess.js") {
		fmt.Println("Note: reading from Cursor install dir — copy the file elsewhere first if Cursor is running.")
	}

	sizeMB := float64(info.Size()) / (1024 * 1024)
	if sizeMB > 1 {
		fmt.Printf("Input size: %.1f MB\n", sizeMB)
	}

	if *outputDir == "" {
		wd, err := os.Getwd()
		bailIf(err)
		*outputDir = filepath.Join(wd, "cursor_proto")
	}

	fmt.Println("Copying source file to temp directory...")
	originalFile, err := os.Open(*inputPath)
	bailIf(err)

	tempFile, err := os.CreateTemp(os.TempDir(), "cursor-source-*.js")
	bailIf(err)
	tempFileName := tempFile.Name()

	_, err = io.Copy(tempFile, originalFile)
	bailIf(err)

	bailIf(originalFile.Close())
	bailIf(tempFile.Close())

	fmt.Printf("Temp file: %s\n", tempFileName)

	if !*skipFormat {
		fmt.Println("Formatting file (typically 10–60s for ~5MB)...")
		if err := runPrettier(tempFileName); err != nil {
			fmt.Printf("Warning: formatting failed: %v\n", err)
			fmt.Println("Continuing without format — message extraction may be incomplete. Use --skip-format to skip explicitly.")
		} else {
			fmt.Println("Formatting complete")
		}
	} else {
		fmt.Println("Skipping formatting (--skip-format)")
	}

	fmt.Println("Extracting Proto definitions...")
	ExtractProtos(tempFileName, *outputDir)

	os.Remove(tempFileName)

	fmt.Printf("\nOutput directory: %s\n", *outputDir)
}
