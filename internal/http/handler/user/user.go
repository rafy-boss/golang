package user

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/azan-boss/posty/internal/handler/auth"
	"github.com/azan-boss/posty/internal/storage"
	"github.com/azan-boss/posty/internal/types"
	"github.com/azan-boss/posty/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body==nil{
			c.JSON(http.StatusBadRequest,response.GeneralError(fmt.Errorf("request body is empty")))
			return
		
		}
	
		var user types.User
		if err:=c.ShouldBindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest,fmt.Errorf("failed to bind json:%s",err))
			return
		}


		v := validator.New()
		// v.RegisterValidation("isChinese", response.IsChinese)
		v.RegisterValidation("strongepwd",response.IsStrongePassword)
		// v.RegisterStructValidation(response.IsChineseStruct)
		// v.RegisterStructValidation(response.)
		if err := v.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, response.ValidationErrors(err.(validator.ValidationErrors)))
			return
		}
		id,err:=storage.RegisterUser(&user)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		slog.Info("User registered successfully", "id", id)
		fmt.Println(id)
		user.Password=""
		c.JSON(http.StatusOK,response.GeneralSuccessPlusData("Account created successfully",user))
		fmt.Println(user)
	}
}


func Login(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body==nil{
			c.JSON(http.StatusBadRequest,response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}
		var user types.User
		if err:=c.ShouldBindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest,fmt.Errorf("failed to bind json:%s",err))
			return
		}
		user, err := storage.Login(&user)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		token, err := auth.GenerateJWT(user.Username,user.ID)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		user.Password=""
		user.JWTToken=token
		c.JSON(http.StatusOK,response.GeneralSuccessPlusData("Login successful",user))	
	}
}

func GetUser(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id:=c.Query("id")
		idCon,_:=strconv.Atoi(id)
		user,err:=storage.GetUser(uint(idCon))
		if err!=nil{
			c.JSON(http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK,response.GeneralSuccessPlusData("User fetched successfully",user))	
	}
}