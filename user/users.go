package user

type User struct {
	Username     string
	PasswordHash string
	Email        string
}

type Users map[string]User

var DB = Users{
	"user": User{
		Username: "user",
		// bcrypt has for "password"
		PasswordHash: "$2a$10$KgFhp4HAaBCRAYbFp5XYUOKrbO90yrpUQte4eyafk4Tu6mnZcNWiK",
		Email:        "user@example.com",
	},
}
