package service

import (
	todo "awesomeProject"
	"awesomeProject/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

// Создание элементов списка
func (r *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := r.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return r.repo.Create(listId, item)
}

func (r *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return r.repo.GetAll(userId, listId)
}

func (r *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return r.repo.GetById(userId, itemId)
}
func (r *TodoItemService) Delete(userId, listId int) error {
	return r.repo.Delete(userId, listId)
}

func (r *TodoItemService) Update(userId, id int, input todo.UpdateItemInput) error {
	return r.repo.Update(userId, id, input)
}
