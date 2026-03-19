package dbmigrations

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DbUpgrade interface {
	InstallMode(context.Context, []string) error
	ExecuteQueries(context.Context, func() error) error
	GetDBVersion(context.Context) (int, error)
	UpdateDBVersion(context.Context, *DBVersionUpgrade) error
}

type dbupgrader struct {
	// logger loggger
	db        *sqlx.DB
	tableName string
}

func NewDbUpgrader(db *sqlx.DB) DbUpgrade {
	return &dbupgrader{
		db: db,
	}
}

func (svc *dbupgrader) setup(ctx context.Context) error {
	// create DB version table
	for i := range installModequeries {
		_, err := svc.db.ExecContext(ctx, installModequeries[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *dbupgrader) InstallMode(ctx context.Context, queries []string) error {

	if err := svc.setup(ctx); err != nil {
		return err
	}

	for i := range queries {
		// run queries
		_, err := svc.db.ExecContext(ctx, queries[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *dbupgrader) ExecuteQueries(ctx context.Context, upgrade func() error) error {
	return upgrade()
}

func (svc *dbupgrader) GetDBVersion(ctx context.Context) (int, error) {
	var version int
	err := svc.db.GetContext(ctx, &version, "SELECT version FROM db_version LIMIT 1")
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return 0, err
	}
	return version, nil
}

func (svc *dbupgrader) UpdateDBVersion(ctx context.Context, dbversion *DBVersionUpgrade) error {

	q := sq.Update(svc.tableName).Set(Version, dbversion.Version)
	query, args, err := q.ToSql()
	if err != nil {
		return err
	}
	_, err = svc.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// insert history table as well
	q2 := sq.Insert(DBVersionHistoryTableName).Columns(DBVersionHistInsertColumns...).
		Values(dbversion.Name, dbversion.Version, dbversion.Description, dbversion.UpdatedAt)
	query2, args2, err := q2.ToSql()
	if err != nil {
		return err
	}
	_, err = svc.db.ExecContext(ctx, query2, args2...)
	if err != nil {
		return err
	}
	return nil
}
