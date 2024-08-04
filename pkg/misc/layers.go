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

package misc

import (
	"fmt"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Layer is a go-containerregistry v1.Layer equivalent, but with optional history attached
type Layer struct {
	v1.Layer

	LayerIdx *int
	History  *v1.History
}

func Layers(image v1.Image) ([]Layer, error) {
	configFile, err := image.ConfigFile()
	if err != nil {
		return nil, err
	}

	// TODO Support returning layers without history
	if len(configFile.History) == 0 {
		return nil, fmt.Errorf("no history available")
	}

	layers, err := image.Layers()
	if err != nil {
		return nil, err
	}

	var ptr = func(i int) *int { return &i }

	var result []Layer

	var nonEmptyLayerIdx int
	for i := len(configFile.History) - 1; i >= 0; i-- {
		var entry = configFile.History[i]

		if entry.EmptyLayer {
			result = append(result, Layer{
				History: &entry,
			})

		} else {
			result = append(result, Layer{
				LayerIdx: ptr(nonEmptyLayerIdx),
				Layer:    layers[nonEmptyLayerIdx],
				History:  &entry,
			})

			nonEmptyLayerIdx++
		}
	}

	return result, nil
}
