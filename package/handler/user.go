package handler

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/folins/biketrackserver"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type emailInput struct {
	Email string `json:"email" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input emailInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	rand.Seed(time.Now().Unix())
	code := 100000 + rand.Intn(999999-100000)

	id, err := h.services.User.Create(input.Email, code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.services.SMTP.SendConfirmCode(input.Email, code)

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.User.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) resetPassword(c *gin.Context) {
	var input emailInput
	var newDetails biketrackserver.UserUpdateInput

	if err := c.BindJSON(&input); err != nil {
		logrus.WithField("handler", "resetPassword").Errorf("Failed to bind sign in structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.GetIdByEmail(input.Email)
	if err != nil {
		logrus.WithField("handler", "resetPassword").Errorf("Failed to get user with input email: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	rand.Seed(time.Now().Unix())
	code := 100000 + rand.Intn(999999-100000)
	newDetails.ConfirmCode = &code

	if err := h.services.User.Update(id, newDetails); err != nil {
		logrus.WithField("handler", "resetPassword").Errorf("Failed update user: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.services.SMTP.SendResetCode(input.Email, code); err != nil {
		logrus.WithField("handler", "resetPassword").Errorf("Failed to send reset code to user: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

type changePasswordInput struct {
	PrevPassword string `json:"prev_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *Handler) changePassword(c *gin.Context) {
	h.userIdentity(c)
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input changePasswordInput
	if err := c.BindJSON(&input); err != nil {
		logrus.WithField("handler", "changePassword").Errorf("Failed to bind sign in structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.User.CheckPassword(userId, input.PrevPassword); err != nil {
		logrus.WithField("handler", "changePassword").Errorf("Wrong user password: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var newUserDetails biketrackserver.UserUpdateInput
	*newUserDetails.Password = input.NewPassword

	if err := h.services.User.Update(userId, newUserDetails); err != nil {
		logrus.WithField("handler", "changePassword").Errorf("Failed update user: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

type verifyCodeInput struct {
	Email string `json:"email" binding:"required"`
	Code int `json:"confirm_code" binding:"required"`
}

func (h *Handler) verifyCode(c *gin.Context) {
	var input verifyCodeInput
	if err := c.BindJSON(&input); err != nil {
		logrus.WithField("handler", "verifyCode").Errorf("Failed to bind input in structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.User.CheckConfirmCode(input.Email, input.Code); err != nil {
		logrus.WithField("handler", "verifyCode").Errorf("Wrong user confirm code: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) checkEmailExistence(c *gin.Context) {
	var input emailInput
	if err := c.BindJSON(&input); err != nil {
		logrus.WithField("handler", "checkEmailExistence").Errorf("Failed to bind sign in structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.User.CheckEmailExistence(input.Email)
	if err != nil {
		logrus.WithField("handler", "checkEmailExistence").Errorf("Failed to get sql query: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: strconv.FormatBool(result),
	})
}

type setPasswordInput struct {
	Email string `json:"email" binding:"required"`
	Code int `json:"confirm_code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) setPassword(c *gin.Context) {
	var input setPasswordInput
	if err := c.BindJSON(&input); err != nil {
		logrus.WithField("handler", "setPassword").Errorf("Failed to bind input in structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Code < 100000 {
		err := errors.New("code should have six numbers")
		logrus.WithField("handler", "setPassword").Errorf("Invalide confirm code: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.User.SetPassword(input.Email, input.Password, input.Code); err != nil {
		logrus.WithField("handler", "setPassword").Errorf("Wrong user confirm code: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}