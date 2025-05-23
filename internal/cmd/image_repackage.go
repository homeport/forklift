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
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/homeport/forklift/pkg/interactive"
	"github.com/homeport/forklift/pkg/misc"
	"github.com/homeport/forklift/pkg/repackage"
	"github.com/spf13/cobra"
)

var repackageCmdSettings struct {
	interactive bool
	target      tag
}

// repackageCmd represents the repackage command
var repackageCmd = &cobra.Command{
	Use:   "repackage",
	Args:  cobra.MinimumNArgs(1),
	Short: "Repackage layers of an image",
	Long:  `Repackage is similar to Git rebase, but for container image layers instead of commits.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ref, err := name.ParseReference(args[0])
		if err != nil {
			return err
		}

		if repackageCmdSettings.target.RegistryStr() == "" {
			defaultTag, err := name.NewTag(ref.String() + "-repackaged")
			if err != nil {
				return err
			}

			repackageCmdSettings.target = tag{defaultTag}
		}

		image, err := misc.LoadImage(cmd.Context(), ref)
		if err != nil {
			return err
		}

		if repackageCmdSettings.interactive {
			layers, err := misc.Layers(image)
			if err != nil {
				return err
			}

			// TODO Extract plan to string and string to plan to helper functions
			var buf bytes.Buffer
			for i, layer := range layers {
				var desc string
				if layer.Layer != nil {
					diffID, err := layer.DiffID()
					if err != nil {
						return err
					}

					desc = diffID.String()
				}

				if desc == "" && layer.History != nil && layer.History.EmptyLayer {
					desc = "(empty layer)"
				}

				fmt.Fprintf(&buf, "%-6s %3d %s\n", repackage.PICK, i, desc)
			}

			planText, err := interactive.Edit(buf.String())
			if err != nil {
				return err
			}

			var plan repackage.Plan
			var scanner = bufio.NewScanner(strings.NewReader(planText))
			for scanner.Scan() {
				parts := strings.Fields(scanner.Text())
				if len(parts) < 2 {
					return fmt.Errorf("plan file entry doesn't match the expected format of <intention> <layer>")
				}

				intent := repackage.Intention(parts[0])

				idx, err := strconv.Atoi(parts[1])
				if err != nil {
					return err
				}

				// TODO Add bound checks
				layer := layers[idx]

				plan = append(plan, repackage.Action{
					Intent:  intent,
					Layer:   layer.Layer,
					History: layer.History,
				})
			}

			if err := scanner.Err(); err != nil {
				return err
			}

			pout("repackage plan (%d entries)\n", len(plan))
			for i := range plan {
				pout("  %s layer=%d (%s)\n",
					plan[i].Intent,
					plan[i].OriginalIdx,
					plan[i].History.CreatedBy,
				)
			}

			repackagedImage, err := repackage.Image(image, plan)
			if err != nil {
				return err
			}

			return misc.SaveImage(repackageCmdSettings.target.Tag, repackagedImage)
		}

		return nil
	},
}

func init() {
	imageCmd.AddCommand(repackageCmd)

	repackageCmd.Flags().SortFlags = false

	repackageCmd.Flags().BoolVarP(&repackageCmdSettings.interactive, "interactive", "i", false, "Interactively decide on repackaging")
	repackageCmd.Flags().VarP(&repackageCmdSettings.target, "target", "t", "target")
}
