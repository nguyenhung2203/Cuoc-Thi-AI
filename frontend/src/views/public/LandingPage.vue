<script setup>
import { useRouter } from 'vue-router'
import { Sparkles, BrainCircuit, Users, Video, Zap, FileText, BarChart, CheckCircle2, ArrowRight } from 'lucide-vue-next'
import { ref, onMounted, onUnmounted } from 'vue'

const router = useRouter()

const handleGetStarted = () => {
  router.push('/login')
}

const currentImageIndex = ref(0)
const images = [
  '/images/hero_office.png',
  '/images/hero_network.png',
  '/images/hero_dashboard.png'
]

let carouselInterval;

// Simple Intersection Observer for scroll animations

// Simple Intersection Observer for scroll animations
onMounted(() => {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible')
      }
    })
  }, { threshold: 0.1 })

  document.querySelectorAll('.animate-on-scroll').forEach((el) => {
    observer.observe(el)
  })

  carouselInterval = setInterval(() => {
    currentImageIndex.value = (currentImageIndex.value + 1) % images.length
  }, 30000) // 30 seconds as requested
})

onUnmounted(() => {
  if (carouselInterval) clearInterval(carouselInterval)
})
</script>

<template>
  <div class="landing-page">
    <!-- Glowing Background Orbs -->
    <div class="glow-orb orb-1"></div>
    <div class="glow-orb orb-2"></div>
    <div class="glow-orb orb-3"></div>

    <!-- 1. Hero Section -->
    <section class="hero-section">
      <!-- Background Carousel -->
      <div class="hero-bg-carousel">
        <div 
          v-for="(img, idx) in images" 
          :key="img"
          class="hero-bg-slide"
          :style="{ backgroundImage: `url(${img})` }"
          :class="{ active: currentImageIndex === idx }"
        ></div>
        <div class="hero-bg-overlay"></div>
      </div>

      <div class="section-inner" style="position: relative; z-index: 2; max-width: 850px; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 24px; padding: 60px 20px;">
        <div class="hero-badge fade-in-up" style="animation-delay: 0.1s; background: rgba(255,255,255,0.15); border-color: rgba(255,255,255,0.3); color: white;">
          <Sparkles size="16" class="text-primary" />
          <span>Giải pháp Chuyển đổi số Tuyển dụng 2026</span>
        </div>
        
        <h1 class="hero-title fade-in-up" style="animation-delay: 0.2s; color: white;">
          Tương lai của Tuyển dụng <br />
          <span style="color: #c4b5fd;">Phỏng vấn thông minh cùng AI</span>
        </h1>
        
        <p class="hero-subtitle fade-in-up" style="animation-delay: 0.3s; color: rgba(255,255,255,0.9);">
          Trải nghiệm hệ thống tự động bóc tách CV, phân tích độ phù hợp với JD và phòng phỏng vấn trực tuyến tích hợp trí tuệ nhân tạo. Chấm dứt kỷ nguyên lọc hồ sơ thủ công.
        </p>

        <div class="hero-actions fade-in-up" style="animation-delay: 0.4s">
          <button class="btn-primary-glow" @click="handleGetStarted">
            Trải nghiệm ngay <ArrowRight size="18" style="margin-left: 8px" />
          </button>
          <button class="btn-glass" @click="handleGetStarted">
            Luyện tập phỏng vấn
          </button>
        </div>
      </div>
    </section>

    <!-- 2. Stats Section -->
    <section class="stats-section animate-on-scroll fade-in-up">
      <div class="section-inner stats-grid">
        <div class="stat-card glass-panel">
          <div class="stat-value gradient-text">98%</div>
          <div class="stat-label">Độ chính xác AI Parsing</div>
        </div>
        <div class="stat-card glass-panel">
          <div class="stat-value gradient-text">< 2s</div>
          <div class="stat-label">Tốc độ xử lý mỗi CV</div>
        </div>
        <div class="stat-card glass-panel">
          <div class="stat-value gradient-text">-80%</div>
          <div class="stat-label">Thời gian lọc hồ sơ</div>
        </div>
        <div class="stat-card glass-panel">
          <div class="stat-value gradient-text">24/7</div>
          <div class="stat-label">Luyện tập Mock Interview</div>
        </div>
      </div>
    </section>

    <!-- 3. How it Works (Workflow) -->
    <section class="workflow-section animate-on-scroll fade-in-up">
      <div class="section-inner">
        <div class="section-header">
          <h2 class="section-title">Vận hành xuyên suốt <span class="gradient-text">4 Bước</span></h2>
          <p class="section-subtitle">Quy trình tuyển dụng được tự động hóa từ khâu tiếp nhận đến khi ra quyết định.</p>
        </div>

        <div class="workflow-steps">
          <div class="step-card">
            <div class="step-number">1</div>
            <div class="step-icon"><FileText size="28" class="text-primary" /></div>
            <h3>Tải lên JD & CV</h3>
            <p>Hệ thống tự động số hóa và chuẩn hóa dữ liệu từ bất kỳ định dạng CV của ứng viên.</p>
          </div>
          <div class="step-connector"></div>
          <div class="step-card">
            <div class="step-number">2</div>
            <div class="step-icon"><BrainCircuit size="28" class="text-primary" /></div>
            <h3>AI Đối chiếu (Matching)</h3>
            <p>Thuật toán NLP đối chiếu chéo kỹ năng ứng viên với JD, đưa ra điểm số (Match Score).</p>
          </div>
          <div class="step-connector"></div>
          <div class="step-card">
            <div class="step-number">3</div>
            <div class="step-icon"><Video size="28" class="text-primary" /></div>
            <h3>Phỏng vấn LiveKit</h3>
            <p>Phòng họp ảo thời gian thực. AI tự động ghi âm, bóc băng và gợi ý câu hỏi liên quan.</p>
          </div>
          <div class="step-connector"></div>
          <div class="step-card">
            <div class="step-number">4</div>
            <div class="step-icon"><BarChart size="28" class="text-primary" /></div>
            <h3>Báo cáo Radar</h3>
            <p>Nhận báo cáo đa chiều phân tích điểm mạnh, yếu và gợi ý quyết định (Hire/Reject).</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 4. Target Audience (Lợi ích 2 chiều) -->
    <section class="audience-section animate-on-scroll fade-in-up">
      <div class="section-inner audience-grid">
        
        <!-- Recruiter Column -->
        <div class="audience-card glass-panel recruiter-card">
          <div class="audience-icon"><Users size="32" color="white" /></div>
          <h2 style="color: white; margin-bottom: 24px">Dành cho Doanh Nghiệp</h2>
          <ul class="benefit-list">
            <li><CheckCircle2 size="20" class="text-success" /> Lọc hàng ngàn CV chỉ trong vài giây.</li>
            <li><CheckCircle2 size="20" class="text-success" /> Xóa bỏ cảm tính trong tuyển dụng.</li>
            <li><CheckCircle2 size="20" class="text-success" /> Ngân hàng câu hỏi AI tự động gợi ý.</li>
            <li><CheckCircle2 size="20" class="text-success" /> Lưu trữ toàn bộ lịch sử phỏng vấn.</li>
          </ul>
        </div>

        <!-- Candidate Column -->
        <div class="audience-card glass-panel candidate-card">
          <div class="audience-icon"><Zap size="32" color="white" /></div>
          <h2 style="color: white; margin-bottom: 24px">Dành cho Ứng Viên</h2>
          <ul class="benefit-list">
            <li><CheckCircle2 size="20" class="text-primary" /> Biết ngay mức độ phù hợp của CV.</li>
            <li><CheckCircle2 size="20" class="text-primary" /> Luyện tập (Mock Interview) 24/7.</li>
            <li><CheckCircle2 size="20" class="text-primary" /> Nhận feedback chi tiết để cải thiện.</li>
            <li><CheckCircle2 size="20" class="text-primary" /> Phỏng vấn mượt mà không cài App.</li>
          </ul>
        </div>

      </div>
    </section>

    <!-- 5. CTA Footer -->
    <section class="cta-section animate-on-scroll fade-in-up">
      <div class="section-inner cta-content glass-panel" style="padding: 60px 40px; border-radius: 24px;">
        <h2>Sẵn sàng để thay đổi cách tuyển dụng?</h2>
        <p>Tham gia cùng hàng ngàn chuyên gia Nhân sự đang sử dụng Interview AI.</p>
        <button class="btn-primary-glow btn-large" @click="handleGetStarted">
          Bắt đầu sử dụng miễn phí ngay hôm nay
        </button>
      </div>
    </section>
  </div>
