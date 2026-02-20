package models

// Struct รับข้อมูลจาก Frontend
type UserPayload struct {
	MBTI      string            `json:"mbti"`
	Aptitude  map[string]string `json:"aptitude"`
	Knowledge string            `json:"knowledge"`
	UserID    int               `json:"user_id"`
}

// Struct สำหรับข้อมูลอาชีพใน Database
type Career struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	ShortDescription string `json:"description"`
	IconKey          string `json:"icon_key"`
	Slug             string `json:"slug"`
}
