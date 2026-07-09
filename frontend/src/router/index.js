import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('../components/layout/CandidateLayout.vue'),
    children: [
      { path: '', component: () => import('../views/public/LandingPage.vue') }
    ]
  },
  { path: '/login', component: () => import('../views/auth/LoginPage.vue') },
  { path: '/register', component: () => import('../views/auth/RegisterPage.vue') },
  { path: '/forgot-password', component: () => import('../views/auth/ForgotPasswordPage.vue') },
  { path: '/interview-consent', component: () => import('../views/candidate/InterviewWaitingRoom.vue') },
  { path: '/interview-expired', component: () => import('../views/shared/InterviewExpiredPage.vue') },
  
  {
    path: '/',
    component: () => import('../components/layout/RecruiterLayout.vue'),
    children: [
      { path: 'dashboard', component: () => import('../views/recruiter/DashboardPage.vue') },
      { path: 'jobs', component: () => import('../views/recruiter/JobListPage.vue') },
      { path: 'jobs/:id', component: () => import('../views/recruiter/JobDetailPage.vue') },
      { path: 'candidates', component: () => import('../views/recruiter/CandidateListPage.vue') },
      { path: 'candidates/:id', component: () => import('../views/recruiter/CandidateDetailPage.vue') },
      { path: 'interviews', component: () => import('../views/recruiter/InterviewListPage.vue') },
      { path: 'interviews/new', component: () => import('../views/recruiter/InterviewSchedulePage.vue') },
      { path: 'interviews/:id', component: () => import('../views/recruiter/InterviewDetailPage.vue') },
      { path: 'interviews/:id/report', component: () => import('../views/recruiter/InterviewReportPage.vue') },
      { path: 'recruiter-room/:interviewId', component: () => import('../views/recruiter/InterviewRoomPage.vue') },
      { path: 'recruiter-room', component: () => import('../views/recruiter/InterviewRoomPage.vue') },
      { path: 'reports', component: () => import('../views/recruiter/InterviewReportPage.vue') },
      { path: 'question-bank', component: () => import('../views/recruiter/QuestionBankPage.vue') },
      { path: 'rubrics', component: () => import('../views/recruiter/RubricPage.vue') },
      { path: 'templates', component: () => import('../views/recruiter/TemplatePage.vue') },
      { path: 'settings', component: () => import('../views/recruiter/SettingsPage.vue') }
    ]
  },
  {
    path: '/',
    component: () => import('../components/layout/CandidateLayout.vue'),
    children: [
      { path: 'home', component: () => import('../views/candidate/CandidateDashboard.vue') },
      { path: 'job-board', component: () => import('../views/candidate/JobBoardPage.vue') },
      { path: 'my-interviews', component: () => import('../views/candidate/MyInterviewsPage.vue') },
      { path: 'candidate-room', component: () => import('../views/candidate/CandidateInterviewRoom.vue') },
      { path: 'mock-setup', component: () => import('../views/candidate/MockSetupPage.vue') },
      { path: 'mock-room', component: () => import('../views/candidate/MockInterviewRoom.vue') },
      { path: 'mock-results', component: () => import('../views/candidate/PracticeHistoryPage.vue') },
      { path: 'mock-results/:id', component: () => import('../views/candidate/MockResultPage.vue') },
      { path: 'profile', component: () => import('../views/candidate/ProfilePage.vue') },
      { path: 'candidate-settings', component: () => import('../views/candidate/CandidateSettingsPage.vue') }
    ]
  },
  {
    path: '/careers/:company_id',
    component: () => import('../components/layout/PublicLayout.vue'),
    children: [
      { path: '', component: () => import('../views/public/CareerPage.vue') },
      { path: 'jobs/:job_id', component: () => import('../views/public/JobApplyPage.vue') }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: () => {
    // If not logged in, redirect to login, else redirect based on role
    const token = localStorage.getItem('access_token')
    if (!token) return '/login'
    return localStorage.getItem('user_role') === 'recruiter' ? '/dashboard' : '/home'
  }}
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Simple auth guard
router.beforeEach((to, from) => {
  const token = localStorage.getItem('access_token')
  const publicPages = ['/', '/login', '/register', '/forgot-password', '/interview-consent', '/interview-expired']
  const isPublicPage = publicPages.includes(to.path) || to.path.startsWith('/careers') || to.path.startsWith('/interview-consent')
  const authRequired = !isPublicPage

  if (authRequired && !token) {
    return { path: '/login', state: { message: 'Vui lòng đăng nhập để sử dụng tính năng này!', type: 'warning' } }
  }

  // Prevent logged in users from visiting login page
  if (!authRequired && token && (to.path === '/login' || to.path === '/register')) {
    const userRole = localStorage.getItem('user_role')
    return userRole === 'recruiter' ? '/dashboard' : '/home'
  }
})

export default router
