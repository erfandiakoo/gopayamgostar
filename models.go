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

type GetRequest struct {
	ID                      string `json:"id"`
	ShowPreviews            bool   `json:"showPreviews"`
	ShowExtendedPreviews    bool   `json:"showExtendedPreviews"`
	IncludeProcessLifePaths bool   `json:"includeProcessLifePaths"`
	IncludeColor            bool   `json:"includeColor"`
	IncludeTags             bool   `json:"includeTags"`
	IncludeListFields       bool   `json:"includeListFields"`
}

type PersonInfo struct {
	FirstName                 string             `json:"firstName"`
	LastName                  string             `json:"lastName"`
	BirthDate                 interface{}        `json:"birthDate"`
	Gender                    string             `json:"gender"`
	PersonPrefix              string             `json:"personPrefix"`
	NationalCode              string             `json:"nationalCode"`
	PreferredContactType      string             `json:"preferredContactType"`
	FacebookUsername          string             `json:"facebookUsername"`
	Organizations             []interface{}      `json:"organizations"`
	NickName                  string             `json:"nickName"`
	PhoneContacts             []PhoneContact     `json:"phoneContacts"`
	AddressContacts           []interface{}      `json:"addressContacts"`
	Email                     string             `json:"email"`
	AlternativeEmail          string             `json:"alternativeEmail"`
	Website                   string             `json:"website"`
	SourceTypeName            string             `json:"sourceTypeName"`
	CustomerNumber            string             `json:"customerNumber"`
	ColorName                 string             `json:"colorName"`
	Classification            string             `json:"classification"`
	CustomerDate              interface{}        `json:"customerDate"`
	Balance                   int64              `json:"balance"`
	IdentityTypeName          string             `json:"identityTypeName"`
	Categories                []Category         `json:"categories"`
	SupportUsername           string             `json:"supportUsername"`
	SaleUsername              string             `json:"saleUsername"`
	OtherUsername             string             `json:"otherUsername"`
	CRMID                     string             `json:"crmId"`
	CRMObjectTypeName         interface{}        `json:"crmObjectTypeName"`
	CRMObjectTypeCode         string             `json:"crmObjectTypeCode"`
	CRMObjectTypeIndex        int64              `json:"crmObjectTypeIndex"`
	CRMObjectTypeID           string             `json:"crmObjectTypeId"`
	ParentCRMObjectID         interface{}        `json:"parentCrmObjectId"`
	ExtendedProperties        []ExtendedProperty `json:"extendedProperties"`
	ProcessLifePaths          []interface{}      `json:"processLifePaths"`
	CreatDate                 time.Time          `json:"creatDate"`
	ModifyDate                time.Time          `json:"modifyDate"`
	RefID                     string             `json:"refId"`
	StageID                   interface{}        `json:"stageId"`
	IdentityID                string             `json:"identityId"`
	Description               string             `json:"description"`
	Subject                   string             `json:"subject"`
	ModifierIDPreview         interface{}        `json:"modifierIdPreview"`
	CreatorIDPreview          interface{}        `json:"creatorIdPreview"`
	CRMObjectTypeIndexPreview interface{}        `json:"crmObjectTypeIndexPreview"`
	IdentityIDPreview         interface{}        `json:"identityIdPreview"`
	AssignedToIDPreview       interface{}        `json:"assignedToIdPreview"`
	IncludedFields            IncludedFields     `json:"includedFields"`
}

type AreasOfInterest struct {
	Name string `json:"Name"`
}

