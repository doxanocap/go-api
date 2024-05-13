package pg

import (
	"auth-api/internal/models"
	"context"
	"github.com/doxanocap/pkg/errs"
	"github.com/jmoiron/sqlx"
)

type WorkspaceUsers struct {
	db *sqlx.DB
}

func InitWorkspaceUsersRepository(db *sqlx.DB) *WorkspaceUsers {
	return &WorkspaceUsers{
		db: db,
	}
}

func (repo *WorkspaceUsers) Create(ctx context.Context, wu *models.WorkspaceUser) (err error) {
	defer errs.WrapIfErr("repo.workspace_users.Create", &err)

	err = repo.db.QueryRowxContext(ctx,
		`insert into workspace_users
			(workspace_idref, user_idref, is_present) 
			values ($1,$2,$3)`,
		wu.WorkspaceIDRef, wu.UserIDRef, wu.IsPresent).Err()
	return
}

func (repo *WorkspaceUsers) GetAll(ctx context.Context) (ws []models.WorkspaceUser, err error) {
	defer errs.WrapIfErr("repo.workspace_users.GetAll", &err)

	ws = []models.WorkspaceUser{}
	err = repo.db.SelectContext(ctx, &ws, `select * from workspace_users`)
	return
}

func (repo *WorkspaceUsers) GetAllByWorkspaceID(ctx context.Context, workspaceID int64) (wu []models.WorkspaceUser, err error) {
	defer errs.WrapIfErr("repo.workspace_users.GetAllByWorkspaceID", &err)

	wu = []models.WorkspaceUser{}
	err = repo.db.SelectContext(ctx, &wu,
		`select * from workspace_users where workspace_idref = $1`, workspaceID)
	return
}

func (repo *WorkspaceUsers) GetAllByUserID(ctx context.Context, userID string) (w []models.Workspace, err error) {
	defer errs.WrapIfErr("repo.workspace_users.GetAllByUserID", &err)

	w = []models.Workspace{}
	err = repo.db.SelectContext(ctx, &w,
		`select 
			w.workspace_id, w.workspace_idcode,
			w.last_editor_id, w.body, w.updated_at, w.created_at
		from workspace_users wu
		join workspaces w on wu.workspace_idref = w.workspace_id
		where wu.user_idref = $1`, userID)
	return
}

func (repo *WorkspaceUsers) SetUserPresent(ctx context.Context, workspaceID int64, userID string) (err error) {
	defer errs.WrapIfErr("repository.workspace_users.SetUserPresent", &err)

	_, err = repo.db.ExecContext(ctx,
		`update workspace_users 
		set is_present = true 
		where workspace_idref = $1 and user_idref = $2`, workspaceID, userID)
	return
}

func (repo *WorkspaceUsers) SetUserAbsent(ctx context.Context, workspaceID int64, userID string) (err error) {
	defer errs.WrapIfErr("repository.workspace_users.SetUserAbsent", &err)

	_, err = repo.db.ExecContext(ctx,
		`update workspace_users 
		set is_present = false 
		where workspace_idref = $1 and user_idref = $2`, workspaceID, userID)
	return
}

func (repo *WorkspaceUsers) DeleteByID(ctx context.Context, id int64) (err error) {
	defer errs.WrapIfErr("repository.session.DeleteByID", &err)

	_, err = repo.db.ExecContext(ctx,
		`delete from workspace_users where workspace_idref = $1`, id)
	return
}
