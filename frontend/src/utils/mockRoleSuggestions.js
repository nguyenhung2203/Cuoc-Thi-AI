/**
 * Suggest mock-interview roles from parsed CV (TopCV-style major categories).
 * Returns primary keyword matches + related neighbors (typically 4–8 cards).
 */

/** @typedef {{ id: string, name: string, desc: string, icon: string, keywords: string[], related?: string[] }} RoleDef */

/** @type {RoleDef[]} */
const ROLE_CATALOG = [
  // —— Công nghệ thông tin ——
  {
    id: 'frontend',
    name: 'Frontend Developer',
    desc: 'React, Vue, UI',
    icon: 'Code',
    keywords: ['react', 'vue', 'angular', 'javascript', 'typescript', 'next.js', 'nuxt', 'html', 'css', 'tailwind', 'frontend', 'front-end', 'front end', 'ui developer', 'lap trinh web'],
    related: ['fullstack', 'mobile', 'qa'],
  },
  {
    id: 'backend',
    name: 'Backend Developer',
    desc: 'API, Server, DB',
    icon: 'Server',
    keywords: ['golang', 'spring', 'nodejs', 'express', 'nestjs', 'django', 'fastapi', 'dotnet', '.net', 'microservices', 'backend', 'back-end', 'back end', 'postgresql', 'mysql', 'mongodb', 'redis', 'kafka', 'lap trinh backend'],
    related: ['fullstack', 'devops', 'ai_ml'],
  },
  {
    id: 'fullstack',
    name: 'Fullstack Engineer',
    desc: 'End-to-End',
    icon: 'Layers',
    keywords: ['fullstack', 'full-stack', 'full stack', 'mern', 'mean'],
    related: ['frontend', 'backend'],
  },
  {
    id: 'mobile',
    name: 'Mobile Developer',
    desc: 'iOS, Android',
    icon: 'Smartphone',
    keywords: ['android', 'ios', 'swift', 'kotlin', 'flutter', 'react native', 'mobile', 'app developer', 'lap trinh di dong'],
    related: ['frontend', 'fullstack'],
  },
  {
    id: 'devops',
    name: 'DevOps & Cloud',
    desc: 'AWS, K8s, CI/CD',
    icon: 'Cloud',
    keywords: ['devops', 'aws', 'azure', 'gcp', 'kubernetes', 'k8s', 'docker', 'ci/cd', 'terraform', 'ansible', 'sre', 'cloud engineer'],
    related: ['backend', 'fullstack'],
  },
  {
    id: 'qa',
    name: 'QA & Automation',
    desc: 'Testing',
    icon: 'ShieldCheck',
    keywords: ['qa', 'tester', 'kiem thu', 'selenium', 'cypress', 'playwright', 'automation test', 'manual test', 'quality assurance', 'test engineer'],
    related: ['frontend', 'backend'],
  },
  {
    id: 'ai_ml',
    name: 'AI / ML Engineer',
    desc: 'LLM, ML, Data Science',
    icon: 'Cpu',
    keywords: ['machine learning', 'deep learning', 'llm', 'rag', 'pytorch', 'tensorflow', 'nlp', 'computer vision', 'ai engineer', 'ml engineer', 'data scientist', 'khoa hoc du lieu'],
    related: ['data', 'backend'],
  },
  {
    id: 'data',
    name: 'Data Engineer / Analyst',
    desc: 'ETL, BI, Analytics',
    icon: 'Layers',
    keywords: ['data engineer', 'data analyst', 'phan tich du lieu', 'etl', 'data warehouse', 'power bi', 'tableau', 'spark', 'airflow', 'bigquery', 'dbt', 'business intelligence'],
    related: ['ai_ml', 'business_analyst'],
  },
  {
    id: 'security',
    name: 'Cyber Security',
    desc: 'Bảo mật thông tin',
    icon: 'ShieldCheck',
    keywords: ['cyber security', 'an ninh mang', 'bao mat', 'infosec', 'penetration', 'soc analyst', 'security engineer'],
    related: ['devops', 'backend'],
  },

  // —— Kinh tế / Tài chính / Ngân hàng / Bảo hiểm ——
  {
    id: 'accountant',
    name: 'Kế toán tổng hợp',
    desc: 'Kế toán / Thuế / Báo cáo',
    icon: 'Calculator',
    keywords: [
      'ke toan', 'kế toán', 'accountant', 'accounting', 'so sach', 'hach toan', 'bao cao tai chinh',
      'ke toan tong hop', 'ke toan noi bo', 'ke toan cong no', 'ke toan thanh toan', 'misa', 'fast accounting',
    ],
    related: ['tax_accountant', 'auditor', 'financial_analyst'],
  },
  {
    id: 'tax_accountant',
    name: 'Kế toán thuế',
    desc: 'Thuế GTGT, TNDN, quyết toán',
    icon: 'Calculator',
    keywords: ['ke toan thue', 'kế toán thuế', 'thue gtgt', 'thue tndn', 'quyet toan thue', 'tax accountant', 'tax'],
    related: ['accountant', 'auditor'],
  },
  {
    id: 'chief_accountant',
    name: 'Kế toán trưởng',
    desc: 'Quản lý bộ máy kế toán',
    icon: 'Briefcase',
    keywords: ['ke toan truong', 'kế toán trưởng', 'chief accountant', 'truong phong ke toan'],
    related: ['accountant', 'financial_controller'],
  },
  {
    id: 'auditor',
    name: 'Kiểm toán',
    desc: 'Internal / External Audit',
    icon: 'ClipboardCheck',
    keywords: ['kiem toan', 'kiểm toán', 'audit', 'auditor', 'internal audit', 'kiem toan vien'],
    related: ['accountant', 'tax_accountant', 'financial_analyst'],
  },
  {
    id: 'financial_analyst',
    name: 'Chuyên viên Phân tích Tài chính',
    desc: 'Financial Analysis / FP&A',
    icon: 'TrendingUp',
    keywords: [
      'phan tich tai chinh', 'phân tích tài chính', 'financial analyst', 'fpa', 'fp&a', 'tai chinh doanh nghiep',
      'dinh gia', 'valuation', 'mo hinh tai chinh', 'financial modeling', 'cfa',
    ],
    related: ['investment_analyst', 'banking_officer', 'accountant'],
  },
  {
    id: 'investment_analyst',
    name: 'Chuyên viên Đầu tư',
    desc: 'Investment / Portfolio',
    icon: 'TrendingUp',
    keywords: ['dau tu', 'đầu tư', 'investment', 'portfolio', 'chung khoan', 'chứng khoán', 'trai phieu', 'co phieu', 'fund'],
    related: ['financial_analyst', 'banking_officer'],
  },
  {
    id: 'banking_officer',
    name: 'Chuyên viên Ngân hàng',
    desc: 'Tín dụng / Giao dịch / RM',
    icon: 'Landmark',
    keywords: [
      'ngan hang', 'ngân hàng', 'banking', 'tin dung', 'tín dụng', 'relationship manager', 'giao dich vien',
      'credit officer', 'retail banking', 'corporate banking', 'the tin dung',
    ],
    related: ['financial_analyst', 'insurance', 'sales'],
  },
  {
    id: 'insurance',
    name: 'Bảo hiểm',
    desc: 'Tư vấn / Underwriting',
    icon: 'ShieldCheck',
    keywords: ['bao hiem', 'bảo hiểm', 'insurance', 'underwriter', 'tu van bao hiem', 'life insurance'],
    related: ['banking_officer', 'sales', 'financial_analyst'],
  },
  {
    id: 'economist',
    name: 'Chuyên viên Kinh tế',
    desc: 'Kinh tế / Nghiên cứu / Dự báo',
    icon: 'Globe',
    keywords: [
      'kinh te', 'kinh tế', 'economics', 'economist', 'vi mo', 'vi mo kinh te', 'macroeconomics', 'microeconomics',
      'kinh te hoc', 'kinh tế học', 'kinh te phat trien', 'kinh te quoc te', 'kinh tế quốc tế',
      'quan tri kinh te', 'quản trị kinh tế', 'bachelor of economics', 'cu nhan kinh te',
    ],
    related: ['financial_analyst', 'business_analyst', 'sales', 'banking_officer'],
  },
  {
    id: 'financial_controller',
    name: 'Kiểm soát Tài chính',
    desc: 'Financial Controller',
    icon: 'Calculator',
    keywords: ['kiem soat tai chinh', 'financial controller', 'controller', 'quan tri tai chinh'],
    related: ['accountant', 'financial_analyst', 'chief_accountant'],
  },

  // —— Kinh doanh / Bán hàng ——
  {
    id: 'sales',
    name: 'Nhân viên Kinh doanh / Sales',
    desc: 'Bán hàng / Doanh số',
    icon: 'Handshake',
    keywords: [
      'kinh doanh', 'ban hang', 'bán hàng', 'sales', 'sales executive', 'sales representative',
      'phat trien kinh doanh', 'doanh so', 'chot sale', 'b2b', 'b2c',
    ],
    related: ['business_development', 'account_manager', 'telesales'],
  },
  {
    id: 'business_development',
    name: 'Business Development',
    desc: 'BD / Mở rộng thị trường',
    icon: 'TrendingUp',
    keywords: ['business development', 'bd ', 'phat trien thi truong', 'mo rong thi truong', 'partnership'],
    related: ['sales', 'account_manager', 'marketing'],
  },
  {
    id: 'account_manager',
    name: 'Account Manager',
    desc: 'Quản lý khách hàng',
    icon: 'Users',
    keywords: ['account manager', 'quan ly khach hang', 'key account', 'kam', 'client manager'],
    related: ['sales', 'customer_service', 'business_development'],
  },
  {
    id: 'telesales',
    name: 'Telesales / Telemarketing',
    desc: 'Bán hàng qua điện thoại',
    icon: 'Smartphone',
    keywords: ['telesales', 'telemarketing', 'ban hang qua dien thoai', 'cold call'],
    related: ['sales', 'customer_service'],
  },
  {
    id: 'retail_sales',
    name: 'Sales Bán lẻ',
    desc: 'Cửa hàng / Siêu thị',
    icon: 'ShoppingBag',
    keywords: ['ban le', 'bán lẻ', 'retail', 'cua hang', 'sieu thi', 'promoter', 'tu van ban hang'],
    related: ['sales', 'customer_service'],
  },

  // —— Marketing / PR / Quảng cáo ——
  {
    id: 'marketing',
    name: 'Marketing Executive',
    desc: 'Marketing tổng hợp',
    icon: 'Megaphone',
    keywords: ['marketing', 'tiep thi', 'tiếp thị', 'marketing executive', 'marketing mix', 'brand marketing'],
    related: ['digital_marketing', 'content_marketing', 'brand'],
  },
  {
    id: 'digital_marketing',
    name: 'Digital Marketing',
    desc: 'Ads, SEO, Social',
    icon: 'Megaphone',
    keywords: ['digital marketing', 'facebook ads', 'google ads', 'seo', 'sem', 'performance marketing', 'chay ads', 'tiktok ads'],
    related: ['marketing', 'content_marketing', 'brand'],
  },
  {
    id: 'content_marketing',
    name: 'Content Marketing',
    desc: 'Nội dung / Copywriting',
    icon: 'FileText',
    keywords: ['content marketing', 'content creator', 'copywriter', 'content writer', 'viet content', 'bien tap vien'],
    related: ['marketing', 'digital_marketing', 'pr'],
  },
  {
    id: 'brand',
    name: 'Brand / Truyền thông thương hiệu',
    desc: 'Branding / Communication',
    icon: 'Award',
    keywords: ['brand', 'thuong hieu', 'thương hiệu', 'branding', 'brand manager', 'truyen thong thuong hieu'],
    related: ['marketing', 'pr', 'digital_marketing'],
  },
  {
    id: 'pr',
    name: 'PR / Quan hệ công chúng',
    desc: 'PR & Media',
    icon: 'Megaphone',
    keywords: ['public relations', 'pr ', 'quan he cong chung', 'media relation', 'truyen thong'],
    related: ['brand', 'marketing', 'content_marketing'],
  },

  // —— Chăm sóc khách hàng / Vận hành ——
  {
    id: 'customer_service',
    name: 'Chăm sóc khách hàng',
    desc: 'CSKH / Customer Service',
    icon: 'Headset',
    keywords: ['cham soc khach hang', 'chăm sóc khách hàng', 'customer service', 'cskh', 'call center', 'hotline', 'support'],
    related: ['sales', 'account_manager', 'admin'],
  },
  {
    id: 'operations',
    name: 'Vận hành / Operations',
    desc: 'Operations Executive',
    icon: 'Settings',
    keywords: ['van hanh', 'vận hành', 'operations', 'ops ', 'dieu phoi', 'operation executive'],
    related: ['customer_service', 'logistics', 'admin'],
  },

  // —— Nhân sự / Hành chính / Pháp chế ——
  {
    id: 'hr_generalist',
    name: 'Nhân sự tổng hợp',
    desc: 'HR Generalist',
    icon: 'Users',
    keywords: ['nhan su', 'nhân sự', 'human resources', 'hr generalist', 'hr ', 'quan tri nhan su'],
    related: ['recruiter', 'payroll', 'admin'],
  },
  {
    id: 'recruiter',
    name: 'Tuyển dụng',
    desc: 'Talent Acquisition',
    icon: 'UserPlus',
    keywords: ['tuyen dung', 'tuyển dụng', 'recruiter', 'talent acquisition', 'sourcing', 'headhunt'],
    related: ['hr_generalist', 'payroll'],
  },
  {
    id: 'payroll',
    name: 'C&B / Payroll',
    desc: 'Lương thưởng & phúc lợi',
    icon: 'Calculator',
    keywords: ['payroll', 'c&b', 'compensation', 'luong thuong', 'phuc loi', 'bhxh', 'bao hiem xa hoi'],
    related: ['hr_generalist', 'accountant'],
  },
  {
    id: 'admin',
    name: 'Hành chính văn phòng',
    desc: 'Admin / Office',
    icon: 'Building2',
    keywords: ['hanh chinh', 'hành chính', 'admin', 'van thu', 'van phong', 'office admin', 'thu ky', 'assistant'],
    related: ['hr_generalist', 'customer_service', 'operations'],
  },
  {
    id: 'legal',
    name: 'Pháp chế / Luật',
    desc: 'Legal Counsel',
    icon: 'Scale',
    keywords: ['phap che', 'pháp chế', 'luat su', 'luật sư', 'legal', 'compliance', 'hop dong', 'luat'],
    related: ['admin', 'hr_generalist'],
  },

  // —— Product / BA ——
  {
    id: 'pm',
    name: 'Product Manager',
    desc: 'Product Strategy',
    icon: 'Briefcase',
    keywords: ['product manager', 'product owner', 'quan ly san pham', 'roadmap', 'prd', 'agile', 'scrum'],
    related: ['business_analyst', 'marketing'],
  },
  {
    id: 'business_analyst',
    name: 'Business Analyst',
    desc: 'Phân tích nghiệp vụ',
    icon: 'Briefcase',
    keywords: [
      'business analyst', 'phan tich nghiep vu', 'phân tích nghiệp vụ', 'requirement', 'business process',
      'use case', 'ba ', 'phan tich kinh doanh',
    ],
    related: ['pm', 'financial_analyst', 'data'],
  },

  // —— Logistics / Thu mua / Kho ——
  {
    id: 'logistics',
    name: 'Logistics / Chuỗi cung ứng',
    desc: 'Supply Chain',
    icon: 'Truck',
    keywords: ['logistics', 'chuoi cung ung', 'supply chain', 'van tai', 'van chuyen', 'freight', 'xnk', 'xuat nhap khau', 'customs'],
    related: ['procurement', 'warehouse', 'operations'],
  },
  {
    id: 'procurement',
    name: 'Thu mua / Procurement',
    desc: 'Purchasing',
    icon: 'ShoppingCart',
    keywords: ['thu mua', 'procurement', 'purchasing', 'buyer', 'mua hang'],
    related: ['logistics', 'warehouse', 'operations'],
  },
  {
    id: 'warehouse',
    name: 'Nhân viên Kho',
    desc: 'Warehouse / Inventory',
    icon: 'Package',
    keywords: ['kho', 'warehouse', 'inventory', 'quan ly kho', 'thu kho'],
    related: ['logistics', 'procurement'],
  },

  // —— Xây dựng / Bất động sản ——
  {
    id: 'construction_engineer',
    name: 'Kỹ sư Xây dựng',
    desc: 'Hiện trường / Công trình',
    icon: 'HardHat',
    keywords: ['xay dung', 'xây dựng', 'ky su xay dung', 'cong trinh', 'giam sat thi cong', 'civil engineer'],
    related: ['architect', 'qs_estimator', 'real_estate'],
  },
  {
    id: 'architect',
    name: 'Kiến trúc sư',
    desc: 'Architecture / Interior',
    icon: 'DraftingCompass',
    keywords: ['kien truc', 'kiến trúc', 'architect', 'thiet ke noi that', 'interior design', 'autocad', 'revit'],
    related: ['construction_engineer', 'graphic_design'],
  },
  {
    id: 'qs_estimator',
    name: 'Dự toán / QS',
    desc: 'Quantity Surveyor',
    icon: 'Calculator',
    keywords: ['du toan', 'dự toán', 'quantity surveyor', 'qs ', 'du toan cong trinh'],
    related: ['construction_engineer', 'accountant'],
  },
  {
    id: 'real_estate',
    name: 'Kinh doanh Bất động sản',
    desc: 'Sales BĐS',
    icon: 'Home',
    keywords: ['bat dong san', 'bất động sản', 'real estate', 'moi gioi bds', 'sales bds', 'can ho', 'chung cu'],
    related: ['sales', 'business_development'],
  },

  // —— Giáo dục ——
  {
    id: 'teacher',
    name: 'Giáo viên / Giảng viên',
    desc: 'Teaching / Training',
    icon: 'GraduationCap',
    keywords: ['giao vien', 'giáo viên', 'giang vien', 'giảng viên', 'teacher', 'gia su', 'dao tao', 'trainer', 'english teacher'],
    related: ['hr_generalist', 'content_marketing'],
  },

  // —— Y tế / Dược ——
  {
    id: 'pharmacist',
    name: 'Dược sĩ / Y tế',
    desc: 'Pharmacy / Healthcare',
    icon: 'HeartPulse',
    keywords: ['duoc', 'dược', 'pharmacist', 'y te', 'y tế', 'dieu duong', 'bac si', 'healthcare', 'clinic'],
    related: ['sales', 'customer_service'],
  },

  // —— Thiết kế ——
  {
    id: 'graphic_design',
    name: 'Thiết kế Đồ họa',
    desc: 'Graphic Design',
    icon: 'Palette',
    keywords: ['thiet ke do hoa', 'thiết kế đồ họa', 'graphic design', 'photoshop', 'illustrator', 'figma', 'ui ux', 'designer'],
    related: ['content_marketing', 'brand', 'frontend'],
  },

  // —— Nhà hàng / Khách sạn / Du lịch ——
  {
    id: 'hospitality',
    name: 'Nhà hàng / Khách sạn',
    desc: 'Hospitality',
    icon: 'Hotel',
    keywords: ['nha hang', 'nhà hàng', 'khach san', 'khách sạn', 'hotel', 'fnb', 'f&b', 'le tan', 'receptionist', 'barista'],
    related: ['tourism', 'customer_service', 'sales'],
  },
  {
    id: 'tourism',
    name: 'Du lịch / Tour',
    desc: 'Travel / Tour guide',
    icon: 'Plane',
    keywords: ['du lich', 'du lịch', 'tourism', 'tour', 'huong dan vien', 'travel agent'],
    related: ['hospitality', 'sales', 'customer_service'],
  },

  // —— Sản xuất ——
  {
    id: 'production',
    name: 'Sản xuất / QC',
    desc: 'Production / Quality',
    icon: 'Factory',
    keywords: ['san xuat', 'sản xuất', 'production', 'qc', 'qa san xuat', 'quan doc', 'co khi', 'co dien'],
    related: ['operations', 'logistics', 'qa'],
  },

  // —— Truyền thông / Báo chí ——
  {
    id: 'media',
    name: 'Báo chí / Biên tập',
    desc: 'Media / Journalism',
    icon: 'Newspaper',
    keywords: ['bao chi', 'báo chí', 'bien tap', 'journalist', 'reporter', 'phong vien', 'truyen hinh'],
    related: ['content_marketing', 'pr', 'brand'],
  },

  // —— Nông nghiệp / Môi trường ——
  {
    id: 'agriculture',
    name: 'Nông nghiệp / Môi trường',
    desc: 'Agri / Environment',
    icon: 'Leaf',
    keywords: ['nong nghiep', 'nông nghiệp', 'moi truong', 'environment', 'nong san', 'agri'],
    related: ['sales', 'logistics'],
  },
]

