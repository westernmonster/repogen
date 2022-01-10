package repo

import (
	"context"
	"database/sql"
	types "github.com/jmoiron/sqlx/types"
	sqalx "github.com/westernmonster/sqalx"
)

type Doctor struct {
	ID                                             int64         `db:"id" json:"id,string"`                                                                                              // ID
	UserID                                         int64         `db:"user_id" json:"user_id,string"`                                                                                    // UserID
	Name                                           string        `db:"name" json:"name"`                                                                                                 // Name
	Title                                          string        `db:"title" json:"title"`                                                                                               // Title
	Description                                    string        `db:"description" json:"description"`                                                                                   // Description
	City                                           string        `db:"city" json:"city"`                                                                                                 // City
	State                                          string        `db:"state" json:"state"`                                                                                               // State
	CountryCode                                    string        `db:"country_code" json:"country_code"`                                                                                 // CountryCode
	Zipcode                                        string        `db:"zipcode" json:"zipcode"`                                                                                           // Zipcode
	SupportPediatric                               types.BitBool `db:"support_pediatric" json:"support_pediatric"`                                                                       // SupportPediatric
	ExternalAppointmentTypeID                      int32         `db:"external_appointment_type_id" json:"external_appointment_type_id"`                                                 // ExternalAppointmentTypeID Acuity appointment type id
	ExternalCalendarID                             int32         `db:"external_calendar_id" json:"external_calendar_id"`                                                                 // ExternalCalendarID Acuity appointment calendar type id
	ExternalCalendarLink                           string        `db:"external_calendar_link" json:"external_calendar_link"`                                                             // ExternalCalendarLink Acuity appointment calendar link
	OnlineConsultationURL                          string        `db:"online_consultation_url" json:"online_consultation_url"`                                                           // OnlineConsultationURL Online consultation link
	ExternalFollowUpAppointmentTypeID              int32         `db:"external_follow_up_appointment_type_id" json:"external_follow_up_appointment_type_id"`                             // ExternalFollowUpAppointmentTypeID Acuity appointment follow up type id
	ExternalPediatricConsultationCategoryOneTypeID int32         `db:"external_pediatric_consultation_category_one_type_id" json:"external_pediatric_consultation_category_one_type_id"` // ExternalPediatricConsultationCategoryOneTypeID
	ExternalPediatricConsultationCategoryTwoTypeID int32         `db:"external_pediatric_consultation_category_two_type_id" json:"external_pediatric_consultation_category_two_type_id"` // ExternalPediatricConsultationCategoryTwoTypeID
	ExternalPediatricFollowUpTypeID                int32         `db:"external_pediatric_follow_up_type_id" json:"external_pediatric_follow_up_type_id"`                                 // ExternalPediatricFollowUpTypeID
	InPersonAppointmentTypeID                      int32         `db:"in_person_appointment_type_id" json:"in_person_appointment_type_id"`                                               // InPersonAppointmentTypeID
	AddressLine1                                   string        `db:"address_line1" json:"address_line1"`                                                                               // AddressLine1
	NavigationInstruction                          string        `db:"navigation_instruction" json:"navigation_instruction"`                                                             // NavigationInstruction
	Lat                                            float64       `db:"lat" json:"lat"`                                                                                                   // Lat
	Lng                                            float64       `db:"lng" json:"lng"`                                                                                                   // Lng
	SupportStates                                  string        `db:"support_states" json:"support_states"`                                                                             // SupportStates
	TaskURL                                        string        `db:"task_url" json:"task_url"`                                                                                         // TaskURL
	ZoomPersonalMeetingID                          string        `db:"zoom_personal_meeting_id" json:"zoom_personal_meeting_id"`                                                         // ZoomPersonalMeetingID
	CreatedAt                                      int64         `db:"created_at" json:"created_at"`                                                                                     // CreatedAt
	UpdatedAt                                      int64         `db:"updated_at" json:"updated_at"`                                                                                     // UpdatedAt
	AddressLine2                                   string        `db:"address_line2" json:"address_line2"`                                                                               // AddressLine2
}

type DDoctorRepository struct{}

