package testing

import (
	"GoAuth/internal/apierr"
	"database/sql"
	"errors"
	"testing"

	"github.com/MintzyG/fail/v3"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func testPGXMapper(t *testing.T) {
	mapper := &apierr.PGXMapper{}

	tests := []struct {
		name       string
		err        error
		expectedID fail.ErrorID
		expectedOk bool
		checkMeta  func(*testing.T, *fail.Error)
	}{
		{
			name:       "SQL No Rows",
			err:        sql.ErrNoRows,
			expectedID: apierr.SQLNotFound,
			expectedOk: true,
		},
		{
			name: "Unique Violation - Schema Version Draft",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "one_version_draft_per_schema",
			},
			expectedID: apierr.SCHEMAVersionDraftAlreadyExists,
			expectedOk: true,
			checkMeta: func(t *testing.T, fe *fail.Error) {
				assert.True(t, apierr.IsUniqueViolationNew(fe))
			},
		},
		{
			name: "Unique Violation - Role Name",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "roles_name_project_id_key",
			},
			expectedID: apierr.ROLENameAlreadyTaken,
			expectedOk: true,
			checkMeta: func(t *testing.T, fe *fail.Error) {
				assert.True(t, apierr.IsUniqueViolationNew(fe))
			},
		},
		/*{
			name: "Unique Violation - Unmatched",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "some_other_unique",
			},
			expectedID: apierr.SQLUnmatchedUniqueViolation,
			expectedOk: true,
			checkMeta: func(t *testing.T, fe *fail.Error) {
				assert.True(t, apierr.IsUniqueViolationNew(fe))
			},
		},*/
		{
			name: "Check Violation - Schema Field Key",
			err: &pgconn.PgError{
				Code:           "23514",
				ConstraintName: "schema_fields_key_check",
			},
			expectedID: apierr.FIELDInvalidCharactersInKey,
			expectedOk: true,
			checkMeta: func(t *testing.T, fe *fail.Error) {
				assert.True(t, apierr.IsCheckViolationNew(fe))
			},
		},
		{
			name: "Check Violation - Unmatched",
			err: &pgconn.PgError{
				Code:           "23514",
				ConstraintName: "some_other_check",
			},
			expectedID: apierr.SQLUnmatchedCheckViolation,
			expectedOk: true,
			checkMeta: func(t *testing.T, fe *fail.Error) {
				assert.True(t, apierr.IsCheckViolationNew(fe))
			},
		},
		{
			name: "Foreign Key Violation",
			err: &pgconn.PgError{
				Code: "23503",
			},
			expectedID: apierr.SQLForeignKeyViolation,
			expectedOk: true,
		},
		{
			name: "Not NULL Violation",
			err: &pgconn.PgError{
				Code: "23502",
			},
			expectedID: apierr.SQLNotNULLViolation,
			expectedOk: true,
		},
		{
			name: "Value Too Long",
			err: &pgconn.PgError{
				Code: "22001",
			},
			expectedID: apierr.SQLValueTooLong,
			expectedOk: true,
		},
		{
			name: "Serialization Failure",
			err: &pgconn.PgError{
				Code: "40001",
			},
			expectedID: apierr.SQLSerializationFailure,
			expectedOk: true,
		},
		{
			name: "Connection Error 08006",
			err: &pgconn.PgError{
				Code: "08006",
			},
			expectedID: apierr.SQLDBConnectionError,
			expectedOk: true,
		},
		{
			name: "Unknown PG Error",
			err: &pgconn.PgError{
				Code: "99999",
			},
			expectedID: apierr.SQLUnknownError,
			expectedOk: true,
		},
		{
			name:       "Non-PG Error",
			err:        errors.New("generic error"),
			expectedID: apierr.SQLInternalDBError,
			expectedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe, ok := mapper.Map(tt.err)
			assert.Equal(t, tt.expectedOk, ok)
			if ok {
				assert.Equal(t, tt.expectedID, fe.ID)
				if tt.checkMeta != nil {
					tt.checkMeta(t, fe)
				}
			}
		})
	}
}
