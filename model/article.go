package model

type Article struct {
	ID        int    `json:"id,omitempty"`
	Title     string `json:"title,omitempty" validate:"required,min=20"`
	Content   string `json:"content,omitempty" validate:"required,min=200"`
	Category  string `json:"category,omitempty" validate:"required,min=3"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Status    string `json:"status,omitempty" validate:"required,oneof=Publish Draft Trash"`
}
