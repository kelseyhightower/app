package user

type User struct {
	Username     string
	PasswordHash string
	Email        string
}

type Users map[string]User

var DB = Users{
	"khightower": User{
		Username: "khightower",
		// bcrypt hash for 123456789
		PasswordHash: "$2y$05$VcnPUSu3n41uY/frKwyraeQpPaZt.rWlROlTuIoajlrBQffHp9GA6",
		Email:        "khightower@example.com",
	},
}
