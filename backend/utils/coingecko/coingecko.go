package coingecko

import (
	"fmt"
	"github.com/bytedance/sonic"
	"net/http"
	"strings"
)

const BaseURL = "https://api.coingecko.com/api/v3"

func GetPrice(coins []string, currencies ...string) (map[string]map[string]float64, error) {
	res, err := http.Get(fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", BaseURL, strings.Join(coins, ","), strings.Join(currencies, ",")))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from coingecko: %s", res.Status)
	}

	var data map[string]map[string]float64
	if err := sonic.ConfigStd.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
