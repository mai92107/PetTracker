package role

type MemberIdentity int

const (
	GUEST MemberIdentity = iota
	MEMBER
	ADMIN
)

func (mi MemberIdentity) ToString() string {
	switch mi {
	case MEMBER:
		return "MEMBER"
	case ADMIN:
		return "ADMIN"
	case GUEST:
		return "GUEST"
	default:
		return "UNKNOWN"
	}
}
