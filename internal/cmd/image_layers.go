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
	"github.com/homeport/forklift/pkg/misc"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var imageLayersCmd = &cobra.Command{
	Use:          "layers <image-reference>",
	Args:         cobra.MinimumNArgs(1),
	Short:        "List layers",
	Long:         `Lists layers of given image`,
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

		layers, err := misc.Layers(image)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorder(false)
		table.SetHeaderLine(true)
		table.SetCenterSeparator("┼")
		table.SetColumnSeparator("│")
		table.SetRowSeparator("─")
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)
		table.SetAutoFormatHeaders(false)

		table.SetHeader([]string{
			"Layer",
			"Size",
			"Created",
			"CreatedBy",
			"Comment",
			"Author",
		})

		for _, layer := range layers {
			if layer.Layer == nil {
				continue
			}

			size, err := layer.Size()
			if err != nil {
				return err
			}

			var createdBy, comment, author, created string
			if layer.History != nil {
				createdBy = layer.History.CreatedBy
				comment = layer.History.Comment
				author = layer.History.Author
				created = layer.History.Created.String()
			}

			table.Append([]string{
				fmt.Sprintf("%d", *layer.LayerIdx),
				humanReadableSize(size),
				created,
				createdBy,
				comment,
				author,
			})
		}

		table.Render()
		return nil
	},
}

func init() {
	imageCmd.AddCommand(imageLayersCmd)
}
