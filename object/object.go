package object

type ObjectType string

type Object interface {
	Type() ObjectType
	Value() interface{}
	Inspect() string
}

func New(input interface{}) Object {
	switch input := convertToNative(input).(type) {
	case float64:
		return &Number{value: float64(input)}
	case string:
		return &String{value: input}
	case bool:
		return &Boolean{value: input}
	default:
		return NULL
	}
}

// Pattern from https://stackoverflow.com/a/40178331
// We want to treat all numbers as float64
func convertToNative(input interface{}) interface{} {
	switch input := input.(type) {
	case uint8:
		return float64(input)
	case int8:
		return float64(input)
	case uint16:
		return float64(input)
	case int16:
		return float64(input)
	case uint32:
		return float64(input)
	case int32:
		return float64(input)
	case uint64:
		return float64(input)
	case int64:
		return float64(input)
	case int:
		return float64(input)
	case float32:
		return float64(input)
	case float64:
		return float64(input)
	default:
		return input
	}
}
