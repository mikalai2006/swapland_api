package service

import (
	"fmt"

	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
)

// Hub — это структура, содержащая всех клиентов и отправляемые им сообщения.
type Hub struct {
	// Зарегистрированные клиенты.
	clients map[string]map[*Client]bool
	// Незарегистрированные клиенты.
	unregister chan *Client
	// Регистрация заявок от клиентов.
	register chan *Client
	// Входящие сообщения от клиентов.
	broadcast chan domain.Message
}

// // Message struct to hold message data
// type Message struct {
// 	Type      string `json:"type"`
// 	Sender    string `json:"sender"`
// 	Recipient string `json:"recipient"`
// 	Content   string `json:"content"`
// 	ID        string `json:"id"`
// }

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		unregister: make(chan *Client),
		register:   make(chan *Client),
		broadcast:  make(chan domain.Message),
	}
}

// Основная функция для запуска хаба
func (h *Hub) Run() {
	for {
		select {
		// Register a client.
		case client := <-h.register:
			h.RegisterNewClient(client)
			// Unregister a client.
		case client := <-h.unregister:
			h.RemoveClient(client)
			// Broadcast a message to all clients.
		case message := <-h.broadcast:

			//Check if the message is a type of "message"
			h.HandleMessage(message)

		}
	}
}

// функция проверяет, существует ли комната, и если нет, создайте ее и добавьте в нее клиента
func (h *Hub) RegisterNewClient(client *Client) {
	connections := h.clients[client.RoomId]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.clients[client.RoomId] = connections
	}
	h.clients[client.RoomId][client] = true

	// user, err :=
	client.Services.User.UpdateUser(client.UserId, &model.User{Online: true})
	// if err == nil {
	// 	h.HandleMessage(domain.Message{Type: "message", Sender: client.UserId, Recipient: "user2", Content: user, ID: "room1", Service: "user"})
	// }

	fmt.Println("Size of clients: ", len(h.clients[client.RoomId]))
}

// function to remvoe client from room
func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.clients[client.RoomId]; ok {
		// user, err :=
		client.Services.User.UpdateUser(client.UserId, &model.User{Online: false})
		// if err == nil {
		// 	h.HandleMessage(domain.Message{Type: "message", Sender: client.UserId, Recipient: "user2", Content: user, ID: "room1", Service: "user"})
		// }

		delete(h.clients[client.RoomId], client)
		close(client.send)
		fmt.Println("Removed client", client.UserId)
	}
}

// function to handle message based on type of message
func (h *Hub) HandleMessage(message domain.Message) {

	//Check if the message is a type of "message"
	if message.Type == "message" {
		clients := h.clients[message.ID]
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients[message.ID], client)
			}
		}
	}

	//Check if the message is a type of "notification"
	if message.Type == "notification" {
		fmt.Println("Notification: ", message.Content)
		clients := h.clients[message.Recipient]
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients[message.Recipient], client)
			}
		}
	}

}
