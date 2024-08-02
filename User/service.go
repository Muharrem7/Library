package User

import (
	"library/Book"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (us *UserService) InsertUser(user User) error {
	return us.userRepo.InsertUser(user)
}

func (us *UserService) GetUsers() ([]User, error) {

	return us.userRepo.GetUsers()

}
func (us *UserService) GetUserIdById(req Book.AssignRequest) (int, error) {
	return us.userRepo.GetUserIdById(req)

}
func (us *UserService) GetAssignedUser() ([]Assignments, error) {
	return us.userRepo.GetAssignedUser()

}
func (us *UserService) GetUserById(req IdRequest) ([]User, error) {
	return us.userRepo.GetUserById(req)
}
