package managers

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
	Name string
}

func GetUserManager() UserManager {
	return userManagerInstance
}

func (um UserManager) AddUser(user User) {
	um.users[user.Name] = user
}
