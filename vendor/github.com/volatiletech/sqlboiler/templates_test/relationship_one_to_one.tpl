{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Tables $dot.Table . -}}
		{{- $varNameSingular := .Table | singular | camelCase -}}
		{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase}}
func test{{$txt.LocalTable.NameGo}}OneToOne{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var foreign {{$txt.ForeignTable.NameGo}}
	var local {{$txt.LocalTable.NameGo}}

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &foreign, {{$foreignVarNameSingular}}DBTypes, true, {{$foreignVarNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$txt.ForeignTable.NameGo}} struct: %s", err)
	}
	if err := randomize.Struct(seed, &local, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$txt.LocalTable.NameGo}} struct: %s", err)
	}

	{{if .ForeignColumnNullable -}}
	foreign.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
	{{- end}}
	{{if .Nullable -}}
	local.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
	{{- end}}

	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreign.{{$txt.Function.ForeignAssignment}} = local.{{$txt.Function.LocalAssignment}}
	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.{{$txt.Function.Name}}(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	{{if $txt.Function.UsesBytes -}}
	if 0 != bytes.Compare(check.{{$txt.Function.ForeignAssignment}}, foreign.{{$txt.Function.ForeignAssignment}}) {
	{{else -}}
	if check.{{$txt.Function.ForeignAssignment}} != foreign.{{$txt.Function.ForeignAssignment}} {
	{{end -}}
		t.Errorf("want: %v, got %v", foreign.{{$txt.Function.ForeignAssignment}}, check.{{$txt.Function.ForeignAssignment}})
	}

	slice := {{$txt.LocalTable.NameGo}}Slice{&local}
	if err = local.L.Load{{$txt.Function.Name}}(tx, false, (*[]*{{$txt.LocalTable.NameGo}})(&slice)); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$txt.Function.Name}} == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.{{$txt.Function.Name}} = nil
	if err = local.L.Load{{$txt.Function.Name}}(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$txt.Function.Name}} == nil {
		t.Error("struct should have been eager loaded")
	}
}

{{end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
