// database/postgres.go
package database

import (
	"context"
	"fmt"
	log "log"
	"os"
	"sync"

	"gocloud.dev/postgres"
	_ "gocloud.dev/postgres/awspostgres"
	logger "neema.co.za/rest/utils/logger"

	"xorm.io/xorm"
	"xorm.io/xorm/core"
	"xorm.io/xorm/names"
	//"gorm.io/driver/postgres"
)

type Database struct {
	*xorm.Engine
}

var engine *Database
var once sync.Once

func GetDatabase() *Database {
	once.Do(func() {
		dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DB"))

		db, err := postgres.Open(context.Background(), dbURL)
		if err != nil {
			return
		}
		//TODO : Add Pooling mechanism
		//defer db.Close()
		coreDB := core.FromDB(db)
		xengine, err := xorm.NewEngineWithDB("postgres", dbURL, coreDB)
		if err != nil {
			log.Fatalf("Error creating XORM engine: %v", err)
		}

		// Enable query logging

		xengine.SetLogger(logger.GetCustomXormLogger())
		xengine.SetMaxIdleConns(1)
		xengine.SetMapper(names.GonicMapper{})

		if err := xengine.Ping(); err != nil {
			log.Fatalf("Error pinging database: %v", err)
		}

		logger.Info("Connected to PostgreSQL database")
		engine = &Database{xengine}
	})
	return engine
}
