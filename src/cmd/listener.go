package main

// Listener is an interface represeting the module
// responsible for receiving requests.
type Listener interface {
	Listen(string) error
}
