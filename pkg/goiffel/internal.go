package goiffel

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func getTimeStampMilliSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func newEiffelEvent(typeName string, data interface{}, links []EiffelLink) EiffelEvent {
	return EiffelEvent{
		Meta: EiffelMeta{
			Id:      uuid.New().String(),
			Type:    typeName,
			Version: EventVersion,
			Time:    getTimeStampMilliSeconds(),
		},
		Data:  data,
		Links: links,
	}
}

func postParseEventData(parsed_data interface{}, evt *EiffelEvent) error {
	// We convert the map[string]interface{}
	// to json encoded data - then we unmarshal it
	// again
	temporary_json, err := json.Marshal(evt.Data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(temporary_json, parsed_data)
	if err != nil {
		return err
	}

	return nil
}

func postReceiveParser(evt *EiffelEvent) error {
	switch t := evt.Meta.Type; t {
	case EiffelArtifactCreatedEvent:
		var parsed_data EiffelArtifactCreatedEventData
		err := postParseEventData(&parsed_data, evt)
		evt.Data = parsed_data
		return err

	case EiffelArtifactPublishedEvent:
		var parsed_data EiffelArtifactPublishedEventData
		err := postParseEventData(&parsed_data, evt)
		evt.Data = parsed_data
		return err

	case EiffelCompositionDefinedEvent:
		var parsed_data EiffelCompositionDefinedEventData
		err := postParseEventData(&parsed_data, evt)
		evt.Data = parsed_data
		return err

	default:
		return nil
	}
}
