# bpks

[![Build Status](https://travis-ci.org/tomdionysus/bpks.svg?branch=master)](https://travis-ci.org/tomdionysus/bpks)
[![Go Report Card](https://goreportcard.com/badge/github.com/tomdionysus/bpks)](https://goreportcard.com/report/github.com/tomdionysus/bpks)
[![GoDoc](https://godoc.org/github.com/tomdionysus/bpks?status.svg)](https://godoc.org/github.com/tomdionysus/bpks)

A B+Tree Key Store in golang.

bpks is a B+Tree that stores arbitary key/value pairs on any device or file that supports io.ReadWriteSeeker.

* String Keys
* Arbitary []byte data