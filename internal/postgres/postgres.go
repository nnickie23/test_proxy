package postgresdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	//pq driver must be imported unnamed
	"github.com/nnickie23/test_proxy/internal/configs"
	"github.com/nnickie23/test_proxy/internal/entities/errors"
	"github.com/nnickie23/test_proxy/internal/logger"
	_ "github.com/jackc/pgx/v4/stdlib" // driver
	"github.com/jmoiron/sqlx"
)

const ErrMsg = "PG-error"
const TransactionCtxKey = "pg_transaction"

type PostgresDb interface {
	Close()
	ContextWithTransaction(ctx context.Context) (context.Context, error)
	CommitContextTransaction(ctx context.Context) error
	RollbackContextTransaction(ctx context.Context)
	InsertReturnId(result interface{}, context context.Context, table string, fields Fields) error
	Get(result interface{}, context context.Context, table string, fields Fields) error
	UpdateFieldFromTo(result interface{}, context context.Context, table, field, from, to string) error
	UpdateByUuid(context context.Context, table, uuid string, fields Fields) error
}

type postgresDb struct {
	logger   logger.Logger
	database *sqlx.DB
}
//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go

// NewPostgresDB - function for creating new postgres database connection.
// Function takes PostgresConfig structure and an input and returns
// *sql.DB instance and error
func Open(logger logger.Logger, config configs.PostgresConfig) (*postgresDb, error) {
	database, err := sqlx.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(config.MaxOpenConns)
	database.SetMaxIdleConns(config.MaxIdleConns)
	database.SetConnMaxLifetime(config.ConnMaxLifetime)

	logger.Info("SUCCESSFULLY CONNNECTED TO DATABASE")

	return &postgresDb{
		logger: logger,
		database: database,
	}, nil
}

func (d *postgresDb) Close() {
	if err := d.database.Close(); err != nil {
		d.logger.Errorf("Error on closing database - '%s'", err.Error())
	}
}

type conSt interface {
	sqlx.ExtContext

	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func (d *postgresDb) getCon(ctx context.Context) conSt {
	if tx := d.getContextTransaction(ctx); tx != nil {
		return tx
	}
	return d.database
}

func (d *postgresDb) getContextTransaction(ctx context.Context) *sqlx.Tx {
	contextV := ctx.Value(TransactionCtxKey)
	if contextV == nil {
		return nil
	}

	switch tx := contextV.(type) {
	case *sqlx.Tx:
		return tx
	default:
		return nil
	}
}

func (d *postgresDb) ContextWithTransaction(ctx context.Context) (context.Context, error) {
	tx, err := d.database.BeginTxx(ctx, nil)
	if err != nil {
		return ctx, err
	}
	ctx = context.WithValue(ctx, TransactionCtxKey, tx)

	return ctx, nil
}

func (d *postgresDb) CommitContextTransaction(ctx context.Context) error {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return nil
	}

	err := tx.Commit()
	if err != nil && err != sql.ErrTxDone {
		d.logger.Errorf("Fail to commit transaction", err)
		return err
	}

	return nil
}

func (d *postgresDb) RollbackContextTransaction(ctx context.Context) {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return
	}

	err := tx.Rollback()
	if err != nil && err != sql.ErrTxDone {
		d.logger.Errorw("Fail to rollback transaction", err)
	}
}

type Fields map[string]interface{}

func (d *postgresDb) InsertReturnId(result interface{}, context context.Context, table string, fields Fields) error {
	con := d.getCon(context)
	
	columns := make([]string, 0, len(fields))
	placeholders := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields))

	i := 1
	for column, arg := range fields {
		columns = append(columns, column)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		args = append(args, arg)
		i++
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (%s)
		RETURNING id
		`,
		table,
		strings.Join(columns, ","),
		strings.Join(placeholders, ","),
	)
	err := con.GetContext(context, result, query, args...)
	if err != nil {
		d.logger.Errorw(ErrMsg, err)
		return err
	}

	return nil
}

func (d *postgresDb) Get(result interface{}, context context.Context, table string, fields Fields) error {
	con := d.getCon(context)
	
	columns := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields))

	i := 1
	for column, arg := range fields {
		columns = append(columns, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, arg)
		i++
	}
	
	query := fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE %s
		`,
		table,
		strings.Join(columns, ","),
	)
	err := con.GetContext(context, result, query, args...)
	if err != nil {
		d.logger.Errorw(ErrMsg, err)
		return err
	}

	return nil
}

func (d *postgresDb) UpdateFieldFromTo(result interface{}, context context.Context, table, field, from, to string) error {
	con := d.getCon(context)
	
	query := fmt.Sprintf(`
		UPDATE %s
		SET %s = $1
		WHERE %s IN (
			SELECT %s
			FROM %s
			WHERE %s = $2
			LIMIT 100
			FOR UPDATE SKIP LOCKED
		)
		RETURNING *
		`,
		table,
		field,
		field,
		field,
		table,
		field,
	)
	err := con.SelectContext(context, result, query, to, from)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.ObjectNotFound
		}
		d.logger.Errorw(ErrMsg, err)
		return err
	}

	return nil
}

func (d *postgresDb) UpdateByUuid(context context.Context, table, uuid string, fields Fields) error {
	con := d.getCon(context)

	columns := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields) + 1)

	i := 2
	args = append(args, uuid)
	for column, arg := range fields {
		columns = append(columns, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, arg)
		i++
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET %s
		WHERE uuid = $1
		`,
		table,
		strings.Join(columns, ","),
	)
	res, err := con.ExecContext(context, query, args...)
	if err != nil {
		d.logger.Errorw(ErrMsg, err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		d.logger.Errorw(ErrMsg, err)
		return err
	}

	if rowsAffected == 0 {
		d.logger.Errorf("UpdateBy - no affected rows by uuid: %s", uuid)
		return errs.ObjectNotFound
	}

	return nil
}