/** Domain → fallback role ids when no keyword hit (never force IT for non-IT CVs). */
const DOMAIN_FALLBACKS = {
  finance: ['economist', 'accountant', 'financial_analyst', 'banking_officer', 'sales'],
  business: ['sales', 'business_development', 'marketing', 'customer_service'],
  marketing: ['marketing', 'digital_marketing', 'content_marketing', 'brand'],
  hr_admin: ['hr_generalist', 'admin', 'recruiter', 'customer_service'],
  logistics: ['logistics', 'procurement', 'warehouse', 'operations'],
  construction: ['construction_engineer', 'architect', 'qs_estimator', 'real_estate'],
  it: ['fullstack', 'backend', 'frontend', 'qa'],
  default: ['sales', 'customer_service', 'admin', 'hr_generalist'],
}

const DOMAIN_PATTERNS = [
  {
    id: 'finance',
    re: /kinh te|economics|tai chinh|ke toan|accounting|ngan hang|banking|kiem toan|audit|thue|bao hiem|dau tu|tin dung|cfa|fpa|fp&a/,
  },
  {
    id: 'marketing',
    re: /marketing|digital marketing|seo|content|quang cao|brand|pr |truyen thong|tiep thi/,
  },
  {
    id: 'business',
    re: /kinh doanh|ban hang|sales|business development|account manager|telesales|ban le/,
  },
  {
    id: 'hr_admin',
    re: /nhan su|human resources|tuyen dung|hanh chinh|payroll|c&b|phap che|thu ky/,
  },
  {
    id: 'logistics',
    re: /logistics|chuoi cung ung|supply chain|thu mua|procurement|kho |warehouse|xuat nhap khau/,
  },
  {
    id: 'construction',
    re: /xay dung|kien truc|bat dong san|real estate|cong trinh|du toan/,
  },
  {
    id: 'it',
    re: /developer|lap trinh|frontend|backend|fullstack|devops|software engineer|ky su phan mem|tester|react\.?js|nodejs|spring boot/,
  },
]

