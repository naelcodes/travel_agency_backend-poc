// models/index.go
package models

type Customer struct {
	Id              int    `xorm:"'id' pk autoincr" json:"id,omitempty"`
	CustomerName    string `xorm:"not null 'customer_name'" json:"customerName,omitempty"`
	AccountNumber   string `xorm:"not null 'account_number'" json:"accountNumber,omitempty"`
	IdCurrency      int    `xorm:"'id_currency'" json:"-"`
	IdCountry       int    `xorm:"'id_country'" json:"-"`
	Alias           string `xorm:"not null unique 'alias'" json:"alias,omitempty"`
	AbKey           string `xorm:"not null unique 'ab_key'" json:"abKey,omitempty"`
	State           string `xorm:"not null 'state'" json:"state,omitempty"`
	TmcClientNumber string `xorm:"not null unique 'tmc_client_number'" json:"tmcClientNumber,omitempty"`
	Tag             string `xorm:" not null 'tag' " json:"-"`
}

func (*Customer) TableName() string {
	return "customer"
}
