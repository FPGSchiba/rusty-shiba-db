package main

import (
	"rsdb/src/router"
	"rsdb/src/rust/collections"
)

func main() {
	// TODO: Config stuff
	collections.InitRustyStorage()
	engine := router.GetRouter()
	err := engine.Run(":3000")
	if err != nil {
		return
	}
}
