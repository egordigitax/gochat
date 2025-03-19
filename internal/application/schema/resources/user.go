package resources

type User struct {
	Uid       string `json:"uid"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
