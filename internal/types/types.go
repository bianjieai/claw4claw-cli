package types

type Reputation struct {
	Level string  `json:"level"`
	Score float64 `json:"score"`
}

type AgentStatus struct {
	DID        string     `json:"did"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	Balance    string     `json:"balance"`
	Reputation Reputation `json:"reputation"`
}

type MarketAgent struct {
	ID           uint        `json:"id"`
	Name         string      `json:"name"`
	Category     string      `json:"category"`
	Price        string      `json:"price"`
	Rating       float64     `json:"rating"`
	Status       string      `json:"status"`
	DID          string      `json:"did"`
	Description  string      `json:"description,omitempty"`
	Staked       string      `json:"staked,omitempty"`
	Uptime       string      `json:"uptime,omitempty"`
	Capabilities []string    `json:"capabilities,omitempty"`
	Reputation   *Reputation `json:"reputation,omitempty"`
}

type MarketAgentList struct {
	Data  []MarketAgent `json:"data"`
	Total int64         `json:"total"`
}

type MarketTask struct {
	ID             int     `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Bounty         string  `json:"bounty"`
	BountyValue    float64 `json:"bountyValue"`
	Status         string  `json:"status"`
	PostedBy       string  `json:"postedBy"`
	PostedByID     int     `json:"postedById"`
	PostedByAvatar string  `json:"postedByAvatar,omitempty"`
	Image          string  `json:"image,omitempty"`
	Category       string  `json:"category"`
	Deadline       *string `json:"deadline,omitempty"`
	CreatedAt      string  `json:"createdAt"`
}

type MarketTaskList struct {
	Data       []MarketTask `json:"data"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	TotalPages int          `json:"totalPages"`
	Categories []string     `json:"categories,omitempty"`
}

type ConsoleTask struct {
	ID                  string  `json:"id"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Status              string  `json:"status"`
	Role                string  `json:"role"`
	Bounty              float64 `json:"bounty"`
	AssociatedAgentID   *string `json:"associatedAgentId,omitempty"`
	AssociatedAgentName *string `json:"associatedAgentName,omitempty"`
	CreatedAt           string  `json:"createdAt"`
	Deadline            *string `json:"deadline,omitempty"`
}

type ConsoleTaskList struct {
	Data  []ConsoleTask `json:"data"`
	Total int           `json:"total"`
}

type PublishTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Bounty      float64 `json:"bounty"`
	Category    string  `json:"category"`
	Deadline    *string `json:"deadline,omitempty"`
}

type PublishTaskResponse struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Status       string  `json:"status"`
	Bounty       float64 `json:"bounty"`
	StakedAmount float64 `json:"stakedAmount"`
	StakeStatus  string  `json:"stakeStatus"`
	CreatedAt    string  `json:"createdAt"`
}

type ApplyTaskRequest struct {
	TaskID            uint   `json:"taskId"`
	Message           string `json:"message,omitempty"`
	EstimatedDuration string `json:"estimatedDuration,omitempty"`
}

type ApplyTaskResponse struct {
	ApplicationID uint   `json:"applicationId"`
	TaskID        uint   `json:"taskId"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
}

type SubmitTaskRequest struct {
	Content     string   `json:"content"`
	Attachments []string `json:"attachments,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

type SubmitTaskResponse struct {
	SubmissionID string `json:"submissionId"`
	TaskID       int    `json:"taskId"`
	Status       string `json:"status"`
	SubmittedAt  string `json:"submittedAt"`
}

type AcceptTaskRequest struct {
	Rating int    `json:"rating"`
	Review string `json:"review,omitempty"`
}

type AcceptTaskResponse struct {
	TaskID       int     `json:"taskId"`
	Status       string  `json:"status"`
	CompletedAt  string  `json:"completedAt"`
	TotalPayment float64 `json:"totalPayment"`
}

type MarketService struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	Provider        string  `json:"provider"`
	ProviderID      int     `json:"providerId"`
	Category        string  `json:"category"`
	Price           string  `json:"price"`
	PriceValue      float64 `json:"priceValue"`
	Rating          float64 `json:"rating"`
	Completed       int     `json:"completed"`
	Image           string  `json:"image"`
	Description     string  `json:"description"`
	InputSchema     string  `json:"inputSchema,omitempty"`
	OutputSchema    string  `json:"outputSchema,omitempty"`
	AvgResponseMs   int     `json:"avgResponseMs"`
	AvgResponseTime string  `json:"avgResponseTime"`
	CreatedAt       string  `json:"createdAt"`
}

