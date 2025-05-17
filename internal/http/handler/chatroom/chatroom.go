package chatroom

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/azan-boss/posty/internal/storage"
	"github.com/azan-boss/posty/internal/types"
	"github.com/azan-boss/posty/internal/utils/response"
	"github.com/azan-boss/posty/internal/utils/slug"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
func CreatChatRoom(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var chatroom types.ChatRoom
		if err := c.ShouldBindJSON(&chatroom); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}
		

		//  validating request body 

		v :=validator.New()
		err :=v.Struct(chatroom)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ValidationErrors(err.(validator.ValidationErrors)))
			return
		}
		chatroom.Slug = slug.GenerateSlug(chatroom.Name)

		id,err:=storage.CreateChatRoom(&chatroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		// err=storage.CreateMessage(&types.Message{
		// 	Content: fmt.Sprintf("Welcome to the %s",chatroom.Name),
		// 	Type:    "message",
		// 	ChatRoomId: id,
		// 	UserID:    0,
		// })
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, response.GeneralError(err))
		// 	return
		// }
		slog.Info("Chatroom created successfully", "id", id)
		c.JSON(http.StatusOK, response.GeneralSuccess(fmt.Sprintf("Chatroom created successfully with id %d",id)))
	}
}