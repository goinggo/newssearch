// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"github.com/goinggo/newssearch/helper"
	"github.com/goinggo/newssearch/plugins/expression"
	"github.com/goinggo/newssearch/rss"
	"github.com/goinggo/utilities/v1/workpool"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//** INTERFACES

// SearchPlugin is implemented by the different search plugins to provide unique functionality
type SearchPlugin interface {
	FindResults(goRoutine string, rssDocument *rss.RSSDocument) (results interface{})
	DisplayResults(goRoutine string, results interface{})
}

//** NEW TYPES

// _Search manages state for the searching engine
type _Search struct {
	WorkPool *workpool.WorkPool // A workPool to process the searches
	FeedList *list.List         // The list of RSS feeds to process
}

// _SearchWork is used to post search work into the work pool
type _SearchWork struct {
	SP           *workpool.WorkPool // Reference to the work pool
	Wait         chan interface{}   // Channel used to receive the result of the search
	Url          string             // The RSS Feed url
	SearchPlugin SearchPlugin       // Reference to the search plugin
}

//** PUBLIC FUNCTIONS

// Run is the processing engine. It is thread safe so multiple searches can be processed
// as the same time
//  goRoutine: The Go routine making the call
//  routines: The number of routines to use to process the search
//  pluginType: The type of search plugin to use
//  parameters: The parameters for the command
func Run(goRoutine string, routines int, pluginType string, parameters []string) {

	defer helper.CatchPanic(nil, goRoutine, "search", "Run")

	helper.WriteStdout(goRoutine, "search", "Run", "Started")

	// Create a search plugin to process the search
	searchPlugin, err := _CreateSearchPlugin(goRoutine, pluginType, parameters)

	if err != nil {

		helper.WriteStdoutf(goRoutine, "search", "Run", "ERROR : %s", err)
		return
	}

	// Capture the list of feeds to process
	feedList, err := _LoadFeedsList(goRoutine)

	if err != nil {

		helper.WriteStdoutf(goRoutine, "search", "Run", "ERROR : %s", err)
		return
	}

	// Create a search object
	search := &_Search{
		WorkPool: workpool.New(routines, int32(feedList.Len())),
		FeedList: feedList,
	}

	// Perform the search
	search.PerformSearch(goRoutine, searchPlugin)

	// Shutdown the search pool
	search.WorkPool.Shutdown(goRoutine)

	helper.WriteStdout(goRoutine, "search", "Run", "Completed")
}

// DisplayHelpExamples displays the examples for each searchPlugin
//  goRoutine: The Go routine making the call
func DisplayHelpExamples(goRoutine string) {

	// TODO: Add new searchPlugin help examples here

	helper.WriteStdoutf(goRoutine, "search", "DisplayHelpExamples", "Example : %s", expression.HelpParameters())
}

//** PRIVATE FUNCTIONS

// _LoadFeedsList reads the feeds.list file and returns the list of feed urls
//  goRoutine: The Go routine making the call
func _LoadFeedsList(goRoutine string) (feedList *list.List, err error) {

	defer helper.CatchPanic(&err, goRoutine, "search", "_LoadFeedsList")

	var file *os.File
	var line string

	helper.WriteStdout(goRoutine, "search", "_LoadFeedsList", "Started")

	// Find the location of the feeds.list file
	strapsFilePath, err := filepath.Abs("feeds.list")

	// Open the feeds.list file
	file, err = os.Open(strapsFilePath)

	// Was there a problem opening the file
	if err != nil {

		helper.WriteStdoutf(goRoutine, "search", "_LoadFeedsList", "ERROR : %s", err)
		return nil, err
	}

	defer func() {

		file.Close()
		helper.WriteStdout(goRoutine, "search", "_LoadFeedsList", "Closing File : Defer Completed")
	}()

	// Create a list to hold the feed urls
	feedList = list.New()

	// Open a reader to the feeds.list file
	reader := bufio.NewReader(file)

	// Read every line and store it
	for {

		// Read a line from the file
		line, err = reader.ReadString('\n')

		// Was there a problem reading the file
		if err != nil {

			if err == io.EOF {

				err = nil
				break
			}

			helper.WriteStdoutf(goRoutine, "search", "_LoadFeedsList", "ERROR : %s", err)
			return nil, err
		}

		uri := strings.TrimRight(line, "\n")

		helper.WriteStdoutf(goRoutine, "search", "_LoadFeedsList", "%s", uri)

		// Store the url
		feedList.PushBack(uri)
	}

	helper.WriteStdout(goRoutine, "search", "_LoadFeedsList", "Completed")

	return feedList, err
}

