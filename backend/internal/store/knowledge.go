package store

import (
	"curriculum/internal/model"
	"fmt"
	"time"
)

const kpFile = "knowledge_points.json"

func (s *Store) ListKnowledgePoints() ([]model.KnowledgePoint, error) {
	var kps []model.KnowledgePoint
	if err := s.loadData(kpFile, &kps); err != nil {
		return nil, err
	}
	if kps == nil {
		kps = []model.KnowledgePoint{}
	}
	return kps, nil
}

func (s *Store) GetKnowledgePoint(id string) (*model.KnowledgePoint, error) {
	kps, err := s.ListKnowledgePoints()
	if err != nil {
		return nil, err
	}
	for _, kp := range kps {
		if kp.ID == id {
			kpc := kp
			return &kpc, nil
		}
	}
	return nil, fmt.Errorf("knowledge point not found: %s", id)
}

func (s *Store) CreateKnowledgePoint(kp *model.KnowledgePoint) error {
	kps, err := s.ListKnowledgePoints()
	if err != nil {
		return err
	}
	now := time.Now()
	kp.CreatedAt = now
	kp.UpdatedAt = now
	if kp.EstimatedLessons == 0 {
		kp.EstimatedLessons = int(kp.Difficulty)
	}
	kps = append(kps, *kp)
	return s.saveData(kpFile, kps)
}

func (s *Store) UpdateKnowledgePoint(kp *model.KnowledgePoint) error {
	kps, err := s.ListKnowledgePoints()
	if err != nil {
		return err
	}
	kp.UpdatedAt = time.Now()
	found := false
	for i, item := range kps {
		if item.ID == kp.ID {
			kps[i] = *kp
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("knowledge point not found: %s", kp.ID)
	}
	return s.saveData(kpFile, kps)
}

func (s *Store) DeleteKnowledgePoint(id string) error {
	kps, err := s.ListKnowledgePoints()
	if err != nil {
		return err
	}
	newList := make([]model.KnowledgePoint, 0, len(kps))
	for _, kp := range kps {
		if kp.ID != id {
			newList = append(newList, kp)
		}
	}
	return s.saveData(kpFile, newList)
}

func (s *Store) TopologicalSort() (*model.TopoSortResult, error) {
	kps, err := s.ListKnowledgePoints()
	if err != nil {
		return nil, err
	}

	kpMap := make(map[string]*model.KnowledgePoint)
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	for i := range kps {
		id := kps[i].ID
		kpMap[id] = &kps[i]
		if _, ok := inDegree[id]; !ok {
			inDegree[id] = 0
		}
	}

	for _, kp := range kps {
		for _, pre := range kp.Prerequisites {
			if _, ok := kpMap[pre]; !ok {
				continue
			}
			graph[pre] = append(graph[pre], kp.ID)
			inDegree[kp.ID]++
		}
	}

	cycle := s.detectCycle(kpMap, graph)
	if cycle.HasCycle {
		return &model.TopoSortResult{Cycle: &cycle}, nil
	}

	queue := make([]string, 0)
	levels := make(map[string]int)

	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
			levels[id] = 1
		}
	}

	sorted := make([]string, 0, len(kps))
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		sorted = append(sorted, curr)

		for _, next := range graph[curr] {
			inDegree[next]--
			if levels[curr]+1 > levels[next] {
				levels[next] = levels[curr] + 1
			}
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	return &model.TopoSortResult{
		Sorted: sorted,
		Levels: levels,
	}, nil
}

func (s *Store) detectCycle(kpMap map[string]*model.KnowledgePoint, graph map[string][]string) model.CycleDetectionResult {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := make([]string, 0)
	var cycle []string

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		for _, next := range graph[node] {
			if !visited[next] {
				if dfs(next) {
					return true
				}
			} else if recStack[next] {
				startIdx := -1
				for i, p := range path {
					if p == next {
						startIdx = i
						break
					}
				}
				if startIdx >= 0 {
					cycle = append(cycle, path[startIdx:]...)
					cycle = append(cycle, next)
				}
				return true
			}
		}

		path = path[:len(path)-1]
		recStack[node] = false
		return false
	}

	for id := range kpMap {
		if !visited[id] {
			if dfs(id) {
				return model.CycleDetectionResult{HasCycle: true, Cycle: cycle}
			}
		}
	}

	return model.CycleDetectionResult{HasCycle: false}
}
