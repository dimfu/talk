package main

type ID int

const (
	MSG ID = iota
)

type command struct {
	id     ID
	sender string
	body   []byte
	client *client
}
