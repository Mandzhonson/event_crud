package main

import "calendar/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
