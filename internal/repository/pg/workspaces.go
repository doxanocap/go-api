package pg

import (
	"auth-api/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Workspace struct {
	db *sqlx.DB
}

func InitWorkspaceRepository(db *sqlx.DB) *Workspace {
	return &Workspace{
		db: db,
	}
}

func (repo *Workspace) Create(ctx context.Context, w *models.Workspace) (err error) {
	defer errs.WrapIfErr("repo.workspace.Create", &err)

	err = repo.db.QueryRowxContext(ctx,
		`insert into workspaces (last_editor_id, workspace_idcode) values ($1, $2) 
		returning workspace_id, created_at`, w.LastEditorID, w.IDCode).
		Scan(&w.ID, &w.CreatedAt)
	return
}

func (repo *Workspace) GetAll(ctx context.Context) (ws []models.Workspace, err error) {
	defer errs.WrapIfErr("repo.workspace.GetAll", &err)

	ws = []models.Workspace{}
	err = repo.db.SelectContext(ctx, &ws, `select * from workspaces`)
	return
}

func (repo *Workspace) FindByID(ctx context.Context, id int64) (w *models.Workspace, err error) {
	defer errs.WrapIfErr("repo.workspace.FindByID", &err)

	w = &models.Workspace{}
	err = repo.db.GetContext(ctx, w, `select * from workspaces where workspace_id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	return
}

func (repo *Workspace) FindByIDCode(ctx context.Context, idcode string) (w *models.Workspace, err error) {
	defer errs.WrapIfErr("repo.workspace.FindByIDCode", &err)

	w = &models.Workspace{}
	err = repo.db.GetContext(ctx, w, `select * from workspaces where workspace_idcode = $1`, idcode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	return
}

func (repo *Workspace) UpdateByID(ctx context.Context, w *models.Workspace) (err error) {
	defer errs.WrapIfErr("repo.workspace.UpdateByID", &err)
	if w.ID == 0 {
		return
	}

	var sb strings.Builder
	baseStr := "update workspaces set"
	sb.WriteString(baseStr)

	if w.UpdatedAt == nil {
		sb.WriteString(fmt.Sprintf("updated_at = %s,", w.UpdatedAt.String()))
	}
	if w.Body == "" {
		sb.WriteString(fmt.Sprintf("body = %s,", w.Body))
	}
	if w.LastEditorID == "" {
		sb.WriteString(fmt.Sprintf("last_editor_id = %s,", w.LastEditorID))
	}

	if sb.Len() == len(baseStr) {
		return
	}

	query := sb.String()
	query = query[:len(query)-1]
	sb.Reset()

	sb.WriteString(query)
	sb.WriteString(fmt.Sprintf(" where workspace_id = %d", w.ID))
	_, err = repo.db.ExecContext(ctx, sb.String())
	return
}

func (repo *Workspace) UpdateBodyByID(ctx context.Context, w *models.Workspace) (err error) {
	defer errs.WrapIfErr("repo.workspace.UpdateBodyByID", &err)

	_, err = repo.db.ExecContext(ctx,
		`update workspaces set body = $1, updated_at = $2, last_editor_id = $3 where workspace_id = $4`,
		w.Body, w.UpdatedAt, w.LastEditorID, w.ID)
	return
}

func (repo *Workspace) DeleteByID(ctx context.Context, id int64) (err error) {
	defer errs.WrapIfErr("repo.workspace.DeleteByID", &err)

	_, err = repo.db.ExecContext(ctx, `delete from workspaces where workspace_id = $1`, id)
	return
}
