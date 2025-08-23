package models

type Payment struct {
	Id                int     `xorm:"'id' pk autoincr"  json:"id,omitempty"`
	PaymentNumber     string  `xorm:"not null 'number'" json:"paymentNumber,omitempty"`
	PaymentDate       string  `xorm:"not null 'date'" json:"paymentDate,omitempty"`
	Balance           float64 `xorm:"not null 'balance'" json:"balance,omitempty"`
	Amount            float64 `xorm:"not null 'amount'" json:"amount,omitempty"`
	BaseAmount        float64 `xorm:"not null 'base_amount'" json:"-"`
	UsedAmount        float64 `xorm:"not null 'used_amount'" json:"usedAmount,omitempty"`
	Type              string  `xorm:"not null 'type'" json:"-"`
	PaymentMode       string  `xorm:"not null 'fop'" json:"paymentMode,omitempty"`
	Status            string  `xorm:"not null 'status'" json:"status,omitempty"`
	IdChartOfAccounts int     `xorm:"not null 'id_chart_of_accounts'" json:"-"`
	IdCurrency        int     `xorm:"not null 'id_currency'" json:"-"`
	IdCustomer        int     `xorm:"'id_customer'" json:"idCustomer,omitempty"`
	Tag               string  `xorm:"not null 'tag'" json:"-"`
}

func (*Payment) TableName() string {
	return "payment_received"
}
