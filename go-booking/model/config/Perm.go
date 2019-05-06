package config

type PermFarmConfig struct {
	Url     string
	UrlTest string
}

func (p *PermFarmConfig) GetUrlData(MockMode bool) string {
	if MockMode {
		return p.UrlTest
	} else {
		return p.Url
	}
}
