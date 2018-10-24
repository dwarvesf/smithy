package domain

// Group .
type Group struct {
	Model
	Name        string       `yaml:"name" json:"name"`
	Description string       `yaml:"description" json:"description"`
	Role        string       `yaml:"role" json:"role"`
	Users       []User       `gorm:"many2many:user_groups;" yaml:"-" json:"-"`
	Permissions []Permission `yaml:"-" json:"-"`
}
