{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range $rel := .Table.ToManyRelationships -}}
		{{- $ltable := $.Aliases.Table $rel.Table -}}
		{{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
		{{- $relAlias := $.Aliases.ManyRelationship $rel.ForeignTable $rel.Name $rel.JoinTable $rel.JoinLocalFKeyName -}}
		{{- $col := $ltable.Column $rel.Column -}}
		{{- $fcol := $ftable.Column $rel.ForeignColumn -}}
		{{- $usesPrimitives := usesPrimitives $.Tables $rel.Table $rel.Column $rel.ForeignTable $rel.ForeignColumn -}}
		{{- $arg := printf "maybe%s" $ltable.UpSingular -}}
		{{- $schemaForeignTable := $rel.ForeignTable | $.SchemaTable}}
// Load{{$relAlias.Local}} allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func ({{$ltable.DownSingular}}L) Load{{$relAlias.Local}}({{if $.NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, singular bool, {{$arg}} interface{}, mods queries.Applicator) error {
	var slice []*{{$ltable.UpSingular}}
	var object *{{$ltable.UpSingular}}

	if singular {
		object = {{$arg}}.(*{{$ltable.UpSingular}})
	} else {
		slice = *{{$arg}}.(*[]*{{$ltable.UpSingular}})
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &{{$ltable.DownSingular}}R{}
		}
		args = append(args, object.{{.Column | titleCase}})
	} else {
		Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &{{$ltable.DownSingular}}R{}
			}

			for _, a := range args {
				{{if $usesPrimitives -}}
				if a == obj.{{$col}} {
				{{else -}}
				if queries.Equal(a, obj.{{$col}}) {
				{{end -}}
					continue Outer
				}
			}

			args = append(args, obj.{{.Column | titleCase}})
		}
	}

		{{if .ToJoinTable -}}
			{{- $schemaJoinTable := .JoinTable | $.SchemaTable -}}
	query := NewQuery(
		qm.Select("{{$schemaForeignTable}}.*, {{id 0 | $.Quotes}}.{{.JoinLocalColumn | $.Quotes}}"),
		qm.From("{{$schemaForeignTable}}"),
		qm.InnerJoin("{{$schemaJoinTable}} as {{id 0 | $.Quotes}} on {{$schemaForeignTable}}.{{.ForeignColumn | $.Quotes}} = {{id 0 | $.Quotes}}.{{.JoinForeignColumn | $.Quotes}}"),
		qm.WhereIn("{{id 0 | $.Quotes}}.{{.JoinLocalColumn | $.Quotes}} in ?", args...),
	)
		{{else -}}
	query := NewQuery(qm.From(`{{if $.Dialect.UseSchema}}{{$.Schema}}.{{end}}{{.ForeignTable}}`), qm.WhereIn(`{{.ForeignColumn}} in ?`, args...))
		{{end -}}
	if mods != nil {
		mods.Apply(query)
	}

	{{if $.NoContext -}}
	results, err := query.Query(e)
	{{else -}}
	results, err := query.QueryContext(ctx, e)
	{{end -}}
	if err != nil {
		return errors.Wrap(err, "failed to eager load {{.ForeignTable}}")
	}

	var resultSlice []*{{$ftable.UpSingular}}
	{{if .ToJoinTable -}}
	{{- $foreignTable := getTable $.Tables .ForeignTable -}}
	{{- $joinTable := getTable $.Tables .JoinTable -}}
	{{- $localCol := $joinTable.GetColumn .JoinLocalColumn}}
	var localJoinCols []{{$localCol.Type}}
	for results.Next() {
		one := new({{$ftable.UpSingular}})
		var localJoinCol {{$localCol.Type}}

		err = results.Scan({{$foreignTable.Columns | columnNames | stringMap $.StringFuncs.titleCase | prefixStringSlice "&one." | join ", "}}, &localJoinCol)
		if err != nil {
			return errors.Wrap(err, "failed to scan eager loaded results for {{.ForeignTable}}")
		}
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice {{.ForeignTable}}")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}
	{{- else -}}
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice {{.ForeignTable}}")
	}
	{{- end}}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on {{.ForeignTable}}")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for {{.ForeignTable}}")
	}

	{{if not $.NoHooks -}}
	if len({{$ftable.DownSingular}}AfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks({{if $.NoContext}}e{{else}}ctx, e{{end -}}); err != nil {
				return err
			}
		}
	}

	{{- end}}
	if singular {
		object.R.{{$relAlias.Local}} = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &{{$ftable.DownSingular}}R{}
			}
			{{if .ToJoinTable -}}
			foreign.R.{{$relAlias.Foreign}} = append(foreign.R.{{$relAlias.Foreign}}, object)
			{{else -}}
			foreign.R.{{$relAlias.Foreign}} = object
			{{end -}}
		}
		return nil
	}

	{{if .ToJoinTable -}}
	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			{{if $usesPrimitives -}}
			if local.{{$col}} == localJoinCol {
			{{else -}}
			if queries.Equal(local.{{$col}}, localJoinCol) {
			{{end -}}
				local.R.{{$relAlias.Local}} = append(local.R.{{$relAlias.Local}}, foreign)
				if foreign.R == nil {
					foreign.R = &{{$ftable.DownSingular}}R{}
				}
				foreign.R.{{$relAlias.Foreign}} = append(foreign.R.{{$relAlias.Foreign}}, local)
				break
			}
		}
	}
	{{else -}}
	for _, foreign := range resultSlice {
		for _, local := range slice {
			{{if $usesPrimitives -}}
			if local.{{$col}} == foreign.{{$fcol}} {
			{{else -}}
			if queries.Equal(local.{{$col}}, foreign.{{$fcol}}) {
			{{end -}}
				local.R.{{$relAlias.Local}} = append(local.R.{{$relAlias.Local}}, foreign)
				if foreign.R == nil {
					foreign.R = &{{$ftable.DownSingular}}R{}
				}
				foreign.R.{{$relAlias.Foreign}} = local
				break
			}
		}
	}
	{{end}}

	return nil
}

{{end -}}{{/* range tomany */}}
{{- end -}}{{/* if IsJoinTable */}}
