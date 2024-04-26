package model

type Status int

const (
	Publish Status = iota
	Draft
	Trash
)

func (s Status) String() string {
	return [...]string{"Publish", "Draft", "Trash"}[s]
}

func (s Status) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Status) UnmarshalText(data []byte) error {
	switch string(data) {
	case `"Publish"`:
		*s = Publish
	case `"Draft"`:
		*s = Draft
	case `"Trash"`:
		*s = Trash
	}
	return nil
}

type Article struct {
	ID        int    `json:"id,omitempty"`
	Title     string `json:"title,omitempty" validate:"required,min=20"`
	Content   string `json:"content,omitempty" validate:"required,min=200"`
	Category  string `json:"category,omitempty" validate:"required,min=3"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Status    string `json:"status,omitempty" validate:"required,oneof=Publish Draft Trash"`
}
