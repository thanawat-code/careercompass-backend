package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanawat-code/careercompass-backend/internal/database"
	"github.com/thanawat-code/careercompass-backend/internal/models"
)

type CareerHandler struct {
	DB *database.DB
}

func NewCareerHandler(db *database.DB) *CareerHandler {
	return &CareerHandler{DB: db}
}

func (h *CareerHandler) RecommendCareer(c *gin.Context) {
	ctx := c.Request.Context()
	var payload models.UserPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. ดึงอาชีพทั้งหมด
	allCareers, err := h.getAllCareers(ctx)
	if err != nil {
		log.Printf("Error fetching careers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch careers"})
		return
	}

	// 2. เรียก AI (ระบบ Smart Retry - อัปเดตรายชื่อโมเดลใหม่)
	recommendedIDs := h.callGeminiAI_SmartRetry(payload, allCareers)

	// 3. บันทึกผลลัพธ์
	userID := 1
	h.saveResult(ctx, userID, payload.MBTI, recommendedIDs)

	// 4. ดึงข้อมูลรายละเอียด
	results, err := h.getCareersByIDs(ctx, recommendedIDs)
	if err != nil {
		log.Printf("Error fetching results: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":             true,
		"recommended_careers": results,
	})
}

// --- AI Function (The Brain 🧠 - Smart Retry Version) ---

func (h *CareerHandler) callGeminiAI_SmartRetry(user models.UserPayload, careers []models.Career) []int {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("⚠️ Warning: GEMINI_API_KEY not found")
		return []int{1, 2, 3}
	}

	// เตรียมข้อมูลอาชีพ
	var careerListStr []string
	for _, c := range careers {
		careerListStr = append(careerListStr, fmt.Sprintf("ID: %d, Name: %s", c.ID, c.Title))
	}

	promptText := fmt.Sprintf(`
		Role: Career Counselor.
		User Profile:
		- MBTI: %s
		- Aptitude Test Result: %v
		- Interests/Skills: %s

		Available Careers Database:
		%s

		Task: Analyze the user profile and select the top 3 most suitable career IDs.
		Response Format: return ONLY a JSON object with a single key "career_ids" containing an array of integers.
		Example: {"career_ids": [1, 5, 8]}
	`, user.MBTI, user.Aptitude, user.Knowledge, strings.Join(careerListStr, "\n"))

	// 🟢 รายชื่อโมเดลที่อัปเดตตามลิสต์ของคุณ (Gemini 2.5 / 2.0 / Latest)
	candidateModels := []string{
		"gemini-2.5-flash",          // ใหม่ล่าสุด
		"gemini-flash-latest",       // ตัวล่าสุด (Alias)
		"gemini-2.0-flash-lite-001", // ตัวเล็ก ประหยัดโควต้า
		"gemini-2.0-flash",          // ตัวมาตรฐาน
		"gemini-2.5-pro",            // ตัวเก่งสุด
		"gemini-pro-latest",         // ตัวเก่งล่าสุด (Alias)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	// Loop ลองทีละโมเดล
	for _, modelName := range candidateModels {
		log.Printf("🔄 Trying model: %s ...", modelName)

		// ใช้ Endpoint v1beta เพราะรุ่น 2.5/2.0 มักจะอยู่ใน Beta
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", modelName, apiKey)

		requestBody, _ := json.Marshal(map[string]interface{}{
			"contents": []interface{}{
				map[string]interface{}{
					"parts": []interface{}{
						map[string]string{"text": promptText},
					},
				},
			},
			"generationConfig": map[string]interface{}{
				"response_mime_type": "application/json",
			},
		})

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("❌ Connection error with %s: %v", modelName, err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
			// ✅ เจอตัวที่ใช่แล้ว!
			log.Printf("✅ Success with model: %s", modelName)

			// แกะผลลัพธ์
			type GeminiResponse struct {
				Candidates []struct {
					Content struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"content"`
				} `json:"candidates"`
			}
			var geminiResp GeminiResponse
			if err := json.Unmarshal(body, &geminiResp); err == nil {
				if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
					jsonText := geminiResp.Candidates[0].Content.Parts[0].Text

					// Clean & Parse JSON
					jsonText = strings.TrimSpace(jsonText)
					if start := strings.Index(jsonText, "{"); start != -1 {
						jsonText = jsonText[start:]
					}
					if end := strings.LastIndex(jsonText, "}"); end != -1 {
						jsonText = jsonText[:end+1]
					}

					type FinalResult struct {
						CareerIDs []int `json:"career_ids"`
					}
					var finalRes FinalResult
					if err := json.Unmarshal([]byte(jsonText), &finalRes); err == nil && len(finalRes.CareerIDs) > 0 {
						log.Printf("🤖 AI Recommended IDs: %v", finalRes.CareerIDs)
						return finalRes.CareerIDs
					}
				}
			}
			break
		} else {
			// ถ้า Error 429 (Quota) หรือ 404 ให้ลองตัวถัดไป
			log.Printf("⚠️ Model %s failed with status: %d", modelName, resp.StatusCode)
		}
	}

	log.Println("❌ All models failed. Using fallback.")
	return []int{1, 2, 3} // Fallback
}

// --- Helper Functions (Database) ---
// ส่วนนี้เหมือนเดิม

func (h *CareerHandler) getAllCareers(ctx context.Context) ([]models.Career, error) {
	rows, err := h.DB.Pool.Query(ctx, "SELECT id, title FROM careers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var careers []models.Career
	for rows.Next() {
		var c models.Career
		if err := rows.Scan(&c.ID, &c.Title); err != nil {
			return nil, err
		}
		careers = append(careers, c)
	}
	return careers, nil
}

func (h *CareerHandler) saveResult(ctx context.Context, userID int, mbti string, careerIDs []int) {
	if len(careerIDs) == 0 {
		return
	}
	query := `INSERT INTO user_assessment_results (user_id, mbti_result, suggested_career_ids) VALUES ($1, $2, $3)`
	_, err := h.DB.Pool.Exec(ctx, query, userID, mbti, careerIDs)
	if err != nil {
		log.Printf("Error saving result to DB: %v", err)
	}
}

func (h *CareerHandler) getCareersByIDs(ctx context.Context, ids []int) ([]models.Career, error) {
	if len(ids) == 0 {
		return []models.Career{}, nil
	}
	query := `SELECT id, title, short_description, icon_key, slug FROM careers WHERE id = ANY($1)`
	rows, err := h.DB.Pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Career
	for rows.Next() {
		var c models.Career
		err := rows.Scan(&c.ID, &c.Title, &c.ShortDescription, &c.IconKey, &c.Slug)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, nil
}
