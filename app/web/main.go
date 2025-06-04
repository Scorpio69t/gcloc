package main

import "github.com/Scorpio69t/gcloc/web"

func main() {
	if err := web.Start(":8080"); err != nil {
		panic(err)
	}
}
