package service

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) SendRegisterMobileCode(mobile string) error {
	return nil
}
