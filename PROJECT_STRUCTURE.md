# Cấu trúc thư mục dự án — AI Interview Platform

> **Tech stack đã chốt:**
> - Frontend: Vue 3 + TailwindCSS + Pinia
> - Backend: Golang (Core & Realtime)
> - AI Service: Python (AI Orchestrator)
> - Database: PostgreSQL + Redis
> - Media: LiveKit SFU
> - Storage: S3-compatible
> - Deployment: Docker

---

## Tổng quan

```
CuocThiAI/
├── docs/                          # Tài liệu dự án (đã có)
├── frontend/                      # Vue 3 app — Owner: Lai
├── backend/                       # Golang services — Owner: Khôi + Hùng
├── ai-service/                    # Python AI Orchestrator — Owner: Khôi
├── deploy/                        # Docker & deployment configs
├── scripts/                       # Scripts tiện ích
├── .github/                       # CI/CD workflows
├── .env.example                   # Biến môi trường mẫu
├── docker-compose.yml             # Docker compose cho dev
├── Makefile                       # Lệnh dev phổ biến
└── README.md                      # Hướng dẫn chạy dự án
```

---

## Chi tiết từng phần

### 1. Frontend (`frontend/`) — Owner: Lai

```
frontend/
├── public/
│   ├── favicon.ico
│   └── index.html
├── src/
│   ├── assets/                        # Static assets
│   │   ├── images/
│   │   ├── icons/
│   │   └── fonts/
│   │
│   ├── styles/                        # Global styles
│   │   ├── main.css                   # TailwindCSS entry
│   │   ├── variables.css              # CSS custom properties
│   │   └── typography.css             # Font imports
│   │
│   ├── components/                    # Shared UI components
│   │   ├── common/                    # Design system primitives
│   │   │   ├── AppButton.vue
│   │   │   ├── AppInput.vue
│   │   │   ├── AppSelect.vue
│   │   │   ├── AppModal.vue
│   │   │   ├── AppTable.vue
│   │   │   ├── AppBadge.vue
│   │   │   ├── AppCard.vue
│   │   │   ├── AppTabs.vue
│   │   │   ├── AppToast.vue
│   │   │   ├── AppAvatar.vue
│   │   │   ├── AppDropdown.vue
│   │   │   ├── AppPagination.vue
│   │   │   ├── AppEmptyState.vue
│   │   │   ├── AppLoadingSkeleton.vue
│   │   │   └── AppErrorState.vue
│   │   │
│   │   ├── layout/                    # Layout components
│   │   │   ├── AppSidebar.vue
│   │   │   ├── AppNavbar.vue
│   │   │   ├── AppFooter.vue
│   │   │   ├── RecruiterLayout.vue
│   │   │   └── CandidateLayout.vue
│   │   │
│   │   ├── interview/                 # Interview room components (Lai + Hùng)
│   │   │   ├── VideoPanel.vue
│   │   │   ├── ChatPanel.vue
│   │   │   ├── TranscriptPanel.vue
│   │   │   ├── AISuggestionPanel.vue
│   │   │   ├── AIScorePanel.vue
│   │   │   ├── QuestionChecklist.vue
│   │   │   ├── RecruiterNotePanel.vue
│   │   │   ├── MediaControlBar.vue
│   │   │   ├── ParticipantList.vue
│   │   │   ├── InterviewTimer.vue
│   │   │   ├── ConnectionStatus.vue
│   │   │   └── ReconnectBanner.vue
│   │   │
│   │   ├── mock/                      # Mock interview components
│   │   │   ├── AIInterviewerPanel.vue
│   │   │   ├── AnswerPanel.vue
│   │   │   ├── MockFeedbackCard.vue
│   │   │   ├── QuestionProgressBar.vue
│   │   │   └── VoiceRecorder.vue
│   │   │
│   │   ├── report/                    # Report components
│   │   │   ├── ReportSummary.vue
│   │   │   ├── ScoreRadarChart.vue
│   │   │   ├── CriterionScoreCard.vue
│   │   │   ├── StrengthWeaknessList.vue
│   │   │   └── RecruiterDecisionForm.vue
│   │   │
│   │   └── candidate/                 # Candidate-specific components
│   │       ├── CVUploader.vue
│   │       ├── CVPreview.vue
│   │       ├── ProfileForm.vue
│   │       └── SkillChips.vue
│   │
│   ├── views/                         # Page views (route targets)
│   │   ├── auth/
│   │   │   ├── LoginPage.vue
│   │   │   ├── RegisterPage.vue
│   │   │   └── ForgotPasswordPage.vue
│   │   │
│   │   ├── recruiter/                 # Recruiter pages
│   │   │   ├── DashboardPage.vue
│   │   │   ├── JobListPage.vue
│   │   │   ├── JobDetailPage.vue
│   │   │   ├── JobCreatePage.vue
│   │   │   ├── CandidateListPage.vue
│   │   │   ├── CandidateDetailPage.vue
│   │   │   ├── InterviewListPage.vue
│   │   │   ├── InterviewSchedulePage.vue
│   │   │   ├── InterviewRoomPage.vue
│   │   │   ├── InterviewReportPage.vue
│   │   │   ├── QuestionBankPage.vue
│   │   │   └── SettingsPage.vue
│   │   │
│   │   ├── candidate/                 # Candidate pages
│   │   │   ├── CandidateDashboard.vue
│   │   │   ├── ProfilePage.vue
│   │   │   ├── CVUploadPage.vue
│   │   │   ├── MyInterviewsPage.vue
│   │   │   ├── InterviewWaitingRoom.vue
│   │   │   ├── CandidateInterviewRoom.vue
│   │   │   ├── MockSetupPage.vue
│   │   │   ├── MockInterviewRoom.vue
│   │   │   ├── MockResultPage.vue
│   │   │   └── PracticeHistoryPage.vue
│   │   │
│   │   └── shared/
│   │       ├── NotFoundPage.vue
│   │       └── UnauthorizedPage.vue
│   │
│   ├── router/                        # Vue Router
│   │   ├── index.js
│   │   ├── recruiter.routes.js
│   │   ├── candidate.routes.js
│   │   ├── auth.routes.js
│   │   └── guards.js                  # Route guards (auth, role)
│   │
│   ├── stores/                        # Pinia stores
│   │   ├── auth.store.js
│   │   ├── user.store.js
│   │   ├── company.store.js
│   │   ├── job.store.js
│   │   ├── candidate.store.js
│   │   ├── interview.store.js
│   │   ├── room.store.js
│   │   ├── transcript.store.js
│   │   ├── chat.store.js
│   │   ├── ai.store.js
│   │   ├── report.store.js
│   │   ├── mock.store.js
│   │   └── notification.store.js
│   │
│   ├── composables/                   # Vue composables (shared logic)
│   │   ├── useAuth.js
│   │   ├── useWebSocket.js            # WebSocket connection (Hùng define)
│   │   ├── useLiveKit.js              # LiveKit media (Hùng define)
│   │   ├── useRoom.js
│   │   ├── useTranscript.js
│   │   ├── useChat.js
│   │   ├── useMediaDevices.js
│   │   ├── useTimer.js
│   │   ├── usePagination.js
│   │   └── useToast.js
│   │
│   ├── services/                      # API service layer
│   │   ├── api.js                     # Axios instance + interceptors
│   │   ├── auth.service.js
│   │   ├── company.service.js
│   │   ├── job.service.js
│   │   ├── candidate.service.js
│   │   ├── interview.service.js
│   │   ├── room.service.js
│   │   ├── transcript.service.js
│   │   ├── ai.service.js
│   │   ├── report.service.js
│   │   ├── mock.service.js
│   │   └── file.service.js
│   │
│   ├── utils/                         # Utilities
│   │   ├── constants.js
│   │   ├── formatters.js              # Date, number, text formatters
│   │   ├── validators.js
│   │   └── helpers.js
│   │
│   ├── App.vue
│   └── main.js
│
├── .env.example
├── .env.development
├── .eslintrc.js
├── tailwind.config.js
├── postcss.config.js
├── vite.config.js
├── package.json
└── README.md
```

