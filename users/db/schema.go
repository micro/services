package db

var (
	accountSchema = `CREATE TABLE IF NOT EXISTS accounts (
id varchar(36) primary key,
username varchar(255),
email varchar(255),
salt varchar(16),
password text,
created integer,
updated integer,
unique (username),
unique (email));`
	sessionSchema = `CREATE TABLE IF NOT EXISTS sessions (
id varchar(255) primary key,
username varchar(255),
created integer,
expires integer);`
)
