package errs

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Common normalized PostgreSQL / pgx error types.
var (
	// === Logical & data-level errors ===
	PgErrNotFound             = errors.New("pg_not_found")
	PgErrDuplicate            = errors.New("pg_duplicate")
	PgErrForeignKey           = errors.New("pg_foreign_key_violation")
	PgErrInvalidInput         = errors.New("pg_invalid_input")
	PgErrConstraint           = errors.New("pg_constraint_violation")
	PgErrNullViolation        = errors.New("pg_not_null_violation")
	PgErrCheckViolation       = errors.New("pg_check_violation")
	PgErrExclusionViolation   = errors.New("pg_exclusion_violation")
	PgErrSerializationFailure = errors.New("pg_serialization_failure")
	PgErrDeadlock             = errors.New("pg_deadlock_detected")

	// === Data / formatting ===
	PgErrNumericOutOfRange  = errors.New("pg_numeric_value_out_of_range")
	PgErrInvalidTextFormat  = errors.New("pg_invalid_text_representation")
	PgErrStringTooLong      = errors.New("pg_string_data_right_truncation")
	PgErrInvalidDatetime    = errors.New("pg_invalid_datetime_format")
	PgErrDivisionByZero     = errors.New("pg_division_by_zero")
	PgErrUntranslatableChar = errors.New("pg_untranslatable_character")

	// === Access / privileges ===
	PgErrInsufficientPrivilege = errors.New("pg_insufficient_privilege")
	PgErrInvalidPassword       = errors.New("pg_invalid_password")
	PgErrUnauthorized          = errors.New("pg_unauthorized")
	PgErrUndefinedTable        = errors.New("pg_undefined_table")
	PgErrUndefinedColumn       = errors.New("pg_undefined_column")

	// === Transaction & connection ===
	PgErrConnection     = errors.New("pg_connection_failure")
	PgErrTxState        = errors.New("pg_invalid_transaction_state")
	PgErrInFailedTx     = errors.New("pg_in_failed_sql_transaction")
	PgErrIdleTxTimeout  = errors.New("pg_idle_in_transaction_timeout")
	PgErrConnectionLost = errors.New("pg_connection_lost")

	// === Resource / system ===
	PgErrOutOfMemory           = errors.New("pg_out_of_memory")
	PgErrDiskFull              = errors.New("pg_disk_full")
	PgErrTooManyConnections    = errors.New("pg_too_many_connections")
	PgErrSystemIO              = errors.New("pg_io_error")
	PgErrSystemInternal        = errors.New("pg_system_error")
	PgErrDataCorrupted         = errors.New("pg_data_corrupted")
	PgErrIndexCorrupted        = errors.New("pg_index_corrupted")
	PgErrConfigurationExceeded = errors.New("pg_configuration_limit_exceeded")

	// === Syntax / schema ===
	PgErrSyntaxError      = errors.New("pg_syntax_error")
	PgErrDatatypeMismatch = errors.New("pg_datatype_mismatch")
	PgErrDuplicateObject  = errors.New("pg_duplicate_object")
	PgErrInvalidSchema    = errors.New("pg_invalid_schema_definition")
	PgErrInvalidTableDef  = errors.New("pg_invalid_table_definition")

	// === Operator / cancel ===
	PgErrQueryCanceled   = errors.New("pg_query_canceled")
	PgErrAdminShutdown   = errors.New("pg_admin_shutdown")
	PgErrCrashShutdown   = errors.New("pg_crash_shutdown")
	PgErrDatabaseDropped = errors.New("pg_database_dropped")

	// === Internal fallback ===
	PgErrInternal = errors.New("pg_internal_error")
	PgErrUnknown  = errors.New("pg_unknown")
)

func ParsePgError(err error) error {
	if err == nil {
		return nil
	}

	// 1. Check for logical pgx errors
	if errors.Is(err, pgx.ErrNoRows) {
		return PgErrNotFound
	}
	if errors.Is(err, pgx.ErrTxClosed) {
		return PgErrTxState
	}

	// 2. Try PostgreSQL-specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {

		// === Class 23 — Integrity Constraint Violation ===
		case "23505":
			return PgErrDuplicate
		case "23503":
			return PgErrForeignKey
		case "23502":
			return PgErrNullViolation
		case "23514":
			return PgErrCheckViolation
		case "23P01":
			return PgErrExclusionViolation

		// === Class 22 — Data Exception ===
		case "22P02":
			return PgErrInvalidInput
		case "22003":
			return PgErrNumericOutOfRange
		case "22001":
			return PgErrStringTooLong
		case "22007":
			return PgErrInvalidDatetime
		case "22012":
			return PgErrDivisionByZero
		case "22P05":
			return PgErrUntranslatableChar

		// === Class 08 — Connection Exception ===
		case "08006", "08001", "08003":
			return PgErrConnection
		case "57P03":
			return PgErrConnectionLost

		// === Class 25 — Invalid Transaction State ===
		case "25000":
			return PgErrTxState
		case "25P02":
			return PgErrInFailedTx
		case "25P03":
			return PgErrIdleTxTimeout

		// === Class 40 — Transaction Rollback ===
		case "40001":
			return PgErrSerializationFailure
		case "40P01":
			return PgErrDeadlock

		// === Class 42 — Syntax Error or Access Rule Violation ===
		case "42601":
			return PgErrSyntaxError
		case "42804":
			return PgErrDatatypeMismatch
		case "42701", "42702", "42703", "42704":
			return PgErrDuplicateObject
		case "42P06":
			return PgErrInvalidSchema
		case "42P17":
			return PgErrInvalidTableDef
		case "42501":
			return PgErrInsufficientPrivilege
		case "28P01":
			return PgErrInvalidPassword

		// === Class 53 — Resource Issues ===
		case "53100":
			return PgErrDiskFull
		case "53200":
			return PgErrOutOfMemory
		case "53300":
			return PgErrTooManyConnections
		case "53400":
			return PgErrConfigurationExceeded

		// === Class 55 — Object State ===
		case "55P03":
			return PgErrSystemIO
		case "55000":
			return PgErrSystemInternal

		// === Class 57 — Operator Intervention ===
		case "57014":
			return PgErrQueryCanceled
		case "57P01":
			return PgErrAdminShutdown
		case "57P02":
			return PgErrCrashShutdown
		case "57P04":
			return PgErrDatabaseDropped

		// === Class 58 — System Errors ===
		case "58030":
			return PgErrSystemIO
		case "XX000":
			return PgErrInternal
		case "XX001":
			return PgErrDataCorrupted
		case "XX002":
			return PgErrIndexCorrupted

		default:
			return err
		}
	}

	// 3. Fallback for unrecognized errors
	return err
}
