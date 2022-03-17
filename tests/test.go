package main

import "fmt"

func main() {
	m := make(map[string]bool)
	m["hah"] = true
	fmt.Println(m["hah"])
	fmt.Println(m["wrong"])
}
