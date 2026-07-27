# DESIGN SYSTEM SKILL — AI Interview Platform

> **Mục đích của file này:** Đây là "skill" thiết kế giao diện dùng chung cho toàn bộ dự án. Bất kỳ ai (người hoặc AI) khi tạo mới / sửa một trang UI **phải đọc file này trước**, và tuân theo đúng token + pattern ở đây. File này biến `docs/GiaoDien.md` (bản mô tả ý tưởng) thành **quy tắc thực thi được**.

> Ngôn ngữ nguồn chân lý: `frontend/src/styles/global.css` là nơi định nghĩa token. File này giải thích cách dùng chúng. Khi hai bên mâu thuẫn → sửa cho khớp, không tạo token mới tùy tiện.

---

## 0. TL;DR — Đọc 60 giây

1. **Dùng token, không hardcode màu.** Viết `var(--primary)`, `var(--accent)`, `var(--highlight)`, không viết `from-blue-600`, `text-indigo-600`, `bg-emerald-100`.
2. **Hệ 3 Màu Chủ Đạo Chuyên Nghiệp (Royal Navy & Tech AI):**
   - **Primary Navy (#1E3A8A - `--primary`)**: Màu chủ đạo cho thương hiệu, header, nút chính, tạo cảm giác uy tín tuyệt đối cho HR Tech / Enterprise.
   - **Secondary / Accent AI Cyan (#06B6D4 - `--accent`)**: Màu công nghệ & AI, dùng cho link, icon nổi bật, border focus và trạng thái active.
   - **Accent Highlight Amber (#F59E0B - `--highlight`)**: CTA quan trọng, badge nổi bật, thành tích cao (KHÔNG dùng thay thế cho warning lỗi, tách biệt rõ ràng).
3. **Thang màu Neutral / Gray Scale chuẩn Tailwind (`gray-50` -> `gray-900`)**: Dùng cho 70-80% giao diện (background, surface, border, text).
4. **Không gradient loạn, không glassmorphism, không emoji trong app.** Nền trắng/gray-50, viền gray-200, đổ bóng rất nhẹ.
5. **Chữ tiêu đề trang 24–30px.** Không dùng `text-3xl`/`text-4xl` gradient text bên trong app.
6. **Recruiter/Admin = dày dữ liệu, quyết đoán. Candidate = nhẹ nhàng, trấn an.** Nhưng dùng chung một hệ token.
7. **Ưu tiên component chung** `src/components/common/App*.vue` thay vì tự chế lại card/button/badge/table.

---

## 1. Chẩn đoán hiện trạng (vì sao cần file này)

Dự án **đã có** một hệ thiết kế tốt trên giấy và trong `global.css`, nhưng các trang thực tế **không tuân theo nó**. Đây là lý do UI trông "chưa hiện đại / thiếu nhất quán":

| Vấn đề | Ví dụ đang tồn tại | Hậu quả |
|---|---|---|
| Bỏ qua token, hardcode Tailwind | `CandidateDashboard.vue`, `MyInterviewsPage.vue`, `JobBoardPage.vue` dùng `from-blue-600 to-indigo-600` | Mỗi trang một sắc độ xanh/tím khác nhau |
| Dùng màu bị cấm | `AdminLayout.vue` dùng gradient `indigo-600 → purple-600` | Trái spec ("no purple"), lệch nhận diện |
| Gradient + glassmorphism tùy hứng | `backdrop-blur-xl`, `bg-white/80`, banner gradient 3 màu | Trông như landing page marketing, không giống HR SaaS nghiêm túc |
| Chữ quá lớn trong app | `text-3xl`/`text-4xl` + gradient text ở dashboard | Phá nhịp "compact, dễ scan" |
| Token "ma" | `AppButton.vue` gọi `var(--text-h)`, `var(--shadow-glow)`, `var(--primary-light)` — **không tồn tại** trong `global.css` | Style rơi về mặc định, khó bảo trì |
| Bóng/viền không nhất quán | chỗ `shadow-xl`, chỗ `shadow-lg`, chỗ `var(--shadow-md)` | Độ nổi các card khác nhau vô lý |

**Kết luận:** vấn đề không phải "thiếu polish" mà là **thiếu nhất quán và lệch brand**. Fix = hợp nhất token → ép các trang dùng token → dọn màu cấm.

---

## 2. Design Tokens (nguồn chân lý)

Toàn bộ token sống trong `:root` của `frontend/src/styles/global.css`. **Không định nghĩa lại màu trong file .vue.** Nếu thiếu token, thêm vào `global.css` rồi mới dùng.

### 2.1 Bảng token chuẩn (bổ sung cho global.css)

```css
:root {
  /* 1. Brand Colors (Core) */
  --primary:        #1E3A8A;  /* Navy — thương hiệu, header, nút chính */
  --primary-hover:  #1D4ED8;
  --primary-light:  #DBEAFE;
  
  --accent:         #06B6D4;  /* Cyan — accent công nghệ/AI, link, icon nổi bật, border focus */
  --accent-bg:      #ECFEFF;

  --highlight:      #F59E0B;  /* Amber — CTA quan trọng, badge nổi bật */
  --highlight-hover:#D97706;
  --highlight-bg:   #FEF3C7;

  /* 2 & 4. Background layers & Neutral / Gray scale */
  --background:     #FFFFFF;  /* Background chính light */
  --surface:        #F8FAFC;  /* Background phụ gray-50 */
  --surface-soft:   #F1F5F9;  /* Background hover/active gray-100 */

  /* 6. Border & Focus */
  --border:         #E2E8F0;  /* Border mặc định gray-200 */
  --border-focus:   #06B6D4;  /* Border focus Cyan */

  /* 5. Text colors */
  --text-main:      #0F172A;  /* Text chính gray-900 (không dùng #000 thuần) */
  --text-secondary: #64748B;  /* Text phụ/mô tả gray-500 */
  --text-muted:     #CBD5E1;  /* Text disabled gray-300 */

  /* 3. Semantic colors (Trạng thái) */
  --success:        #10B981;  /* Match thành công, duyệt, hoàn thành */
  --warning:        #F97316;  /* Cảnh báo nhẹ */
  --danger:         #EF4444;  /* Lỗi validation, xóa, từ chối */
  --info:           #3B82F6;  /* Thông báo trung tính */

  /* ---- Radius ---- */
  --radius:         12px;
  --radius-lg:      16px;
  --radius-full:    9999px;

  /* ---- Shadow (rất nhẹ, KHÔNG glow) ---- */
  --shadow-sm: 0 1px 2px 0 rgba(15, 23, 42, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(15, 23, 42, 0.05), 0 2px 4px -2px rgba(15, 23, 42, 0.05);
  --shadow-lg: 0 10px 15px -3px rgba(15, 23, 42, 0.08), 0 4px 6px -4px rgba(15, 23, 42, 0.05);

  /* ---- Typography ---- */
  --sans: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
```

> **Việc cần làm ngay:** `AppButton.vue` đã được chuẩn hóa dùng solid `--primary` (`#1E3A8A`) và `--shadow-sm/md`. Toàn bộ các trang (Dashboard, Candidate, Recruiter, Admin, Auth) được ép dùng token chung này.

### 2.2 Quy tắc dùng bộ 3 màu chủ đạo & thang neutral

- **Primary Navy (`--primary: #1E3A8A`)**: Thương hiệu, header, nút chính (CTA chính), tạo cảm giác uy tín cao cấp cho HR/Enterprise.
- **Secondary Cyan (`--accent: #06B6D4`)**: Accent công nghệ/AI, link, icon nổi bật, trạng thái active/focus, viền focus input.
- **Accent Highlight Amber (`--highlight: #F59E0B`)**: CTA quan trọng đặc biệt, badge nổi bật, điểm số cao. **Tách biệt hoàn toàn với màu Warning (`--warning: #F97316`)**.
- **Thang Neutral (`gray-50` -> `gray-900`)**: Chiếm 70-80% giao diện, không dùng đen `#000` thuần cho chữ, luôn dùng `--text-main` (`#0F172A`) hoặc `--text-secondary` (`#64748B`).
- **Trạng thái nghiệp vụ**: dùng đúng `--success (#10B981) / --warning (#F97316) / --danger (#EF4444)`. Badge nền = màu @ 10% alpha, chữ = màu đặc (đã có sẵn `.badge-*`).
- **CẤM**: tím, indigo-làm-chính, hồng, neon, màu gaming/crypto loè loẹt, gradient trong app nội bộ.

---

## 3. Typography

| Vai trò | Size | Weight | Token/class |
|---|---|---|---|
| Tiêu đề trang | 24–30px | 600 | `.text-h1` (24px) |
| Tiêu đề section / card | 16–18px | 600 | `.text-h2` / `.card-title` |
| Body | 14px | 400–500 | mặc định / `.text-body` |
| Helper / phụ | 12–13px | 400 | `.text-helper` |
| Nút | 14px | 600 | trong `.btn` |
| Số liệu KPI | 28px | 700 | `.stat-value` |

**Quy tắc:**
- Font duy nhất: **Inter** (`var(--sans)`).
- **Không** `text-3xl`/`text-4xl`/`text-5xl` bên trong màn hình app. Con số KPI tối đa 28–32px.
- **Không** gradient text (`bg-clip-text text-transparent`) trong app. Tiêu đề dùng `var(--text-main)`.
- Không dùng emoji thay icon. Dùng **Lucide** (`lucide-vue-next`), size 18–20px, màu `currentColor`.

---

## 4. Layout & Spacing (hệ 8px)

| Thành phần | Giá trị |
|---|---|
| Sidebar (recruiter/admin) | 240px (thu gọn: 64px) |
| Header / navbar | 64px (candidate top-nav hiện 72px — chấp nhận, giữ nhất quán trong nhóm candidate) |
| Padding trang | 24px hoặc 32px |
| Padding card | 20px hoặc 24px |
| Border radius | 12px (`--radius`) đến 16px (`--radius-lg`) |
| Chiều cao dòng bảng | 56–64px |
| Chiều cao input | 40–44px |
| Gap grid/section | 16–24px |

- Bố cục **desktop-first**, nhưng responsive: KPI grid `repeat(auto-fit, minmax(200px,1fr))`, nội dung chính `lg:grid-cols-3` với cột chính chiếm 2.
- Tránh khoảng trắng lớn vô nghĩa **và** panel chật cứng. Ưu tiên "dễ scan".

---

## 5. Component System — dùng lại, đừng chế lại

Đã có sẵn trong `src/components/common/`. **Luôn ưu tiên import các component này** thay vì viết lại markup + class Tailwind:

| Cần | Dùng | Ghi chú |
|---|---|---|
| Nút | `AppButton` (`variant`: primary/secondary/ghost/danger) | Không tự viết `<button class="bg-blue-600...">` |
| Thẻ chứa | `AppCard` | Nền trắng, viền `--border`, `--shadow-sm`, radius-lg |
| Input / textarea | `AppInput` | Có label, error, focus ring xanh |
| Select | `AppSelect` | |
| Badge trạng thái | `AppBadge` (`type`: success/warning/danger/info/neutral) | |
| Bảng | `AppTable` | th uppercase muted, dòng 64px |
| Modal | `AppModal` (`size`: sm/md/lg/xl) | Xác nhận trước hành động nguy hiểm |
| Toast | `AppToast` | |
| Rỗng dữ liệu | `AppEmptyState` | icon + message + CTA |
| Lỗi | `AppErrorState` | |
| Loading | `AppLoadingSkeleton` | **Dùng skeleton thay vì spinner** cho tải dữ liệu |
| Tabs | `AppTabs` | |
| Phân trang | `AppPagination` | |
| Avatar | `AppAvatar` | |

> Các component còn ở dạng "stub" (`AppEmptyState`, `AppLoadingSkeleton`, `AppTabs`, `AppSelect`, `AppPagination`, `AppDropdown`) cần được hoàn thiện bằng token trước khi dùng rộng. Đây là ưu tiên hạ tầng.

### 5.1 Component đặc thù sản phẩm (cần chuẩn hóa)
`AI insight card`, `Candidate card`, `Job card`, `Interview card`, `Report card`, `Score badge`, `Rubric score`, `Transcript item`, `Video tile`, `Voice recording button`, `Audio waveform`, `AI thinking indicator`, `Consent notice`, `Device test`, `CV preview card`, `Timeline`, `Stepper`, `Notification item`.
→ Tất cả phải dựng trên token. Khối AI dùng `--accent`/`--accent-bg`, khối thường dùng `--primary`.

---

## 6. Khác biệt theo vai trò (cùng token, khác cảm giác)

| | Recruiter | Candidate | Admin |
|---|---|---|---|
| Điều hướng | Sidebar 240px trái, viền phải mảnh | Top-nav ngang, nhẹ, thân thiện | Sidebar thu gọn được |
| Tông | Năng suất, dày dữ liệu, quyết đoán | Trấn an, khích lệ, thoáng | Kỹ thuật, gọn, ưu tiên bảng/log |
| Mật độ | Cao (bảng, KPI, AI panel) | Vừa (card lớn, ít cột) | Cao |
| Nhấn màu | `--primary` xanh | `--primary` xanh + đôi chỗ mềm hơn | `--primary` xanh — **bỏ gradient tím hiện tại** |
| AI | Panel trợ lý bên phải phòng phỏng vấn | Feedback sau buổi mock | — |

**Nguyên tắc bất biến:** khác nhau về **mật độ và giọng điệu**, KHÔNG khác nhau về **bảng màu**. Cả ba dùng chung `--primary` xanh. Admin **phải bỏ** gradient `indigo→purple` trong `AdminLayout.vue`, thay bằng `--primary` đặc hoặc nền xanh nhạt `--primary-light`.

---

## 7. UX States bắt buộc

Mọi màn hình có dữ liệu async phải xử lý đủ:
- **Loading** → `AppLoadingSkeleton` (không spinner trần).
- **Empty** → `AppEmptyState` (icon + câu giải thích + CTA).
- **Error** → `AppErrorState` (thông báo + nút thử lại).
- **Trạng thái nghiệp vụ** hiển thị bằng badge: Scheduled / Waiting / Active / Completed / Cancelled / Report Ready / AI Analyzing / Transcript Enabled.
- **Xác nhận** bằng `AppModal` trước hành động khó hoàn tác (kết thúc phỏng vấn, xóa).
- **Consent notice** trước khi ghi hình/transcript.
- **Cảnh báo** khi AI thiếu bằng chứng để chấm điểm.

Bảo mật hiển thị (candidate KHÔNG được thấy): điểm AI nội bộ, ghi chú riêng recruiter, gợi ý/đề xuất AI nội bộ, điểm rubric, quyết định tuyển trong lúc phỏng vấn.

---

## 8. Checklist review UI (dán vào PR)

Trước khi merge một trang, tự trả lời **CÓ** cho tất cả:

- [ ] Không có mã màu hardcode (`#hex`, `from-blue-*`, `text-indigo-*`, `bg-emerald-*`) — chỉ `var(--...)` hoặc class token.
- [ ] Không dùng tím/purple/indigo-làm-chính, không gradient nhiều màu, không `backdrop-blur` glass trong app.
- [ ] Tiêu đề ≤ 30px, không gradient text, không emoji thay icon.
- [ ] Dùng `AppButton/AppCard/AppInput/AppBadge/AppTable...` thay vì tự chế.
- [ ] Có đủ Loading (skeleton) / Empty / Error states.
- [ ] Bóng dùng `--shadow-sm/md/lg`, radius dùng `--radius*`, spacing theo hệ 8px.
- [ ] Icon là Lucide, size 18–20px, `currentColor`.
- [ ] Nếu là khối AI: dùng `--accent`/`--accent-bg`, gắn nhãn rõ (VD "Gợi ý từ AI").
- [ ] Candidate không lộ dữ liệu đánh giá nội bộ.
- [ ] Kiểm tra trên màn nhỏ (responsive) — grid co giãn, không tràn ngang.

---

## 9. Anti-patterns (tuyệt đối tránh)

```html
<!-- ❌ SAI: hardcode màu, gradient tím, chữ khổng lồ, glass -->
<h1 class="text-4xl bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600">...</h1>
<div class="bg-white/80 backdrop-blur-xl shadow-xl rounded-2xl">
  <div class="bg-gradient-to-tr from-indigo-600 to-purple-600">A</div>
</div>
```

```html
<!-- ✅ ĐÚNG: token, phẳng, calm, dùng component chung -->
<h1 class="text-h1">Tổng quan</h1>
<AppCard>
  <div class="kpi-icon"><Briefcase :size="18" /></div>
  <div class="text-helper">Jobs đang mở</div>
  <div class="stat-value">{{ stats.jobs }}</div>
</AppCard>

<style scoped>
.kpi-icon {
  display: inline-flex; padding: 8px;
  border-radius: var(--radius);
  background: var(--primary-light); color: var(--primary);
}
</style>
```

---

## 10. Lộ trình áp dụng (đề xuất, chưa thực hiện)

Ưu tiên theo thứ tự để đạt nhất quán nhanh nhất:

1. **Hạ tầng token**: hợp nhất token vào `global.css` (mục 2.1), sửa token "ma" trong `AppButton.vue`.
2. **Hoàn thiện common stubs**: `AppEmptyState`, `AppLoadingSkeleton`, `AppSelect`, `AppTabs`, `AppPagination`, `AppDropdown` bằng token.
3. **Dọn màu cấm**: bỏ gradient indigo→purple ở `AdminLayout.vue`; đổi các `indigo-*` ở recruiter dashboard về `--primary`.
4. **Refactor trang candidate** (đang lệch nhiều nhất): Dashboard → JobBoard → MyInterviews → MockSetup → CVUpload → Profile → PracticeHistory.
5. **Refactor recruiter** rồi **admin** theo cùng checklist mục 8.
6. Mỗi trang refactor xong chạy qua checklist §8.

> Các bước 1–6 là **thay đổi code** — sẽ làm ở phần triển khai riêng, không nằm trong phạm vi tạo tài liệu này.

---

## 11. Tham chiếu

- Ý tưởng thiết kế gốc (mô tả từng màn hình): `docs/GiaoDien.md`
- Token thực thi: `frontend/src/styles/global.css`
- Component chung: `frontend/src/components/common/App*.vue`
- Layout: `frontend/src/components/layout/{Candidate,Recruiter,Admin}Layout.vue`
