package models

type Invoice struct {
	Id            int     `xorm:"'id' pk autoincr" json:"id,omitempty"`
	CreationDate  string  `xorm:"not null 'creation_date'" json:"creationDate,omitempty"`
	InvoiceNumber string  `xorm:"not null 'invoice_number'" json:"invoiceNumber,omitempty"`
	DueDate       string  `xorm:"not null 'due_date'" json:"dueDate,omitempty"`
	Status        string  `xorm:"not null 'status'" json:"status,omitempty"`
	Amount        float64 `xorm:"not null 'amount'" json:"amount,omitempty"`
	Balance       float64 `xorm:"not null 'balance'" json:"balance,omitempty"`
	NetAmount     float64 `xorm:"not null 'net_amount'" json:"-"`
	BaseAmount    float64 `xorm:"not null 'base_amount'" json:"-"`
	CreditApply   float64 `xorm:"not null 'credit_apply'" json:"creditApply,omitempty"`
	Tag           string  `xorm:" not null 'tag' " json:"-"`

	IdCustomer int `xorm:"'id_customer'" json:"idCustomer,omitempty"`
}

func (*Invoice) TableName() string {
	return "invoice"
}
