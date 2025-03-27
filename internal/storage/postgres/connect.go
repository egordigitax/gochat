package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type PostgresClient struct {
	C_RW *sqlx.DB
	C_RO *sqlx.DB
}

func NewPostgresClient() *PostgresClient {

	getDSN := func() string {
		return fmt.Sprintf(
			"host=%s "+
				"user=%s "+
				"password=%s "+
				"dbname=%s "+
				"port=%d "+
				"sslmode=%s "+
				"TimeZone=%s",
			viper.GetString("database.host"),
			viper.GetString("database.user"),
			viper.GetString("database.password"),
			viper.GetString("database.dbname"),
			viper.GetInt("database.port"),
			viper.GetString("database.sslmode"),
			viper.GetString("database.timezone"),
		)
	}

	client := &PostgresClient{
		C_RO: sqlx.MustConnect("postgres", getDSN()),
		C_RW: sqlx.MustConnect("postgres", getDSN()),
	}

	client.C_RO.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
	client.C_RW.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
	client.C_RO.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	client.C_RW.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	client.C_RO.SetConnMaxLifetime(viper.GetDuration("database.conn_max_lifetime") * time.Second)
	client.C_RW.SetConnMaxLifetime(viper.GetDuration("database.conn_max_lifetime") * time.Second)

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
