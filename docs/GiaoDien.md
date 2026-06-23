Design a complete, premium, bright, modern SaaS UI for an **AI Interview Platform**.

This product is an online interview and AI-assisted hiring platform with two main user sides:

1. **Candidate side**

   * Candidate joins real interviews with a real recruiter.
   * Candidate practices mock interviews with an AI interviewer.
   * Candidate uploads CV, manages profile, tracks interview schedule, and receives AI feedback.

2. **Interviewer / Recruiter side**

   * Recruiter creates jobs, manages candidates, schedules interviews, joins real interview rooms, uses AI real-time assistant, reviews interview reports, and makes hiring decisions.

The UI must be suitable for a serious HR/recruitment environment, not a toy AI app.

---

# 1. Product Concept

The product has two interview modes:

## Mode 1: Real Interview

Candidate talks to a **real recruiter/interviewer** through an online interview room.

AI works in the background to help the recruiter:

* Live transcript
* Suggested follow-up questions
* Candidate answer analysis
* Rubric scoring
* Interview summary
* Final interview report

Important:

* Candidate must not see internal AI scoring.
* Candidate must not see recruiter private notes.
* Candidate can see only public interview information, chat, own notes, and optional transcript if enabled.

## Mode 2: Mock Interview

Candidate talks to an **AI interviewer**.

AI behaves like a professional interviewer:

* Asks questions using text and voice
* Optionally has a professional avatar/video tile
* Listens to the candidate’s answer
* Gives feedback
* Asks follow-up questions
* Scores candidate performance
* Creates a practice report

Important:

* AI interviewer should look professional and calm.
* Do not make AI look like a fantasy character, chatbot toy, game avatar, or overly large human image.
* Mock interview should feel like serious interview practice.

---

# 2. Visual Direction

Create a premium HR SaaS interface.

The UI should feel like:

* Linear
* Notion
* Ashby
* Greenhouse
* Deel
* Modern AI productivity tools
* Professional enterprise HR software

The interface must be:

* Bright
* Clean
* Premium
* Trustworthy
* Calm
* Professional
* Easy to scan
* Suitable for long working sessions
* Suitable for HR teams and job candidates

Avoid:

* Generic admin dashboard look
* Purple-heavy UI
* Dark theme as the main theme
* Neon colors
* Gaming style
* Crypto dashboard style
* Cartoon AI style
* Oversized AI human faces
* Random gradient cards
* Too many colorful icons
* Marketing-style icon circles
* Empty large blank areas
* Crowded unreadable panels

---

# 3. Color Palette

Use a consistent bright theme.

Primary colors:

* Primary blue: #2563EB
* Primary dark: #1D4ED8
* Accent cyan: #0891B2
* Background: #F8FAFC
* Surface: #FFFFFF
* Soft surface: #F1F5F9
* Border: #E2E8F0
* Text primary: #0F172A
* Text secondary: #475569
* Text muted: #94A3B8

Status colors:

* Success: #16A34A
* Warning: #D97706
* Danger: #DC2626
* Info: #2563EB

Rules:

* Use blue only for primary actions, active navigation, selected tabs, and key highlights.
* Use cyan only as a subtle AI accent.
* Use white cards with subtle borders.
* Use very light shadows.
* Do not use purple as the main UI color.
* Do not use heavy gradients.
* Do not use colorful cards randomly.

---

# 4. Typography

Use a professional sans-serif font:

* Inter
* SF Pro
* Geist
* IBM Plex Sans

Typography rules:

* Main page title: 24–30px, semibold
* Section title: 16–18px, semibold
* Body text: 14px
* Helper text: 12–13px
* Button text: 14px, medium
* Avoid huge text inside app screens
* Keep text readable and compact

---

# 5. Layout System

Use an 8px spacing system.

General rules:

* Sidebar width: 240px
* Header height: 64px
* Page padding: 24px or 32px
* Card padding: 20px or 24px
* Border radius: 12px to 16px
* Table row height: 56px to 64px
* Form input height: 40px to 44px
* Use clean alignment
* Use enough spacing but avoid excessive blank space
* Desktop-first layout, but responsive-friendly

