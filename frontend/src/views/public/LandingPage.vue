<script setup>
import { useRouter } from 'vue-router'
import { Sparkles, BrainCircuit, Video, FileText, BarChart3, CheckCircle2, ArrowRight, ShieldCheck, Bot, ChevronDown, Star } from 'lucide-vue-next'
import { ref, onMounted, onUnmounted } from 'vue'

const router = useRouter()
const handleGetStarted = () => router.push('/login')

const currentImageIndex = ref(0)
const images = ['/images/hero_office.png', '/images/hero_network.png', '/images/hero_dashboard.png']
let carouselInterval

// 3D mouse parallax/tilt state
const tiltX = ref(0)
const tiltY = ref(0)
const parX = ref(0)
const parY = ref(0)
const heroRef = ref(null)
let reduceMotion = false

const onHeroMove = (e) => {
  if (reduceMotion || !heroRef.value) return
  const r = heroRef.value.getBoundingClientRect()
  const px = (e.clientX - r.left) / r.width - 0.5   // -0.5 .. 0.5
  const py = (e.clientY - r.top) / r.height - 0.5
  tiltY.value = px * 16      // rotateY
  tiltX.value = -py * 16     // rotateX
  parX.value = px * 40       // orb parallax
  parY.value = py * 40
}
const onHeroLeave = () => { tiltX.value = 0; tiltY.value = 0; parX.value = 0; parY.value = 0 }

// FAQ accordion
const openFaq = ref(0)
const faqs = [
  { q: 'Hệ thống bóc tách CV chính xác đến đâu?', a: 'AI đọc được hầu hết định dạng PDF/DOCX, trích xuất kỹ năng, kinh nghiệm và học vấn với độ chính xác ~98% trên dữ liệu tiếng Việt và tiếng Anh.' },
  { q: 'Ứng viên có cần cài phần mềm gì không?', a: 'Không. Phòng phỏng vấn chạy trực tiếp trên trình duyệt qua LiveKit — chỉ cần micro và camera, không cài đặt thêm.' },
  { q: 'Dữ liệu phỏng vấn được bảo mật thế nào?', a: 'Toàn bộ bản ghi, transcript và điểm đánh giá được lưu riêng theo từng doanh nghiệp, phân quyền chặt chẽ. Ứng viên không thấy điểm nội bộ của nhà tuyển dụng.' },
  { q: 'Tôi có thể luyện phỏng vấn miễn phí không?', a: 'Có. Ứng viên luyện tập với AI Interviewer 24/7 không giới hạn số lần và nhận feedback chi tiết sau mỗi buổi.' },
]

onMounted(() => {
  reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

  // Stagger children reveal via CSS custom prop index (đặt trước khi observe)
  document.querySelectorAll('[data-stagger]').forEach((container) => {
    Array.from(container.children).forEach((child, i) => {
      child.style.setProperty('--i', i)
    })
  })

  // Reveal on scroll — theo dõi cả section và lưới con, cuộn lên/xuống đều hiện lại
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) entry.target.classList.add('visible')
      else entry.target.classList.remove('visible')
    })
  }, { threshold: 0.12, rootMargin: '0px 0px -6% 0px' })
  document.querySelectorAll('.animate-on-scroll, .reveal-children').forEach((el) => observer.observe(el))

  carouselInterval = setInterval(() => {
    currentImageIndex.value = (currentImageIndex.value + 1) % images.length
  }, 8000)
})

onUnmounted(() => { if (carouselInterval) clearInterval(carouselInterval) })
</script>