// _CreateSearchPlugin will create an object of the specified searchPlugin type
//  goRoutine: The Go routine making the call
//  pluginType: The type of search plugin to use
//  parameters: The parameters for the command
func _CreateSearchPlugin(goRoutine string, pluginType string, parameters []string) (searchPlugin SearchPlugin, err error) {

	// TODO: Add new search plugin types here

	switch pluginType {

	case "expression":
		return expression.New(goRoutine, parameters)
	}

	// Display the help examples
	DisplayHelpExamples(goRoutine)

	return nil, errors.New("Unknown Command")
}

//** PRIVATE MEMBER FUNCTIONS

// PerformSearch performs all business logic related to searching the feeds
//  goRoutine: The Go routine making the call
//  search: The searchPlugin to use when performing the search
func (this *_Search) PerformSearch(goRoutine string, searchPlugin SearchPlugin) {

	defer helper.CatchPanic(nil, goRoutine, "search", "PerformSearch")

	helper.WriteStdout(goRoutine, "search._Search", "PerformSearch", "Started")

	// Capture the number of feeds to process
	numberOfFeeds := this.FeedList.Len()

	// Create an array to hold the results
	searchResults := make([]interface{}, numberOfFeeds)

	// Channel used to wait for all work to be completed
	// The results are sent back on this channel
	wait := make(chan interface{}, numberOfFeeds)

	// Post search work for each feed in the list
	for element := this.FeedList.Front(); element != nil; element = element.Next() {

		// Prepare to run the search
		searchWork := &_SearchWork{
			SP:           this.WorkPool,
			Wait:         wait,
			Url:          element.Value.(string),
			SearchPlugin: searchPlugin,
		}

		// Post the search into the search pool
		this.WorkPool.PostWork(goRoutine, searchWork)
	}

	helper.WriteStdout(goRoutine, "search._Search", "PerformSearch", "Info : Waiting For Feeds To Complete")

	// Wait for each feed to signal they are done
	for feed := 0; feed < numberOfFeeds; feed++ {

		searchResults[feed] = <-wait
	}

	// Display the results
	searchPlugin.DisplayResults(goRoutine, searchResults)

	helper.WriteStdout(goRoutine, "search._Search", "PerformSearch", "Completed")
}

// PerformSearch is called to perform the work from the search pool
//  workRoutine: The internal id of the work routine making the call
func (this *_SearchWork) DoSearch(workRoutine int) {

	var results interface{}

	defer func() {

		// Signal that we are done
		this.Wait <- results
	}()

	// Create the name of the Go routine this is being processed in
	goRoutine := fmt.Sprintf("Work Routine %d", workRoutine)

	helper.WriteStdoutf(goRoutine, "search._Search", "DoSearch", "Started : QW: %d  AR: %d", this.SP.QueuedWork(), this.SP.ActiveRoutines())

	// Retrieve the RSS feed document
	rssDocument, err := rss.RetrieveRssFeed(goRoutine, this.Url)

	if err != nil {

		helper.WriteStdoutf(goRoutine, "search._Search", "DoSearch", "ERROR - Completed : %s", err)
		return
	}

	// Use the search plugin to find the results
	results = this.SearchPlugin.FindResults(goRoutine, rssDocument)

	helper.WriteStdoutf(goRoutine, "search._Search", "DoSearch", "Completed : QW: %d  AR: %d", this.SP.QueuedWork(), this.SP.ActiveRoutines())
}