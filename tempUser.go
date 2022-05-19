package biketrackserver

type tempUser struct {
	Id          int    `json:"-" db:"id"`
	Email       string `json:"email" binding:"required"`
	ConfirmCode int    `json:"confirm_code"`
}