</template>

<style scoped>
.landing-page {
  position: relative;
  min-height: calc(100vh - 64px);
  overflow-x: hidden;
  padding: 0;
}

section {
  padding: 60px 20px;
  width: 100%;
}

.section-inner {
  max-width: 1200px;
  margin: 0 auto;
}

/* Section Headers */
.section-header {
  text-align: center;
  margin-bottom: 48px;
}

.section-title {
  font-size: 2.5rem;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 16px;
}

.section-subtitle {
  font-size: 1.1rem;
  color: var(--text-secondary);
  max-width: 600px;
  margin: 0 auto;
}

/* Glass Panel Base */
.glass-panel {
  background: var(--surface-soft);
  border: 1px solid var(--border);
  border-radius: 24px;
  backdrop-filter: blur(20px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.05);
}

/* Text & Colors */
.gradient-text {
  background: linear-gradient(135deg, var(--primary) 0%, #a855f7 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.text-primary { color: var(--primary); }
.text-success { color: #10b981; }

/* 1. Hero Section */
.hero-section {
  position: relative;
  width: 100%;
  overflow: hidden;
  margin-top: -64px; /* Pull up to cover the area under navbar if needed, or just let it sit */
}

.hero-bg-carousel {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.hero-bg-slide {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  opacity: 0;
  transition: opacity 2s ease-in-out, transform 15s linear;
  transform: scale(1);
}

.hero-bg-slide.active {
  opacity: 1;
  transform: scale(1.05);
}

.hero-bg-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.9) 0%, rgba(30, 58, 138, 0.7) 100%);
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 20px;
  background: rgba(79, 70, 229, 0.1);
  border: 1px solid rgba(79, 70, 229, 0.2);
  border-radius: 100px;
  color: var(--primary);
  font-weight: 600;
  font-size: 14px;
  backdrop-filter: blur(10px);
}

.hero-title {
  font-size: 4rem;
  line-height: 1.1;
  font-weight: 800;
  color: var(--text-main);
  margin: 0;
  letter-spacing: -1px;
}

.hero-subtitle {
  font-size: 1.25rem;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0;
  max-width: 700px;
}

.hero-actions {
  display: flex;
  gap: 16px;
  margin-top: 24px;
}

/* 2. Stats Section */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 24px;
}

