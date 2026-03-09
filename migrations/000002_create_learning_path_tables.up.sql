-- Learning Paths table
CREATE TABLE IF NOT EXISTS learning_paths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    career_name VARCHAR(255) NOT NULL,
    description TEXT,
    total_stages INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Stages table
CREATE TABLE IF NOT EXISTS stages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learning_path_id UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    stage_number INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255),
    position_top VARCHAR(50),
    position_left VARCHAR(50),
    position_right VARCHAR(50),
    position_bottom VARCHAR(50),
    position_transform VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(learning_path_id, stage_number)
);

-- Courses table
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stage_id UUID NOT NULL REFERENCES stages(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    subtitle TEXT,
    url VARCHAR(1000),
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- User Progress table
CREATE TABLE IF NOT EXISTS user_stage_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stage_id UUID NOT NULL REFERENCES stages(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'locked',  -- locked, in-progress, completed
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, stage_id)
);

-- Indexes
CREATE INDEX idx_stages_learning_path_id ON stages(learning_path_id);
CREATE INDEX idx_courses_stage_id ON courses(stage_id);
CREATE INDEX idx_user_stage_progress_user_id ON user_stage_progress(user_id);
CREATE INDEX idx_user_stage_progress_stage_id ON user_stage_progress(stage_id);

