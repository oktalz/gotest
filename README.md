# gotest
Go test tool to help run only tests with tags

## installation

go install github.com/oktalz/gotest@latest

## local installation

task install

## usage

gotest --tags tag1,tag2

## why

when you run `go test --tags=integration`
you will run both test that have `integration` and tests that does not have any build tags.
This is sometimes not desired

## todo

support all go test params
