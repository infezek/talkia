package util_pagination

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func New(page, perPage int32, orderBy OrderBy) *Pagination {
	offset := (page - 1) * perPage
	return &Pagination{
		Limit:   perPage,
		Offset:  offset,
		Page:    page,
		OrderBy: string(orderBy),
	}
}

func NewWithParams(c *fiber.Ctx, params Pagination) *Pagination {
	paramPerPage := c.Query("per_page", fmt.Sprintf("%d", params.Limit), "10")
	paramPage := c.Query("page", fmt.Sprintf("%d", params.Offset), "1")
	paramOrderBy := c.Query("order_by", string(params.OrderBy), "DESC")
	perPage, err := strconv.Atoi(paramPerPage)
	if err != nil || perPage < 1 {
		perPage = 10
	}
	if err != nil || perPage > 50 {
		perPage = 50
	}
	page, err := strconv.Atoi(paramPage)
	if err != nil || page < 1 {
		page = 1
	}
	orderBy := strings.ToUpper(paramOrderBy)
	if orderBy != "DESC" && orderBy != "ASC" {
		paramOrderBy = "DESC"
	}
	logrus.WithFields(logrus.Fields{
		"perPage": perPage,
		"page":    page,
		"orderBy": orderBy,
	}).Info("NewWithParams")
	return New(int32(page), int32(perPage), OrderBy(orderBy))
}

type Pagination struct {
	Offset  int32
	Page    int32
	Limit   int32
	OrderBy string
}
type OrderBy string

var (
	OrderByDesc OrderBy = "DESC"
	OrderByAsc  OrderBy = "ASC"
)
