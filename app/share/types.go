package share

type XmlEntry struct {
	Uid       string `xml:"uid"`
	FirstName string `xml:"firstName"`
	LastName  string `xml:"lastName"`
	SdnType   string `xml:"sdnType"`
}

type JsonEntry struct {
	Uid       int    `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
