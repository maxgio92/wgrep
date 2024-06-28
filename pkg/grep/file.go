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
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	d "github.com/gocolly/colly/debug"
	"github.com/pkg/errors"
)

// crawlFiles returns a list of file names found from the seed URL, filtered by file name regex.
//
//nolint:funlen,cyclop
func (o *Options) crawlFiles() (*Result, error) {
	seeds := []*url.URL{}

	err := o.Validate()
	if err != nil {
		return nil, err
	}

	for _, v := range o.SeedURLs {
		u, _ := url.Parse(v)

		seeds = append(seeds, u)
	}

	var contains func(string, string) bool
	if o.CaseInsensitive {
		contains = func(text, pattern string) bool {
			return strings.Contains(strings.ToLower(text), strings.ToLower(pattern))
		}
	} else {
		contains = func(text, pattern string) bool {
			return strings.Contains(strings.ToLower(text), strings.ToLower(pattern))
			//return strings.Contains(text, pattern)
		}
	}

	var files, urls []string

	allowedDomains := getHostnamesFromURLs(seeds)

	// Create the collector settings
	coOptions := []func(*colly.Collector){
		colly.AllowedDomains(allowedDomains...),
		colly.Async(o.Async),
		colly.MaxBodySize(o.MaxBodySize),
	}

	if o.Verbose {
		coOptions = append(coOptions, colly.Debugger(&d.LogDebugger{}))
	}

	// Create the collector.
	co := colly.NewCollector(coOptions...)
	if o.ClientTransport != nil {
		co.WithTransport(o.ClientTransport)
	}

	// Add the callback to Visit the linked resource, for each HTML element found
	co.OnHTML(HTMLTagLink, func(e *colly.HTMLElement) {
		href := e.Attr(HTMLAttrRef)

		includeLinked := true
		if o.IncludeRegexp != "" {
			includeLinked = false
			if match, _ := regexp.MatchString(o.IncludeRegexp, href); match == true {
				includeLinked = true
			}
		}

		// Traverse the folder hierarchy in top-down order.
		if o.Recursive && !(strings.Contains(href, UpDir)) && href != RootDir && includeLinked {
			//nolint:errcheck
			co.Visit(e.Request.AbsoluteURL(e.Attr(HTMLAttrRef)))
		}
	})

	// Manage errors.
	co.OnError(o.handleError)

	co.OnResponse(func(r *colly.Response) {
		// Parse the HTML content with Goquery
		doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
		if err != nil {
			fmt.Println("Error parsing HTML:", err)
			return
		}

		// Grep in the filtered elements the pattern.
		doc.Find(o.ElementFilter).Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if contains(text, o.Pattern) {
				fmt.Printf("%s: %s\n\n", r.Request.URL.String(), text)
			}
		})
	})

	// Visit each root folder.
	for _, seedURL := range seeds {
		err := co.Visit(seedURL.String())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error scraping file with URL %s", seedURL.String()))
		}
	}

	// Wait until colly goroutines are finished.
	co.Wait()

	return &Result{BaseNames: files, URLs: urls}, nil
}
