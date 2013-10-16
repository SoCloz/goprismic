package goprismic

type Ref struct {
	Ref         string `json:"ref"`
	Label       string `json:"label"`
	IsMasterRef bool   `json:"isMasterRef"`
}
