package database

import (
	"database/sql"
)

func makeErrorMessage(err error) string {
	return "failed to setup database: " + err.Error()
}

func handleError(tx *sql.Tx, err error) {
	if err != nil {
		message := makeErrorMessage(err)
		if err := tx.Rollback(); err != nil {
			message += " (transaction rollback failed: " + err.Error() + ")"
		} else {
			message += " (transaction rollback successful)"
		}
		panic(message)
	}
}

func Setup(conn *MongoConnection) {
	tx, err := conn.Begin()
	if err != nil {
		panic(makeErrorMessage(err))
	}

	_, err = tx.Exec(`SET TIME ZONE 'UTC';`)
	handleError(tx, err)

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS "user" (
		id SERIAL PRIMARY KEY,
		username VARCHAR(25) NOT NULL,
		display_name VARCHAR(25) NOT NULL,
		key BYTEA NOT NULL,
		created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
		last_updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
	)`)
	handleError(tx, err)

	_, err = tx.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS user_username_unique_idx ON "user" (username);
	`)
	handleError(tx, err)

	_, err = tx.Exec(`
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_type WHERE typname = 'chattype'
		) THEN
			CREATE TYPE CHATTYPE AS ENUM (
				'note',
				'private',
				'group'
			);
		END IF;
	END$$;`)
	handleError(tx, err)

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS chat (
		id SERIAL PRIMARY KEY,
		type CHATTYPE NOT NULL,
		name VARCHAR(25) NOT NULL,
		created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
		last_updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
	)`)
	handleError(tx, err)

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS chat_users (
		id SERIAL PRIMARY KEY,
		chat_id INT NOT NULL REFERENCES chat(id),
		user_id INT NOT NULL REFERENCES "user"(id),
		created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
		CONSTRAINT chat_users_unique_constraint UNIQUE (chat_id, user_id)
	)`)
	handleError(tx, err)

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS message (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES "user"(id),
		chat_id INT NOT NULL REFERENCES chat(id),
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
		last_updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
	)`)
	handleError(tx, err)

	err = tx.Commit()
	handleError(tx, err)
}
