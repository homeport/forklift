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

package repackage_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/homeport/forklift/pkg/misc"
	"github.com/homeport/forklift/pkg/repackage"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/empty"
)

var _ = Describe("Repackage", func() {
	It("should produce the same output image as the input image if all layers are picked", func() {
		sampleImage := pullDaemonImage("test:me")

		layers, err := misc.Layers(sampleImage)
		Expect(err).ToNot(HaveOccurred())

		var plan = repackage.Plan{}
		for i := range layers {
			plan = append(plan, repackage.Action{
				OriginalIdx: i,
				Intent:      repackage.PICK,
				Layer:       layers[i].Layer,
				History:     layers[i].History,
			})
		}

		result, err := repackage.Image(sampleImage, plan)
		Expect(err).ToNot(HaveOccurred())

		tag, err := name.NewTag("test:" + random(6))
		Expect(err).ToNot(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "result image tag: %s\n", tag)

		pushDaemonImage(tag, result)
		defer pushDaemonImage(tag, empty.Image)

		result = pullDaemonImage(tag.String())
		Expect(result).To(BeImage(sampleImage))
	})

	It("should combine layers into one", func() {
		sampleImage := pullDaemonImage("test:me")

		layers, err := misc.Layers(sampleImage)
		Expect(err).ToNot(HaveOccurred())

		var tmp = repackage.Plan{}
		for i := range layers {
			var t = repackage.PICK
			if strings.Contains(layers[i].History.CreatedBy, "COPY run") {
				t = repackage.FIXUP
			}

			tmp = append(tmp, repackage.Action{
				OriginalIdx: i,
				Intent:      t,
				Layer:       layers[i].Layer,
				History:     layers[i].History,
			})
		}

		var plan repackage.Plan
		var envs []repackage.Action
		for i := range tmp {
			if strings.HasPrefix(tmp[i].History.CreatedBy, "ENV") {
				envs = append(envs, tmp[i])
			} else {
				plan = append(plan, tmp[i])
			}
		}

		plan = append(plan, envs...)

		fmt.Fprintf(GinkgoWriter, "replackage plan (%d entries)\n", len(plan))
		for i := range plan {
			fmt.Fprintf(GinkgoWriter, "  %s layer=%d (%s)\n",
				plan[i].Intent,
				plan[i].OriginalIdx,
				plan[i].History.CreatedBy,
			)
		}

		result, err := repackage.Image(sampleImage, plan)
		Expect(err).ToNot(HaveOccurred())

		tag, err := name.NewTag("test:" + random(6))
		Expect(err).ToNot(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "result image tag: %s\n", tag)

		pushDaemonImage(tag, result)
		defer pushDaemonImage(tag, empty.Image)

		result = pullDaemonImage(tag.String())

		layers, err = misc.Layers(result)
		Expect(err).ToNot(HaveOccurred())
		Expect(layers).To(HaveLen(4))
	})
})
