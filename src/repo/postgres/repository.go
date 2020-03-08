package postgres

import (
	"context"
	"database/sql"
	"time"
	"user/model"
	"user/repo"
	"user/utils"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	errs "github.com/pkg/errors"
)

type repository struct {
	db      *sqlx.DB
	timeout time.Duration
	sql     querys
}

type querys struct {
	insert   string
	findByID string
}

func connect(url string, timeout time.Duration) (*sqlx.DB, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return sqlx.ConnectContext(ctx, "pgx", url)
}

func New(url string, timeout int) repo.Repository {

	_timeout := time.Duration(timeout) * time.Second
	db, err := connect(url, _timeout)

	if err != nil {
		panic(err)
	}

	return &repository{
		db:      db,
		timeout: _timeout,

		sql: querys{
			insert:   utils.ParseFile("sql/insert_one.sql"),
			findByID: utils.ParseFile("sql/find_by_id.sql"),
		},
	}
}

func (it *repository) findByID(ctx context.Context, user *model.User) error {

	return it.db.GetContext(ctx, user, it.sql.findByID, user.ID)
}

func (it *repository) FindBy(key string, user *model.User) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), it.timeout)
	defer cancel()

	switch key {
	case "ID":
		err = it.findByID(ctx, user)
	}

	if err == sql.ErrNoRows {
		return model.ErrUserNotFound
	}

	return err
}

func (it *repository) Store(user *model.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), it.timeout)
	defer cancel()

	opt := &sql.TxOptions{Isolation: sql.LevelSerializable}
	tx, err := it.db.BeginTxx(ctx, opt)
	if err != nil {

		return err
	}

	// === Insert ===
	_, err = tx.NamedExec(it.sql.insert, user)
	if err != nil {
		tx.Rollback()

		if isErrIntegrityConstraint(err) {

			return errs.WithMessage(model.ErrDBConstraint, err.Error())
		}

		return err
	}

	// === Get Inserted Data ===
	err = tx.Get(user, it.sql.findByID, user.ID)
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (it *repository) Delete(user *model.User) error {
	// TODO

	return nil
}

func isErrIntegrityConstraint(err error) bool {

	pgerr, ok := err.(pgx.PgError)

	if !ok {
		return false
	}

	errlist := [...]string{
		pgerrcode.IntegrityConstraintViolation,
		pgerrcode.RestrictViolation,
		pgerrcode.NotNullViolation,
		pgerrcode.ForeignKeyViolation,
		pgerrcode.UniqueViolation,
		pgerrcode.CheckViolation,
		pgerrcode.ExclusionViolation,
	}

	for _, code := range errlist {

		if code == pgerr.Code {
			return true
		}
	}

	return false
}

func isErr(err error, code string) bool {

	pgerr, ok := err.(pgx.PgError)

	return ok && pgerr.Code == code
}
