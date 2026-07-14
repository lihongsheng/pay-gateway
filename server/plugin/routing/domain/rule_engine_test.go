package domain

import (
	"testing"

	"github.com/lihongsheng/pay-gateway/plugin/routing/dto"
)

func newTestEngine() *RuleEngine {
	return NewRuleEngine(NewRuleSchemaService())
}

// 辅助：构建 RuleContext
func makeContext(data map[string]any, accounts []dto.RuleAccountContext) dto.RuleContext {
	return dto.RuleContext{Data: data, Accounts: accounts}
}

// 辅助：默认测试账号
func defaultAccounts() []dto.RuleAccountContext {
	return []dto.RuleAccountContext{
		{AccountNo: "ACC001", DayAmount: 5000.00, DayOrder: 10},
		{AccountNo: "ACC002", DayAmount: 15000.50, DayOrder: 50},
		{AccountNo: "ACC003", DayAmount: 500.00, DayOrder: 2},
	}
}

// ─── Validate 测试 ───────────────────────────────────────────

func TestValidate_ValidSimpleRule(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{
				"id":       "currency",
				"operator": "equal",
				"value":    "USD",
			},
		},
	}
	if err := engine.Validate(rule, 0); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingCondition(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"rules": []any{},
	}
	if err := engine.Validate(rule, 0); err == nil {
		t.Fatal("expected error for missing condition")
	}
}

func TestValidate_MissingRules(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
	}
	if err := engine.Validate(rule, 0); err == nil {
		t.Fatal("expected error for missing rules")
	}
}

func TestValidate_InvalidCondition(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "XOR",
		"rules":     []any{map[string]any{"id": "currency", "operator": "equal", "value": "USD"}},
	}
	if err := engine.Validate(rule, 0); err == nil {
		t.Fatal("expected error for invalid condition")
	}
}

func TestValidate_EmptyRules(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules":     []any{},
	}
	if err := engine.Validate(rule, 0); err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestValidate_InvalidFieldId(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "nonexistent", "operator": "equal", "value": "x"},
		},
	}
	if err := engine.Validate(rule, 0); err == nil {
		t.Fatal("expected error for invalid field id")
	}
}

func TestValidate_InvalidOperator(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "invalid_op", "value": "USD"},
		},
	}
	if err := engine.Validate(rule, 0); err == nil {
		t.Fatal("expected error for invalid operator")
	}
}

func TestValidate_OperatorNotApplicableToField(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "greater", "value": "10"},
		},
	}
	if err := engine.Validate(rule, 0); err == nil {
		t.Fatal("expected error for operator not applicable to field")
	}
}

func TestValidate_NumericOperatorWithNonNumericValue(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "order_amount", "operator": "greater", "value": "abc"},
		},
	}
	if err := engine.Validate(rule, 0); err == nil {
		t.Fatal("expected error for non-numeric value with numeric operator")
	}
}

func TestValidate_IsNullNoValueRequired(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "is_null"},
		},
	}
	if err := engine.Validate(rule, 0); err != nil {
		t.Fatalf("is_null should not require value, got: %v", err)
	}
}

func TestValidate_NestedRule(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{
				"condition": "OR",
				"rules": []any{
					map[string]any{"id": "currency", "operator": "equal", "value": "USD"},
					map[string]any{"id": "currency", "operator": "equal", "value": "EUR"},
				},
			},
			map[string]any{"id": "order_amount", "operator": "greater", "value": "100"},
		},
	}
	if err := engine.Validate(rule, 0); err != nil {
		t.Fatalf("expected no error for nested rule, got: %v", err)
	}
}

func TestValidate_MaxDepthExceeded(t *testing.T) {
	engine := newTestEngine()
	inner := map[string]any{
		"id": "currency", "operator": "equal", "value": "USD",
	}
	for i := 0; i < maxDepth+2; i++ {
		inner = map[string]any{
			"condition": "AND",
			"rules":     []any{inner},
		}
	}
	if err := engine.Validate(inner, 0); err == nil {
		t.Fatal("expected error for max depth exceeded")
	}
}

// ─── Evaluate 测试 ───────────────────────────────────────────

