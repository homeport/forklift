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

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/homeport/forklift/pkg/misc"
	"github.com/spf13/cobra"
)

var lookupMediaType = &cobra.Command{
	Use:          "mediatype <reference>",
	Aliases:      []string{"media-type"},
	Args:         cobra.MinimumNArgs(1),
	Short:        "Look-up media type",
	Long:         `Look-up media type of given index or image reference`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ref, err := name.ParseReference(args[0])
		if err != nil {
			return err
		}

		opts, err := misc.RemoteOptionsFromRef(cmd.Context(), ref)
		if err != nil {
			return err
		}

		desc, err := remote.Head(ref, opts...)
		if err != nil {
			return err
		}

		fmt.Println(desc.MediaType)

		return nil
	},
}

func init() {
	lookupCmd.AddCommand(lookupMediaType)
}
