package types

import (
	"log/slog"
	"time"

	"github.com/azan-boss/posty/internal/handler/auth"
	"github.com/azan-boss/posty/internal/utils/method"
	"gorm.io/gorm"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type status string
const (
	Busy = "busy" 
	offline = "offline" 
	Online = "Online"
)
type User struct {
	gorm.Model
	Username string `json:"username" validate:"required,min=6,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20,strongepwd"`
	Role     Role   `json:"role" validate:"required,oneof=admin user"`
	JWTToken string `json:"jwt_token,omitempty"`
	Posts []Post `json:"posts,omitempty"`
	Status	status `json:"status" validate:"oneof=busy offline online"`
	LastActive time.Time `json:"lastseen"`
	ChatRooms []ChatRoom `json:"chatroom" gorm:"many2many:user_chatrooms;"`
}


 type ChatRoom struct{
	gorm.Model
	Name string `json:"name" validate:"required,min=5"`
	Description string `json:"description"`
	Members []User `json:"members" gorm:"many2many:user_chatrooms;"`
	Slug string `json:"slug"`
 }

type Message struct{
	gorm.Model
	Content string 
	Type string
	ChatRoomId uint
	UserID uint
	User User `gorm:"foreignKey:UserID"`
	ChatRoom ChatRoom `gorm:"foreignKey:ChatRoomId"`
}


type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title" validate:"required,min=10"`
	Content string `json:"content" validate:"required,min=10,max=2000"`
	UserId uint   `json:"user_id" validate:"required"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	hash, err := method.GenerateHashPassword(u.Password)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	u.Password = hash
	slog.Info("Password hashed successfully", "password", u.Password)
	return nil
}

func (u *User) AfterSave(tx *gorm.DB) (err error) {
	jwtToken, err := auth.GenerateJWT(u.Username,u.ID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	u.JWTToken = jwtToken
	return nil
}

