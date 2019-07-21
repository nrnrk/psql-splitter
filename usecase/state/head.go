package state

import (
	"fmt"
)

// These 3 characters are all you need to take care of when splitting sqls not parsing them.
// ref. https://www.postgresql.jp/document/9.2/html/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS
const Semicolumn = byte(';')
const StrConstEncloser = byte('\'')
const IdentifierEncloser = byte('"')
const LineCommentStarterFirst = byte('-')
const LineCommentStarterSecond = LineCommentStarterFirst
const LineCommentFinisher = byte('\n')
const BlockCommentStarterFirst = byte('/')
const BlockCommentStarterSecond = byte('*')
const BlockCommentFinisherFirst = BlockCommentStarterSecond
const BlockCommentFinisherSecond = BlockCommentStarterFirst

type Head interface {
	Continue(b byte)
	IsEndStmt() bool
	Restart()
}

type head struct {
	state         *State
	depth         int
	stateTransMap map[State]*transMap
}

func NewHead() Head {
	s := Normal
	return &head{
		state:         &s,
		depth:         0,
		stateTransMap: stateTransMap,
	}
}

func (h *head) Continue(b byte) {
	tMap, ok := h.stateTransMap[*h.state]
	if !ok {
		fmt.Printf("%v", *h.state)
		panic(`unsupported state`)
	}

	nextState := tMap.get(b)
	if h.depth == 0 {
		h.state = &nextState
		return
	}

	// For next block comment (/* /* comment */ */)
	if *h.state == PreBlockCommentStart {
		if nextState == InBlockComment {
			h.depth++
			h.state = &nextState
			return
		}
		if nextState == Normal {
			nState := InBlockComment
			h.state = &nState
			return
		}
	}

	if *h.state == PreBlockCommentEnd {
		if nextState == Normal {
			h.depth--
			if h.depth == 0 {
				h.state = &nextState
				return
			}
			nState := InBlockComment
			h.state = &nState
			return
		}
	}
}

func (h *head) IsEndStmt() bool {
	return *h.state == StatementEnd
}

func (h *head) Restart() {
	initState := Normal
	h.state = &initState
}

type transMap struct {
	tMap map[byte]State
	d    State
}

func (t *transMap) get(b byte) State {
	v, ok := t.tMap[b]
	if !ok {
		return t.d
	}
	return v
}

var (
	normalMap transMap = transMap{
		tMap: map[byte]State{
			StrConstEncloser:         InStrConst,
			IdentifierEncloser:       InIdentifier,
			LineCommentStarterFirst:  PreLineCommentStart,
			BlockCommentStarterFirst: PreBlockCommentStart,
			Semicolumn:               StatementEnd,
		},
		d: Normal,
	}

	inStrConstMap = transMap{
		tMap: map[byte]State{
			StrConstEncloser: Normal,
		},
		d: InStrConst,
	}

	inIdentifierMap = transMap{
		tMap: map[byte]State{
			IdentifierEncloser: Normal,
		},
		d: InIdentifier,
	}

	preLineCommentStartMap = transMap{
		tMap: map[byte]State{
			LineCommentStarterSecond: InLineComment,
		},
		d: Normal,
	}

	inLineCommentMap = transMap{
		tMap: map[byte]State{
			LineCommentFinisher: Normal,
		},
		d: InLineComment,
	}

	preBlockCommentStartMap = transMap{
		tMap: map[byte]State{
			BlockCommentStarterSecond: InBlockComment,
		},
		d: Normal,
	}

	inBlockCommentMap = transMap{
		tMap: map[byte]State{
			BlockCommentFinisherFirst: PreBlockCommentEnd,
			BlockCommentStarterFirst:  PreBlockCommentStart,
		},
		d: InBlockComment,
	}

	preBlockCommentEndMap = transMap{
		tMap: map[byte]State{
			BlockCommentFinisherSecond: Normal,
		},
		d: InBlockComment,
	}

	stateTransMap map[State]*transMap = map[State]*transMap{
		Normal:               &normalMap,
		InStrConst:           &inStrConstMap,
		InIdentifier:         &inIdentifierMap,
		PreLineCommentStart:  &preLineCommentStartMap,
		InLineComment:        &inLineCommentMap,
		PreBlockCommentStart: &preBlockCommentStartMap,
		InBlockComment:       &inBlockCommentMap,
		PreBlockCommentEnd:   &preBlockCommentEndMap,
	}
)
