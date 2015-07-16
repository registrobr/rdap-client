package protocol

// https://tools.ietf.org/html/rfc7483#section-10.2.1
const (
	RemarkTypeResultTruncatedAuthorization        RemarkType = "result set truncated due to authorization"
	RemarkTypeResultTruncatedExcessiveLoad        RemarkType = "result set truncated due to excessive load"
	RemarkTypeResultTruncatedUnexplainableReasons RemarkType = "result set truncated due to unexplainable reasons"
	RemarkTypeObjectTruncatedAuthorization        RemarkType = "object truncated due to authorization"
	RemarkTypeObjectTruncatedExcessiveLoad        RemarkType = "object truncated due to excessive load"
	RemarkTypeObjectTruncatedUnexplainableReasons RemarkType = "object truncated due to unexplainable reasons"
)

type RemarkType string

type Remark struct {
	Type        string   `json:"type,omitempty"`
	Description []string `json:"description,omitempty"`
}
