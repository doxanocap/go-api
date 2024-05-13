package service

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

const (
	storageTicker = 5 * time.Second
)

type Workspace struct {
	log     *zap.Logger
	manager interfaces.IManager
}

func InitWorkspaceService(manager interfaces.IManager, log *zap.Logger) *Workspace {
	return &Workspace{
		log:     log,
		manager: manager,
	}
}

func (w *Workspace) Create(ctx context.Context, editorID string) (*models.Workspace, error) {
	workspace := &models.Workspace{
		LastEditorID: editorID,
		IDCode:       uuid.NewString(),
	}

	err := w.manager.Repository().Workspace().Create(ctx, workspace)
	if err != nil {
		return nil, err
	}

	err = w.manager.Repository().WorkspaceUsers().Create(ctx, &models.WorkspaceUser{
		WorkspaceIDRef: workspace.ID,
		UserIDRef:      editorID,
		IsPresent:      true,
	})
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (w *Workspace) GetByIDCode(ctx context.Context, workspaceIDCode string) (*models.Workspace, error) {
	workspace, err := w.manager.Repository().Workspace().FindByIDCode(ctx, workspaceIDCode)
	if err != nil {
		return nil, err
	}

	cacheWorkspace := w.manager.Service().Websocket().GetByIDCode(workspace.IDCode)
	if cacheWorkspace == nil {
		return nil, models.ErrWorkspaceNotFound
	}

	if cacheWorkspace.UpdatedAt != workspace.UpdatedAt {
		err = w.manager.Repository().Workspace().UpdateBodyByID(ctx, cacheWorkspace)
		if err != nil {
			return nil, err
		}
	}
	return cacheWorkspace, nil
}

func (w *Workspace) GetAllByUserID(ctx context.Context, userID string) ([]models.Workspace, error) {
	workspaces, err := w.manager.Repository().WorkspaceUsers().GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (w *Workspace) HandleNewClient(ctx context.Context,
	conn *websocket.Conn, clientID, workspaceIDCode string) (err error) {
	ws := w.manager.Service().Websocket()

	defer func() {
		if err != nil {
			ws.WriteConnErr(conn, err)
		}
	}()

	workspace, err := w.manager.Repository().Workspace().FindByIDCode(ctx, workspaceIDCode)
	if err != nil {
		return err
	}
	if workspace == nil {
		return models.ErrWorkspaceNotFound
	}

	_, err = w.manager.Processor().APIs().AuthAPI().GetUserByID(ctx, clientID)
	if err != nil {
		return err
	}

	workspace = ws.GetByIDCode(workspace.IDCode)
	if workspace == nil {
		return models.ErrWorkspaceNotFound
	}

	wUsers, err := w.manager.Repository().WorkspaceUsers().GetAllByWorkspaceID(ctx, workspace.ID)
	if err != nil {
		return
	}

	found := false
	for _, wu := range wUsers {
		if wu.UserIDRef == clientID {
			err = w.manager.Repository().WorkspaceUsers().SetUserPresent(ctx, workspace.ID, clientID)
			if err != nil {
				return err
			}
			found = true
		}
	}

	if !found {
		err = w.manager.Repository().WorkspaceUsers().Create(ctx, &models.WorkspaceUser{
			WorkspaceIDRef: workspace.ID,
			UserIDRef:      clientID,
			IsPresent:      true,
		})
		if err != nil {
			return err
		}
	}

	sig := make(chan struct{})

	go w.asyncUpdate(ctx, sig, workspace.IDCode)
	err = ws.HandleNewClient(clientID, workspace.IDCode, conn)

	sig <- struct{}{}
	return err
}

func (w *Workspace) asyncUpdate(ctx context.Context, sig chan struct{}, workspaceIDCode string) {
	var err error
	defer func() {
		if err != nil {
			w.log.Info(fmt.Sprintf("workspace.asyncUpdate: %s", err))
		}
	}()

	ticker := time.NewTicker(storageTicker)

	workspace := w.manager.Service().Websocket().GetByIDCode(workspaceIDCode)
	lastUpdatedAt := time.Time{}
	if workspace.UpdatedAt != nil {
		lastUpdatedAt = *workspace.UpdatedAt
	}

	for {
		select {
		case <-ticker.C:
			updWorkspace := w.manager.Service().Websocket().GetByIDCode(workspaceIDCode)
			if updWorkspace == nil {
				err = models.ErrWorkspaceNotFound
				return
			}
			if lastUpdatedAt == *workspace.UpdatedAt {
				continue
			}

			lastUpdatedAt = *workspace.UpdatedAt
			err = w.manager.Repository().Workspace().UpdateBodyByID(ctx, updWorkspace)
			if err != nil {
				return
			}
		case <-sig:
			err = w.manager.Repository().Workspace().UpdateBodyByID(ctx, workspace)
			return
		}
	}

}
