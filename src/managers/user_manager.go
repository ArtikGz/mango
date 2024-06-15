package managers

var userManagerInstance userManager

func init() {
	userManagerInstance = userManager{
		users: make([]User, 0),
	}
}

type userManager struct {
	users []User
}

type User struct {
	Name string
}

func GetUserManager() userManager {
	return userManagerInstance
}

func (um userManager) AddUser(user User) {
	um.users = append(um.users, user)
}
