package turn

type Turn int

const (
	Black Turn = iota
	White
)

func (t Turn) Flip() Turn { return t ^ 1 }

func (t Turn) String() string {
	switch t {
	case Black:
		return "先手"
	case White:
		return "後手"
	}
	return ""
}

func (t Turn) Symbol() string {
	switch t {
	case Black:
		return "▲"
	case White:
		return "△"
	}
	return ""
}
