package users

import (
	"context"
	"database/sql"
	"time"

	"practice3go/internal/repository/_postgres"
	"practice3go/pkg/modules"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: 5 * time.Second,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var users []modules.User
	err := r.db.DB.SelectContext(ctx, &users, `
		select id, name, email, age, created_at
		from users
		order by id
	`)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var u modules.User
	err := r.db.DB.GetContext(ctx, &u, `
		select id, name, email, age, created_at
		from users
		where id = $1
	`, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, modules.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) CreateUser(u modules.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var id int
	err := r.db.DB.QueryRowxContext(ctx, `
		insert into users (name, email, age)
		values ($1, $2, $3)
		returning id
	`, u.Name, u.Email, u.Age).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateUser(id int, u modules.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(ctx, `
		update users
		set name = $1, email = $2, age = $3
		where id = $4
	`, u.Name, u.Email, u.Age, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return modules.ErrUserNotFound
	}

	return nil
}

func (r *Repository) DeleteUserByID(id int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(ctx, `
		delete from users where id = $1
	`, id)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected == 0 {
		return 0, modules.ErrUserNotFound
	}

	return affected, nil
}
