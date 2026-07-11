import { reactive, ref } from 'vue'

const getInitialLang = () => {
  return localStorage.getItem('app_lang') || 'vi'
}

const currentLang = ref(getInitialLang())

const dictionary = {
  vi: {
    // Top Navigation
    nav: {
      about: 'Giới thiệu',
      overview: 'Tổng quan',
      jobs: 'Tìm việc',
      interviews: 'Phỏng vấn',
      aiPractice: 'Luyện tập AI',
      results: 'Kết quả',
      notifications: 'Thông báo',
      markAllRead: 'Đánh dấu đã đọc',
      noNotifications: 'Chưa có thông báo nào.',
      viewAllNotifications: 'Xem tất cả thông báo',
      myProfile: 'Hồ sơ của tôi',
      settings: 'Cài đặt',
      logout: 'Đăng xuất',
      langTitle: 'Ngôn ngữ Giao diện',
      justNow: 'Vừa xong',
      minutesAgo: 'phút trước',
      hoursAgo: 'giờ trước',
      daysAgo: 'ngày trước',
      appliedJobs: 'Đã ứng tuyển',
      savedJobs: 'Đã lưu'
    },
    // Mock Setup Page (`/mock-setup`)
    setup: {
      brandTag: 'Phòng Luyện tập Phỏng vấn AI',
      title: 'Cấu hình Phòng Luyện tập Phỏng vấn AI',
      subtitle: 'Tùy chỉnh kỹ năng và cấp độ mong muốn. AI sẽ đóng vai chuyên gia phỏng vấn thực tế và chấm điểm bạn tự động.',
      startNow: 'Bắt đầu Luyện tập ngay',
      starting: 'Đang khởi động...',
      step1Title: '1. Chọn Vị trí & Chức danh Ứng tuyển',
      step1Desc: 'AI sẽ tự động chuẩn bị bộ câu hỏi chuyên sâu theo đúng vị trí được lựa chọn',
      step2Title: '2. Cấp độ Kỹ năng & Kinh nghiệm',
      step2Desc: 'Độ khó câu hỏi và tiêu chuẩn chấm điểm sẽ được điều chỉnh cho phù hợp với level',
      step3Title: '3. Lựa chọn Hướng & Kỹ năng Phỏng vấn',
      step3Desc: 'Tập trung luyện sâu vào khía cạnh bạn muốn nâng cấp nhất',
      step4Title: '4. Thời lượng & Ngôn ngữ Phỏng vấn',
      step4Desc: 'Chọn thời gian và ngôn ngữ trả lời thoải mái nhất',
      step5Title: '5. Chuyên gia AI & Nguồn dữ liệu CV',
      step5Desc: 'Chọn phong cách người phỏng vấn và sử dụng CV để AI hỏi sát thực tế dự án',
      realtimeScoring: 'Chấm điểm tự động real-time • Ghi âm rảnh tay hoặc gõ phím',
      btnStartStudio: 'Bắt đầu Luyện tập Phỏng vấn Ngay',
      customJD: 'Tự nhập chức danh / Mô tả công việc (JD)',
      pasteJD: 'Dán trực tiếp JD (Job Description) vào đây...',
      aiSummary: 'Tổng quan Cấu hình AI Studio',
      selectedRole: 'Vị trí phỏng vấn',
      selectedLevel: 'Cấp độ',
      selectedTrack: 'Hướng phỏng vấn',
      selectedTime: 'Thời lượng & Ngôn ngữ',
      selectedPersona: 'Phong cách AI',
      selectedCV: 'Dữ liệu CV'
    },
    // Practice History Page (`/mock-results`)
    history: {
      brandTag: 'Trung tâm Phân tích Năng lực AI',
      title: 'Kết quả & Lịch sử Phỏng vấn AI',
      subtitle: 'Xem lại điểm số chi tiết các lần phỏng vấn và nhận gợi ý cải thiện kỹ năng từ AI Coach.',
      startNew: 'Bắt đầu Luyện tập mới',
      filterTitle: 'Bộ lọc & Tìm kiếm phiên phỏng vấn',
      searchPlaceholder: 'Tìm theo chức danh, chủ đề, câu hỏi...',
      allStatus: 'Tất cả trạng thái',
      completed: 'Hoàn thành',
      inProgress: 'Đang diễn ra',
      allScores: 'Tất cả điểm số',
      highScore: 'Xuất sắc (≥ 80đ)',
      midScore: 'Khá tốt (60 - 79đ)',
      lowScore: 'Cần cải thiện (< 60đ)',
      noResults: 'Không tìm thấy phiên luyện tập nào phù hợp.',
      totalSessions: 'Tổng số phiên luyện tập',
      avgScore: 'Điểm số trung bình',
      rubricRadar: 'Biểu đồ Năng lực Rubric Đa chiều'
    },
    // Footer
    footer: {
      desc: 'Nền tảng phỏng vấn & tuyển dụng thông minh tích hợp Trí tuệ Nhân tạo real-time thế hệ mới. Tự động bóc tách CV, phỏng vấn trực tuyến và đánh giá ứng viên chuẩn Rubric chính xác tuyệt đối.',
      solutionsTitle: 'Giải pháp AI',
      sol1: 'Phỏng vấn Trực tuyến AI',
      sol2: 'Bóc tách & Chấm điểm CV',
      sol3: 'Ngân hàng Câu hỏi Chuẩn',
      sol4: 'Thang điểm Rubric Tự động',
      candTitle: 'Dành cho Ứng viên',
      cand1: 'Cơ hội Việc làm AI',
      cand2: 'Luyện tập Phỏng vấn 1-1',
      cand3: 'Lịch sử Chấm điểm AI',
      cand4: 'Hồ sơ Năng lực Trực tuyến',
      secTitle: 'Bảo mật & Hỗ trợ',
      terms: 'Điều khoản dịch vụ',
      privacy: 'Chính sách bảo mật 2FA',
      help: 'Trợ giúp'
    },
    // Dashboard (`/home`)
    dashboard: {
      welcome: 'Chào mừng trở lại',
      welcomeSub: 'Theo dõi lịch phỏng vấn, luyện tập với AI và cải thiện kỹ năng trả lời của bạn.',
      upcoming: 'Lịch sắp tới',
      completedAI: 'Luyện tập AI đã xong',
      avgScore: 'Điểm AI trung bình',
      profileComplete: 'Mức độ hoàn thiện CV',
      upcomingSection: 'Lịch phỏng vấn sắp tới',
      viewAll: 'Xem tất cả',
      noUpcoming: 'Bạn chưa có lịch phỏng vấn nào sắp tới.',
      findJobsNow: 'Tìm việc ngay',
      realMode: 'Phỏng vấn thật',
      mockMode: 'Phỏng vấn thử',
      joinNow: 'Tham gia ngay',
      bannerTitle: 'Sẵn sàng vượt qua mọi câu hỏi phỏng vấn?',
      bannerDesc: 'Trải nghiệm phỏng vấn 1-kèm-1 với AI Interviewer. Luyện tập không giới hạn, nhận phản hồi ngay lập tức.',
      bannerCta: 'Bắt đầu luyện tập',
      prepTitle: 'Hành trang ứng viên',
      uploadCV: 'Tải lên CV',
      uploadCVSub: 'Bắt buộc để AI phân tích',
      addSkills: 'Thêm Kỹ năng',
      addSkillsSub: 'Giúp nhà tuyển dụng tìm thấy bạn',
      coachTitle: 'AI Career Coach',
      coachText: 'Dựa trên kết quả phỏng vấn gần đây, tốc độ nói của bạn rất tốt, tuy nhiên bạn nên luyện tập thêm cách trả lời rành mạch các câu hỏi về ',
      coachBold: 'Kỹ năng chuyên môn sâu',
      coachCta: 'Luyện chủ đề này'
    },
    // Job Board Page (`/job-board`)
    jobs: {
      heroTitle: 'Khám phá cơ hội nghề nghiệp',
      heroDesc: 'Tìm kiếm hàng ngàn việc làm phù hợp với kỹ năng và định hướng phát triển của bạn.',
      searchPlaceholder: 'Nhập chức danh, từ khóa hoặc công ty...',
      searchBtn: 'Tìm việc ngay',
      loadingText: 'Đang tìm kiếm việc làm phù hợp...',
      emptyTitle: 'Chưa tìm thấy kết quả',
      emptyDesc: 'Rất tiếc, chúng tôi không tìm thấy vị trí nào phù hợp với từ khóa của bạn lúc này. Vui lòng thử lại với từ khóa khác.',
      clearSearch: 'Xóa tìm kiếm',
      recommendTitle: 'Việc làm đề xuất cho bạn',
      resultsFor: 'Kết quả cho',
      pageChip: 'Trang',
      viewDetail: 'Xem chi tiết',
      prevPage: 'Trang trước',
      nextPage: 'Trang sau'
    },
    // My Interviews Page (`/my-interviews`)
    interviews: {
      title: 'Phỏng vấn của tôi',
      subtitle: 'Quản lý các lịch phỏng vấn sắp tới và lịch sử phỏng vấn.',
      loadingText: 'Đang tải danh sách phỏng vấn...',
      emptyTitle: 'Chưa có lịch phỏng vấn',
      emptyDesc: 'Bạn hiện chưa có lịch phỏng vấn nào sắp tới. Khi nhà tuyển dụng gửi lời mời, lịch sẽ xuất hiện tại đây.',
      realMode: 'Phỏng vấn thật',
      mockMode: 'Phỏng vấn thử',
      completed: 'Đã hoàn thành',
      cancelled: 'Đã hủy',
      joinBtn: 'Tham gia'
    },
    // Profile Page (`/profile`)
    profile: {
      title: 'Hồ sơ cá nhân & CV',
      subtitle: 'Cập nhật thông tin để AI có thể đưa ra bài luyện tập chính xác nhất.',
      basicInfo: 'Thông tin cơ bản',
      fullName: 'Họ và Tên',
      phone: 'Số điện thoại',
      careerOrient: 'Định hướng nghề nghiệp',
      targetRole: 'Vị trí mục tiêu (Target Role)',
      level: 'Cấp độ hiện tại',
      saveBtn: 'Lưu hồ sơ',
      savingBtn: 'Đang lưu...',
      skillsTitle: 'Kỹ năng chuyên môn',
      skillsPlaceholder: 'Ví dụ: ReactJS, NodeJS, TypeScript...',
      skillsHelper: 'Phân cách các kỹ năng bằng dấu phẩy (,)',
      cvTitle: 'CV của bạn',
      dropzoneTitle: 'Tải CV lên (Nhiều file)',
      dropzoneHelper: 'PDF, DOCX (Tối đa 5MB/file)',
      uploadedTitle: 'Danh sách CV đã tải lên',
      analyzing: 'Đang phân tích',
      aiNote: 'CV của bạn sẽ được AI dùng làm ngữ cảnh (context) để đặt câu hỏi sát với kinh nghiệm thực tế của bạn trong phần Luyện phỏng vấn Mock Interview.'
    },
    // Landing Page (`/about` or `/`)
    landing: {
      badge: 'Giải pháp chuyển đổi số tuyển dụng 2026',
      heroTitle: 'Tương lai của tuyển dụng',
      heroAccent: 'Phỏng vấn thông minh cùng AI',
      heroSub: 'Tự động bóc tách CV, chấm độ phù hợp với JD và phỏng vấn trực tuyến tích hợp trí tuệ nhân tạo. Chấm dứt kỷ nguyên lọc hồ sơ thủ công.',
      btnTry: 'Trải nghiệm ngay',
      btnPractice: 'Luyện tập phỏng vấn',
      trusted: 'Được tin dùng bởi hàng ngàn chuyên gia HR',
      matchScore: 'Match Score',
      chipCv: 'CV.pdf đã bóc tách',
      chipRadar: 'Báo cáo Radar',
      statsTitle1: 'Độ chính xác AI Parsing',
      statsTitle2: 'Tốc độ xử lý mỗi CV',
      statsTitle3: 'Thời gian lọc hồ sơ',
      statsTitle4: 'Luyện tập Mock Interview',
      featuresHeading: 'Tất cả trong',
      featuresHeadingAccent: 'một nền tảng',
      featuresSub: 'Từ sàng lọc hồ sơ đến quyết định tuyển dụng — mọi công đoạn đều được AI hỗ trợ.',
      f1Title: 'Bóc tách CV tự động',
      f1Desc: 'Số hóa mọi định dạng CV, trích xuất kỹ năng & kinh nghiệm chỉ trong vài giây.',
      f2Title: 'AI đối chiếu JD',
      f2Desc: 'Thuật toán NLP chấm điểm độ phù hợp giữa ứng viên và mô tả công việc.',
      f3Title: 'Phỏng vấn LiveKit',
      f3Desc: 'Phòng họp ảo thời gian thực, tự động ghi âm và bóc băng transcript.',
      f4Title: 'AI Interviewer 24/7',
      f4Desc: 'Luyện phỏng vấn thử không giới hạn, nhận feedback chi tiết ngay lập tức.',
      f5Title: 'Báo cáo Radar đa chiều',
      f5Desc: 'Phân tích điểm mạnh, điểm yếu và gợi ý quyết định Hire / Reject.',
      f6Title: 'Bảo mật & phân quyền',
      f6Desc: 'Dữ liệu tách biệt theo doanh nghiệp, phân quyền chặt chẽ theo vai trò.',
      showcaseBadge: 'Giao diện trực quan',
      showcaseTitle: 'Mọi thứ bạn cần,',
      showcaseTitleAccent: 'trong một màn hình',
      showcaseSub: 'Theo dõi lịch phỏng vấn, kết quả luyện tập và độ phù hợp hồ sơ — tất cả gọn gàng trên một bảng điều khiển.',
      sc1: 'Bảng điều khiển tổng quan realtime',
      sc2: 'Lịch sử luyện tập & điểm số theo thời gian',
      sc3: 'Thông báo lịch phỏng vấn sắp tới',
      btnExplore: 'Khám phá ngay',
      wfTitle: 'Vận hành xuyên suốt',
      wfTitleAccent: '4 Bước',
      wfSub: 'Quy trình tuyển dụng được tự động hóa từ khâu tiếp nhận đến khi ra quyết định.',
      w1Title: 'Tải lên JD & CV',
      w1Desc: 'Hệ thống tự động số hóa và chuẩn hóa dữ liệu từ mọi định dạng CV.',
      w2Title: 'AI Đối chiếu',
      w2Desc: 'NLP đối chiếu kỹ năng ứng viên với JD, đưa ra Match Score.',
      w3Title: 'Phỏng vấn LiveKit',
      w3Desc: 'Phòng ảo thời gian thực, AI ghi âm, bóc băng và gợi ý câu hỏi.',
      w4Title: 'Báo cáo Radar',
      w4Desc: 'Báo cáo đa chiều về điểm mạnh, yếu và gợi ý quyết định.',
      audTitle: 'Đồng hành cùng bạn',
      audTitleAccent: 'trước & sau',
      audTitleEnd: 'phỏng vấn',
      audSub: 'Từ lúc chuẩn bị hồ sơ đến khi bước vào buổi phỏng vấn thật, bạn luôn có AI hỗ trợ.',
      audRecruiter: 'Chuẩn bị hồ sơ',
      ar1: 'Tải CV lên, AI bóc tách kỹ năng & kinh nghiệm.',
      ar2: 'Biết ngay mức độ phù hợp với từng tin tuyển dụng.',
      ar3: 'Gợi ý điểm cần bổ sung để hồ sơ nổi bật hơn.',
      ar4: 'Ứng tuyển nhanh, theo dõi trạng thái mọi lúc.',
      audCandidate: 'Luyện tập & tự tin',
      ac1: 'Luyện phỏng vấn với AI Interviewer 24/7.',
      ac2: 'Câu hỏi bám sát vị trí và cấp độ bạn chọn.',
      ac3: 'Nhận feedback chi tiết cho từng câu trả lời.',
      ac4: 'Phỏng vấn ngay trên trình duyệt, không cài App.',
      tipsTitle: 'Mẹo giúp bạn',
      tipsTitleAccent: 'ghi điểm',
      tipsSub: 'Vài nguyên tắc đơn giản giúp buổi phỏng vấn của bạn thuyết phục hơn.',
      t1Title: 'Trả lời theo cấu trúc STAR',
      t1Desc: 'Với câu hỏi tình huống, hãy nêu rõ Situation – Task – Action – Result để câu trả lời mạch lạc và có kết quả cụ thể.',
      t2Title: 'Gắn câu trả lời với JD',
      t2Desc: 'Đọc kỹ mô tả công việc và dẫn chứng đúng kỹ năng họ cần. Match Score của bạn sẽ cho biết nên nhấn mạnh điều gì.',
      t3Title: 'Luyện trước khi phỏng vấn thật',
      t3Desc: 'Chạy vài buổi Mock Interview với AI, đọc kỹ feedback và cải thiện điểm yếu trước khi bước vào buổi phỏng vấn chính thức.',
      faqTitle: 'Câu hỏi',
      faqTitleAccent: 'thường gặp',
      ctaTitle: 'Sẵn sàng để thay đổi cách tuyển dụng?',
      ctaSub: 'Tham gia cùng hàng ngàn chuyên gia Nhân sự đang sử dụng Interview AI.',
      ctaBtn: 'Bắt đầu miễn phí ngay hôm nay'
    }
  },
  en: {
    // Top Navigation
    nav: {
      about: 'About Us',
      overview: 'Overview',
      jobs: 'Find Jobs',
      interviews: 'Interviews',
      aiPractice: 'AI Practice',
      results: 'Results',
      notifications: 'Notifications',
      markAllRead: 'Mark all as read',
      noNotifications: 'No notifications yet.',
      viewAllNotifications: 'View all notifications',
      myProfile: 'My Profile',
      settings: 'Settings',
      logout: 'Log Out',
      langTitle: 'Interface Language',
      justNow: 'Just now',
      minutesAgo: 'minutes ago',
      hoursAgo: 'hours ago',
      daysAgo: 'days ago',
      appliedJobs: 'Applied Jobs',
      savedJobs: 'Saved Jobs'
    },
    // Mock Setup Page (`/mock-setup`)
    setup: {
      brandTag: 'Enterprise AI Interview Practice Studio',
      title: 'AI Interview Practice Room Setup',
      subtitle: 'Customize your skills and level. AI will act as a real interviewer and evaluate you automatically.',
      startNow: 'Start Practice Now',
      starting: 'Initializing...',
      step1Title: '1. Select Job Role & Target Title',
      step1Desc: 'AI automatically curates in-depth technical & behavioral questions tailored to the chosen role',
      step2Title: '2. Skill Level & Seniority',
      step2Desc: 'Question difficulty and scoring benchmarks adapt dynamically to your experience level',
      step3Title: '3. Interview Track & Focus Area',
      step3Desc: 'Deep dive into specialized domains and essential core competencies',
      step4Title: '4. Duration & Spoken Language',
      step4Desc: 'Choose your preferred interview length and communication language',
      step5Title: '5. AI Interviewer Persona & CV Source',
      step5Desc: 'Customize interviewer tone and integrate your resume for personalized project drill-downs',
      realtimeScoring: 'Real-time automatic rubric scoring • Hands-free voice or keyboard input',
      btnStartStudio: 'Start AI Mock Interview Now',
      customJD: 'Custom Job Title / Paste Job Description (JD)',
      pasteJD: 'Paste the exact Job Description (JD) here...',
      aiSummary: 'AI Studio Configuration Summary',
      selectedRole: 'Target Role',
      selectedLevel: 'Seniority Level',
      selectedTrack: 'Interview Track',
      selectedTime: 'Duration & Lang',
      selectedPersona: 'AI Persona Tone',
      selectedCV: 'Resume Context'
    },
    // Practice History Page (`/mock-results`)
    history: {
      brandTag: 'Enterprise AI Competency Analytics Center',
      title: 'AI Interview Results & Practice History',
      subtitle: 'Review detailed interview scores and receive actionable skill improvement suggestions from your AI Coach.',
      startNew: 'Start New Practice',
      filterTitle: 'Filter & Search Interview Sessions',
      searchPlaceholder: 'Search by role, topic, question...',
      allStatus: 'All Statuses',
      completed: 'Completed',
      inProgress: 'In Progress',
      allScores: 'All Scores',
      highScore: 'Excellent (≥ 80 pts)',
      midScore: 'Good (60 - 79 pts)',
      lowScore: 'Needs Improvement (< 60 pts)',
      noResults: 'No matching practice sessions found.',
      totalSessions: 'Total Sessions Practiced',
      avgScore: 'Average Rubric Score',
      rubricRadar: 'Multi-Dimensional Rubric Radar Chart'
    },
    // Footer
    footer: {
      desc: 'Next-generation intelligent interviewing & recruitment platform powered by Real-time AI. Automated resume parsing, live interactive interviews, and high-precision multi-dimensional Rubric evaluation.',
      solutionsTitle: 'AI Solutions',
      sol1: 'Live Real-time AI Interview',
      sol2: 'Automated Resume Parser',
      sol3: 'Standardized Question Bank',
      sol4: 'Automated Rubric Evaluation',
      candTitle: 'For Candidates',
      cand1: 'AI Job Board Opportunities',
      cand2: '1-on-1 AI Mock Practice',
      cand3: 'AI Scoring Analytics',
      cand4: 'Online Competency Profile',
      secTitle: 'Security & Support',
      terms: 'Terms of Service',
      privacy: '2FA Privacy Policy',
      help: 'Help Center'
    },
    // Dashboard (`/home`)
    dashboard: {
      welcome: 'Welcome back',
      welcomeSub: 'Track upcoming interview schedules, practice with AI, and elevate your response skills.',
      upcoming: 'Upcoming Schedule',
      completedAI: 'Completed AI Practices',
      avgScore: 'Average AI Score',
      profileComplete: 'Profile Completeness',
      upcomingSection: 'Upcoming Scheduled Interviews',
      viewAll: 'View All',
      noUpcoming: 'You have no upcoming scheduled interviews.',
      findJobsNow: 'Find Jobs Now',
      realMode: 'Real Interview',
      mockMode: 'Mock Practice',
      joinNow: 'Join Now',
      bannerTitle: 'Ready to master every interview question?',
      bannerDesc: 'Experience 1-on-1 practice with AI Interviewer. Unlimited sessions and instant actionable feedback.',
      bannerCta: 'Start Practice Now',
      prepTitle: 'Candidate Toolkit',
      uploadCV: 'Upload Resume',
      uploadCVSub: 'Required for AI analysis',
      addSkills: 'Add Core Skills',
      addSkillsSub: 'Help employers discover you',
      coachTitle: 'AI Career Coach',
      coachText: 'Based on recent sessions, your speaking pacing is excellent; consider deeper practice on clear structured articulation for ',
      coachBold: 'Core Technical Competencies',
      coachCta: 'Practice This Topic'
    },
    // Job Board Page (`/job-board`)
    jobs: {
      heroTitle: 'Explore Career Opportunities',
      heroDesc: 'Discover thousands of jobs perfectly tailored to your core technical skills and career aspirations.',
      searchPlaceholder: 'Enter job title, keyword, or company...',
      searchBtn: 'Search Jobs Now',
      loadingText: 'Searching matching career opportunities...',
      emptyTitle: 'No Results Found',
      emptyDesc: 'We could not find any positions matching your keywords at this time. Please try searching with different keywords.',
      clearSearch: 'Clear Search',
      recommendTitle: 'Recommended Jobs for You',
      resultsFor: 'Results for',
      pageChip: 'Page',
      viewDetail: 'View Details',
      prevPage: 'Previous Page',
      nextPage: 'Next Page'
    },
    // My Interviews Page (`/my-interviews`)
    interviews: {
      title: 'My Scheduled Interviews',
      subtitle: 'Manage your upcoming interview schedules and past practice history.',
      loadingText: 'Loading interview schedules...',
      emptyTitle: 'No Scheduled Interviews',
      emptyDesc: 'You currently have no upcoming scheduled interviews. Invitations from recruiters and active rooms will appear right here.',
      realMode: 'Real Interview',
      mockMode: 'Mock Practice',
      completed: 'Completed',
      cancelled: 'Cancelled',
      joinBtn: 'Join Room'
    },
    // Profile Page (`/profile`)
    profile: {
      title: 'Personal Profile & Resume',
      subtitle: 'Keep your information updated so AI can curate the most accurate practice questions.',
      basicInfo: 'Basic Information',
      fullName: 'Full Name',
      phone: 'Phone Number',
      careerOrient: 'Career Orientation',
      targetRole: 'Target Role',
      level: 'Current Level',
      saveBtn: 'Save Profile',
      savingBtn: 'Saving...',
      skillsTitle: 'Professional Skills',
      skillsPlaceholder: 'Example: ReactJS, NodeJS, TypeScript...',
      skillsHelper: 'Separate skills with commas (,)',
      cvTitle: 'Your Resume (CV)',
      dropzoneTitle: 'Upload Resume (Multiple files)',
      dropzoneHelper: 'PDF, DOCX (Max 5MB/file)',
      uploadedTitle: 'Uploaded Resumes',
      analyzing: 'Analyzing data',
      aiNote: 'Your resume will be utilized by AI as real context to ask highly tailored questions during Mock Interviews.'
    },
    // Landing Page (`/about` or `/`)
    landing: {
      badge: 'Recruitment Digital Transformation Solution 2026',
      heroTitle: 'Future of Hiring',
      heroAccent: 'Intelligent AI Interviewing',
      heroSub: 'Automated CV extraction, precise JD matching, and real-time live AI interviews. End the era of manual resume screening.',
      btnTry: 'Experience Now',
      btnPractice: 'Practice Interview',
      trusted: 'Trusted by thousands of HR professionals & enterprise recruiters',
      matchScore: 'Match Score',
      chipCv: 'CV.pdf Parsed',
      chipRadar: 'Radar Report',
      statsTitle1: 'AI Parsing Accuracy',
      statsTitle2: 'Processing Speed per CV',
      statsTitle3: 'Time Saved in Screening',
      statsTitle4: '24/7 Mock Interview Availability',
      featuresHeading: 'All-in-one',
      featuresHeadingAccent: 'Unified Platform',
      featuresSub: 'From automated screening to hiring decisions — every step is seamlessly empowered by AI.',
      f1Title: 'Automated CV Parsing',
      f1Desc: 'Digitize any resume format, extracting skills and work experience in mere seconds.',
      f2Title: 'AI JD Matching',
      f2Desc: 'Advanced NLP algorithms calculate exact fit scores between candidates and job descriptions.',
      f3Title: 'LiveKit Virtual Rooms',
      f3Desc: 'Real-time virtual interview rooms with automatic audio recording and live transcriptions.',
      f4Title: '24/7 AI Interviewer',
      f4Desc: 'Unlimited mock interview practice with instant, structured feedback after every session.',
      f5Title: 'Multi-dimensional Radar Report',
      f5Desc: 'Analyze strengths, weaknesses, and clear actionable Hire / Reject recommendations.',
      f6Title: 'Enterprise Security & Roles',
      f6Desc: 'Strict data isolation and role-based access control per enterprise workspace.',
      showcaseBadge: 'Intuitive Dashboard',
      showcaseTitle: 'Everything you need,',
      showcaseTitleAccent: 'in one screen',
      showcaseSub: 'Track interview schedules, practice outcomes, and resume match scores — all neatly organized on a unified dashboard.',
      sc1: 'Real-time comprehensive overview dashboard',
      sc2: 'Practice history and performance tracking over time',
      sc3: 'Instant alerts for upcoming scheduled interviews',
      btnExplore: 'Explore Now',
      wfTitle: 'Streamlined',
      wfTitleAccent: '4-Step Workflow',
      wfSub: 'The entire recruitment lifecycle automated from intake to final hiring decisions.',
      w1Title: 'Upload JD & CV',
      w1Desc: 'System automatically parses and standardizes data from any resume format.',
      w2Title: 'AI Matching Engine',
      w2Desc: 'NLP evaluates candidate competencies against JD, producing exact Match Scores.',
      w3Title: 'LiveKit Interview',
      w3Desc: 'Virtual room with live transcription, AI recording, and real-time question cues.',
      w4Title: 'Radar Assessment',
      w4Desc: 'Comprehensive multi-dimensional report detailing competencies and hiring advice.',
      audTitle: 'Your dedicated partner',
      audTitleAccent: 'before & after',
      audTitleEnd: 'the interview',
      audSub: 'From preparation to entering the real interview room, AI is always supporting you.',
      audRecruiter: 'Resume Preparation',
      ar1: 'Upload CV, AI instantly extracts your core skills and experience.',
      ar2: 'Know your exact compatibility score for every job opening.',
      ar3: 'Get actionable tips on how to make your resume stand out.',
      ar4: 'Fast one-click application and real-time status tracking.',
      audCandidate: 'Practice & Confidence',
      ac1: 'Practice unlimited sessions with AI Interviewer 24/7.',
      ac2: 'Questions calibrated specifically to your target role and seniority.',
      ac3: 'Receive structured feedback and rubrics for each answer.',
      ac4: 'Run interviews right in your browser with zero app installation.',
      tipsTitle: 'Pro Tips to',
      tipsTitleAccent: 'Excel & Score',
      tipsSub: 'Simple guiding principles that make your interview performance far more persuasive.',
      t1Title: 'Structure with the STAR Method',
      t1Desc: 'For behavioral questions, clearly articulate Situation – Task – Action – Result for compelling, structured answers.',
      t2Title: 'Align Answers to the JD',
      t2Desc: 'Carefully study the job description and highlight exact competencies needed. Your Match Score highlights key focus areas.',
      t3Title: 'Practice Before the Real Interview',
      t3Desc: 'Run a few Mock Interview rounds with AI, analyze feedback, and refine weaknesses before stepping into the real room.',
      faqTitle: 'Frequently Asked',
      faqTitleAccent: 'Questions',
      ctaTitle: 'Ready to transform how you hire?',
      ctaSub: 'Join thousands of HR leaders and engineers leveraging Interview AI today.',
      ctaBtn: 'Start Free Today'
    }
  }
}

