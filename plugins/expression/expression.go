// Copyright 2013 Ardan Studios. All rights reserved.
// Use of expression source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Search Plugins must be thread safe. Only one instance of a Search Plugin will be instantiated.
The FindResults interface method can be called simultaneously by many routines.
The DisplayResults interface will be called only once.
*/

// Package expression provides an expression based plugin for searching news.
package expression

import (
	"container/list"
	"errors"
	"github.com/goinggo/newssearch/helper"
	"github.com/goinggo/newssearch/rss"
	"regexp"
)

// Expression performs regular expression matching.
type Expression struct {
	parameters []string
}

// searchMatch contains the result of a successful match.
type searchMatch struct {
	Expression string   // The expression that was matched.
	Field      string   // The field that was matched in the rss document.
	Document   rss.Item // The item document that was matched.
}

// searchResult contains the result of a search.
type searchResult struct {
	URI     string     // The feed URI that was processed.
	Matches *list.List // A list of the searchMatch objects.
}

// New creates a new expression search plugin
func New(goRoutine string, parameters []string) (*Expression, error) {
	// Check the parameters are valid
	result := checkParameters(parameters)
	if result == false {
		return nil, errors.New("Invalid _Parameters")
	}

	// Create an expression search plugin
	expression := Expression{
		parameters: parameters,
	}

	return &expression, nil
}

// HelpParameters provides command useage help.
func HelpParameters() (message string) {
	return "newssearch 8 expression (?i)spy"
}

// CheckParameters reviews the parameters to make sure they are compatible.
func checkParameters(parameters []string) (result bool) {
	// We accept any parameters
	return true
}

// FindResults implements the SearchPlugin interface for processing expression searches.
// https://code.google.com/p/re2/wiki/Syntax
// https://code.google.com/p/re2/source/browse/re2/re2.h
func (expression *Expression) FindResults(goRoutine string, document *rss.Document) interface{} {
	defer helper.CatchPanic(nil, goRoutine, "expression.Expression", "FindResults")

	// Create a search result for expression run
	searchResult := searchResult{
		URI:     document.URI,
		Matches: list.New(),
	}

	for _, rssItem := range document.Channel.Item {
		for _, expression := range expression.parameters {
			// Check the title
			matched, err := regexp.MatchString(expression, rssItem.Title)
			if err != nil {
				helper.WriteStdoutf(goRoutine, "expression.Expression", "FindResults", "ERROR : Title Search : URI[%s] Expression[%s] Title[%s] : %s", document.URI, expression, rssItem.Title, err)
				return searchResult
			}

			if matched == true {
				// Add the match
				searchResult.Matches.PushBack(&searchMatch{
					Expression: expression,
					Field:      "Title",
					Document:   rssItem,
				})

				helper.WriteStdoutf(goRoutine, "expression.Expression", "FindResults", "MATCH : Title Search : URI[%s] Expression[%s] Title[%s]", document.URI, expression, rssItem.Title)
				continue
			}

			// Check the description
			matched, err = regexp.MatchString(expression, rssItem.Description)
			if err != nil {
				helper.WriteStdoutf(goRoutine, "expression.Expression", "FindResults", "ERROR : Description Search : URI[%s] Expression[%s] Desc[%s] : %s", document.URI, expression, rssItem.Description, err)
				return searchResult
			}

			if matched == true {
				// Add the match
				searchResult.Matches.PushBack(&searchMatch{
					Expression: expression,
					Field:      "Description",
					Document:   rssItem,
				})

				helper.WriteStdoutf(goRoutine, "expression.Expression", "FindResults", "MATCH : Description Search : URI[%s] Expression[%s] Title[%s]", document.URI, expression, rssItem.Title)
				continue
			}
		}
	}

	if searchResult.Matches.Len() == 0 {
		helper.WriteStdoutf(goRoutine, "expression.Expression", "FindResults", "NOT MATCHED : URI[%s]", document.URI)
	}

	return &searchResult
}

// DisplayResults implements the SearchPlugin interface for displaying search results.
func (expression *Expression) DisplayResults(goRoutine string, results interface{}) {
	helper.WriteStdout(goRoutine, "expression.Expression", "DisplayResults", "Show Results\n")

	// Iterate through the array of search results
	for _, result := range results.([]interface{}) {
		// Verify we received a result from expression feed
		if result == nil {
			continue
		}

		// Cast the search result
		searchResult := result.(*searchResult)

		for element := searchResult.Matches.Front(); element != nil; element = element.Next() {
			match := element.Value.(*searchMatch)

			helper.WriteStdoutf(goRoutine, "expression.Expression", "DisplayResults", "MATCH : *************************************")
			helper.WriteStdoutf(goRoutine, "expression.Expression", "DisplayResults", "MATCH : URI[%s] Expression[%s] Field[%s]\n%s\n%s", searchResult.URI, match.Expression, match.Field, match.Document.Title, match.Document.Description)
			helper.WriteStdoutf(goRoutine, "expression.Expression", "DisplayResults", "MATCH : *************************************\n\n")
		}
	}
}
