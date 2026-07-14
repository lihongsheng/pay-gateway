package enum

type RuleStatus int

const (
	RuleStatusInactive RuleStatus = iota + 1
	RuleStatusActive
)

func (s RuleStatus) Label() string {
	switch s {
	case RuleStatusActive:
		return "激活中"
	case RuleStatusInactive:
		return "未激活"
	default:
		return ""
	}
}

func (s RuleStatus) TagType() string {
	switch s {
	case RuleStatusActive:
		return "success"
	case RuleStatusInactive:
		return "info"
	default:
		return ""
	}
}

func (s RuleStatus) String() string {
	switch s {
	case RuleStatusActive:
		return "active"
	case RuleStatusInactive:
		return "inactive"
	default:
		return ""
	}
}

func AllRuleStatuses() []RuleStatus {
	return []RuleStatus{RuleStatusInactive, RuleStatusActive}
}

func RuleStatusFromString(v string) RuleStatus {
	switch v {
	case "active":
		return RuleStatusActive
	case "inactive":
		return RuleStatusInactive
	default:
		return RuleStatusInactive
	}
}