<template>
  <div class="landing-page">
    <!-- 1. HERO with 3D mouse parallax -->
    <section
      ref="heroRef"
      class="hero-section"
      @mousemove="onHeroMove"
      @mouseleave="onHeroLeave"
    >
      <!-- Background carousel -->
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

      <!-- Parallax glow orbs (move opposite to mouse) -->
      <div class="glow-orb orb-1" :style="{ transform: `translate3d(${-parX}px, ${-parY}px, 0)` }"></div>
      <div class="glow-orb orb-2" :style="{ transform: `translate3d(${parX}px, ${parY}px, 0)` }"></div>
      <div class="glow-orb orb-3" :style="{ transform: `translate3d(${parX * 0.5}px, ${-parY}px, 0)` }"></div>

      <div class="hero-grid">
        <!-- Left: copy -->
        <div class="hero-copy">
          <div class="hero-badge fade-in-up" style="animation-delay: 0.05s">
            <Sparkles :size="16" />
            <span>Giải pháp Chuyển đổi số Tuyển dụng 2026</span>
          </div>

          <h1 class="hero-title fade-in-up" style="animation-delay: 0.15s">
            Tương lai của Tuyển dụng
            <span class="hero-title-accent">Phỏng vấn thông minh cùng AI</span>
          </h1>

          <p class="hero-subtitle fade-in-up" style="animation-delay: 0.25s">
            Tự động bóc tách CV, chấm độ phù hợp với JD và phỏng vấn trực tuyến tích hợp trí tuệ nhân tạo. Chấm dứt kỷ nguyên lọc hồ sơ thủ công.
          </p>

          <div class="hero-actions fade-in-up" style="animation-delay: 0.35s">
            <button class="btn-primary-glow sheen" @click="handleGetStarted">
              Trải nghiệm ngay <ArrowRight :size="18" style="margin-left: 8px" />
            </button>
            <button class="btn-glass" @click="handleGetStarted">Luyện tập phỏng vấn</button>
          </div>

          <div class="hero-trust fade-in-up" style="animation-delay: 0.45s">
            <div class="trust-avatars">
              <span class="t-av">A</span><span class="t-av">N</span><span class="t-av">H</span><span class="t-av">+</span>
            </div>
            <div class="trust-text">
              <div class="trust-stars"><Star :size="14" fill="currentColor" /><Star :size="14" fill="currentColor" /><Star :size="14" fill="currentColor" /><Star :size="14" fill="currentColor" /><Star :size="14" fill="currentColor" /></div>
              <span>Được tin dùng bởi hàng ngàn chuyên gia HR</span>
            </div>
          </div>
        </div>

        <!-- Right: 3D tilting visual card -->
        <div class="hero-visual fade-in-up" style="animation-delay: 0.3s">
          <div
            class="tilt-scene"
            :style="{ transform: `perspective(1000px) rotateX(${tiltX}deg) rotateY(${tiltY}deg)` }"
          >
            <!-- Floating layers -->
            <div class="float-card layer-back" style="transform: translateZ(20px);">
              <div class="mini-head"><Bot :size="18" /><span>AI Interviewer</span><span class="live-pill"><span class="live-dot"></span>Live</span></div>
              <div class="mini-line w-80"></div>
              <div class="mini-line w-60"></div>
              <div class="mini-line w-40"></div>
            </div>

            <div class="float-card layer-front" style="transform: translateZ(70px);">
              <div class="score-head">
                <span>Match Score</span>
                <span class="score-badge">92%</span>
              </div>
              <div class="bar"><i style="width: 92%"></i></div>
              <div class="score-rows">
                <div class="score-row"><span>Kỹ năng</span><b>Xuất sắc</b></div>
                <div class="score-row"><span>Kinh nghiệm</span><b>Phù hợp</b></div>
                <div class="score-row"><span>Văn hoá</span><b>Tốt</b></div>
              </div>
            </div>

            <div class="float-chip chip-cv" style="transform: translateZ(110px);">
              <FileText :size="16" /> CV.pdf đã bóc tách
            </div>
            <div class="float-chip chip-radar" style="transform: translateZ(120px);">
              <BarChart3 :size="16" /> Báo cáo Radar
            </div>
          </div>
        </div>
      </div>

      <div class="hero-scroll-hint"><ChevronDown :size="22" /></div>
    </section>

    <!-- 2. STATS -->
    <section class="stats-section animate-on-scroll fade-in-up">
      <div class="section-inner stats-grid reveal-children" data-stagger>
        <div class="stat-card"><div class="stat-value gradient-text">98%</div><div class="stat-label">Độ chính xác AI Parsing</div></div>
        <div class="stat-card"><div class="stat-value gradient-text">&lt; 2s</div><div class="stat-label">Tốc độ xử lý mỗi CV</div></div>
        <div class="stat-card"><div class="stat-value gradient-text">-80%</div><div class="stat-label">Thời gian lọc hồ sơ</div></div>
        <div class="stat-card"><div class="stat-value gradient-text">24/7</div><div class="stat-label">Luyện tập Mock Interview</div></div>
      </div>
    </section>

    <!-- 3. FEATURES grid (thông tin hơn) -->
    <section class="features-section animate-on-scroll fade-in-up">
      <div class="section-inner">
        <div class="section-header">
          <h2 class="section-title">Tất cả trong <span class="gradient-text">một nền tảng</span></h2>
          <p class="section-subtitle">Từ sàng lọc hồ sơ đến quyết định tuyển dụng — mọi công đoạn đều được AI hỗ trợ.</p>
        </div>
        <div class="feature-grid reveal-children" data-stagger>
          <div class="feature-card tilt-3d">
            <div class="feature-icon"><FileText :size="24" /></div>
            <h3>Bóc tách CV tự động</h3>
            <p>Số hóa mọi định dạng CV, trích xuất kỹ năng &amp; kinh nghiệm chỉ trong vài giây.</p>
          </div>
          <div class="feature-card tilt-3d">
            <div class="feature-icon is-accent"><BrainCircuit :size="24" /></div>
            <h3>AI đối chiếu JD</h3>
            <p>Thuật toán NLP chấm điểm độ phù hợp giữa ứng viên và mô tả công việc.</p>
          </div>
          <div class="feature-card tilt-3d">
            <div class="feature-icon"><Video :size="24" /></div>
            <h3>Phỏng vấn LiveKit</h3>
            <p>Phòng họp ảo thời gian thực, tự động ghi âm và bóc băng transcript.</p>
          </div>
          <div class="feature-card tilt-3d">
            <div class="feature-icon is-accent"><Bot :size="24" /></div>
            <h3>AI Interviewer 24/7</h3>
            <p>Luyện phỏng vấn thử không giới hạn, nhận feedback chi tiết ngay lập tức.</p>
          </div>
          <div class="feature-card tilt-3d">
            <div class="feature-icon"><BarChart3 :size="24" /></div>
            <h3>Báo cáo Radar đa chiều</h3>
            <p>Phân tích điểm mạnh, điểm yếu và gợi ý quyết định Hire / Reject.</p>
          </div>
          <div class="feature-card tilt-3d">
            <div class="feature-icon is-accent"><ShieldCheck :size="24" /></div>
            <h3>Bảo mật &amp; phân quyền</h3>
            <p>Dữ liệu tách biệt theo doanh nghiệp, phân quyền chặt chẽ theo vai trò.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 4. WORKFLOW -->
    <section class="workflow-section animate-on-scroll fade-in-up">
      <div class="section-inner">
        <div class="section-header">
          <h2 class="section-title">Vận hành xuyên suốt <span class="gradient-text">4 Bước</span></h2>
          <p class="section-subtitle">Quy trình tuyển dụng được tự động hóa từ khâu tiếp nhận đến khi ra quyết định.</p>
        </div>
        <div class="workflow-steps">
          <div class="step-card">
            <div class="step-number">1</div>
            <div class="step-icon"><FileText :size="28" /></div>
            <h3>Tải lên JD &amp; CV</h3>
            <p>Hệ thống tự động số hóa và chuẩn hóa dữ liệu từ mọi định dạng CV.</p>
          </div>
          <div class="step-connector"></div>
          <div class="step-card">
            <div class="step-number">2</div>
            <div class="step-icon"><BrainCircuit :size="28" /></div>
            <h3>AI Đối chiếu</h3>
            <p>NLP đối chiếu kỹ năng ứng viên với JD, đưa ra Match Score.</p>
          </div>
          <div class="step-connector"></div>
          <div class="step-card">
            <div class="step-number">3</div>
            <div class="step-icon"><Video :size="28" /></div>
            <h3>Phỏng vấn LiveKit</h3>
            <p>Phòng ảo thời gian thực, AI ghi âm, bóc băng và gợi ý câu hỏi.</p>
          </div>
          <div class="step-connector"></div>
          <div class="step-card">
            <div class="step-number">4</div>
            <div class="step-icon"><BarChart3 :size="28" /></div>
            <h3>Báo cáo Radar</h3>
            <p>Báo cáo đa chiều về điểm mạnh, yếu và gợi ý quyết định.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 5. LỢI ÍCH CHO ỨNG VIÊN -->
    <section class="audience-section animate-on-scroll fade-in-up">
      <div class="section-inner">
        <div class="section-header">
          <h2 class="section-title">Đồng hành cùng bạn <span class="gradient-text">trước &amp; sau</span> phỏng vấn</h2>
          <p class="section-subtitle">Từ lúc chuẩn bị hồ sơ đến khi bước vào buổi phỏng vấn thật, bạn luôn có AI hỗ trợ.</p>
        </div>
        <div class="audience-grid">
          <div class="audience-card recruiter-card">
            <div class="audience-icon"><FileText :size="30" color="white" /></div>
            <h2>Chuẩn bị hồ sơ</h2>
            <ul class="benefit-list">
              <li><CheckCircle2 :size="20" /> Tải CV lên, AI bóc tách kỹ năng &amp; kinh nghiệm.</li>
              <li><CheckCircle2 :size="20" /> Biết ngay mức độ phù hợp với từng tin tuyển dụng.</li>
              <li><CheckCircle2 :size="20" /> Gợi ý điểm cần bổ sung để hồ sơ nổi bật hơn.</li>
              <li><CheckCircle2 :size="20" /> Ứng tuyển nhanh, theo dõi trạng thái mọi lúc.</li>
            </ul>
          </div>
          <div class="audience-card candidate-card">
            <div class="audience-icon is-light"><Bot :size="30" /></div>
            <h2>Luyện tập &amp; tự tin</h2>
            <ul class="benefit-list">
              <li><CheckCircle2 :size="20" /> Luyện phỏng vấn với AI Interviewer 24/7.</li>
              <li><CheckCircle2 :size="20" /> Câu hỏi bám sát vị trí và cấp độ bạn chọn.</li>
              <li><CheckCircle2 :size="20" /> Nhận feedback chi tiết cho từng câu trả lời.</li>
              <li><CheckCircle2 :size="20" /> Phỏng vấn ngay trên trình duyệt, không cài App.</li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- 6. MẸO PHỎNG VẤN -->
    <section class="testi-section animate-on-scroll fade-in-up">
      <div class="section-inner">
        <div class="section-header">
          <h2 class="section-title">Mẹo giúp bạn <span class="gradient-text">ghi điểm</span></h2>
          <p class="section-subtitle">Vài nguyên tắc đơn giản giúp buổi phỏng vấn của bạn thuyết phục hơn.</p>
        </div>
        <div class="testi-grid reveal-children" data-stagger>
          <div class="testi-card tilt-3d">
            <div class="tip-num">01</div>
            <h3 class="tip-title">Trả lời theo cấu trúc STAR</h3>
            <p>Với câu hỏi tình huống, hãy nêu rõ Situation – Task – Action – Result để câu trả lời mạch lạc và có kết quả cụ thể.</p>
          </div>
          <div class="testi-card tilt-3d">
            <div class="tip-num">02</div>
            <h3 class="tip-title">Gắn câu trả lời với JD</h3>
            <p>Đọc kỹ mô tả công việc và dẫn chứng đúng kỹ năng họ cần. Match Score của bạn sẽ cho biết nên nhấn mạnh điều gì.</p>
          </div>
          <div class="testi-card tilt-3d">
            <div class="tip-num">03</div>
            <h3 class="tip-title">Luyện trước khi phỏng vấn thật</h3>
            <p>Chạy vài buổi Mock Interview với AI, đọc kỹ feedback và cải thiện điểm yếu trước khi bước vào buổi phỏng vấn chính thức.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 7. FAQ -->
    <section class="faq-section animate-on-scroll fade-in-up">
      <div class="section-inner faq-inner">
        <div class="section-header">
          <h2 class="section-title">Câu hỏi <span class="gradient-text">thường gặp</span></h2>
        </div>
        <div class="faq-list">
          <div v-for="(f, i) in faqs" :key="i" class="faq-item" :class="{ open: openFaq === i }">
            <button class="faq-q" @click="openFaq = openFaq === i ? -1 : i">
              <span>{{ f.q }}</span>
              <ChevronDown :size="20" class="faq-chevron" />
            </button>
            <div class="faq-a"><p>{{ f.a }}</p></div>
          </div>
        </div>
      </div>
    </section>

    <!-- 8. CTA -->
    <section class="cta-section animate-on-scroll fade-in-up">
      <div class="section-inner">
        <div class="cta-content brand-banner">
          <h2>Sẵn sàng để thay đổi cách tuyển dụng?</h2>
          <p>Tham gia cùng hàng ngàn chuyên gia Nhân sự đang sử dụng Interview AI.</p>
          <button class="cta-btn sheen" @click="handleGetStarted">
            Bắt đầu miễn phí ngay hôm nay <ArrowRight :size="18" />
          </button>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.landing-page { position: relative; min-height: calc(100vh - 64px); overflow-x: clip; font-family: var(--sans); }
