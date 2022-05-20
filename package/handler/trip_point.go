package handler

import (
	"net/http"
	"strconv"

	"github.com/folins/biketrackserver"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createPoint(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid trip id param")
		return
	}

	var input biketrackserver.TripPoint
	if err := c.BindJSON(&input); err != nil { 
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TripPoint.Create(userId, tripId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}


func (h *Handler) getAllPoints(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid trip id param")
		return
	}

	points, err := h.services.TripPoint.GetAll(userId, tripId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, points)
}

func (h *Handler) getPointById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid trip id param")
		return
	}

	pointId, err := strconv.Atoi(c.Param("point_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid point id param")
		return
	}

	points, err := h.services.TripPoint.GetById(userId, tripId, pointId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, points)
}