---

### 2. Backend (`backend/`) — Owner: Khôi (core) + Hùng (realtime)

```
backend/
├── cmd/                               # Entry points
│   ├── api/
│   │   └── main.go                    # REST API server — Khôi
│   └── realtime/
│       └── main.go                    # WebSocket realtime gateway — Hùng
│
├── internal/                          # Private application code
│   ├── config/                        # Configuration
│   │   ├── config.go
│   │   ├── database.go
│   │   ├── redis.go
│   │   ├── livekit.go
│   │   └── ai.go
│   │
│   ├── middleware/                     # HTTP middleware — Khôi
│   │   ├── auth.go                    # JWT validation
│   │   ├── rbac.go                    # Role-based access control
│   │   ├── company_scope.go           # Tenant isolation
│   │   ├── rate_limit.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── recovery.go
│   │
│   ├── models/                        # Database models — Khôi
│   │   ├── user.go
│   │   ├── company.go
│   │   ├── company_member.go
│   │   ├── job.go
│   │   ├── candidate.go
│   │   ├── candidate_job.go
│   │   ├── interview.go
│   │   ├── interview_template.go
│   │   ├── interview_question.go
│   │   ├── interview_room.go
│   │   ├── interview_participant.go
│   │   ├── interview_transcript.go
│   │   ├── interview_score.go
│   │   ├── interview_report.go
│   │   ├── ai_suggestion.go
│   │   ├── ai_prompt_template.go
│   │   ├── ai_request_log.go
│   │   ├── mock_interview.go
│   │   ├── mock_interview_answer.go
│   │   ├── file_upload.go
│   │   ├── audit_log.go
│   │   └── notification.go
│   │
│   ├── repository/                    # Data access layer — Khôi
│   │   ├── user_repo.go
│   │   ├── company_repo.go
│   │   ├── job_repo.go
│   │   ├── candidate_repo.go
│   │   ├── interview_repo.go
│   │   ├── transcript_repo.go
│   │   ├── score_repo.go
│   │   ├── report_repo.go
│   │   ├── mock_repo.go
│   │   ├── file_repo.go
│   │   └── audit_repo.go
│   │
│   ├── service/                       # Business logic — Khôi
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── company_service.go
│   │   ├── job_service.go
│   │   ├── candidate_service.go
│   │   ├── interview_service.go
│   │   ├── transcript_service.go
│   │   ├── report_service.go
│   │   ├── mock_service.go
│   │   ├── file_service.go
│   │   └── audit_service.go
│   │
│   ├── handler/                       # HTTP handlers (REST API) — Khôi
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── company_handler.go
│   │   ├── job_handler.go
│   │   ├── candidate_handler.go
│   │   ├── interview_handler.go
│   │   ├── room_handler.go
│   │   ├── transcript_handler.go
│   │   ├── ai_handler.go
│   │   ├── report_handler.go
│   │   ├── mock_handler.go
│   │   ├── file_handler.go
│   │   └── audit_handler.go
│   │
│   ├── ai/                            # AI client (Golang → Python) — Khôi
│   │   ├── client.go                  # HTTP client cho AI Orchestrator
│   │   ├── jd_analyzer.go
│   │   ├── cv_analyzer.go
│   │   ├── question_generator.go
│   │   ├── scoring_engine.go
│   │   ├── report_generator.go
│   │   └── types.go                   # AI request/response types
│   │
│   ├── queue/                         # Redis queue — Khôi
│   │   ├── worker.go
│   │   ├── dispatcher.go
│   │   ├── jobs/
│   │   │   ├── generate_report_job.go
│   │   │   ├── analyze_cv_job.go
│   │   │   └── batch_transcript_job.go
│   │   └── types.go
│   │
│   ├── realtime/                      # Realtime gateway — Hùng
│   │   ├── server.go                  # WebSocket server setup
│   │   ├── handler.go                 # HTTP → WS upgrade handler
│   │   ├── connection.go              # Client connection struct
│   │   ├── connection_manager.go      # Connection pool
│   │   ├── message_router.go          # Event routing
│   │   ├── auth.go                    # WS token authentication
│   │   ├── room.go                    # Room struct + methods
│   │   ├── room_manager.go            # Room pool management
│   │   ├── participant.go             # Participant struct
│   │   ├── broadcast.go               # Broadcast functions
│   │   ├── presence.go                # Heartbeat + presence tracking
│   │   ├── config.go                  # Realtime config
│   │   │
│   │   ├── events/                    # Event definitions
│   │   │   ├── event_types.go         # Event name constants
│   │   │   ├── event_envelope.go      # Envelope struct
│   │   │   ├── room_events.go         # Room event payloads
│   │   │   ├── interview_events.go    # Interview event payloads
│   │   │   ├── chat_events.go         # Chat event payloads
│   │   │   ├── media_events.go        # Media status payloads
│   │   │   ├── transcript_events.go   # Transcript event payloads
│   │   │   ├── ai_events.go           # AI event payloads
│   │   │   └── visibility.go          # Role visibility rules
│   │   │
│   │   └── handlers/                  # Event handlers
│   │       ├── room_handler.go        # Join/leave/heartbeat
│   │       ├── interview_handler.go   # Start/end interview
│   │       ├── chat_handler.go        # Chat messages
│   │       ├── media_handler.go       # Media status
│   │       ├── transcript_handler.go  # Transcript pipeline
│   │       ├── ai_handler.go          # AI suggestion/score bridge
│   │       └── note_handler.go        # Recruiter notes
│   │
│   ├── livekit/                       # LiveKit integration — Hùng
│   │   ├── client.go                  # LiveKit server SDK client
│   │   ├── token.go                   # Access token generation
│   │   ├── room.go                    # Room management
│   │   └── audio_hook.go             # Audio stream hook cho AI
│   │
│   ├── storage/                       # File storage (S3) — Khôi
│   │   ├── s3_client.go
│   │   ├── uploader.go
│   │   └── signed_url.go
│   │
│   ├── dto/                           # Data transfer objects
│   │   ├── request/
│   │   │   ├── auth_request.go
│   │   │   ├── job_request.go
│   │   │   ├── candidate_request.go
│   │   │   ├── interview_request.go
│   │   │   └── ai_request.go
│   │   └── response/
│   │       ├── common.go             # Envelope, pagination, error
│   │       ├── auth_response.go
│   │       ├── job_response.go
│   │       ├── candidate_response.go
│   │       ├── interview_response.go
│   │       ├── report_response.go
│   │       └── ai_response.go
│   │
│   └── pkg/                           # Shared utilities
│       ├── jwt/
│       │   └── jwt.go
│       ├── validator/
│       │   └── validator.go
│       ├── logger/
│       │   └── logger.go
│       ├── pagination/
│       │   └── pagination.go
│       └── errors/
│           └── errors.go
│
├── migrations/                        # Database migrations — Khôi
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_companies.up.sql
│   ├── 000002_create_companies.down.sql
│   ├── 000003_create_jobs.up.sql
│   ├── 000003_create_jobs.down.sql
│   ├── 000004_create_candidates.up.sql
│   ├── 000004_create_candidates.down.sql
│   ├── 000005_create_interviews.up.sql
│   ├── 000005_create_interviews.down.sql
│   ├── 000006_create_interview_rooms.up.sql
│   ├── 000006_create_interview_rooms.down.sql
│   ├── 000007_create_transcripts.up.sql
│   ├── 000007_create_transcripts.down.sql
│   ├── 000008_create_scores.up.sql
│   ├── 000008_create_scores.down.sql
│   ├── 000009_create_reports.up.sql
│   ├── 000009_create_reports.down.sql
│   ├── 000010_create_mock_interviews.up.sql
│   ├── 000010_create_mock_interviews.down.sql
│   ├── 000011_create_audit_logs.up.sql
│   ├── 000011_create_audit_logs.down.sql
│   └── README.md
│
├── seeds/                             # Seed data for dev
│   ├── seed_users.sql
│   ├── seed_companies.sql
│   └── seed_jobs.sql
│
├── tests/                             # Integration tests
│   ├── auth_test.go
│   ├── job_test.go
│   ├── interview_test.go
│   └── realtime_test.go
│
├── go.mod
├── go.sum
├── .air.toml                          # Hot reload config
└── README.md
```