type MarketServiceList struct {
	Data       []MarketService `json:"data"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"totalPages"`
	Categories []string        `json:"categories,omitempty"`
}

type ConsoleService struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	ProviderAgentID string `json:"providerAgentId"`
	Category        string `json:"category"`
	Price           int    `json:"price"`
	TotalCalls      int    `json:"totalCalls"`
	TotalEarnings   int    `json:"totalEarnings"`
	Status          string `json:"status"`
	CreatedAt       string `json:"createdAt"`
}

type ConsoleServiceList struct {
	Data  []ConsoleService `json:"data"`
	Total int              `json:"total"`
}

type ConsoleServiceDetails struct {
	ID                string              `json:"id"`
	Title             string              `json:"title"`
	ProviderAgentID   string              `json:"providerAgentId"`
	ProviderAgentName string              `json:"providerAgentName"`
	Price             int                 `json:"price"`
	TotalCalls        int                 `json:"totalCalls"`
	TotalEarnings     int                 `json:"totalEarnings"`
	AvgResponseTime   string              `json:"avgResponseTime"`
	Status            string              `json:"status"`
	CreatedAt         string              `json:"createdAt"`
	ErrorRate         string              `json:"errorRate"`
	RecentLogs        []ConsoleServiceLog `json:"recentLogs"`
}

type ConsoleServiceLog struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Duration  string `json:"duration"`
	Cost      int    `json:"cost"`
}

type PublishServiceRequest struct {
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Category      string                 `json:"category"`
	Price         float64                `json:"price"`
	AvgResponseMs int                    `json:"avgResponseMs"`
	InputSchema   map[string]interface{} `json:"inputSchema"`
	OutputSchema  map[string]interface{} `json:"outputSchema"`
}

type PublishServiceResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type MyServiceItem struct {
	ID             int     `json:"id"`
	Title          string  `json:"title"`
	Status         string  `json:"status"`
	Price          float64 `json:"price"`
	CompletedCount int     `json:"completedCount"`
	Rating         float64 `json:"rating"`
	CreatedAt      string  `json:"createdAt"`
}

type MyServicesListResponse struct {
	Services []MyServiceItem `json:"services"`
	Total    int             `json:"total"`
}

type UpdateServiceRequest struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
}

type UpdateServiceResponse struct {
	Message string `json:"message"`
}

type ActivateServiceResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type DeactivateServiceResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type InvokeServiceRequest struct {
	ServiceID      int                    `json:"serviceId"`
	Input          map[string]interface{} `json:"input"`
	MaxPrice       *float64               `json:"maxPrice,omitempty"`
	IdempotencyKey string                 `json:"idempotencyKey,omitempty"`
}

type InvokeServiceResponse struct {
	InvocationID int     `json:"invocationId"`
	ServiceID    int     `json:"serviceId"`
	Status       string  `json:"status"`
	Price        float64 `json:"price"`
	CreatedAt    string  `json:"createdAt"`
}

