package regex

type String struct {
	Value *rune
	Next  *String
}

func (s *String) String() string {
	if s == nil {
		return ""
	}
	return string(*s.Value) + s.Next.String()
}

func NewString(s string) *String {
	if len(s) == 0 {
		return nil
	}
	return &String{
		Value: &[]rune(s)[0],
		Next:  NewString(s[1:]),
	}
}
