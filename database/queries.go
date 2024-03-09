package database

const (
	queryCreateTables string = `BEGIN;

	CREATE TABLE IF NOT EXISTS users(
		id BIGINT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
		is_deleted BOOLEAN NOT NULL DEFAULT false,
		is_disabled BOOLEAN NOT NULL DEFAULT false,
		password_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions(
		id BIGINT PRIMARY KEY,
		uid BIGINT NOT NULL,
		refresh_token VARCHAR(64) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
		expires_at TIMESTAMP NOT NULL,
		valid BOOLEAN NOT NULL DEFAULT true,
		CONSTRAINT fk_sessions_uid FOREIGN KEY (uid) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS urls(
		id BIGINT PRIMARY KEY,
		uid BIGINT NOT NULL,
		url_path TEXT NOT NULL,
		link TEXT NOT NULL,
		name VARCHAR(20),
		created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
		valid_until TIMESTAMP,
		is_deleted BOOLEAN NOT NULL DEFAULT false,
		CONSTRAINT fk_urls_uid FOREIGN KEY (uid) REFERENCES users(id)
	);
	
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE is_deleted = false;

	CREATE INDEX IF NOT EXISTS idx_urls_url_path ON urls(url_path) WHERE is_deleted = false;
	CREATE INDEX IF NOT EXISTS idx_urls_uid ON urls(uid) WHERE is_deleted = false;
	CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls(created_at) WHERE is_deleted = false;

	COMMIT;	
	`

	queryInsertUser string = `INSERT INTO users(id, name, email, password_hash) VALUES ($1, $2, $3, $4);`

	queryGetUserByEmail string = `SELECT * FROM users WHERE email = $1 LIMIT 1;`
	queryGetUserById    string = `SELECT * FROM users WHERE id = $1 LIMIT 1;`

	queryInsertSession string = `INSERT INTO sessions(id, uid, refresh_token, expires_at) VALUES ($1, $2, $3, $4);`

	queryInsertUrl     string = `INSERT INTO urls(id, uid, link, url_path, name, valid_until) VALUES ($1, $2, $3, $4, $5, $6);`
	queryGetLinkByPath string = `SELECT id, link FROM urls WHERE url_path = $1 LIMIT 1;`
	queryGetUserLinks  string = `SELECT * FROM urls WHERE uid = $1 AND is_deleted = false LIMIT 50;`
)
