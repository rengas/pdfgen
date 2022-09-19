package design

type Design struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ProfileId string `json:"profileId"`
	Fields    []byte `json:"fields"`
	Template  []byte `json:"design"`
}
