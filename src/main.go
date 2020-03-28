package main

import (
	"api/env"
	"api/framework/cache"
	"api/framework/postgres"
	"api/framework/server"
	"api/service/game"
	"api/service/order"
	"api/service/user"

	"github.com/go-chi/chi"
)

func main() {

	// === Framework ===
	env := env.New()
	cache := cache.Get()
	db := postgres.New(env.Postgres.ToURL())
	it := server.New(env)

	// === Handler ===
	game := game.New(it, env, db, cache)
	user := user.New(it, env, db, cache)
	order := order.New(it, env, db, cache)

	it.Route("/"+env.API.Version, func(server chi.Router) {
		// === Game ===
		server.Route("/games", func(server chi.Router) {
			server.Get("/", game.GET_ALL)
			server.Get("/{name}", game.GET)
			server.Post("/", game.POST)
			server.Put("/{name}", game.PUT)
		})

		// === User ===
		server.Post("/token", user.POST)
		server.Get("/auth", user.Auth)

		// === Order ===
		server.Route("/orders", func(server chi.Router) {
			server.Post("/", order.POST)
			server.Put("/{order_id}", order.PUT)
		})
	})

	it.Listen(env.Service.Port)
}
