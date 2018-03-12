package goiffel

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
	CustomData      map[string]string `required:"false"`
}

type EiffelEvent struct {
        Meta            EiffelMeta      `required:"true"`
        Data            interface{}     `required:"true"`
        Links           []EiffelLink
}

type OnEiffelEventReceived func(event EiffelEvent) ()

type EiffelChannel struct {
        ChannelData     interface{}     `required:"true"`
}
