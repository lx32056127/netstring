# netstring
Simple netstring server and client in go. Netstring is a simple protocol. See [Wikipedia](https://en.wikipedia.org/wiki/Netstring) to more info.

## Installation

    go get  https://github.com/elektro79/netstring

## How to use

See test for complete use. You must implement NetStringProcessor. You also must change netstring.MaxMsg to avoid a malicius client eat all ram