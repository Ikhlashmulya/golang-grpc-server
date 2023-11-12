package model

type CreateTodoRequest struct {
	Name string `json:"name,omitempty"`
}

type TodoResponse struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
