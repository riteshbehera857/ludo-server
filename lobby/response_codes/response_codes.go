package response_codes

type ResponseCodeDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var ResponseCodes = map[string]ResponseCodeDetails{
	"BOARD_LIST_FETCHED_SUCCESSFULLY": {
		Code:    "B200",
		Message: "Board list fetched successfully",
	},
}

// Helper function to get response detail
func GetResponseCodeDetails(key string) ResponseCodeDetails {
	if detail, exists := ResponseCodes[key]; exists {
		return detail
	}
	return ResponseCodeDetails{
		Code:    "S101",
		Message: "Something went wrong please try again",
	}
}
