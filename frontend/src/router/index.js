import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/login' },
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
      { path: 'recruiter-room', component: () => import('../views/recruiter/InterviewRoomPage.vue') },
      { path: 'reports', component: () => import('../views/recruiter/InterviewReportPage.vue') },
      { path: 'question-bank', component: () => import('../views/recruiter/QuestionBankPage.vue') },
      { path: 'templates', component: () => import('../views/recruiter/TemplatePage.vue') },
      { path: 'settings', component: () => import('../views/recruiter/SettingsPage.vue') }
    ]
  },
  {
    path: '/',
    component: () => import('../components/layout/CandidateLayout.vue'),
    children: [
      { path: 'home', component: () => import('../views/candidate/CandidateDashboard.vue') },
      { path: 'my-interviews', component: () => import('../views/candidate/MyInterviewsPage.vue') },
      { path: 'candidate-room', component: () => import('../views/candidate/CandidateInterviewRoom.vue') },
      { path: 'mock-setup', component: () => import('../views/candidate/MockSetupPage.vue') },
      { path: 'mock-room', component: () => import('../views/candidate/MockInterviewRoom.vue') },
      { path: 'mock-results', component: () => import('../views/candidate/PracticeHistoryPage.vue') },
      { path: 'mock-results/:id', component: () => import('../views/candidate/MockResultPage.vue') },
      { path: 'profile', component: () => import('../views/candidate/ProfilePage.vue') },
      { path: 'cv', component: () => import('../views/candidate/MyCvPage.vue') }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: () => localStorage.getItem('role') === 'recruiter' ? '/dashboard' : '/home' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Simple auth guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const publicPages = ['/login', '/register', '/forgot-password', '/interview-consent', '/interview-expired']
  const authRequired = !publicPages.includes(to.path)

  if (authRequired && !token) {
    return next('/login')
  }

  next()
})

export default router
