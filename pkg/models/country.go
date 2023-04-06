package models

// CountryOutput - структура на выход, в которую апи кладет список стран
type CountryOutput struct {
	UUID string `json:"-"`
	Name string
}
