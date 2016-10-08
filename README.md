# bpks

[![Build Status](https://travis-ci.org/tomdionysus/bpks.svg?branch=master)](https://travis-ci.org/tomdionysus/bpks)
[![Coverage Status](https://coveralls.io/repos/tomdionysus/bpks/badge.svg?branch=master&service=github)](https://coveralls.io/github/tomdionysus/bpks?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/tomdionysus/bpks)](https://goreportcard.com/report/github.com/tomdionysus/bpks)
[![GoDoc](https://godoc.org/github.com/tomdionysus/bpks?status.svg)](https://godoc.org/github.com/tomdionysus/bpks)

A B+Tree Key Store in golang.

bpks is a B+Tree that stores arbitary key/value pairs on any device or file that supports io.ReadWriteSeeker.

* 128-bit or String Keys
* Arbitary []byte data

## Status / TODO

bpks is currently ALPHA and should not be used in production. The major missing component is a real free space manager (allocator/deallocator) for blocks.

Done:

* Read/write/remove key/value

Needs:

* Multi-block data
* A real bitmap/range based allocator
* Index Block merging on Remove key
* Block caching

## License

bpks is licensed under the Open Source MIT license. Please see the [License File](LICENSE.txt) for more details.

## Code Of Conduct

The bkps project supports and enforces [The Contributor Covenant](http://contributor-covenant.org/) code of conduct. Please read [the code of conduct](CODE_OF_CONDUCT.md) before contributing.