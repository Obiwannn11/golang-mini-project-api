package utils

import (
	"rakamin-evermos/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginationInput struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginationResult struct {
    Data       interface{} `json:"data"`
    TotalData  int64       `json:"total_data"`
    TotalPage  int64       `json:"total_page"`
    CurrentPage int         `json:"current_page"`
    NextPage   *int        `json:"next_page"`
    PrevPage   *int        `json:"prev_page"`
}

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit

		return db.Offset(offset).Limit(limit)
	}
}

// get 'page' dan 'limit' from query URL
func GetPaginationFromQuery(c *gin.Context) PaginationInput {
	// Set default
	page := 1
	limit := 10

	// take query param if exist
	pageQuery := c.Query("page")
	if pageQuery != "" {
		p, err := strconv.Atoi(pageQuery)
		if err == nil && p > 0 {
			page = p
		}
	}

	limitQuery := c.Query("limit")
	if limitQuery != "" {
		l, err := strconv.Atoi(limitQuery)
		if err == nil && l > 0 {
			limit = l
		}
	}

	return PaginationInput{Page: page, Limit: limit}
}

//  count and format pagination result
func GeneratePaginationResult(data interface{}, totalData int64, page, limit int) PaginationResult {
    totalPage := totalData / int64(limit)
    if totalData % int64(limit) != 0 {
        totalPage++
    }

    var next, prev *int
    if page < int(totalPage) {
        n := page + 1
        next = &n
    }
    if page > 1 {
        p := page - 1
        prev = &p
    }

    return PaginationResult{
        Data:        data,
        TotalData:   totalData,
        TotalPage:   totalPage,
        CurrentPage: page,
        NextPage:    next,
        PrevPage:    prev,
    }
}