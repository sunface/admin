package misc

type User struct {
	ID         int    `json:"id" db:"id"`
	UserName   string `json:"username" db:"username"`
	PW         string `json:"-" db:"pw"`
	Priv       string `json:"priv" db:"priv"`
	CreateDate string `json:"create_date" db:"create_date"`
	ModifyDate string `json:"modify_date" db:"modify_date"`
}
