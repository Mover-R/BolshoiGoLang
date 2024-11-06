package main

import (
	"BolshiGoLang/fileutils"
	"BolshiGoLang/internal/pkg/server"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

const(
	createStorageDB = `CREATE TABLE IF NOT EXISTS core (
						version bigserial PRIMARY KEY,
						timestamp bigint NOT NULL,
						payload JSONB NOT NULL
	)`
)

func insertState(db *sql.DB, timestamp int64, payload []byte) error {
    _, err := db.Exec(`
        INSERT INTO core (timestamp, payload) VALUES ($1, $2);
    `, timestamp, payload)
	
    return err
}

func main() {
	dburl := os.Getenv("DATABASE_URL") + "?sslmode=disable"

	db, err := sql.Open("postgres", dburl)
	if err!=nil{
		log.Fatal("open", err)
		return
	}

	defer db.Close()

	if err:=db.Ping(); err!=nil{
		log.Fatal("ping", err)
	}

	_, err = db.Exec(createStorageDB)
	if err!=nil{
		log.Fatal(err)
	}

	r, err := fileutils.DataStorageFileRead()
	if err != nil {
		panic(err)
	}

	port:="8090" //os.Getenv("BASIC_SERVER_PORT")
	s := server.NewServer(":" + port, r)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(2*time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			data, err := json.Marshal(r)
			if err!=nil{
				return
			}
			err = insertState(db, time.Now().Unix(), data)
			if err!=nil{
				log.Fatal()
			}
		}
	}()

	go func() {
		s.Start()
		_, err = db.Exec(`
			CREATE OR REPLACE FUNCTION maintain_latest_five_versions()
			RETURNS TRIGGER AS $$
			BEGIN
				DELETE FROM core
				WHERE version NOT IN (
					SELECT version FROM core
					ORDER BY version DESC
					LIMIT 5
				);
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;

			DROP TRIGGER IF EXISTS after_insert_core ON core;

			CREATE TRIGGER after_insert_core
			AFTER INSERT ON core
			FOR EACH ROW EXECUTE FUNCTION maintain_latest_five_versions();
		`)
		if err != nil {
			log.Fatal(err)
		}	
	}()

	<-signalChan

	err = fileutils.DataStorageFileWrite(r)
	if err != nil {
		return
	}
}