.stat-card {
  padding: 32px 24px;
  text-align: center;
  transition: transform 0.3s ease;
}

.stat-card:hover { transform: translateY(-5px); }
.stat-value { font-size: 3rem; font-weight: 800; margin-bottom: 8px; }
.stat-label { font-size: 1rem; color: var(--text-secondary); font-weight: 500; }

/* 3. Workflow Section */
.workflow-steps {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  position: relative;
}

.step-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  position: relative;
  z-index: 2;
}

.step-number {
  width: 32px;
  height: 32px;
  background: var(--primary);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 14px;
  margin-bottom: 16px;
  box-shadow: 0 4px 10px rgba(79, 70, 229, 0.4);
}

.step-icon {
  width: 80px;
  height: 80px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  box-shadow: 0 10px 20px rgba(0,0,0,0.05);
  transition: transform 0.3s ease;
}

.step-card:hover .step-icon { transform: translateY(-5px) scale(1.05); }

.step-card h3 { font-size: 1.2rem; font-weight: 700; color: var(--text-main); margin-bottom: 12px; }
.step-card p { font-size: 0.95rem; color: var(--text-secondary); line-height: 1.5; padding: 0 10px; }

.step-connector {
  flex: 0.5;
  height: 2px;
  background: linear-gradient(90deg, var(--border) 0%, var(--primary) 50%, var(--border) 100%);
  margin-top: 56px;
  opacity: 0.5;
}

