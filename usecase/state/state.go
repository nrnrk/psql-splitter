package state

type State string

const (
	Normal               State = `normal`
	InStrConst           State = `inStrConst`
	InIdentifier         State = `inIdentifier`
	PreLineCommentStart  State = `preLineCommentStart`
	InLineComment        State = `inLineComment`
	PreBlockCommentStart State = `preBlockCommentStart`
	InBlockComment       State = `inBlockComment`
	PreBlockCommentEnd   State = `preBlockCommentEnd`
	StatementEnd         State = `statementEnd`
)
