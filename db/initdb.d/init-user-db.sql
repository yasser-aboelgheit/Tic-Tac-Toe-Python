--- Create permify Database
CREATE USER permify WITH
	LOGIN
	INHERIT
	PASSWORD 'permify';
CREATE DATABASE permify WITH
	OWNER permify;
GRANT ALL PRIVILEGES ON DATABASE permify TO permify;
