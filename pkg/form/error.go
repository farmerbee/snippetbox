package form

type errors map[string][]string

func (e errors) Add(key, msg string) {
	e[key] = append(e[key], msg)
}

func (e errors) Get(key string) string {
	msg := e[key]
	if len(msg) > 0 {
		return msg[0]
	}

	return ""
}
