package order

import (
	"api/env"
	"api/framework/cache"
	"api/framework/postgres"
	"api/model"
	gamerepo "api/repo/game"
	orderrepo "api/repo/order"
	userrepo "api/repo/user"
	"api/utils"
	"api/utils/json"
	"database/sql"
	"fmt"

	"time"
)

type Usecase struct {
	env   *env.Env
	order *orderrepo.Repo
	user  *userrepo.Repo
	game  *gamerepo.Repo
}

func New(env *env.Env, db *postgres.DB, c *cache.Cache) *Usecase {

	return &Usecase{
		env:   env,
		order: orderrepo.New(db, c),
		user:  userrepo.New(db, c),
		game:  gamerepo.New(db, c),
	}
}

func (it *Usecase) sendBet(order *model.Order) error {

	user := model.User{
		ID: order.UserID,
	}
	if err := it.user.FindBy("ID", &user); err != nil {
		return err
	}

	game, err := it.game.FindByID(order.GameID)
	if err != nil {
		return err
	}

	url := it.env.Agent.Domain + it.env.Agent.API + "/transaction/game/bet"

	req := map[string]interface{}{
		"account":    user.Username,
		"created_at": order.CreatedAt.Time,
		"gamename":   game.Name,
		"roundid":    order.ID,
		"amount":     order.Bet,
	}

	headers := map[string]string{
		"Content-Type":       "application/json",
		"organization_token": it.env.Agent.Token,
		"session":            user.Session,
	}

	resp, err := utils.Post(url, req, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	res := map[string]interface{}{}
	json.Parse(resp.Body, &res)

	if resp.StatusCode != 200 {
		status := res["status"].(map[string]interface{})

		return fmt.Errorf("%s", status["message"])
	}

	data := res["data"].(map[string]interface{})
	balance := data["balance"].(float64)

	user.Balance = uint64(balance)

	if err := it.user.Store("Cache", &user); err != nil {
		return err
	}

	return nil
}

func (it *Usecase) Create(order *model.Order) error {

	order.ID = utils.UUID()
	order.State = model.Pending
	order.CreatedAt = sql.NullTime{time.Now(), true}

	if err := it.sendBet(order); err != nil {
		return err
	}

	return it.order.Store("Cache", order)
}

func (it *Usecase) sendEndRound(order *model.Order) error {

	user := model.User{
		ID: order.UserID,
	}
	if err := it.user.FindBy("ID", &user); err != nil {
		return err
	}

	game, err := it.game.FindByID(order.GameID)
	if err != nil {
		return err
	}

	url := it.env.Agent.Domain + it.env.Agent.API + "/transaction/game/endround"

	req := map[string]interface{}{
		"account":      user.Username,
		"created_at":   time.Now(),
		"gamename":     game.Name,
		"roundid":      order.ID,
		"completed_at": order.CompletedAt,
	}

	headers := map[string]string{
		"Content-Type":       "application/json",
		"organization_token": it.env.Agent.Token,
		"session":            user.Session,
	}

	resp, err := utils.Post(url, req, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	res := map[string]interface{}{}
	json.Parse(resp.Body, &res)

	if resp.StatusCode != 200 {
		status := res["status"].(map[string]interface{})

		return fmt.Errorf("%s", status["message"])
	}

	data := res["data"].(map[string]interface{})
	balance := data["balance"].(float64)

	user.Balance = uint64(balance)

	if err := it.user.Store("Cache", &user); err != nil {
		return err
	}

	return nil
}

func (it *Usecase) Checkout(order *model.Order) error {

}