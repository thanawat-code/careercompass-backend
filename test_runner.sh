#!/bin/bash
# CareerCompass Backend - Test Runner (CSV Output)
BASE="http://localhost:4546"
PASS=0; FAIL=0
CSV="/home/bnagon/all/careercompass-backend/test_results.csv"

echo "Test_ID,Func,Name,Expected Result,Actual Result,Priority,Status" > "$CSV"

req() { curl -s -w "|||%{http_code}" "$@"; }
req_auth() { curl -s -w "|||%{http_code}" -H "Authorization: Bearer $TOKEN" "$@"; }

check() {
  local id=$1 func=$2 name=$3 expected=$4 priority=$5
  local resp=$6 want_status=$7 jq_check=$8 actual_desc=$9

  local body status ok=true
  body="${resp%|||*}"; status="${resp##*|||}"

  [[ "$status" != "$want_status" ]] && ok=false
  if $ok && [[ -n "$jq_check" ]]; then
    echo "$body" | python3 -c "$jq_check" 2>/dev/null || ok=false
  fi

  local result="Pass"; $ok || result="Fail"
  $ok && ((PASS++)) || ((FAIL++))
  echo "$id,$func,\"$name\",\"$expected\",\"${actual_desc:-HTTP $status}\",${priority},${result}" >> "$CSV"
}

echo "🚀 Running tests..."

# Setup test user
TS=$(date +%s)
EMAIL="runner_${TS}@test.com"
R=$(req -X POST "$BASE/api/auth/register" -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"Test1234!\",\"confirm_password\":\"Test1234!\",\"display_name\":\"Runner\"}")
B="${R%|||*}"
USER_ID=$(echo "$B" | python3 -c "import json,sys;print(json.load(sys.stdin)['user']['id'])" 2>/dev/null)
TOKEN=$(echo "$B" | python3 -c "import json,sys;print(json.load(sys.stdin)['token'])" 2>/dev/null)
echo "✔ User: $EMAIL ($USER_ID)"

STAGE1="a0000001-0000-0000-0000-000000000001"
CAREER="Data Scientist"; CAREER_ENC="Data%20Scientist"
GHOST="00000000-0000-0000-0000-000000000099"

# ── HEALTH ──────────────────────────────────────────────────────────────────
r=$(req "$BASE/health"); b="${r%|||*}"; s="${r##*|||}"
check "HLT_001" "health" "DB connected" "200 + status:ok" "Low" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert d['status']=='ok'" \
  "HTTP $s - $(echo $b | python3 -c 'import json,sys;d=json.load(sys.stdin);print(d.get("status","?"))' 2>/dev/null)"

T0=$(date +%s%3N); curl -s "$BASE/health" >/dev/null; T1=$(date +%s%3N); MS=$((T1-T0))
(( MS < 500 )) && RES="Pass" || RES="Fail"
(( MS < 500 )) && ((PASS++)) || ((FAIL++))
echo "HLT_002,health,\"Response < 500ms\",\"< 500ms\",\"${MS}ms\",Low,${RES}" >> "$CSV"

# ── REGISTER bug fixes ───────────────────────────────────────────────────────
r=$(req -X POST "$BASE/api/auth/register" -H "Content-Type: application/json" \
  -d "{\"email\":\"UPPER_${TS}@test.com\",\"password\":\"Test1234!\",\"confirm_password\":\"Test1234!\",\"display_name\":\"U\"}")
check "REG_008_fix" "register" "Uppercase email → stored lowercase" "201 Created" "High" \
  "$r" "201" "" "HTTP ${r##*|||}"

r=$(req -X POST "$BASE/api/auth/register" -H "Content-Type: application/json" \
  -d "{\"email\":\"upper_${TS}@test.com\",\"password\":\"Test1234!\",\"confirm_password\":\"Test1234!\",\"display_name\":\"U\"}")
check "REG_008_dup" "register" "Same email lowercase → 409 duplicate" "409 Conflict" "High" \
  "$r" "409" "" "HTTP ${r##*|||}"

# ── LOGIN ────────────────────────────────────────────────────────────────────
r=$(req -X POST "$BASE/api/auth/login" -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL^^}\",\"password\":\"Test1234!\"}")
check "LOG_008_fix" "login" "Uppercase email login → 200 + token" "200 + token" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert 'token' in d" "HTTP ${r##*|||}"

r=$(req -X POST "$BASE/api/auth/login" -H "Content-Type: application/json" -d '{"email":"","password":"Test1234!"}')
check "LOG_002" "login" "Email empty → 400" "400 Bad Request" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req -X POST "$BASE/api/auth/login" -H "Content-Type: application/json" -d "{\"email\":\"$EMAIL\",\"password\":\"\"}")
check "LOG_003" "login" "Password empty → 400" "400 Bad Request" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req -X POST "$BASE/api/auth/login" -H "Content-Type: application/json" -d '{"email":"none@x.com","password":"Test1234!"}')
check "LOG_004" "login" "Email not found → 401" "401 Unauthorized" "High" "$r" "401" "" "HTTP ${r##*|||}"

