package utils

type PaginationResult struct {
	Page       int
	PageSize   int
	Conditions map[string]interface{}
}

// GetOffsetFields 从请求体中提取分页参数和条件字段。
// 该函数接收一个映射，该映射包含分页信息（如页码和页大小）以及其他条件字段。
// 它会解析出分页参数和条件字段，并返回一个PaginationResult对象，该对象包含这些信息。
// 如果解析过程中没有遇到错误，它将返回nil作为错误值。
func GetOffsetFields(body map[string]interface{}) (*PaginationResult, error) {
	// 初始化页码和页大小变量。
	var (
		page     = 0
		pageSize = 0
	)

	// 尝试从body中提取页码信息。
	if p, ok := body["page"].(int); ok {
		page = p
	} else if pFloat, ok := body["page"].(float64); ok {
		page = int(pFloat)
	}

	// 尝试从body中提取页大小信息。
	if ps, ok := body["page_size"].(int); ok {
		pageSize = ps
	} else if psFloat, ok := body["page_size"].(float64); ok {
		pageSize = int(psFloat)
	}

	// 创建一个映射来存储除分页信息外的条件字段。
	conditions := make(map[string]interface{})
	for key, value := range body {
		// 忽略分页相关的字段，因为它们已经被处理过了。
		if key == "page" || key == "page_size" {
			continue
		}
		conditions[key] = value
	}

	// 返回包含分页信息和条件字段的PaginationResult对象。
	return &PaginationResult{
		Page:       page,
		PageSize:   pageSize,
		Conditions: conditions,
	}, nil
}
