package post

import (
	"fmt"
	"net/http"

	"github.com/azan-boss/posty/internal/storage"
	"github.com/azan-boss/posty/internal/types"
	"github.com/azan-boss/posty/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		var post types.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		username , _ := c.Get("username")
		userId ,_ :=c.Get("userId")
		fmt.Println("username",username)
		fmt.Println("userId",userId)
		post.UserId = userId.(uint)
		// Validate that UserId is provided
		if post.UserId == 0 {
			c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("user ID is required")))
			return
		}

		v := validator.New()
		if err := v.Struct(post); err != nil {
			c.JSON(http.StatusBadRequest, response.ValidationErrors(err.(validator.ValidationErrors)))
			return
		}

		if err := storage.CreatePost(&post); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, response.GeneralSuccess("Post created successfully"))
	}
}