type FormInfo struct {
	CRMID                     string              `json:"CrmId"`
	CRMObjectTypeIndexPreview AssignedToIDPreview `json:"CrmObjectTypeIndexPreview"`
	CRMObjectTypeIndex        int64               `json:"CrmObjectTypeIndex"`
	CRMObjectTypeName         AssignedToIDPreview `json:"CrmObjectTypeName"`
	CRMObjectTypeID           string              `json:"CrmObjectTypeId"`
	CRMObjectTypeCode         string              `json:"CrmObjectTypeCode"`
	ParentCRMObjectID         interface{}         `json:"ParentCrmObjectId"`
	ExtendedProperties        []ExtendedProperty  `json:"ExtendedProperties"`
	Tags                      []interface{}       `json:"Tags"`
	RefID                     string              `json:"RefId"`
	StageID                   interface{}         `json:"StageId"`
	IdentityIDPreview         AssignedToIDPreview `json:"IdentityIdPreview"`
	IdentityID                string              `json:"IdentityId"`
	Description               string              `json:"Description"`
	Subject                   string              `json:"Subject"`
	ProcessLifePaths          []interface{}       `json:"ProcessLifePaths"`
	Color                     interface{}         `json:"Color"`
	ModifierIDPreview         AssignedToIDPreview `json:"ModifierIdPreview"`
	ModifierID                string              `json:"ModifierId"`
	CreatorIDPreview          AssignedToIDPreview `json:"CreatorIdPreview"`
	CreatorID                 string              `json:"CreatorId"`
	AssignedToIDPreview       AssignedToIDPreview `json:"AssignedToIdPreview"`
	AssignedToID              interface{}         `json:"AssignedToId"`
}

type AssignedToIDPreview struct {
	Name string `json:"Name"`
}

type CreatePurchase struct {
	CrmId              string             `json:"crmId,omitempty"`
	CRMObjectTypeCode  string             `json:"crmObjectTypeCode"`
	Details            []Detail           `json:"details"`
	Discount           int64              `json:"discount"`
	FinalValue         int64              `json:"finalValue"`
	Toll               int64              `json:"toll"`
	TotalValue         int64              `json:"totalValue"`
	Vat                int64              `json:"vat"`
	ParentCRMObjectID  *string            `json:"parentCrmObjectId"`
	ExtendedProperties []ExtendedProperty `json:"extendedProperties"`
	Tags               *[]string          `json:"tags"`
	RefID              *string            `json:"refId"`
	StageID            *string            `json:"stageId"`
	ColorID            int64              `json:"colorId"`
	IdentityID         string             `json:"identityId"`
	Description        *string            `json:"description"`
	Subject            *string            `json:"subject"`
	AssignedToUserName *string            `json:"assignedToUserName"`
	Number             *string            `json:"number"`
	PriceListName      *string            `json:"priceListName"`
	AdditionalCosts    *string            `json:"additionalCosts"`
	InvoiceDate        *string            `json:"invoiceDate"`
	ExpireDate         *string            `json:"expireDate"`
	DiscountPercent    *string            `json:"discountPercent"`
	RelatedQuoteID     *string            `json:"relatedQuoteId"`
}

type Detail struct {
	IsService           bool   `json:"isService"`
	BaseUnitPrice       int64  `json:"baseUnitPrice"`
	FinalUnitPrice      int64  `json:"finalUnitPrice"`
	Count               int64  `json:"count"`
	ReturnedCount       int64  `json:"returnedCount"`
	TotalUnitPrice      int64  `json:"totalUnitPrice"`
	TotalDiscount       int64  `json:"totalDiscount"`
	TotalVat            int64  `json:"totalVat"`
	TotalToll           int64  `json:"totalToll"`
	ProductCode         string `json:"productCode"`
	ProductID           string `json:"productId"`
	ProductName         string `json:"productName"`
	DiscountPercent     string `json:"discountPercent"`
	DetailDescription   string `json:"detailDescription"`
	ProductUnitTypeName string `json:"productUnitTypeName"`
}

type DeleteRequest struct {
	Id     string `json:"id"`
	Option int    `json:"option"`
}

type FindResponse struct {
	Data  []PersonInfo `json:"data"`
	Total int64        `json:"total"`
}

type FindRequest struct {
	TypeKey    string  `json:"typeKey"`
	Queries    []Query `json:"queries"`
	PageNumber int64   `json:"pageNumber"`
	PageSize   int64   `json:"pageSize"`
}