section { padding: 80px 20px; width: 100%; }
.section-inner { max-width: 1200px; margin: 0 auto; }

.section-header { text-align: center; margin-bottom: 52px; }
.section-title { font-size: 2.3rem; font-weight: 800; color: var(--text-main); margin-bottom: 14px; letter-spacing: -0.02em; }
.section-subtitle { font-size: 1.05rem; color: var(--text-secondary); max-width: 600px; margin: 0 auto; line-height: 1.6; }

.gradient-text { background: var(--gradient-brand); -webkit-background-clip: text; background-clip: text; -webkit-text-fill-color: transparent; }

/* ===== 1. HERO ===== */
.hero-section { position: relative; overflow: hidden; margin-top: -64px; padding: 150px 20px 90px; min-height: 92vh; display: flex; align-items: center; }
.hero-bg-carousel { position: absolute; inset: 0; z-index: 0; }
.hero-bg-slide { position: absolute; inset: 0; background-size: cover; background-position: center; opacity: 0; transition: opacity 2s ease-in-out, transform 9s ease-out; transform: scale(1); }
.hero-bg-slide.active { opacity: 1; transform: scale(1.08); }
.hero-bg-overlay { position: absolute; inset: 0; background: linear-gradient(120deg, rgba(10,20,40,0.94) 0%, rgba(12,35,70,0.82) 55%, rgba(8,60,90,0.72) 100%); }

