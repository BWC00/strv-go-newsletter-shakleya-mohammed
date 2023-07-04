package subscription

type Subscription struct {
	Email string        `json:"email" form:"required,email"`
	NewsletterID uint32 `json:"newsletter_id,string" form:"required"`
}