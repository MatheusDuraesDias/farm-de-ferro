package handler

import (
	"algorithm/mod/algoritmo"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Algo algoritmo.Algo
}

func (h *Handler) GetRecommendedSongs() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("userId")
		limit, err := strconv.Atoi(c.QueryParam("size"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		viewedPosts := strings.Split(c.QueryParam("excludedPostIds"), ",")

		recommendedSongs, err := h.Algo.Algoritmo(id, limit, viewedPosts)
		return c.JSON(http.StatusOK, recommendedSongs)
	}
}