/* 4. Audience Section */
.audience-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
}

.audience-card {
  padding: 48px;
  position: relative;
  overflow: hidden;
}

.recruiter-card {
  background: linear-gradient(135deg, rgba(79, 70, 229, 0.8) 0%, rgba(168, 85, 247, 0.8) 100%);
  border: none;
}

.candidate-card {
  background: rgba(30, 41, 59, 0.7);
}

.audience-icon {
  width: 64px;
  height: 64px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  backdrop-filter: blur(10px);
}

.benefit-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 16px; }
.benefit-list li { display: flex; align-items: flex-start; gap: 12px; font-size: 1.05rem; font-weight: 500; }
.recruiter-card .benefit-list li { color: rgba(255, 255, 255, 0.9); }
.candidate-card .benefit-list li { color: var(--text-main); }

/* 5. CTA Section */
.cta-section {
  padding: 80px 20px;
  background: linear-gradient(180deg, var(--background) 0%, rgba(79, 70, 229, 0.05) 100%);
  border-bottom: 4px solid var(--primary);
}

.cta-content { text-align: center; }
.cta-content h2 { font-size: 2.5rem; font-weight: 800; color: var(--text-main); margin-bottom: 16px; }
.cta-content p { font-size: 1.2rem; color: var(--text-secondary); margin-bottom: 40px; }

/* Buttons */
.btn-primary-glow {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 14px 28px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--primary) 0%, #6366f1 100%);
  color: white;
  font-weight: 600;
  font-size: 16px;
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 10px 20px -10px var(--primary);
}

.btn-primary-glow:hover { transform: translateY(-2px); box-shadow: 0 15px 25px -10px var(--primary); }
.btn-large { padding: 18px 40px; font-size: 1.1rem; border-radius: 16px; }

.btn-glass {
  padding: 14px 28px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-main);
  font-weight: 600;
  font-size: 16px;
  border: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.btn-glass:hover { background: rgba(255, 255, 255, 0.1); transform: translateY(-2px); }

/* Animations */
@keyframes float {
  0% { transform: translate(0, 0) scale(1); }
  100% { transform: translate(50px, -50px) scale(1.1); }
}

.fade-in-up { opacity: 0; transform: translateY(30px); transition: opacity 0.8s ease, transform 0.8s ease; }
.fade-in-up.visible { opacity: 1; transform: translateY(0); }

/* Keyframe for elements immediately visible without scroll */
.hero-section .fade-in-up {
  animation: fadeInUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeInUp { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }

/* Responsive */
@media (max-width: 992px) {
  .workflow-steps { flex-direction: column; gap: 40px; }
  .step-connector { display: none; }
  .step-card { flex-direction: row; text-align: left; align-items: flex-start; gap: 20px; }
  .step-number { position: absolute; left: 15px; top: -15px; }
  .step-icon { width: 60px; height: 60px; margin-bottom: 0; }
  
  .audience-grid { grid-template-columns: 1fr; }
}

@media (max-width: 768px) {
  .hero-title { font-size: 2.5rem; }
  .hero-actions { flex-direction: column; width: 100%; }
  .hero-actions button { width: 100%; }
  .audience-card { padding: 32px; }
}
</style>
