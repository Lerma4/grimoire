package service

import (
	"context"

	"github.com/Lerma4/grimoire/internal/domain"
)

// SearchResult groups matching tasks and notes for a query.
type SearchResult struct {
	Query string
	Tasks []domain.Task
	Notes []domain.Note
}

// SearchService finds tasks and notes matching free-text queries.
type SearchService struct {
	tasks *TaskService
	notes *NoteService
}

// NewSearchService builds a SearchService.
func NewSearchService(tasks *TaskService, notes *NoteService) *SearchService {
	return &SearchService{tasks: tasks, notes: notes}
}

// Search runs the query across tasks and notes and returns both buckets.
func (s *SearchService) Search(ctx context.Context, query string) (SearchResult, error) {
	res := SearchResult{Query: query}
	if query == "" {
		return res, nil
	}
	tasks, err := s.tasks.List(ctx, TaskFilter{Search: query})
	if err != nil {
		return res, err
	}
	notes, err := s.notes.List(ctx, NoteFilter{Search: query})
	if err != nil {
		return res, err
	}
	res.Tasks = tasks
	res.Notes = notes
	return res, nil
}
