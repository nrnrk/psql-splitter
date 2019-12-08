package head

import (
	log "github.com/sirupsen/logrus"
)

// These 3 characters are all you need to take care of when splitting sqls not parsing them.
// ref. https://www.postgresql.jp/document/9.2/html/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS
const semicolumn = byte(';')
const strConstEncloser = byte('\'')
const identifierEncloser = byte('"')
const lineCommentStarterFirst = byte('-')
const lineCommentStarterSecond = lineCommentStarterFirst
const lineCommentFinisher = byte('\n')
const blockCommentStarterFirst = byte('/')
const blockCommentStarterSecond = byte('*')
const blockCommentFinisherFirst = blockCommentStarterSecond
const blockCommentFinisherSecond = blockCommentStarterFirst

type Head interface {
	Continue(b byte)
	IsEndStmt() bool
	Restart()
}

type head struct {
	state         *state
	depth         int
	stateTransMap map[state]*transMap
}

func NewHead() Head {
	s := normal
	return &head{
		state:         &s,
		depth:         0,
		stateTransMap: stateTransMap,
	}
}

func (h *head) Continue(b byte) {
	tMap, ok := h.stateTransMap[*h.state]
	if !ok {
		log.WithFields(log.Fields{
			"state": h.state,
		}).Error("Unsupported state")
		panic(`unsupported state`)
	}

	nextState := tMap.get(b)
	if h.depth == 0 {
		h.state = &nextState
		return
	}

	// For next block comment (/* /* comment */ */)
	if *h.state == preBlockCommentStart {
		if nextState == inBlockComment {
			h.depth++
			h.state = &nextState
			return
		}
		if nextState == normal {
			nState := inBlockComment
			h.state = &nState
			return
		}
	}

	if *h.state == preBlockCommentEnd {
		if nextState == normal {
			h.depth--
			if h.depth == 0 {
				h.state = &nextState
				return
			}
			nState := inBlockComment
			h.state = &nState
			return
		}
	}
}

func (h *head) IsEndStmt() bool {
	return *h.state == statementEnd
}

func (h *head) Restart() {
	initState := normal
	h.state = &initState
}

type transMap struct {
	tMap map[byte]state
	d    state
}

func (t *transMap) get(b byte) state {
	v, ok := t.tMap[b]
	if !ok {
		return t.d
	}
	return v
}

var (
	normalMap transMap = transMap{
		tMap: map[byte]state{
			strConstEncloser:         inStrConst,
			identifierEncloser:       inIdentifier,
			lineCommentStarterFirst:  preLineCommentStart,
			blockCommentStarterFirst: preBlockCommentStart,
			semicolumn:               statementEnd,
		},
		d: normal,
	}

	inStrConstMap = transMap{
		tMap: map[byte]state{
			strConstEncloser: normal,
		},
		d: inStrConst,
	}

	inIdentifierMap = transMap{
		tMap: map[byte]state{
			identifierEncloser: normal,
		},
		d: inIdentifier,
	}

	preLineCommentStartMap = transMap{
		tMap: map[byte]state{
			lineCommentStarterSecond: inLineComment,
		},
		d: normal,
	}

	inLineCommentMap = transMap{
		tMap: map[byte]state{
			lineCommentFinisher: normal,
		},
		d: inLineComment,
	}

	preBlockCommentStartMap = transMap{
		tMap: map[byte]state{
			blockCommentStarterSecond: inBlockComment,
		},
		d: normal,
	}

	inBlockCommentMap = transMap{
		tMap: map[byte]state{
			blockCommentFinisherFirst: preBlockCommentEnd,
			blockCommentStarterFirst:  preBlockCommentStart,
		},
		d: inBlockComment,
	}

	preBlockCommentEndMap = transMap{
		tMap: map[byte]state{
			blockCommentFinisherSecond: normal,
		},
		d: inBlockComment,
	}

	stateTransMap map[state]*transMap = map[state]*transMap{
		normal:               &normalMap,
		inStrConst:           &inStrConstMap,
		inIdentifier:         &inIdentifierMap,
		preLineCommentStart:  &preLineCommentStartMap,
		inLineComment:        &inLineCommentMap,
		preBlockCommentStart: &preBlockCommentStartMap,
		inBlockComment:       &inBlockCommentMap,
		preBlockCommentEnd:   &preBlockCommentEndMap,
	}
)
