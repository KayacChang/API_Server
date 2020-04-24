package token

import (
	"api/env"
	"api/framework/server"

	token "api/usecase/token"
)

type Handler struct {
	server.Server
	env     env.Env
	usecase token.Usecase
}

func New(server server.Server, env env.Env, usecase token.Usecase) Handler {

	return Handler{server, env, usecase}
}

func (it Handler) getHref(url string) string {

	return it.env.Service.Domain + "/" + it.env.API.Version + url
}
