package repository

import (
	todo "awesomeProject"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title , description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", itemsListsTable)

	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id 
    INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, itemsListsTable, usersListsTable)
	if err := r.db.Select(&items, query, listId, userId); err != nil {
		fmt.Print(err, items)
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var items todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id 
    INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, itemsListsTable, usersListsTable)
	if err := r.db.Get(&items, query, itemId, userId); err != nil {
		fmt.Print(err, items)
		return items, err
	}
	return items, nil
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {

	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li,%s ul WHERE ti.id = li.item_id and ul.list_id = li.list_id and ti.id=$1 and ul.user_id=$2`,
		todoItemsTable, itemsListsTable, usersListsTable)
	if _, err := r.db.Exec(query, itemId, userId); err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}
func (r *TodoItemPostgres) Update(userId, id int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s ti SET %s FROM %s ul, %s li WHERE ti.id = li.item_id AND ul.list_id = li.list_id AND ul.user_id=$%d AND li.list_id=$%d",
		todoItemsTable, setQuery, usersListsTable, itemsListsTable, argId, argId+1)

	args = append(args, userId, id)

	_, err := r.db.Exec(query, args...)

	return err
}
