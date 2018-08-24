package qm

import (
	"strings"

	"github.com/volatiletech/sqlboiler/queries"
)

// QueryMod modifies a query object.
type QueryMod interface {
	Apply(q *queries.Query)
}

// The QueryModFunc type is an adapter to allow the use
// of ordinary functions for query modifying. If f is a
// function with the appropriate signature,
// QueryModFunc(f) is a QueryMod that calls f.
type QueryModFunc func(q *queries.Query)

// Apply calls f(q).
func (f QueryModFunc) Apply(q *queries.Query) {
	f(q)
}

type queryMods []QueryMod

// Apply applies the query mods to a query, satisfying
// the applicator interface in queries. This "clever"
// inversion of dependency is because suddenly the
// eager loading needs to be able to store query mods
// in the query object, which before - never knew about
// query mods.
func (m queryMods) Apply(q *queries.Query) {
	Apply(q, m...)
}

// Apply the query mods to the Query object
func Apply(q *queries.Query, mods ...QueryMod) {
	for _, mod := range mods {
		mod.Apply(q)
	}
}

type sqlQueryMod struct {
	sql  string
	args []interface{}
}

// Apply implements QueryMod.Apply.
func (qm sqlQueryMod) Apply(q *queries.Query) {
	queries.SetSQL(q, qm.sql, qm.args...)
}

// SQL allows you to execute a plain SQL statement
func SQL(sql string, args ...interface{}) QueryMod {
	return sqlQueryMod{
		sql:  sql,
		args: args,
	}
}

type loadQueryMod struct {
	relationship string
	mods         []QueryMod
}

// Apply implements QueryMod.Apply.
func (qm loadQueryMod) Apply(q *queries.Query) {
	queries.AppendLoad(q, qm.relationship)

	if len(qm.mods) != 0 {
		queries.SetLoadMods(q, qm.relationship, queryMods(qm.mods))
	}
}

// Load allows you to specify foreign key relationships to eager load
// for your query. Passed in relationships need to be in the format
// MyThing or MyThings.
// Relationship name plurality is important, if your relationship is
// singular, you need to specify the singular form and vice versa.
//
// In the following example we see how to eager load a users's videos
// and the video's tags comments, and publisher during a query to find users.
//
//   models.Users(qm.Load("Videos.Tags"))
//
// In order to filter better on the query for the relationships you can additionally
// supply query mods.
//
//   models.Users(qm.Load("Videos.Tags", Where("deleted = ?", isDeleted)))
//
// Keep in mind the above only sets the query mods for the query on the last specified
// relationship. In this case, only Tags will get the query mod. If you want to do
// intermediate relationships with query mods you must specify them separately:
//
//   models.Users(
//     qm.Load("Videos", Where("deleted = false"))
//     qm.Load("Videos.Tags", Where("deleted = ?", isDeleted))
//   )
func Load(relationship string, mods ...QueryMod) QueryMod {
	return loadQueryMod{
		relationship: relationship,
		mods:         mods,
	}
}

type innerJoinQueryMod struct {
	clause string
	args   []interface{}
}

// Apply implements QueryMod.Apply.
func (qm innerJoinQueryMod) Apply(q *queries.Query) {
	queries.AppendInnerJoin(q, qm.clause, qm.args...)
}

// InnerJoin on another table
func InnerJoin(clause string, args ...interface{}) QueryMod {
	return innerJoinQueryMod{
		clause: clause,
		args:   args,
	}
}

type selectQueryMod struct {
	columns []string
}

// Apply implements QueryMod.Apply.
func (qm selectQueryMod) Apply(q *queries.Query) {
	queries.AppendSelect(q, qm.columns...)
}

// Select specific columns opposed to all columns
func Select(columns ...string) QueryMod {
	return selectQueryMod{
		columns: columns,
	}
}

type whereQueryMod struct {
	clause string
	args   []interface{}
}

// Apply implements QueryMod.Apply.
func (qm whereQueryMod) Apply(q *queries.Query) {
	queries.AppendWhere(q, qm.clause, qm.args...)
}

// Where allows you to specify a where clause for your statement
func Where(clause string, args ...interface{}) QueryMod {
	return whereQueryMod{
		clause: clause,
		args:   args,
	}
}

type andQueryMod struct {
	clause string
	args   []interface{}
}

// Apply implements QueryMod.Apply.
func (qm andQueryMod) Apply(q *queries.Query) {
	queries.AppendWhere(q, qm.clause, qm.args...)
}

// And allows you to specify a where clause separated by an AND for your statement
// And is a duplicate of the Where function, but allows for more natural looking
// query mod chains, for example: (Where("a=?"), And("b=?"), Or("c=?")))
func And(clause string, args ...interface{}) QueryMod {
	return andQueryMod{
		clause: clause,
		args:   args,
	}
}

type orQueryMod struct {
	clause string
	args   []interface{}
}

