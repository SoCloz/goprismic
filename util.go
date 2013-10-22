package goprismic

func stripQuery(query string) string {
	if len(query) < 2 {
		return query
	}
	return query[1 : len(query)-1]
}
