package config

var DB_NAME = "test"
var DB_USER = "postgres"
var DB_PASS = "pass123"
var DB_PORT = "3242"
var DB_HOST = "localhost"

var TIME_NEXT_TAKINGS = 120 //в минутах

var ConnStr = "postgres://" + DB_USER + ":" + DB_PASS + "@" + DB_HOST + ":" + DB_PORT + "/" + DB_NAME + "?sslmode=disable"