// Global Auto DOM Translator Dictionary (Maps any remaining/un-bound Vietnamese strings directly across the entire DOM tree)
const globalViToEnMap = {
  'Tương lai của tuyển dụng': 'Future of Hiring',
  'Tương lai của Tuyển dụng': 'Future of Hiring',
  'Phỏng vấn thông minh cùng AI': 'Intelligent AI Interviewing',
  'Giải pháp chuyển đổi số tuyển dụng 2026': 'Recruitment Digital Transformation Solution 2026',
  'Giải pháp Chuyển đổi số Tuyển dụng 2026': 'Recruitment Digital Transformation Solution 2026',
  'Tự động bóc tách CV, chấm độ phù hợp với JD và phỏng vấn trực tuyến tích hợp trí tuệ nhân tạo. Chấm dứt kỷ nguyên lọc hồ sơ thủ công.': 'Automated CV extraction, precise JD matching, and real-time live AI interviews. End the era of manual resume screening.',
  'Trải nghiệm ngay': 'Experience Now',
  'Luyện tập phỏng vấn': 'Practice Interview',
  'Được tin dùng bởi hàng ngàn chuyên gia HR': 'Trusted by thousands of HR professionals',
  'CV.pdf đã bóc tách': 'CV.pdf Parsed',
  'Báo cáo Radar': 'Radar Report',
  'Độ chính xác AI Parsing': 'AI Parsing Accuracy',
  'Tốc độ xử lý mỗi CV': 'Processing Speed per CV',
  'Thời gian lọc hồ sơ': 'Time Saved in Screening',
  'Luyện tập Mock Interview': '24/7 Mock Interview Availability',
  'Tất cả trong': 'All-in-one',
  'một nền tảng': 'Unified Platform',
  'Từ sàng lọc hồ sơ đến quyết định tuyển dụng — mọi công đoạn đều được AI hỗ trợ.': 'From automated screening to hiring decisions — every step is seamlessly empowered by AI.',
  'Bóc tách CV tự động': 'Automated CV Parsing',
  'Số hóa mọi định dạng CV, trích xuất kỹ năng & kinh nghiệm chỉ trong vài giây.': 'Digitize any resume format, extracting skills and work experience in mere seconds.',
  'AI đối chiếu JD': 'AI JD Matching',
  'Thuật toán NLP chấm điểm độ phù hợp giữa ứng viên và mô tả công việc.': 'Advanced NLP algorithms calculate exact fit scores between candidates and job descriptions.',
  'Phỏng vấn LiveKit': 'LiveKit Virtual Rooms',
  'Phòng họp ảo thời gian thực, tự động ghi âm và bóc băng transcript.': 'Real-time virtual interview rooms with automatic audio recording and live transcriptions.',
  'AI Interviewer 24/7': '24/7 AI Interviewer',
  'Luyện phỏng vấn thử không giới hạn, nhận feedback chi tiết ngay lập tức.': 'Unlimited mock interview practice with instant, structured feedback after every session.',
  'Báo cáo Radar đa chiều': 'Multi-dimensional Radar Report',
  'Phân tích điểm mạnh, điểm yếu và gợi ý quyết định Hire / Reject.': 'Analyze strengths, weaknesses, and clear actionable Hire / Reject recommendations.',
  'Bảo mật & phân quyền': 'Enterprise Security & Roles',
  'Dữ liệu tách biệt theo doanh nghiệp, phân quyền chặt chẽ theo vai trò.': 'Strict data isolation and role-based access control per enterprise workspace.',
  'Giao diện trực quan': 'Intuitive Dashboard',
  'Mọi thứ bạn cần,': 'Everything you need,',
  'trong một màn hình': 'in one screen',
  'Theo dõi lịch phỏng vấn, kết quả luyện tập và độ phù hợp hồ sơ — tất cả gọn gàng trên một bảng điều khiển.': 'Track interview schedules, practice outcomes, and resume match scores — all neatly organized on a unified dashboard.',
  'Bảng điều khiển tổng quan realtime': 'Real-time comprehensive overview dashboard',
  'Lịch sử luyện tập & điểm số theo thời gian': 'Practice history and performance tracking over time',
  'Thông báo lịch phỏng vấn sắp tới': 'Instant alerts for upcoming scheduled interviews',
  'Khám phá ngay': 'Explore Now',
  'Vận hành xuyên suốt': 'Streamlined',
  '4 Bước': '4-Step Workflow',
  'Quy trình tuyển dụng được tự động hóa từ khâu tiếp nhận đến khi ra quyết định.': 'The entire recruitment lifecycle automated from intake to final hiring decisions.',
  'Tải lên JD & CV': 'Upload JD & CV',
  'Hệ thống tự động số hóa và chuẩn hóa dữ liệu từ mọi định dạng CV.': 'System automatically parses and standardizes data from any resume format.',
  'AI Đối chiếu': 'AI Matching Engine',
  'NLP đối chiếu kỹ năng ứng viên với JD, đưa ra Match Score.': 'NLP evaluates candidate competencies against JD, producing exact Match Scores.',
  'Phỏng vấn LiveKit': 'LiveKit Interview',
  'Phòng ảo thời gian thực, AI ghi âm, bóc băng và gợi ý câu hỏi.': 'Virtual room with live transcription, AI recording, and real-time question cues.',
  'Báo cáo multi-dimensional về điểm mạnh, yếu và gợi ý quyết định.': 'Comprehensive multi-dimensional report detailing competencies and hiring advice.',
  'Đồng hành cùng bạn': 'Your dedicated partner',
  'trước & sau': 'before & after',
  'phỏng vấn': 'the interview',
  'Từ lúc chuẩn bị hồ sơ đến khi bước vào buổi phỏng vấn thật, bạn luôn có AI hỗ trợ.': 'From preparation to entering the real interview room, AI is always supporting you.',
  'Chuẩn bị hồ sơ': 'Resume Preparation',
  'Tải CV lên, AI bóc tách kỹ năng & kinh nghiệm.': 'Upload CV, AI instantly extracts your core skills and experience.',
  'Biết ngay mức độ phù hợp với từng tin tuyển dụng.': 'Know your exact compatibility score for every job opening.',
  'Gợi ý điểm cần bổ sung để hồ sơ nổi bật hơn.': 'Get actionable tips on how to make your resume stand out.',
  'Ứng tuyển nhanh, theo dõi trạng thái mọi lúc.': 'Fast one-click application and real-time status tracking.',
  'Luyện tập & tự tin': 'Practice & Confidence',
  'Luyện phỏng vấn với AI Interviewer 24/7.': 'Practice unlimited sessions with AI Interviewer 24/7.',
  'Câu hỏi bám sát vị trí và cấp độ bạn chọn.': 'Questions calibrated specifically to your target role and seniority.',
  'Nhận feedback chi tiết cho từng câu trả lời.': 'Receive structured feedback and rubrics for each answer.',
  'Phỏng vấn ngay trên trình duyệt, không cài App.': 'Run interviews right in your browser with zero app installation.',
  'Mẹo giúp bạn': 'Pro Tips to',
  'ghi điểm': 'Excel & Score',
  'Vài nguyên tắc đơn giản giúp buổi phỏng vấn của bạn thuyết phục hơn.': 'Simple guiding principles that make your interview performance far more persuasive.',
  'Trả lời theo cấu trúc STAR': 'Structure with the STAR Method',
  'Với câu hỏi tình huống, hãy nêu rõ Situation – Task – Action – Result để câu trả lời mạch lạc và có kết quả cụ thể.': 'For behavioral questions, clearly articulate Situation – Task – Action – Result for compelling, structured answers.',
  'Gắn câu trả lời với JD': 'Align Answers to the JD',
  'Đọc kỹ mô tả công việc và dẫn chứng đúng kỹ năng họ cần. Match Score của bạn sẽ cho biết nên nhấn mạnh điều gì.': 'Carefully study the job description and highlight exact competencies needed. Your Match Score highlights key focus areas.',
  'Luyện trước khi phỏng vấn thật': 'Practice Before the Real Interview',
  'Chạy vài buổi Mock Interview với AI, đọc kỹ feedback và cải thiện điểm yếu trước khi bước vào buổi phỏng vấn chính thức.': 'Run a few Mock Interview rounds with AI, analyze feedback, and refine weaknesses before stepping into the real room.',
  'Câu hỏi': 'Frequently Asked',
  'thường gặp': 'Questions',
  'Sẵn sàng để thay đổi cách tuyển dụng?': 'Ready to transform how you hire?',
  'Tham gia cùng hàng ngàn chuyên gia Nhân sự đang sử dụng Interview AI.': 'Join thousands of HR leaders and engineers leveraging Interview AI today.',
  'Bắt đầu miễn phí ngay hôm nay': 'Start Free Today',
  'Khám phá cơ hội nghề nghiệp': 'Explore Career Opportunities',
  'Tìm kiếm hàng ngàn việc làm phù hợp với kỹ năng và định hướng phát triển của bạn.': 'Discover thousands of jobs perfectly tailored to your core technical skills and career aspirations.',
  'Tìm việc ngay': 'Search Jobs Now',
  'Nhập chức danh, từ khóa hoặc công ty...': 'Enter job title, keyword, or company...',
  'Đang tìm kiếm việc làm phù hợp...': 'Searching matching career opportunities...',
  'Chưa tìm thấy kết quả': 'No Results Found',
  'Rất tiếc, chúng tôi không tìm thấy vị trí nào phù hợp với từ khóa của bạn lúc này. Vui lòng thử lại với từ khóa khác.': 'We could not find any positions matching your keywords at this time. Please try searching with different keywords.',
  'Xóa tìm kiếm': 'Clear Search',
  'Việc làm đề xuất cho bạn': 'Recommended Jobs for You',
  'Xem chi tiết': 'View Details',
  'Trang trước': 'Previous Page',
  'Trang sau': 'Next Page',
  'Phỏng vấn của tôi': 'My Scheduled Interviews',
  'Quản lý các lịch phỏng vấn sắp tới và lịch sử phỏng vấn.': 'Manage your upcoming interview schedules and past practice history.',
  'Đang tải danh sách phỏng vấn...': 'Loading interview schedules...',
  'Chưa có lịch phỏng vấn': 'No Scheduled Interviews',
  'Bạn hiện chưa có lịch phỏng vấn nào sắp tới. Khi nhà tuyển dụng gửi lời mời, lịch sẽ xuất hiện tại đây.': 'You currently have no upcoming scheduled interviews. Invitations from recruiters and active rooms will appear right here.',
  'Phỏng vấn thật': 'Real Interview',
  'Phỏng vấn thử': 'Mock Practice',
  'Đã hoàn thành': 'Completed',
  'Đã hủy': 'Cancelled',
  'Tham gia': 'Join Room',
  'Hồ sơ cá nhân & CV': 'Personal Profile & Resume',
  'Cập nhật thông tin để AI có thể đưa ra bài luyện tập chính xác nhất.': 'Keep your information updated so AI can curate the most accurate practice questions.',
  'Thông tin cơ bản': 'Basic Information',
  'Họ và Tên': 'Full Name',
  'Số điện thoại': 'Phone Number',
  'Định hướng nghề nghiệp': 'Career Orientation',
  'Vị trí mục tiêu (Target Role)': 'Target Role',
  'Cấp độ hiện tại': 'Current Level',
  'Lưu hồ sơ': 'Save Profile',
  'Đang lưu...': 'Saving...',
  'Kỹ năng chuyên môn': 'Professional Skills',
  'Ví dụ: ReactJS, NodeJS, TypeScript...': 'Example: ReactJS, NodeJS, TypeScript...',
  'Phân cách các kỹ năng bằng dấu phẩy (,)': 'Separate skills with commas (,)',
  'CV của bạn': 'Your Resume (CV)',
  'Tải CV lên (Nhiều file)': 'Upload Resume (Multiple files)',
  'PDF, DOCX (Tối đa 5MB/file)': 'PDF, DOCX (Max 5MB/file)',
  'Danh sách CV đã tải lên': 'Uploaded Resumes',
  'AI đang đọc thông tin...': 'AI is extracting data...',
  'Đang phân tích': 'Analyzing data',
  'CV của bạn sẽ được AI dùng làm ngữ cảnh (context) để đặt câu hỏi sát với kinh nghiệm thực tế của bạn trong phần Luyện phỏng vấn Mock Interview.': 'Your resume will be utilized by AI as real context to ask highly tailored questions during Mock Interviews.',
  'Đăng nhập / Đăng ký': 'Sign In / Register',
  'Hủy': 'Cancel',
  'Xác nhận': 'Confirm',
  'Đóng': 'Close',

  // MockSetupPage (`/mock-setup`) & AI Studio Tile
  'Nhập Vị trí khác / Dán JD tùy chỉnh': 'Enter Custom Position / Paste Custom JD',
  'Tùy chỉnh 100% câu hỏi theo mô tả công việc cụ thể bạn đang ứng tuyển': '100% tailored questions based on your specific job description',
  'Tên Chức danh / Vị trí cụ thể': 'Specific Job Title / Target Position',
  'Ví dụ: Senior Frontend Engineer (React/Vue)': 'Example: Senior Frontend Engineer (React/Vue)',
  'Dán toàn bộ Nội dung JD (Mô tả công việc) vào đây': 'Paste Complete Job Description (JD) Content Here',
  'AI sẽ tự động đọc hiểu yêu cầu kỹ năng, dự án và trách nhiệm từ JD của bạn để đặt câu hỏi sát thực tế 100%.': 'AI will automatically analyze required skills, projects, and responsibilities from your JD to generate 100% realistic interview questions.',
  'Dán JD (Job Description) vào đây...': 'Paste Job Description (JD) here...',
  'Khuyến nghị 5 - 10 câu để đủ chiều sâu và đánh giá đa chiều nhất.': 'Recommended 5 - 10 questions for maximum depth and multi-dimensional evaluation.',
  'Chuẩn Đánh giá Rubric Enterprise': 'Enterprise Rubric Evaluation Standards',
  'AI tự động chấm điểm đa chiều theo 6 trục:': 'AI automatically evaluates across 6 dimensions:',
  'Chuyên môn': 'Technical Competency',
  'Cấu trúc STAR': 'STAR Methodology',
  'Trade-offs': 'System Trade-offs',
  'Giao tiếp': 'Communication & Clarity',
  'Áp lực & Phản xạ': 'Pressure & Quick Reflexes',
  'Thuật ngữ IT': 'IT Terminology & Vocabulary',
  'Hướng dẫn Thực chiến': 'Live Practice Guide & Tips',
  'Micro rõ ràng, không gian tĩnh.': 'Ensure clear microphone audio and quiet surroundings.',
  'Nhấn giữ phím hoặc Gõ chữ đều được.': 'Press and hold spacebar to speak or type your answers.',
  'Nhấn Kết thúc sớm để nhận Report.': 'Click End Early anytime to immediately generate your Report.',
  'Đồng bộ & hỏi dựa theo CV thực tế trong hồ sơ ứng viên': 'Sync & ask questions verified against candidate profile resume (CV)',
  '(Hồ sơ chưa tải CV lên)': '(No resume/CV uploaded in profile yet)',
  '✓ Đã sẵn sàng': '✓ Ready & Synchronized',
  'AI sẽ đọc kinh nghiệm, dự án (Projects) và công nghệ ghi trong CV của bạn để đặt câu hỏi xác thực': 'AI will read work experience, projects, and tech stack from your resume to verify real-world skills',
  'AI Interviewer Pro': 'AI Interviewer Pro',
  'Được phát triển trên nền tảng Gemini Real-time API': 'Powered by Gemini Real-time API Enterprise Engine',
  'STUDIO ONLINE': 'STUDIO ONLINE',

  // Profile Page (`/profile`)
  'Hồ sơ cá nhân & CV': 'Personal Profile & Resume',
  'Cập nhật thông tin để AI có thể đưa ra bài luyện tập chính xác nhất.': 'Keep your profile updated so AI can curate the most accurate practice questions.',
  'Ảnh đại diện': 'Profile Avatar',
  'Hỗ trợ JPG, PNG. Tối đa 2MB.': 'Supports JPG, PNG. Max file size 2MB.',
  'Giới thiệu ngắn (Bio)': 'Short Bio & Summary',
  'Viết một vài dòng giới thiệu về bản thân, mục tiêu, hoặc kinh nghiệm nổi bật...': 'Write a few lines about yourself, career objectives, or highlighted achievements...',
  'Ngôn ngữ phỏng vấn (AI)': 'AI Interview Language',
  'Tiếng Việt': 'Vietnamese',
  'English': 'English',
  'Số năm kinh nghiệm': 'Years of Experience',
  'Ví dụ: 2.5': 'Example: 2.5',
  'Tải CV lên (Nhiều file)': 'Upload Resume / CV (Multiple files)',
  'PDF, DOCX (Tối đa 5MB/file)': 'PDF, DOCX (Max 5MB/file)',
  
  // Interview Waiting Room & Camera/Mic check
  'Hệ thống sẽ yêu cầu quyền truy cập Camera và Micro ở bước tiếp theo để tiến hành phỏng vấn.': 'The system will request access to your Camera and Microphone in the next step to conduct the live interview.',
  'Tôi đã đọc, hiểu rõ và đồng ý với việc sử dụng hệ thống AI phân tích và ghi âm trong buổi phỏng vấn này.': 'I have read, understood, and agree to the use of AI analysis and audio recording during this interview.',
  'Từ chối & Quay lại': 'Decline & Go Back',
  'Tham gia phỏng vấn': 'Join Live Interview',
  
  // Practice History Page & My Applications
  'Đơn ứng tuyển của tôi': 'My Job Applications',
  'Theo dõi tiến độ các vị trí bạn đã ứng tuyển và lịch phỏng vấn sắp tới.': 'Track the progress of your job applications and upcoming interview schedules.',

  // My Applications Page (`/my-applications`) & Dropdown Menu items
  'Việc làm đã ứng tuyển': 'Applied Jobs',
  'Đã ứng tuyển': 'Applied Jobs',
  'Quản lý và theo dõi tiến độ chi tiết từng hồ sơ bạn đã nộp cho doanh nghiệp.': 'Manage and track detailed progress of every job application submitted to employers.',
  'Khám phá việc làm mới': 'Explore New Opportunities',
  'Khám phá việc làm ngay': 'Explore Jobs Now',
  'Tìm kiếm theo vị trí công việc hoặc tên công ty...': 'Search by job title or company name...',
  'Tất cả': 'All',
  'Đang chờ duyệt': 'Pending Review',
  'HR đã xem': 'HR Reviewed',
  'HR đã xem hồ sơ': 'HR Viewed Resume',
  'Đang phỏng vấn / Test': 'Interview / Assessment',
  'Đang phỏng vấn': 'Interviewing',
  'Chưa phù hợp': 'Not Suitable at this Time',
  'Hồ sơ đã được gửi đến bộ phận nhân sự và đang trong quá trình tiếp nhận.': 'Your resume has been delivered to HR and is currently pending review.',
  'Nhà tuyển dụng đã mở xem CV và hồ sơ năng lực của bạn.': 'The employer has opened and viewed your resume and profile.',
  'Bạn đã vượt qua vòng hồ sơ và đang tham gia phỏng vấn đánh giá.': 'You have passed the screening stage and are now participating in interviews.',
  'Nhà tuyển dụng đã phản hồi hồ sơ chưa phù hợp với vị trí lúc này.': 'The employer has reviewed your profile and marked it not suitable at this time.',
  'Không tìm thấy hồ sơ ứng tuyển nào': 'No job applications found',
  'Bạn chưa nộp hồ sơ vào vị trí nào trong danh mục này hoặc từ khóa tìm kiếm chưa khớp.': 'You haven\'t applied to any positions in this category, or no keywords matched.',
  'Xóa bộ lọc': 'Clear Filters',
  'Ngày nộp:': 'Applied Date:',
  'CV đã nộp:': 'Submitted CV:',
  'Xem tin': `View Job`,
  'Rút hồ sơ': 'Withdraw Application',
  'Bạn có chắc chắn muốn rút hồ sơ ứng tuyển vị trí': 'Are you sure you want to withdraw your application for position',

  // Saved Jobs Page (`/saved-jobs`) & Dropdown
  'Việc làm đã lưu': 'Saved Jobs',
  'Đã lưu': 'Saved Jobs',
  'Danh sách các cơ hội nghề nghiệp bạn quan tâm để chuẩn bị ứng tuyển.': 'List of career opportunities you bookmarked for preparation and applying.',
  'việc làm': 'jobs',
  'Tìm thêm việc làm': 'Find More Jobs',
  'Tìm theo tên công việc, kỹ năng hoặc công ty...': 'Search by job title, skills, or company...',
  'Địa điểm:': 'Location:',
  'Tất cả khu vực': 'All Locations',
  'Hà Nội': 'Hanoi',
  'TP. HCM': 'Ho Chi Minh City',
  'Đà Nẵng': 'Da Nang',
  'Làm việc từ xa (Remote)': 'Remote Work',
  'Xóa tất cả': 'Clear All',
  'Chưa có việc làm nào phù hợp': 'No matching saved jobs',
  'Bạn chưa lưu công việc nào vào danh mục quan tâm. Hãy lướt xem bảng tin tuyển dụng để tìm vị trí ưng ý nhé!': 'You haven\'t bookmarked any jobs yet. Browse the job board to find roles you like!',
  'Không tìm thấy việc làm nào khớp với từ khóa tìm kiếm hoặc bộ lọc hiện tại.': 'No saved jobs match your current search query or location filter.',
  'Đã lưu:': 'Saved on:',
  'Ứng tuyển': 'Apply Now',
  'Bỏ lưu công việc này': 'Remove bookmark',
  'Thỏa thuận': 'Negotiable',
  'Từ': 'From',
  'Đến': 'Up to',
  'Ứng viên': 'Candidate',
  'Chưa cập nhật': 'Not Updated'
}

