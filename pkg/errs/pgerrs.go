package errs

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Common normalized PostgreSQL / pgx error types.
var (
	// === Logical & data-level errors ===
	ErrPgNotFound             = errors.New("pg_not_found")
	ErrPgDuplicate            = errors.New("pg_duplicate")
	ErrPgForeignKey           = errors.New("pg_foreign_key_violation")
	ErrPgInvalidInput         = errors.New("pg_invalid_input")
	ErrPgConstraint           = errors.New("pg_constraint_violation")
	ErrPgNullViolation        = errors.New("pg_not_null_violation")
	ErrPgCheckViolation       = errors.New("pg_check_violation")
	ErrPgExclusionViolation   = errors.New("pg_exclusion_violation")
	ErrPgSerializationFailure = errors.New("pg_serialization_failure")
	ErrPgDeadlock             = errors.New("pg_deadlock_detected")

	// === Data / formatting ===
	ErrPgNumericOutOfRange  = errors.New("pg_numeric_value_out_of_range")
	ErrPgInvalidTextFormat  = errors.New("pg_invalid_text_representation")
	ErrPgStringTooLong      = errors.New("pg_string_data_right_truncation")
	ErrPgInvalidDatetime    = errors.New("pg_invalid_datetime_format")
	ErrPgDivisionByZero     = errors.New("pg_division_by_zero")
	ErrPgUntranslatableChar = errors.New("pg_untranslatable_character")

	// === Access / privileges ===
	ErrPgInsufficientPrivilege = errors.New("pg_insufficient_privilege")
	ErrPgInvalidPassword       = errors.New("pg_invalid_password")
	ErrPgUnauthorized          = errors.New("pg_unauthorized")
	ErrPgUndefinedTable        = errors.New("pg_undefined_table")
	ErrPgUndefinedColumn       = errors.New("pg_undefined_column")

	// === Transaction & connection ===
	ErrPgConnection     = errors.New("pg_connection_failure")
	ErrPgTxState        = errors.New("pg_invalid_transaction_state")
	ErrPgInFailedTx     = errors.New("pg_in_failed_sql_transaction")
	ErrPgIdleTxTimeout  = errors.New("pg_idle_in_transaction_timeout")
	ErrPgConnectionLost = errors.New("pg_connection_lost")

	// === Resource / system ===
	ErrPgOutOfMemory           = errors.New("pg_out_of_memory")
	ErrPgDiskFull              = errors.New("pg_disk_full")
	ErrPgTooManyConnections    = errors.New("pg_too_many_connections")
	ErrPgSystemIO              = errors.New("pg_io_error")
	ErrPgSystemInternal        = errors.New("pg_system_error")
	ErrPgDataCorrupted         = errors.New("pg_data_corrupted")
	ErrPgIndexCorrupted        = errors.New("pg_index_corrupted")
	ErrPgConfigurationExceeded = errors.New("pg_configuration_limit_exceeded")

	// === Syntax / schema ===
	ErrPgSyntaxError      = errors.New("pg_syntax_error")
	ErrPgDatatypeMismatch = errors.New("pg_datatype_mismatch")
	ErrPgDuplicateObject  = errors.New("pg_duplicate_object")
	ErrPgInvalidSchema    = errors.New("pg_invalid_schema_definition")
	ErrPgInvalidTableDef  = errors.New("pg_invalid_table_definition")

	// === Operator / cancel ===
	ErrPgQueryCanceled   = errors.New("pg_query_canceled")
	ErrPgAdminShutdown   = errors.New("pg_admin_shutdown")
	ErrPgCrashShutdown   = errors.New("pg_crash_shutdown")
	ErrPgDatabaseDropped = errors.New("pg_database_dropped")

	// === Internal fallback ===
	ErrPgInternal = errors.New("pg_internal_error")
	ErrPgUnknown  = errors.New("pg_unknown")
)

func ParsePgError(err error) error {
	if err == nil {
		return nil
	}

	// 1. Check for logical pgx errors
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrPgNotFound
	}
	if errors.Is(err, pgx.ErrTxClosed) {
		return ErrPgTxState
	}

	// 2. Try PostgreSQL-specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {

		// === Class 23 — Integrity Constraint Violation ===
		case "23505":
			return ErrPgDuplicate
		case "23503":
			return ErrPgForeignKey
		case "23502":
			return ErrPgNullViolation
		case "23514":
			return ErrPgCheckViolation
		case "23P01":
			return ErrPgExclusionViolation

		// === Class 22 — Data Exception ===
		case "22P02":
			return ErrPgInvalidInput
		case "22003":
			return ErrPgNumericOutOfRange
		case "22001":
			return ErrPgStringTooLong
		case "22007":
			return ErrPgInvalidDatetime
		case "22012":
			return ErrPgDivisionByZero
		case "22P05":
			return ErrPgUntranslatableChar

		// === Class 08 — Connection Exception ===
		case "08006", "08001", "08003":
			return ErrPgConnection
		case "57P03":
			return ErrPgConnectionLost

		// === Class 25 — Invalid Transaction State ===
		case "25000":
			return ErrPgTxState
		case "25P02":
			return ErrPgInFailedTx
		case "25P03":
			return ErrPgIdleTxTimeout

		// === Class 40 — Transaction Rollback ===
		case "40001":
			return ErrPgSerializationFailure
		case "40P01":
			return ErrPgDeadlock

		// === Class 42 — Syntax Error or Access Rule Violation ===
		case "42601":
			return ErrPgSyntaxError
		case "42804":
			return ErrPgDatatypeMismatch
		case "42701", "42702", "42703", "42704":
			return ErrPgDuplicateObject
		case "42P06":
			return ErrPgInvalidSchema
		case "42P17":
			return ErrPgInvalidTableDef
		case "42501":
			return ErrPgInsufficientPrivilege
		case "28P01":
			return ErrPgInvalidPassword

		// === Class 53 — Resource Issues ===
		case "53100":
			return ErrPgDiskFull
		case "53200":
			return ErrPgOutOfMemory
		case "53300":
			return ErrPgTooManyConnections
		case "53400":
			return ErrPgConfigurationExceeded

		// === Class 55 — Object State ===
		case "55P03":
			return ErrPgSystemIO
		case "55000":
			return ErrPgSystemInternal

		// === Class 57 — Operator Intervention ===
		case "57014":
			return ErrPgQueryCanceled
		case "57P01":
			return ErrPgAdminShutdown
		case "57P02":
			return ErrPgCrashShutdown
		case "57P04":
			return ErrPgDatabaseDropped

		// === Class 58 — System Errors ===
		case "58030":
			return ErrPgSystemIO
		case "XX000":
			return ErrPgInternal
		case "XX001":
			return ErrPgDataCorrupted
		case "XX002":
			return ErrPgIndexCorrupted

		default:
			return err
		}
	}

	// 3. Fallback for unrecognized errors
	return err
}
