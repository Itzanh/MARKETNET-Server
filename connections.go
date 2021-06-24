package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Id            string    `json:"id"`
	Address       string    `json:"address"`
	User          int16     `json:"user"`
	DateConnected time.Time `json:"dateConnected"`
	ws            *websocket.Conn
}

func (c *Connection) addConnection() {
	c.Id = uuid.New().String()
	c.DateConnected = time.Now()

	connections = append(connections, *c)
}

func (c *Connection) deleteConnection() {
	for i := 0; i < len(connections); i++ {
		if connections[i].Id == c.Id {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

func disconnectConnection(id string) bool {
	for i := 0; i < len(connections); i++ {
		if connections[i].Id == id {
			connections[i].ws.Close()
			return true
		}
	}
	return false
}

type ConnectionWeb struct {
	Id            string    `json:"id"`
	Address       string    `json:"address"`
	User          string    `json:"user"`
	DateConnected time.Time `json:"dateConnected"`
}

func getConnections() []ConnectionWeb {
	conn := make([]ConnectionWeb, 0)

	for i := 0; i < len(connections); i++ {
		var userName string

		sqlStatement := `SELECT username FROM "user" WHERE id=$1`
		row := db.QueryRow(sqlStatement, connections[i].User)
		row.Scan(&userName)

		conn = append(conn, ConnectionWeb{Id: connections[i].Id, Address: connections[i].Address, User: userName, DateConnected: connections[i].DateConnected})
	}

	return conn
}
