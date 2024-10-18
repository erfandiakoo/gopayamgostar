package gopayamgostar

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
)

// GetQueryParams converts the struct to map[string]string
// The fields tags must have `json:"<name>,string,omitempty"` format for all types, except strings
// The string fields must have: `json:"<name>,omitempty"`. The `json:"<name>,string,omitempty"` tag for string field
// will add additional double quotes.
// "string" tag allows to convert the non-string fields of a structure to map[string]string.
// "omitempty" allows to skip the fields with default values.
func GetQueryParams(s interface{}) (map[string]string, error) {
	// if obj, ok := s.(GetGroupsParams); ok {
	// 	obj.OnMarshal()
	// 	s = obj
	// }
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	var res map[string]string
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// StringOrArray represents a value that can either be a string or an array of strings
type StringOrArray []string

// UnmarshalJSON unmarshals a string or an array object from a JSON array or a JSON string
func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	if len(data) > 1 && data[0] == '[' {
		var obj []string
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}
		*s = StringOrArray(obj)
		return nil
	}

	var obj string
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*s = StringOrArray([]string{obj})
	return nil
}

// MarshalJSON converts the array of strings to a JSON array or JSON string if there is only one item in the array
func (s *StringOrArray) MarshalJSON() ([]byte, error) {
	if len(*s) == 1 {
		return json.Marshal([]string(*s)[0])
	}
	return json.Marshal([]string(*s))
}

// EnforcedString can be used when the expected value is string but Keycloak in some cases gives you mixed types
type EnforcedString string

// UnmarshalJSON modify data as string before json unmarshal
func (s *EnforcedString) UnmarshalJSON(data []byte) error {
	if data[0] != '"' {
		// Escape unescaped quotes
		data = bytes.ReplaceAll(data, []byte(`"`), []byte(`\"`))
		data = bytes.ReplaceAll(data, []byte(`\\"`), []byte(`\"`))

		// Wrap data in quotes
		data = append([]byte(`"`), data...)
		data = append(data, []byte(`"`)...)
	}

	var val string
	err := json.Unmarshal(data, &val)
	*s = EnforcedString(val)
	return err
}

// MarshalJSON return json marshal
func (s *EnforcedString) MarshalJSON() ([]byte, error) {
	return json.Marshal(*s)
}

// APIErrType is a field containing more specific API error types
// that may be checked by the receiver.
type APIErrType string

const (
	// APIErrTypeUnknown is for API errors that are not strongly
	// typed.
	APIErrTypeUnknown APIErrType = "unknown"

	// APIErrTypeInvalidGrant corresponds with Keycloak's
	// OAuthErrorException due to "invalid_grant".
	APIErrTypeInvalidGrant = "oauth: invalid grant"
)

// ParseAPIErrType is a convenience method for returning strongly
// typed API errors.
func ParseAPIErrType(err error) APIErrType {
	if err == nil {
		return APIErrTypeUnknown
	}
	switch {
	case strings.Contains(err.Error(), "invalid_grant"):
		return APIErrTypeInvalidGrant
	default:
		return APIErrTypeUnknown
	}
}

// APIError holds message and statusCode for api errors
type APIError struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Type    APIErrType `json:"type"`
}

// Error stringifies the APIError
func (apiError APIError) Error() string {
	return apiError.Message
}

type AuthRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	DeviceId     string `json:"deviceId"`
	PlatformType int    `json:"platformType"`
}

