package firecache

type NoData struct{}

func (e *NoData) Error() string {
	return "no data"
}
