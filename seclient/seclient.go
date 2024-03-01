package seclient

import (
	"github.com/idoberko2/semonitor/general"
	"github.com/pkg/errors"
	"time"

	"github.com/imroc/req/v3"
)

type SEClient interface {
	GetEnergy(start time.Time, end time.Time) ([]general.Energy, error)
}

type seclient struct {
	apiKey string
	siteId string
	client *req.Client
}

func NewSEClient(client *req.Client, apiKey string, siteId string) SEClient {
	return &seclient{
		apiKey: apiKey,
		siteId: siteId,
		client: client,
	}
}

func (c *seclient) GetEnergy(start time.Time, end time.Time) ([]general.Energy, error) {
	var energyRes EnergyResponseContainer
	_, err := c.client.R().SetQueryParams(map[string]string{
		"api_key":   c.apiKey,
		"startDate": formatDateForRequest(start),
		"endDate":   formatDateForRequest(end),
		"timeUnit":  "QUARTER_OF_AN_HOUR",
	}).SetSuccessResult(&energyRes).Get("https://monitoringapi.solaredge.com/site/" + c.siteId + "/energy")
	if err != nil {
		return nil, err
	}

	return parseEnergyResponse(energyRes)
}

func parseEnergyResponse(energyRes EnergyResponseContainer) ([]general.Energy, error) {
	var result []general.Energy

	for _, val := range energyRes.Values {
		energy, err := parseEnergyResponseValue(val)
		if err != nil {
			return nil, err
		}

		result = append(result, energy)
	}

	return result, nil
}

func parseEnergyResponseValue(entry EnergyResponseValue) (general.Energy, error) {
	dt, err := time.Parse("2006-01-02 15:04:05", entry.Date)
	if err != nil {
		return general.Energy{}, errors.Wrap(err, "failed to parse energy datetime")
	}

	return general.Energy{
		DateTime: dt,
		Value:    int(entry.Value),
	}, nil
}

func formatDateForRequest(dt time.Time) string {
	return dt.Format("2006-01-02")
}

type EnergyResponse struct {
	Values []EnergyResponseValue
}

type EnergyResponseValue struct {
	Date  string
	Value float64
}

type EnergyResponseContainer struct {
	EnergyResponse `json:"energy"`
}
