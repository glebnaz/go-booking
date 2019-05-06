package config

type RostaConfig struct {
	USER     string
	PASS     string
	URL      string
	URL_TEST string
}

func (r *RostaConfig) GetURLData(MockMode bool) (string, string, string) {
	if MockMode {
		return r.URL_TEST, r.USER, r.PASS
	} else {
		return r.URL, r.USER, r.PASS
	}
}
