package mock

import (
	"errors"

	"github.com/JohanVong/online_bazaar/pkg/models"
)

type CountryModel struct{}

var countryList = []*models.CountryOutput{
	{
		UUID: "uuid.v6[1]",
		Name: "TestCountry",
	},
	{
		UUID: "uuid.v6[2]",
		Name: "TestCountry2",
	},
	{
		UUID: "uuid.v6[3]",
		Name: "TestCountry3",
	},
}

func (c *CountryModel) GetList() ([]*models.CountryOutput, error) {
	return countryList, nil
}

func (c *CountryModel) GetByName(name string) (string, error) {
	for _, v := range countryList {
		if name == v.Name {
			return v.UUID, nil
		}
	}

	return "", errors.New("Provided country does not exist")
}
