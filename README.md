# Go-clip

[![codecov](https://codecov.io/gh/prashantgupta24/go-clip/branch/master/graph/badge.svg?token=PSO715XHBI)](https://codecov.io/gh/prashantgupta24/go-clip) [![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/go-clip)](https://goreportcard.com/report/github.com/prashantgupta24/go-clip) [![version][version-badge]][releases] ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/prashantgupta24/go-clip)

[version-badge]: https://img.shields.io/github/v/release/prashantgupta24/go-clip
[releases]: https://github.com/prashantgupta24/go-clip/releases

A minimalistic clipboard manager for Mac.

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [How to install](#how-to-install)
  - [Install from binary](#install-from-binary)
  - [Install from source](#install-from-source)

<!-- /code_chunk_output -->

## How to install

### Install from binary

1. Download the latest `go-clip.app.zip` from the [releases](https://github.com/prashantgupta24/go-clip/releases) page, unzip it, and copy the .app to your `Applications` folder like any other application.

1. Since the application is not notarized, you will need to right click on the `go-clip.app` and choose `Open`.

1. You will see a scary message that warns you about all the bad things that the app can do to your computer. If you are paranoid (fair enough, you don't really know me that well) then you can skip to the section which builds the app from source. That way you can see what exactly the app does (You can check that the application makes no connections to the internet whatsoever).

1. In case you do trust me, once you click on `Open`, you should see the clipboard icon on your system tray. Copy a piece of text to see it on the app!

### Install from source

Make sure you have `go` installed. Once that is done, clone this repo and run `Make`, it should create the `go-clip.app` and open the folder where it was built for you. Copy the .app to your `Applications` folder like any other application.

Double click on the app, and the clipboard icon should start displaying on your system tray. Copy a piece of text to see it on the app!
