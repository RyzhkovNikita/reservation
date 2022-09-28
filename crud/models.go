package crud

type Profile struct {
	Id          int64  `orm:"auto"`
	Name        string `orm:"size(50)"`
	Description string `orm:"size(400)"`
	Address     string `orm:"size(100)"`
	LogoUrl     string
	Credentials *Credentials `orm:"rel(one);on_delete(cascade)"`
}

type Credentials struct {
	Id           int64    `orm:"auto"`
	Email        string   `orm:"size(30)"`
	PasswordHash string   `orm:"size(100)"`
	Profile      *Profile `orm:"reverse(one)"`
}
