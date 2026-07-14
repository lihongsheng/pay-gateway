package domain

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/lihongsheng/pay-gateway/plugin/routing/dto"
	"github.com/lihongsheng/pay-gateway/plugin/routing/enum"
)

const maxDepth = 10

// RuleEngine 规则引擎，负责验证和执行规则
type RuleEngine struct {
	schemaService *RuleSchemaService
}

// NewRuleEngine 创建规则引擎
func NewRuleEngine(schemaService *RuleSchemaService) *RuleEngine {
	return &RuleEngine{schemaService: schemaService}
}

// Validate 验证规则 JSON 结构合法性
func (e *RuleEngine) Validate(ruleJson map[string]any, depth int) error {
	if depth > maxDepth {
		return fmt.Errorf("规则嵌套深度超过限制（最大 %d 层）", maxDepth)
	}

	// 顶层必须有 condition 和 rules
	conditionVal, hasCondition := ruleJson["condition"]
	rulesVal, hasRules := ruleJson["rules"]
	if !hasCondition || !hasRules {
		return fmt.Errorf("规则必须包含 condition 和 rules 字段")
	}

	condition := strings.ToUpper(fmt.Sprintf("%v", conditionVal))
	if condition != "AND" && condition != "OR" {
		return fmt.Errorf("condition 必须为 AND 或 OR")
	}

	rules, ok := rulesVal.([]any)
	if !ok || len(rules) == 0 {
		return fmt.Errorf("rules 必须为非空数组")
	}

	validFieldIds := e.schemaService.GetValidFieldIds()

	for i, rule := range rules {
		ruleMap, ok := rule.(map[string]any)
		if !ok {
			return fmt.Errorf("rules[%d] 必须为对象", i)
		}

		if _, hasSubCondition := ruleMap["condition"]; hasSubCondition {
			// 递归验证子条件组
			if err := e.Validate(ruleMap, depth+1); err != nil {
				return err
			}
		} else {
			// 验证叶子条件
			if err := e.validateLeaf(ruleMap, i, validFieldIds); err != nil {
				return err
			}
		}
	}

	return nil
}

// Evaluate 执行规则，对每个支付账号注入账户统计数据后评估，返回匹配的账号列表
func (e *RuleEngine) Evaluate(ruleJson map[string]any, context dto.RuleContext) []dto.RuleAccountContext {
	if len(context.Accounts) == 0 {
		return nil
	}

	var matched []dto.RuleAccountContext
	for _, account := range context.Accounts {
		// 将业务数据与当前账户数据合并为评估上下文
		evalCtx := e.buildEvalContext(context.Data, account)
		if e.evaluateRule(ruleJson, evalCtx) {
			matched = append(matched, account)
		}
	}
	return matched
}

// evaluateRule 递归评估规则，返回 bool
func (e *RuleEngine) evaluateRule(ruleJson map[string]any, context map[string]any) bool {
	conditionRaw, _ := ruleJson["condition"]
	conditionStr, _ := conditionRaw.(string)
	condition := enum.ConditionLogicFromWithDefault(&conditionStr)

	rulesVal, _ := ruleJson["rules"]
	rules, _ := rulesVal.([]any)

	var results []bool
	for _, rule := range rules {
		ruleMap, ok := rule.(map[string]any)
		if !ok {
			continue
		}
		if _, hasSubCondition := ruleMap["condition"]; hasSubCondition {
			results = append(results, e.evaluateRule(ruleMap, context))
		} else {
			results = append(results, e.evaluateLeaf(ruleMap, context))
		}
	}

	return condition.Evaluate(results)
}

// buildEvalContext 将业务数据与账户统计数据合并为评估上下文
// 账户字段 account.[].day_order / account.[].day_amount 映射为 account_day_order / account_day_amount
func (e *RuleEngine) buildEvalContext(data map[string]any, account dto.RuleAccountContext) map[string]any {
	ctx := make(map[string]any, len(data)+2)

	// 复制业务数据
	for k, v := range data {
		ctx[k] = v
	}

	// 注入账户统计数据
	ctx["account_day_order"] = account.DayOrder
	ctx["account_day_amount"] = account.DayAmount
	ctx["account_no"] = account.AccountNo

	return ctx
}

// resolveFieldId 将 schema 中的字段 ID 映射为上下文中的 key
// account.[].day_order → account_day_order
func resolveFieldId(fieldId string) string {
	if strings.HasPrefix(fieldId, "account.[].") {
		return "account_" + strings.TrimPrefix(fieldId, "account.[].")
	}
	return fieldId
}

