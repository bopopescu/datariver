package canal

type Holder interface {
	PushBack(data *EventData) error
}
