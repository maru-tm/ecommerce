package domain

type UserRepository interface {
	CreateUser(user *User) (*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	ListUsers() ([]User, error)
	UpdateUser(user *User) (*User, error)
	DeleteUser(id string) error
}

type UserUseCase interface {
	CreateUser(user *User) (*User, error)
	GetUserByID(id string) (*User, error)
	ListUsers() ([]User, error)
	UpdateUser(id string, user *User) (*User, error)
	DeleteUser(id string) error
}
