package models

type Task struct {
	MilestoneId        int64  `json:"milestone_id,omitempty"`
	EpicID             int64  `json:"epic_id,omitempty"`
	ActionPlanID       int64  `json:"action_plan_id,omitempty"`
	AssigneeId         string `json:"assignee_id,omitempty"`
	BrandId            string `json:"brand_id,omitempty"`
	BucketID           string `json:"bucket_id,omitempty"`
	CreatorID          string `json:"creator_id,omitempty"`
	ProjectID          string `json:"project_id,omitempty"`
	EndTime            string `json:"end_time,omitempty"`
	ID                 int64  `json:"id"`
	NumberOfAttachment int64  `json:"number_of_attachment,omitempty"`
	NumberOfComments   int64  `json:"number_of_comments,omitempty"`
	Priority           string `json:"priority,omitempty"`
	ResolvedTime       string `json:"resolved_time,omitempty"`
	StartTime          string `json:"start_time,omitempty"`
	Status             string `json:"status,omitempty"`
	StatusID           int64  `json:"status_id,omitempty"`
	Title              string `json:"title,omitempty"`
}

type TaskReq struct {
	MilestoneID int64  `json:"milestone_id"`
	EpicID      int64  `json:"epic_id"`
	GroupId     string `json:"group_id"`
	ProjectID   string `json:"project_id"`
	CompanyID   string `json:"company_id"`
	CreatorID   string `json:"creator_id"`
	ReporterID  string `json:"reporter_id"`
	AssigneeID  string `json:"assignee_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}
