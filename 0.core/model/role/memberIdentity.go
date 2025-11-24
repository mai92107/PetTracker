package role

type MemberIdentity int

const (
	MEMBER MemberIdentity = iota
	ADMIN
)

func (mi MemberIdentity) ToString() string {
	switch mi {
	case MEMBER:
		return "MEMBER"
	case ADMIN:
		return "ADMIN"
	default:
		return "UNKNOWN"
	}
}