package goprismic

import "time"

type Ref struct {
	Id          string `json:"id"`
	Ref         string `json:"ref"`
	Label       string `json:"label"`
	IsMasterRef bool   `json:"isMasterRef"`
	ScheduledAt int64  `json:"scheduledAt"`
}

func (r *Ref) ScheduledTime() *time.Time {
	if r.ScheduledAt == 0 {
		return nil
	}
	sec := r.ScheduledAt / 1000
	nsec := (r.ScheduledAt % 1000) * 1000
	date := time.Unix(sec, nsec)
	return &date
}
