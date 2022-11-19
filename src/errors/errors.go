package errors

type ReturnCode string

const (
	ReturnCodeStateChange                       = "Display state change"
	ReturnCodeDisplayDoesntExistNewDisplayAdded = "Display doesn't exist, new display added"
	ReturnCodeFunctionToDisplayNotFound         = "Function to display not found in function table"

	ErrorCodeAddDisplayNode = "Error: AddDisplayNode"
)
