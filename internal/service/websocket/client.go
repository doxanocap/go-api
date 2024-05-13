package websocket

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"encoding/json"
	"errors"
	"github.com/doxanocap/pkg/errs"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

const (
	pingPeriod         = 54 * time.Second
	pongWait           = 60 * time.Second
	awaitWriteDuration = 10 * time.Second
)

type Client struct {
	userID      string
	workspaceID string

	// errCh channel with errors | close only from outside!
	errCh         chan error
	incMessagesCh chan []byte

	poolProcessor interfaces.IWebsocketPoolService
	conn          *websocket.Conn

	log *zap.Logger
}

// Reader reads all messages incoming from WS poolProcessor
func (c *Client) Reader() {
	var message *models.Message
	var err error

	defer func() {
		httpErr := errs.UnmarshalError(err)
		if httpErr.StatusCode == 0 {
			err = errs.Wrap("ws.Reader", err)
		}
		c.writeError(err)
		c.unregisterClient()
	}()

	err = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}

	c.conn.SetPingHandler(func(appData string) error {
		return nil
	})

	c.conn.SetPongHandler(func(appData string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		message, err = c.readMessage()
		if err != nil {
			return
		}
		if message.UserID != c.userID {
			err = models.ErrMsgFromDifferentUser
			return
		}

		if message != nil {
			c.poolProcessor.Send(message)
		}
	}
}

// Writer writes message into WS poolProcessor
func (c *Client) Writer() {
	ticker := time.NewTicker(pingPeriod)
	var err error

	defer func() {
		ticker.Stop()
		if err != nil {
			httpErr := errs.UnmarshalError(err)
			if httpErr.StatusCode == 0 {
				err = errs.Wrap("ws.Writer", err)
			}
			c.writeError(err)
		}
	}()

	for {
		select {
		case message, ok := <-c.incMessagesCh:
			if !ok {
				return
			}
			if err = c.setWriteDeadline(); err != nil {
				return
			}

			if err = c.writeMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err = c.setWriteDeadline(); err != nil {
				return
			}

			if err = c.writeMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readMessage() (*models.Message, error) {
	messageType, body, err := c.conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			return nil, models.ErrConnGracefullyClosed
		}
		return nil, errs.Wrap("readMessage", err)
	}

	switch messageType {
	case websocket.PingMessage:
		return nil, nil
	case websocket.CloseMessage:
		return nil, models.ErrConnGracefullyClosed
	}

	message := &models.Message{
		UserID:      c.userID,
		WorkspaceID: c.workspaceID,
		MsgType:     consts.SendMessage,
	}

	if err = json.Unmarshal(body, message); err != nil {
		return nil, errs.Wrap("readMessage: unmarshal", err)
	}

	return message, nil
}

func (c *Client) setWriteDeadline() error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(awaitWriteDuration)); err != nil {
		return errs.Wrap("setWriteDeadline", err)
	}
	return nil
}

func (c *Client) writeMessage(messageType int, message []byte) (err error) {
	defer func() {
		if err != nil {
			err = errs.Wrap("writeMessage", err)
		}
	}()
	if err = c.conn.WriteMessage(messageType, message); err != nil {
		if errors.Is(err, websocket.ErrCloseSent) {
			return nil
		}

		switch messageType {
		case websocket.PingMessage, websocket.PongMessage:
			err = errs.Wrap("ping/pong", err)
		case websocket.CloseMessage:
			err = errs.Wrap("close", err)
		}
		return err
	}

	return nil
}

func (c *Client) writeError(err error) {
	c.log.Named("[ERROR]").Error(err.Error())
	c.errCh <- err
}

func (c *Client) unregisterClient() {
	c.poolProcessor.Unregister(c.userID)
}

func (c *Client) registerClient() error {
	err := c.poolProcessor.Register(c.userID, c.workspaceID, c.incMessagesCh)
	if err != nil {
		return err
	}
	return nil
}
