-- 1. ล้างตารางเก่าทิ้งก่อน (เพื่อให้ข้อมูลใหม่เข้าได้)
DROP TABLE IF EXISTS user_assessment_results CASCADE;
DROP TABLE IF EXISTS careers CASCADE;

-- 2. สร้างตาราง careers ใหม่
CREATE TABLE careers (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    short_description VARCHAR(255),
    full_description TEXT,
    icon_key VARCHAR(50),
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. สร้างตารางผลลัพธ์ ใหม่
CREATE TABLE user_assessment_results (
    id SERIAL PRIMARY KEY,
    user_id INT,
    mbti_result VARCHAR(10),
    suggested_career_ids INT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. ยัดข้อมูลอาชีพทั้ง 10 อาชีพ (Seed Data)
INSERT INTO careers (title, short_description, icon_key, slug) VALUES 
('Software Engineer', 'พัฒนาซอฟต์แวร์และแอปพลิเคชัน', 'cpu', 'software-engineer'),
('Data Scientist', 'วิเคราะห์ข้อมูลขนาดใหญ่ด้วย AI', 'chart-bar', 'data-scientist'),
('UX/UI Designer', 'ออกแบบหน้าตาและการใช้งานแอปฯ', 'pen-tool', 'ux-ui-designer'),
('System Analyst', 'วิเคราะห์และวางระบบ IT องค์กร', 'server', 'system-analyst'),
('Digital Marketer', 'ทำการตลาดออนไลน์และวิเคราะห์เทรนด์', 'briefcase', 'digital-marketer'),
('Project Manager', 'บริหารจัดการโปรเจกต์ให้สำเร็จตามเป้า', 'briefcase', 'project-manager'),
('Database Administrator', 'ดูแลและจัดการฐานข้อมูลให้ปลอดภัย', 'database', 'database-admin'),
('Network Engineer', 'ดูแลระบบเครือข่ายและการสื่อสาร', 'server', 'network-engineer'),
('Content Creator', 'สร้างสรรค์คอนเทนต์ให้น่าสนใจ', 'pen-tool', 'content-creator'),
('AI Engineer', 'สร้างและเทรนโมเดลปัญญาประดิษฐ์', 'cpu', 'ai-engineer');

-- 5. เช็คผลลัพธ์ทันที
SELECT * FROM careers;