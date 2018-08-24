{{- $alias := .Aliases.Table .Table.Name}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{if .AddGlobal -}}
// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *{{$alias.UpSingular}}) UpsertG({{if not .NoContext}}ctx context.Context, {{end -}} updateColumns, insertColumns boil.Columns) error {
	return o.Upsert({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, updateColumns, insertColumns)
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *{{$alias.UpSingular}}) UpsertGP({{if not .NoContext}}ctx context.Context, {{end -}} updateColumns, insertColumns boil.Columns) {
	if err := o.Upsert({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, updateColumns, insertColumns); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if .AddPanic -}}
// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *{{$alias.UpSingular}}) UpsertP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, updateColumns, insertColumns boil.Columns) {
	if err := o.Upsert({{if not .NoContext}}ctx, {{end -}} exec, updateColumns, insertColumns); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *{{$alias.UpSingular}}) Upsert({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for upsert")
	}

	{{- template "timestamp_upsert_helper" . }}

	{{if not .NoHooks -}}
	if err := o.doBeforeUpsertHooks({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		return err
	}
	{{- end}}

	nzDefaults := queries.NonZeroDefaultSet({{$alias.DownSingular}}ColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	{{$alias.DownSingular}}UpsertCacheMut.RLock()
	cache, cached := {{$alias.DownSingular}}UpsertCache[key]
	{{$alias.DownSingular}}UpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			{{$alias.DownSingular}}Columns,
			{{$alias.DownSingular}}ColumnsWithDefault,
			{{$alias.DownSingular}}ColumnsWithoutDefault,
			nzDefaults,
		)
		insert = strmangle.SetComplement(insert, {{$alias.DownSingular}}ColumnsWithAuto)
		for i, v := range insert {
			if strmangle.ContainsAny({{$alias.DownSingular}}PrimaryKeyColumns, v) && strmangle.ContainsAny({{$alias.DownSingular}}ColumnsWithDefault, v) {
				insert = append(insert[:i], insert[i+1:]...)
			}
		}
		if len(insert) == 0 {
			return errors.New("{{.PkgName}}: unable to upsert {{.Table.Name}}, could not build insert column list")
		}

		ret = strmangle.SetMerge(ret, {{$alias.DownSingular}}ColumnsWithAuto)
		ret = strmangle.SetMerge(ret, {{$alias.DownSingular}}ColumnsWithDefault)

		update := updateColumns.UpdateColumnSet(
			{{$alias.DownSingular}}Columns,
			{{$alias.DownSingular}}PrimaryKeyColumns,
		)
		update = strmangle.SetComplement(update, {{$alias.DownSingular}}ColumnsWithAuto)

		if len(update) == 0 {
			return errors.New("{{.PkgName}}: unable to upsert {{.Table.Name}}, could not build update column list")
		}

		cache.query = buildUpsertQueryMSSQL(dialect, "{{.Table.Name}}", {{$alias.DownSingular}}PrimaryKeyColumns, update, insert, ret)

		whitelist := make([]string, len({{$alias.DownSingular}}PrimaryKeyColumns))
		copy(whitelist, {{$alias.DownSingular}}PrimaryKeyColumns)
		whitelist = append(whitelist, update...)
		whitelist = append(whitelist, insert...)

		cache.valueMapping, err = queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		{{if .NoContext -}}
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		{{else -}}
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		{{end -}}
		if err == sql.ErrNoRows {
			err = nil // MSSQL doesn't return anything when there's no update
		}
	} else {
		{{if .NoContext -}}
		_, err = exec.Exec(cache.query, vals...)
		{{else -}}
		_, err = exec.ExecContext(ctx, cache.query, vals...)
		{{end -}}
	}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to upsert {{.Table.Name}}")
	}

	if !cached {
		{{$alias.DownSingular}}UpsertCacheMut.Lock()
		{{$alias.DownSingular}}UpsertCache[key] = cache
		{{$alias.DownSingular}}UpsertCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return o.doAfterUpsertHooks({{if not .NoContext}}ctx, {{end -}} exec)
	{{- else -}}
	return nil
	{{- end}}
}
