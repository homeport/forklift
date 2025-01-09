// Copyright © 2024 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package interactive

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func Edit(input string) (output string, err error) {
	tmpFile, err := os.CreateTemp("", "interactive-edit")
	if err != nil {
		return "", err
	}

	var filename = tmpFile.Name()

	if err := os.WriteFile(filename, []byte(input), os.FileMode(0600)); err != nil {
		return "", err
	}

	if err := tmpFile.Close(); err != nil {
		return "", err
	}

	editor, args, err := editor()
	if err != nil {
		return "", err
	}

	// TODO Introduce {} style replacement option
	args = append(args, filename)

	fmt.Fprintf(os.Stderr, "Running: %s %v\n", editor, args)
	var cmd = exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func editor() (name string, args []string, err error) {
	// Helper function to split a given input string into command and arguments
	var split = func(in string) (string, []string, error) {
		var ifs, ok = os.LookupEnv("IFS")
		if !ok {
			ifs = " \t\n"
		}

		var runes = []rune(ifs)
		var parts = strings.FieldsFunc(in, func(r rune) bool {
			return slices.Contains(runes, r)
		})

		switch len(parts) {
		case 0:
			return "", nil, fmt.Errorf("failed to split %q into command and arguments", in)

		case 1:
			return parts[0], nil, nil

		default:
			return parts[0], parts[1:], nil
		}
	}

	// Check the EDITOR environment variable
	if editor, ok := os.LookupEnv("EDITOR"); ok && editor != "" {
		return split(editor)
	}

	// If EDITOR is not set or not found, try fallback editors
	for _, editor := range []string{"vim", "vi", "nano"} {
		// Check if the fallback editor is available in the PATH
		if _, err := exec.LookPath(editor); err == nil {
			return split(editor)
		}
	}

	// No suitable editor found
	return "", nil, fmt.Errorf("no suitable editor found, use EDITOR environment variable to set one")
}
