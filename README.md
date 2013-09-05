# News Search

Copyright 2013 Ardan Studios. All rights reserved.  
Use of this source code is governed by a BSD-style license that can be found in the LICENSE handle.

This application processes a set of RSS feed Url and looks for keywords and regular expressions. Any matches that are found is returned to the user. The application provides frameworks for creating plugable code, channel communication, and the use of ArdanStudios/workpool.

Ardan Studios  
12973 SW 112 ST, Suite 153  
Miami, FL 33186  
bill@ardanstudios.com

GoingGo.net Post:  
http://www.goinggo.net/2013/07/an-rss-feed-searching-framework-using-go.html

	-- Get, build and install the code
	export GOPATH=$HOME/goinggo
	go get github.com/goinggo/newssearch
	
	-- Run the code
	cd $GOPATH/bin
	** Copy the feeds.list file to bin
	./newssearch