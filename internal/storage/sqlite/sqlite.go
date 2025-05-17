package sqlite

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/azan-boss/posty/internal/config"
	"github.com/azan-boss/posty/internal/types"
	"github.com/azan-boss/posty/internal/utils/method"
	"github.com/google/uuid"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Sqlite struct {
	db *gorm.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	gormDB, err := gorm.Open(sqliteDriver.Open(cfg.Storage.Database), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err)
	}
	err = gormDB.AutoMigrate(&types.Post{}, &types.User{}, &types.ChatRoom{}, &types.Message{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %s", err)
	}

	return &Sqlite{db: gormDB}, nil
}

func (s *Sqlite) RegisterUser(user *types.User) (uint, error) {
	if s.db == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}

	// Check if user already exists
	var existingUser types.User
	if result := s.db.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser); result.Error == nil {
		return 0, fmt.Errorf("user with username or email already exists")
	}

	result := s.db.Create(user)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to register user: %v", result.Error)
	}

	return user.ID, nil
}

func (s *Sqlite) Login(user *types.User) (types.User, error) {
	var dbUser types.User
	result := s.db.Where("username = ?", user.Username).First(&dbUser)
	if result.Error != nil {
		return types.User{}, fmt.Errorf("user not found")
	}
	if err := method.CompareHashAndPassword(dbUser.Password, user.Password); err != nil {
		return types.User{}, fmt.Errorf("invalid password")
	}
	return dbUser, nil
}

func (s *Sqlite) CreatePost(post *types.Post) error {
	// If no ID is provided, generate a unique one
	if post.ID == "" {
		post.ID = uuid.New().String()
	}

	// Use GORM's Create method which handles ID generation
	result := s.db.Create(post)
	if result.Error != nil {
		return fmt.Errorf("failed to create post: %s", result.Error)
	}

	return nil
}

func (s *Sqlite) GetUser(id uint) (types.User, error) {

	var user types.User
	result := s.db.Where("id = ?", id).Find(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return types.User{}, fmt.Errorf("user not found")
		} else {
			return types.User{}, fmt.Errorf("db error %s", string(result.Error.Error()))
		}
	}
	return user, result.Error

}

func (s *Sqlite) CreateChatRoom(chatroom *types.ChatRoom) (uint, error) {

	result := s.db.Create(chatroom)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create chatroom: %s", result.Error)
	}
	return chatroom.ID, nil
}

func (s *Sqlite) GetChatRoom(slug string) (types.ChatRoom, error) {
	var chatroom types.ChatRoom
	result := s.db.Where("slug = ?", slug).First(&chatroom)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return types.ChatRoom{}, fmt.Errorf("chatroom not found")
		} else {
			return types.ChatRoom{}, fmt.Errorf("db error %s", string(result.Error.Error()))
		}
	}
	return chatroom, nil
}

func (s *Sqlite) GetUseChatRoom(slug, id string) (types.ChatRoom, error) {
	var chatroom types.ChatRoom
	result := s.db.Preload("Members", "id = ?", id).Where("slug = ?", slug).First(&chatroom)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return types.ChatRoom{}, fmt.Errorf("chatroom not found")
		} else {
			return types.ChatRoom{}, fmt.Errorf("db error %s", string(result.Error.Error()))
		}
	}
	return chatroom, nil
}

func (s *Sqlite) FindUserChatRoomOrJoin(slug string, id string) (types.ChatRoom, error) {
	// First check if the chatroom exists and preload the user if they're a member
	var chatroom types.ChatRoom
	result := s.db.Preload("Members", "id = ?", id).Where("slug = ?", slug).First(&chatroom)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return types.ChatRoom{}, fmt.Errorf("chatroom not found")
		} else {
			return types.ChatRoom{}, fmt.Errorf("db error %s", result.Error.Error())
		}
	}

	// If Members is empty, the user is not in the chatroom yet, so add them
	if len(chatroom.Members) == 0 {
		userID, err := strconv.Atoi(id)
		if err != nil {
			return types.ChatRoom{}, fmt.Errorf("invalid user ID: %s", err.Error())
		}

		user, err := s.GetUser(uint(userID))
		if err != nil {
			return types.ChatRoom{}, fmt.Errorf("user not found: %s", err.Error())
		}

		// Add the user to the chatroom
		err = s.db.Model(&chatroom).Association("Members").Append(&user)
		if err != nil {
			return types.ChatRoom{}, fmt.Errorf("failed to add user to chatroom: %s", err.Error())
		}
	}

	return chatroom, nil
}




func (s *Sqlite)UpdateStatus(user *types.User) error {
		result:=s.db.Model(&user).Updates(types.User{Status:user.Status,LastActive: time.Now()})

		if result.Error!=nil{
			return fmt.Errorf("failed to update user status: %s", result.Error.Error())
		}
		return nil
	}



func (s *Sqlite)GetMessageByChatRoomId(id string)([]types.Message ,error)  {
	var messages []types.Message
	result:=s.db.Where("chat_room_id = ?",id).Find(&messages)
	fmt.Println(result)
	if result.Error!= nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			slog.Error("message not found")
			return nil, fmt.Errorf("message not found")
		}else {
			slog.Error("db error", "error", result.Error)
			return nil, fmt.Errorf("db error %s", string(result.Error.Error()))
		}
	}else{
		slog.Info("Message history successfully loaded", "amount", result.RowsAffected)
		return messages,nil
	}
}






func (s *Sqlite)CreateMessage(message *types.Message) error  {
	result:=s.db.Create(message)
	if result.Error!= nil {
		return fmt.Errorf("failed to create message: %s", result.Error.Error())
	}
	return nil
}