package services

type IUserRepository interface {
}

type UserService struct {
	usrRepo IUserRepository
}

func NewUser(usrRepo IUserRepository) *UserService {
	return &UserService{
		usrRepo: usrRepo,
	}
}