func TestEvaluate_EqualString(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "equal", "value": "USD"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_EqualStringCaseInsensitive(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "equal", "value": "usd"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_NotEqual(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "not_equal", "value": "EUR"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_Less(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "order_amount", "operator": "less", "value": "200"},
		},
	}
	ctx := makeContext(map[string]any{"order_amount": "100"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_GreaterOrEqual(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "order_amount", "operator": "greater_or_equal", "value": "100"},
		},
	}
	ctx := makeContext(map[string]any{"order_amount": "100"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_InWithStringArray(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "in", "value": []any{"USD", "EUR"}},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_InWithCommaString(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "in", "value": "USD,EUR"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "EUR"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_NotIn(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "not_in", "value": []any{"USD", "EUR"}},
		},
	}
	ctx := makeContext(map[string]any{"currency": "GBP"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_IsNull(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "is_null"},
		},
	}
	ctx := makeContext(map[string]any{"currency": nil}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_IsNullEmptyString(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "is_null"},
		},
	}
	ctx := makeContext(map[string]any{"currency": ""}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_IsNotNull(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "is_not_null"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_Contains(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "merchant_region", "operator": "contains", "value": "北京"},
		},
	}
	ctx := makeContext(map[string]any{"merchant_region": "北京市朝阳区"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_NotContains(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "merchant_region", "operator": "not_contains", "value": "上海"},
		},
	}
	ctx := makeContext(map[string]any{"merchant_region": "北京市朝阳区"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_BeginsWith(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "merchant_region", "operator": "begins_with", "value": "北京"},
		},
	}
	ctx := makeContext(map[string]any{"merchant_region": "北京市朝阳区"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_EndsWith(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "merchant_region", "operator": "ends_with", "value": "朝阳区"},
		},
	}
	ctx := makeContext(map[string]any{"merchant_region": "北京市朝阳区"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_ANDLogic(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "equal", "value": "USD"},
			map[string]any{"id": "order_amount", "operator": "greater", "value": "50"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD", "order_amount": "100"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}

	ctx2 := makeContext(map[string]any{"currency": "EUR", "order_amount": "100"}, defaultAccounts())
	result2 := engine.Evaluate(rule, ctx2)
	if len(result2) != 0 {
		t.Fatalf("expected 0 matched accounts, got %d", len(result2))
	}
}

func TestEvaluate_ORLogic(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "OR",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "equal", "value": "USD"},
			map[string]any{"id": "currency", "operator": "equal", "value": "EUR"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "EUR"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}

	ctx2 := makeContext(map[string]any{"currency": "GBP"}, defaultAccounts())
	result2 := engine.Evaluate(rule, ctx2)
	if len(result2) != 0 {
		t.Fatalf("expected 0 matched accounts, got %d", len(result2))
	}
}

func TestEvaluate_NestedConditions(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{
				"condition": "OR",
				"rules": []any{
					map[string]any{"id": "currency", "operator": "equal", "value": "USD"},
					map[string]any{"id": "currency", "operator": "equal", "value": "EUR"},
				},
			},
			map[string]any{"id": "order_amount", "operator": "greater", "value": "100"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "EUR", "order_amount": "200"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}

	ctx2 := makeContext(map[string]any{"currency": "EUR", "order_amount": "50"}, defaultAccounts())
	result2 := engine.Evaluate(rule, ctx2)
	if len(result2) != 0 {
		t.Fatalf("expected 0 matched accounts, got %d", len(result2))
	}
}

func TestEvaluate_EqualNumeric(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "order_amount", "operator": "equal", "value": "100"},
		},
	}
	ctx := makeContext(map[string]any{"order_amount": "100"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_InWithNumericValues(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "order_amount", "operator": "in", "value": []any{"100", "200", "300"}},
		},
	}
	ctx := makeContext(map[string]any{"order_amount": "200"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_MissingContextField(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "equal", "value": "USD"},
		},
	}
	ctx := makeContext(map[string]any{}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 0 {
		t.Fatalf("expected 0 matched accounts when context field is missing, got %d", len(result))
	}
}

func TestEvaluate_DefaultConditionIsAND(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "equal", "value": "USD"},
			map[string]any{"id": "order_amount", "operator": "greater", "value": "50"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD", "order_amount": "100"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 3 {
		t.Fatalf("expected 3 matched accounts with AND as default, got %d", len(result))
	}
}

// ─── 账户统计字段测试 ───────────────────────────────────────────

func TestEvaluate_AccountDayAmountFilter(t *testing.T) {
	engine := newTestEngine()
	// 规则：账户日金额 < 1000 的账号才匹配
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "account.[].day_amount", "operator": "less", "value": "1000"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	// ACC003: DayAmount=500 < 1000 ✓, ACC001: 5000 ✗, ACC002: 15000.50 ✗
	if len(result) != 1 || result[0].AccountNo != "ACC003" {
		t.Fatalf("expected [ACC003], got %v", result)
	}
}

func TestEvaluate_AccountDayOrderFilter(t *testing.T) {
	engine := newTestEngine()
	// 规则：账户日交易数 >= 10 的账号才匹配
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "account.[].day_order", "operator": "greater_or_equal", "value": "10"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	// ACC001: DayOrder=10 ✓, ACC002: DayOrder=50 ✓, ACC003: DayOrder=2 ✗
	if len(result) != 2 {
		t.Fatalf("expected 2 matched accounts, got %d", len(result))
	}
	accountNos := make(map[string]bool)
	for _, acc := range result {
		accountNos[acc.AccountNo] = true
	}
	if !accountNos["ACC001"] || !accountNos["ACC002"] {
		t.Fatalf("expected ACC001 and ACC002, got %v", result)
	}
}

