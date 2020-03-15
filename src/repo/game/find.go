package game

import (
	"api/model"
	"api/utils"

	"github.com/jmoiron/sqlx"
)

func (it *Repo) findCache(ids ...string) ([]string, []*model.Game) {

	games := []*model.Game{}

	finds := []string{}

	for _, id := range ids {

		if game, found := it.cache.Get(id); found {

			if game, ok := game.(model.Game); ok {

				games = append(games, &game)

				finds = append(finds, id)
			}
		}
	}

	return utils.Diff(ids, finds), games
}

func (it *Repo) FindByID(id string) (*model.Game, error) {

	game := model.Game{}

	if _, games := it.findCache(id); len(games) == 1 {
		return games[0], nil
	}

	if err := it.db.Get(&game, it.sql.findByID, id); err != nil {
		return nil, err
	}

	// === Save to Cache ===
	defer it.storeCache(&game)

	return &game, nil
}

func (it *Repo) FindByIDs(ids []string) ([]*model.Game, error) {

	remains, games := it.findCache(ids...)

	if len(remains) == 0 {
		return games, nil
	}

	query, args, err := sqlx.In(it.sql.findByID, remains)
	if err != nil {
		return nil, err
	}

	err = it.db.Select(&games, it.db.Rebind(query), args...)

	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return nil, model.ErrGameNotFound
	}

	// === Save to Cache ===
	defer it.storeCache(games...)

	return games, nil
}