package service

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/mocks"
)

func newDatasourceSvc(t *testing.T) (DatasourceService, *mocks.MockDatasourceDB) {
	t.Helper()
	db := mocks.NewMockDatasourceDB(t)
	svc := &datasourceService{
		db:     db,
		logger: zerolog.Nop(),
	}
	return svc, db
}

// LIST SCHEMAS
func TestDatasourceService_ListSchemas_Success(t *testing.T) {
	svc, db := newDatasourceSvc(t)
	ctx := context.Background()

	dbOut := dto.ListSchemasResult{}

	db.EXPECT().ListSchemas(mock.Anything).Return(dbOut, nil).Once()

	res, err := svc.ListSchemas(ctx)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)
}

func TestDatasourceService_ListSchemas_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db := newDatasourceSvc(t)
	ctx := context.Background()

	dbOut := dto.ListSchemasResult{}
	dbErr := errs.PgErrDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().ListSchemas(mock.Anything).Return(dbOut, dbErr).Once()

	res, err := svc.ListSchemas(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dbOut, res)
}

// LIST TABLES
func TestDatasourceService_ListTables_Success(t *testing.T) {
	svc, db := newDatasourceSvc(t)
	ctx := context.Background()

	in := dto.ListTablesParams{Schema: "public"}
	dbOut := dto.ListTablesResult{}

	db.EXPECT().ListTables(mock.Anything, in).Return(dbOut, nil).Once()

	res, err := svc.ListTables(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)
}

func TestDatasourceService_ListTables_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db := newDatasourceSvc(t)
	ctx := context.Background()

	in := dto.ListTablesParams{Schema: "public"}
	dbOut := dto.ListTablesResult{}

	dbErr := errs.PgErrDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().ListTables(mock.Anything, in).Return(dbOut, dbErr).Once()

	res, err := svc.ListTables(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dbOut, res)
}

// LIST COLUMNS
func TestDatasourceService_ListColumns_Success(t *testing.T) {
	svc, db := newDatasourceSvc(t)
	ctx := context.Background()

	in := dto.ListColumnsParams{Schema: "public", Table: "users"}
	dbOut := dto.ListColumnsResult{}

	db.EXPECT().ListColumns(mock.Anything, in).Return(dbOut, nil).Once()

	res, err := svc.ListColumns(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)
}

func TestDatasourceService_ListColumns_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db := newDatasourceSvc(t)
	ctx := context.Background()

	in := dto.ListColumnsParams{Schema: "public", Table: "users"}
	dbOut := dto.ListColumnsResult{}

	dbErr := errs.PgErrDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().ListColumns(mock.Anything, in).Return(dbOut, dbErr).Once()

	res, err := svc.ListColumns(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dbOut, res)
}
