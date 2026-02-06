package dto

type CreateTodoInput struct {
	Title   string `json:"title" binding:"required,min=2"` 
	Content string `json:"content"`
}

type UpdateTodoInput struct {
	Title   *string `json:"title" binding:"omitnil,min=2"` 
	Content *string `json:"content"` 
	Done    *bool   `json:"done"`
}