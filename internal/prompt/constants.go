package prompt

const (
	lineDelimiter = '\n'

	readInputFormat = "read input: %w"
)

var yesAnswers = map[string]bool{
	"y":   true,
	"yes": true,
}