.glow-orb { position: absolute; border-radius: 50%; filter: blur(70px); pointer-events: none; z-index: 1; transition: transform 0.35s cubic-bezier(0.2,0.8,0.2,1); }
.orb-1 { width: 360px; height: 360px; background: rgba(37,99,235,0.45); top: -60px; left: -40px; }
.orb-2 { width: 300px; height: 300px; background: rgba(8,145,178,0.4); bottom: -60px; right: 0; }
.orb-3 { width: 240px; height: 240px; background: rgba(59,130,246,0.35); top: 40%; left: 55%; }

.hero-grid { position: relative; z-index: 2; max-width: 1200px; margin: 0 auto; width: 100%; display: grid; grid-template-columns: 1.05fr 0.95fr; gap: 48px; align-items: center; }
.hero-copy { display: flex; flex-direction: column; align-items: flex-start; gap: 22px; }
.hero-badge { display: inline-flex; align-items: center; gap: 8px; padding: 8px 18px; background: rgba(255,255,255,0.14); border: 1px solid rgba(255,255,255,0.28); border-radius: var(--radius-full); color: #fff; font-weight: 600; font-size: 13.5px; backdrop-filter: blur(10px); }
.hero-badge :deep(svg) { color: #67e8f9; }
.hero-title { font-size: 3.6rem; line-height: 1.08; font-weight: 800; color: #fff; margin: 0; letter-spacing: -1.5px; }
.hero-title-accent { display: block; margin-top: 6px; background: linear-gradient(120deg, #60a5fa, #22d3ee); -webkit-background-clip: text; background-clip: text; -webkit-text-fill-color: transparent; }
.hero-subtitle { font-size: 1.18rem; color: rgba(255,255,255,0.9); line-height: 1.65; margin: 0; max-width: 560px; }
.hero-actions { display: flex; gap: 14px; margin-top: 6px; flex-wrap: wrap; }
.hero-trust { display: flex; align-items: center; gap: 14px; margin-top: 14px; }
.trust-avatars { display: flex; }
.t-av { width: 34px; height: 34px; border-radius: 50%; background: var(--gradient-brand); color: #fff; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 13px; border: 2px solid rgba(255,255,255,0.5); margin-left: -8px; }
.t-av:first-child { margin-left: 0; }
.trust-text { display: flex; flex-direction: column; gap: 2px; }
.trust-stars { display: flex; gap: 2px; color: #fbbf24; }
.trust-text span { font-size: 13px; color: rgba(255,255,255,0.85); }

/* 3D tilt visual */
.hero-visual { perspective: 1000px; display: flex; justify-content: center; }
.tilt-scene { position: relative; width: 380px; height: 360px; transform-style: preserve-3d; transition: transform 0.2s ease-out; }
.float-card { position: absolute; background: rgba(255,255,255,0.96); border: 1px solid rgba(255,255,255,0.6); border-radius: var(--radius-lg); box-shadow: 0 30px 60px -20px rgba(0,0,0,0.5); padding: 20px; }
.layer-back { top: 0; left: 0; width: 300px; }
.layer-front { bottom: 0; right: 0; width: 290px; }
.mini-head { display: flex; align-items: center; gap: 8px; font-weight: 700; color: var(--text-main); font-size: 14px; margin-bottom: 16px; }
.mini-head :deep(svg) { color: var(--primary); }
.live-pill { margin-left: auto; display: inline-flex; align-items: center; gap: 5px; font-size: 11px; font-weight: 700; color: var(--success); background: rgba(22,163,74,0.12); padding: 3px 9px; border-radius: var(--radius-full); }
.live-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--success); animation: blink 1.6s ease-in-out infinite; }
.mini-line { height: 9px; border-radius: 999px; background: var(--surface-soft); margin-bottom: 10px; }
.w-80 { width: 80%; } .w-60 { width: 60%; } .w-40 { width: 40%; }
.score-head { display: flex; justify-content: space-between; align-items: center; font-weight: 700; color: var(--text-main); margin-bottom: 12px; }
.score-badge { background: var(--gradient-brand); color: #fff; padding: 4px 12px; border-radius: var(--radius-full); font-size: 14px; }
.bar { height: 8px; border-radius: 999px; background: var(--surface-soft); overflow: hidden; margin-bottom: 16px; }
.bar i { display: block; height: 100%; background: var(--gradient-brand); border-radius: 999px; }
.score-rows { display: flex; flex-direction: column; gap: 9px; }
.score-row { display: flex; justify-content: space-between; font-size: 13px; color: var(--text-secondary); }
.score-row b { color: var(--text-main); }
.float-chip { position: absolute; display: inline-flex; align-items: center; gap: 7px; background: #fff; border: 1px solid var(--border); border-radius: var(--radius-full); padding: 9px 15px; font-size: 12.5px; font-weight: 600; color: var(--text-main); box-shadow: 0 16px 30px -12px rgba(0,0,0,0.4); }
.float-chip :deep(svg) { color: var(--primary); }
.chip-cv { top: -18px; right: 10px; animation: floaty 4.5s ease-in-out infinite; }
.chip-radar { bottom: 40px; left: -26px; animation: floaty 5.5s ease-in-out infinite reverse; }

.hero-scroll-hint { position: absolute; bottom: 22px; left: 50%; transform: translateX(-50%); z-index: 2; color: rgba(255,255,255,0.7); animation: bob 2s ease-in-out infinite; }

@keyframes floaty { 0%,100% { transform: translateZ(110px) translateY(0); } 50% { transform: translateZ(110px) translateY(-10px); } }
@keyframes bob { 0%,100% { transform: translate(-50%, 0); } 50% { transform: translate(-50%, 8px); } }
@keyframes blink { 0%,100% { opacity: 1; } 50% { opacity: 0.3; } }

/* Buttons */
.btn-primary-glow { display: inline-flex; align-items: center; justify-content: center; padding: 14px 26px; border-radius: var(--radius); background: var(--gradient-brand); color: #fff; font-weight: 700; font-size: 15px; border: none; cursor: pointer; transition: transform 0.25s ease, box-shadow 0.25s ease; box-shadow: var(--shadow-glow); }
.btn-primary-glow:hover { transform: translateY(-2px); box-shadow: 0 14px 30px -8px rgba(37,99,235,0.5); }
.btn-glass { padding: 14px 26px; border-radius: var(--radius); background: rgba(255,255,255,0.1); color: #fff; font-weight: 600; font-size: 15px; border: 1px solid rgba(255,255,255,0.35); cursor: pointer; transition: all 0.25s ease; backdrop-filter: blur(10px); }
.btn-glass:hover { background: rgba(255,255,255,0.2); transform: translateY(-2px); }

/* ===== 2. STATS ===== */
.stats-section { background: var(--surface); }
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 24px; }
.stat-card { padding: 30px 24px; text-align: center; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); transition: transform 0.3s ease, box-shadow 0.3s ease; }
.stat-card:hover { transform: translateY(-5px); box-shadow: var(--shadow-lg); }
.stat-value { font-size: 2.8rem; font-weight: 800; margin-bottom: 6px; letter-spacing: -0.02em; }
.stat-label { font-size: 0.95rem; color: var(--text-secondary); font-weight: 500; }

/* ===== 3. FEATURES ===== */
.features-section { background: var(--background); }
.feature-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 24px; }
.feature-card { padding: 28px; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); }
.feature-icon { width: 52px; height: 52px; border-radius: var(--radius); background: var(--primary-light); color: var(--primary); display: flex; align-items: center; justify-content: center; margin-bottom: 18px; transition: transform 0.3s ease; }
.feature-icon.is-accent { background: var(--accent-bg); color: var(--accent); }
.feature-card:hover .feature-icon { transform: scale(1.08); }
.feature-card h3 { font-size: 1.15rem; font-weight: 700; color: var(--text-main); margin-bottom: 8px; }
.feature-card p { font-size: 0.95rem; color: var(--text-secondary); line-height: 1.55; }

/* ===== 4. WORKFLOW ===== */
.workflow-section { background: var(--surface); }
.workflow-steps { display: flex; justify-content: space-between; align-items: flex-start; }
.step-card { flex: 1; display: flex; flex-direction: column; align-items: center; text-align: center; position: relative; z-index: 2; }
.step-number { width: 30px; height: 30px; background: var(--primary); color: #fff; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: 800; font-size: 13px; margin-bottom: 14px; box-shadow: 0 4px 10px rgba(37,99,235,0.35); }
.step-icon { width: 74px; height: 74px; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); display: flex; align-items: center; justify-content: center; margin-bottom: 18px; box-shadow: var(--shadow-sm); transition: transform 0.3s ease; color: var(--primary); }
.step-card:hover .step-icon { transform: translateY(-5px) scale(1.05); box-shadow: var(--shadow-md); }
.step-card h3 { font-size: 1.1rem; font-weight: 700; color: var(--text-main); margin-bottom: 10px; }
.step-card p { font-size: 0.9rem; color: var(--text-secondary); line-height: 1.5; padding: 0 8px; }
.step-connector { flex: 0.5; height: 2px; background: linear-gradient(90deg, var(--border), var(--primary), var(--border)); margin-top: 52px; opacity: 0.5; }

/* ===== 5. AUDIENCE ===== */
.audience-section { background: var(--background); }
.audience-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 28px; }
.audience-card { padding: 44px; border-radius: var(--radius-lg); position: relative; overflow: hidden; }
.recruiter-card { background: var(--gradient-brand); box-shadow: var(--shadow-glow); }
.recruiter-card h2 { color: #fff; }
.candidate-card { background: var(--surface); border: 1px solid var(--border); box-shadow: var(--shadow-sm); }
.candidate-card h2 { color: var(--text-main); }
.audience-card h2 { font-size: 1.5rem; font-weight: 800; margin-bottom: 22px; }
.audience-icon { width: 60px; height: 60px; background: rgba(255,255,255,0.2); border-radius: var(--radius); display: flex; align-items: center; justify-content: center; margin-bottom: 22px; }
.audience-icon.is-light { background: var(--primary-light); color: var(--primary); }
.benefit-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 14px; }
.benefit-list li { display: flex; align-items: flex-start; gap: 12px; font-size: 1rem; font-weight: 500; }
.recruiter-card .benefit-list li { color: rgba(255,255,255,0.95); }
.recruiter-card .benefit-list :deep(svg) { color: #a7f3d0; flex-shrink: 0; }
.candidate-card .benefit-list li { color: var(--text-secondary); }
.candidate-card .benefit-list :deep(svg) { color: var(--primary); flex-shrink: 0; }

/* ===== 6. TESTIMONIALS ===== */
.testi-section { background: var(--surface); }
.testi-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 24px; }
.testi-card { padding: 28px; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); }
.quote-icon { color: var(--primary); opacity: 0.5; margin-bottom: 12px; }
.tip-num { font-size: 1.6rem; font-weight: 800; background: var(--gradient-brand); -webkit-background-clip: text; background-clip: text; -webkit-text-fill-color: transparent; margin-bottom: 8px; }
.tip-title { font-size: 1.1rem; font-weight: 700; color: var(--text-main); margin-bottom: 8px; }
.testi-card p { font-size: 0.98rem; color: var(--text-secondary); line-height: 1.6; margin-bottom: 0; }
.testi-author { display: flex; align-items: center; gap: 12px; }
.testi-author .t-av { margin: 0; border: none; width: 40px; height: 40px; }
.testi-author b { display: block; font-size: 14px; color: var(--text-main); }
.testi-author small { font-size: 12px; color: var(--text-muted); }

/* ===== 7. FAQ ===== */
.faq-section { background: var(--background); }
.faq-inner { max-width: 760px; }
.faq-list { display: flex; flex-direction: column; gap: 12px; }
.faq-item { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); overflow: hidden; transition: border-color 0.2s ease; }
.faq-item.open { border-color: var(--primary-light); }
.faq-q { width: 100%; display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 18px 22px; background: transparent; border: none; cursor: pointer; font-size: 1rem; font-weight: 600; color: var(--text-main); text-align: left; }
.faq-chevron { color: var(--text-muted); transition: transform 0.3s ease; flex-shrink: 0; }
.faq-item.open .faq-chevron { transform: rotate(180deg); color: var(--primary); }
.faq-a { max-height: 0; overflow: hidden; transition: max-height 0.35s ease; }
.faq-item.open .faq-a { max-height: 220px; }
.faq-a p { padding: 0 22px 20px; color: var(--text-secondary); line-height: 1.6; margin: 0; }

/* ===== 8. CTA ===== */
.cta-section { background: var(--surface); }
.cta-content { text-align: center; padding: 60px 40px; }
.cta-content h2 { font-size: 2.2rem; font-weight: 800; color: #fff; margin-bottom: 14px; }
.cta-content p { font-size: 1.1rem; color: rgba(255,255,255,0.9); margin-bottom: 32px; }
.cta-btn { display: inline-flex; align-items: center; gap: 8px; padding: 16px 36px; border-radius: var(--radius); background: #fff; color: var(--primary); font-weight: 700; font-size: 1.05rem; border: none; cursor: pointer; transition: transform 0.25s ease; }
.cta-btn:hover { transform: translateY(-2px); }

/* ===== REVEAL ON SCROLL (nhẹ nhàng, 2 chiều) ===== */
.fade-in-up { opacity: 0; transform: translateY(24px); transition: opacity 0.7s cubic-bezier(0.16,1,0.3,1), transform 0.7s cubic-bezier(0.16,1,0.3,1); }
.fade-in-up.visible { opacity: 1; transform: translateY(0); }
.hero-section .fade-in-up { animation: fadeInUp 0.8s cubic-bezier(0.16,1,0.3,1) forwards; opacity: 0; }
@keyframes fadeInUp { to { opacity: 1; transform: translateY(0); } }

/* Stagger children when parent becomes visible */
.reveal-children > * { opacity: 0; transform: translateY(24px); transition: opacity 0.6s cubic-bezier(0.16,1,0.3,1), transform 0.6s cubic-bezier(0.16,1,0.3,1); transition-delay: calc(var(--i, 0) * 80ms); }
.reveal-children.visible > * { opacity: 1; transform: translateY(0); }

/* 3D tilt on cards (dùng class chung từ global.css: .tilt-3d) */

/* ===== RESPONSIVE ===== */
@media (max-width: 992px) {
  .hero-grid { grid-template-columns: 1fr; gap: 40px; text-align: center; }
  .hero-copy { align-items: center; }
  .hero-visual { margin-top: 10px; }
  .workflow-steps { flex-direction: column; gap: 36px; }
  .step-connector { display: none; }
  .step-card { flex-direction: row; text-align: left; align-items: flex-start; gap: 18px; }
  .step-icon { margin-bottom: 0; }
  .audience-grid { grid-template-columns: 1fr; }
}
@media (max-width: 768px) {
  .hero-section { padding: 120px 20px 70px; min-height: auto; }
  .hero-title { font-size: 2.5rem; }
  .hero-actions { flex-direction: column; width: 100%; }
  .hero-actions button { width: 100%; }
  .tilt-scene { width: 300px; height: 320px; transform: none !important; }
  .audience-card { padding: 32px; }
  .section-title { font-size: 1.9rem; }
}

/* Reduce motion */
@media (prefers-reduced-motion: reduce) {
  .hero-bg-slide, .glow-orb, .chip-cv, .chip-radar, .hero-scroll-hint, .live-dot, .tilt-scene { animation: none !important; transition: none !important; }
  .tilt-scene { transform: none !important; }
  .fade-in-up, .reveal-children > * { transition: none !important; opacity: 1 !important; transform: none !important; }
  .hero-section .fade-in-up { animation: none !important; opacity: 1 !important; }
}
</style>
