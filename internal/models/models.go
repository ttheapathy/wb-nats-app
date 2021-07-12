package models

type Message struct {
	Id   int64  `db:"id"`
	Text string `db:"text"`
}

type ErrorMessage struct {
	Message string
	Detail  string
}
