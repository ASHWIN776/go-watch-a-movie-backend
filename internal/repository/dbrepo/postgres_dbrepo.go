package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"time"
)

// This type will have methods that interact with the database - in order to do that, they need a DB connection
type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

// To return the connection
func (repo *PostgresDBRepo) Connection() *sql.DB {
	return repo.DB
}

// To get all the movies from the database
func (repo *PostgresDBRepo) AllMovies() ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	stmt := `
		select 
			id, title, release_date, runtime, mpaa_rating,
			description, coalesce(image, ''), 
			created_at, updated_at
		from
			movies
	`

	rows, err := repo.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.Runtime,
			&movie.MpaaRating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

// To get a user using email
func (repo *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	stmt := `
		select 
			id, email, first_name, last_name, 
			password, created_at, updated_at
		from users 
		where
			email = $1
	`

	row := repo.DB.QueryRowContext(ctx, stmt, email)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

// To get a user by id
func (repo *PostgresDBRepo) GetUserByID(userId int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	stmt := `
		select 
			id, email, first_name, last_name, 
			password, created_at, updated_at
		from users 
		where
			id = $1
	`

	row := repo.DB.QueryRowContext(ctx, stmt, userId)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
