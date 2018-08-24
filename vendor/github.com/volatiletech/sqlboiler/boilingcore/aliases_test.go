package boilingcore

import (
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/drivers"
)

func TestAliasesTables(t *testing.T) {
	t.Parallel()

	tables := []drivers.Table{
		{
			Name: "videos",
			Columns: []drivers.Column{
				{Name: "id"},
				{Name: "name"},
			},
		},
	}

	t.Run("Automatic", func(t *testing.T) {
		expect := TableAlias{
			UpPlural:     "Videos",
			UpSingular:   "Video",
			DownPlural:   "videos",
			DownSingular: "video",
			Columns: map[string]string{
				"id":   "ID",
				"name": "Name",
			},
			Relationships: make(map[string]RelationshipAlias),
		}

		a := Aliases{}
		FillAliases(&a, tables)

		if got := a.Tables["videos"]; !reflect.DeepEqual(expect, got) {
			t.Errorf("it should fill in the blanks: %#v", got)
		}
	})

	t.Run("UserOverride", func(t *testing.T) {
		expect := TableAlias{
			UpPlural:     "NotVideos",
			UpSingular:   "NotVideo",
			DownPlural:   "notVideos",
			DownSingular: "notVideo",
			Columns: map[string]string{
				"id":   "NotID",
				"name": "NotName",
			},
			Relationships: make(map[string]RelationshipAlias),
		}

		a := Aliases{}
		a.Tables = map[string]TableAlias{"videos": expect}
		FillAliases(&a, tables)

		if !reflect.DeepEqual(expect, a.Tables["videos"]) {
			t.Error("it should not alter things that were specified by user")
		}
	})
}

