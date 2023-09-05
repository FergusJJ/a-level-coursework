package utils

import (
	"time"

	http "github.com/useflyent/fhttp"

	fhttp "github.com/AlienRecall/fhttp"

	sensorutils "github.com/FergusJJ/go-sensor/utils"
)

type ZalandoSession struct {
	XSRF      string
	UserAgent string // will be passed to api
	Profile   *Profile
}

//myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
type AsosSession struct {
	UserAgent string

	Cid              string
	Idxsrf           string
	StLoginToken     string
	LoginURL         string
	LoginAccessToken string

	BagId      string
	VariantID  int
	CustomerId string

	IsCheckedOut bool
	IsLoggedIn   bool

	Profile             *Profile
	AkamaiClient        *fhttp.Client
	AsosClient          *http.Client
	AbckDevice          *sensorutils.Device
	Bearer              *AsosBearerResponse
	SecuredTouchSession *SecuredTouchSession
	WebhookInfo         *WebhookInfo
}

type WebhookInfo struct {
	ItemName  string
	Pid       string
	Price     string
	Size      string
	ImageLink string
}

type SecuredTouchSession struct {
	StDeviceId  string
	StToken     string
	StSessionId string
	Index       int
	Username    string
}

type AsosBearerResponse struct {
	Id_token      string
	Access_token  string
	Token_type    string
	Expires_in    string
	Issued        string
	Expires       string
	Scope         string
	State         string
	Session_state string
}

type Profile struct {
	Store        string
	Mode         string
	ProductURL   string
	Size         string
	Delay        string
	Proxy        string
	FirstName    string
	LastName     string
	Email        string
	Password     string
	AddressLine1 string
	AddressLine2 string
	City         string
	Province     string
	Postcode     string
	CountryCode  string
	CC           PaymentDetails
	Phone        string
	TaskID       string
}

type PaymentDetails struct {
	CardNumber  string
	ExpiryMonth string
	ExpiryYear  string
	CVC         string
}

type CartValues struct {
	Bag struct {
		Total struct {
			TotalSalesTax interface{} `json:"totalSalesTax"`
			ItemsSubTotal struct {
				Xrp   float64 `json:"xrp"`
				Value float64 `json:"value"`
				Text  string  `json:"text"`
			} `json:"itemsSubTotal"`
			TotalDiscount struct {
				Xrp   float64 `json:"xrp"`
				Value float64 `json:"value"`
				Text  string  `json:"text"`
			} `json:"totalDiscount"`
			TotalDelivery struct {
				Xrp   float64 `json:"xrp"`
				Value float64 `json:"value"`
				Text  string  `json:"text"`
			} `json:"totalDelivery"`
			Total struct {
				Xrp   float64 `json:"xrp"`
				Value float64 `json:"value"`
				Text  string  `json:"text"`
			} `json:"total"`
		} `json:"total"`
		Delivery struct {
			Options []struct {
				Price struct {
					SalesTax interface{} `json:"salesTax"`
					Current  struct {
						Xrp   float64 `json:"xrp"`
						Value float64 `json:"value"`
						Text  string  `json:"text"`
					} `json:"current"`
					Discount   interface{} `json:"discount"`
					PriceToPay struct {
						Xrp   float64 `json:"xrp"`
						Value float64 `json:"value"`
						Text  string  `json:"text"`
					} `json:"priceToPay"`
					PremierPrice struct {
						Xrp   float64 `json:"xrp"`
						Value float64 `json:"value"`
						Text  string  `json:"text"`
					} `json:"premierPrice"`
				} `json:"price"`
				DeliveryOptionID       int           `json:"deliveryOptionId"`
				Name                   string        `json:"name"`
				AvailableDeliveryDates []interface{} `json:"availableDeliveryDates"`
				EstimatedDeliveryDate  string        `json:"estimatedDeliveryDate"`
				ExpiryDate             time.Time     `json:"expiryDate"`
				Messages               struct {
					Proposition struct {
						Message        string   `json:"message"`
						MessageContext []string `json:"messageContext"`
					} `json:"proposition"`
					Information struct {
						Message        string        `json:"message"`
						MessageContext []interface{} `json:"messageContext"`
					} `json:"information"`
					Delay                        interface{} `json:"delay"`
					BagAutoUpgradePrompt         interface{} `json:"bagAutoUpgradePrompt"`
					BagAutoUpgradeQualified      interface{} `json:"bagAutoUpgradeQualified"`
					CheckoutAutoUpgradeQualified interface{} `json:"checkoutAutoUpgradeQualified"`
					CheckoutAutoUpgradePrompt    interface{} `json:"checkoutAutoUpgradePrompt"`
				} `json:"messages"`
				IsDefault                    bool `json:"isDefault"`
				DeliveryMethodID             int  `json:"deliveryMethodId"`
				IsApplicableForPremierSaving bool `json:"isApplicableForPremierSaving"`
			} `json:"options"`
			ExcludedDeliveryMethods  []interface{} `json:"excludedDeliveryMethods"`
			SelectedDeliveryOptionID int           `json:"selectedDeliveryOptionId"`
			DutyMessage              interface{}   `json:"dutyMessage"`
			NominatedDate            interface{}   `json:"nominatedDate"`
		} `json:"delivery"`
	} `json:"bag"`
	Messages []interface{} `json:"messages"`
}

