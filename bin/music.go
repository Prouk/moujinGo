package bin

type Music struct {
	Url      string
	Title    string
	Channel  string
	Duration int
	AddedBy  struct {
		Name string
		Icon string
	}
}

func (m *Music) SetInfosFromWeb() {
	
}