func TestAliasesRelationships(t *testing.T) {
	t.Parallel()

	tables := []drivers.Table{
		{
			Name: "videos",
			Columns: []drivers.Column{
				{Name: "id"},
				{Name: "name"},
			},
			FKeys: []drivers.ForeignKey{
				{
					Name:          "fkey_1",
					Table:         "videos",
					Column:        "user_id",
					ForeignTable:  "users",
					ForeignColumn: "id",
				},
				{
					Name:          "fkey_2",
					Table:         "videos",
					Column:        "publisher_id",
					ForeignTable:  "users",
					ForeignColumn: "id",
				},
				{
					Name:          "fkey_3",
					Table:         "videos",
					Column:        "one_id",
					Unique:        true,
					ForeignTable:  "ones",
					ForeignColumn: "one",
				},
			},
		},
	}

	t.Run("Automatic", func(t *testing.T) {
		expect1 := RelationshipAlias{
			Local:   "Videos",
			Foreign: "User",
		}
		expect2 := RelationshipAlias{
			Local:   "PublisherVideos",
			Foreign: "Publisher",
		}
		expect3 := RelationshipAlias{
			Local:   "Video",
			Foreign: "One",
		}

		a := Aliases{}
		FillAliases(&a, tables)

		table := a.Tables["videos"]
		if got := table.Relationships["fkey_1"]; !reflect.DeepEqual(expect1, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := table.Relationships["fkey_2"]; !reflect.DeepEqual(expect2, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := table.Relationships["fkey_3"]; !reflect.DeepEqual(expect3, got) {
			t.Errorf("bad values: %#v", got)
		}
	})

	t.Run("UserOverride", func(t *testing.T) {
		expect1 := RelationshipAlias{
			Local:   "Videos",
			Foreign: "TheUser",
		}
		expect2 := RelationshipAlias{
			Local:   "PublishedVideos",
			Foreign: "Publisher",
		}
		expect3 := RelationshipAlias{
			Local:   "AwesomeOneVideo",
			Foreign: "TheOne",
		}

		a := Aliases{
			Tables: map[string]TableAlias{
				"videos": {
					Relationships: map[string]RelationshipAlias{
						"fkey_1": {Foreign: "TheUser"},
						"fkey_2": {Local: "PublishedVideos"},
						"fkey_3": {Local: "AwesomeOneVideo", Foreign: "TheOne"},
					},
				},
			},
		}
		FillAliases(&a, tables)

		table := a.Tables["videos"]
		if got := table.Relationships["fkey_1"]; !reflect.DeepEqual(expect1, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := table.Relationships["fkey_2"]; !reflect.DeepEqual(expect2, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := table.Relationships["fkey_3"]; !reflect.DeepEqual(expect3, got) {
			t.Errorf("bad values: %#v", got)
		}
	})
}

func TestAliasesRelationshipsJoinTable(t *testing.T) {
	t.Parallel()

	tables := []drivers.Table{
		{
			Name: "videos",
		},
		{
			Name: "tags",
		},
		{
			Name:        "video_tags",
			IsJoinTable: true,
			FKeys: []drivers.ForeignKey{
				{
					Name:          "fk_video_id",
					Table:         "video_tags",
					Column:        "video_id",
					ForeignTable:  "videos",
					ForeignColumn: "id",
				},
				{
					Name:          "fk_tag_id",
					Table:         "video_tags",
					Column:        "tags_id",
					ForeignTable:  "tags",
					ForeignColumn: "id",
				},
			},
		},
	}

	t.Run("Automatic", func(t *testing.T) {
		expect1 := RelationshipAlias{
			Local:   "Tags",
			Foreign: "Videos",
		}
		expect2 := RelationshipAlias{
			Local:   "Videos",
			Foreign: "Tags",
		}

		a := Aliases{}
		FillAliases(&a, tables)

		table := a.Tables["video_tags"]
		if got := table.Relationships["fk_video_id"]; !reflect.DeepEqual(expect1, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := table.Relationships["fk_tag_id"]; !reflect.DeepEqual(expect2, got) {
			t.Errorf("bad values: %#v", got)
		}
	})

	t.Run("UserOverride", func(t *testing.T) {
		expect1 := RelationshipAlias{
			Local:   "NotTags",
			Foreign: "NotVideos",
		}
		expect2 := RelationshipAlias{
			Local:   "NotVideos",
			Foreign: "NotTags",
		}

		a := Aliases{
			Tables: map[string]TableAlias{
				"video_tags": {
					Relationships: map[string]RelationshipAlias{
						"fk_video_id": {Local: "NotTags", Foreign: "NotVideos"},
					},
				},
			},
		}
		FillAliases(&a, tables)

		table := a.Tables["video_tags"]
		if got := table.Relationships["fk_video_id"]; !reflect.DeepEqual(expect1, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := table.Relationships["fk_tag_id"]; !reflect.DeepEqual(expect2, got) {
			t.Errorf("bad values: %#v", got)
		}
	})
}

func TestAliasHelpers(t *testing.T) {
	t.Parallel()

	a := Aliases{
		Tables: map[string]TableAlias{
			"videos": {
				UpPlural: "Videos",
				Columns: map[string]string{
					"id": "NotID",
				},
				Relationships: map[string]RelationshipAlias{
					"fk_user_id": {Local: "Videos", Foreign: "User"},
				},
			},
			"video_tags": {
				Relationships: map[string]RelationshipAlias{
					"fk_video_id": {Local: "NotTags", Foreign: "NotVideos"},
				},
			},
		},
	}

	if got := a.Table("videos").UpPlural; got != "Videos" {
		t.Error("videos upPlural wrong:", got)
	}

	if got := a.Table("videos").Relationship("fk_user_id"); got.Local != "Videos" {
		t.Error("videos relationship wrong:", got)
	}

	got := a.ManyRelationship("videos", "fk_user_id", "video_tags", "fk_video_id")
	if got.Foreign != "NotVideos" {
		t.Error("relationship wrong:", got)
	}

	got = a.ManyRelationship("videos", "fk_user_id", "", "")
	if got.Foreign != "User" {
		t.Error("relationship wrong:", got)
	}
}
