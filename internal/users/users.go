package users

import (
	"database/sql"
	"log"

	database "github.com/soniceroy/hackernewsGQL/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
)

// User database representation
type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

// Create a user, checking password validity first
func (user *User) Create() {
	statement, err := database.Db.Prepare("INSERT INTO Users(Username, Password) VALUES(?,?)")
	print(statement)
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := HashPassword(user.Password)
	_, err = statement.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Authenticate user by finding user and comparing password
func (user *User) Authenticate() bool {
	statement, err := database.Db.Prepare("SELECT Password from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(user.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return CheckPasswordHash(user.Password, hashedPassword)
}

// CheckPasswordHash compares raw password with its hashed value
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserIDByUsername check if a user exists and get ID
func GetUserIDByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("SELECT ID from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var ID int
	err = row.Scan(&ID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}
	return ID, nil
}
