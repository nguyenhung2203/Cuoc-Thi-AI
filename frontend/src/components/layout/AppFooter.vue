<script setup>
import { useRouter } from 'vue-router'
import { Sparkles, Mail, Phone, MapPin, Github, Linkedin, Facebook, Globe, ShieldCheck, Heart, Bot } from 'lucide-vue-next'
import { langStore } from '../../stores/lang.store'
import { usePlatformStore } from '../../stores/platform.store'
import AppLogo from '../common/AppLogo.vue'

const router = useRouter()
const platformStore = usePlatformStore()
</script>

<template>
  <footer class="app-footer">
    <div class="footer-container">
      <div class="footer-grid">
        
        <!-- Brand Column (2 cols) -->
        <div class="brand-col">
          <div @click="router.push('/')" class="cursor-pointer w-fit">
            <AppLogo size="lg" showSubtitle subtitle="Nền tảng Phỏng vấn & Tuyển dụng AI" />
          </div>
          
          <p class="brand-desc mt-2">
            {{ langStore.t('footer', 'desc') }}
          </p>

          <div class="social-links">
            <a href="https://github.com" target="_blank" aria-label="GitHub" class="social-item">
              <Github :size="18" />
            </a>
            <a href="https://linkedin.com" target="_blank" aria-label="LinkedIn" class="social-item">
              <Linkedin :size="18" />
            </a>
            <a href="https://facebook.com" target="_blank" aria-label="Facebook" class="social-item">
              <Facebook :size="18" />
            </a>
            <a href="https://interview-ai.vn" target="_blank" aria-label="Website" class="social-item">
              <Globe :size="18" />
            </a>
          </div>
        </div>

        <!-- Solutions Column -->
        <div class="footer-col">
          <h3 class="col-heading">
            <Sparkles :size="15" class="text-primary" />
            <span>{{ langStore.t('footer', 'solutionsTitle') }}</span>
          </h3>
          <ul class="footer-links">
            <li><a @click.prevent="router.push('/login')" href="#">{{ langStore.t('footer', 'sol1') }}</a></li>
            <li><a @click.prevent="router.push('/login')" href="#">{{ langStore.t('footer', 'sol2') }}</a></li>
            <li><a @click.prevent="router.push('/login')" href="#">{{ langStore.t('footer', 'sol3') }}</a></li>
            <li><a @click.prevent="router.push('/login')" href="#">{{ langStore.t('footer', 'sol4') }}</a></li>
          </ul>
        </div>

        <!-- Candidates Column -->
        <div class="footer-col">
          <h3 class="col-heading">
            <span>{{ langStore.t('footer', 'candTitle') }}</span>
          </h3>
          <ul class="footer-links">
            <li><a @click.prevent="router.push('/job-board')" href="#">{{ langStore.t('footer', 'cand1') }}</a></li>
            <li><a @click.prevent="router.push('/mock-setup')" href="#">{{ langStore.t('footer', 'cand2') }}</a></li>
            <li><a @click.prevent="router.push('/mock-results')" href="#">{{ langStore.t('footer', 'cand3') }}</a></li>
            <li><a @click.prevent="router.push('/profile')" href="#">{{ langStore.t('footer', 'cand4') }}</a></li>
          </ul>
        </div>

        <!-- Contact & Security Column -->
        <div class="footer-col">
          <h3 class="col-heading">
            <ShieldCheck :size="15" class="text-accent" />
            <span>{{ langStore.t('footer', 'secTitle') }}</span>
          </h3>
          <ul class="footer-contact">
            <li>
              <Mail :size="15" class="contact-icon" />
              <span>{{ platformStore.supportEmail }}</span>
            </li>
            <li>
              <Phone :size="15" class="contact-icon" />
              <span>+84 (0) 28 8888 2026</span>
            </li>
            <li class="contact-address">
              <MapPin :size="15" class="contact-icon top-offset" />
              <span>Khu Công nghệ Cao, TP. Thủ Đức, TP. Hồ Chí Minh</span>
            </li>
          </ul>
        </div>

      </div>

      <!-- Bottom Bar -->
      <div class="footer-bottom">
        <p class="copyright-text">
          <span>© 2026 {{ platformStore.brandName }} {{ platformStore.brandBadge }} Platform. Made with</span>
          <Heart :size="14" class="text-heart" />
          <span>& AI Real-time Technology.</span>
        </p>
        <div class="legal-links">
          <a @click.prevent="router.push('/')" href="#">{{ langStore.t('footer', 'terms') }}</a>
          <a @click.prevent="router.push('/')" href="#">{{ langStore.t('footer', 'privacy') }}</a>
          <a @click.prevent="router.push('/')" href="#">{{ langStore.t('footer', 'help') }}</a>
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.app-footer {
  background-color: var(--surface, #ffffff);
  border-top: 1px solid var(--border, #e2e8f0);
  color: var(--text-secondary, #475569);
  font-family: var(--sans, 'Inter', -apple-system, sans-serif);
  font-size: 0.875rem;
  transition: background-color 0.3s ease, border-color 0.3s ease;
}

/* Main Container */
.footer-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 48px 32px 32px;
}

.footer-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1.25fr;
  gap: 40px;
  margin-bottom: 48px;
}

/* Brand Column */
.brand-col {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.brand-logo-group {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  width: fit-content;
}
.brand-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: var(--primary, #2563eb);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
}
.brand-titles {
  display: flex;
  flex-direction: column;
}
.brand-main {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--text-main, #0f172a);
  letter-spacing: -0.02em;
  line-height: 1.2;
}
.brand-sub {
  font-size: 0.7rem;
  font-weight: 700;
  color: var(--text-muted, #64748b);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin-top: 2px;
}
.brand-desc {
  font-size: 0.85rem;
  color: var(--text-secondary, #475569);
  line-height: 1.6;
  max-width: 380px;
}

.social-links {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 4px;
}
.social-item {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background-color: var(--surface-soft, #f1f5f9);
  color: var(--text-secondary, #475569);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}
.social-item:hover {
  background-color: var(--primary, #2563eb);
  color: #ffffff;
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.2);
}

/* Columns */
.footer-col {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.col-heading {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-main, #0f172a);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
}
.text-primary { color: var(--primary, #2563eb); }
.text-accent { color: var(--accent, #0891b2); }

.footer-links {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.footer-links li a {
  color: var(--text-secondary, #475569);
  text-decoration: none;
  font-size: 0.85rem;
  font-weight: 500;
  transition: color 0.2s ease, padding-left 0.2s ease;
  display: inline-block;
}
.footer-links li a:hover {
  color: var(--primary, #2563eb);
  padding-left: 4px;
}

.footer-contact {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
  font-size: 0.85rem;
}
.footer-contact li {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary, #475569);
}
.footer-contact li.contact-address {
  align-items: flex-start;
}
.contact-icon {
  color: var(--primary, #2563eb);
  flex-shrink: 0;
}
.contact-icon.top-offset {
  margin-top: 3px;
}

/* Bottom Bar */
.footer-bottom {
  padding-top: 24px;
  border-top: 1px solid var(--border, #e2e8f0);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 0.8rem;
  color: var(--text-muted, #64748b);
}
.copyright-text {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
}
.text-heart {
  color: #ef4444;
  fill: #ef4444;
}
.legal-links {
  display: flex;
  align-items: center;
  gap: 24px;
  font-weight: 600;
}
.legal-links a {
  color: var(--text-muted, #64748b);
  text-decoration: none;
  transition: color 0.2s ease;
}
.legal-links a:hover {
  color: var(--primary, #2563eb);
}

/* Responsive */
@media (max-width: 1024px) {
  .footer-grid {
    grid-template-columns: 1fr 1fr;
    gap: 32px;
  }
}
@media (max-width: 640px) {
  .footer-grid {
    grid-template-columns: 1fr;
    gap: 28px;
  }
  .footer-bottom {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  .legal-links {
    gap: 16px;
    flex-wrap: wrap;
  }
}
</style>