---

### 3. AI Service (`ai-service/`) — Owner: Khôi

```
ai-service/
├── app/
│   ├── main.py                        # FastAPI entry point
│   ├── config.py                      # Configuration
│   │
│   ├── api/                           # API routes
│   │   ├── __init__.py
│   │   ├── jd_routes.py               # POST /ai/analyze-jd
│   │   ├── cv_routes.py               # POST /ai/analyze-cv
│   │   ├── question_routes.py         # POST /ai/generate-questions
│   │   ├── suggestion_routes.py       # POST /ai/suggest-follow-up
│   │   ├── scoring_routes.py          # POST /ai/score-answer
│   │   ├── report_routes.py           # POST /ai/generate-report
│   │   ├── mock_routes.py             # POST /ai/mock/*
│   │   └── health_routes.py           # GET /health
│   │
│   ├── services/                      # Business logic
│   │   ├── __init__.py
│   │   ├── jd_analyzer.py
│   │   ├── cv_analyzer.py
│   │   ├── question_generator.py
│   │   ├── follow_up_suggester.py
│   │   ├── answer_scorer.py
│   │   ├── report_generator.py
│   │   ├── mock_interviewer.py
│   │   └── stt_service.py            # Whisper STT
│   │
│   ├── prompts/                       # Versioned prompt templates
│   │   ├── __init__.py
│   │   ├── jd_analysis_v1.py
│   │   ├── cv_analysis_v1.py
│   │   ├── question_generation_v1.py
│   │   ├── follow_up_v1.py
│   │   ├── scoring_v1.py
│   │   ├── report_v1.py
│   │   └── mock_v1.py
│   │
│   ├── llm/                           # LLM provider abstraction
│   │   ├── __init__.py
│   │   ├── base.py                    # Abstract LLM client
│   │   ├── gemini_client.py           # Google Gemini
│   │   └── response_parser.py         # JSON response parsing
│   │
│   ├── models/                        # Pydantic models
│   │   ├── __init__.py
│   │   ├── jd_models.py
│   │   ├── cv_models.py
│   │   ├── question_models.py
│   │   ├── scoring_models.py
│   │   ├── report_models.py
│   │   └── mock_models.py
│   │
│   ├── utils/                         # Utilities
│   │   ├── __init__.py
│   │   ├── text_processing.py
│   │   ├── json_validator.py
│   │   └── guardrails.py             # Bias/hallucination checks
│   │
│   └── db/                            # Database access (nếu cần)
│       ├── __init__.py
│       └── connection.py
│
├── tests/
│   ├── test_jd_analyzer.py
│   ├── test_cv_analyzer.py
│   ├── test_scoring.py
│   ├── test_report.py
│   └── test_mock.py
│
├── requirements.txt
├── Dockerfile
└── README.md
```

