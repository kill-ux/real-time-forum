// models/user.go
package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"forum/db"
	"forum/utils"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	UUID      any    `json:"uuid"`
	UUID_EXP  int64  `json:"uuid_exp"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	CreatedAt int64  `json:"created_at"`
	LastSeen  int64  `json:"last_seen"`
	Image     string `json:"image"`
}

func (u *User) BeforeCreate() error {
	u.CreatedAt = time.Now().Unix()

	// trim spaces
	u.Nickname = strings.ToLower(strings.TrimSpace(u.Nickname))
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Gender = strings.ToLower(strings.TrimSpace(u.Gender))

	// nickname
	nicknamePattern := regexp.MustCompile(`^[a-zA-Z0-9_-]{3,40}$`)
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._+-]{3,20}@[a-zA-Z0-9.-]{3,20}\.[a-zA-Z]{2,10}$`)
	passPattern := regexp.MustCompile(`^.{8,100}$`)
	firstNamePattern := regexp.MustCompile(`^[A-Za-z]{3,40}$`)
	lastNamePattern := regexp.MustCompile(`^[A-Za-z]{3,40}$`)
	gender := regexp.MustCompile(`^(male|female)$`)

	if !nicknamePattern.MatchString(u.Nickname) ||
		!emailPattern.MatchString(u.Email) ||
		!passPattern.MatchString(u.Password) ||
		!firstNamePattern.MatchString(u.FirstName) ||
		!lastNamePattern.MatchString(u.LastName) ||
		!gender.MatchString(u.Gender) ||
		u.Age < 10 || u.Age > 200 {
		return errors.New("errors in Patterns")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Image = strings.ToUpper(string(u.FirstName[0])) + ".png"
	return nil
}

func (u *User) VerifyPassword(password string) bool {
	passPattern := regexp.MustCompile(`^.{8,100}$`)
	return passPattern.MatchString(u.Password) && bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password)) == nil
}

func (user *User) CreateUser() error {
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	_, err := db.DB.Exec(`INSERT INTO users VALUES (NULL,NULL,0,?,?,?,?,?,?,?,?,0,?)`, utils.GetExecFields(user, "ID", "UUID", "UUID_EXP", "LastSeen")...)
	return err
}

func (user *User) UpdateUuid() error {
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	_, err := db.DB.Exec(`UPDATE users SET uuid = ? , uuid_exp = ? WHERE id = ?`, user.UUID, user.UUID_EXP, user.ID)
	return err
}

func Logout(id int) error {
	_, err := db.DB.Exec(`UPDATE users SET uuid = NULL , uuid_exp = 0 WHERE id = ?`, id)
	return err
}

// SQL Injection Prevention
func GetUserBy(id string) (*User, error) {
	id = strings.TrimSpace(strings.ToLower(id))
	nicknamePattern := regexp.MustCompile(`^[a-zA-Z0-9_-]{3,40}$`)
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._+-]{3,20}@[a-zA-Z0-9.-]{3,20}\.[a-zA-Z]{2,10}$`)
	if !nicknamePattern.MatchString(id) && !emailPattern.MatchString(id) {
		return nil, errors.New("error in fileds")
	}
	user := &User{}
	err := db.DB.QueryRow("SELECT * FROM users WHERE email = ? OR nickname = ?", id, id).Scan(utils.GetScanFields(user)...)
	fmt.Println(err)

	return user, err
}
