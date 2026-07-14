package enum

type InputType string

const (
	InputTypeText     InputType = "text"
	InputTypeSelect   InputType = "select"
	InputTypeRadio    InputType = "radio"
	InputTypeCheckbox InputType = "checkbox"
	InputTypeTextarea InputType = "textarea"
)

func (i InputType) Label() string {
	switch i {
	case InputTypeText:
		return "文本输入"
	case InputTypeSelect:
		return "下拉选择"
	case InputTypeRadio:
		return "单选"
	case InputTypeCheckbox:
		return "多选"
	case InputTypeTextarea:
		return "文本域"
	default:
		return ""
	}
}
