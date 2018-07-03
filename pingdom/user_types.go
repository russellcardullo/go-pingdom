package pingdom

type UserSms struct {
	Severity string `json:"severity"`
	CountryCode int64 `json:"country_code"`
	Number string `json:"number"`
	Provider string `json:"provider"`
}

type UserEmail struct {
	Severity string `json:"severity"`
	Address string `json:"address"`
}

// MaintenanceWindow represents a Pingdom Maintenance Window.
type User struct {
	Id    		   int64  `json:"id"`
	Paused         int64  `json:"paused,omitempty"`
	Username       string `json:"name,omitempty"`
	Sms			   []UserSmsResponse `json:"sms,omitempty"`
	Email 		   []UserEmailResponse `json:"email,omitempty"`
}


//func (u *User) PutParams() map[string]string {
//
//}
//
//func (u *User) PutContactParams() map[string]string {
//
//}
//
//func (u *User) PostParams() map[string]string {
//
//}
//
//func (u *User) PostContactParams() map[string]string {
//
//}
//
//func (u *User) DeleteParams() map[string]string {
//
//}
//
//func (u *User) DeleteContactParams() map[string]string {
//
//}
//
//func (u *User) Valid() error {
//
//}