package config

type Rigla struct {
	Token string
	URL   string
}

func (r Rigla) GetUrlData(MOKE_MODE bool) string {
	if MOKE_MODE {
		return r.URL
	} else {
		return r.URL
	}
}
