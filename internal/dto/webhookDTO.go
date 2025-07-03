package dto

type IPChangeNotification struct {
	UserGUID string `json:"user_guid"`
	OldIP    string `json:"old_ip"`
	NewIP    string `json:"new_ip"`
	Time     string `json:"time"`
}
