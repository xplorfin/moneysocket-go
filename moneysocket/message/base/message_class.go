package base

type MessageClass int

const (
	Notification MessageClass = 0
	Request      MessageClass = iota
)

const (
	NotificationName = "NOTIFICATION"
	RequestName      = "REQUEST"
)

func (m MessageClass) ToString() string {
	switch m {
	case Notification:
		return NotificationName
	case Request:
		return RequestName
	}
	panic("message not found")
}

func MessageClassFromString(class string) MessageClass {
	switch class {
	case NotificationName:
		return Notification
	case RequestName:
		return Request
	}
	panic("message not recognized")
}
