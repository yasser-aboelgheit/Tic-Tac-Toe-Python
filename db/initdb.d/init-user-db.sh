#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE ROLE permify LOGIN PASSWORD 'permify';
	CREATE DATABASE permify;
	GRANT ALL ON permify TO permify;
EOSQL
