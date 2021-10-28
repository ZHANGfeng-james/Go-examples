package builtin

type Size uint8

const (
	small Size = iota
	medium
	large
	extraLarge
)

func (s Size) String() string {
	switch s {
	case small:
		return "Small"
	case medium:
		return "Medium"
	case large:
		return "Large"
	case extraLarge:
		return "ExtraLarge"
	default:
		return ""
	}
}

const (
	a = 100 + iota
	b
	c
)
