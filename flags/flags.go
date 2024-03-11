package flags

type Flag uint16

func Has[T ~uint16](flag, other T) bool {
	return (flag & other) != 0
}

func Set[T ~uint16](flag *T, other T) {
	*flag |= other
}

func Clear[T ~uint16](flag *T, other T) {
	*flag &= ^other
}

func Toggle[T ~uint16](flag *T, other T) {
	*flag ^= other
}
