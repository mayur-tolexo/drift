package lib

import (
	"github.com/go-pg/pg"
	"github.com/jordan-wright/email"
	"github.com/olivere/elastic"
)

//Error Model
type Error struct {
	Code     int    `json:"code"`
	Msg      string `json:"Message"`
	Trace    string `json:"trace,omitempty"`
	DebugMsg string `json:"debug_msg,omitempty"`
}

//Resp : Response Model
type Resp struct {
	Data   interface{} `json:"data,omitempty"`
	Msg    interface{} `json:"message,omitempty"`
	Error  error       `json:"error,omitempty"`
	Status bool        `json:"status"`
}

//MsgResp : Message response for create/update/delete api
type MsgResp struct {
	Data interface{}
	Msg  string
}

//TestResp : Response Model for test files
type TestResp struct {
	Data   interface{} `json:"data"`
	Msg    interface{} `json:"message"`
	Error  Error       `json:"error"`
	Status bool        `json:"status"`
}

//AttributeCodeOption : Attribute Code and Attribute Option ID fetch model from database
type AttributeCodeOption struct {
	AttributeCode string      `sql:"attribute_code"`
	OptionID      interface{} `sql:"attribute_option_id"`
}

//SelectFilter : Get request filter
type SelectFilter struct {
	Alias  string
	Filter map[string]interface{}
	Offset int
	Limit  int
}

//FilterDet : select query filter alias name of field, operator to use and value
type FilterDet struct {
	Key             string
	Alias           string
	Operator        string
	Value           string
	ValidTag        string
	DefaultOperator string
	DefaultValue    string
	Sort            bool
	Skip            bool
	Or              bool
	Unescape        bool
	ESFilterType    string
}

//AdvFilter : advance select filter
type AdvFilter struct {
	Alias         string
	Filter        map[string]FilterDet
	Offset        int
	Limit         int
	IncludeFields string
	ExcludeFields string
}

//DB is a wrap struct of go-pg
type DB struct {
	*pg.DB
	query       string
	value       []interface{}
	bindValue   []interface{}
	condition   map[string][]interface{}
	orCondition map[string][]interface{}
	group       []string
	order       []string
	limit       string
	offset      string
}

//Email struct for mail
type Email struct {
	*email.Email
}

//AttributeInfo : attribute info struct
type AttributeInfo struct {
	AttrCodeString string `sql:"attr_codes"`
	AttrLevel      string `sql:"attr_level"`
}

//RemoveAttr : remove attribute code request format
type RemoveAttr struct {
	Attribute []string `json:"attribute" validate:"required,min=1,dive,alphanumunderscore"`
	CreatedBy int      `json:"-"`
}

//ESJob : elasticsearch index job format
// type can be product or article
type ESJob struct {
	ID   []int
	Type string
}

//CategoryAttribute : model of category and mandatory attribute fetch from database
type CategoryAttribute struct {
	CategoryID  int    `sql:"fk_category_id"`
	AttributeID string `sql:"fk_attribute_id"`
}

//ES is a wrap struct of elastic client
type ES struct {
	*elastic.Client
}

//wrapESLogger is wrapped elastic logger
type wrapESLogger struct {
	elastic.Logger
}
