package httpclient

type HttpClientConfig struct {
	Timeout int `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

func (cl *HttpClientConfig) ApplyDefaults() {
	if cl.Timeout == 0 {
		cl.Timeout = 60
	}
}