type countryCodeName struct {
	CountryCode string
	Name        string
}

var CountryCodeNameMap = []countryCodeName{{
	CountryCode: "AF",
	Name:        "Afghanistan",
}, {
	CountryCode: "AX",
	Name:        "Aland Islands",
}, {
	CountryCode: "AL",
	Name:        "Albania",
}, {
	CountryCode: "DZ",
	Name:        "Algeria",
}, {
	CountryCode: "AS",
	Name:        "American Samoa",
}, {
	CountryCode: "AD",
	Name:        "Andorra",
}, {
	CountryCode: "AO",
	Name:        "Angola",
}, {
	CountryCode: "AI",
	Name:        "Anguilla",
}, {
	CountryCode: "AQ",
	Name:        "Antarctica",
}, {
	CountryCode: "AG",
	Name:        "Antigua and Barbuda",
}, {
	CountryCode: "AR",
	Name:        "Argentina",
}, {
	CountryCode: "AM",
	Name:        "Armenia",
}, {
	CountryCode: "AW",
	Name:        "Aruba",
}, {
	CountryCode: "AT",
	Name:        "Austria",
}, {
	CountryCode: "AZ",
	Name:        "Azerbaijan",
}, {
	CountryCode: "BS",
	Name:        "Bahamas",
}, {
	CountryCode: "BH",
	Name:        "Bahrain",
}, {
	CountryCode: "BD",
	Name:        "Bangladesh",
}, {
	CountryCode: "BB",
	Name:        "Barbados",
}, {
	CountryCode: "BY",
	Name:        "Belarus",
}, {
	CountryCode: "BE",
	Name:        "Belgium",
}, {
	CountryCode: "BZ",
	Name:        "Belize",
}, {
	CountryCode: "BJ",
	Name:        "Benin",
}, {
	CountryCode: "BM",
	Name:        "Bermuda",
}, {
	CountryCode: "BT",
	Name:        "Bhutan",
}, {
	CountryCode: "BO",
	Name:        "Bolivia, Plurinational State of",
}, {
	CountryCode: "BQ",
	Name:        "Bonaire, Sint Eustatius and Saba",
}, {
	CountryCode: "BA",
	Name:        "Bosnia and Herzegovina",
}, {
	CountryCode: "BW",
	Name:        "Botswana",
}, {
	CountryCode: "BR",
	Name:        "Brazil",
}, {
	CountryCode: "IO",
	Name:        "British Indian Ocean Territory",
}, {
	CountryCode: "BN",
	Name:        "Brunei Darussalam",
}, {
	CountryCode: "BG",
	Name:        "Bulgaria",
}, {
	CountryCode: "BF",
	Name:        "Burkina Faso",
}, {
	CountryCode: "BI",
	Name:        "Burundi",
}, {
	CountryCode: "KH",
	Name:        "Cambodia",
}, {
	CountryCode: "CM",
	Name:        "Cameroon",
}, {
	CountryCode: "CA",
	Name:        "Canada",
}, {
	CountryCode: "CV",
	Name:        "Cape Verde",
}, {
	CountryCode: "KY",
	Name:        "Cayman Islands",
}, {
	CountryCode: "CF",
	Name:        "Central African Republic",
}, {
	CountryCode: "TD",
	Name:        "Chad",
}, {
	CountryCode: "CL",
	Name:        "Chile",
}, {
	CountryCode: "CX",
	Name:        "Christmas Island (Australia)",
}, {
	CountryCode: "CC",
	Name:        "Cocos (Keeling) Islands",
}, {
	CountryCode: "CO",
	Name:        "Colombia",
}, {
	CountryCode: "KM",
	Name:        "Comoros",
}, {
	CountryCode: "CD",
	Name:        "Congo, the Democratic Republic of the",
}, {
	CountryCode: "CG",
	Name:        "Congo, the Republic of",
}, {
	CountryCode: "CK",
	Name:        "Cook Islands",
}, {
	CountryCode: "CR",
	Name:        "Costa Rica",
}, {
	CountryCode: "CI",
	Name:        "Cote d'Ivoire",
}, {
	CountryCode: "HR",
	Name:        "Croatia",
}, {
	CountryCode: "CU",
	Name:        "Cuba",
}, {
	CountryCode: "CW",
	Name:        "Curacao",
}, {
	CountryCode: "CY",
	Name:        "Cyprus",
}, {
	CountryCode: "CZ",
	Name:        "Czech Republic",
}, {
	CountryCode: "KP",
	Name:        "Democratic People's Republic of Korea (North)",
}, {
	CountryCode: "DK",
	Name:        "Denmark",
}, {
	CountryCode: "DJ",
	Name:        "Djibouti",
}, {
	CountryCode: "DM",
	Name:        "Dominica",
}, {
	CountryCode: "DO",
	Name:        "Dominican Republic",
}, {
	CountryCode: "EC",
	Name:        "Ecuador",
}, {
	CountryCode: "EG",
	Name:        "Egypt",
}, {
	CountryCode: "SV",
	Name:        "El Salvador",
}, {
	CountryCode: "GQ",
	Name:        "Equatorial Guinea",
}, {
	CountryCode: "ER",
	Name:        "Eritrea",
}, {
	CountryCode: "EE",
	Name:        "Estonia",
}, {
	CountryCode: "SZ",
	Name:        "Eswatini",
}, {
	CountryCode: "ET",
	Name:        "Ethiopia",
}, {
	CountryCode: "FK",
	Name:        "Falkland Islands (Malvinas)",
}, {
	CountryCode: "FO",
	Name:        "Faroe Islands",
}, {
	CountryCode: "FJ",
	Name:        "Fiji",
}, {
	CountryCode: "FI",
	Name:        "Finland",
}, {
	CountryCode: "FR",
	Name:        "France",
}, {
	CountryCode: "GF",
	Name:        "French Guiana (Guyane)",
}, {
	CountryCode: "PF",
	Name:        "French Polynesia",
}, {
	CountryCode: "TF",
	Name:        "French Southern Territories",
}, {
	CountryCode: "GA",
	Name:        "Gabon",
}, {
	CountryCode: "GM",
	Name:        "Gambia",
}, {
	CountryCode: "GE",
	Name:        "Georgia",
}, {
	CountryCode: "DE",
	Name:        "Germany",
}, {
	CountryCode: "GH",
	Name:        "Ghana",
}, {
	CountryCode: "GI",
	Name:        "Gibraltar",
}, {
	CountryCode: "GR",
	Name:        "Greece",
}, {
	CountryCode: "GL",
	Name:        "Greenland",
}, {
	CountryCode: "GD",
	Name:        "Grenada",
}, {
	CountryCode: "GP",
	Name:        "Guadeloupe",
}, {
	CountryCode: "GU",
	Name:        "Guam",
}, {
	CountryCode: "GT",
	Name:        "Guatemala",
}, {
	CountryCode: "GN",
	Name:        "Guinea",
}, {
	CountryCode: "GW",
	Name:        "Guinea-Bissau",
}, {
	CountryCode: "GY",
	Name:        "Guyana, Co-operative Republic of",
}, {
	CountryCode: "HT",
	Name:        "Haiti",
}, {
	CountryCode: "VA",
	Name:        "Holy See (Vatican City State)",
}, {
	CountryCode: "HN",
	Name:        "Honduras",
}, {
	CountryCode: "HK",
	Name:        "Hong Kong",
}, {
	CountryCode: "HU",
	Name:        "Hungary",
}, {
	CountryCode: "IS",
	Name:        "Iceland",
}, {
	CountryCode: "IN",
	Name:        "India",
}, {
	CountryCode: "ID",
	Name:        "Indonesia",
}, {
	CountryCode: "IR",
	Name:        "Iran, Islamic Republic of",
}, {
	CountryCode: "IQ",
	Name:        "Iraq",
}, {
	CountryCode: "IE",
	Name:        "Ireland, Republic of",
}, {
	CountryCode: "IL",
	Name:        "Israel",
}, {
	CountryCode: "IT",
	Name:        "Italy",
}, {
	CountryCode: "JM",
	Name:        "Jamaica",
}, {
	CountryCode: "JP",
	Name:        "Japan",
}, {
	CountryCode: "JO",
	Name:        "Jordan",
}, {
	CountryCode: "KZ",
	Name:        "Kazakhstan",
}, {
	CountryCode: "KE",
	Name:        "Kenya",
}, {
	CountryCode: "KI",
	Name:        "Kiribati",
}, {
	CountryCode: "KR",
	Name:        "Korea, Republic of (South Korea)",
}, {
	CountryCode: "XK",
	Name:        "Kosovo",
}, {
	CountryCode: "KW",
	Name:        "Kuwait",
}, {
	CountryCode: "KG",
	Name:        "Kyrgyzstan",
}, {
	CountryCode: "LA",
	Name:        "Lao People's Democratic Republic",
}, {
	CountryCode: "LV",
	Name:        "Latvia",
}, {
	CountryCode: "LB",
	Name:        "Lebanon",
}, {
	CountryCode: "LS",
	Name:        "Lesotho",
}, {
	CountryCode: "LR",
	Name:        "Liberia",
}, {
	CountryCode: "LY",
	Name:        "Libya",
}, {
	CountryCode: "LI",
	Name:        "Liechtenstein",
}, {
	CountryCode: "LT",
	Name:        "Lithuania",
}, {
	CountryCode: "LU",
	Name:        "Luxembourg",
}, {
	CountryCode: "MO",
	Name:        "Macao",
}, {
	CountryCode: "MG",
	Name:        "Madagascar",
}, {
	CountryCode: "MW",
	Name:        "Malawi",
}, {
	CountryCode: "MY",
	Name:        "Malaysia",
}, {
	CountryCode: "MV",
	Name:        "Maldives",
}, {
	CountryCode: "ML",
	Name:        "Mali",
}, {
	CountryCode: "MT",
	Name:        "Malta",
}, {
	CountryCode: "MH",
	Name:        "Marshall Islands",
}, {
	CountryCode: "MQ",
	Name:        "Martinique",
}, {
	CountryCode: "MR",
	Name:        "Mauritania",
}, {
	CountryCode: "MU",
	Name:        "Mauritius",
}, {
	CountryCode: "YT",
	Name:        "Mayotte",
}, {
	CountryCode: "MX",
	Name:        "Mexico",
}, {
	CountryCode: "FM",
	Name:        "Micronesia, Federated States of",
}, {
	CountryCode: "MD",
	Name:        "Moldova, Republic of",
}, {
	CountryCode: "MC",
	Name:        "Monaco",
}, {
	CountryCode: "MN",
	Name:        "Mongolia",
}, {
	CountryCode: "ME",
	Name:        "Montenegro",
}, {
	CountryCode: "MS",
	Name:        "Montserrat",
}, {
	CountryCode: "MA",
	Name:        "Morocco",
}, {
	CountryCode: "MZ",
	Name:        "Mozambique",
}, {
	CountryCode: "MM",
	Name:        "Myanmar",
}, {
	CountryCode: "NA",
	Name:        "Namibia",
}, {
	CountryCode: "NR",
	Name:        "Nauru",
}, {
	CountryCode: "NP",
	Name:        "Nepal",
}, {
	CountryCode: "NL",
	Name:        "Netherlands",
}, {
	CountryCode: "NC",
	Name:        "New Caledonia",
}, {
	CountryCode: "NZ",
	Name:        "New Zealand",
}, {
	CountryCode: "NI",
	Name:        "Nicaragua",
}, {
	CountryCode: "NE",
	Name:        "Niger",
}, {
	CountryCode: "NG",
	Name:        "Nigeria",
}, {
	CountryCode: "NU",
	Name:        "Niue",
}, {
	CountryCode: "NF",
	Name:        "Norfolk Island",
}, {
	CountryCode: "MK",
	Name:        "North Macedonia",
}, {
	CountryCode: "MP",
	Name:        "Northern Mariana Islands",
}, {
	CountryCode: "NO",
	Name:        "Norway",
}, {
	CountryCode: "OM",
	Name:        "Oman",
}, {
	CountryCode: "PK",
	Name:        "Pakistan",
}, {
	CountryCode: "PW",
	Name:        "Palau",
}, {
	CountryCode: "PS",
	Name:        "Palestine",
}, {
	CountryCode: "PA",
	Name:        "Panama",
}, {
	CountryCode: "PG",
	Name:        "Papua New Guinea",
}, {
	CountryCode: "PY",
	Name:        "Paraguay",
}, {
	CountryCode: "PE",
	Name:        "Peru",
}, {
	CountryCode: "PH",
	Name:        "Philippines",
}, {
	CountryCode: "PN",
	Name:        "Pitcairn",
}, {
	CountryCode: "PL",
	Name:        "Poland",
}, {
	CountryCode: "PT",
	Name:        "Portugal",
}, {
	CountryCode: "PR",
	Name:        "Puerto Rico",
}, {
	CountryCode: "QA",
	Name:        "Qatar",
}, {
	CountryCode: "RE",
	Name:        "Reunion",
}, {
	CountryCode: "RO",
	Name:        "Romania",
}, {
	CountryCode: "RW",
	Name:        "Rwanda",
}, {
	CountryCode: "BL",
	Name:        "Saint Barthelemy",
}, {
	CountryCode: "SH",
	Name:        "Saint Helena, Ascension and Tristan da Cunha",
}, {
	CountryCode: "KN",
	Name:        "Saint Kitts and Nevis",
}, {
	CountryCode: "LC",
	Name:        "Saint Lucia",
}, {
	CountryCode: "MF",
	Name:        "Saint Martin (French part)",
}, {
	CountryCode: "PM",
	Name:        "Saint Pierre and Miquelon",
}, {
	CountryCode: "VC",
	Name:        "Saint Vincent and the Grenadines",
}, {
	CountryCode: "WS",
	Name:        "Samoa",
}, {
	CountryCode: "SM",
	Name:        "San Marino",
}, {
	CountryCode: "ST",
	Name:        "Sao Tome and Principe",
}, {
	CountryCode: "SA",
	Name:        "Saudi Arabia",
}, {
	CountryCode: "SN",
	Name:        "Senegal",
}, {
	CountryCode: "RS",
	Name:        "Serbia",
}, {
	CountryCode: "SC",
	Name:        "Seychelles",
}, {
	CountryCode: "SL",
	Name:        "Sierra Leone",
}, {
	CountryCode: "SG",
	Name:        "Singapore",
}, {
	CountryCode: "SX",
	Name:        "Sint Maarten (Dutch part)",
}, {
	CountryCode: "SK",
	Name:        "Slovakia",
}, {
	CountryCode: "SI",
	Name:        "Slovenia",
}, {
	CountryCode: "SB",
	Name:        "Solomon Islands",
}, {
	CountryCode: "SO",
	Name:        "Somalia",
}, {
	CountryCode: "ZA",
	Name:        "South Africa",
}, {
	CountryCode: "GS",
	Name:        "South Georgia and the South Sandwich Islands",
}, {
	CountryCode: "SS",
	Name:        "South Sudan",
}, {
	CountryCode: "ES",
	Name:        "Spain",
}, {
	CountryCode: "LK",
	Name:        "Sri Lanka",
}, {
	CountryCode: "SD",
	Name:        "Sudan",
}, {
	CountryCode: "SR",
	Name:        "Suriname",
}, {
	CountryCode: "SJ",
	Name:        "Svalbard and Jan Mayen",
}, {
	CountryCode: "SE",
	Name:        "Sweden",
}, {
	CountryCode: "CH",
	Name:        "Switzerland",
}, {
	CountryCode: "SY",
	Name:        "Syrian Arab Republic",
}, {
	CountryCode: "TW",
	Name:        "Taiwan",
}, {
	CountryCode: "TJ",
	Name:        "Tajikistan",
}, {
	CountryCode: "TZ",
	Name:        "Tanzania, United Republic of",
}, {
	CountryCode: "TH",
	Name:        "Thailand",
}, {
	CountryCode: "TL",
	Name:        "Timor-Leste",
}, {
	CountryCode: "TG",
	Name:        "Togo",
}, {
	CountryCode: "TK",
	Name:        "Tokelau",
}, {
	CountryCode: "TO",
	Name:        "Tonga",
}, {
	CountryCode: "TT",
	Name:        "Trinidad and Tobago",
}, {
	CountryCode: "TN",
	Name:        "Tunisia",
}, {
	CountryCode: "TR",
	Name:        "Turkey",
}, {
	CountryCode: "TM",
	Name:        "Turkmenistan",
}, {
	CountryCode: "TC",
	Name:        "Turks and Caicos Islands",
}, {
	CountryCode: "TV",
	Name:        "Tuvalu",
}, {
	CountryCode: "UG",
	Name:        "Uganda",
}, {
	CountryCode: "GB",
	Name:        "UK",
}, {
	CountryCode: "UA",
	Name:        "Ukraine",
}, {
	CountryCode: "AE",
	Name:        "United Arab Emirates",
}, {
	CountryCode: "UM",
	Name:        "United States Minor Outlying Islands",
}, {
	CountryCode: "UY",
	Name:        "Uruguay",
}, {
	CountryCode: "UZ",
	Name:        "Uzbekistan",
}, {
	CountryCode: "VU",
	Name:        "Vanuatu",
}, {
	CountryCode: "VE",
	Name:        "Venezuela, Bolivarian Republic of",
}, {
	CountryCode: "VN",
	Name:        "Vietnam",
}, {
	CountryCode: "VG",
	Name:        "Virgin Islands, British",
}, {
	CountryCode: "VI",
	Name:        "Virgin Islands, U.S.",
}, {
	CountryCode: "WF",
	Name:        "Wallis and Futuna",
}, {
	CountryCode: "EH",
	Name:        "Western Sahara",
}, {
	CountryCode: "YE",
	Name:        "Yemen",
}, {
	CountryCode: "ZM",
	Name:        "Zambia",
}, {
	CountryCode: "ZW",
	Name:        "Zimbabwe",
}}
