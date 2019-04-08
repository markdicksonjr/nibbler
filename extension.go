package nibbler

type NoOpExtension struct {
}

func (s *NoOpExtension) Init(app *Application) error {
	return nil
}

func (s *NoOpExtension) Destroy(app *Application) error {
	return nil
}

func (s *NoOpExtension) AddRoutes(app *Application) error {
	return nil
}