const LEVEL_HINTS = [
  { id: 'fresher', patterns: [/thuc tap|thực tập|intern|fresher|sinh vien|sinh viên|moi ra truong|mới ra trường/i] },
  { id: 'junior', patterns: [/junior|1\s*-\s*2|duoi\s*2|dưới\s*2/i] },
  { id: 'middle', patterns: [/middle|mid-level|2\s*-\s*4|3\s*-\s*5/i] },
  { id: 'senior', patterns: [/senior|5\+|hon\s*5|hơn\s*5|tren\s*5|trên\s*5/i] },
  { id: 'lead', patterns: [/tech lead|team lead|principal|architect|truong nhom|trưởng nhóm|truong phong|trưởng phòng|lead /i] },
]

function normalizeVi(text) {
  return String(text || '')
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/đ/g, 'd')
    .replace(/Đ/g, 'd')
}

function asList(value) {
  if (!value) return []
  if (Array.isArray(value)) {
    return value.map((item) => {
      if (typeof item === 'string') return item
      if (item && typeof item === 'object') {
        return item.name || item.label || item.title || item.role || item.degree || item.school || ''
      }
      return String(item || '')
    }).filter(Boolean)
  }
  if (typeof value === 'string') return [value]
  return []
}

function collectCorpus(parsed) {
  if (!parsed || typeof parsed !== 'object') return ''
  const chunks = []
  chunks.push(...asList(parsed.skills))
  chunks.push(...asList(parsed.skill_tags))
  chunks.push(...asList(parsed.technologies))
  chunks.push(...asList(parsed.tech_stack))
  chunks.push(...asList(parsed.education))
  chunks.push(...asList(parsed.potential_strengths))
  chunks.push(...asList(parsed.potential_concerns))
  if (parsed.title) chunks.push(parsed.title)
  if (parsed.role) chunks.push(parsed.role)
  if (parsed.target_role) chunks.push(parsed.target_role)
  if (parsed.summary) chunks.push(parsed.summary)
  if (parsed.experience) {
    chunks.push(typeof parsed.experience === 'string' ? parsed.experience : JSON.stringify(parsed.experience))
  }
  if (Array.isArray(parsed.work_experience)) {
    parsed.work_experience.forEach((job) => {
      if (!job) return
      chunks.push(job.role || job.title || '')
      chunks.push(job.company || '')
      chunks.push(job.description || '')
      chunks.push(job.duration || '')
      chunks.push(...asList(job.highlights))
      chunks.push(...asList(job.skills))
    })
  }
  if (Array.isArray(parsed.projects)) {
    parsed.projects.forEach((p) => {
      if (!p) return
      chunks.push(p.name || p.title || '')
      chunks.push(p.description || '')
      chunks.push(...asList(p.technologies || p.tech_stack))
    })
  }
  return normalizeVi(chunks.join(' | '))
}

