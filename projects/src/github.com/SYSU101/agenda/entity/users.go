package entity

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/SYSU101/agenda/fmterror"
	"github.com/SYSU101/agenda/logger"
)

const emailRegExp = "^[-!#$%&'*+/0-9=?A-Z^_`a-z{|}~]+(\\.[-!#$%&'*+/0-9=?A-Z^_`a-z{|}~]+)*@[A-Za-z0-9]([A-Za-z0-9-]{0,61}[A-Za-z0-9])?(\\.[A-Za-z0-9]([A-Za-z0-9-]{0,61}[A-Za-z0-9])?)*$"
const userFilePath = "./user.json"

var CurrentUser *User = nil

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	LoggedIn bool   `json:"logged_in"`
}

func (user *User) Validate() error {
	if user.Username == "" {
		return fmterror.New("User's information is invalid: username is required.")
	}
	if user.Email == "" {
		return fmterror.New("User's information is invalid: email is required.")
	}
	if _, isExist := users[user.Username]; isExist {
		return fmterror.New("User's information is invalid: username `%v` has been existed.", user.Username)
	}
	if matched, err := regexp.Match(emailRegExp, []byte(user.Email)); err != nil {
		logger.Printf("[error] %v\n", err)
		return fmterror.New("Internal error happened, please contact author.")
	} else if !matched {
		return fmterror.New("User's information is invalid: user's email is invalid.")
	}
	return nil
}

var users = make(map[string]*User)

func QueryUser(username string) (*User, error) {
	if user, isExist := users[username]; isExist {
		return user, nil
	} else {
		return nil, fmterror.New("User `%v` does not exist.", username)
	}
}

func Register(user *User) error {
	if err := user.Validate(); err != nil {
		return err
	} else {
		user.Password = fmt.Sprintf("%x", md5.Sum([]byte(user.Password)))
		users[user.Username] = user
		return nil
	}
}

func Login(username, password string) error {
	if CurrentUser != nil {
		return fmterror.New("There has been other user logged in, pleas log out first")
	}
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	if user, err := QueryUser(username); err != nil {
		return err
	} else if password != user.Password {
		return fmterror.New("Password incorrect, please try again")
	} else {
		user.LoggedIn = true
		CurrentUser = user
		return nil
	}
}

func Logout() error {
	if CurrentUser == nil {
		return fmterror.New("There are no user logged in")
	}
	CurrentUser.LoggedIn = false
	CurrentUser = nil
	return nil
}

func SaveUser() error {
	userFile, err := os.Create(userFilePath)
	if err != nil {
		logger.Printf("[error] %v\n", err)
		return fmterror.New("Failed when saving user's data: can not open data file")
	}
	data := make([]*User, 0, len(users))
	for _, user := range users {
		if user != nil {
			data = append(data, user)
		}
	}
	if databytes, err := json.MarshalIndent(data, "", "\t"); err == nil {
		userFile.Write(databytes)
	} else {
		logger.Printf("[error] %v\n", err)
		return fmterror.New("Failed when saving user's data: serialize data")
	}
	return nil
}

func init() {
	userFile, err := os.OpenFile(userFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Printf("[error] %v\n", err)
		return
	}
	databytes, err := ioutil.ReadAll(userFile)
	if err != nil {
		logger.Printf("[error] %v\n", err)
	}
	data := make([]*User, 0)
	if len(databytes) == 0 {
		databytes = append(databytes, '[', ']')
	}
	if err = json.Unmarshal(databytes, &data); err != nil {
		logger.Printf("[error] %v\n", err)
		return
	}
	for _, user := range data {
		if user != nil {
			if user.LoggedIn {
				CurrentUser = user
			}
			users[user.Username] = user
		}
	}
	return
}
