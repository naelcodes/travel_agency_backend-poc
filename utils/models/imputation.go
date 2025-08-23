package models

type Imputation struct {
	Id                int     `xorm:"'id' pk autoincr" json:"id,omitempty"`
	AmountApplied     float64 `xorm:"amount_apply" json:"amountApplied,omitempty"`
	InvoiceAmount     float64 `xorm:"invoice_amount" json:"invoiceAmount,omitempty"`
	PaymentAmount     float64 `xorm:"payment_amount" json:"paymentAmount,omitempty"`
	IdInvoice         int     `xorm:"'id_invoice'" json:"-"`
	IdPaymentReceived int     `xorm:"'id_payment_received'" json:"-"`
	Tag               string  `xorm:" not null 'tag' " json:"-"`
}

func (*Imputation) TableName() string {
	return "invoice_payment_received"
}
