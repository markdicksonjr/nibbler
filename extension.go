package nibbler

// NoOpExtension implements nibbler.Extension, but all methods do nothing
// This is provided as a base class that can be used to more tersely define extensions
type NoOpExtension struct {
	Logger Logger
}

func (s *NoOpExtension) Init(app *Application) error {
	return nil
}

func (s *NoOpExtension) Destroy(app *Application) error {
	return nil
}

func (s *NoOpExtension) PostInit(app *Application) error {
	return nil
}

func (s *NoOpExtension) GetName() string {
	return "nameless"
}
