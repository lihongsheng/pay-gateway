package enum

type Operator string

const (
	OperatorEqual          Operator = "equal"
	OperatorNotEqual       Operator = "not_equal"
	OperatorLess           Operator = "less"
	OperatorLessOrEqual    Operator = "less_or_equal"
	OperatorGreater        Operator = "greater"
	OperatorGreaterOrEqual Operator = "greater_or_equal"
	OperatorIn             Operator = "in"
	OperatorNotIn          Operator = "not_in"
	OperatorIsNull         Operator = "is_null"
	OperatorIsNotNull      Operator = "is_not_null"
	OperatorContains       Operator = "contains"
	OperatorNotContains    Operator = "not_contains"
	OperatorBeginsWith     Operator = "begins_with"
	OperatorEndsWith       Operator = "ends_with"
)

func (o Operator) Label() string {
	switch o {
	case OperatorEqual:
		return "等于"
	case OperatorNotEqual:
		return "不等于"
	case OperatorLess:
		return "小于"
	case OperatorLessOrEqual:
		return "小于等于"
	case OperatorGreater:
		return "大于"
	case OperatorGreaterOrEqual:
		return "大于等于"
	case OperatorIn:
		return "包含"
	case OperatorNotIn:
		return "不包含"
	case OperatorIsNull:
		return "为空"
	case OperatorIsNotNull:
		return "不为空"
	case OperatorContains:
		return "包含字符"
	case OperatorNotContains:
		return "不包含字符"
	case OperatorBeginsWith:
		return "开头是"
	case OperatorEndsWith:
		return "结尾是"
	default:
		return ""
	}
}

func (o Operator) NeedsValue() bool {
	return o != OperatorIsNull && o != OperatorIsNotNull
}

func (o Operator) AcceptsArrayValue() bool {
	return o == OperatorIn || o == OperatorNotIn
}

func (o Operator) AcceptsNumericValue() bool {
	return o == OperatorLess || o == OperatorLessOrEqual ||
		o == OperatorGreater || o == OperatorGreaterOrEqual
}

func (o Operator) AcceptsStringValue() bool {
	return o == OperatorEqual || o == OperatorNotEqual ||
		o == OperatorContains || o == OperatorNotContains ||
		o == OperatorBeginsWith || o == OperatorEndsWith
}

func (o Operator) AcceptFieldTypes() []string {
	switch o {
	case OperatorIsNull, OperatorIsNotNull:
		return []string{"string", "integer", "double", "boolean", "date"}
	case OperatorIn, OperatorNotIn:
		return []string{"string", "integer", "double", "select"}
	case OperatorLess, OperatorLessOrEqual, OperatorGreater, OperatorGreaterOrEqual:
		return []string{"integer", "double"}
	case OperatorContains, OperatorNotContains, OperatorBeginsWith, OperatorEndsWith:
		return []string{"string"}
	default:
		return []string{"string", "integer", "double", "select"}
	}
}

func OperatorFromLegacy(legacy string) Operator {
	switch legacy {
	case "eq":
		return OperatorEqual
	case "ne":
		return OperatorNotEqual
	case "gt":
		return OperatorGreater
	case "gte", ">=":
		return OperatorGreaterOrEqual
	case "lt":
		return OperatorLess
	case "lte", "<=":
		return OperatorLessOrEqual
	case "is_empty":
		return OperatorIsNull
	default:
		if op := Operator(legacy); op.IsValid() {
			return op
		}
		return OperatorEqual
	}
}

func (o Operator) ToLegacy() string {
	switch o {
	case OperatorEqual:
		return "eq"
	case OperatorNotEqual:
		return "ne"
	case OperatorGreater:
		return "gt"
	case OperatorGreaterOrEqual:
		return "gte"
	case OperatorLess:
		return "lt"
	case OperatorLessOrEqual:
		return "lte"
	case OperatorIsNull:
		return "is_empty"
	case OperatorIsNotNull:
		return "is_empty"
	default:
		return string(o)
	}
}

func AllOperators() []Operator {
	return []Operator{
		OperatorEqual, OperatorNotEqual,
		OperatorLess, OperatorLessOrEqual,
		OperatorGreater, OperatorGreaterOrEqual,
		OperatorIn, OperatorNotIn,
		OperatorIsNull, OperatorIsNotNull,
		OperatorContains, OperatorNotContains,
		OperatorBeginsWith, OperatorEndsWith,
	}
}

func (o Operator) IsValid() bool {
	switch o {
	case OperatorEqual, OperatorNotEqual,
		OperatorLess, OperatorLessOrEqual,
		OperatorGreater, OperatorGreaterOrEqual,
		OperatorIn, OperatorNotIn,
		OperatorIsNull, OperatorIsNotNull,
		OperatorContains, OperatorNotContains,
		OperatorBeginsWith, OperatorEndsWith:
		return true
	default:
		return false
	}
}
