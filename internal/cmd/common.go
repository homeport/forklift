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

package cmd

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/pflag"
)

type tag struct {
	name.Tag
}

var _ pflag.Value = &tag{}

func (t *tag) String() string {
	return t.Tag.String()
}

func (t *tag) Set(s string) (err error) {
	t.Tag, err = name.NewTag(s)
	return err
}

func (t *tag) Type() string {
	return "image tag"
}

func pout(format string, a ...any) {
	// TODO Make target configurable via root command-line flag
	_, _ = fmt.Fprintf(os.Stdout, format, a...)
}

func humanReadableSize(bytes int64) string {
	var mods = []string{"Byte", "KiB", "MiB", "GiB", "TiB"}

	value := float64(bytes)
	i := 0
	for value > 1023.99999 {
		value /= 1024.0
		i++
	}

	return fmt.Sprintf("%.1f %s", value, mods[i])
}
