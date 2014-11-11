package goprismic

type Ref struct {
	Id          string `json:"id"`
	Ref         string `json:"ref"`
	Label       string `json:"label"`
	IsMasterRef bool   `json:"isMasterRef"`
}
