package types

type Category string

const (
	CategoryWriting          Category = "writing"
	CategoryCustomerService  Category = "customer_service"
	CategoryDataAnalysis     Category = "data_analysis"
	CategoryMarketing        Category = "marketing"
	CategoryOfficeAutomation Category = "office_automation"
	CategoryProgramming      Category = "programming"
	CategoryDesign           Category = "design"
	CategoryConsulting       Category = "consulting"
	CategoryResearch         Category = "research"
)

var ValidCategories = []Category{
	CategoryWriting,
	CategoryCustomerService,
	CategoryDataAnalysis,
	CategoryMarketing,
	CategoryOfficeAutomation,
	CategoryProgramming,
	CategoryDesign,
	CategoryConsulting,
	CategoryResearch,
}

var CategoryLabels = map[Category]string{
	CategoryWriting:          "写作",
	CategoryCustomerService:  "客服",
	CategoryDataAnalysis:     "数据分析",
	CategoryMarketing:        "营销",
	CategoryOfficeAutomation: "办公自动化",
	CategoryProgramming:      "编程开发",
	CategoryDesign:           "设计",
	CategoryConsulting:       "咨询",
	CategoryResearch:         "研究",
}

func IsValidCategory(category string) bool {
	for _, c := range ValidCategories {
		if string(c) == category {
			return true
		}
	}
	return false
}

func GetAllCategories() []Category {
	return ValidCategories
}

func GetCategoryLabel(category Category) string {
	if label, ok := CategoryLabels[category]; ok {
		return label
	}
	return string(category)
}
