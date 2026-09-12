package repository

import (
	"database/sql"
	"fmt"
)

func InitDB(db *sql.DB) {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT,
			bio TEXT,
			avatar_path TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		fmt.Println("Erreur:", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS recipes (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			image_path TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		fmt.Println("Erreur:", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ingredients (
    		id SERIAL PRIMARY KEY,
    		recipe_id INTEGER NOT NULL,
    		name TEXT NOT NULL,
    		quantity TEXT NOT NULL,
    		unit TEXT NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id)
		)	
	`)
	if err != nil {
		fmt.Println("Erreur:", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS steps (
			id SERIAL PRIMARY KEY,
    		recipe_id INTEGER NOT NULL,
			step_order INTEGER NOT NULL,
			description TEXT NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id)
		)
	`)
	if err != nil {
		fmt.Println("Erreur:", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS friendship (
			id SERIAL PRIMARY KEY,
    		requester_id INTEGER NOT NULL,
			receiver_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (requester_id) REFERENCES users(id),
			FOREIGN KEY (receiver_id) REFERENCES users(id),
			UNIQUE(requester_id, receiver_id)
		)
	`)
	if err != nil {
		fmt.Println("Erreur:", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS likes (
			id SERIAL PRIMARY KEY,
    		user_id INTEGER NOT NULL,
			recipe_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (recipe_id) REFERENCES recipes(id),
			UNIQUE(user_id, recipe_id)
		)
	`)
	if err != nil {
		fmt.Println("Erreur:", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS favorites (
			id SERIAL PRIMARY KEY,
    		user_id INTEGER NOT NULL,
			recipe_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (recipe_id) REFERENCES recipes(id),
			UNIQUE(user_id, recipe_id)
		)
	`)
	if err != nil {
		fmt.Println("Erreur:", err)
	}
	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password)
	VALUES (0, 'system', 'system@cookapp.local', '')
	ON CONFLICT (id) DO NOTHING
	`)

	if err != nil {
		fmt.Println("Erreur:", err)
	}
}
