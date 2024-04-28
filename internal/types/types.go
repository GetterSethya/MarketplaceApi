package types

type ServerResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AppError struct {
	Error  error
	Status int
}

type ListQueryTransactionValid struct {
	Seller bool
	Limit  int
	Offset int
	Sort   string
	Order  string
	Search string
}

type ListQueryTransaction struct {
	Seller string
	Limit  string
	Offset string
	Sort   string
	Order  string
	Search string
}

type ListQuery struct {
	UserOnly       string
	Limit          string
	Offset         string
	Tags           []string
	Condition      string
	ShowEmptyStock string
	MaxPrice       string
	MinPrice       string
	Sort           string
	Order          string
	Search         string
}

type ListQueryValid struct {
	UserOnly       string
	Limit          int
	Offset         int
	Tags           []string
	Condition      string
	ShowEmptyStock string
	MaxPrice       float64
	MinPrice       float64
	Sort           string
	Order          string
	Search         string
}
