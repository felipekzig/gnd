package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Task      string
	Pings     uint16
	Done      bool      `gorm:"default:false"`
	DueDate   time.Time `gorm:"default:null"`
	CreatedAt int
	UpdatedAt int
	DeletedAt int
}

type TaskService interface {
	Add(task string, dueDate time.Time) (Task, error)
	Ping(taskID uint) (Task, error)
	List(all bool) ([]Task, error)
	Complete(taskID uint) (Task, error)
}

type taskService struct {
	db *gorm.DB
}

func (s *taskService) Add(task string, dueDate time.Time) (Task, error) {
	e := &Task{Task: task, DueDate: dueDate}
	r := s.db.Create(&e)
	if r.Error != nil {
		return Task{}, r.Error
	}

	return *e, nil
}

func (s *taskService) Ping(taskID uint) (Task, error) {
	r := s.db.Model(&Task{}).Where("id = ?", taskID).Update("pings", gorm.Expr("pings + 1"))
	if r.Error != nil {
		return Task{}, r.Error
	}

	e := &Task{}
	r = s.db.First(&e, taskID)
	if r.Error != nil {
		return Task{}, r.Error
	}

	if r.RowsAffected == 0 {
		return Task{}, fmt.Errorf("no task found with id <%d>", taskID)
	}

	return *e, nil
}

func (s *taskService) List(all bool) ([]Task, error) {
	var tasks []Task

	q := s.db.Order("due_date DESC, pings DESC")

	if !all {
		q = q.Where("done IS FALSE")
	}

	r := q.Find(&tasks)
	if r.Error != nil {
		return nil, r.Error
	}

	return tasks, nil
}

func (s *taskService) Complete(taskID uint) (Task, error) {
	r := s.db.Model(&Task{}).Where("id = ?", taskID).Update("done", 1)
	if r.Error != nil {
		return Task{}, r.Error
	}

	if r.RowsAffected == 0 {
		return Task{}, fmt.Errorf("no task found with id <%d>", taskID)
	}

	e := &Task{}
	r = s.db.First(&e, taskID)
	if r.Error != nil {
		return Task{}, r.Error
	}

	return *e, nil
}

func NewTaskService(db *gorm.DB) TaskService {
	db.AutoMigrate(&Task{})

	return &taskService{db}
}
