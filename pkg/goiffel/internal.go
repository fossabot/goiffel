package goiffel

import (
	"time"
	"log"
        "encoding/json"
	"github.com/google/uuid"
)

func getTimeStampMilliSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func newEiffelEvent(typeName string, data interface{}, links []EiffelLink) EiffelEvent {
	return EiffelEvent{
		Meta: EiffelMeta {
			Id: uuid.New().String(),
			Type: typeName,
			Version: EventVersion,
			Time: getTimeStampMilliSeconds(),
		},
		Data: data,
		Links: links,
	}
}

func postParseEventData(parsed_data interface{}, evt *EiffelEvent) error {
	// We convert the map[string]interface{}
	// to json encoded data - then we unmarshal it
	// again
	temporary_json, err := json.Marshal(evt.Data)
	if err != nil { return err }

	err = json.Unmarshal(temporary_json, parsed_data)
	if err != nil { return err }

	return nil
}

func postReceiveParser(evt *EiffelEvent) error {
	switch t := evt.Meta.Type; t {

	case EiffelArtifactCreatedEvent:
		log.Printf("--- parsing EiffelArtifactCreatedEventData...")
		var artifactCreated EiffelArtifactCreatedEventData
		err := postParseEventData(&artifactCreated, evt)
		if err != nil {
			evt.Data = artifactCreated
		}
		return err

	case EiffelArtifactPublishedEvent:
		log.Printf("--- parsing EiffelArtifactPublishedEventData...")
		var artifactPublished EiffelArtifactPublishedEventData
		err := postParseEventData(&artifactPublished, evt)
		if err != nil {
			evt.Data = artifactPublished
		}
		return err

	case EiffelCompositionDefinedEvent:
		log.Printf("--- parsing EiffelCompositionDefinedEventData...")
		var compositionDefined EiffelCompositionDefinedEventData
		err := postParseEventData(&compositionDefined, evt)
		if err != nil {
			evt.Data = compositionDefined
		}
		return err

	default:
		return nil
	}
}
