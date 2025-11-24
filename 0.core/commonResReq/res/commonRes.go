package response

import request "batchLog/0.core/commonResReq/req"

func GetPageResponse(req request.PageInfo, count, pages int64) map[string]interface{} {
	return map[string]interface{}{
		"page":req.Page,
		"size":req.Size,
		"total": count,
		"totalPages": pages,
	}
}