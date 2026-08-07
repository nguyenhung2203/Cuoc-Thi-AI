import { createRouter, createWebHistory } from 'vue-router'
import { authStore } from '../stores/auth.store'
import { langStore } from '../stores/lang.store'

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
  { path: '/403', component: () => import('../views/shared/ForbiddenPage.vue') },

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
      { path: 'profile', component: () => import('../views/candidate/ProfilePage.vue') },
      { path: 'job-board', component: () => import('../views/candidate/JobBoardPage.vue') },
      { path: 'my-applications', component: () => import('../views/candidate/MyApplicationsPage.vue') },
      { path: 'saved-jobs', component: () => import('../views/candidate/SavedJobsPage.vue') },
      { path: 'my-interviews', component: () => import('../views/candidate/MyInterviewsPage.vue') },
      { path: 'candidate-room', component: () => import('../views/candidate/CandidateInterviewRoom.vue') },
      { path: 'mock-setup', component: () => import('../views/candidate/MockSetupPage.vue') },
      { path: 'mock-room', component: () => import('../views/candidate/MockInterviewRoom.vue') },
      { path: 'mock-results', component: () => import('../views/candidate/PracticeHistoryPage.vue') },
      { path: 'mock-results/:id', component: () => import('../views/candidate/MockResultPage.vue') },
      { path: 'candidate-settings', redirect: '/profile?tab=security' }
    ]
  },
  {
    path: '/careers/:company_id',
    component: () => import('../components/layout/CandidateLayout.vue'),
    children: [
      { path: '', component: () => import('../views/public/CareerPage.vue') },
      { path: 'jobs/:job_id', component: () => import('../views/public/JobApplyPage.vue') }
    ]
  },
  {
    path: '/admin/login',
    component: () => import('../views/admin/AdminLoginPage.vue')
  },
  {
    path: '/admin',
    component: () => import('../components/layout/AdminLayout.vue'),
    children: [
      { path: '', redirect: 'dashboard' },
      { path: 'dashboard', component: () => import('../views/admin/AdminDashboard.vue') },
      { path: 'users', component: () => import('../views/admin/AdminUsersPage.vue') },
      { path: 'companies', component: () => import('../views/admin/AdminCompaniesPage.vue') },
      { path: 'reports', component: () => import('../views/admin/AdminReportsPage.vue') },
      { path: 'logs', component: () => import('../views/admin/AdminLogsPage.vue') },
      { path: 'settings', component: () => import('../views/admin/AdminSettingsPage.vue') }
    ]
  },
  {
    path: '/:pathMatch(.*)*', redirect: () => {
      const token = localStorage.getItem('access_token')
      if (!token) return '/login'
      return localStorage.getItem('user_role') === 'recruiter' ? '/dashboard' : '/'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Auth & Role Access Guard
router.beforeEach((to, from) => {
  const token = localStorage.getItem('access_token')
  const userRole = localStorage.getItem('user_role') || authStore.user?.role
  // Guest may browse job board / careers / landing. Candidate home & portal tools require login.
  const publicPages = ['/', '/job-board', '/login', '/register', '/forgot-password', '/interview-consent', '/interview-expired', '/admin/login', '/403']
  const isPublicPage = publicPages.includes(to.path) || to.path.startsWith('/careers') || to.path.startsWith('/interview-consent')
  const authRequired = !isPublicPage

  if (authRequired && !token) {
    if (to.path.startsWith('/admin')) {
      return { path: '/admin/login', state: { message: 'Vui lòng đăng nhập Quản trị viên (Admin) để truy cập trang này!', type: 'warning' } }
    }
    return { path: '/login', query: { redirect: to.fullPath }, state: { message: 'Vui lòng đăng nhập để sử dụng tính năng này!', type: 'warning' } }
  }

  // Admin chỉ được vào khu vực /admin
  if (token && userRole === 'admin' && !to.path.startsWith('/admin') && !isPublicPage) {
    return { path: '/admin/dashboard' }
  }

  // Prevent non-admins from accessing ANY /admin route
  if (token && to.path.startsWith('/admin') && userRole && userRole !== 'admin') {
    return { path: '/403', query: { reason: 'admin_required', attempted: to.path } }
  }

  // Prevent Candidates from accessing Recruiter-only management routes
  if (token && userRole === 'candidate') {
    const recruiterPrefixes = ['/dashboard', '/jobs', '/candidates', '/interviews', '/reports', '/question-bank', '/rubrics', '/templates', '/settings', '/recruiter-room']
    const isRecruiterRoute = recruiterPrefixes.some(prefix => to.path === prefix || to.path.startsWith(prefix + '/'))
    if (isRecruiterRoute && !to.path.startsWith('/candidate-settings')) {
      return { path: '/403', query: { reason: 'recruiter_required', attempted: to.path } }
    }
  }

  // Prevent Recruiters from accessing Candidate-only practice / mock interview routes
  if (token && userRole === 'recruiter') {
    const candidatePrefixes = ['/mock-setup', '/mock-room', '/mock-results', '/my-interviews', '/candidate-room']
    const isCandidateRoute = candidatePrefixes.some(prefix => to.path === prefix || to.path.startsWith(prefix + '/'))
    if (isCandidateRoute) {
      return { path: '/403', query: { reason: 'candidate_required', attempted: to.path } }
    }
  }

  // Prevent logged in users from visiting login/register page
  if (!authRequired && token && (to.path === '/login' || to.path === '/register' || (to.path === '/admin/login' && userRole === 'admin'))) {
    if (userRole === 'admin') return '/admin/dashboard'
    return userRole === 'recruiter' ? '/dashboard' : '/'
  }
})

router.afterEach(() => {
  if (typeof window !== 'undefined') {
    setTimeout(() => {
      langStore.runGlobalDomTranslator()
    }, 200)
  }
})

export default router
