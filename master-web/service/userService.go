package service

import "fmt"

type UserService struct {

}

func NewUserService() ServiceInterface {
	return &UserService{
	}
}

func (user *UserService) AAA()  {
	fmt.Printf("dasdf")
}