func TestEvaluate_CombinedBusinessAndAccountFilter(t *testing.T) {
	engine := newTestEngine()
	// 规则：币种=USD 且 账户日金额 < 10000
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "equal", "value": "USD"},
			map[string]any{"id": "account.[].day_amount", "operator": "less", "value": "10000"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	// ACC001: 5000 < 10000 ✓, ACC003: 500 < 10000 ✓, ACC002: 15000.50 ✗
	if len(result) != 2 {
		t.Fatalf("expected 2 matched accounts, got %d", len(result))
	}
	accountNos := make(map[string]bool)
	for _, acc := range result {
		accountNos[acc.AccountNo] = true
	}
	if !accountNos["ACC001"] || !accountNos["ACC003"] {
		t.Fatalf("expected ACC001 and ACC003, got %v", result)
	}
}

func TestEvaluate_AccountFilterWithOR(t *testing.T) {
	engine := newTestEngine()
	// 规则：账户日金额 < 1000 或 日交易数 >= 50
	rule := map[string]any{
		"condition": "OR",
		"rules": []any{
			map[string]any{"id": "account.[].day_amount", "operator": "less", "value": "1000"},
			map[string]any{"id": "account.[].day_order", "operator": "greater_or_equal", "value": "50"},
		},
	}
	ctx := makeContext(map[string]any{}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	// ACC003: 500 < 1000 ✓, ACC002: DayOrder=50 ✓, ACC001: neither ✗
	if len(result) != 2 {
		t.Fatalf("expected 2 matched accounts, got %d", len(result))
	}
	accountNos := make(map[string]bool)
	for _, acc := range result {
		accountNos[acc.AccountNo] = true
	}
	if !accountNos["ACC002"] || !accountNos["ACC003"] {
		t.Fatalf("expected ACC002 and ACC003, got %v", result)
	}
}

func TestEvaluate_NoAccounts(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "equal", "value": "USD"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, nil)
	result := engine.Evaluate(rule, ctx)
	if result != nil {
		t.Fatalf("expected nil for no accounts, got %v", result)
	}
}

func TestEvaluate_AllAccountsFiltered(t *testing.T) {
	engine := newTestEngine()
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "currency", "operator": "equal", "value": "JPY"},
		},
	}
	ctx := makeContext(map[string]any{"currency": "USD"}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 0 {
		t.Fatalf("expected 0 matched accounts, got %d", len(result))
	}
}

func TestEvaluate_AccountNoFilter(t *testing.T) {
	engine := newTestEngine()
	// 规则：指定账号 = ACC002
	rule := map[string]any{
		"condition": "AND",
		"rules": []any{
			map[string]any{"id": "account_no", "operator": "equal", "value": "ACC002"},
		},
	}
	ctx := makeContext(map[string]any{}, defaultAccounts())
	result := engine.Evaluate(rule, ctx)
	if len(result) != 1 || result[0].AccountNo != "ACC002" {
		t.Fatalf("expected [ACC002], got %v", result)
	}
}

// ─── resolveFieldId 测试 ───────────────────────────────────────────

func TestResolveFieldId(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"currency", "currency"},
		{"order_amount", "order_amount"},
		{"account.[].day_order", "account_day_order"},
		{"account.[].day_amount", "account_day_amount"},
		{"merchant_region", "merchant_region"},
	}
	for _, tt := range tests {
		got := resolveFieldId(tt.input)
		if got != tt.expected {
			t.Errorf("resolveFieldId(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// ─── buildEvalContext 测试 ───────────────────────────────────────────

func TestBuildEvalContext(t *testing.T) {
	engine := newTestEngine()
	data := map[string]any{"currency": "USD", "order_amount": "100"}
	account := dto.RuleAccountContext{AccountNo: "ACC001", DayAmount: 5000.0, DayOrder: 10}

	ctx := engine.buildEvalContext(data, account)

	if ctx["currency"] != "USD" {
		t.Fatal("expected currency=USD in context")
	}
	if ctx["order_amount"] != "100" {
		t.Fatal("expected order_amount=100 in context")
	}
	if ctx["account_no"] != "ACC001" {
		t.Fatal("expected account_no=ACC001 in context")
	}
	if ctx["account_day_order"] != int64(10) {
		t.Fatalf("expected account_day_order=10, got %v", ctx["account_day_order"])
	}
	if ctx["account_day_amount"] != 5000.0 {
		t.Fatalf("expected account_day_amount=5000.0, got %v", ctx["account_day_amount"])
	}
}

func TestBuildEvalContextDoesNotMutateOriginal(t *testing.T) {
	engine := newTestEngine()
	data := map[string]any{"currency": "USD"}
	account := dto.RuleAccountContext{AccountNo: "ACC001", DayAmount: 100, DayOrder: 5}

	ctx := engine.buildEvalContext(data, account)
	ctx["currency"] = "EUR"

	if data["currency"] != "USD" {
		t.Fatal("buildEvalContext should not mutate original data map")
	}
}
