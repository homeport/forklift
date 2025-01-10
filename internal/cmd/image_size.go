// Copyright © 2025 The Homeport Team
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
	"io"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/homeport/forklift/pkg/misc"
	"github.com/spf13/cobra"
)

type imageSizeCmdOpts struct {
	humanReadable bool
	uncompressed  bool
}

var imageSizeCmdSettings imageSizeCmdOpts

var imageSizeCmd = &cobra.Command{
	Use:          "size <image-reference>",
	Args:         cobra.MinimumNArgs(1),
	Short:        "Determine image size",
	Long:         `Determine the image size of a given image`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ref, err := name.ParseReference(args[0])
		if err != nil {
			return err
		}

		image, err := misc.LoadImage(cmd.Context(), ref)
		if err != nil {
			return err
		}

		size, err := image.Size()
		if err != nil {
			return err
		}

		layers, err := misc.Layers(image)
		if err != nil {
			return err
		}

		for _, layer := range layers {
			if layer.Layer == nil {
				continue
			}

			if imageSizeCmdSettings.uncompressed {
				r, err := layer.Uncompressed()
				if err != nil {
					return err
				}

				written, err := io.Copy(io.Discard, r)
				if err != nil {
					return err
				}

				size += written

			} else {
				layerSize, err := layer.Size()
				if err != nil {
					return err
				}

				size += layerSize
			}
		}

		if imageSizeCmdSettings.humanReadable {
			fmt.Println(humanReadableSize(size))

		} else {
			fmt.Printf("%d\n", size)
		}

		return nil
	},
}

func init() {
	imageCmd.AddCommand(imageSizeCmd)

	imageSizeCmd.Flags().BoolVarP(&imageSizeCmdSettings.humanReadable, "human-readable", "H", false, "Show sizes in human readable ranges")
	imageSizeCmd.Flags().BoolVarP(&imageSizeCmdSettings.uncompressed, "uncompressed", "u", false, "Determine the uncompressed image size")
}
