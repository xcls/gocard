# Gocard

Web-app for sharing and using spaced repetition flash cards.

## Development

```
go get -u github.com/xcls/gocard
cd $GOPATH/src/github.com/xcls/gocard
make server
make autotest
```

## Features / Roadmap

* Easy to share decks/cards
* Decks can be added partially, and the parts can be activated/deactivated
  whenever
* Decks can have a description, similar to a README on Github. They can have
  external links to the subject/book/course they are about.
* An API that allows other apps to integrate the decks in, for example, video
  courses

## Inspiration

* Anki: http://ankisrs.net/
* Spaced Repetition and Learning: http://www.gwern.net/Spaced%20repetition

## Side-project motivation

* https://jvns.ca/blog/2016/09/19/getting-things-done/
* https://jvns.ca/blog/2016/08/16/how-do-you-work-on-something-important/