type PersonInfo struct {
	FirstName                 string          `json:"FirstName"`
	LastName                  string          `json:"LastName"`
	BirthDate                 interface{}     `json:"BirthDate"`
	Gender                    AreasOfInterest `json:"Gender"`
	GenderIndex               interface{}     `json:"GenderIndex"`
	PersonPrefix              AreasOfInterest `json:"PersonPrefix"`
	PersonPrefixIndex         interface{}     `json:"PersonPrefixIndex"`
	Degree                    AreasOfInterest `json:"Degree"`
	DegreeIndex               interface{}     `json:"DegreeIndex"`
	CreditType                AreasOfInterest `json:"CreditType"`
	CreditTypeIndex           interface{}     `json:"CreditTypeIndex"`
	NationalCode              string          `json:"NationalCode"`
	Spouse                    string          `json:"Spouse"`
	PreferredContactType      AreasOfInterest `json:"PreferredContactType"`
	PreferredContactTypeIndex interface{}     `json:"PreferredContactTypeIndex"`
	PaymentStatusType         AreasOfInterest `json:"PaymentStatusType"`
	PaymentStatusTypeIndex    interface{}     `json:"PaymentStatusTypeIndex"`
	AreasOfInterest           AreasOfInterest `json:"AreasOfInterest"`
	AreasOfInterestIndex      interface{}     `json:"AreasOfInterestIndex"`
	Hobbies                   string          `json:"Hobbies"`
	Children                  string          `json:"Children"`
	MannerType                AreasOfInterest `json:"MannerType"`
	MannerTypeIndex           interface{}     `json:"MannerTypeIndex"`
	FacebookUsername          string          `json:"FacebookUsername"`
	Organizations             []interface{}   `json:"Organizations"`
	NickName                  string          `json:"NickName"`
	PhoneContacts             []PhoneContact  `json:"PhoneContacts"`
	AddressContacts           []interface{}   `json:"AddressContacts"`
	Email                     string          `json:"Email"`
	AlternativeEmail          string          `json:"AlternativeEmail"`
	Website                   string          `json:"Website"`
	CustomerNumber            string          `json:"CustomerNumber"`
	CustomerDate              interface{}     `json:"CustomerDate"`
	Balance                   interface{}     `json:"Balance"`
	Categories                []Category      `json:"Categories"`
	DontSMS                   interface{}     `json:"DontSms"`
	DontSocialSMS             interface{}     `json:"DontSocialSms"`
	DontPhoneCall             interface{}     `json:"DontPhoneCall"`
	DontEmail                 interface{}     `json:"DontEmail"`
	DontFax                   interface{}     `json:"DontFax"`
	SupportOperatorUserID     interface{}     `json:"SupportOperatorUserId"`
	SaleOperatorUserID        interface{}     `json:"SaleOperatorUserId"`
	OtherOperatorUserID       interface{}     `json:"OtherOperatorUserId"`
	SourceTypeIndex           interface{}     `json:"SourceTypeIndex"`
	ClassificationID          interface{}     `json:"ClassificationId"`
	IdentityTypeName          AreasOfInterest `json:"IdentityTypeName"`
	Classification            AreasOfInterest `json:"Classification"`
	ColorName                 AreasOfInterest `json:"ColorName"`
	SourceTypeName            AreasOfInterest `json:"SourceTypeName"`
	SupportOperatorUsername   AreasOfInterest `json:"SupportOperatorUsername"`
	SaleOperatorUsername      AreasOfInterest `json:"SaleOperatorUsername"`
	OtherOperatorUsername     AreasOfInterest `json:"OtherOperatorUsername"`
	CRMID                     string          `json:"CrmId"`
	CRMObjectTypeIndexPreview AreasOfInterest `json:"CrmObjectTypeIndexPreview"`
	CRMObjectTypeIndex        int64           `json:"CrmObjectTypeIndex"`
	CRMObjectTypeName         AreasOfInterest `json:"CrmObjectTypeName"`
	CRMObjectTypeID           string          `json:"CrmObjectTypeId"`
	CRMObjectTypeCode         string          `json:"CrmObjectTypeCode"`
	ParentCRMObjectID         interface{}     `json:"ParentCrmObjectId"`
	ExtendedProperties        []interface{}   `json:"ExtendedProperties"`
	CreatDate                 time.Time       `json:"CreatDate"`
	ModifyDate                time.Time       `json:"ModifyDate"`
	Tags                      []interface{}   `json:"Tags"`
	RefID                     string          `json:"RefId"`
	StageID                   interface{}     `json:"StageId"`
	IdentityIDPreview         AreasOfInterest `json:"IdentityIdPreview"`
	IdentityID                interface{}     `json:"IdentityId"`
	Description               string          `json:"Description"`
	Subject                   string          `json:"Subject"`
	ProcessLifePaths          []interface{}   `json:"ProcessLifePaths"`
	Color                     interface{}     `json:"Color"`
	ModifierIDPreview         AreasOfInterest `json:"ModifierIdPreview"`
	ModifierID                string          `json:"ModifierId"`
	CreatorIDPreview          AreasOfInterest `json:"CreatorIdPreview"`
	CreatorID                 string          `json:"CreatorId"`
	AssignedToIDPreview       AreasOfInterest `json:"AssignedToIdPreview"`
	AssignedToID              interface{}     `json:"AssignedToId"`
}

type AreasOfInterest struct {
	Name string `json:"Name"`
}

type Category struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
	Key  string `json:"Key"`
	Type string `json:"Type"`
}

type PhoneContact struct {
	PhoneType       string `json:"PhoneType"`
	PhoneNumber     string `json:"PhoneNumber"`
	ContinuedNumber string `json:"ContinuedNumber"`
	Extension       string `json:"Extension"`
	ID              string `json:"Id"`
	Default         bool   `json:"Default"`
}
