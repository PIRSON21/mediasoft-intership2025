package postgresql

import (
	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	conn *pgx.Conn
}

func NewPostgres() (*Postgres, error) {
	// connOpts := createPostgresOpts()

	// conn, err := pgx.ConnectConfig(context.TODO(), connOpts)
	// if err != nil {
	// 	return nil, err
	// }

	// return &Postgres{
	// 	conn: conn,
	// }, nil

	return &Postgres{}, nil
}

func createPostgresOpts() *pgx.ConnConfig {
	cfg, _ := pgx.ParseConfig("")
	return cfg
}