// In-memory + LocalStorage Persistent Auto-Translation Cache
const loadAutoCache = () => {
  if (typeof window === 'undefined') return {}
  try {
    const saved = localStorage.getItem('app_vi_en_auto_cache')
    return saved ? JSON.parse(saved) : {}
  } catch (e) { return {} }
}
const autoCache = loadAutoCache()

const saveToAutoCache = (viText, enText) => {
  if (!viText || !enText || typeof window === 'undefined') return
  autoCache[viText] = enText
  try {
    localStorage.setItem('app_vi_en_auto_cache', JSON.stringify(autoCache))
  } catch (e) {}
}

// Smart Algorithmic Phrase & Rule Synthesizer for instant 0ms translation of common UI/HR/Recruitment terms
const instantPhraseSynthesizer = (text) => {
  // Check exact maps first
  if (globalViToEnMap[text]) return globalViToEnMap[text]
  if (autoCache[text]) return autoCache[text]

  // Embedded comprehensive UI / Recruiter / Candidate vocabulary dictionary
  const quickTerms = {
    'Cài đặt Hệ thống': 'System Settings',
    'Tùy chỉnh thông tin tài khoản và cấu hình hệ thống.': 'Customize account information and system configuration.',
    'Tài khoản': 'Account',
    'Thông báo': 'Notifications',
    'Quyền riêng tư & Bảo mật': 'Privacy & Security',
    'Cấu hình AI Workspace': 'AI Workspace Configuration',
    'Thông tin tài khoản': 'Account Information',
    'Tên người dùng': 'Full Name',
    'Email đăng nhập': 'Login Email',
    'Đổi mật khẩu': 'Change Password',
    'Mật khẩu hiện tại': 'Current Password',
    'Mật khẩu mới': 'New Password',
    'Xác nhận mật khẩu mới': 'Confirm New Password',
    'Lưu tùy chọn': 'Save Options',
    'Đang lưu...': 'Saving...',
    'Đang đổi...': 'Changing...',
    'Xóa tài khoản': 'Delete Account',
    'Quản lý dữ liệu': 'Data Management',
    'Hiển thị hồ sơ': 'Profile Visibility',
    'Cho phép các nhà tuyển dụng khác xem hồ sơ của tôi (Public Profile)': 'Allow recruiters to view my profile (Public Profile)',
    'Chia sẻ ẩn danh kết quả Mock Interview để cải thiện AI': 'Anonymously share Mock Interview results to improve AI',
    'Mô hình AI mặc định': 'Default AI Model',
    'Tự động bóc tách JD thành Rubric': 'Automatically parse JD into Rubrics',
    'Tự động Suggest câu hỏi follow-up trong lúc phỏng vấn': 'Automatically suggest follow-up questions during live interviews',
    'Lưu cấu hình': 'Save Configuration',
    'Tuỳ chỉnh cách AI Assistant hoạt động trong không gian làm việc của công ty bạn.': 'Customize how AI Assistant operates in your enterprise workspace.',
    'Nhận email thông báo khi có lịch phỏng vấn mới': 'Receive email notifications when new interviews are scheduled',
    'Nhận email nhắc nhở trước 1 tiếng khi diễn ra phỏng vấn': 'Receive email reminder 1 hour prior to interview',
    'Hiển thị thông báo trên trình duyệt (Browser push)': 'Display browser push notifications',
    'Thông báo & Cảnh báo': 'Notifications & Alerts',
    'Qua Email': 'Via Email',
    'Thông báo đẩy (Push Notifications)': 'Push Notifications',
    'Đăng nhập': 'Login',
    'Đăng ký': 'Register',
    'Quên mật khẩu': 'Forgot Password',
    'Tạo câu hỏi': 'Create Question',
    'Ngân hàng câu hỏi': 'Question Bank',
    'Tiêu chí đánh giá': 'Rubrics & Evaluation Criteria',
    'Quản lý tin tuyển dụng': 'Manage Job Listings',
    'Quản lý ứng viên': 'Manage Candidates',
    'Báo cáo phỏng vấn': 'Interview Reports',
    'Phòng phỏng vấn ảo': 'Virtual Interview Room',
    'Mẫu email': 'Email Templates',
    'Thêm mới': 'Add New',
    'Chỉnh sửa': 'Edit',
    'Xóa': 'Delete',
    'Tìm kiếm': 'Search',
    'Bộ lọc': 'Filters',
    'Tất cả': 'All',
    'Trạng thái': 'Status',
    'Hành động': 'Actions',
    'Ngày tạo': 'Created Date',
    'Đang tải...': 'Loading...',
    'Đã lưu thành công': 'Saved successfully!',
    'Không có dữ liệu': 'No data available',
    'Xem báo cáo': 'View Report',
    'Chi tiết': 'Details',
    'Mô tả công việc': 'Job Description',
    'Kỹ năng yêu cầu': 'Required Skills',
    'Mức lương': 'Salary Range',
    'Địa điểm': 'Location',
    'Kinh nghiệm': 'Experience',
    'Cấp độ': 'Level',
    'Hạn nộp hồ sơ': 'Application Deadline',
    'Ứng tuyển': 'Apply Now',
    'Gửi lời mời': 'Send Invitation',
    'Lịch phỏng vấn': 'Interview Schedule',
    'Điểm đánh giá': 'Evaluation Score',
    'Nhận xét của AI': 'AI Feedback & Comments'
  }
  if (quickTerms[text]) return quickTerms[text]

  // Rule-based prefixes and structural translations
  if (text.startsWith('Trang ')) return text.replace('Trang ', 'Page ')
  if (text.startsWith('Kết quả cho ')) return text.replace('Kết quả cho ', 'Results for ')
  if (text.startsWith('Danh sách ')) return text.replace('Danh sách ', 'List of ')
  if (text.startsWith('Quản lý ')) return text.replace('Quản lý ', 'Manage ')
  if (text.startsWith('Chi tiết ')) return text.replace('Chi tiết ', 'Details of ')
  if (text.startsWith('Thêm ')) return text.replace('Thêm ', 'Add ')

  return null
}

