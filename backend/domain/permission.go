package domain

// Permission .
type Permission struct {
	Model
	DatabaseName string `yaml:"database" json:"database"`
	TableName    string `yaml:"table" json:"table"`
	Select       bool   `yaml:"select" json:"select"`
	Insert       bool   `yaml:"insert" json:"insert"`
	Update       bool   `yaml:"update" json:"update"`
	Delete       bool   `yaml:"delete" json:"delete"`
	UserID       UUID   `yaml:"user_id" json:"user_id"`
	GroupID      UUID   `yaml:"group_id" json:"group_id"`
}

// AND .
func (p Permission) AND(q Permission) Permission {
	p.Select = p.Select && q.Select
	p.Insert = p.Insert && q.Insert
	p.Update = p.Update && q.Update
	p.Delete = p.Delete && q.Delete
	return p
}
