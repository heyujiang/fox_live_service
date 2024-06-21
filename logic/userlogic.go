package logic

import (
	"fmt"
	"time"
)

var ULogic = &UserLogic{
	Now: time.Now(),
}

func NewUserLogic() *UserLogic {
	return &UserLogic{
		Now: time.Now(),
	}
}

type UserLogic struct {
	Now time.Time
}

func (ul *UserLogic) Login() {
	fmt.Println(ul.Now)
}