---

### 4. Deploy (`deploy/`) — Shared

```
deploy/
├── docker/
│   ├── frontend.Dockerfile
│   ├── backend-api.Dockerfile
│   ├── backend-realtime.Dockerfile
│   ├── ai-service.Dockerfile
│   └── nginx.conf                     # Reverse proxy
│
├── docker-compose.yml                 # Full stack local dev
├── docker-compose.prod.yml            # Production
│
├── nginx/
│   ├── default.conf                   # API gateway / routing
│   └── ssl/                           # SSL certificates
│
├── livekit/
│   ├── livekit.yaml                   # LiveKit server config
│   └── docker-compose.livekit.yml
│
└── scripts/
    ├── setup-dev.sh                   # Setup dev environment
    ├── reset-db.sh                    # Reset database
    └── seed-data.sh                   # Seed sample data
```

---

### 5. Scripts (`scripts/`)

```
scripts/
├── generate-api-docs.sh               # Generate API documentation
├── run-migrations.sh                  # Run DB migrations
├── create-migration.sh                # Create new migration file
└── lint-all.sh                        # Lint all projects
```

---

### 6. Docs (`docs/`) — Đã có

```
docs/
├── README.md
├── PRODUCT_REQUIREMENTS.md
├── API_SPEC.md
├── DATABASE_DESIGN.md
├── REALTIME_EVENTS.md
├── AI_PROMPT_AND_SCORING.md
├── GiaoDien.md
├── AI_INTERVIEW_PROJECT_RULES.md
├── AI_INTERVIEW_PROJECT_SKILL.md
├── SYSTEM_ANALYSIS.md
│
├── Hung/
│   ├── MODULE_HUNG_REALTIME_AI.md
│   ├── CHECKLIST.md
│   └── task/                          # 24 task files
│
├── Khoi/
│   ├── MODULE_KHOI_CORE_AI_BACKEND.md
│   ├── CHECKLIST.md                   # (cần tạo)
│   └── task/                          # (cần tạo)
│
└── Lai/
    ├── MODULE_LAI_PRODUCT_UI_BUSINESS.md
    ├── CHECKLIST.md                   # (cần tạo)
    └── task/                          # (cần tạo)
```

---

### 7. Root files

```
CuocThiAI/
├── .gitignore
├── .env.example
├── docker-compose.yml
├── Makefile
├── PROJECT_STRUCTURE.md               # File này
└── README.md
```
