package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Id            string `json:"id"`
	Address       string `json:"address"`
	User          int32  `json:"user"`
	enterprise    int32
	DateConnected time.Time `json:"dateConnected"`
	ws            *websocket.Conn
}

func (c *Connection) addConnection() {
	connectionsMutex.Lock()
	c.Id = uuid.New().String()
	c.DateConnected = time.Now()

	connections = append(connections, *c)
	connectionsMutex.Unlock()
}

func (c *Connection) deleteConnection() {
	connectionsMutex.Lock()
	for i := 0; i < len(connections); i++ {
		if connections[i].Id == c.Id {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}
	connectionsMutex.Unlock()
}

func disconnectConnection(id string, enterpriseId int32) bool {
	connectionsMutex.Lock()
	for i := 0; i < len(connections); i++ {
		if connections[i].Id == id && connections[i].enterprise == enterpriseId {
			connections[i].ws.Close()
			connectionsMutex.Unlock()
			return true
		}
	}
	connectionsMutex.Unlock()
	return false
}

func disconnectAllConnections(enterpriseId int32) {
	connectionsMutex.Lock()
	for i := 0; i < len(connections); i++ {
		if connections[i].enterprise == enterpriseId {
			connections[i].ws.Close()
		}
	}
	connectionsMutex.Unlock()
}

type ConnectionWeb struct {
	Id            string    `json:"id"`
	Address       string    `json:"address"`
	User          string    `json:"user"`
	DateConnected time.Time `json:"dateConnected"`
}

func getConnections(enterpriseId int32) []ConnectionWeb {
	conn := make([]ConnectionWeb, 0)

	connectionsMutex.Lock()
	for i := 0; i < len(connections); i++ {
		if connections[i].enterprise != enterpriseId {
			continue
		}
		var userName string
		// get a user's username from the user id in the connection using dbOrm
		dbOrm.Model(&User{}).Where("id = ?", connections[i].User).Pluck("name", &userName)
		conn = append(conn, ConnectionWeb{Id: connections[i].Id, Address: connections[i].Address, User: userName, DateConnected: connections[i].DateConnected})
	}

	connectionsMutex.Unlock()
	return conn
}
