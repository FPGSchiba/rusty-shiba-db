package documents

import "rsdb/src/util"

type documentCreateRequest struct {
	// json tag to de-serialize json body
	Collection string                 `uri:"collection" binding:"required"`
	Data       map[string]interface{} `json:"data" binding:"required"`
}

type documentCreateResponse struct {
	util.Response
	DocumentId string `json:"document_id"`
}

type documentReadRequest struct {
}
