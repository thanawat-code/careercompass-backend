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
INSERT INTO careers (title, short_description, full_description, icon_key, slug) VALUES 
('Software Engineer', 'พัฒนาซอฟต์แวร์และแอปพลิเคชัน', 'ออกแบบ เขียนโค้ด และทดสอบระบบซอฟต์แวร์ให้ทำงานได้อย่างมีประสิทธิภาพ', 'cpu', 'software-engineer'),
('Data Scientist', 'วิเคราะห์ข้อมูลขนาดใหญ่เพื่อหา Insights', 'ใช้คณิตศาสตร์และ AI วิเคราะห์ข้อมูลเพื่อช่วยองค์กรตัดสินใจ', 'chart-bar', 'data-scientist'),
('AI/ML Engineer', 'สร้างและเทรนโมเดลปัญญาประดิษฐ์', 'พัฒนาระบบ Machine Learning ให้คอมพิวเตอร์เรียนรู้และตัดสินใจได้เอง', 'cpu', 'ai-ml-engineer'),
('System Analyst', 'วิเคราะห์และวางระบบ IT ให้องค์กร', 'ออกแบบโครงสร้างระบบคอมพิวเตอร์ให้ตอบโจทย์ธุรกิจ', 'server', 'system-analyst'),
('Cybersecurity Specialist', 'ปกป้องระบบคอมพิวเตอร์จากการถูกแฮ็ก', 'ตรวจสอบช่องโหว่และป้องกันความปลอดภัยของข้อมูลองค์กร', 'shield', 'cybersecurity-specialist'),
('Cloud Architect', 'ออกแบบสถาปัตยกรรมระบบคลาวด์', 'ดูแลและวางโครงสร้าง Server บน Cloud อย่าง AWS หรือ Google Cloud', 'cloud', 'cloud-architect'),
('Management Consultant', 'ที่ปรึกษาด้านการบริหารจัดการธุรกิจ', 'วิเคราะห์ปัญหาองค์กรและเสนอแนวทางแก้ไขเพื่อเพิ่มกำไร', 'briefcase', 'management-consultant'),
('Business Analyst', 'นักวิเคราะห์ธุรกิจ', 'นำข้อมูลทางธุรกิจมาวิเคราะห์เพื่อหาโอกาสและปรับปรุงกระบวนการทำงาน', 'pie-chart', 'business-analyst'),
('UX/UI Designer', 'ออกแบบประสบการณ์และหน้าตาแอปพลิเคชัน', 'ทำความเข้าใจผู้ใช้เพื่อออกแบบเว็บและแอปให้สวยงาม ใช้งานง่าย', 'pen-tool', 'ux-ui-designer'),
('Content Creator', 'สร้างสรรค์คอนเทนต์และสื่อดิจิทัล', 'ผลิตวิดีโอ บทความ หรือภาพกราฟิกเพื่อสื่อสารกับผู้คนบนโลกออนไลน์', 'edit-3', 'content-creator'),
('HR Specialist', 'ดูแลและพัฒนาบุคลากรในองค์กร', 'สรรหาพนักงานใหม่ และดูแลสวัสดิการรวมถึงการฝึกอบรม', 'users', 'hr-specialist'),
('Psychologist / Counselor', 'นักจิตวิทยาให้คำปรึกษา', 'พูดคุยและบำบัดผู้ที่มีความเครียดหรือต้องการคำแนะนำในการใช้ชีวิต', 'heart', 'psychologist'),
('Copywriter', 'นักเขียนโฆษณาและบทความ', 'ใช้ภาษาดึงดูดใจเพื่อสร้างยอดขายหรือสร้างภาพลักษณ์ให้แบรนด์', 'edit', 'copywriter'),
('PR Specialist', 'นักประชาสัมพันธ์', 'ดูแลภาพลักษณ์ขององค์กรและสื่อสารกับสื่อมวลชน', 'mic', 'pr-specialist'),
('Educator / Teacher', 'ครูหรือผู้สอน', 'ถ่ายทอดความรู้และพัฒนาศักยภาพของผู้เรียน', 'book-open', 'educator'),
('Database Administrator', 'ผู้ดูแลฐานข้อมูลองค์กร', 'จัดการ สำรองข้อมูล และดูแล Database ให้ทำงานรวดเร็วและปลอดภัย', 'database', 'database-admin'),
('Project Manager', 'ผู้บริหารจัดการโปรเจกต์', 'วางแผน ควบคุมงบประมาณและเวลา เพื่อให้งานเสร็จตามเป้าหมาย', 'clipboard', 'project-manager'),
('Quality Assurance (QA)', 'ผู้ทดสอบและควบคุมคุณภาพซอฟต์แวร์', 'หา Bug และทดสอบระบบเพื่อความสมบูรณ์ก่อนปล่อยให้ผู้ใช้จริง', 'check-circle', 'qa-tester'),
('Accountant', 'นักบัญชี', 'จัดการงบประมาณ ตรวจสอบรายรับรายจ่าย และภาษีขององค์กร', 'file-text', 'accountant'),
('Financial Analyst', 'นักวิเคราะห์ทางการเงิน', 'ประเมินความคุ้มค่าในการลงทุนและจัดการความเสี่ยงทางการเงิน', 'dollar-sign', 'financial-analyst'),
('Supply Chain Manager', 'ผู้จัดการห่วงโซ่อุปทาน', 'ดูแลการผลิต การจัดซื้อ และการขนส่งสินค้าให้ราบรื่น', 'truck', 'supply-chain-manager'),
('Healthcare Administrator', 'ผู้บริหารจัดการสถานพยาบาล', 'ดูแลระบบการทำงานภายในโรงพยาบาลหรือคลินิกให้มีประสิทธิภาพ', 'activity', 'healthcare-admin'),
('Civil Engineer', 'วิศวกรโยธา', 'ออกแบบและควบคุมการก่อสร้างอาคาร ถนน และโครงสร้างพื้นฐาน', 'map', 'civil-engineer'),
('Network Engineer', 'ดูแลระบบเครือข่ายและการเชื่อมต่อ', 'ติดตั้งและแก้ไขปัญหาระบบ Network ให้องค์กรสื่อสารกันได้ไม่สะดุด', 'wifi', 'network-engineer'),
('Game Developer', 'นักพัฒนาเกม', 'เขียนโปรแกรมและออกแบบระบบภายในเกม', 'play-circle', 'game-developer'),
('Digital Marketer', 'นักการตลาดดิจิทัล', 'ทำโฆษณาออนไลน์ SEO และวิเคราะห์พฤติกรรมผู้บริโภคบน Social Media', 'trending-up', 'digital-marketer'),
('Graphic Designer', 'นักออกแบบกราฟิก', 'สร้างสรรค์ผลงานภาพเพื่อใช้ในสื่อต่างๆ', 'image', 'graphic-designer'),
('Event Planner', 'นักจัดงานอีเวนต์', 'ออกแบบและคุมงานจัดแสดงสินค้า งานแต่ง หรืองานคอนเสิร์ต', 'calendar', 'event-planner'),
('Sales Representative', 'ตัวแทนขาย', 'นำเสนอสินค้า เจรจาต่อรอง และสร้างความสัมพันธ์กับลูกค้า', 'shopping-cart', 'sales-representative'),
('Chef / Culinary Artist', 'เชฟและผู้เชี่ยวชาญด้านอาหาร', 'คิดค้นเมนูและประกอบอาหารด้วยศิลปะและความคิดสร้างสรรค์', 'coffee', 'chef');

-- 5. เช็คผลลัพธ์ทันที
SELECT * FROM careers;