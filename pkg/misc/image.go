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
	"context"
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func LoadImage(ctx context.Context, ref name.Reference) (v1.Image, error) {
	if image, err := daemon.Image(ref, daemon.WithContext(ctx)); err == nil {
		return image, nil
	}

	auth, err := authn.DefaultKeychain.ResolveContext(ctx, ref.Context())
	if err != nil {
		return nil, err
	}

	return remote.Image(ref,
		remote.WithContext(ctx),
		remote.WithAuth(auth),
	)
}

func SaveImage(tag name.Tag, img v1.Image, options ...daemon.Option) error {
	response, err := daemon.Write(tag, img, options...)
	if err != nil {
		fmt.Fprintln(os.Stderr, response)
		return fmt.Errorf("failed to write image: %w", err)
	}

	return nil
}
