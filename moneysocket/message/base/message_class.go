package base

// MessageClass is a type used to determine a message class
type MessageClass int

const (
	// Notification is the MessageClass of the Notification
	Notification MessageClass = 0
	// Request is the MessageClass of the Request
	Request MessageClass = iota
)

const (
	// NotificationName is the name of a notification class
	NotificationName = "NOTIFICATION"
	// RequestName is the name of a request class
	RequestName = "REQUEST"
)

// ToString converts a message class to the string (either NotificationName or RequestName)
func (m MessageClass) ToString() string {
	switch m {
	case Notification:
		return NotificationName
	case Request:
		return RequestName
	}
	panic("message not found")
}

// MessageClassFromString determines the MessageClass from a given string
func MessageClassFromString(class string) MessageClass {
	switch class {
	case NotificationName:
		return Notification
	case RequestName:
		return Request
	}
	panic("message not recognized")
}
