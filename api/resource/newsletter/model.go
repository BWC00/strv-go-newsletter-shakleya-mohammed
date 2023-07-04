package newsletter

import (
	"time"
)

type Newsletter struct {
	ID  		  uint32    `json:"newsletter_id"`
	EditorId      uint32    `json:"editor_id"`
	Name		  string    `json:"name" form:"required,max=255"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

type Newsletters []*Newsletter
