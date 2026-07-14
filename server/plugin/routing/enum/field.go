package enum

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInteger FieldType = "integer"
	FieldTypeDouble  FieldType = "double"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeDate    FieldType = "date"
)

func (f FieldType) Label() string {
	switch f {
	case FieldTypeString:
		return "字符串"
	case FieldTypeInteger:
		return "整数"
	case FieldTypeDouble:
		return "小数"
	case FieldTypeBoolean:
		return "布尔"
	case FieldTypeDate:
		return "日期"
	default:
		return ""
	}
}

func (f FieldType) IsNumeric() bool {
	return f == FieldTypeInteger || f == FieldTypeDouble
}