---

# 6. Global Component System

Design a consistent component system for both candidate and recruiter.

Components needed:

* Button
* Icon button
* Input
* Select
* Textarea
* Search bar
* Filter chip
* Status badge
* Score badge
* AI insight card
* Candidate card
* Job card
* Interview card
* Report card
* Table
* Tabs
* Modal
* Drawer
* Toast
* Tooltip
* Empty state
* Loading skeleton
* AI thinking indicator
* Transcript item
* Video tile
* Voice recording button
* Audio waveform
* Progress bar
* Stepper
* Rubric score component
* Consent notice
* Device test component
* CV preview card
* Timeline
* Notification item

---

# 7. Main Navigation

Create separate navigation for recruiter and candidate.

## Recruiter sidebar

Logo:

* Interview AI
* Small subtitle: Enterprise Tier

Navigation items:

* Tổng quan
* Việc làm
* Ứng viên
* Lịch phỏng vấn
* Báo cáo
* Kho câu hỏi
* Mẫu AI
* Cài đặt

Sidebar style:

* White background
* Subtle right border
* Active item has soft blue background
* Active item text is blue
* Icons are simple line icons
* No colored icon circles
* User profile at bottom

## Candidate navigation

Candidate side can use a simpler sidebar or top navigation.

Items:

* Tổng quan
* Phỏng vấn của tôi
* Luyện phỏng vấn AI
* Kết quả luyện tập
* Hồ sơ của tôi
* CV của tôi
* Cài đặt

Candidate navigation should feel lighter, calmer, and friendlier than recruiter navigation.

---

# 8. Recruiter / Interviewer Screens

Design the full recruiter-side experience.

---

## 8.1 Recruiter Dashboard

Create a polished dashboard for interviewers and HR teams.

Header:

* Page title: “Tổng quan”
* Subtitle: “Theo dõi lịch phỏng vấn, ứng viên và báo cáo AI hôm nay.”
* Global search
* Notification icon
* Help icon
* User avatar
* Primary button: “Tạo lịch phỏng vấn”

Top KPI cards:

1. Jobs đang mở
2. Ứng viên mới
3. Lịch phỏng vấn hôm nay
4. Báo cáo chờ xem

KPI card rules:

* White card
* Subtle border
* Small line icon
* No colorful icon circle
* Clear number
* Small change indicator if needed

Main layout:

Left column:

* “Lịch phỏng vấn sắp tới”
* Interview items:

  * Time
  * Candidate name
  * Job title
  * Interview type
  * Status badge
  * Join button

Right/main column:

* “AI Insights”

  * 3 ứng viên có mức phù hợp cao
  * 2 buổi phỏng vấn cần review
  * 1 JD có thể tối ưu thêm

* “Ứng viên mới nhất” table:

  * Candidate
  * Job
  * AI fit score
  * Status
  * Last interview
  * Action

* “Thao tác nhanh”:

  * Tạo job mới
  * Thêm ứng viên
  * Tạo lịch phỏng vấn
  * Mở kho câu hỏi

The dashboard should look like a real HR SaaS product, not a simple admin template.

---

## 8.2 Job List Page

Create a job management screen.

Header:

* Title: “Việc làm”
* Subtitle: “Quản lý các vị trí tuyển dụng và pipeline ứng viên.”
* Button: “Tạo job mới”

Content:

* Search bar
* Filters:

  * Status
  * Department
  * Level
  * Location
* Job list as table or high-quality cards

Each job shows:

* Job title
* Department
* Level
* Location
* Number of candidates
* Number of interviews
* Average fit score
* Status badge
* Actions

Include a right-side AI panel:

* “AI gợi ý tối ưu JD”
* “Job thiếu tiêu chí đánh giá”
* “Nên thêm câu hỏi technical cho vị trí này”

---

## 8.3 Job Detail Page

Create a clean and premium job detail page.

Header:

* Breadcrumb: Jobs > Product Manager
* Job title
* Status badge
* Department
* Level
* Location
* Buttons:

  * Edit Job
  * Add Candidate
  * Schedule Interview

Tabs:

* Tổng quan
* Ứng viên
* Bộ câu hỏi AI
* Rubric đánh giá
* Lịch phỏng vấn
* Báo cáo

Overview layout:

Main content:

* Candidate pipeline summary:

  * New
  * Screening
  * Interview
  * Offer
  * Rejected

* Job Description card

* Requirements card

* Interview rubric card

* Candidate pipeline table

Right sidebar:

* AI JD Summary
* AI improvement suggestions
* Missing information checklist
* Job health score

AI cards should look serious and professional, not like colorful stickers.

---

## 8.4 Candidate List Page

Create a candidate management screen.

Header:

* Title: “Ứng viên”
* Subtitle: “Theo dõi ứng viên, điểm phù hợp và lịch sử phỏng vấn.”
* Button: “Thêm ứng viên”

Content:

* Search by name, email, skill
* Filters:

  * Job
  * Status
  * Fit score
  * Source
  * Interview status

Candidate table columns:

* Candidate avatar and name
* Email
* Applied role
* Skills
* AI fit score
* Status
* Last interview
* Action

Fit score should be shown as a clean score badge or small progress bar.

---

## 8.5 Candidate Detail Page for Recruiter

Create a premium candidate profile page for recruiters.

Header:

* Candidate avatar
* Candidate name
* Applied role
* Email
* Phone
* LinkedIn
* Current status
* Buttons:

  * Schedule interview
  * View report
  * Send email
  * Update status

Main layout:

Left/main area:

* CV Preview card

  * Looks like a document preview
  * Clean and readable
  * Not too small

* Experience timeline

* Skills extracted from CV

* Interview history

* Recruiter notes

Right sidebar:

* AI Recommendation card:

  * Fit level
  * Short reason
  * Confidence
  * Recommended next step

* AI CV Summary

* Key skills

* Suggested interview questions

* Risk flags if any

Do not make this screen too empty. It should feel useful and decision-oriented.

---

## 8.6 Interview Scheduling Page

Create a wizard-style scheduling page.

Purpose:

Recruiter schedules a real interview with a candidate.

Steps:

1. Select job and candidate
2. Select interviewer and time
3. Select interview template and AI settings
4. Preview invitation email
5. Send invitation

Fields:

* Job
* Candidate
* Interviewer
* Date and time
* Duration
* Interview type:

  * HR Interview
  * Technical Interview
  * Behavioral Interview
  * Final Interview
* AI assistant enabled
* Transcript enabled
* Recording enabled
* Candidate consent required

Preview card:

* Candidate email invitation
* Interview link
* Preparation checklist

Button:

* “Gửi lời mời phỏng vấn”

---

## 8.7 Recruiter Real Interview Room

This is the most important recruiter screen.

Purpose:

Recruiter interviews a candidate in real-time while AI assists in the background.

Top header:

* Job title
* Candidate name
* Interview timer
* Recording indicator
* Transcript indicator
* AI status indicator
* End Interview button

Main layout:

Left 70%:

* Video area

  * Candidate video tile large
  * Recruiter self-view smaller
  * Professional clean video layout
  * No oversized decorative images

* Bottom control bar:

  * Mic
  * Camera
  * Screen share
  * Chat
  * Notes
  * End call

Right 30%:

Tabbed AI panel:

Tabs:

1. AI Assistant
2. CV Summary
3. Rubric
4. Notes

AI Assistant tab includes:

* Real-time analysis:

  * Độ rõ ràng
  * Mức liên quan
  * Độ tự tin
* Suggested follow-up questions
* AI warning if answer lacks evidence
* “Hỏi ngay” action for each suggested question
* AI thinking state

CV Summary tab includes:

* Short CV summary
* Key skills
* Relevant experience
* Possible gaps

Rubric tab includes:

* Criteria list
* Current AI scoring
* Evidence required
* Score confidence

Notes tab includes:

* Recruiter private notes
* Tags
* Important moments

Bottom area:

* Live Transcript panel
* Speaker labels
* Timestamps
* Important quote highlight
* AI-generated short summary

Rules:

* Candidate must not see recruiter-only AI scoring.
* AI panel should be helpful but not visually noisy.
* This screen should feel like a serious interview command center.
* Make it premium and focused.