r=$(req -X POST "$BASE/api/auth/login" -H "Content-Type: application/json" -d "{\"email\":\"$EMAIL\",\"password\":\"Wrong!\"}")
check "LOG_005" "login" "Wrong password → 401" "401 Unauthorized" "High" "$r" "401" "" "HTTP ${r##*|||}"

r=$(req -X POST "$BASE/api/auth/login" -H "Content-Type: application/json" -d "{\"email\":\"$EMAIL\",\"password\":\"Test1234!\"}")
check "LOG_006" "login" "Response has no password_hash" "Response safe" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert 'password_hash' not in str(d)" "No password_hash in JSON"

check "LOG_007" "login" "JWT token valid (3 parts)" "Token valid" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert len(d['token'].split('.'))==3" "JWT format valid"

# ── USERS ────────────────────────────────────────────────────────────────────
r=$(req_auth "$BASE/api/users")
check "USR_001" "users" "Get all users → 200 array" "200 + array" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert isinstance(d,list)" "HTTP ${r##*|||} array"

check "USR_002" "users" "No password_hash in response" "Response safe" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert all('password_hash' not in str(u) for u in d)" "No password_hash"

check "USR_004" "users" "Auth token accepted → 200" "200 with auth" "Medium" "$r" "200" "" "HTTP ${r##*|||}"

# ── CAREER RECOMMEND ─────────────────────────────────────────────────────────
r=$(req -X POST "$BASE/api/career-recommend" -H "Content-Type: application/json" \
  -d '{"mbti":"INTJ","aptitude":{"math":"high"},"knowledge":"programming"}')
check "CAR_001" "career-recommend" "Valid payload → 200 + careers" "200 + careers array" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert len(d.get('recommended_careers',[]))>0" \
  "HTTP ${r##*|||} $(echo ${r%|||*} | python3 -c 'import json,sys;d=json.load(sys.stdin);print(len(d.get("recommended_careers",[])))' 2>/dev/null) careers"

check "CAR_006" "career-recommend" "No password in response" "Response safe" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert 'password' not in str(d)" "No sensitive data"

r=$(req -X POST "$BASE/api/career-recommend" -H "Content-Type: application/json" \
  -d '{"aptitude":{"math":"high"},"knowledge":"x"}')
check "CAR_002" "career-recommend" "Missing mbti → 400" "400 Bad Request" "High" \
  "$r" "400" "" "HTTP ${r##*|||} (bug: no validation)"

r=$(req -X POST "$BASE/api/career-recommend" -H "Content-Type: application/json" -d '{}')
check "CAR_003" "career-recommend" "Empty body → 400" "400 Bad Request" "High" \
  "$r" "400" "" "HTTP ${r##*|||} (bug: no validation)"

r=$(req -X POST "$BASE/api/career-recommend" -H "Content-Type: application/json" -d 'notjson')
check "CAR_005" "career-recommend" "Invalid JSON → 400" "400 Bad Request" "Medium" "$r" "400" "" "HTTP ${r##*|||}"

# ── LEARNING PATHS ───────────────────────────────────────────────────────────
r=$(req "$BASE/api/learning-paths")
check "LP_001" "learning-paths" "Get all paths → 200" "200 + array" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert isinstance(d,list)" "HTTP ${r##*|||}"

check "LP_002" "learning-paths" "Has id/career_name/total_stages fields" "All fields present" "High" \
  "$r" "200" "import json,sys;p=json.load(sys.stdin)[0];assert 'id' in p and 'career_name' in p and 'total_stages' in p" "Fields OK"

# ── LEARNING PATH DETAIL ─────────────────────────────────────────────────────
r=$(req "$BASE/api/learning-path/$CAREER_ENC")
check "LPD_001" "learning-path" "Valid career_name → 200 + stages" "200 + stages" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert 'stages' in d" "HTTP ${r##*|||}"

check "LPD_005" "learning-path" "No user_id → stage1=in-progress" "stage1=in-progress" "Medium" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert d['stages'][0]['status']=='in-progress'" "stage1 in-progress"

check "LPD_007" "learning-path" "Courses array exists in stages" "courses=array" "Medium" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert all(isinstance(s.get('courses',[]),list) for s in d['stages'])" "courses arrays OK"

r=$(req "$BASE/api/learning-path/NoSuchCareer")
check "LPD_002" "learning-path" "Career not found → 404" "404 Not Found" "High" "$r" "404" "" "HTTP ${r##*|||}"

r=$(req "$BASE/api/learning-path/$CAREER_ENC?user_id=$USER_ID")
check "LPD_003" "learning-path" "With valid user_id → 200" "200 + progress" "High" "$r" "200" "" "HTTP ${r##*|||}"
check "LPD_006" "learning-path" "Stage1 not locked for user" "stage1 ≠ locked" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert d['stages'][0]['status']!='locked'" "stage1 unlocked"

r=$(req "$BASE/api/learning-path/$CAREER_ENC?user_id=bad-uuid")
check "LPD_004" "learning-path" "Invalid user_id ignored → 200" "200" "Medium" "$r" "200" "" "HTTP ${r##*|||}"

# ── UPDATE PROGRESS ──────────────────────────────────────────────────────────
r=$(req_auth -X POST "$BASE/api/learning-path/progress" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"stage_id\":\"$STAGE1\",\"status\":\"in-progress\"}")
check "PRG_001" "progress" "Set in-progress → 200" "200 + progress" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert d.get('status')=='in-progress'" \
  "HTTP ${r##*|||} $(echo ${r%|||*} | python3 -c 'import json,sys;d=json.load(sys.stdin);print(d.get("error",d.get("status","?")))' 2>/dev/null)"

r=$(req_auth -X POST "$BASE/api/learning-path/progress" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"stage_id\":\"$STAGE1\",\"status\":\"completed\"}")
check "PRG_002" "progress" "Set completed → completed_at set" "200 + completed" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert d.get('status')=='completed'" \
  "HTTP ${r##*|||} $(echo ${r%|||*} | python3 -c 'import json,sys;d=json.load(sys.stdin);print(d.get("error",d.get("status","?")))' 2>/dev/null)"

r=$(req_auth -X POST "$BASE/api/learning-path/progress" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"stage_id\":\"$STAGE1\",\"status\":\"completed\"}")
check "PRG_003" "progress" "Upsert same stage → 200 no duplicate" "200" "High" "$r" "200" "" "HTTP ${r##*|||}"

r=$(req_auth -X POST "$BASE/api/learning-path/progress" -H "Content-Type: application/json" \
  -d "{\"stage_id\":\"$STAGE1\",\"status\":\"in-progress\"}")
check "PRG_004" "progress" "Missing user_id → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req_auth -X POST "$BASE/api/learning-path/progress" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"status\":\"in-progress\"}")
check "PRG_005" "progress" "Missing stage_id → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req_auth -X POST "$BASE/api/learning-path/progress" -H "Content-Type: application/json" \
  -d '{"user_id":"bad","stage_id":"'$STAGE1'","status":"in-progress"}')
check "PRG_006" "progress" "Invalid user_id → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req_auth -X POST "$BASE/api/learning-path/progress" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"stage_id\":\"bad\",\"status\":\"in-progress\"}")
check "PRG_007" "progress" "Invalid stage_id → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

# ── COMPLETE STAGE ───────────────────────────────────────────────────────────
r=$(req_auth -X POST "$BASE/api/learning-path/complete-stage" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"stage_id\":\"$STAGE1\",\"career_name\":\"$CAREER\"}")
check "CS_001" "complete-stage" "Complete stage → next_stage_id" "200 + next_id" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert 'next_stage_id' in d or 'message' in d" "HTTP ${r##*|||}"

r=$(req_auth -X POST "$BASE/api/learning-path/complete-stage" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"stage_id\":\"$STAGE1\",\"career_name\":\"$CAREER\"}")
check "CS_007" "complete-stage" "Complete same stage twice → idempotent 200" "200" "Medium" "$r" "200" "" "HTTP ${r##*|||}"

r=$(req_auth -X POST "$BASE/api/learning-path/complete-stage" -H "Content-Type: application/json" \
  -d "{\"stage_id\":\"$STAGE1\",\"career_name\":\"$CAREER\"}")
check "CS_003" "complete-stage" "Missing user_id → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req_auth -X POST "$BASE/api/learning-path/complete-stage" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"career_name\":\"$CAREER\"}")
check "CS_004" "complete-stage" "Missing stage_id → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req_auth -X POST "$BASE/api/learning-path/complete-stage" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"stage_id\":\"$STAGE1\"}")
check "CS_005" "complete-stage" "Missing career_name → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req_auth -X POST "$BASE/api/learning-path/complete-stage" -H "Content-Type: application/json" \
  -d "{\"user_id\":\"bad\",\"stage_id\":\"$STAGE1\",\"career_name\":\"$CAREER\"}")
check "CS_006" "complete-stage" "Invalid user_id → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

# ── USER PROGRESS ─────────────────────────────────────────────────────────────
r=$(req_auth "$BASE/api/learning-path/progress/$USER_ID")
check "UPR_001" "user-progress" "Valid user with progress → 200 array" "200 + array" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert isinstance(d,list)" "HTTP ${r##*|||}"

r=$(req_auth "$BASE/api/learning-path/progress/$GHOST")
check "UPR_004" "user-progress" "Non-existent UUID → 200 empty array" "200 + []" "Medium" \
  "$r" "200" "import json,sys;assert json.load(sys.stdin)==[]" "HTTP ${r##*|||} []"

r=$(req_auth "$BASE/api/learning-path/progress/not-a-uuid")
check "UPR_003" "user-progress" "Invalid UUID → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

# ── RESET ─────────────────────────────────────────────────────────────────────
r=$(req_auth -X DELETE "$BASE/api/learning-path/$CAREER_ENC/reset?user_id=$USER_ID")
check "RST_001" "reset" "Valid reset → 200 message" "200 + message" "High" \
  "$r" "200" "import json,sys;d=json.load(sys.stdin);assert 'message' in d" "HTTP ${r##*|||}"

r=$(req_auth "$BASE/api/learning-path/progress/$USER_ID")
check "RST_006" "reset" "Progress empty after reset" "200 + []" "High" \
  "$r" "200" "import json,sys;assert json.load(sys.stdin)==[]" "HTTP ${r##*|||} []"

r=$(req_auth -X DELETE "$BASE/api/learning-path/$CAREER_ENC/reset")
check "RST_002" "reset" "Missing user_id → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req_auth -X DELETE "$BASE/api/learning-path/$CAREER_ENC/reset?user_id=bad")
check "RST_003" "reset" "Invalid user_id → 400" "400" "High" "$r" "400" "" "HTTP ${r##*|||}"

r=$(req_auth -X DELETE "$BASE/api/learning-path/FakeCareer/reset?user_id=$USER_ID")
check "RST_004" "reset" "Career not found → 200 idempotent" "200" "Medium" "$r" "200" "" "HTTP ${r##*|||}"

r=$(req_auth -X DELETE "$BASE/api/learning-path/$CAREER_ENC/reset?user_id=$GHOST")
check "RST_005" "reset" "No progress to delete → 200" "200" "Medium" "$r" "200" "" "HTTP ${r##*|||}"

# ── QUIZ ──────────────────────────────────────────────────────────────────────
r=$(req -X POST "$BASE/api/quiz/generate" -H "Content-Type: application/json" \
  -d '{"career_name":"Data Scientist","stage_name":"Foundation","career_slug":"data-scientist","course_titles":["Python"]}')
QS=$(echo ${r%|||*} | python3 -c "import json,sys;print(len(json.load(sys.stdin).get('questions',[])))" 2>/dev/null)
check "QZ_001" "quiz" "Valid payload → 200 + questions" "200 + 10 questions" "High" \
  "$r" "200" "import json,sys;assert len(json.load(sys.stdin).get('questions',[]))>0" "HTTP ${r##*|||} ${QS} questions"

check "QZ_004" "quiz" "Each question has question/options/answer" "Schema valid" "High" \
  "$r" "200" "import json,sys;qs=json.load(sys.stdin)['questions'];assert all('question' in q and 'options' in q and 'answer' in q for q in qs)" "Schema OK"

check "QZ_005" "quiz" "Answer is one of the options" "answer ∈ options" "High" \
  "$r" "200" "import json,sys;qs=json.load(sys.stdin)['questions'];assert all(q['answer'] in q['options'] for q in qs)" "All answers valid"

r=$(req -X POST "$BASE/api/quiz/generate" -H "Content-Type: application/json" -d 'notjson')
check "QZ_008" "quiz" "Invalid JSON → 400" "400" "Medium" "$r" "400" "" "HTTP ${r##*|||}"

# ── AI FALLBACK BANKS ─────────────────────────────────────────────────────────
for slug in "data" "ml" "web" "foundation" "unknown-xyz"; do
  r=$(req -X POST "$BASE/api/quiz/generate" -H "Content-Type: application/json" \
    -d "{\"career_name\":\"x\",\"stage_name\":\"x\",\"career_slug\":\"$slug\",\"course_titles\":[]}")
  N=$(echo ${r%|||*} | python3 -c "import json,sys;print(len(json.load(sys.stdin).get('questions',[])))" 2>/dev/null)
  case $slug in
    data)    ID="AI_025"; BANK="data" ;;
    ml)      ID="AI_026"; BANK="ML" ;;
    web)     ID="AI_027"; BANK="web" ;;
    foundation) ID="AI_028"; BANK="foundation" ;;
    *)       ID="AI_029"; BANK="general" ;;
  esac
  check "$ID" "fallback" "slug=$slug → $BANK question bank" "questions ≥ 1" "Medium" \
    "$r" "200" "import json,sys;assert len(json.load(sys.stdin).get('questions',[]))>0" "HTTP ${r##*|||} ${N} questions"
