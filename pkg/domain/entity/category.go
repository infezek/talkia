package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Bots      []*Bot    `json:"bots,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}

func NewCategory(id uuid.UUID, name string, active bool, createdAt *time.Time) *Category {
	if id == uuid.Nil {
		id, _ = uuid.NewV7()
	}
	if createdAt == nil {
		now := time.Now()
		createdAt = &now
	}
	return &Category{
		ID:        id,
		Name:      name,
		Active:    active,
		CreatedAt: *createdAt,
	}
}

func (c *Category) Update(name string, active bool) {
	c.Name = name
	c.Active = active
	c.UpdateAt = time.Now()
}

func (c *Category) AddBots(bots []*Bot) {
	c.Bots = append(c.Bots, bots...)
}
