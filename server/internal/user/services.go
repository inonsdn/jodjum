package user

type UserService struct {
	userRepo *UserRepo
}

func NewService(userRepo *UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}
