package managers

import "mango/src/logger"

var userManagerInstance UserManager

func init() {
	userManagerInstance = UserManager{
		users: make(map[string]User),
	}
}

type UserManager struct {
	users map[string]User
}

type User struct {
	EntityId int32
	UUID     []byte
	Position UserPosition

	Name string
}

type UserPosition struct {
	X     float64
	Y     float64
	Z     float64
	Yaw   uint8
	Pitch uint8
}

func GetUserManager() *UserManager {
	return &userManagerInstance
}

func (um *UserManager) AddUser(user User) {
	logger.Debug("userManager.AddUser() = %+v", user)
	um.users[user.Name] = user
}

func (um *UserManager) GetUser(username string) *User {
	if user, ok := um.users[username]; ok {
		return &user
	}
	return nil
}

func (um *UserManager) UpdateUser(user User) {
	um.users[user.Name] = user
}
