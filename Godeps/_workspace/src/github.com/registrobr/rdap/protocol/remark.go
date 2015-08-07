package protocol

// https://tools.ietf.org/html/rfc7483#section-10.2.1
const (
	// RemarkTypeResultTruncatedAuthorization the list of results does not
	// contain all results due to lack of authorization.  This may indicate to
	// some clients that proper authorization will yield a longer result set
	RemarkTypeResultTruncatedAuthorization RemarkType = "result set truncated due to authorization"

	// RemarkTypeResultTruncatedExcessiveLoad the list of results does not
	// contain all results due to an excessively heavy load on the server. This
	// may indicate to some clients that requerying at a later time will yield a
	// longer result set
	RemarkTypeResultTruncatedExcessiveLoad RemarkType = "result set truncated due to excessive load"

	// RemarkTypeResultTruncatedUnexplainableReasons the list of results does not
	// contain all results for an unexplainable reason. This may indicate to some
	// clients that requerying for any reason will not yield a longer result set
	RemarkTypeResultTruncatedUnexplainableReasons RemarkType = "result set truncated due to unexplainable reasons"

	// RemarkTypeObjectTruncatedAuthorization The object does not contain all
	// data due to lack of authorization
	RemarkTypeObjectTruncatedAuthorization RemarkType = "object truncated due to authorization"

	// RemarkTypeObjectTruncatedExcessiveLoad the object does not contain all
	// data due to an excessively heavy load on the server. This may indicate to
	// some clients that requerying at a later time will yield all data of the
	// object
	RemarkTypeObjectTruncatedExcessiveLoad RemarkType = "object truncated due to excessive load"

	// RemarkTypeObjectTruncatedUnexplainableReasons the object does not
	// contain all data for an unexplainable reason
	RemarkTypeObjectTruncatedUnexplainableReasons RemarkType = "object truncated due to unexplainable reasons"
)

// RemarkType stores one of the possible remark types as listed in RFC 7483,
// section 10.2.1
type RemarkType string

// Remark describes Remarks as it is in RFC 7483, section 4.3
type Remark struct {
	Type        string   `json:"type,omitempty"`
	Description []string `json:"description,omitempty"`
}
