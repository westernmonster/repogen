package service

import (
	jsoniter "jsoniter"
	"strings"
	vin "vin"
)

type DoctorItem struct {
	ID                                             int64   `json:"id,string"`
	UserID                                         int64   `json:"user_id,string"`
	Name                                           string  `json:"name"`
	Title                                          string  `json:"title"`
	Description                                    string  `json:"description"`
	City                                           string  `json:"city"`
	State                                          string  `json:"state"`
	CountryCode                                    string  `json:"country_code"`
	Zipcode                                        string  `json:"zipcode"`
	SupportPediatric                               bool    `json:"support_pediatric"`
	ExternalAppointmentTypeID                      int32   `json:"external_appointment_type_id"`
	ExternalCalendarID                             int32   `json:"external_calendar_id"`
	ExternalCalendarLink                           string  `json:"external_calendar_link"`
	OnlineConsultationURL                          string  `json:"online_consultation_url"`
	ExternalFollowUpAppointmentTypeID              int32   `json:"external_follow_up_appointment_type_id"`
	ExternalPediatricConsultationCategoryOneTypeID int32   `json:"external_pediatric_consultation_category_one_type_id"`
	ExternalPediatricConsultationCategoryTwoTypeID int32   `json:"external_pediatric_consultation_category_two_type_id"`
	ExternalPediatricFollowUpTypeID                int32   `json:"external_pediatric_follow_up_type_id"`
	InPersonAppointmentTypeID                      int32   `json:"in_person_appointment_type_id"`
	AddressLine1                                   string  `json:"address_line1"`
	NavigationInstruction                          string  `json:"navigation_instruction"`
	Lat                                            float64 `json:"lat"`
	Lng                                            float64 `json:"lng"`
	SupportStates                                  string  `json:"support_states"`
	TaskURL                                        string  `json:"task_url"`
	ZoomPersonalMeetingID                          string  `json:"zoom_personal_meeting_id"`
	CreatedAt                                      int64   `json:"created_at"`
	Decfield1                                      float64 `json:"decfield1"`
	Charfield                                      string  `json:"charfield"`
	UpdatedAt                                      int64   `json:"updated_at"`
	AddressLine2                                   string  `json:"address_line2"`
}

type DoctorListResp struct {
	Items []*DoctorItem `json:"items"`
	Total int           `json:"total"`
}

func getAllDoctors(c *vin.Context) {
	c.JSON(srv.GetAllDoctors(c))
}

func getDoctorsListPaged(c *vin.Context) {
	filters := c.Query("filters")

	cond := make(map[string]interface{})
	if strings.TrimSpace(filters) != "" {
		if err := jsoniter.Unmarshal([]byte(filters), &cond); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	page := c.QueryIntDefault("page", 1)
	pageSize := c.QueryIntDefault("page_size", 10)

	c.JSON(srv.GetDoctorsListPaged(c, cond, page, pageSize))
}

func getDoctorByID(c *vin.Context) {
	id, err := c.QueryInt64("id")
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.GetDoctorByID(c, id))
}
