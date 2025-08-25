package skills

import (
	"errors"
	"sort"
	"strings"
)

func normalize(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

type Registry interface {
	Register(Skill) error
	Get(name string) (Skill, bool)
	List() []string
}

type memoryRegistry struct {
	store map[string]Skill
}

func (m *memoryRegistry) Register(skill Skill) error {
	if skill == nil {
		return errors.New("skill is nil")
	}
	name := normalize(skill.Name())
	if name == "" {
		return errors.New("skill name is empty")
	}
	_, exists := m.store[name]
	if exists {
		return ErrDuplicateSkill
	}
	m.store[name] = skill
	return nil
}

func (m *memoryRegistry) Get(name string) (Skill, bool) {
	key := normalize(name)
	if key == "" {
		return nil, false
	}
	skill, exists := m.store[key]
	return skill, exists
}

func (m *memoryRegistry) List() []string {
	names := make([]string, 0, len(m.store))
	for name := range m.store {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func NewRegistry() Registry {
	return &memoryRegistry{
		store: make(map[string]Skill),
	}
}
