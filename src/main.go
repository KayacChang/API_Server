package main

import (
	"api/env"
	"api/framework/cache"
	"api/framework/postgres"
	"api/user"
	"sync"
)

func main() {
	e := env.New()

	c := cache.Get()
	db := postgres.New(e.Postgres.ToURL(), 30)

	var wg sync.WaitGroup

	wg.Add(3)

	// go func() {
	// 	game.New(e)

	// 	wg.Done()
	// }()

	go func() {
		user.New(e, db, c)

		wg.Done()
	}()

	// go func() {
	// 	order.New(e)

	// 	wg.Done()
	// }()

	wg.Wait()
}
