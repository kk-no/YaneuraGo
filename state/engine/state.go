package engine

type State int

const (
	WaitConnecting State = iota + 1
	Connected
	WaitReadyOk
	WaitCommand
	WaitBestMove
	WaitOneLine
	Disconnected State = 99
)

func (s State) String() string {
	switch s {
	case WaitConnecting:
		return "WaitConnecting"
	case Connected:
		return "Connected"
	case WaitReadyOk:
		return "WaitReadyOk"
	case WaitCommand:
		return "WaitCommand"
	case WaitBestMove:
		return "WaitBestMove"
	case WaitOneLine:
		return "WaitOneLine"
	case Disconnected:
		return "Disconnected"
	}
	return ""
}