---

## 8.8 Interview Report Page

Create a polished interview report for recruiters.

Header:

* Candidate name
* Job title
* Interview date
* Interviewer
* Final recommendation badge
* Export PDF button
* Share button

Top summary cards:

* Overall score
* AI recommendation
* Main strengths
* Main risks

Main content:

Left:

* Score breakdown by rubric:

  * Technical Knowledge
  * Problem Solving
  * Communication
  * Experience Relevance
  * Culture Fit
* Use clean progress bars or radar chart, but keep it minimal.

Right/main:

* Evidence from transcript
* Strengths with supporting quotes
* Weaknesses with supporting quotes
* Risks and follow-up recommendation
* Recruiter final decision area

Transcript section:

* Collapsible transcript
* Speaker labels
* Highlighted evidence
* AI extracted insights

The report should feel credible, clear, and suitable for real hiring decisions.

---

# 9. Candidate Screens

Design the full candidate-side experience.

---

## 9.1 Candidate Dashboard

Create a friendly and professional candidate dashboard.

Header:

* Greeting: “Chào mừng trở lại, Nguyễn Văn A”
* Subtitle: “Theo dõi lịch phỏng vấn, luyện tập với AI và cải thiện kỹ năng trả lời.”

Top cards:

1. Lịch phỏng vấn sắp tới
2. Buổi luyện tập đã hoàn thành
3. Điểm trung bình gần đây
4. Kỹ năng cần cải thiện

Main layout:

Left column:

* Upcoming Interview card:

  * Job title
  * Company name
  * Interview date/time
  * Interview type
  * Join button
  * Prepare button

* Mock Interview CTA card:

  * “Sẵn sàng luyện phỏng vấn với AI?”
  * Button: “Bắt đầu luyện tập”

Right column:

* Profile completion card:

  * CV uploaded
  * Skills added
  * Experience completed
  * Target roles selected

* AI Career Coach card:

  * Suggested next action
  * Example: “Bạn nên luyện thêm câu hỏi về React performance.”

Bottom:

* Recent practice results
* Score trend
* Recommended practice plan

The candidate dashboard should feel calm and encouraging.

---

## 9.2 Candidate Profile Page

Create a clean career profile editor.

Sections:

* Personal information:

  * Full name
  * Email
  * Phone
  * Location
  * LinkedIn
  * Portfolio/GitHub

* Career target:

  * Target role
  * Target level
  * Expected salary
  * Preferred work type
  * Preferred location

* Skills:

  * Technical skills
  * Soft skills
  * Tools
  * Languages

* Experience:

  * Company
  * Role
  * Duration
  * Description

* Education

* Certificates

Right sidebar:

* Profile completeness
* AI profile suggestions
* Missing information checklist

---

## 9.3 CV Upload and CV Analysis Page

Create a CV management page.

Content:

* Upload CV area:

  * Drag and drop
  * PDF/DOC/DOCX support
  * File size note

* CV preview:

  * Clean document preview

* AI CV Summary:

  * Candidate summary
  * Main skills
  * Experience level
  * Strengths
  * Missing information
  * Suggested improvements

* Extracted skills chips

* Suggested target roles

Buttons:

* “Phân tích lại CV”
* “Cập nhật hồ sơ từ CV”

AI feedback should feel helpful and supportive, not judgmental.

---

## 9.4 My Interviews Page

Create a page for candidates to manage real interview invitations.

Content:

* List of interviews
* Filters:

  * Upcoming
  * Completed
  * Cancelled
  * Waiting confirmation

Each interview card/table row:

* Company name
* Job title
* Interview date/time
* Interview type
* Recruiter name
* Status
* Actions:

  * Join interview
  * View details
  * Reschedule request
  * Prepare with AI

Interview detail panel:

* Job information
* Company information
* Interview schedule
* Meeting link
* Preparation checklist
* Notes for candidate
* Device test button

---

## 9.5 Candidate Real Interview Waiting Room

Design the screen before candidate enters a real interview.

Purpose:

Candidate checks device and confirms readiness before joining recruiter.

Content:

* Interview title
* Company name
* Recruiter name
* Scheduled time
* Camera preview
* Microphone test
* Speaker test
* Network status
* Consent notice:

  * “Buổi phỏng vấn có thể được ghi transcript để hỗ trợ đánh giá.”
* Checkbox:

  * “Tôi đồng ý tham gia buổi phỏng vấn”
* Button:

  * “Vào phòng phỏng vấn”

Design should reduce stress and feel calm.

---

## 9.6 Candidate Real Interview Room

Create the candidate view of a real interview room.

Important:

Candidate is talking to a real recruiter.

Candidate can see:

* Recruiter video
* Candidate self-view
* Timer
* Connection status
* Chat
* Job info
* Their own notes
* Optional transcript if enabled

Candidate cannot see:

* AI scoring
* Recruiter private notes
* Internal AI recommendation
* Rubric scores
* Hiring decision during interview

Layout:

Top header:

* Job title
* Company name
* Timer
* Connection status

Main area:

* Recruiter video large
* Candidate self-view smaller
* Control bar:

  * Mic
  * Camera
  * Chat
  * Screen share if allowed
  * Leave interview

Right panel tabs:

* Thông tin buổi phỏng vấn
* JD tóm tắt
* Ghi chú cá nhân
* Chat

Add a calm message:

* “Hãy trả lời tự nhiên. Nhà tuyển dụng sẽ dẫn dắt buổi phỏng vấn.”

This screen should be simpler and less complex than the recruiter interview room.

---

## 9.7 AI Mock Interview Setup Page

Create a setup page for candidate mock interviews.

Title:

* “Luyện phỏng vấn cùng AI”

Subtitle:

* “AI sẽ đóng vai người phỏng vấn thật, đặt câu hỏi theo vị trí ứng tuyển và đưa feedback chi tiết sau buổi luyện tập.”

Setup form:

* Select target role:

  * Frontend Developer
  * Backend Developer
  * Product Manager
  * UI/UX Designer
  * Sales
  * Marketing

* Select level:

  * Intern
  * Fresher
  * Junior
  * Middle
  * Senior

* Select interview type:

  * HR Interview
  * Technical Interview
  * Behavioral Interview
  * Final Interview

* Upload CV

* Choose answer mode:

  * Voice
  * Text
  * Voice + Text

* Choose AI interviewer style:

  * Friendly
  * Professional
  * Challenging

Button:

* “Bắt đầu Mock Interview”

Right side:

* Preview card of AI interviewer
* Practice tips
* Estimated duration
* Number of questions

---

## 9.8 AI Mock Interview Room

This is the most important candidate practice screen.

Candidate is interviewing with AI interviewer.

The AI should feel like a real interviewer with voice and optional professional avatar/video tile.

Top header:

* Target role
* Interview type
* Question progress
* Timer
* End practice button

Main layout:

Left side:

* AI Interviewer panel:

  * Professional avatar or video-like tile
  * AI name: “AI Interviewer”
  * Speaking state:

    * “AI đang hỏi…”
    * “AI đang lắng nghe…”
    * “AI đang phân tích…”

* Current question card:

  * Question text
  * Question type
  * Difficulty
  * Replay voice button

Right side:

* Candidate answer panel:

  * Voice recording button
  * Text answer area
  * Submit answer button
  * Answer timer
  * Recording waveform

Bottom area:

* Conversation transcript
* Progress through questions
* Quick tips

After candidate answers:

* AI gives short feedback:

  * Clarity
  * Relevance
  * Structure
  * Missing evidence
* AI can ask a follow-up question if needed

Style rules:

* Professional, not cartoon
* AI avatar should be subtle and realistic
* Do not use huge unrealistic human face
* Keep focus on question and answer
* Make the candidate feel guided, not judged

---

## 9.9 Mock Interview Result Page

Create the result page after candidate finishes AI mock interview.

Header:

* “Kết quả luyện phỏng vấn”
* Target role
* Interview type
* Date
* Overall score

Top summary:

* Overall score
* Communication score
* Technical score
* Structure score
* Confidence score

Main sections:

* Điểm mạnh
* Điểm cần cải thiện
* Câu trả lời tốt nhất
* Câu trả lời cần cải thiện
* Gợi ý trả lời tốt hơn
* Lộ trình luyện tập tiếp theo

Question-by-question review:

For each question:

* Question
* Candidate answer
* AI feedback
* Suggested better answer
* Score
* Tags:

  * Good structure
  * Needs evidence
  * Too short
  * Strong example

Actions:

* Practice again
* Download report
* Save to profile
* Share result if allowed

---

## 9.10 Practice History Page

Create a page for candidate practice history.

Content:

* List of previous mock interviews
* Role
* Level
* Date
* Score
* Improvement trend
* Main weakness
* Action:

  * View result
  * Practice similar questions

Charts:

* Score trend over time
* Skill improvement
* Common weak areas

---

# 10. Shared Interview Room Logic

Make sure the UI clearly separates these two rooms:

## Real Interview Room

* Candidate talks to real recruiter.
* AI assists recruiter behind the scenes.
* Recruiter has AI assistant panel.
* Candidate does not see internal scoring.

## Mock Interview Room

* Candidate talks directly to AI interviewer.
* AI asks questions.
* Candidate answers.
* AI gives feedback.
* Candidate sees AI scoring and improvement suggestions.

This distinction must be obvious in the UI.

---

# 11. Vietnamese UI Labels

Use Vietnamese labels in the UI.

Examples:

* Tổng quan
* Việc làm
* Ứng viên
* Lịch phỏng vấn
* Báo cáo
* Kho câu hỏi
* Hồ sơ của tôi
* CV của tôi
* Luyện phỏng vấn AI
* Kết quả luyện tập
* Gợi ý từ AI
* Phân tích real-time
* Câu hỏi tiếp theo
* Bắt đầu phỏng vấn
* Kết thúc phỏng vấn
* Vào phòng phỏng vấn
* Kiểm tra thiết bị
* AI đang lắng nghe
* AI đang phân tích câu trả lời
* Câu hỏi hiện tại
* Trả lời bằng giọng nói
* Trả lời bằng văn bản
* Gửi câu trả lời
* Xem feedback
* Luyện lại
* Tải báo cáo
* Xuất PDF
* Chia sẻ
* Mời ứng viên
* Tạo lịch phỏng vấn
* Điểm cần cải thiện
* Gợi ý trả lời tốt hơn

---

# 12. UX Rules

The UI must follow these UX rules:

* Recruiter side should feel productive, data-rich, and decision-oriented.
* Candidate side should feel calm, friendly, and confidence-building.
* Interview rooms must be focused and not distracting.
* AI must feel helpful, not gimmicky.
* AI feedback must be clear and explainable.
* Scoring must show evidence where appropriate.
* Candidate must not see private recruiter evaluation.
* Always show clear states:

  * Scheduled
  * Waiting
  * Active
  * Completed
  * Cancelled
  * Report Ready
  * AI Analyzing
  * Transcript Enabled
* Use loading skeletons for async data.
* Use empty states when no data exists.
* Use confirmation modal before ending interview.
* Use consent notice before recording/transcript.
* Use warning state when AI does not have enough evidence to score.

---

# 13. Screen Priority

Generate a complete polished UI concept, but prioritize screens in this order:

1. Recruiter Real Interview Room
2. Candidate Real Interview Room
3. AI Mock Interview Room
4. Interview Report Page
5. Recruiter Dashboard
6. Candidate Dashboard
7. Candidate Detail Page for Recruiter
8. Candidate Profile Page
9. Job Detail Page
10. AI Mock Interview Setup Page
11. Mock Interview Result Page
12. CV Upload and AI CV Analysis Page

All screens must look visually consistent and production-ready.

---

# 14. Final Quality Bar

The final design must look like a real product that can be shown to investors, HR teams, and candidates.

It should feel:

* Premium
* Clean
* Serious
* Modern
* Bright
* Trustworthy
* Useful
* Professional

Do not create a generic admin dashboard.
Do not create a toy AI chatbot interface.
Do not create a fantasy AI avatar product.
Create a serious, polished AI-powered interview platform.