// Apply implements QueryMod.Apply.
func (qm orQueryMod) Apply(q *queries.Query) {
	queries.AppendWhere(q, qm.clause, qm.args...)
	queries.SetLastWhereAsOr(q)
}

// Or allows you to specify a where clause separated by an OR for your statement
func Or(clause string, args ...interface{}) QueryMod {
	return orQueryMod{
		clause: clause,
		args:   args,
	}
}

// Apply implements QueryMod.Apply.
type whereInQueryMod struct {
	clause string
	args   []interface{}
}

func (qm whereInQueryMod) Apply(q *queries.Query) {
	queries.AppendIn(q, qm.clause, qm.args...)
}

// WhereIn allows you to specify a "x IN (set)" clause for your where statement
// Example clauses: "column in ?", "(column1,column2) in ?"
func WhereIn(clause string, args ...interface{}) QueryMod {
	return whereInQueryMod{
		clause: clause,
		args:   args,
	}
}

type andInQueryMod struct {
	clause string
	args   []interface{}
}

// Apply implements QueryMod.Apply.
func (qm andInQueryMod) Apply(q *queries.Query) {
	queries.AppendIn(q, qm.clause, qm.args...)
}

// AndIn allows you to specify a "x IN (set)" clause separated by an AndIn
// for your where statement. AndIn is a duplicate of the WhereIn function, but
// allows for more natural looking query mod chains, for example:
// (WhereIn("column1 in ?"), AndIn("column2 in ?"), OrIn("column3 in ?"))
func AndIn(clause string, args ...interface{}) QueryMod {
	return andInQueryMod{
		clause: clause,
		args:   args,
	}
}

type orInQueryMod struct {
	clause string
	args   []interface{}
}

// Apply implements QueryMod.Apply.
func (qm orInQueryMod) Apply(q *queries.Query) {
	queries.AppendIn(q, qm.clause, qm.args...)
	queries.SetLastInAsOr(q)
}

// OrIn allows you to specify an IN clause separated by
// an OR for your where statement
func OrIn(clause string, args ...interface{}) QueryMod {
	return orInQueryMod{
		clause: clause,
		args:   args,
	}
}

type groupByQueryMod struct {
	clause string
}

// Apply implements QueryMod.Apply.
func (qm groupByQueryMod) Apply(q *queries.Query) {
	queries.AppendGroupBy(q, qm.clause)
}

// GroupBy allows you to specify a group by clause for your statement
func GroupBy(clause string) QueryMod {
	return groupByQueryMod{
		clause: clause,
	}
}

type orderByQueryMod struct {
	clause string
}

// Apply implements QueryMod.Apply.
func (qm orderByQueryMod) Apply(q *queries.Query) {
	queries.AppendOrderBy(q, qm.clause)
}

// OrderBy allows you to specify a order by clause for your statement
func OrderBy(clause string) QueryMod {
	return orderByQueryMod{
		clause: clause,
	}
}

type havingQueryMod struct {
	clause string
	args   []interface{}
}

// Apply implements QueryMod.Apply.
func (qm havingQueryMod) Apply(q *queries.Query) {
	queries.AppendHaving(q, qm.clause, qm.args...)
}

// Having allows you to specify a having clause for your statement
func Having(clause string, args ...interface{}) QueryMod {
	return havingQueryMod{
		clause: clause,
		args:   args,
	}
}

type fromQueryMod struct {
	from string
}

// Apply implements QueryMod.Apply.
func (qm fromQueryMod) Apply(q *queries.Query) {
	queries.AppendFrom(q, qm.from)
}

// From allows to specify the table for your statement
func From(from string) QueryMod {
	return fromQueryMod{
		from: from,
	}
}

type limitQueryMod struct {
	limit int
}

// Apply implements QueryMod.Apply.
func (qm limitQueryMod) Apply(q *queries.Query) {
	queries.SetLimit(q, qm.limit)
}

// Limit the number of returned rows
func Limit(limit int) QueryMod {
	return limitQueryMod{
		limit: limit,
	}
}

type offsetQueryMod struct {
	offset int
}

// Apply implements QueryMod.Apply.
func (qm offsetQueryMod) Apply(q *queries.Query) {
	queries.SetOffset(q, qm.offset)
}

// Offset into the results
func Offset(offset int) QueryMod {
	return offsetQueryMod{
		offset: offset,
	}
}

type forQueryMod struct {
	clause string
}

// Apply implements QueryMod.Apply.
func (qm forQueryMod) Apply(q *queries.Query) {
	queries.SetFor(q, qm.clause)
}

// For inserts a concurrency locking clause at the end of your statement
func For(clause string) QueryMod {
	return forQueryMod{
		clause: clause,
	}
}

// Rels is an alias for strings.Join to make it easier to use relationship name
// constants in Load.
func Rels(r ...string) string {
	return strings.Join(r, ".")
}
