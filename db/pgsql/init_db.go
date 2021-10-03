package pgsql

import (
	"database/sql"
	"fmt"
	"urlshortener/app/starter"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return db, fmt.Errorf("error connection db: %v", err)
	}

	err = AddTables(db)
	if err != nil {
		return db, fmt.Errorf("error create tables: %v", err)
	}

	return db, nil
}

func InitStorers(db *sql.DB) *starter.Storers {
	u := NewUserStore(db)
	l := NewLinkStorer(db)
	lt := NewLinkTransitStorer(db)

	return starter.NewStorers(u, l, lt)
}

func AddTables(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE public.users (
		id uuid NOT NULL,
		login varchar(100) NOT NULL,
		"password" varchar(255) NOT NULL,
		CONSTRAINT firstkey PRIMARY KEY (id)
	);
	
	CREATE TABLE public.links (
		id uuid NOT NULL,
		created_at date NOT NULL,
		short_link varchar(255) NOT NULL,
		long_link varchar(255) NOT NULL,
		owner_id uuid NOT NULL,
		CONSTRAINT linkkey PRIMARY KEY (id)
	);
	
	ALTER TABLE public.links ADD CONSTRAINT links_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id);
	
	CREATE TABLE public.link_transitions (
		id uuid NOT NULL,
		link_id uuid NOT NULL,
		ip varchar(50) NOT NULL,
		used_count int4 NOT NULL,
		date timestamp NOT NULL,
		CONSTRAINT linktkey PRIMARY KEY (id)
	);
	
	ALTER TABLE public.link_transitions ADD CONSTRAINT link_transitions_link_id_fkey FOREIGN KEY (link_id) REFERENCES public.links(id);
	`)

	if err != nil {
		return err
	}

	return nil
}
