# go-tun2socks-mobile

[![Build Status](https://travis-ci.com/eycorsican/go-tun2socks-mobile.svg?branch=master)](https://travis-ci.com/eycorsican/go-tun2socks-mobile)

Demo for building and using `go-tun2socks` on iOS and Android.

> If you're looking for an easy to use `tun2socks` implementation for iOS, you might be interested in [`leaf`](https://github.com/eycorsican/leaf) and [`ileaf`](https://github.com/eycorsican/ileaf).
> `leaf` [implements `tun2socks`](https://github.com/eycorsican/leaf/tree/master/leaf/src/proxy/tun/netstack) and it's written in Rust, with significantly less memory usage and significantly better performance compares to the Go version.

## Prerequisites

- macOS (iOS)
- Xcode (iOS)
- SDK (Android)
- NDK (Android)
- make
- Go >= 1.11
- A C compiler (e.g.: clang, gcc)
- gomobile (https://github.com/golang/go/wiki/Mobile)
- Other common utilities (e.g.: git)

## Build
```bash
go get -d ./...

# Build an AAR
make android

# Build a Framework
make ios

# Both
make
```
