package store

import (
	"curriculum/internal/model"
	"fmt"
	"time"
)

const revisionFile = "revision_requests.json"

func (s *Store) ListRevisionRequests() ([]model.RevisionRequest, error) {
	var requests []model.RevisionRequest
	if err := s.loadData(revisionFile, &requests); err != nil {
		return nil, err
	}
	if requests == nil {
		requests = []model.RevisionRequest{}
	}
	return requests, nil
}

func (s *Store) GetRevisionRequest(id string) (*model.RevisionRequest, error) {
	requests, err := s.ListRevisionRequests()
	if err != nil {
		return nil, err
	}
	for _, r := range requests {
		if r.ID == id {
			rc := r
			return &rc, nil
		}
	}
	return nil, fmt.Errorf("revision request not found: %s", id)
}

func (s *Store) CreateRevisionRequest(req *model.RevisionRequest) error {
	requests, err := s.ListRevisionRequests()
	if err != nil {
		return err
	}
	now := time.Now()
	req.CreatedAt = now
	req.Status = model.RevisionPending
	requests = append(requests, *req)
	return s.saveData(revisionFile, requests)
}

func (s *Store) UpdateRevisionRequest(req *model.RevisionRequest) error {
	requests, err := s.ListRevisionRequests()
	if err != nil {
		return err
	}
	found := false
	for i, item := range requests {
		if item.ID == req.ID {
			req.CreatedAt = item.CreatedAt
			requests[i] = *req
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("revision request not found: %s", req.ID)
	}
	return s.saveData(revisionFile, requests)
}
