package api_code

import "math"

type PageData struct {
	Page      uint64 `json:"page"`
	PageSize  uint64 `json:"page_size"`
	Total     int64  `json:"total"`
	TotalPage int64  `json:"total_page"`
}

func (p *PageData) GetTotalPadge() {
	if p.PageSize == 0 {
		return
	}
	p.TotalPage = int64(math.Ceil(float64(p.Total) / float64(p.PageSize)))

}