-- Seed Data: Data Scientist Learning Path
INSERT INTO learning_paths (id, career_name, description, total_stages) VALUES
('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'Data Scientist', 'เส้นทางการเรียนรู้สู่การเป็น Data Scientist ระดับมืออาชีพ', 6);

-- Seed Stages
INSERT INTO stages (id, learning_path_id, stage_number, title, subtitle, position_top, position_left, position_right, position_bottom, position_transform) VALUES
('a0000001-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 1, 'Foundation Stage', 'พื้นฐาน', '1%', '62%', NULL, NULL, 'translateX(-50%)'),
('a0000001-0000-0000-0000-000000000002', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 2, 'Core Knowledge Stage', 'องค์ความรู้หลัก', '17%', NULL, '55%', NULL, NULL),
('a0000001-0000-0000-0000-000000000003', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 3, 'Essential Skills Stage', 'ทักษะหลัก', '26%', NULL, '15%', NULL, NULL),
('a0000001-0000-0000-0000-000000000004', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 4, 'Specialized Stage', 'ทักษะเฉพาะทาง', '38%', '23%', NULL, NULL, NULL),
('a0000001-0000-0000-0000-000000000005', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 5, 'Portfolio & Project Stage', 'โครงงานและผลงาน', NULL, NULL, '23%', '33%', NULL),
('a0000001-0000-0000-0000-000000000006', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 6, 'Career Launch Stage', 'เริ่มต้นอาชีพ', NULL, '25%', NULL, '3%', NULL);

-- Seed Courses for Stage 1 (Foundation)
INSERT INTO courses (stage_id, title, subtitle, url, sort_order) VALUES
('a0000001-0000-0000-0000-000000000001', 'Introduction to Programming', 'เรียนรู้พื้นฐานการเขียนโปรแกรม', 'https://www.coursera.org/learn/python', 1),
('a0000001-0000-0000-0000-000000000001', 'Mathematics Basics (Algebra, Calculus)', 'เรียนรู้พื้นฐานคณิตศาสตร์', 'https://www.coursera.org/learn/data-science-math-skills', 2),
('a0000001-0000-0000-0000-000000000001', 'Statistics Fundamentals', 'เรียนรู้พื้นฐานสถิติ', 'https://www.coursera.org/learn/statistics', 3),
('a0000001-0000-0000-0000-000000000001', 'Basic English for the Workplace', 'เรียนรู้พื้นฐานภาษาอังกฤษ', 'https://www.coursera.org/learn/english-for-business-and-entrepreneurship', 4),
('a0000001-0000-0000-0000-000000000001', 'Math for Machine Learning - Basics', 'เรียนรู้พื้นฐานคณิตสำหรับ ML', 'https://www.coursera.org/specializations/mathematics-machine-learning', 5),
('a0000001-0000-0000-0000-000000000001', 'Excel for Data', 'เรียนรู้การใช้ Excel สำหรับ Data', 'https://www.coursera.org/specializations/excel', 6),
('a0000001-0000-0000-0000-000000000001', 'Introduction to Cybersecurity Foundations', 'เรียนรู้พื้นฐาน Cybersecurity', 'https://www.coursera.org/specializations/intro-cyber-security', 7),
('a0000001-0000-0000-0000-000000000001', 'Introduction to Database and SQL', 'เรียนรู้พื้นฐาน Database and SQL', 'https://www.coursera.org/learn/sql-for-data-science', 8),
('a0000001-0000-0000-0000-000000000001', 'Introduction to Data Analysis', 'เรียนรู้พื้นฐาน Data Analysis', 'https://www.coursera.org/learn/data-analysis-with-python', 9);

-- Seed Courses for Stage 2 (Core Knowledge)
INSERT INTO courses (stage_id, title, subtitle, url, sort_order) VALUES
('a0000001-0000-0000-0000-000000000002', 'Python for Data Science', 'เรียนรู้ Python สำหรับ Data Science', 'https://www.coursera.org/learn/python-for-applied-data-science-ai', 1),
('a0000001-0000-0000-0000-000000000002', 'Data Wrangling with Pandas', 'จัดการข้อมูลด้วย Pandas', 'https://www.coursera.org/learn/data-cleaning-processing', 2),
('a0000001-0000-0000-0000-000000000002', 'Data Visualization with Matplotlib', 'สร้างกราฟและ Visualization', 'https://www.coursera.org/learn/python-for-data-visualization', 3),
('a0000001-0000-0000-0000-000000000002', 'SQL for Data Analysis', 'ใช้ SQL วิเคราะห์ข้อมูล', 'https://www.coursera.org/learn/sql-for-data-science-capstone', 4);

-- Seed Courses for Stage 3 (Essential Skills)
INSERT INTO courses (stage_id, title, subtitle, url, sort_order) VALUES
('a0000001-0000-0000-0000-000000000003', 'Machine Learning Fundamentals', 'พื้นฐาน Machine Learning', 'https://www.coursera.org/learn/machine-learning', 1),
('a0000001-0000-0000-0000-000000000003', 'Scikit-Learn for ML', 'ใช้งาน Scikit-Learn', 'https://www.coursera.org/learn/machine-learning-algorithms', 2),
('a0000001-0000-0000-0000-000000000003', 'Feature Engineering', 'การทำ Feature Engineering', 'https://www.coursera.org/learn/feature-engineering', 3);

-- Seed Courses for Stage 4 (Specialized)
INSERT INTO courses (stage_id, title, subtitle, url, sort_order) VALUES
('a0000001-0000-0000-0000-000000000004', 'Deep Learning with TensorFlow', 'เรียนรู้ Deep Learning', 'https://www.coursera.org/specializations/deep-learning', 1),
('a0000001-0000-0000-0000-000000000004', 'Natural Language Processing', 'การประมวลผลภาษาธรรมชาติ', 'https://www.coursera.org/specializations/natural-language-processing', 2),
('a0000001-0000-0000-0000-000000000004', 'Computer Vision Basics', 'พื้นฐาน Computer Vision', 'https://www.coursera.org/learn/introduction-computer-vision-watson-opencv', 3);

-- Seed Courses for Stage 5 (Portfolio)
INSERT INTO courses (stage_id, title, subtitle, url, sort_order) VALUES
('a0000001-0000-0000-0000-000000000005', 'Building Data Science Portfolio', 'สร้าง Portfolio สำหรับ Data Scientist', 'https://www.coursera.org/learn/data-science-capstone', 1),
('a0000001-0000-0000-0000-000000000005', 'Kaggle Competitions', 'เข้าร่วม Kaggle Competition', 'https://www.coursera.org/learn/competitive-data-science', 2);

-- Seed Courses for Stage 6 (Career Launch)
INSERT INTO courses (stage_id, title, subtitle, url, sort_order) VALUES
('a0000001-0000-0000-0000-000000000006', 'Job Interview Preparation', 'เตรียมตัวสัมภาษณ์งาน', 'https://www.coursera.org/learn/advanced-interviewing-techniques', 1),
('a0000001-0000-0000-0000-000000000006', 'Resume & LinkedIn Optimization', 'ปรับปรุง Resume และ LinkedIn', 'https://www.coursera.org/projects/optimize-linkedin-profile', 2);