type Query struct {
	LogicalOperator     int    `json:"logicalOperator"`
	Operator            int    `json:"operator"`
	LeafNegate          bool   `json:"leafNegate,omitempty"`
	Field               string `json:"field"`
	FieldOperator       int    `json:"fieldOperator,omitempty"`
	Value               string `json:"value"`
	LeafLogicalOperator int    `json:"leafLogicalOperator,omitempty"`
}

type Datum struct {
	FirstName                 string             `json:"firstName"`
	LastName                  string             `json:"lastName"`
	BirthDate                 interface{}        `json:"birthDate"`
	Gender                    string             `json:"gender"`
	PersonPrefix              string             `json:"personPrefix"`
	NationalCode              string             `json:"nationalCode"`
	PreferredContactType      string             `json:"preferredContactType"`
	FacebookUsername          string             `json:"facebookUsername"`
	Organizations             []interface{}      `json:"organizations"`
	NickName                  string             `json:"nickName"`
	PhoneContacts             []PhoneContact     `json:"phoneContacts"`
	AddressContacts           []interface{}      `json:"addressContacts"`
	Email                     string             `json:"email"`
	AlternativeEmail          string             `json:"alternativeEmail"`
	Website                   string             `json:"website"`
	SourceTypeName            string             `json:"sourceTypeName"`
	CustomerNumber            string             `json:"customerNumber"`
	ColorName                 string             `json:"colorName"`
	Classification            string             `json:"classification"`
	CustomerDate              interface{}        `json:"customerDate"`
	Balance                   int64              `json:"balance"`
	IdentityTypeName          string             `json:"identityTypeName"`
	Categories                []Category         `json:"categories"`
	SupportUsername           string             `json:"supportUsername"`
	SaleUsername              string             `json:"saleUsername"`
	OtherUsername             string             `json:"otherUsername"`
	CRMID                     string             `json:"crmId"`
	CRMObjectTypeName         interface{}        `json:"crmObjectTypeName"`
	CRMObjectTypeCode         string             `json:"crmObjectTypeCode"`
	CRMObjectTypeIndex        int64              `json:"crmObjectTypeIndex"`
	CRMObjectTypeID           string             `json:"crmObjectTypeId"`
	ParentCRMObjectID         interface{}        `json:"parentCrmObjectId"`
	ExtendedProperties        []ExtendedProperty `json:"extendedProperties"`
	ProcessLifePaths          []interface{}      `json:"processLifePaths"`
	CreatDate                 time.Time          `json:"creatDate"`
	ModifyDate                time.Time          `json:"modifyDate"`
	RefID                     string             `json:"refId"`
	StageID                   interface{}        `json:"stageId"`
	IdentityID                string             `json:"identityId"`
	Description               string             `json:"description"`
	Subject                   string             `json:"subject"`
	ModifierIDPreview         interface{}        `json:"modifierIdPreview"`
	CreatorIDPreview          interface{}        `json:"creatorIdPreview"`
	CRMObjectTypeIndexPreview interface{}        `json:"crmObjectTypeIndexPreview"`
	IdentityIDPreview         interface{}        `json:"identityIdPreview"`
	AssignedToIDPreview       interface{}        `json:"assignedToIdPreview"`
	IncludedFields            IncludedFields     `json:"includedFields"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	Type string `json:"type"`
}

type ExtendedProperty struct {
	Value   string      `json:"value"`
	UserKey string      `json:"userKey"`
	Preview interface{} `json:"preview"`
}

type IncludedFields struct {
}

type PhoneContact struct {
	PhoneType       string `json:"phoneType"`
	PhoneNumber     string `json:"phoneNumber"`
	ContinuedNumber string `json:"continuedNumber"`
	Extension       string `json:"extension"`
	ID              string `json:"id"`
	Default         bool   `json:"default"`
}
