package models

type TravelItem struct {
	Id              int     `json:"id,omitempty" xorm:"'id' pk autoincr"`
	TotalPrice      float64 `xorm:"not null 'total_price'" json:"totalPrice,omitempty"`
	Itinerary       string  `xorm:"not null 'itinerary'" json:"itinerary,omitempty"`
	TravelerName    string  `xorm:"not null 'traveler_name'" json:"travelerName,omitempty"`
	TicketNumber    int     `xorm:"not null 'ticket_number'" json:"ticketNumber,omitempty"`
	TransactionType string  `xorm:"not null 'transaction_type'" json:"-"`
	ProductType     string  `xorm:"not null 'product_type'" json:"-"`
	Status          string  `xorm:"not null 'status'" json:"-"`
	IdInvoice       int     `xorm:"'id_invoice'" json:"-"`
}

func (*TravelItem) TableName() string {
	return "air_booking"
}
