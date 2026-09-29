package uikit

func AskParsed[T any](s Session, value, label string, parse func(string) (T, error)) (T, error) {
	raw, err := s.ValueOrAsk(value, label, func(input string) error {
		_, err := parse(input)
		return Displayable(err)
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return parse(raw)
}
