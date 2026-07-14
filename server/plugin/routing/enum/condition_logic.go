package enum

import "strings"

type ConditionLogic string

const (
	ConditionLogicAnd ConditionLogic = "AND"
	ConditionLogicOr  ConditionLogic = "OR"
)

func (c ConditionLogic) Label() string {
	switch c {
	case ConditionLogicAnd:
		return "且"
	case ConditionLogicOr:
		return "或"
	default:
		return ""
	}
}

func (c ConditionLogic) Evaluate(results []bool) bool {
	if c == ConditionLogicAnd {
		for _, r := range results {
			if !r {
				return false
			}
		}
		return true
	}
	for _, r := range results {
		if r {
			return true
		}
	}
	return false
}

func ConditionLogicFromWithDefault(value *string) ConditionLogic {
	if value == nil {
		return ConditionLogicAnd
	}
	v := strings.ToUpper(*value)
	switch ConditionLogic(v) {
	case ConditionLogicAnd, ConditionLogicOr:
		return ConditionLogic(v)
	default:
		return ConditionLogicAnd
	}
}
