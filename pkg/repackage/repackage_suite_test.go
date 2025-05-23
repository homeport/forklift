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
	"bytes"
	"fmt"
	"math/rand/v2"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"gopkg.in/yaml.v3"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
)

var lcs = []rune("abcdefghijklmnopqrstuvwxyz")

func TestRepackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repackage Suite")
}

func random(n int) string {
	var buf = make([]rune, n)
	for i := range buf {
		buf[i] = lcs[rand.IntN(len(lcs))]
	}

	return string(buf)
}

func pullDaemonImage(s string) v1.Image {
	GinkgoHelper()

	ref, err := name.ParseReference(s)
	Expect(err).ToNot(HaveOccurred())

	image, err := daemon.Image(ref)
	Expect(err).ToNot(HaveOccurred())

	return image
}

func pushDaemonImage(tag name.Tag, img v1.Image, options ...daemon.Option) {
	GinkgoHelper()

	response, err := daemon.Write(tag, img, options...)
	Expect(err).ToNot(HaveOccurred(), response)
}

func BeImage(expected v1.Image) types.GomegaMatcher {
	return &BeImageMatcher{expected: expected}
}

type BeImageMatcher struct {
	expected v1.Image
	report   dyff.Report
}

func (matcher *BeImageMatcher) Match(actual interface{}) (bool, error) {
	switch actual := actual.(type) {
	case v1.Image:
		var translate = func(obj any) (*yaml.Node, error) {
			out, err := yaml.Marshal(obj)
			if err != nil {
				return nil, err
			}

			var node yaml.Node
			if err := yaml.Unmarshal(out, &node); err != nil {
				return nil, err
			}

			return &node, nil
		}

		aCfgFile, err := actual.ConfigFile()
		if err != nil {
			return false, err
		}

		aCfgFileNode, err := translate(aCfgFile)
		if err != nil {
			return false, err
		}

		bCfgFile, err := matcher.expected.ConfigFile()
		if err != nil {
			return false, err
		}

		bCfgFileNode, err := translate(bCfgFile)
		if err != nil {
			return false, err
		}

		matcher.report, err = dyff.CompareInputFiles(
			ytbx.InputFile{Documents: []*yaml.Node{
				bCfgFileNode,
			}},
			ytbx.InputFile{Documents: []*yaml.Node{
				aCfgFileNode,
			}},
		)

		if err != nil {
			return false, err
		}

		return len(matcher.report.Diffs) == 0, nil

	default:
		return false, fmt.Errorf("BeImage matcher expected an image, not %T", actual)
	}
}

func (matcher *BeImageMatcher) FailureMessage(actual interface{}) string {
	reporter := &dyff.HumanReport{
		Report:          matcher.report,
		Indent:          2,
		NoTableStyle:    true,
		OmitHeader:      true,
		UseGoPatchPaths: true,
		PrefixMultiline: true,
	}

	var buf bytes.Buffer
	if err := reporter.WriteReport(&buf); err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("Expected images to match, but differences were found comparing expected with actual: %s", buf.String())
}

func (matcher *BeImageMatcher) NegatedFailureMessage(actual interface{}) string {
	return "Expected images not to match, but no differences were found comparing expected with actual"
}
