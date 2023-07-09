package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"real-state-finder/pkg/entities"

	"github.com/pkg/errors"
)

type Storage interface {
	Load() error
	Dump() error
	Get(id string) (entities.RealState, bool)
	Save(rs entities.RealState) error
	Exists(id string) bool
	GetList() []entities.RealState
	ResetNew()
}

type storage struct {
	filepath   string
	realStates map[string]entities.RealState
}

func NewStorage(filepath string) Storage {
	return &storage{filepath: filepath}
}

func (s *storage) Load() error {
	rs := make([]entities.RealState, 0)
	mp := make(map[string]entities.RealState, 0)
	bArr, err := ioutil.ReadFile(s.filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bArr, &rs)
	if err != nil {
		return err
	}

	for _, r := range rs {
		mp[r.Id] = r
	}

	s.realStates = mp
	return nil
}

func (s *storage) Dump() error {
	bArr, err := json.Marshal(s.GetList())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.filepath, bArr, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) Get(id string) (entities.RealState, bool) {
	if m, ok := s.realStates[id]; ok {
		return m, true
	}

	return entities.RealState{}, false
}

func (s *storage) Save(r entities.RealState) error {
	if s.Exists(r.Id) {
		return errors.Errorf("Real state with id %s already exists", r.Id)
	}
	s.realStates[r.Id] = r
	return nil
}

func (s *storage) Exists(id string) bool {
	_, ok := s.realStates[id]
	return ok
}

func (s *storage) GetList() []entities.RealState {
	rs := make([]entities.RealState, 0)
	for k := range s.realStates {
		rs = append(rs, s.realStates[k])
		fmt.Println(fmt.Printf("%+v", s.realStates[k]))
	}

	return rs
}

func (s *storage) ResetNew() {
	for k := range s.realStates {
		rs := s.realStates[k]
		rs.IsNew = false
		s.realStates[k] = rs
	}
}
