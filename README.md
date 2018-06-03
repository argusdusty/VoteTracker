# VoteTracker

VoteTracker implements a simple and efficient webserver for tracking election results and combining results together from multiple sources (e.g. a state's official election results page, Clarity, etc.). No sources or results are included in this repo except for example data. The full code is hosted at https://elections.argusdusty.com with live election results every election night and past results for hundreds of elections

## Features

* https, automatic certificate generation with Let's Encrypt
* API json data by appending .json to page requests
* Access-Control-Allow-Origin: * for json requests to allow people to use the live data on their own websites
* Automating page refreshing with http-equiv="refresh"
* Page/file caching to improve server performance
* Intelligent source caching with If-Modified-Since to improve source loading performance
* Intelligent results merging that combines county data together to get results that are more up-to-date than any individual source
* Seperates the server code from the updater code so that the updater can be modified without having to restart the server
* Multi-updaters so that all the updaters for a single election night can be run with a single program
* Forecasts that can be pre-programmed with polling data and live updated with election results (more improvements planned)
* Election maps

## Install

Pull and install the code

```sh
go get -u github.com/argusdusty/VoteTracker
go install VoteTracker
```

## Run the Server

Simply run the compiled VoteTracker program

```sh
VoteTracker
```

## Track a Race

First, create a source (or multiple sources) under the Sources directory. You'll need a url that has the results, and a function that parses it to a Go-usable struct. Give your source some params that identify what race/url you're loading, so you can re-use that source in future elections. Then use `LoadURL` to load and cache the URL (subsequent requests will send `If-Modified-Since` header and use the cached result if Not Modified status is returned).

```Go
data, err := LoadURL(url, parseSourceToStructFunction)
```

Next, turn that struct into a `Summary` just like `ExampleSource` in Sources/example.go does. Regions are generally keyed on their FIPS county code so you can use them in maps with topo.json, however you can use any keys you want as long as the topo.json file matches it (or you aren't using a topo.json file). You can get a state's topo.json file from: https://github.com/deldersveld/topojson

Finally, create your race(s) just like in ExampleDate, in a date/race/source tree structure which supports multiple races per date and multiple sources per race. Just setup the updater.go files under the date/race/source, the date/races.json file, and index.json in the root folder. Then just the updater:

```sh
cd date
go run updater.go
```

And you're now tracking an election.