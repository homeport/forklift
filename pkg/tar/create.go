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

package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Create(directory string) (*os.File, error) {
	target, err := os.CreateTemp("", "tarball")
	if err != nil {
		return nil, err
	}

	tw := tar.NewWriter(target)
	defer func() {
		tw.Flush()
		tw.Close()
	}()

	return target, filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		switch {
		case info.Mode().IsDir():
			return tw.WriteHeader(header)

		case info.Mode().IsRegular():
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			return write(tw, path)

		case info.Mode()&os.ModeSymlink == os.ModeSymlink:
			deref, info, err := followSymLink(path)
			if err != nil {
				return err
			}

			header, err = tar.FileInfoHeader(info, deref)
			if err != nil {
				return err
			}

			header.Name, err = filepath.Rel(directory, path)
			if err != nil {
				return err
			}

			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			return write(tw, deref)

		default:
			return fmt.Errorf("unsupported file type: %s", path)
		}
	})
}

func write(w io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(w, file)
	return err
}