func (p *Dao) GetAllDoctors(c context.Context, node sqalx.Node) (items []*model.Doctor, err error) {
	items = make([]*model.Doctor, 0)
	sqlSelect := "SELECT a.id,a.user_id,a.name,a.title,a.description,a.city,a.state,a.country_code,a.zipcode,a.support_pediatric,a.external_appointment_type_id,a.external_calendar_id,a.external_calendar_link,a.online_consultation_url,a.external_follow_up_appointment_type_id,a.external_pediatric_consultation_category_one_type_id,a.external_pediatric_consultation_category_two_type_id,a.external_pediatric_follow_up_type_id,a.in_person_appointment_type_id,a.address_line1,a.navigation_instruction,a.lat,a.lng,a.support_states,a.task_url,a.zoom_personal_meeting_id,a.created_at,a.updated_at,a.address_line2 FROM d_doctor a WHERE 1=1 ORDER BY a.created_at DESC "

	if err = node.SelectContext(c, &items, sqlSelect); err != nil {
		log.For(c).Errorf("dao.GetAllDoctors err(%+v)", err)
		return
	}
	return
}

func (p *Dao) GetDoctorsPaged(c context.Context, node sqalx.Node, cond map[string]interface{}, limit int, offset int) (total int, items []*model.Doctor, err error) {
	items = make([]*model.Doctor, 0)

	sqlCount := "SELECT COUNT(1) as count FROM d_doctor a WHERE 1=1 %s"
	sqlSelect := "SELECT a.id,a.user_id,a.name,a.title,a.description,a.city,a.state,a.country_code,a.zipcode,a.support_pediatric,a.external_appointment_type_id,a.external_calendar_id,a.external_calendar_link,a.online_consultation_url,a.external_follow_up_appointment_type_id,a.external_pediatric_consultation_category_one_type_id,a.external_pediatric_consultation_category_two_type_id,a.external_pediatric_follow_up_type_id,a.in_person_appointment_type_id,a.address_line1,a.navigation_instruction,a.lat,a.lng,a.support_states,a.task_url,a.zoom_personal_meeting_id,a.created_at,a.updated_at,a.address_line2 FROM d_doctor a WHERE 1=1 ORDER BY %s a.id DESC LIMIT ?,?"

	condition := make([]interface{}, 0)
	clause := ""

	if val, ok := cond["created_at[gt]"]; ok {
		clause += " AND a.created_at > ?"
		condition = append(condition, val)
	}

	if val, ok := cond["created_at[lt]"]; ok {
		clause += " AND a.created_at < ?"
		condition = append(condition, val)
	}

	if val, ok := cond["created_at[gte]"]; ok {
		clause += " AND a.created_at >= ?"
		condition = append(condition, val)
	}

	if val, ok := cond["created_at[lte]"]; ok {
		clause += " AND a.created_at <= ?"
		condition = append(condition, val)
	}

	sqlCount = fmt.Sprintf(sqlCount, clause)
	if err = node.GetContext(c, &total, sqlCount, condition...); err != nil {
		log.For(c).Errorf("dao.GetDoctorsPaged err(%+v) condition(%+v)", err, cond)
		return
	}

	sqlSelect = fmt.Sprintf(sqlSelect, clause)
	condition = append(condition, offset, limit)

	if err = node.SelectContext(c, &items, sqlSelect, condition...); err != nil {
		log.For(c).Errorf("dao.GetDoctorsPaged err(%+v) condition(%+v)", err, cond)
		return
	}
	return
}

func (p *Dao) GetDoctorByID(c context.Context, node sqalx.Node, id int64) (item *model.Doctor, err error) {
	item = new(model.Doctor)
	sqlSelect := "SELECT a.id,a.user_id,a.name,a.title,a.description,a.city,a.state,a.country_code,a.zipcode,a.support_pediatric,a.external_appointment_type_id,a.external_calendar_id,a.external_calendar_link,a.online_consultation_url,a.external_follow_up_appointment_type_id,a.external_pediatric_consultation_category_one_type_id,a.external_pediatric_consultation_category_two_type_id,a.external_pediatric_follow_up_type_id,a.in_person_appointment_type_id,a.address_line1,a.navigation_instruction,a.lat,a.lng,a.support_states,a.task_url,a.zoom_personal_meeting_id,a.created_at,a.updated_at,a.address_line2 FROM d_doctor a WHERE a.id=?"

	if err = node.GetContext(c, item, sqlSelect, id); err != nil {
		if err == sql.ErrNoRows {
			item = nil
			err = nil
			return
		}
		log.For(c).Errorf("dao.GetDoctorByID err(%+v), id(%+v)", err, id)
	}

	return
}

