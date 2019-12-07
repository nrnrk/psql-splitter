package head

type state string

const (
	normal               state = `normal`
	inStrConst           state = `inStrConst`
	inIdentifier         state = `inIdentifier`
	preLineCommentStart  state = `preLineCommentStart`
	inLineComment        state = `inLineComment`
	preBlockCommentStart state = `preBlockCommentStart`
	inBlockComment       state = `inBlockComment`
	preBlockCommentEnd   state = `preBlockCommentEnd`
	statementEnd         state = `statementEnd`
)
