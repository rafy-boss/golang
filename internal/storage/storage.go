package storage

import "github.com/azan-boss/posty/internal/types"

type Storage interface {
	CreatePost(post *types.Post) error
	RegisterUser(user *types.User) (uint, error)
	Login(user *types.User) (types.User, error)
	GetUser(id uint) (types.User, error)
	CreateChatRoom(chatroom *types.ChatRoom) (uint,error)
	GetChatRoom(slug string) (types.ChatRoom, error)
	GetUseChatRoom(slug string ,id string) (types.ChatRoom, error)
	FindUserChatRoomOrJoin(slug string ,id string) (types.ChatRoom, error)
	UpdateStatus(user *types.User) error
	GetMessageByChatRoomId(id string)([]types.Message ,error)
	CreateMessage(message *types.Message) error
}