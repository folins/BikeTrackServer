package biketrackserver

type User struct {
	Id          int    `json:"-" db:"id"`
	Name        string `json:"name"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password"`
	ConfirmCode int    `json:"confirm_code" db:"confirm_code"`
}

type UserUpdateInput struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	ConfirmCode *int
}

