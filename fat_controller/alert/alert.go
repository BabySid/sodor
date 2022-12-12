package alert

type Alert interface {
	GetName() string
	GiveAlarm(content string) error
}
