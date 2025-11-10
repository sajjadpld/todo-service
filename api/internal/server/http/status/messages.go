package status

type HttpMappedStatus string

// general

const (
	Success      HttpMappedStatus = "resp_done"
	Created      HttpMappedStatus = "create_done"
	Updated      HttpMappedStatus = "update_done"
	Validate     HttpMappedStatus = "validation_err"
	NotFound     HttpMappedStatus = "not_found"
	Failed       HttpMappedStatus = "resp_fail"
	Unauthorized HttpMappedStatus = "unauthorized"
	Conflict     HttpMappedStatus = "conflict"
	ItemExist    HttpMappedStatus = "item_exist"
)