type SubmitServiceResultRequest struct {
	Status      string                 `json:"status"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Attachments []string               `json:"attachments,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

type SubmitServiceResultResponse struct {
	InvocationID int    `json:"invocationId"`
	Status       string `json:"status"`
	CompletedAt  string `json:"completedAt"`
}

type ReviewServiceInvocationRequest struct {
	Rating *int   `json:"rating,omitempty"`
	Review string `json:"review,omitempty"`
}

type ReviewServiceInvocationResponse struct {
	InvocationID int    `json:"invocationId"`
	Message      string `json:"message"`
}

type ServiceInvocationItem struct {
	ID           int     `json:"id"`
	ServiceID    int     `json:"serviceId"`
	ServiceTitle string  `json:"serviceTitle"`
	Role         string  `json:"role"`
	Status       string  `json:"status"`
	Price        float64 `json:"price"`
	CreatedAt    string  `json:"createdAt"`
	CompletedAt  *string `json:"completedAt,omitempty"`
}

type ServiceInvocationListResponse struct {
	Invocations []ServiceInvocationItem `json:"invocations"`
	Total       int                     `json:"total"`
	Page        int                     `json:"page"`
	Limit       int                     `json:"limit"`
}

type ServiceInvocationDetail struct {
	ID           int                    `json:"id"`
	ServiceID    int                    `json:"serviceId"`
	ServiceTitle string                 `json:"serviceTitle"`
	Role         string                 `json:"role"`
	Status       string                 `json:"status"`
	Input        map[string]interface{} `json:"input,omitempty"`
	Output       map[string]interface{} `json:"output,omitempty"`
	Attachments  []string               `json:"attachments,omitempty"`
	Price        float64                `json:"price"`
	Rating       *int                   `json:"rating,omitempty"`
	Review       string                 `json:"review,omitempty"`
	TimeoutAt    *string                `json:"timeoutAt,omitempty"`
	CreatedAt    string                 `json:"createdAt"`
	CompletedAt  *string                `json:"completedAt,omitempty"`
}

type TaskApplication struct {
	ID                      string  `json:"id"`
	TaskID                  string  `json:"taskId"`
	AgentID                 string  `json:"agentId"`
	AgentName               string  `json:"agentName"`
	AgentRating             float64 `json:"agentRating"`
	Message                 string  `json:"message"`
	EstimatedCompletionTime *string `json:"estimatedCompletionTime,omitempty"`
	Status                  string  `json:"status"`
	CreatedAt               string  `json:"createdAt"`
}

type TaskApplicationList struct {
	Data  []TaskApplication `json:"data"`
	Total int               `json:"total"`
}

type AcceptApplicantRequest struct {
	Message string `json:"message,omitempty"`
}

type AcceptApplicantResponse struct {
	TaskID          int    `json:"taskId"`
	SelectedAgentID string `json:"selectedAgentId"`
	Status          string `json:"status"`
	StartedAt       string `json:"startedAt"`
}

type CancelTaskResponse struct {
	TaskID      int    `json:"taskId"`
	Status      string `json:"status"`
	CancelledAt string `json:"cancelledAt"`
}

type RegisterAgentReq struct {
	Name         string                 `json:"name"`
	Category     string                 `json:"category"`
	Description  string                 `json:"description"`
	Capabilities []string               `json:"capabilities"`
	Preferences  map[string]interface{} `json:"preferences"`
}

type AgentInfo struct {
	ID                uint                   `json:"id"`
	Name              string                 `json:"name"`
	Category          string                 `json:"category"`
	Status            string                 `json:"status"`
	Description       string                 `json:"description"`
	Capabilities      []string               `json:"capabilities"`
	Preferences       map[string]interface{} `json:"preferences"`
	Staked            float64                `json:"staked"`
	Earned            float64                `json:"earned"`
	Rating            float64                `json:"rating"`
	CompletedTasks    int                    `json:"completedTasks"`
	CreatedAt         string                 `json:"createdAt"`
	LastActiveAt      string                 `json:"lastActiveAt"`
	MarketVisibility  string                 `json:"marketVisibility"`
	MarketPublishedAt *string                `json:"marketPublishedAt,omitempty"`
}

type UpdateAgentReq struct {
	Name         *string                `json:"name,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Capabilities []string               `json:"capabilities,omitempty"`
	Preferences  map[string]interface{} `json:"preferences,omitempty"`
}

type SetStatusReq struct {
	Status string `json:"status"`
}

type RegisterAgentResp struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type UpdateAgentResp struct {
	Message string `json:"message"`
}

type SetStatusResp struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type UnpublishAgentResponse struct {
	ID               uint   `json:"id"`
	MarketVisibility string `json:"marketVisibility"`
	Message          string `json:"message"`
}

type MyTask struct {
	ID                  string  `json:"id"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Status              string  `json:"status"`
	Role                string  `json:"role"`
	Bounty              float64 `json:"bounty"`
	StakedAmount        float64 `json:"stakedAmount"`
	StakeStatus         string  `json:"stakeStatus"`
	Deadline            *string `json:"deadline,omitempty"`
	AssociatedAgentID   *string `json:"associatedAgentId,omitempty"`
	AssociatedAgentName *string `json:"associatedAgentName,omitempty"`
	CreatedAt           string  `json:"createdAt"`
}

type MyTaskList struct {
	Tasks []MyTask `json:"tasks"`
	Total int      `json:"total"`
}

type MyTaskApplication struct {
	ID                string  `json:"id"`
	TaskID            string  `json:"taskId"`
	TaskTitle         string  `json:"taskTitle"`
	TaskBounty        float64 `json:"taskBounty"`
	TaskStatus        string  `json:"taskStatus"`
	Status            string  `json:"status"`
	Message           string  `json:"message"`
	EstimatedDuration string  `json:"estimatedDuration,omitempty"`
	CreatedAt         string  `json:"createdAt"`
}

type MyTaskApplicationList struct {
	Applications []MyTaskApplication `json:"applications"`
	Total        int                 `json:"total"`
	Page         int                 `json:"page"`
	Limit        int                 `json:"limit"`
}

type AcceptedTask struct {
	ID                 string  `json:"id"`
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	Status             string  `json:"status"`
	Bounty             float64 `json:"bounty"`
	Deadline           *string `json:"deadline,omitempty"`
	PublisherAgentID   string  `json:"publisherAgentId"`
	PublisherAgentName string  `json:"publisherAgentName"`
	ApplicationID      string  `json:"applicationId"`
	StartedAt          string  `json:"startedAt"`
}

type AcceptedTaskList struct {
	Tasks []AcceptedTask `json:"tasks"`
	Total int            `json:"total"`
}

type TaskDetail struct {
	ID               uint    `json:"id"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Category         string  `json:"category"`
	Status           string  `json:"status"`
	Bounty           float64 `json:"bounty"`
	Deadline         *string `json:"deadline,omitempty"`
	PublisherAgentID uint    `json:"publisherAgentId"`
	WorkerAgentID    *uint   `json:"workerAgentId,omitempty"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

type WithdrawApplicationResponse struct {
	Message string `json:"message"`
}

type PublishAgentReq struct {
	Preferences *PublishPreferences `json:"preferences,omitempty"`
}

type PublishPreferences struct {
	ExpectedSalary int      `json:"expectedSalary,omitempty"`
	WorkHours      string   `json:"workHours,omitempty"`
	PreferredTasks []string `json:"preferredTasks,omitempty"`
}

type PublishAgentResp struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	MarketVisibility  string `json:"marketVisibility"`
	MarketPublishedAt string `json:"marketPublishedAt"`
	Message           string `json:"message"`
}

type EmploymentStatus string

const (
	EmploymentStatusPending    EmploymentStatus = "pending"
	EmploymentStatusActive     EmploymentStatus = "active"
	EmploymentStatusTerminated EmploymentStatus = "terminated"
	EmploymentStatusCompleted  EmploymentStatus = "completed"
	EmploymentStatusRejected   EmploymentStatus = "rejected"
)

type StakeStatus string

const (
	StakeStatusFrozen  StakeStatus = "frozen"
	StakeStatusSettled StakeStatus = "settled"
)

type CreateEmploymentRequest struct {
	EmployeeAgentID uint    `json:"employeeAgentId"`
	Salary          float64 `json:"salary"`
	Duration        string  `json:"duration,omitempty"`
	StakeAmount     float64 `json:"stakeAmount,omitempty"`
}

type CreateEmploymentResponse struct {
	ID              uint        `json:"id"`
	EmployerAgentID uint        `json:"employerAgentId"`
	EmployeeAgentID uint        `json:"employeeAgentId"`
	Salary          float64     `json:"salary"`
	StakedAmount    float64     `json:"stakedAmount"`
	StakeStatus     StakeStatus `json:"stakeStatus"`
	Status          string      `json:"status"`
	CreatedAt       string      `json:"createdAt"`
}

type AcceptEmploymentRequest struct {
	Message string `json:"message,omitempty"`
}

type AcceptEmploymentResponse struct {
	ID        uint   `json:"id"`
	Status    string `json:"status"`
	StartTime string `json:"startTime"`
	Message   string `json:"message"`
}

type RejectEmploymentRequest struct {
	Reason string `json:"reason,omitempty"`
}

type RejectEmploymentResponse struct {
	ID      uint   `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type TerminateEmploymentRequest struct {
	Reason string `json:"reason,omitempty"`
}

type TerminateEmploymentResponse struct {
	ID            uint    `json:"id"`
	Status        string  `json:"status"`
	TotalDuration int64   `json:"totalDuration"`
	BilledHours   int     `json:"billedHours"`
	TotalPayment  float64 `json:"totalPayment"`
	RefundAmount  float64 `json:"refundAmount"`
	EndTime       string  `json:"endTime"`
	Message       string  `json:"message"`
}

type MyEmploymentsQueryParams struct {
	Role   string
	Status string
	Page   int
	Limit  int
}

type EmploymentListItem struct {
	ID                uint        `json:"id"`
	Role              string      `json:"role"`
	EmployerAgentID   uint        `json:"employerAgentId"`
	EmployerAgentName string      `json:"employerAgentName"`
	EmployeeAgentID   uint        `json:"employeeAgentId"`
	EmployeeAgentName string      `json:"employeeAgentName"`
	Salary            float64     `json:"salary"`
	StakedAmount      float64     `json:"stakedAmount"`
	StakeStatus       StakeStatus `json:"stakeStatus"`
	Status            string      `json:"status"`
	Duration          string      `json:"duration,omitempty"`
	TotalDuration     int64       `json:"totalDuration"`
	StartTime         *string     `json:"startTime,omitempty"`
	CreatedAt         string      `json:"createdAt"`
}

type MyEmploymentsListResponse struct {
	Employments []EmploymentListItem `json:"employments"`
	Total       int                  `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}

type EmploymentDetail struct {
	ID                  uint        `json:"id"`
	Role                string      `json:"role"`
	EmployerAgentID     uint        `json:"employerAgentId"`
	EmployerAgentName   string      `json:"employerAgentName"`
	EmployerAgentAvatar string      `json:"employerAgentAvatar,omitempty"`
	EmployeeAgentID     uint        `json:"employeeAgentId"`
	EmployeeAgentName   string      `json:"employeeAgentName"`
	EmployeeAgentAvatar string      `json:"employeeAgentAvatar,omitempty"`
	Salary              float64     `json:"salary"`
	StakedAmount        float64     `json:"stakedAmount"`
	StakeStatus         StakeStatus `json:"stakeStatus"`
	Status              string      `json:"status"`
	Duration            string      `json:"duration,omitempty"`
	TotalDuration       int64       `json:"totalDuration"`
	StartTime           *string     `json:"startTime,omitempty"`
	EndTime             *string     `json:"endTime,omitempty"`
	CreatedAt           string      `json:"createdAt"`
}

type WalletTransaction struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Balance     int64  `json:"balance"`
	Description string `json:"description"`
	RelatedID   string `json:"relatedId,omitempty"`
	RelatedType string `json:"relatedType,omitempty"`
	CreatedAt   string `json:"createdAt"`
}

type WalletSummary struct {
	Balance            int64               `json:"balance"`
	FrozenAmount       int64               `json:"frozenAmount"`
	TotalAssets        int64               `json:"totalAssets"`
	BalanceInCNY       float64             `json:"balanceInCNY"`
	InvoicableAmount   float64             `json:"invoicableAmount"`
	PendingRefunds     int                 `json:"pendingRefunds"`
	RecentTransactions []WalletTransaction `json:"recentTransactions"`
}

type BudgetInfo struct {
	BudgetType      string  `json:"budgetType"`
	BudgetAmount    *int64  `json:"budgetAmount,omitempty"`
	BudgetUsed      int64   `json:"budgetUsed"`
	BudgetRemaining *int64  `json:"budgetRemaining,omitempty"`
	BudgetPeriod    string  `json:"budgetPeriod,omitempty"`
	BudgetResetAt   *string `json:"budgetResetAt,omitempty"`
}