// validateLeaf 验证叶子条件节点
func (e *RuleEngine) validateLeaf(rule map[string]any, index int, validFieldIds []string) error {
	idVal, _ := rule["id"]
	id, _ := idVal.(string)
	if id == "" {
		return fmt.Errorf("rules[%d] 缺少 id 字段", index)
	}

	if !containsString(validFieldIds, id) {
		return fmt.Errorf("rules[%d].id='%s' 不是合法的字段", index, id)
	}

	opVal, _ := rule["operator"]
	op, _ := opVal.(string)
	if op == "" {
		return fmt.Errorf("rules[%d] 缺少 operator 字段", index)
	}

	operator := enum.Operator(op)
	if !operator.IsValid() {
		return fmt.Errorf("rules[%d].operator='%s' 不是合法的操作符", index, op)
	}

	// 验证操作符是否适用于该字段
	fieldDef := e.schemaService.GetFieldById(id)
	if fieldDef != nil && len(fieldDef.Operators) > 0 {
		if !containsString(fieldDef.Operators, op) {
			return fmt.Errorf("操作符 '%s' 不适用于字段 '%s'", op, id)
		}
	}

	// is_null / is_not_null 不需要 value
	if !operator.NeedsValue() {
		return nil
	}

	// in / not_in 允许字符串，后端会处理
	if operator.AcceptsArrayValue() {
		return nil
	}

	// 数值操作符需要 value 可转为数字
	if operator.AcceptsNumericValue() {
		valueVal, hasValue := rule["value"]
		if hasValue {
			valueStr := fmt.Sprintf("%v", valueVal)
			if valueStr != "" {
				if _, err := strconv.ParseFloat(valueStr, 64); err != nil {
					return fmt.Errorf("rules[%d].value 必须为数字（字段: %s）", index, id)
				}
			}
		}
	}

	return nil
}

// evaluateLeaf 执行叶子条件比较
func (e *RuleEngine) evaluateLeaf(rule map[string]any, context map[string]any) bool {
	fieldId, _ := rule["id"].(string)
	// 将 schema 字段 ID 映射为上下文 key
	contextKey := resolveFieldId(fieldId)

	opStr, _ := rule["operator"].(string)
	operator := enum.Operator(opStr)
	if !operator.IsValid() {
		operator = enum.OperatorEqual
	}
	expected := rule["value"]
	actual, _ := context[contextKey]

	switch operator {
	case enum.OperatorEqual:
		return compareEqual(actual, expected)
	case enum.OperatorNotEqual:
		return !compareEqual(actual, expected)
	case enum.OperatorLess:
		return compareLess(actual, expected)
	case enum.OperatorLessOrEqual:
		return compareLessOrEqual(actual, expected)
	case enum.OperatorGreater:
		return compareGreater(actual, expected)
	case enum.OperatorGreaterOrEqual:
		return compareGreaterOrEqual(actual, expected)
	case enum.OperatorIn:
		return compareIn(actual, expected)
	case enum.OperatorNotIn:
		return !compareIn(actual, expected)
	case enum.OperatorIsNull:
		return actual == nil || toString(actual) == ""
	case enum.OperatorIsNotNull:
		return actual != nil && toString(actual) != ""
	case enum.OperatorContains:
		return strings.Contains(toString(actual), toString(expected))
	case enum.OperatorNotContains:
		return !strings.Contains(toString(actual), toString(expected))
	case enum.OperatorBeginsWith:
		return strings.HasPrefix(toString(actual), toString(expected))
	case enum.OperatorEndsWith:
		return strings.HasSuffix(toString(actual), toString(expected))
	default:
		return false
	}
}

// ─── 比较函数 ───────────────────────────────────────────

func compareEqual(actual, expected any) bool {
	actualStr := toString(actual)
	if actual == nil || actualStr == "" {
		expectedStr := toString(expected)
		return expected == nil || expectedStr == ""
	}

	actualNum, actualIsNum := toFloat64(actual)
	expectedNum, expectedIsNum := toFloat64(expected)
	if actualIsNum && expectedIsNum {
		return actualNum == expectedNum
	}

	return strings.ToUpper(actualStr) == strings.ToUpper(toString(expected))
}

func compareLess(actual, expected any) bool {
	a, aOk := toFloat64(actual)
	b, bOk := toFloat64(expected)
	return aOk && bOk && a < b
}

func compareLessOrEqual(actual, expected any) bool {
	a, aOk := toFloat64(actual)
	b, bOk := toFloat64(expected)
	return aOk && bOk && a <= b
}

func compareGreater(actual, expected any) bool {
	a, aOk := toFloat64(actual)
	b, bOk := toFloat64(expected)
	return aOk && bOk && a > b
}

func compareGreaterOrEqual(actual, expected any) bool {
	a, aOk := toFloat64(actual)
	b, bOk := toFloat64(expected)
	return aOk && bOk && a >= b
}

func compareIn(actual, expected any) bool {
	actualStr := toString(actual)
	if actual == nil || actualStr == "" {
		return false
	}

	var values []string
	switch v := expected.(type) {
	case []any:
		for _, item := range v {
			values = append(values, toString(item))
		}
	case string:
		for _, s := range strings.Split(v, ",") {
			values = append(values, strings.TrimSpace(s))
		}
	default:
		values = []string{toString(expected)}
	}

	actualNum, actualIsNum := toFloat64(actual)
	for _, v := range values {
		vNum, vIsNum := toFloat64(v)
		if actualIsNum && vIsNum {
			if actualNum == vNum {
				return true
			}
		} else if strings.ToUpper(actualStr) == strings.ToUpper(v) {
			return true
		}
	}

	return false
}

// ─── 工具函数 ───────────────────────────────────────────

func toString(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == math.Trunc(val) {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case string:
		f, err := strconv.ParseFloat(val, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func containsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
