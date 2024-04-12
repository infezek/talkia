package controller_bots

type TypeImage string

var (
	TypeImageAvatar     TypeImage = "avatar"
	TypeImageBackground TypeImage = "background"
)

func (t TypeImage) String() string {
	return string(t)
}
