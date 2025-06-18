package utils

type PaginationResult struct {
	Page       int
	PageSize   int
	Conditions map[string]interface{}
}

func GetFields(body map[string]interface{}) (*PaginationResult, error) {
	var (
		page     = 0
		pageSize = 0
	)

	if p, ok := body["page"].(int); ok {
		page = p
	} else if pFloat, ok := body["page"].(float64); ok {
		page = int(pFloat)
	}

	if ps, ok := body["page_size"].(int); ok {
		pageSize = ps
	} else if psFloat, ok := body["page_size"].(float64); ok {
		pageSize = int(psFloat)
	}

	conditions := make(map[string]interface{})
	for key, value := range body {
		if key == "page" || key == "page_size" {
			continue
		}
		conditions[key] = value
	}

	return &PaginationResult{
		Page:       page,
		PageSize:   pageSize,
		Conditions: conditions,
	}, nil
}
