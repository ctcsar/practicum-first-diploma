-- +goose Up
CREATE TABLE users (
	id VARCHAR(255) PRIMARY KEY,
	login VARCHAR(255) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	registerDate TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;
