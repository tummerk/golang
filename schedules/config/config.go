package config

var DB_NAME = "test"
var DB_USER = "postgres"
var DB_PASS = "pass123"
var DB_PORT = "3242"
var DB_HOST = "localhost"

var TIME_NEXT_TAKINGS = 1000 //в минутах

var gRPC_PORT = "12345"

var ConnStr = "postgres://" + DB_USER + ":" + DB_PASS + "@" + DB_HOST + ":" + DB_PORT + "/" + DB_NAME + "?sslmode=disable"

var Key = []byte("8a1f3d9c7b2e45f60a9e8d2b4c3fds76")
