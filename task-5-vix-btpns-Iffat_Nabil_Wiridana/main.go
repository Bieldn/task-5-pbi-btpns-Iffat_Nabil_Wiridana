package main

import (
	"vix-btpns/database"
	"vix-btpns/helpers"
	"vix-btpns/router"
)

func main() {
	helpers.LoadEnv()

	database.InitDB()

	var rt router.Route
	rt.Init()
}
