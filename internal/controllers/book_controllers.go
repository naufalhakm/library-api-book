package controllers

import (
	"library-api-book/internal/commons/response"
	"library-api-book/internal/models"
	"library-api-book/internal/params"
	"library-api-book/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookController interface {
	CreateBook(ctx *gin.Context)
	GetDetailBook(ctx *gin.Context)
	UpdateBook(ctx *gin.Context)
	DeleteBook(ctx *gin.Context)
	GetAllBooks(ctx *gin.Context)
	GetRecommendationBook(ctx *gin.Context)
}

type BookControllerImpl struct {
	BookService services.BookService
}

func NewBookController(bookService services.BookService) BookController {
	return &BookControllerImpl{
		BookService: bookService,
	}
}

func (controller *BookControllerImpl) CreateBook(ctx *gin.Context) {
	var req = new(params.BookRequest)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}

	custErr := controller.BookService.CreateBook(ctx, req)
	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	resp := response.CreatedSuccessWithPayload("Success create data book")
	ctx.JSON(resp.StatusCode, resp)
}

func (controller *BookControllerImpl) GetDetailBook(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}

	result, custErr := controller.BookService.GetDetailBook(ctx, uint64(id))

	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get detail book", result)
	ctx.JSON(resp.StatusCode, resp)
}

func (controller *BookControllerImpl) UpdateBook(ctx *gin.Context) {
	var req = new(params.BookRequest)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))

	custErr := controller.BookService.UpdateBook(ctx, uint64(id), req)
	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success update data book", nil)
	ctx.JSON(resp.StatusCode, resp)
}

func (controller *BookControllerImpl) DeleteBook(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}

	custErr := controller.BookService.DeleteBook(ctx, uint64(id))
	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success delete data book", nil)
	ctx.JSON(resp.StatusCode, resp)
}

func (controller *BookControllerImpl) GetAllBooks(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	search := ctx.Query("search")

	pageNum := 1
	limitSize := 5

	if page != "" {
		parsedPage, err := strconv.Atoi(page)
		if err == nil && parsedPage > 0 {
			pageNum = parsedPage
		}
	}

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err == nil && parsedLimit > 0 {
			limitSize = parsedLimit
		}
	}

	pagination := models.Pagination{
		Page:     pageNum,
		Offset:   (pageNum - 1) * limitSize,
		PageSize: limitSize,
	}

	result, custErr := controller.BookService.GetAllBooks(ctx, &pagination, search)

	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	type Response struct {
		Books      interface{} `json:"books"`
		Pagination interface{} `json:"pagination"`
	}

	var responses Response
	responses.Books = result
	responses.Pagination = pagination

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get data books", responses)
	ctx.JSON(resp.StatusCode, resp)

}
func (controller *BookControllerImpl) GetRecommendationBook(ctx *gin.Context) {
	authId := ctx.GetInt("authId")

	result, custErr := controller.BookService.GetRecommendationBook(ctx, uint64(authId))

	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get data recomendation books", result)
	ctx.JSON(resp.StatusCode, resp)
}
