package postgres

import (
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresClient struct {
	C_RW *sqlx.DB
	C_RO *sqlx.DB
}

func NewPostgresClient() *PostgresClient {
	client := &PostgresClient{
		C_RO: sqlx.MustConnect("postgres", os.Getenv("POSTGRES_URI_RO")+" application_name=backoffice_service"),
		C_RW: sqlx.MustConnect("postgres", os.Getenv("POSTGRES_URI_RW")+" application_name=backoffice_service"),
	}

	client.C_RO.SetMaxOpenConns(4)
	client.C_RW.SetMaxOpenConns(4)
	client.C_RO.SetMaxIdleConns(4)
	client.C_RW.SetMaxIdleConns(4)
	client.C_RO.SetConnMaxLifetime(60 * time.Second)
	client.C_RW.SetConnMaxLifetime(60 * time.Second)

	log.Println("PostgresDB connected")
	return client
}

func (c *PostgresClient) Ping() {
	if err := c.C_RO.Ping(); err != nil {
		log.Fatal(err)
	}
}

func (c *PostgresClient) Close() {
	if err := c.C_RO.Close(); err != nil {
		log.Fatal(err)
	}
	if err := c.C_RW.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("PostgresDB disconnected")
}
