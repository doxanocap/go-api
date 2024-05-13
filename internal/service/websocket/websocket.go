package websocket

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Websocket struct {
	log     *zap.Logger
	config  *models.Config
	manager interfaces.IManager

	poolProcessor interfaces.IWebsocketPoolService
}

func NewWebsocketService(manager interfaces.IManager, config *models.Config, log *zap.Logger) *Websocket {
	w := &Websocket{
		log:     log,
		config:  config,
		manager: manager,
	}
	if err := w.initPool(); err != nil {
		log.Fatal(err.Error())
	}
	return w
}

func (w *Websocket) HandleNewClient(userID, workspaceID string, conn *websocket.Conn) error {
	c := &Client{
		userID:      userID,
		workspaceID: workspaceID,

		errCh:         make(chan error),
		incMessagesCh: make(chan []byte),
		poolProcessor: w.poolProcessor,
		conn:          conn,
		log:           w.log,
	}

	if err := c.registerClient(); err != nil {
		return err
	}

	go c.Reader()
	go c.Writer()

	err := <-c.errCh
	close(c.errCh)
	close(c.incMessagesCh)
	return err
}

func (w *Websocket) GetByIDCode(workspaceID string) *models.Workspace {
	return w.poolProcessor.GetByID(workspaceID)
}

func (w *Websocket) WriteConnErr(conn *websocket.Conn, err error) {
	log := w.log.Named("[ERROR]")
	if errors.Is(err, models.ErrConnGracefullyClosed) {
		return
	}

	log.Error(err.Error())
	message := models.Message{}

	httpErr := errs.UnmarshalError(err)
	if httpErr.StatusCode != 0 {
		message.ErrorMessage = httpErr.Error()
	} else {
		message.ErrorMessage = models.HttpInternalServerError.Error()
	}

	rawErrMsg, _ := json.Marshal(message)
	if err = conn.WriteMessage(websocket.TextMessage, rawErrMsg); err != nil {
		w.log.Error(fmt.Sprintf("writeConnErr: msg: %s | err: %s", string(rawErrMsg), err))
	}
}

func (w *Websocket) initPool() (err error) {
	defer errs.WrapIfErr("initPool", &err)

	ctx := context.Background()
	workspaces, err := w.manager.Repository().Workspace().GetAll(ctx)
	if err != nil {
		return err
	}

	workspacesMap := map[string]*models.Workspace{}
	idToIDCode := map[int64]string{}
	for _, ws := range workspaces {
		workspacesMap[ws.IDCode] = &ws
		idToIDCode[ws.ID] = ws.IDCode
	}

	workspaceUsers, err := w.manager.Repository().WorkspaceUsers().GetAll(ctx)
	if err != nil {
		return err
	}

	workspaceUsersMap := map[string][]string{}
	for _, wsu := range workspaceUsers {
		idcode := idToIDCode[wsu.WorkspaceIDRef]

		_, ok := workspaceUsersMap[idcode]
		if !ok {
			workspaceUsersMap[idcode] = []string{wsu.UserIDRef}
			continue
		}

		workspaceUsersMap[idcode] = append(workspaceUsersMap[idcode], wsu.UserIDRef)
	}

	w.poolProcessor = InitPool(w.log)
	fmt.Println(workspacesMap)
	w.poolProcessor.SetWorkspaces(workspacesMap)
	w.poolProcessor.SetWorkspaceClients(workspaceUsersMap)
	go w.poolProcessor.Handle()
	return nil
}
