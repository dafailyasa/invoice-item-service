package pagination

type PaginationRequest struct {
	Limit   int    `json:"limit,omitempty" query:"limit"`
	Page    int    `json:"page,omitempty" query:"page"`
	Sort    string `json:"sort,omitempty" query:"sort"`
	Keyword string `json:"keyword,omitempty" query:"keyword"`
	Total   int64  `json:"total,omitempty" example:"100"`

	// invoices
	InvoiceID string `json:"invoiceId,omitempty" query:"invoiceId"`
	ItemCount string `json:"itemCount,omitempty" query:"itemCount"`
	Customer  string `json:"customer,omitempty" query:"customer"`
	IssueDate string `json:"issueDate,omitempty" query:"issueDate"`
	DueDate   string `json:"dueDate,omitempty" query:"dueDate"`
	Status    string `json:"status,omitempty" query:"status"`
}

func (pr PaginationRequest) Validate() (err error) {
	pr.Limit = pr.GetLimit()
	pr.Page = pr.GetPage()
	pr.Sort = pr.GetSort()
	return
}

func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PaginationRequest) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *PaginationRequest) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *PaginationRequest) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}