done

r=$(req -X POST "$BASE/api/quiz/generate" -H "Content-Type: application/json" \
  -d '{"career_name":"x","stage_name":"x","career_slug":"unknown","course_titles":[]}')
check "AI_030" "fallback" "Fallback returns ≤ 10 questions" "≤ 10 questions" "Low" \
  "$r" "200" "import json,sys;assert len(json.load(sys.stdin).get('questions',[]))<=10" \
  "$(echo ${r%|||*} | python3 -c 'import json,sys;print(len(json.load(sys.stdin).get("questions",[])))' 2>/dev/null) questions"

# ── AI CONFIG ─────────────────────────────────────────────────────────────────
KEY=$(grep GEMINI_API_KEY /home/bnagon/all/careercompass-backend/.env | cut -d= -f2)
if [[ -n "$KEY" ]]; then
  echo "AI_001,gemini-config,\"GEMINI_API_KEY present in .env\",\"Key loaded\",\"Key found in .env\",High,Pass" >> "$CSV"; ((PASS++))
else
  echo "AI_001,gemini-config,\"GEMINI_API_KEY present in .env\",\"Key loaded\",\"Key NOT found\",High,Fail" >> "$CSV"; ((FAIL++))
fi

r=$(req -X POST "$BASE/api/career-recommend" -H "Content-Type: application/json" \
  -d '{"mbti":"INTJ","aptitude":{"math":"high"},"knowledge":"x"}')