function scoreRole(role, corpus) {
  let score = 0
  const hits = []
  for (const kw of role.keywords) {
    const needle = normalizeVi(kw).trim()
    if (!needle || needle.length < 2) continue
    // Avoid tiny tokens matching inside unrelated words (e.g. "qa" in "equal")
    const hit = needle.length <= 3
      ? new RegExp(`(?:^|[^a-z0-9])${needle.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}(?:[^a-z0-9]|$)`).test(corpus)
      : corpus.includes(needle)
    if (hit) {
      score += needle.length > 6 ? 3 : needle.length > 4 ? 2 : 1
      if (hits.length < 4) hits.push(kw)
    }
  }
  return { score, hits }
}

function detectDomain(corpus) {
  for (const d of DOMAIN_PATTERNS) {
    if (d.re.test(corpus)) return d.id
  }
  return 'default'
}

function relatedIdsOf(role) {
  return Array.isArray(role.related) ? role.related : []
}

/**
 * @param {object|null} parsedData
 * @param {{ max?: number, minPrimary?: number }} [opts]
 * @returns {{ roles: Array, suggestedLevel: string|null, matchedSkills: string[], domain: string }}
 */
export function suggestRolesFromParsedCV(parsedData, opts = {}) {
  const max = opts.max ?? 8
  const minPrimary = opts.minPrimary ?? 1
  const corpus = collectCorpus(parsedData)
  const domain = detectDomain(corpus)

  const scored = ROLE_CATALOG.map((role) => {
    const { score, hits } = scoreRole(role, corpus)
    return { ...role, score, matchedKeywords: hits, source: score > 0 ? 'cv_match' : 'related' }
  }).sort((a, b) => b.score - a.score)

  const primary = scored.filter((r) => r.score > 0)
  const selected = []
  const seen = new Set()

  for (const role of primary) {
    if (selected.length >= max) break
    selected.push(role)
    seen.add(role.id)
  }

  const relatedTarget = Math.min(max, Math.max(4, primary.length + 2))
  for (const role of primary.slice(0, 3)) {
    for (const relId of relatedIdsOf(role)) {
      if (selected.length >= relatedTarget) break
      if (seen.has(relId)) continue
      const rel = ROLE_CATALOG.find((r) => r.id === relId)
      if (!rel) continue
      selected.push({
        ...rel,
        score: 0,
        matchedKeywords: [],
        source: 'related',
        relatedTo: role.name,
      })
      seen.add(relId)
    }
  }

  if (selected.length < minPrimary) {
    const fallbackIds = DOMAIN_FALLBACKS[domain] || DOMAIN_FALLBACKS.default
    fallbackIds.forEach((id) => {
      if (seen.has(id) || selected.length >= 5) return
      const role = ROLE_CATALOG.find((r) => r.id === id)
      if (!role) return
      selected.push({ ...role, score: 0, matchedKeywords: [], source: 'fallback' })
      seen.add(id)
    })
  }

  const matchedSkills = [...new Set(primary.flatMap((r) => r.matchedKeywords))].slice(0, 12)
  return {
    roles: selected.slice(0, max),
    suggestedLevel: suggestLevelFromParsedCV(parsedData, corpus),
    matchedSkills,
    domain,
  }
}

export function suggestLevelFromParsedCV(parsedData, corpusInput) {
  const corpus = corpusInput || collectCorpus(parsedData) || ''
  for (const level of LEVEL_HINTS) {
    if (level.patterns.some((re) => re.test(corpus))) return level.id
  }
  const yearMatch = corpus.match(/(\d+(?:\.\d+)?)\s*(?:\+|nam|năm|years?|yrs?)/i)
  if (yearMatch) {
    const years = Number(yearMatch[1])
    if (years < 1) return 'fresher'
    if (years < 2) return 'junior'
    if (years < 5) return 'middle'
    if (years < 8) return 'senior'
    return 'lead'
  }
  return null
}

/** Full multi-industry role list for manual selection (not CV-matched). */
export function getRoleCatalog() {
  return ROLE_CATALOG.map(({ id, name, desc, icon }) => ({ id, name, desc, icon }))
}

export function getRoleById(id) {
  return ROLE_CATALOG.find((r) => r.id === id) || null
}
