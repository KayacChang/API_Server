package game

import (
	"api/model"
	"fmt"
)

// FindByID find game by id in db
func (it Repo) FindByID(id string) (*model.Game, error) {

	game := model.Game{}

	sql := fmt.Sprintf("SELECT * FROM %s WHERE game_id=$1", table)

	if err := it.db.Ping(); err != nil {

		return nil, err
	}

	if err := it.db.Get(&game, sql, id); err != nil {

		return nil, err
	}

	return &game, nil
}

// FindByName find game by name in db
func (it Repo) FindByName(name string) (*model.Game, error) {

	game := model.Game{}

	sql := fmt.Sprintf("SELECT * FROM %s WHERE name=$1", table)

	if err := it.db.Ping(); err != nil {

		return nil, err
	}

	if err := it.db.Get(&game, sql, name); err != nil {

		return nil, err
	}

	return &game, nil
}

// FindAll find all games in db
func (it Repo) FindAll() ([]model.Game, error) {

	games := []model.Game{}

	sql := fmt.Sprintf("SELECT * FROM %s", table)

	if err := it.db.Ping(); err != nil {

		return nil, err
	}

	if err := it.db.Select(&games, sql); err != nil {

		return nil, err
	}

	return games, nil
}
