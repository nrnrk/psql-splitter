package split

type SplittedStatements struct {
	Statements string
	Order      int
}

func (ss *SplittedStatements) AddNewLine() {
	ss.Statements += "\n"
}
