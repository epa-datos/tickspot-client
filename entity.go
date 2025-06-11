package tickspot

type TickEntry struct {
	ID     int     `json:"id,omitempty"`
	Date   string  `json:"date"`
	Hours  float64 `json:"hours"`
	Notes  string  `json:"notes"`
	TaskID int     `json:"task_id"`
	UserID int     `json:"user_id"`
}

type UsersTick struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
