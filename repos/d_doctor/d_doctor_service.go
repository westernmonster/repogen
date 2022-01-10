package service

import "context"

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
	UpdatedAt                                      int64   `json:"updated_at"`
	AddressLine2                                   string  `json:"address_line2"`
}

type DoctorListResp struct {
	Items []*DoctorItem `json:"items"`
	Total int           `json:"total"`
}

func (p *Service) GetAllDoctors(c context.Context) (items []*model.DoctorItem, err error) {
	var data []*model.Doctor
	if data, err = p.d.GetAllDoctors(c, p.d.DB()); err != nil {
		return
	}

	items = make([]*model.DoctorItem, len(data))
	for i, v := range data {
		items[i] = &model.DoctorItem{
			ID:                                v.ID,
			UserID:                            v.UserID,
			Name:                              v.Name,
			Title:                             v.Title,
			Description:                       v.Description,
			City:                              v.City,
			State:                             v.State,
			CountryCode:                       v.CountryCode,
			Zipcode:                           v.Zipcode,
			SupportPediatric:                  bool(v.SupportPediatric),
			ExternalAppointmentTypeID:         v.ExternalAppointmentTypeID,
			ExternalCalendarID:                v.ExternalCalendarID,
			ExternalCalendarLink:              v.ExternalCalendarLink,
			OnlineConsultationURL:             v.OnlineConsultationURL,
			ExternalFollowUpAppointmentTypeID: v.ExternalFollowUpAppointmentTypeID,
			ExternalPediatricConsultationCategoryOneTypeID: v.ExternalPediatricConsultationCategoryOneTypeID,
			ExternalPediatricConsultationCategoryTwoTypeID: v.ExternalPediatricConsultationCategoryTwoTypeID,
			ExternalPediatricFollowUpTypeID:                v.ExternalPediatricFollowUpTypeID,
			InPersonAppointmentTypeID:                      v.InPersonAppointmentTypeID,
			AddressLine1:                                   v.AddressLine1,
			NavigationInstruction:                          v.NavigationInstruction,
			Lat:                                            v.Lat,
			Lng:                                            v.Lng,
			SupportStates:                                  v.SupportStates,
			TaskURL:                                        v.TaskURL,
			ZoomPersonalMeetingID:                          v.ZoomPersonalMeetingID,
			CreatedAt:                                      v.CreatedAt,
			UpdatedAt:                                      v.UpdatedAt,
			AddressLine2:                                   v.AddressLine2,
		}
	}
	return
}

func (p *Service) GetDoctorsListPaged(c context.Context, cond map[string]interface{}, page, pageSize int) (resp *model.DoctorListResp, err error) {
	offset := (page - 1) * pageSize

	var data []*model.Doctor
	var total int
	if total, data, err = p.d.GetDoctorsPaged(c, p.d.DB(), cond, pageSize, offset); err != nil {
		return
	}

	resp = &model.DoctorListResp{
		Items: make([]*model.DoctorItem, len(data)),
		Total: total,
	}

	for i, v := range data {
		resp.Items[i] = &model.DoctorItem{
			ID:                                v.ID,
			UserID:                            v.UserID,
			Name:                              v.Name,
			Title:                             v.Title,
			Description:                       v.Description,
			City:                              v.City,
			State:                             v.State,
			CountryCode:                       v.CountryCode,
			Zipcode:                           v.Zipcode,
			SupportPediatric:                  bool(v.SupportPediatric),
			ExternalAppointmentTypeID:         v.ExternalAppointmentTypeID,
			ExternalCalendarID:                v.ExternalCalendarID,
			ExternalCalendarLink:              v.ExternalCalendarLink,
			OnlineConsultationURL:             v.OnlineConsultationURL,
			ExternalFollowUpAppointmentTypeID: v.ExternalFollowUpAppointmentTypeID,
			ExternalPediatricConsultationCategoryOneTypeID: v.ExternalPediatricConsultationCategoryOneTypeID,
			ExternalPediatricConsultationCategoryTwoTypeID: v.ExternalPediatricConsultationCategoryTwoTypeID,
			ExternalPediatricFollowUpTypeID:                v.ExternalPediatricFollowUpTypeID,
			InPersonAppointmentTypeID:                      v.InPersonAppointmentTypeID,
			AddressLine1:                                   v.AddressLine1,
			NavigationInstruction:                          v.NavigationInstruction,
			Lat:                                            v.Lat,
			Lng:                                            v.Lng,
			SupportStates:                                  v.SupportStates,
			TaskURL:                                        v.TaskURL,
			ZoomPersonalMeetingID:                          v.ZoomPersonalMeetingID,
			CreatedAt:                                      v.CreatedAt,
			UpdatedAt:                                      v.UpdatedAt,
			AddressLine2:                                   v.AddressLine2,
		}
	}
	return
}

func (p *Service) GetDoctorByID(c context.Context, id int64) (item *model.Doctor, err error) {
	return p.getDoctorByID(c, id)
}

func (p *Service) getDoctorByID(c context.Context, id int64) (item *model.Doctor, err error) {
	if item, err = p.d.GetDoctorByID(c, p.d.DB(), id); err != nil {
		return
	} else if item == nil {
		err = ecode.NothingFound
		return
	}

	return
}
