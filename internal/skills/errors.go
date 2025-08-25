package skills

import "errors"

var (
	ErrDuplicateSkill = errors.New("duplicate skill")
	ErrSkillNotFound  = errors.New("skill not found")
)
