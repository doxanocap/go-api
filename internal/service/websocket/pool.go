package websocket

import (
	"auth-api/internal/models"
	"auth-api/internal/pkg/tools"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Pool struct {
	mu  *sync.RWMutex
	log *zap.Logger

	users          map[string]chan []byte
	workspaces     map[string]*models.Workspace
	workspaceUsers map[string][]string

	broadcast chan models.Message
}

func InitPool(log *zap.Logger) *Pool {
	return &Pool{
		log: log,
		mu:  &sync.RWMutex{},

		users:          make(map[string]chan []byte),
		workspaces:     map[string]*models.Workspace{},
		workspaceUsers: map[string][]string{},

		broadcast: make(chan models.Message),
	}
}

// Handle handles all messages and writes to users
// to whom message should be delivered
func (p *Pool) Handle() {
	for {
		message, ok := <-p.broadcast
		if !ok {
			return
		}

		p.handleMessage(&message)
	}
}

func (p *Pool) Register(userID, workspaceID string, ch chan []byte) (err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, ok := p.workspaces[workspaceID]
	if !ok {
		return models.ErrWorkspaceNotFound
	}
	wClients, ok := p.workspaceUsers[workspaceID]
	if !ok {
		return models.ErrWorkspaceNotFound
	}

	found := false
	p.users[userID] = ch
	for _, cID := range wClients {
		if cID == userID {
			found = true
			break
		}
	}

	if !found {
		wClients = append(wClients, userID)
	}
	return nil
}

func (p *Pool) Unregister(userID string) {
	p.log.Info(fmt.Sprintf("unregistered: user_id: %s", userID))

	p.mu.RLock()
	defer p.mu.RUnlock()
	delete(p.users, userID)
}

// Send sends message to pool
func (p *Pool) Send(message *models.Message) {
	if message.ErrorMessage == "" {
		p.log.Info(fmt.Sprintf("send: user_id: %s | workspace_id: %d | body: %s",
			message.UserID, message.WorkspaceID, message.Body))
	}
	p.broadcast <- *message
}

func (p *Pool) GetByID(workspaceID string) *models.Workspace {
	p.mu.RLock()
	defer p.mu.RUnlock()

	v, ok := p.workspaces[workspaceID]
	if !ok {
		return nil
	}
	return v
}

func (p *Pool) SetWorkspaces(m map[string]*models.Workspace) {
	p.workspaces = m
}

func (p *Pool) SetWorkspaceClients(m map[string][]string) {
	p.workspaceUsers = m
}

func (p *Pool) handleMessage(msg *models.Message) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// set message into body
	w := p.workspaces[msg.WorkspaceID]
	w.Body = msg.Body
	w.UpdatedAt = tools.GetPtr(time.Now())
	w.LastEditorID = msg.UserID
	p.workspaces[msg.WorkspaceID] = w

	rawMessage, _ := json.Marshal(msg)
	for _, userID := range p.workspaceUsers[msg.WorkspaceID] {
		incMessagesCh, ok := p.users[userID]
		// if user absent || self messaging
		if !ok || userID == msg.UserID {
			continue
		}
		incMessagesCh <- rawMessage
	}
}
