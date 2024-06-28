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
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/maxgio92/wgrep/internal/network"
	"github.com/maxgio92/wgrep/internal/output"
	"github.com/maxgio92/wgrep/pkg/grep"
)

type Command struct {
	ConnectionTimeout   int
	KeepAliveInterval   int
	TLSHandshakeTimeout int
	IdleConnTimeout     int
	ConnPoolSize        int
	ConnPoolPerHostSize int
	*grep.Options
}

// NewCmd returns a new grep command.
func NewCmd() *cobra.Command {
	o := &Command{
		Options: &grep.Options{},
	}

	cmd := &cobra.Command{
		Use:               "wgrep PATTERN URL",
		Short:             "Wgrep print lines that match patterns in web pages",
		DisableAutoGenTag: true,
		Args:              cobra.MinimumNArgs(2),
		RunE:              o.Run,
	}

	// General flags.
	cmd.Flags().BoolVarP(&o.CaseInsensitive, "ignore-case", "i", false,
		"Whether to search for the pattern case insensitive")
	cmd.Flags().StringVar(&o.ElementFilter, "element", "p",
		"The HTLM element to filer on while searching for the pattern")
	cmd.Flags().StringVar(&o.IncludeRegexp, "include", "",
		"Search only pages whose URL matches teh include regular expression.")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", false,
		"Inspect all web pages recursively by following each hypertext reference.")
	cmd.Flags().BoolVar(&o.Async, "async", true,
		"Whether to scrape with asynchronous jobs.")
	cmd.Flags().BoolVarP(&o.Verbose, "verbose", "v", false,
		"Enable verbosity to log all visited HTTP(s) files")

	// Timeouts flags.
	cmd.Flags().IntVar(&o.ConnectionTimeout, "connection-timeout", network.DefaultTimeout,
		"The maximum amount of time in milliseconds a dial will wait for a connect to complete.")
	cmd.Flags().IntVar(&o.KeepAliveInterval, "keep-alive-interval", network.DefaultKeepAlive,
		"The interval between keep-alive probes for an active network connection.")
	cmd.Flags().IntVar(&o.TLSHandshakeTimeout, "tls-handshake-timeout", network.DefaultTLSHandshakeTimeout,
		"The maximum amount of time in milliseconds a connection will wait for a TLS handshake.")
	cmd.Flags().IntVar(&o.IdleConnTimeout, "idle-connection-timeout", network.DefaultIdleConnTimeout,
		"The maximum amount of time in milliseconds a connection will remain idle before closing itself.")

	// Sizes flags.
	cmd.Flags().IntVar(&o.ConnPoolSize, "connection-pool-size", network.DefaultMaxIdleConns,
		"The maximum number of idle connections across all hosts.")
	cmd.Flags().IntVar(&o.ConnPoolPerHostSize, "connection-pool-size-per-host", network.DefaultMaxIdleConnsPerHost,
		"The maximum number of idle connections across for each host.")
	cmd.Flags().IntVar(&o.MaxBodySize, "max-body-size", grep.DefaultMaxBodySize,
		"The maximum size in bytes a response body is read for each request.")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := NewCmd()
	output.ExitOnErr(cmd.Execute())
}

func (o *Command) validate() error {
	if err := o.Validate(); err != nil {
		return errors.Wrap(err, "error validating Command")
	}

	return nil
}

func (o *Command) Run(_ *cobra.Command, args []string) error {
	var pattern, seed string
	if len(args) > 1 {
		pattern = args[0]
		seed = args[1]
	}

	o.Pattern = pattern
	o.SeedURLs = append(o.SeedURLs, seed)

	if err := o.validate(); err != nil {
		return err
	}

	// Network client dialer.
	dialer := network.NewDialer(
		network.WithTimeout(time.Duration(o.ConnectionTimeout)*time.Millisecond),
		network.WithKeepAlive(time.Duration(o.KeepAliveInterval)*time.Millisecond),
	)

	// HTTP client transport.
	transport := network.NewTransport(
		network.WithDialer(dialer),
		network.WithIdleConnsTimeout(time.Duration(o.IdleConnTimeout)*time.Millisecond),
		network.WithTLSHandshakeTimeout(time.Duration(o.TLSHandshakeTimeout)*time.Millisecond),
		network.WithMaxIdleConns(o.ConnPoolSize),
		network.WithMaxIdleConnsPerHost(o.ConnPoolPerHostSize),
	)

	// Wgrep matcher.
	matcher := grep.NewGrep(
		grep.WithSeedURLs(o.SeedURLs),
		grep.WithPattern(o.Pattern),
		grep.WithCaseInsensitive(o.CaseInsensitive),
		grep.WithElementFilter(o.ElementFilter),
		grep.WithIncludeRegexp(o.IncludeRegexp),
		grep.WithRecursive(o.Recursive),
		grep.WithVerbosity(o.Verbose),
		grep.WithAsync(o.Async),
		grep.WithMaxBodySize(o.MaxBodySize),
		grep.WithClientTransport(transport),
		grep.WithContextDeadlineRetryBackOff(grep.DefaultExponentialBackOffOptions),
		grep.WithConnResetRetryBackOff(grep.DefaultExponentialBackOffOptions),
		grep.WithConnTimeoutRetryBackOff(grep.DefaultExponentialBackOffOptions),
	)

	found, err := matcher.Grep()
	if err != nil {
		return errors.Wrap(err, "error finding the pattern")
	}

	for _, v := range found.URLs {
		output.Print(v)
	}

	return nil
}
