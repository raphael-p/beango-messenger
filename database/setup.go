package database

func handleError(err error) {
	if err != nil {
		panic("failed to setup database: " + err.Error())
	}
}

func Setup(conn *MongoConnection) {
	_, err := conn.Exec(`SET TIME ZONE 'UTC';`)
	handleError(err)

	_, err = conn.Exec(`
	CREATE TABLE IF NOT EXISTS "user" (
		id SERIAL PRIMARY KEY,
		username VARCHAR(25) NOT NULL,
		display_name VARCHAR(25) NOT NULL,
		key BYTEA NOT NULL,
		created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
		last_updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
	)`)
	handleError(err)

	_, err = conn.Exec(`
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_type WHERE typname = 'chattype'
		) THEN
			CREATE TYPE CHATTYPE AS ENUM (
				'private',
				'group'
			);
		END IF;
	END$$;`)
	handleError(err)

	_, err = conn.Exec(`
	CREATE TABLE IF NOT EXISTS chat (
		id SERIAL PRIMARY KEY,
		type CHATTYPE NOT NULL,
		name VARCHAR(25),
		created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
		last_updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
	)`)
	handleError(err)

	_, err = conn.Exec(`
	CREATE TABLE IF NOT EXISTS chat_users (
		id SERIAL PRIMARY KEY,
		chat_id UUID NOT NULL REFERENCES chat(id),
		user_id UUID NOT NULL REFERENCES "user"(id),
		created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
		CONSTRAINT unique_combination UNIQUE (chat_id, user_id)
	)`)
	handleError(err)

	_, err = conn.Exec(`
	CREATE TABLE IF NOT EXISTS message (
		id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES "user"(id),
		chat_id UUID NOT NULL REFERENCES chat(id),
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
		last_updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
	)`)
	handleError(err)
}
