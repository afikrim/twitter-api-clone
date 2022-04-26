package domains

type Todo struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateTodoDto struct {
	Title string `json:"title"`
}

type UpdateTodoDto struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}
