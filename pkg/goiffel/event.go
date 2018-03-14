package goiffel

const EventVersion string = "1.1.0"

const EiffelArtifactCreatedEvent string = "EiffelArtifactCreatedEvent"
const EiffelArtifactPublishedEvent string = "EiffelArtifactPublishedEvent"
const DefaultEiffelEvent string = "DefaultEiffelEvent"

type CustomData map[string]string

type EiffelMeta struct {
        Id              string          `required:"true"`
        Type            string          `required:"true"`
        Version         string          `required:"true"`
        Time            int64           `required:"true"` // Creation, milliseconds since epoch
        Tags            []string        `required:"false"`
}

type EiffelLink struct {
        Type            string
        Target          string
}

type EiffelAnnouncementPublishedEventData struct {
	Heading         string            `required:"true"`
	Body            string            `required:"true"`
	Uri             string            `required:"true"`
	Severity        string            `required:"true"`
	CustomData      CustomData        `required:"false"`
}

type Gav struct {
	GroupId         string            `required:"true"`
	ArtifactId      string            `required:"true"`
	Version         string            `required:"true"`
}

type FileInformation struct {
	Classifier      string            `required:"true"`
	Extension       string            `required:"true"`
}

type EiffelArtifactCreatedEventData struct {
	Gav                     Gav                 `required:"true"`
	FileInformation         []FileInformation   `required:"false"`
	BuildCommand            string              `required:"false"`
	RequiresImplementation  string              `required:"false"`
	Implements              []Gav               `required:"false"`
	DependsOn               []Gav               `required:"false"`
	Name                    string              `required:"false"`
	CustomData              CustomData          `required:"false"`
}

type Location struct {
	Type            string          `required:"true"`
	Uri             string          `required:"true"`
}

type EiffelArtifactPublishedEventData struct {
	Locations       []Location      `required:"true"`
}

type EiffelEvent struct {
        Meta            EiffelMeta      `required:"true"`
        Data            interface{}     `required:"true"`
        Links           []EiffelLink
}

type OnEiffelEventReceived func(event EiffelEvent) ()

type EventCallbacks map[string]OnEiffelEventReceived

type EiffelChannel struct {
        ChannelData     interface{}     `required:"true"`
}

/***********************************
 *
 * Event Creation functions
 *
 ***********************************/

func InitiateEiffelArtifactCreatedEvent(
	data EiffelArtifactCreatedEventData,
	links []EiffelLink) EiffelEvent {
	return newEiffelEvent(EiffelArtifactCreatedEvent, data, links)
}

func InitiateEiffelArtifactPublishedEvent(
	data EiffelArtifactPublishedEventData,
	links []EiffelLink) EiffelEvent {
	return newEiffelEvent(EiffelArtifactPublishedEvent, data, links)
}
