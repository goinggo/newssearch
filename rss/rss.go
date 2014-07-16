// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rss provide support for retrieving and reading rss feeds.
package rss

import (
	"bytes"
	"encoding/xml"
	"errors"
	"github.com/goinggo/newssearch/helper"
	"io/ioutil"
	"net/http"
)

// Item defines the fields associated with the item tag in the buoy RSS document.
type Item struct {
	XMLName     xml.Name `xml:"item"`
	PubDate     string   `xml:"pubDate"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	GeoRssPoint string   `xml:"georss:point"`
}

// Image defines the fields associated with the image tag in the buoy RSS document.
type Image struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
}

// Channel defines the fields associated with the channel tag in the buoy RSS document.
type Channel struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`
	Description    string   `xml:"description"`
	Link           string   `xml:"link"`
	PubDate        string   `xml:"pubDate"`
	LastBuildDate  string   `xml:"lastBuildDate"`
	TTL            string   `xml:"ttl"`
	Language       string   `xml:"language"`
	ManagingEditor string   `xml:"managingEditor"`
	WebMaster      string   `xml:"webMaster"`
	Image          Image    `xml:"image"`
	Item           []Item   `xml:"item"`
}

// Document defines the fields associated with the buoy RSS document.
type Document struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
	URI     string
}

// RetrieveRssFeed performs a HTTP Get request to the RSS feed and serializes the results.
func RetrieveRssFeed(goRoutine string, uri string) (*Document, error) {
	helper.WriteStdoutf(goRoutine, "rss", "RetrieveRssFeed", "Started : Uri[%s]", uri)

	if uri == "" {
		helper.WriteStdout(goRoutine, "rss", "RetrieveRssFeed", "Completed : No RSS Feed Uri Provided")
		return nil, errors.New("No RSS Feed Uri Provided")
	}

	resp, err := http.Get(uri)
	if err != nil {
		helper.WriteStdoutf(goRoutine, "rss", "RetrieveRssFeed", "ERROR - Completed : HTTP Get : %s : %s", uri, err)
		return nil, err
	}

	defer resp.Body.Close()

	rawDocument, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helper.WriteStdoutf(goRoutine, "rss", "RetrieveRssFeed", "ERROR - Completed : Read Resp : %s : %s", uri, err)
		return nil, err
	}

	var document Document
	if err = xml.NewDecoder(bytes.NewReader(rawDocument)).Decode(&document); err != nil {
		helper.WriteStdoutf(goRoutine, "rss", "RetrieveRssFeed", "ERROR - Completed : Decode : %s : %s", uri, err)
		return nil, err
	}

	// Save the uri to the feed.
	document.URI = uri

	helper.WriteStdoutf(goRoutine, "rss", "RetrieveRssFeed", "Completed : Uri[%s] Title[%s]", uri, document.Channel.Title)
	return &document, err
}
