package response

type MonthlyUsageItem struct {
	Month string `json:"month"`
	Usage int64  `json:"usage"`
}

type GrowthItem struct {
	Period string `json:"period"`
	Users  int    `json:"users"`
}

type AdminReports struct {
	TotalUsers         int                `json:"total_users"`
	TotalCandidates    int                `json:"total_candidates"`
	TotalRecruiters    int                `json:"total_recruiters"`
	TotalCompanies     int                `json:"total_companies"`
	TotalInterviews    int                `json:"total_interviews"`
	TotalTokenUsage    int64              `json:"total_token_usage"`
	InputTokens        int64              `json:"input_tokens"`
	OutputTokens       int64              `json:"output_tokens"`
	SystemUptimeHours  float64            `json:"system_uptime_hours"`
	CpuUsagePercent    float64            `json:"cpu_usage_percent"`
	RamUsagePercent    float64            `json:"ram_usage_percent"`
	RedisMemoryMb      float64            `json:"redis_memory_mb"`
	ServerLatencyMs    int                `json:"server_latency_ms"`
	TokenUsageByModel  map[string]int64   `json:"token_usage_by_model"`
	MonthlyTokenUsage  []MonthlyUsageItem `json:"monthly_token_usage"`
	UserGrowthTrend    []GrowthItem       `json:"user_growth_trend"`
}
