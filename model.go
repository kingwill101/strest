package strest

//Payload entire yaml configuration
type Payload struct {
	Version string             `yaml:"version"`
	Request map[string]Request `yaml:"requests"`
}

//RequestData Request data parameters
type RequestData struct {
	Params  map[string]interface{} `yaml:"params"`
	Headers map[string]string      `yaml:"headers"`
	Form    map[string]string      `yaml:"form"`
	JSON    map[string]interface{} `yaml:"json"`
	Raw     string                 `yaml:"raw"`
}

//Request Information about the current request
type Request struct {
	FailOnError bool        `yaml:"failOnError"`
	URL         string      `yaml:"url"`
	Timeout     int         `yaml:"timeout"`
	Method      string      `yaml:"method"`
	Data        RequestData `yaml:"data"`
	Log         bool        `yaml:"log"`
	Validation  Validation  `yaml:"validation"`
	Repeat      int         `yaml:"repeat"`
	Delay       int         `yaml:"delay"`
}

//Validation types of validation available on each request
type Validation struct {
	Body       string `yaml:"body"`
	Status     string `yaml:"status"`
	StatusCode int    `yaml:"statusCode"`
}
