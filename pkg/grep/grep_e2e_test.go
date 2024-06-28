//go:build all_tests || e2e_tests

/*
Copyright © 2024 maxgio92 me@maxgio.me

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package grep_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/maxgio92/wgrep/internal/network"
	grep "github.com/maxgio92/wgrep/pkg/grep"
)

const (
	seedURL         = "https://mirrors.edge.kernel.org/centos/8-stream"
	fileRegexp      = "repomd.xml$"
	expectedResults = 155
)

var _ = Describe("File crawling", func() {
	Context("Async", func() {
		var (
			search = grep.NewFind(
				grep.WithAsync(true),
				grep.WithSeedURLs([]string{seedURL}),
				grep.WithClientTransport(network.DefaultClientTransport),
				grep.WithRecursive(true),
				grep.WithMaxBodySize(grep.DefaultMaxBodySize),
				grep.WithContextDeadlineRetryBackOff(grep.DefaultExponentialBackOffOptions),
				grep.WithConnTimeoutRetryBackOff(grep.DefaultExponentialBackOffOptions),
				grep.WithConnResetRetryBackOff(grep.DefaultExponentialBackOffOptions),
			)
			actual        *grep.Result
			err           error
			expectedCount = expectedResults
		)
		BeforeEach(func() {
			actual, err = search.Find()
		})
		It("Should not fail", func() {
			Expect(err).To(BeNil())
		})
		It("Should stage results", func() {
			Expect(actual.URLs).ToNot(BeEmpty())
			Expect(actual.URLs).ToNot(BeNil())
		})
		It("Should stage exact result count", func() {
			Expect(len(actual.URLs)).To(Equal(expectedCount))
		})
	})
})
