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

package repackage

import (
	"compress/gzip"
	"fmt"
	"os"
	"strings"

	"github.com/homeport/forklift/pkg/tar"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

type Intention string

const (
	PICK  Intention = "pick"
	FIXUP Intention = "fixup"
)

type Action struct {
	OriginalIdx int
	Intent      Intention
	Layer       v1.Layer
	History     *v1.History
}

type Plan []Action

func Image(input v1.Image, plan Plan) (v1.Image, error) {
	configFile, err := input.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to read image config file: %w", err)
	}

	// reset config details so that it can be used in a fresh image
	configFile.RootFS.DiffIDs = []v1.Hash{}
	configFile.History = []v1.History{}

	// create a fresh empty image using the input image's config file
	result, err := mutate.ConfigFile(empty.Image, configFile)
	if err != nil {
		return nil, err
	}

	type repkgStage struct {
		layer     *v1.Layer
		history   *v1.History
		directory *string
		created   *v1.Time
		createdBy []string
	}

	var stage *repkgStage

	var flush = func() (err error) {
		if stage == nil {
			return nil
		}

		// TODO Remove me eventually
		if stage.directory == nil && stage.layer == nil {
			panic("invalid state")
		}

		if stage.directory != nil {
			// TODO Remove temporary file at the end, defer won't work
			tmpball, err := tar.Create(*stage.directory)
			if err != nil {
				return err
			}

			stage.directory = nil

			// TODO Make compression configurable
			layer, err := tarball.LayerFromFile(tmpball.Name(), tarball.WithCompressionLevel(gzip.DefaultCompression))
			if err != nil {
				return err
			}

			stage.layer = &layer
			stage.history = &v1.History{
				Author:    "forklift",
				Comment:   "combined layers",
				Created:   *stage.created,
				CreatedBy: strings.Join(stage.createdBy, ", "),
			}
		}

		addendum := mutate.Addendum{Layer: *stage.layer}

		if stage.history != nil {
			addendum.History = *stage.history
		}

		result, err = mutate.Append(result, addendum)
		stage = nil
		return err
	}

	for i := range plan {
		var action = plan[i]
		switch action.Intent {

		// TODO Write comment
		case PICK:
			if err := flush(); err != nil {
				return nil, err
			}

			stage = &repkgStage{
				layer:   &action.Layer,
				history: action.History,
			}

		// TODO Write comment
		case FIXUP:
			if action.Layer == nil {
				return nil, fixupOnEmptyLayer
			}

			if stage.directory == nil {
				dir, err := os.MkdirTemp("", "fixup")
				if err != nil {
					return nil, err
				}

				stage.directory = &dir

				if err := tar.ExtractLayer(*stage.layer, *stage.directory); err != nil {
					return nil, err
				}

				stage.createdBy = append(stage.createdBy, stage.history.CreatedBy)
			}

			if err := tar.ExtractLayer(action.Layer, *stage.directory); err != nil {
				return nil, err
			}

			if action.History != nil {
				stage.created = &action.History.Created
				stage.createdBy = append(stage.createdBy, action.History.CreatedBy)
			}
		}
	}

	if err := flush(); err != nil {
		return nil, err
	}

	return result, nil
}
