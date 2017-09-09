package parameters

// RESTSettings defines the settings for the HTTP rest server.
type RESTSettings struct {
	Port     uint   `json:"port"`
	CORS     bool   `json:"cors"`
	Frontend string `json:"frontend"`
}
