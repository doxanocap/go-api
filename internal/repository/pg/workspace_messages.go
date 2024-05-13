package pg

import (
	"auth-api/internal/models"
	"context"
	"github.com/doxanocap/pkg/errs"
	"github.com/jmoiron/sqlx"
)

type WorkspaceMessages struct {
	db *sqlx.DB
}

func InitWorkspaceMessagesRepository(db *sqlx.DB) *WorkspaceMessages {
	return &WorkspaceMessages{
		db: db,
	}
}

func (repo *WorkspaceMessages) Create(ctx context.Context, wm *models.WorkspaceMessage) (err error) {
	defer errs.WrapIfErr("repo.workspace_messages.Create", &err)

	err = repo.db.QueryRowxContext(ctx,
		`insert into workspace_messages
			(workspace_idref, user_idref, content) 
			values ($1,$2,$3) returning message_id`,
		wm.WorkspaceIDRef, wm.UserIDRef, wm.Content).
		Scan(&wm.MessageID)
	return
}

func (repo *WorkspaceMessages) GetAllByWorkspaceID(ctx context.Context, workspaceID int64) (ws []models.WorkspaceMessage, err error) {
	defer errs.WrapIfErr("repo.workspace_messages.GetAllByWorkspaceID", &err)

	ws = []models.WorkspaceMessage{}
	err = repo.db.SelectContext(ctx, ws,
		`select * from workspace_messages where workspace_idref = $1 order by ts`, workspaceID)
	return
}

func (repo *WorkspaceMessages) DeleteByWorkspaceID(ctx context.Context, id int64) (err error) {
	defer errs.WrapIfErr("repository.session.DeleteByWorkspaceID", &err)

	_, err = repo.db.ExecContext(ctx, `delete from workspace_messages where workspace_idref = $1`, id)
	return
}