// Background Async Universal Translation Engine (Translates ANY unmapped Vietnamese text dynamically using free web translation API and caches instantly)
const pendingTranslations = new Set()
const asyncUniversalTranslate = async (viText, callback) => {
  if (!viText || pendingTranslations.has(viText)) return
  pendingTranslations.add(viText)
  try {
    // Call free translation API (MyMemory / Google GT fallback)
    const url = `https://api.mymemory.translated.net/get?q=${encodeURIComponent(viText)}&langpair=vi|en`
    const res = await fetch(url)
    const data = await res.json()
    if (data && data.responseData && data.responseData.translatedText) {
      let en = data.responseData.translatedText
      // Clean up common translation artifacts
      en = en.replace(/&quot;/g, '"').replace(/&#39;/g, "'").trim()
      if (en && en.toLowerCase() !== viText.toLowerCase()) {
        saveToAutoCache(viText, en)
        callback(en)
      }
    }
  } catch (err) {
    // Ignore network failures gracefully
  } finally {
    pendingTranslations.delete(viText)
  }
}

export const langStore = reactive({
  lang: currentLang.value,

  setLang(code) {
    if (code !== 'vi' && code !== 'en') return
    this.lang = code
    currentLang.value = code
    localStorage.setItem('app_lang', code)
    // Vue reactivity handles t() bindings automatically.
    // Schedule DOM translator for static text with a safe delay after Vue re-renders.
    if (typeof window !== 'undefined') {
      setTimeout(() => this.runGlobalDomTranslator(), 120)
    }
  },

  t(section, key) {
    const currentDict = dictionary[this.lang] || dictionary.vi
    if (currentDict[section] && currentDict[section][key]) {
      return currentDict[section][key]
    }
    // Fallback to Vietnamese if key missing in EN
    return (dictionary.vi[section] && dictionary.vi[section][key]) || key
  },

  // Universal Auto DOM Translator Engine (Scans #app and translates visible text/placeholders automatically across EVERY page without manual configuration)
  _isTranslating: false,
  _observer: null,
  _observerInitialized: false,

  runGlobalDomTranslator() {
    if (typeof window === 'undefined' || !document.body || this._isTranslating) return
    const appEl = document.getElementById('app') || document.body
    if (!appEl) return
    // Skip nodes inside Vue-managed template bindings (they contain {{ }})
    const skipTags = new Set(['SCRIPT', 'STYLE', 'NOSCRIPT', 'SVG', 'PATH', 'TEXTAREA', 'CODE', 'PRE'])

    this._isTranslating = true
    if (this._observer) {
      this._observer.disconnect()
    }

    try {
      const walkAndTranslate = (root) => {
        // 1. Walk Text Nodes
        const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT, null, false)
        let node
        while ((node = walker.nextNode())) {
          // Skip nodes inside non-translatable elements
          if (node.parentElement && skipTags.has(node.parentElement.tagName)) continue
          const text = node.nodeValue.trim()
          if (!text || text.length < 2) continue

          // Store original Vietnamese text if not already stored
          if (!node._originalViText) {
            const isKnownOrVi = /[àáạảãâầấậẩẫăằắặẳẵèéẹẻẽêềếệểễìíịỉĩòóọỏõôồốộổỗơờớợởỡùúụủũưừứựửữỳýỵỷỹđ]/i.test(text) || globalViToEnMap[text] || autoCache[text] || instantPhraseSynthesizer(text)
            if (isKnownOrVi) {
              node._originalViText = node.nodeValue
            }
          }

          if (this.lang === 'en' && node._originalViText) {
            // Skip if already translated to EN and not modified back by Vue
            if (node._translatedLang === 'en' && node.nodeValue === node._translatedEnText) {
              continue
            }

            const original = node._originalViText.trim()
            const instantResult = instantPhraseSynthesizer(original)

            if (instantResult) {
              const targetVal = node._originalViText.replace(original, instantResult)
              if (node.nodeValue !== targetVal) {
                node.nodeValue = targetVal
              }
              node._translatedEnText = targetVal
              node._translatedLang = 'en'
            } else if (/[àáạảãâầấậẩẫăằắặẳẵèéẹẻẽêềếệểễìíịỉĩòóọỏõôồốộổỗơờớợởỡùúụủũưừứựửữỳýỵỷỹđ]/i.test(original)) {
              asyncUniversalTranslate(original, (translatedEn) => {
                if (this.lang === 'en' && node && node.nodeValue) {
                  const targetVal = node._originalViText.replace(original, translatedEn)
                  if (node.nodeValue !== targetVal) {
                    this._isTranslating = true
                    if (this._observer) this._observer.disconnect()
                    node.nodeValue = targetVal
                    node._translatedEnText = targetVal
                    node._translatedLang = 'en'
                    if (this._observer) this._observer.observe(appEl, { childList: true, subtree: true })
                    this._isTranslating = false
                  }
                }
              })
            }
          } else if (this.lang === 'vi' && node._originalViText) {
            // Skip if already reverted to VI
            if (node._translatedLang === 'vi' && node.nodeValue === node._originalViText) {
              continue
            }
            if (node.nodeValue !== node._originalViText) {
              node.nodeValue = node._originalViText
            }
            node._translatedLang = 'vi'
          }
        }

        // 2. Walk Placeholders
        const inputs = root.querySelectorAll('input[placeholder], textarea[placeholder]')
        inputs.forEach(el => {
          if (!el._originalViPlaceholder) {
            el._originalViPlaceholder = el.placeholder
          }
          if (this.lang === 'en' && el._originalViPlaceholder) {
            if (el._translatedLang === 'en' && el.placeholder === el._translatedEnPlaceholder) {
              return
            }
            const orig = el._originalViPlaceholder.trim()
            const instantResult = instantPhraseSynthesizer(orig)
            if (instantResult) {
              const targetVal = el._originalViPlaceholder.replace(orig, instantResult)
              if (el.placeholder !== targetVal) el.placeholder = targetVal
              el._translatedEnPlaceholder = targetVal
              el._translatedLang = 'en'
            } else if (/[àáạảãâầấậẩẫăằắặẳẵèéẹẻẽêềếệểễìíịỉĩòóọỏõôồốộổỗơờớợởỡùúụủũưừứựửữỳýỵỷỹđ]/i.test(orig)) {
              asyncUniversalTranslate(orig, (translatedEn) => {
                if (this.lang === 'en' && el) {
                  const targetVal = el._originalViPlaceholder.replace(orig, translatedEn)
                  if (el.placeholder !== targetVal) {
                    el.placeholder = targetVal
                    el._translatedEnPlaceholder = targetVal
                    el._translatedLang = 'en'
                  }
                }
              })
            }
          } else if (this.lang === 'vi' && el._originalViPlaceholder) {
            if (el._translatedLang === 'vi' && el.placeholder === el._originalViPlaceholder) {
              return
            }
            if (el.placeholder !== el._originalViPlaceholder) el.placeholder = el._originalViPlaceholder
            el._translatedLang = 'vi'
          }
        })
      }

      walkAndTranslate(appEl)
    } finally {
      this._isTranslating = false
      if (this._observer) {
        this._observer.observe(appEl, { childList: true, subtree: true })
      }
    }
  },

  initAutoTranslator() {
    if (typeof window === 'undefined' || !document.body || this._observerInitialized) return
    this._observerInitialized = true

    this.runGlobalDomTranslator()

    let debounceTimeout = null
    this._observer = new MutationObserver(() => {
      if (this._isTranslating) return
      if (debounceTimeout) clearTimeout(debounceTimeout)
      debounceTimeout = setTimeout(() => {
        this.runGlobalDomTranslator()
      }, 50)
    })

    const appEl = document.getElementById('app') || document.body
    this._observer.observe(appEl, { childList: true, subtree: true })
  }
})
