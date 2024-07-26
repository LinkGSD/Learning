package main

type name struct {
	a func() string
}

func main() {
	n := &name{}
	n.a()
}
