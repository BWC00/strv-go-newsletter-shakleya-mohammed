package subscription

type Subscription struct {
	Email string `json:"email" form:"required,email"`
	ID uint32 	 `json:"id,string" form:"required"`
}