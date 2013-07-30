package Server

import (
	"fmt"
	m "github.com/jashper/texas-coinem/Message"
)

type Parser struct {
	connection *Connection
	context    *ServerContext
}

func (this *Parser) Init(connection *Connection, context *ServerContext) {
	this.connection = connection
	this.context = context
}

func (this *Parser) getStrings(paramCount int) []string {
	data := make([]string, paramCount)

	var length [1]byte
	for i := 0; i < paramCount; i++ {
		this.connection.Read(length[:])
		temp := make([]byte, int(length[0]))
		this.connection.Read(temp)
		data[i] = string(temp)
	}

	return data
}

func (this *Parser) appendStrings(message []byte, params []string) {
	for i := 0; i < len(params); i++ {
		str := params[i]
		if len(str) > 255 {
			fmt.Println("CRITICAL : Tried to send a string message greater than 255 characters")
			return
		}
		message = append(message, byte(len(str)))
		message = append(message, []byte(str)...)
	}
}

func (this *Parser) Message(mType m.ServerMessage) {
	if mType == m.SM_LOGIN_REGISTER {
		data := this.getStrings(2)
		username := data[0]
		password := data[1]

		err := RegisterUser(username, password, this.context)
		var message m.ClientMessage
		if err == nil {
			message = m.CM_LOGIN_REGISTER_SUCCESS
		} else {
			message = m.CM_LOGIN_REGISTER_FAILURE
		}

		toSend := []byte{byte(message)}
		this.connection.Write(toSend)
	}
}
