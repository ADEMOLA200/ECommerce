package models

type Config struct {
	Port int
	Env  string
	API  string
	DB   struct {
		DSN string
	}
	Stripe struct {
		Secret string
		Key    string
	}
}
