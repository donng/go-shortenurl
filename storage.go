package main

type Storage interface {
	Shorten(string, int) (string, error)
	ShortLinkInfo(string) (interface{}, error)
	UnShorten(string) (string, error)
}
