package model

// Router ...
type Router struct {
	ID       uint   `db:"id"`
	Name     string `db:"name"`
	Address  string `db:"address"`
	Login    string `db:"login"`
	Password string `db:"password"`
}
