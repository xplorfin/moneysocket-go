[![Coverage Status](https://coveralls.io/repos/github/xplorfin/moneysocket-go/badge.svg?branch=master)](https://coveralls.io/github/xplorfin/moneysocket-go?branch=master)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://app.renovatebot.com/dashboard#github/xplorfin/moneysocket-go)
[![Build status](https://github.com/xplorfin/moneysocket-go/workflows/test/badge.svg)](https://github.com/xplorfin/moneysocket-go/actions?query=workflow%3Atest)
[![Build status](https://github.com/xplorfin/moneysocket-go/workflows/goreleaser/badge.svg)](https://github.com/xplorfin/moneysocket-go/actions?query=workflow%3Agoreleaser)
[![](https://godoc.org/github.com/xplorfin/moneysocket-go?status.svg)](https://godoc.org/github.com/xplorfin/moneysocket-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/xplorfin/moneysocket-go)](https://goreportcard.com/report/github.com/xplorfin/moneysocket-go)

# Status

This is a golang implementation of the [moneysocket protocol](https://github.com/moneysocket/py-moneysocket) that aims to have full integration testing & parity with py & js moneysocket. This project is not yet usable in any state beyond tvl encoding/decoding. The code is a bit of a mess as well since parts of this were torn from private [xplorfin](https://entropy.rocks/) repos. If you want to check it out anyway, you can run terminus

 To Do:
 
 1. General cleanup
 2. Cleanup mixed tlv encoding/decoding
 3. use uuid or string consistently

# Architecture
   
Testing:
 1. create two accounts with an arbitrary number of sats:
    - `terminus-cli create`: (returns account-0)
    - `terminus-cli create`: (returns account-1)
 1. `./terminus-cli listen account-1` -> copy and paste account-1 into the seller wallet consumer
 1. `./terminus-cli listen account-0` -> copy and paste into wallet consumer
 1. Wallet -> provider -> generate beacon -> connect -> copy beacon
 1. paste into my wallet consumer
 1. generate buyer consumer beacon
 1. connect seller app provider to buyer consumer
 1. seller opens store
 
 
