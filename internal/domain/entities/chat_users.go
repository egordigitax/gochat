package entities

type ChatUsers struct {
	UID     string              `db:"uid"`
	Users   map[string]User     `db:"users"`
	Roles   map[string]ChatRole `db:"roles"`
	Blocked map[string]User     `db:"blocked"`
}
