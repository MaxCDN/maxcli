package maxcdn

import (
	"encoding/json"
	"net/http"
)

// Error represent a maxcnd error.
type Error struct {
	Message string `json:"message,omitempty"`
	Type    string `json:"type,omitempty"`
}

// Error implements go's error interface.
func (e Error) Error() string {
	return e.Type + ": " + e.Message
}

// Response object for all json requests.
type Response struct {
	Code  int             `json:"code,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error Error           `json:"error,omitempty"`

	// Non-JSON data.
	Headers http.Header `json:"-"`
}

// Generic is the generic data type for JSON responses from API calls.
type Generic map[string]interface{}

// LogRecord holds the data of a single record.
type LogRecord struct {
	Bytes           int     `json:"bytes"`
	CacheStatus     string  `json:"cache_status"`
	ClientAsn       string  `json:"client_asn"`
	ClientCity      string  `json:"client_city"`
	ClientContinent string  `json:"client_continent"`
	ClientCountry   string  `json:"client_country"`
	ClientDma       string  `json:"client_dma"`
	ClientIP        string  `json:"client_ip"`
	ClientLatitude  float64 `json:"client_latitude"`
	ClientLongitude float64 `json:"client_longitude"`
	ClientState     string  `json:"client_state"`
	CompanyID       int     `json:"company_id"`
	Hostname        string  `json:"hostname"`
	Method          string  `json:"method"`
	OriginTime      float64 `json:"origin_time"`
	Pop             string  `json:"pop"`
	Protocol        string  `json:"protocol"`
	QueryString     string  `json:"query_string"`
	Referer         string  `json:"referer"`
	Scheme          string  `json:"scheme"`
	Status          int     `json:"status"`
	Time            string  `json:"time"`
	URI             string  `json:"uri"`
	UserAgent       string  `json:"user_agent"`
	ZoneID          int     `json:"zone_id"`
}

// Logs .
type Logs struct {
	Limit       int         `json:"limit"`
	NextPageKey string      `json:"next_page_key"`
	Page        int         `json:"page"`
	Records     []LogRecord `json:"records"`
	RequestTime int         `json:"request_time"`
}