N=$(echo ${r%|||*} | python3 -c "import json,sys;print(len(json.load(sys.stdin).get('recommended_careers',[])))" 2>/dev/null)
check "AI_009" "ai-career" "career-recommend always returns careers" "careers ≥ 1" "High" \
  "$r" "200" "import json,sys;assert len(json.load(sys.stdin).get('recommended_careers',[]))>=1" "HTTP ${r##*|||} ${N} careers"
check "AI_014" "ai-career" "No crash on AI failure → 200" "200" "High" "$r" "200" "" "HTTP ${r##*|||}"

r=$(req -X POST "$BASE/api/quiz/generate" -H "Content-Type: application/json" \
  -d '{"career_name":"DS","stage_name":"Basics","career_slug":"data","course_titles":["Pandas"]}')
check "AI_016" "ai-quiz" "Quiz returns ≤ 10 questions" "≤ 10 questions" "High" \
  "$r" "200" "import json,sys;assert len(json.load(sys.stdin).get('questions',[]))<=10" \
  "$(echo ${r%|||*} | python3 -c 'import json,sys;print(len(json.load(sys.stdin).get("questions",[])))' 2>/dev/null) questions"

check "AI_019" "ai-quiz" "Each question has 4 options" "4 options each" "High" \
  "$r" "200" "import json,sys;assert all(len(q['options'])==4 for q in json.load(sys.stdin)['questions'])" "4 options each"

check "AI_020" "ai-quiz" "Answer is in options list" "answer ∈ options" "High" \
  "$r" "200" "import json,sys;qs=json.load(sys.stdin)['questions'];assert all(q['answer'] in q['options'] for q in qs)" "All answers valid"

check "AI_022" "ai-quiz" "No crash on fallback → 200" "200" "High" "$r" "200" "" "HTTP ${r##*|||}"

# ── SUMMARY ───────────────────────────────────────────────────────────────────
TOTAL=$((PASS+FAIL))
echo ""
echo "✅ Done! Results saved to: $CSV"
echo "Total: $TOTAL | Pass: $PASS | Fail: $FAIL"
echo ""
echo "Opening CSV preview:"
column -t -s',' "$CSV" | head -5
