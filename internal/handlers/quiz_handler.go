package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// QuizQuestion represents a single quiz question with 4 options and 1 answer
type QuizQuestion struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}

// GenerateQuiz generates 10 quiz questions using Gemini AI based on career and stage
// POST /api/quiz/generate
func GenerateQuiz(c *gin.Context) {
	var req struct {
		CareerName string `json:"career_name"`
		StageName  string `json:"stage_name"`
		CareerSlug string `json:"career_slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questions, err := callGeminiForQuiz(req.CareerName, req.StageName, req.CareerSlug)
	if err != nil {
		log.Printf("❌ Gemini failed, using fallback questions: %v", err)
		questions = getFallbackQuestions(req.CareerName, req.StageName, req.CareerSlug)
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

// getFallbackQuestions returns topic-appropriate static questions when AI is unavailable
func getFallbackQuestions(careerName, stageName, careerSlug string) []QuizQuestion {
	slug := strings.ToLower(careerSlug + " " + careerName + " " + stageName)

	var pool []QuizQuestion
	if strings.Contains(slug, "data") || strings.Contains(slug, "ข้อมูล") {
		pool = dataQuestions
	} else if strings.Contains(slug, "machine") || strings.Contains(slug, "ml") || strings.Contains(slug, "ai") {
		pool = mlQuestions
	} else if strings.Contains(slug, "web") || strings.Contains(slug, "frontend") || strings.Contains(slug, "backend") {
		pool = webQuestions
	} else if strings.Contains(slug, "foundation") || strings.Contains(slug, "พื้นฐาน") || strings.Contains(slug, "basic") {
		pool = foundationQuestions
	} else {
		pool = generalQuestions
	}

	// Shuffle and return 10
	shuffled := make([]QuizQuestion, len(pool))
	copy(shuffled, pool)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := int(time.Now().UnixNano()) % (i + 1)
		if j < 0 {
			j = -j
		}
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	if len(shuffled) > 10 {
		shuffled = shuffled[:10]
	}
	return shuffled
}

// ─── Fallback Question Banks ──────────────────────────────────────────────────

var foundationQuestions = []QuizQuestion{
	{"อัลกอริทึม (Algorithm) คืออะไร?", []string{"ชุดคำสั่งแก้ปัญหาแบบทีละขั้นตอน", "ภาษาโปรแกรม", "ฮาร์ดแวร์", "ระบบปฏิบัติการ"}, "ชุดคำสั่งแก้ปัญหาแบบทีละขั้นตอน"},
	{"ตัวแปร (Variable) คืออะไร?", []string{"พื้นที่จัดเก็บข้อมูลในหน่วยความจำ", "ฟังก์ชัน", "คลาส", "Loop"}, "พื้นที่จัดเก็บข้อมูลในหน่วยความจำ"},
	{"Loop (การวนซ้ำ) ใช้สำหรับอะไร?", []string{"ทำซ้ำชุดคำสั่ง", "ประกาศตัวแปร", "Import library", "สร้าง Database"}, "ทำซ้ำชุดคำสั่ง"},
	{"Function คืออะไร?", []string{"บล็อกโค้ดที่สามารถเรียกใช้ซ้ำได้", "ตัวแปร", "ลูป", "เงื่อนไข"}, "บล็อกโค้ดที่สามารถเรียกใช้ซ้ำได้"},
	{"Conditional Statement (if-else) ทำหน้าที่อะไร?", []string{"ตัดสินใจตามเงื่อนไข", "วนซ้ำ", "เก็บข้อมูล", "Import โมดูล"}, "ตัดสินใจตามเงื่อนไข"},
	{"Array คืออะไร?", []string{"โครงสร้างข้อมูลที่เก็บค่าหลายค่าในลำดับ", "ชนิดข้อมูลตัวเลข", "ฟังก์ชัน", "ตัวแปรเดียว"}, "โครงสร้างข้อมูลที่เก็บค่าหลายค่าในลำดับ"},
	{"String คืออะไร?", []string{"ชนิดข้อมูลที่เก็บข้อความ", "ตัวเลขทศนิยม", "ค่า boolean", "อาเรย์"}, "ชนิดข้อมูลที่เก็บข้อความ"},
	{"Boolean มีค่าได้กี่ค่า?", []string{"2 ค่า (True/False)", "3 ค่า", "ไม่จำกัด", "1 ค่า"}, "2 ค่า (True/False)"},
	{"Debugging คืออะไร?", []string{"กระบวนการค้นหาและแก้ไข Bug", "การเขียนโค้ดใหม่", "การ Deploy", "การ Test"}, "กระบวนการค้นหาและแก้ไข Bug"},
	{"Comment ในโค้ดคืออะไร?", []string{"ข้อความอธิบายที่ไม่ถูก Execute", "ตัวแปร", "ฟังก์ชัน", "Loop"}, "ข้อความอธิบายที่ไม่ถูก Execute"},
	{"IDE ย่อมาจากอะไร?", []string{"Integrated Development Environment", "Internet Data Exchange", "Internal Design Engine", "Indexed Data Element"}, "Integrated Development Environment"},
	{"Git ใช้สำหรับอะไร?", []string{"ระบบควบคุมเวอร์ชันของโค้ด", "Database", "Web Server", "Testing Framework"}, "ระบบควบคุมเวอร์ชันของโค้ด"},
}

var dataQuestions = []QuizQuestion{
	{"Pandas ใน Python ใช้สำหรับอะไร?", []string{"จัดการข้อมูลเชิงตาราง", "สร้าง UI", "เขียน API", "เชื่อมต่อ Database"}, "จัดการข้อมูลเชิงตาราง"},
	{"Mean (ค่าเฉลี่ย) คำนวณอย่างไร?", []string{"ผลรวม ÷ จำนวน", "ค่าที่เกิดบ่อยสุด", "ค่ากลาง", "ค่าสูงสุด-ต่ำสุด"}, "ผลรวม ÷ จำนวน"},
	{"Outlier คืออะไร?", []string{"ค่าที่อยู่ห่างจากกลุ่มข้อมูล", "ค่าเฉลี่ย", "ค่ามัธยฐาน", "ค่าฐานนิยม"}, "ค่าที่อยู่ห่างจากกลุ่มข้อมูล"},
	{"Median คืออะไร?", []string{"ค่ากลางของข้อมูลเมื่อเรียงลำดับ", "ค่าเฉลี่ย", "ค่าสูงสุด", "ผลรวม"}, "ค่ากลางของข้อมูลเมื่อเรียงลำดับ"},
	{"DataFrame ใน Pandas คืออะไร?", []string{"ตารางข้อมูล 2 มิติ (มี row และ column)", "รายการ 1 มิติ", "Dictionary", "String"}, "ตารางข้อมูล 2 มิติ (มี row และ column)"},
	{"Data Cleaning คืออะไร?", []string{"การแก้ไขข้อมูลผิดพลาด/ขาดหาย", "การวิเคราะห์ข้อมูล", "การแสดงกราฟ", "การจัดเก็บข้อมูล"}, "การแก้ไขข้อมูลผิดพลาด/ขาดหาย"},
	{"Correlation คืออะไร?", []string{"ความสัมพันธ์ระหว่างตัวแปร", "ค่าเฉลี่ย", "ความแปรปรวน", "ค่ามัธยฐาน"}, "ความสัมพันธ์ระหว่างตัวแปร"},
	{"Matplotlib ใช้สำหรับอะไร?", []string{"สร้างกราฟและการแสดงข้อมูล", "Machine Learning", "Web Scraping", "Database"}, "สร้างกราฟและการแสดงข้อมูล"},
	{"Missing Value จัดการด้วยวิธีใดได้บ้าง?", []string{"ลบออก, แทนด้วยค่าเฉลี่ย, หรือ Imputation", "คัดลอก", "Transpose", "Sort"}, "ลบออก, แทนด้วยค่าเฉลี่ย, หรือ Imputation"},
	{"SQL SELECT ใช้สำหรับอะไร?", []string{"ดึงข้อมูลจากตาราง", "เพิ่มข้อมูล", "ลบข้อมูล", "อัปเดตข้อมูล"}, "ดึงข้อมูลจากตาราง"},
	{"Standard Deviation บอกอะไร?", []string{"การกระจายของข้อมูลรอบค่าเฉลี่ย", "ค่าเฉลี่ย", "ค่าสูงสุด", "จำนวนข้อมูล"}, "การกระจายของข้อมูลรอบค่าเฉลี่ย"},
	{"ETL ย่อมาจากอะไร?", []string{"Extract, Transform, Load", "Edit, Test, Launch", "Evaluate, Track, Log", "Encode, Train, Label"}, "Extract, Transform, Load"},
}

var mlQuestions = []QuizQuestion{
	{"Supervised Learning คืออะไร?", []string{"เรียนรู้จากข้อมูลที่มี label", "เรียนรู้โดยไม่มี label", "เรียนรู้จากรางวัล", "เรียนรู้แบบลึก"}, "เรียนรู้จากข้อมูลที่มี label"},
	{"Overfitting คืออะไร?", []string{"โมเดลจำ train data มากเกินไป", "โมเดลที่ดี", "ข้อมูลไม่พอ", "Learning rate สูง"}, "โมเดลจำ train data มากเกินไป"},
	{"Feature ใน ML คืออะไร?", []string{"ตัวแปร input ของโมเดล", "ผลลัพธ์ของโมเดล", "Activation function", "จำนวน layer"}, "ตัวแปร input ของโมเดล"},
	{"Training Set และ Test Set แตกต่างกันอย่างไร?", []string{"Train ใช้สอนโมเดล, Test ใช้ประเมินผล", "Test ใช้สอนโมเดล, Train ใช้ประเมินผล", "เหมือนกัน", "Test คือ Validation"}, "Train ใช้สอนโมเดล, Test ใช้ประเมินผล"},
	{"Accuracy คืออะไร?", []string{"สัดส่วนการทำนายถูกต้อง", "ความเร็วของโมเดล", "ขนาดของโมเดล", "จำนวน feature"}, "สัดส่วนการทำนายถูกต้อง"},
	{"Decision Tree คืออะไร?", []string{"โมเดลที่แบ่งข้อมูลด้วยเงื่อนไขแบบ tree", "Neural Network", "Regression", "Clustering"}, "โมเดลที่แบ่งข้อมูลด้วยเงื่อนไขแบบ tree"},
	{"Gradient Descent ทำงานอย่างไร?", []string{"ปรับ parameter เพื่อลด Loss function", "เพิ่ม Loss function", "สร้าง feature ใหม่", "ลบ outlier"}, "ปรับ parameter เพื่อลด Loss function"},
	{"Unsupervised Learning ต่างจาก Supervised อย่างไร?", []string{"ไม่มี label", "มี label ครบ", "ใช้ reinforcement", "ใช้ deep learning"}, "ไม่มี label"},
	{"Cross-validation ใช้สำหรับอะไร?", []string{"ประเมินโมเดลโดยใช้ข้อมูลหลายส่วน", "เพิ่มข้อมูล", "ลด feature", "เพิ่ม layer"}, "ประเมินโมเดลโดยใช้ข้อมูลหลายส่วน"},
	{"Scikit-learn คืออะไร?", []string{"Python library สำหรับ Machine Learning", "Database", "Web Framework", "Visualization tool"}, "Python library สำหรับ Machine Learning"},
	{"Neural Network ได้รับแรงบันดาลใจจากอะไร?", []string{"สมองมนุษย์", "DNA", "อินเทอร์เน็ต", "ตาราง Excel"}, "สมองมนุษย์"},
	{"Epoch ในการ train ML คืออะไร?", []string{"การผ่านข้อมูล train ครบ 1 รอบ", "จำนวน layer", "Learning rate", "Batch size"}, "การผ่านข้อมูล train ครบ 1 รอบ"},
}

var webQuestions = []QuizQuestion{
	{"HTML ย่อมาจากอะไร?", []string{"HyperText Markup Language", "High Transfer Markup Language", "Hyperlink Text Machine Language", "HyperText Machine Language"}, "HyperText Markup Language"},
	{"CSS ใช้สำหรับอะไร?", []string{"จัดรูปแบบและตกแต่งหน้าเว็บ", "สร้าง Database", "เขียน Logic", "จัดการ Server"}, "จัดรูปแบบและตกแต่งหน้าเว็บ"},
	{"REST API คืออะไร?", []string{"รูปแบบสถาปัตยกรรม API ที่ใช้ HTTP", "Database", "CSS Framework", "Testing tool"}, "รูปแบบสถาปัตยกรรม API ที่ใช้ HTTP"},
	{"HTTP GET ต่างจาก POST อย่างไร?", []string{"GET ดึงข้อมูล, POST ส่งข้อมูล", "GET ส่งข้อมูล, POST ดึงข้อมูล", "เหมือนกัน", "GET ลบข้อมูล"}, "GET ดึงข้อมูล, POST ส่งข้อมูล"},
	{"JSON ย่อมาจากอะไร?", []string{"JavaScript Object Notation", "Java Standard Object Name", "JavaScript Online Notation", "Java Syntax Object Network"}, "JavaScript Object Notation"},
	{"DOM คืออะไร?", []string{"โครงสร้างเอกสาร HTML ที่ browser สร้างขึ้น", "Database", "ภาษา styling", "Server"}, "โครงสร้างเอกสาร HTML ที่ browser สร้างขึ้น"},
	{"React คืออะไร?", []string{"JavaScript library สำหรับสร้าง UI", "Database", "CSS Framework", "Backend framework"}, "JavaScript library สำหรับสร้าง UI"},
	{"HTTPS ต่างจาก HTTP อย่างไร?", []string{"HTTPS มีการเข้ารหัส (SSL/TLS)", "HTTP เร็วกว่า", "เหมือนกัน", "HTTPS ใช้ port 80"}, "HTTPS มีการเข้ารหัส (SSL/TLS)"},
	{"Responsive Design คืออะไร?", []string{"การออกแบบที่ปรับตามขนาดหน้าจอ", "การออกแบบสำหรับ Desktop เท่านั้น", "การออกแบบสำหรับมือถือเท่านั้น", "การออกแบบ 3D"}, "การออกแบบที่ปรับตามขนาดหน้าจอ"},
	{"npm ใช้สำหรับอะไร?", []string{"จัดการ JavaScript packages", "สร้าง Database", "เขียน CSS", "Deploy เว็บ"}, "จัดการ JavaScript packages"},
	{"SQL injection คืออะไร?", []string{"การโจมตีโดยแทรก SQL ใน input", "การ query ปกติ", "การ backup", "การ optimize"}, "การโจมตีโดยแทรก SQL ใน input"},
	{"Cookie ใช้สำหรับอะไร?", []string{"เก็บข้อมูลขนาดเล็กบน browser ของผู้ใช้", "เก็บข้อมูลบน server", "ส่ง email", "เชื่อมต่อ API"}, "เก็บข้อมูลขนาดเล็กบน browser ของผู้ใช้"},
}

var generalQuestions = []QuizQuestion{
	{"Python คืออะไร?", []string{"ภาษาโปรแกรมมิ่ง", "ฐานข้อมูล", "ระบบปฏิบัติการ", "เว็บเบราว์เซอร์"}, "ภาษาโปรแกรมมิ่ง"},
	{"Machine Learning คืออะไร?", []string{"การเรียนรู้ของเครื่องจากข้อมูล", "ฮาร์ดแวร์", "ระบบปฏิบัติการ", "โปรแกรมคำนวณ"}, "การเรียนรู้ของเครื่องจากข้อมูล"},
	{"VS Code คืออะไร?", []string{"Code editor", "ระบบปฏิบัติการ", "แอนตี้ไวรัส", "ฐานข้อมูล"}, "Code editor"},
	{"Cloud Computing คืออะไร?", []string{"บริการคอมพิวเตอร์ผ่านอินเทอร์เน็ต", "คอมพิวเตอร์ธรรมดา", "ระบบ LAN", "Offline software"}, "บริการคอมพิวเตอร์ผ่านอินเทอร์เน็ต"},
	{"API คืออะไร?", []string{"ตัวกลางเชื่อมต่อระหว่างซอฟต์แวร์", "ภาษาโปรแกรม", "Database", "ระบบปฏิบัติการ"}, "ตัวกลางเชื่อมต่อระหว่างซอฟต์แวร์"},
	{"Open Source คืออะไร?", []string{"ซอฟต์แวร์ที่เปิดให้ดูและแก้ไข source code", "ซอฟต์แวร์ฟรี", "ซอฟต์แวร์เชิงพาณิชย์", "ซอฟต์แวร์ทดลองใช้"}, "ซอฟต์แวร์ที่เปิดให้ดูและแก้ไข source code"},
	{"Linux คืออะไร?", []string{"ระบบปฏิบัติการ Open Source", "ภาษาโปรแกรม", "Database", "Browser"}, "ระบบปฏิบัติการ Open Source"},
	{"Encryption คืออะไร?", []string{"การเข้ารหัสข้อมูล", "การลบข้อมูล", "การ backup", "การ compress"}, "การเข้ารหัสข้อมูล"},
	{"Agile คืออะไร?", []string{"รูปแบบการพัฒนาซอฟต์แวร์แบบยืดหยุ่น", "ภาษาโปรแกรม", "Framework CSS", "Database"}, "รูปแบบการพัฒนาซอฟต์แวร์แบบยืดหยุ่น"},
	{"DevOps คืออะไร?", []string{"แนวปฏิบัติรวม Dev และ Ops เข้าด้วยกัน", "ภาษาโปรแกรม", "Hardware", "Network protocol"}, "แนวปฏิบัติรวม Dev และ Ops เข้าด้วยกัน"},
	{"Cybersecurity คืออะไร?", []string{"การปกป้องระบบและข้อมูลจากการโจมตี", "การพัฒนาเว็บ", "Machine Learning", "การออกแบบ UI"}, "การปกป้องระบบและข้อมูลจากการโจมตี"},
	{"UX Design คืออะไร?", []string{"การออกแบบประสบการณ์ผู้ใช้", "การเขียนโค้ด", "การจัดการ Database", "การ Deploy"}, "การออกแบบประสบการณ์ผู้ใช้"},
}

func callGeminiForQuiz(careerName, stageName, careerSlug string) ([]QuizQuestion, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	promptText := fmt.Sprintf(`คุณเป็นผู้เชี่ยวชาญด้าน %s โดยเฉพาะในเรื่อง "%s"

สร้างข้อสอบแบบปรนัย 10 ข้อ เกี่ยวกับหัวข้อ "%s" สำหรับสายงาน "%s"

กฎสำคัญ:
1. แต่ละข้อต้องมี 4 ตัวเลือก (options)
2. ต้องมีคำตอบถูกเพียง 1 ข้อ และต้องตรงกับหนึ่งใน options พอดี
3. คำถามต้องหลากหลาย ครอบคลุมทั้งทฤษฎีและปฏิบัติ
4. ภาษาไทยเป็นหลัก แต่คำศัพท์เทคนิคภาษาอังกฤษใช้ได้
5. ความยากง่ายต้องระดับเริ่มต้นถึงกลาง

ตอบกลับเป็น JSON เท่านั้น ในรูปแบบนี้:
{
  "questions": [
    {
      "question": "คำถามที่ 1?",
      "options": ["ตัวเลือก A", "ตัวเลือก B", "ตัวเลือก C", "ตัวเลือก D"],
      "answer": "ตัวเลือก A"
    }
  ]
}`, careerName, stageName, stageName, careerName)

	type modelEntry struct {
		name       string
		version    string
		maxRetries int
	}

	// gemini-2.0 uses v1beta, gemini-1.5 uses v1
	models := []modelEntry{
		{"gemini-2.0-flash", "v1beta", 3},
		{"gemini-2.0-flash-lite", "v1beta", 2},
		{"gemini-1.5-flash", "v1", 2},
		{"gemini-1.5-flash-8b", "v1", 2},
		{"gemini-1.5-pro", "v1", 1},
	}

	client := &http.Client{Timeout: 45 * time.Second}

	for _, m := range models {
		for attempt := 1; attempt <= m.maxRetries; attempt++ {
			log.Printf("🔄 [Quiz] Trying model: %s (attempt %d/%d) ...", m.name, attempt, m.maxRetries)

			apiURL := fmt.Sprintf(
				"https://generativelanguage.googleapis.com/%s/models/%s:generateContent?key=%s",
				m.version, m.name, apiKey,
			)

			// response_mime_type is only supported in v1beta
			genConfig := map[string]interface{}{
				"temperature": 0.7,
			}
			if m.version == "v1beta" {
				genConfig["response_mime_type"] = "application/json"
			}

			requestBody, _ := json.Marshal(map[string]interface{}{
				"contents": []interface{}{
					map[string]interface{}{
						"parts": []interface{}{
							map[string]string{"text": promptText},
						},
					},
				},
				"generationConfig": genConfig,
			})

			httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
			if err != nil {
				log.Printf("❌ [Quiz] Failed to create request: %v", err)
				break
			}
			httpReq.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(httpReq)
			if err != nil {
				log.Printf("❌ [Quiz] Connection error with %s: %v", m.name, err)
				break
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			// Rate limited — wait with backoff and retry same model
			if resp.StatusCode == 429 {
				waitSec := time.Duration(attempt*5) * time.Second
				log.Printf("⏳ [Quiz] Model %s rate limited (attempt %d), waiting %v before retry...", m.name, attempt, waitSec)
				time.Sleep(waitSec)
				continue
			}

			// Non-200 (not rate limit) → skip to next model
			if resp.StatusCode != 200 {
				bodySnap := body
				if len(bodySnap) > 300 {
					bodySnap = bodySnap[:300]
				}
				log.Printf("⚠️ [Quiz] Model %s failed with status: %d — %s", m.name, resp.StatusCode, string(bodySnap))
				break
			}

			// Parse Gemini response envelope
			type GeminiResp struct {
				Candidates []struct {
					Content struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"content"`
				} `json:"candidates"`
			}
			var geminiResp GeminiResp
			if err := json.Unmarshal(body, &geminiResp); err != nil {
				log.Printf("⚠️ [Quiz] Failed to parse Gemini response: %v", err)
				break
			}

			if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
				log.Printf("⚠️ [Quiz] Empty candidates from model %s", m.name)
				break
			}

			// Extract and clean JSON text
			jsonText := geminiResp.Candidates[0].Content.Parts[0].Text
			jsonText = strings.TrimSpace(jsonText)
			jsonText = strings.TrimPrefix(jsonText, "```json")
			jsonText = strings.TrimPrefix(jsonText, "```")
			jsonText = strings.TrimSuffix(jsonText, "```")
			jsonText = strings.TrimSpace(jsonText)

			if start := strings.Index(jsonText, "{"); start != -1 {
				jsonText = jsonText[start:]
			}
			if end := strings.LastIndex(jsonText, "}"); end != -1 {
				jsonText = jsonText[:end+1]
			}

			type QuizResponse struct {
				Questions []QuizQuestion `json:"questions"`
			}
			var quizResp QuizResponse
			if err := json.Unmarshal([]byte(jsonText), &quizResp); err != nil {
				log.Printf("⚠️ [Quiz] Failed to parse quiz JSON: %v\nRaw snippet: %.200s", err, jsonText)
				break
			}

			if len(quizResp.Questions) < 5 {
				log.Printf("⚠️ [Quiz] Only %d questions returned, skipping model", len(quizResp.Questions))
				break
			}

			// Trim to exactly 10
			if len(quizResp.Questions) > 10 {
				quizResp.Questions = quizResp.Questions[:10]
			}

			log.Printf("✅ [Quiz] Generated %d questions with model %s", len(quizResp.Questions), m.name)
			return quizResp.Questions, nil
		}

		log.Printf("⚠️ [Quiz] Exhausted retries for model %s, moving to next", m.name)
	}

	return nil, fmt.Errorf("all models failed to generate quiz")
}
