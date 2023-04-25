package models

type User struct {
	Type     int    `json:"type"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
}

type Usertype struct {
	UserType string `json:"usertype"`
	ID       int    `json:"id" gorm:"primary_key"`
}
type EmailType struct {
	Email    string `json:"email"`
	UserType string `json:"userType"`
}

// t.id,b.buildingName,c.type,p.type,r.roomNumber,s.severity
type Ticket struct {
	Id         int    `json:"id"`
	Building   string `json:"building"`
	Category   string `json:"category"`
	Progress   string `json:"progress"`
	RoomNum    string `json:"roomNum"`
	TimeSubmit string `json:"timeSubmit"`
	Issue      string `json:"issue"`
	User       string `json:"user"`
}

type TicketDetails struct {
	Id               int    `json:"id"`
	Building         string `json:"building"`
	Category         string `json:"category"`
	Progress         string `json:"progress"`
	RoomNum          string `json:"roomNum"`
	TimeSubmit       string `json:"timeSubmit"`
	Image            string `json:"image"`
	UserComments     string `json:"userComments"`
	InternalComments string `json:"internalComments"`
	TimeFinished     string `json:"timeFinished"`
	User             string `json:"user"`
	Severity         string `json:"severity"`
	Issue            string `json:"issue"`
	SeverityID       string `json:"severityID"`
	ProgressID       string `json:"progressID"`
}
type SubmitTicket struct {
	Building     int    `json:"building"`
	Category     int    `json:"category"`
	Progress     int    `json:"progress"`
	RoomNum      int    `json:"roomNum"`
	Image        string `json:"image"`
	UserComments string `json:"userComments"`
	User         string `json:"user"`
}
type UpdateTicket struct {
	Building int `json:"building"`
	Category int `json:"category"`
	Progress int `json:"progress"`
	RoomNum  int `json:"roomNum"`
}