func (p *Dao) AddDoctor(c context.Context, node sqalx.Node, item *model.Doctor) (err error) {
	sqlInsert := "INSERT INTO d_doctor( id,user_id,name,title,description,city,state,country_code,zipcode,support_pediatric,external_appointment_type_id,external_calendar_id,external_calendar_link,online_consultation_url,external_follow_up_appointment_type_id,external_pediatric_consultation_category_one_type_id,external_pediatric_consultation_category_two_type_id,external_pediatric_follow_up_type_id,in_person_appointment_type_id,address_line1,navigation_instruction,lat,lng,support_states,task_url,zoom_personal_meeting_id,created_at,updated_at,address_line2) VALUES ( ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	if _, err = node.ExecContext(c, sqlInsert, item.ID, item.UserID, item.Name, item.Title, item.Description, item.City, item.State, item.CountryCode, item.Zipcode, item.SupportPediatric, item.ExternalAppointmentTypeID, item.ExternalCalendarID, item.ExternalCalendarLink, item.OnlineConsultationURL, item.ExternalFollowUpAppointmentTypeID, item.ExternalPediatricConsultationCategoryOneTypeID, item.ExternalPediatricConsultationCategoryTwoTypeID, item.ExternalPediatricFollowUpTypeID, item.InPersonAppointmentTypeID, item.AddressLine1, item.NavigationInstruction, item.Lat, item.Lng, item.SupportStates, item.TaskURL, item.ZoomPersonalMeetingID, item.CreatedAt, item.UpdatedAt, item.AddressLine2); err != nil {
		log.For(c).Errorf("dao.AddDoctor err(%+v), item(%+v)", err, item)
		return
	}

	return
}

func (p *Dao) UpdateDoctor(c context.Context, node sqalx.Node, item *model.Doctor) (err error) {
	sqlUpdate := "UPDATE d_doctor SET user_id=?,name=?,title=?,description=?,city=?,state=?,country_code=?,zipcode=?,support_pediatric=?,external_appointment_type_id=?,external_calendar_id=?,external_calendar_link=?,online_consultation_url=?,external_follow_up_appointment_type_id=?,external_pediatric_consultation_category_one_type_id=?,external_pediatric_consultation_category_two_type_id=?,external_pediatric_follow_up_type_id=?,in_person_appointment_type_id=?,address_line1=?,navigation_instruction=?,lat=?,lng=?,support_states=?,task_url=?,zoom_personal_meeting_id=?,updated_at=?,address_line2=? WHERE id=?"

	_, err = node.ExecContext(c, sqlUpdate, item.UserID, item.Name, item.Title, item.Description, item.City, item.State, item.CountryCode, item.Zipcode, item.SupportPediatric, item.ExternalAppointmentTypeID, item.ExternalCalendarID, item.ExternalCalendarLink, item.OnlineConsultationURL, item.ExternalFollowUpAppointmentTypeID, item.ExternalPediatricConsultationCategoryOneTypeID, item.ExternalPediatricConsultationCategoryTwoTypeID, item.ExternalPediatricFollowUpTypeID, item.InPersonAppointmentTypeID, item.AddressLine1, item.NavigationInstruction, item.Lat, item.Lng, item.SupportStates, item.TaskURL, item.ZoomPersonalMeetingID, item.UpdatedAt, item.AddressLine2, item.ID)
	if err != nil {
		log.For(c).Errorf("dao.UpdateDoctor err(%+v), item(%+v)", err, item)
		return
	}

	return
}

func (p *Dao) DelDoctor(c context.Context, node sqalx.Node, id int64) (err error) {
	sqlDelete := "DELETE FROM d_doctor WHERE id=? "

	if _, err = node.ExecContext(c, sqlDelete, id); err != nil {
		log.For(c).Errorf("dao.DelDoctor err(%+v), item(%+v)", err, id)
		return
	}

	return
}
