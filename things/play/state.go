package play

type State uint8

const (
	StateIdle State = iota
	StateLoading
	StatePlaying
	StatePaused
	StateStopped
	StateEnded
	StateError
)

func (s State) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateLoading:
		return "loading"
	case StatePlaying:
		return "playing"
	case StatePaused:
		return "paused"
	case StateStopped:
		return "stopped"
	case StateEnded:
		return "ended"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}
