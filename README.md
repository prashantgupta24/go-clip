# Go-clip

[![codecov](https://codecov.io/gh/prashantgupta24/go-clip/branch/master/graph/badge.svg?token=PSO715XHBI)](https://codecov.io/gh/prashantgupta24/go-clip) [![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/go-clip)](https://goreportcard.com/report/github.com/prashantgupta24/go-clip) [![version][version-badge]][releases] ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/prashantgupta24/go-clip)

[version-badge]: https://img.shields.io/github/v/release/prashantgupta24/go-clip
[releases]: https://github.com/prashantgupta24/go-clip/releases

A minimalistic clipboard manager for Mac.

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Demo](#demo)
- [Disclaimer](#disclaimer)
- [Basic functionality](#basic-functionality)
  - [Pin](#pin)
  - [Obfuscate](#obfuscate)
- [How to install](#how-to-install)
  - [Install from binary](#install-from-binary)
  - [Install from source](#install-from-source)
- [Future work](#future-work)

<!-- /code_chunk_output -->

## Demo

![](https://github.com/prashantgupta24/go-clip/blob/master/demo/go-clip-demo.gif)

## Disclaimer

1. This application is intended to help you manage multiple short-lived clippings, it is **not meant as a password store or a place to save important data for later use**.

1. The application **DOES NOT SAVE YOUR TEXT** anywhere on your computer. It is all in-memory, you lose all your data if you quit the application. So do not worry about your copied passwords being saved in plain-text anywhere.

1. The application has **NO connectivity to the internet**, so do not worry about your copied text being used for nefarious purposes. Feel free to go through the code to make sure of that.

## Basic functionality

### Pin

Pinning clipping will prevent them from being over-written once all existing clippings are used up.

### Obfuscate

In case you copy a password and want to keep it in the tool for sometime, you can obfuscate it. It will mark it with `****`. You can still copy the whole password once you click on it, it just won't show up on the application.

This is useful in case you don't want your password to be displayed as whole while you are sharing your screen with someone.

## How to install

### Install from binary

1. Download the latest `go-clip.app.zip` from the [releases](https://github.com/prashantgupta24/go-clip/releases) page, unzip it, and copy the .app to your `Applications` folder like any other application.

1. Since the application is not notarized, you will need to right click on the `go-clip.app` and choose `Open`.

1. You will see a scary message that warns you about all the bad things that the app can do to your computer. If you are paranoid (fair enough, you don't really know me that well) then you can skip to the section which builds the app from source. That way you can see what exactly the app does (You can check that the application makes no connections to the internet whatsoever).

1. In case you do trust me, once you click on `Open`, you should see the clipboard icon on your system tray. Copy a piece of text to see it on the app!

### Install from source

Make sure you have `go` installed. Once that is done, clone this repo and run `Make`, it should create the `go-clip.app` and open the folder where it was built for you. Copy the .app to your `Applications` folder like any other application.

Double click on the app, and the clipboard icon should start displaying on your system tray. Copy a piece of text to see it on the app!

## Future work

Thanks to `github.com/getlantern/systray`, I am able to provide a system tray functionality. I am waiting for `fyne.io` to release their integration with systray, so that I can offer a pop-out clipboard experience along with systray. Right now both systray and external app can't exist simultaneously, since both need to be executed on the main thread. (Ref https://github.com/fyne-io/fyne/issues/283)
