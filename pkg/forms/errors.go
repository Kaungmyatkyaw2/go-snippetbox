package forms

type errors map[string][]string

func (e errors) Add(field, msg string) {
	e[field] = append(e[field], msg)
}

func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
