package forms

type errors map[string][]string

func (e errors) AddError(tagId, message string) {
	e[tagId] = append(e[tagId], message)
}
func (e errors) GetError(tagId string) string {
	es := e[tagId]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
