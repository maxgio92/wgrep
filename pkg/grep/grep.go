/*
Copyright © 2023 maxgio92 me@maxgio.me

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

package grep

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	log "github.com/rs/zerolog"

	"github.com/maxgio92/wgrep/internal/network"
)

// Result represents the output of the Grep job.
type Result struct {
	// BaseNames are the path base of the files found.
	BaseNames []string

	// URLs are the universal resource location of the files found.
	URLs []string
}

// Options represents the options for the Grep job.
type Options struct {
	// SeedURLs are the URLs used as root URLs from which to scrape.
	SeedURLs []string

	// Pattern is a string pattern to verify against web page content.
	Pattern string

	// Recursive enables the Grep job to examine files referenced to by the seeds files recursively.
	Recursive bool

	// CaseInsensitive compare content with page content case-insensitive.
	CaseInsensitive bool

	// ElementFilter is the HTML element to filter while looking for the pattern.
	ElementFilter string

	// IncludeRegexp is the string that page URL should match
	IncludeRegexp string

	// Verbose enables the Grep job verbosity printing every visited URL.
	Verbose bool

	Logger log.Logger

	// Async represetns the option to scrape with multiple asynchronous coroutines.
	Async bool

	// ClientTransport represents the Transport used for the HTTP client.
	ClientTransport http.RoundTripper

	// MaxBodySize is the limit in bytes of each of the retrieved response body.
	MaxBodySize int

	// ContextDeadlineRetryBackOff controls the error handling on responses.
	// If not nil, when the request context deadline exceeds, the request
	// is retried with an exponential backoff interval.
	ContextDeadlineRetryBackOff *ExponentialBackOffOptions

	// ConnResetRetryBackOff controls the error handling on responses.
	// If not nil, when the connection is reset by the peer (TCP RST), the request
	// is retried with an exponential backoff interval.
	ConnResetRetryBackOff *ExponentialBackOffOptions

	// TimeoutRetryBackOff controls the error handling on responses.
	// If not nil, when the connection times out (based on client timeout), the request
	// is retried with an exponential backoff interval.
	TimeoutRetryBackOff *ExponentialBackOffOptions
}

type Option func(opts *Options)

func WithSeedURLs(seedURLs []string) Option {
	return func(opts *Options) {
		opts.SeedURLs = seedURLs
	}
}

func WithPattern(pattern string) Option {
	return func(opts *Options) {
		opts.Pattern = pattern
	}
}

func WithRecursive(recursive bool) Option {
	return func(opts *Options) {
		opts.Recursive = recursive
	}
}

func WithElementFilter(element string) Option {
	return func(opts *Options) {
		opts.ElementFilter = element
	}
}

func WithCaseInsensitive(insensitive bool) Option {
	return func(opts *Options) {
		opts.CaseInsensitive = insensitive
	}
}

func WithIncludeRegexp(include string) Option {
	return func(opts *Options) {
		opts.IncludeRegexp = include
	}
}

func WithVerbosity(verbosity bool) Option {
	return func(opts *Options) {
		opts.Verbose = verbosity
	}
}

func WithLogger(logger log.Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}

func WithAsync(async bool) Option {
	return func(opts *Options) {
		opts.Async = async
	}
}

func WithClientTransport(transport http.RoundTripper) Option {
	return func(opts *Options) {
		opts.ClientTransport = transport
	}
}

func WithMaxBodySize(maxBodySize int) Option {
	return func(opts *Options) {
		opts.MaxBodySize = maxBodySize
	}
}

func WithContextDeadlineRetryBackOff(backoff *ExponentialBackOffOptions) Option {
	return func(opts *Options) {
		opts.ContextDeadlineRetryBackOff = backoff
	}
}

func WithConnResetRetryBackOff(backoff *ExponentialBackOffOptions) Option {
	return func(opts *Options) {
		opts.ConnResetRetryBackOff = backoff
	}
}

func WithConnTimeoutRetryBackOff(backoff *ExponentialBackOffOptions) Option {
	return func(opts *Options) {
		opts.TimeoutRetryBackOff = backoff
	}
}

// NewGrep returns a new Grep object to grep files over HTTP and HTTPS.
func NewGrep(opts ...Option) *Options {
	o := &Options{}

	for _, f := range opts {
		f(o)
	}

	o.init()

	return o
}

// Validate validates the Grep job options and returns an error.
func (o *Options) init() {
	if o.ClientTransport == nil {
		o.ClientTransport = network.DefaultClientTransport
	}
	if o.MaxBodySize == 0 {
		// Set max body size to 100 KB.
		o.MaxBodySize = 100 * 1024
	}
}

// Validate validates the Grep job options and returns an error.
func (o *Options) Validate() error {
	if o.Pattern == "" {
		return ErrPatternNotSpecified
	}
	if len(o.SeedURLs) == 0 {
		return ErrSeedURLsNotSpecified
	}

	for k, v := range o.SeedURLs {
		_, err := url.Parse(v)
		if err != nil {
			return ErrSeedURLNotValid
		}

		if !strings.HasSuffix(v, "/") {
			o.SeedURLs[k] = v + "/"
		}
	}

	return nil
}

func (o *Options) Grep() (*Result, error) {
	if err := o.Validate(); err != nil {
		return nil, errors.Wrap(err, "error validating grep options")
	}

	return o.crawlFiles()
}
