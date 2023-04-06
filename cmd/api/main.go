package main

import (
	_ "github.com/lib/pq"

	"github.com/JohanVong/online_bazaar/internal/app"
)

func main() {
	app.AssembleAndGo()
}
