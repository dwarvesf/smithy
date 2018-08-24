package queries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/drivers"
)

// joinKind is the type of join
type joinKind int

// Join type constants
const (
	JoinInner joinKind = iota
	JoinOuterLeft
	JoinOuterRight
	JoinNatural
)

// Query holds the state for the built up query
type Query struct {
	dialect *drivers.Dialect
	rawSQL  rawSQL

	load     []string
	loadMods map[string]Applicator

	delete     bool
	update     map[string]interface{}
	selectCols []string
	count      bool
	from       []string
	joins      []join
	where      []where
	in         []in
	groupBy    []string
	orderBy    []string
	having     []having
	limit      int
	offset     int
	forlock    string
}

// Applicator exists only to allow
// query mods into the query struct around
// eager loaded relationships.
type Applicator interface {
	Apply(*Query)
}

type where struct {
	clause      string
	orSeparator bool
	args        []interface{}
}

type in struct {
	clause      string
	orSeparator bool
	args        []interface{}
}

type having struct {
	clause string
	args   []interface{}
}

type rawSQL struct {
	sql  string
	args []interface{}
}

type join struct {
	kind   joinKind
	clause string
	args   []interface{}
}

// Raw makes a raw query, usually for use with bind
func Raw(query string, args ...interface{}) *Query {
	return &Query{
		rawSQL: rawSQL{
			sql:  query,
			args: args,
		},
	}
}

// RawG makes a raw query using the global boil.Executor, usually for use with bind
func RawG(query string, args ...interface{}) *Query {
	return Raw(query, args...)
}

// Exec executes a query that does not need a row returned
func (q *Query) Exec(exec boil.Executor) (sql.Result, error) {
	qs, args := BuildQuery(q)
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, qs)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	return exec.Exec(qs, args...)
}

// QueryRow executes the query for the One finisher and returns a row
func (q *Query) QueryRow(exec boil.Executor) *sql.Row {
	qs, args := BuildQuery(q)
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, qs)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	return exec.QueryRow(qs, args...)
}

// Query executes the query for the All finisher and returns multiple rows
func (q *Query) Query(exec boil.Executor) (*sql.Rows, error) {
	qs, args := BuildQuery(q)
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, qs)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	return exec.Query(qs, args...)
}

// ExecContext executes a query that does not need a row returned
func (q *Query) ExecContext(ctx context.Context, exec boil.ContextExecutor) (sql.Result, error) {
	qs, args := BuildQuery(q)
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, qs)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	return exec.ExecContext(ctx, qs, args...)
}

// QueryRowContext executes the query for the One finisher and returns a row
func (q *Query) QueryRowContext(ctx context.Context, exec boil.ContextExecutor) *sql.Row {
	qs, args := BuildQuery(q)
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, qs)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	return exec.QueryRowContext(ctx, qs, args...)
}

// QueryContext executes the query for the All finisher and returns multiple rows
func (q *Query) QueryContext(ctx context.Context, exec boil.ContextExecutor) (*sql.Rows, error) {
	qs, args := BuildQuery(q)
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, qs)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	return exec.QueryContext(ctx, qs, args...)
}

// ExecP executes a query that does not need a row returned
// It will panic on error
func (q *Query) ExecP(exec boil.Executor) sql.Result {
	res, err := q.Exec(exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return res
}

// QueryP executes the query for the All finisher and returns multiple rows
// It will panic on error
func (q *Query) QueryP(exec boil.Executor) *sql.Rows {
	rows, err := q.Query(exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rows
}

// SetDialect on the query.
func SetDialect(q *Query, dialect *drivers.Dialect) {
	q.dialect = dialect
}

// SetSQL on the query.
func SetSQL(q *Query, sql string, args ...interface{}) {
	q.rawSQL = rawSQL{sql: sql, args: args}
}

// SetArgs is primarily for re-use of a query so that the
// query text does not need to be re-generated, useful
// if you're performing the same query with different arguments
// over and over.
func SetArgs(q *Query, args ...interface{}) {
	q.rawSQL.args = args
}

// SetLoad on the query.
func SetLoad(q *Query, relationships ...string) {
	q.load = append([]string(nil), relationships...)
}

// AppendLoad on the query.
func AppendLoad(q *Query, relationships string) {
	q.load = append(q.load, relationships)
}

// SetLoadMods on the query.
func SetLoadMods(q *Query, rel string, appl Applicator) {
	if q.loadMods == nil {
		q.loadMods = make(map[string]Applicator)
	}

	q.loadMods[rel] = appl
}

// SetSelect on the query.
func SetSelect(q *Query, sel []string) {
	q.selectCols = sel
}

// GetSelect from the query
func GetSelect(q *Query) []string {
	return q.selectCols
}

// SetCount on the query.
func SetCount(q *Query) {
	q.count = true
}

// SetDelete on the query.
func SetDelete(q *Query) {
	q.delete = true
}

// SetLimit on the query.
func SetLimit(q *Query, limit int) {
	q.limit = limit
}

// SetOffset on the query.
func SetOffset(q *Query, offset int) {
	q.offset = offset
}

// SetFor on the query.
func SetFor(q *Query, clause string) {
	q.forlock = clause
}

// SetUpdate on the query.
func SetUpdate(q *Query, cols map[string]interface{}) {
	q.update = cols
}

// AppendSelect on the query.
func AppendSelect(q *Query, columns ...string) {
	q.selectCols = append(q.selectCols, columns...)
}

// AppendFrom on the query.
func AppendFrom(q *Query, from ...string) {
	q.from = append(q.from, from...)
}

// SetFrom replaces the current from statements.
func SetFrom(q *Query, from ...string) {
	q.from = append([]string(nil), from...)
}

// AppendInnerJoin on the query.
func AppendInnerJoin(q *Query, clause string, args ...interface{}) {
	q.joins = append(q.joins, join{clause: clause, kind: JoinInner, args: args})
}

// AppendHaving on the query.
func AppendHaving(q *Query, clause string, args ...interface{}) {
	q.having = append(q.having, having{clause: clause, args: args})
}

// AppendWhere on the query.
func AppendWhere(q *Query, clause string, args ...interface{}) {
	q.where = append(q.where, where{clause: clause, args: args})
}

// AppendIn on the query.
func AppendIn(q *Query, clause string, args ...interface{}) {
	q.in = append(q.in, in{clause: clause, args: args})
}

// SetLastWhereAsOr sets the or separator for the tail "WHERE" in the slice
func SetLastWhereAsOr(q *Query) {
	if len(q.where) == 0 {
		return
	}

	q.where[len(q.where)-1].orSeparator = true
}

// SetLastInAsOr sets the or separator for the tail "IN" in the slice
func SetLastInAsOr(q *Query) {
	if len(q.in) == 0 {
		return
	}

	q.in[len(q.in)-1].orSeparator = true
}

// AppendGroupBy on the query.
func AppendGroupBy(q *Query, clause string) {
	q.groupBy = append(q.groupBy, clause)
}

// AppendOrderBy on the query.
func AppendOrderBy(q *Query, clause string) {
	q.orderBy = append(q.orderBy, clause)
}
