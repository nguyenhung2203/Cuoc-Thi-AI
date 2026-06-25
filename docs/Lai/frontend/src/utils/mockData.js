// src/utils/mockData.js

// Khởi tạo data giả vào localStorage để có thể CRUD
if (!localStorage.getItem('mockJobs')) {
  localStorage.setItem('mockJobs', JSON.stringify([
    { id: 'job_1', title: 'Frontend Developer', status: 'Open', created: '2026-06-20', applicants: 5, description: 'We need a React developer.' },
    { id: 'job_2', title: 'Backend Engineer', status: 'Closed', created: '2026-05-15', applicants: 12, description: 'Node.js and PostgreSQL.' },
  ]));
}

if (!localStorage.getItem('mockCandidates')) {
  localStorage.setItem('mockCandidates', JSON.stringify([
    { id: 'cand_1', name: 'John Doe', email: 'john@example.com', appliedJob: 'Frontend Developer', status: 'Interviewing', cv: 'john_cv.pdf' },
    { id: 'cand_2', name: 'Jane Smith', email: 'jane@example.com', appliedJob: 'Backend Engineer', status: 'Rejected', cv: 'jane_cv.pdf' },
  ]));
}

if (!localStorage.getItem('mockInterviews')) {
  localStorage.setItem('mockInterviews', JSON.stringify([
    { id: 'int_1', jobTitle: 'Frontend Developer', candidateName: 'John Doe', datetime: '2026-06-28T10:00', status: 'Scheduled', link: 'https://interview.ai/room/abc-123' }
  ]));
}

// Hàm giả lập API delay
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

export const mockApi = {
  jobs: {
    getAll: async () => {
      await delay(500);
      return JSON.parse(localStorage.getItem('mockJobs'));
    },
    getById: async (id) => {
      await delay(300);
      const jobs = JSON.parse(localStorage.getItem('mockJobs'));
      return jobs.find(j => j.id === id);
    },
    create: async (jobData) => {
      await delay(600);
      const jobs = JSON.parse(localStorage.getItem('mockJobs'));
      const newJob = { ...jobData, id: 'job_' + Date.now(), applicants: 0 };
      localStorage.setItem('mockJobs', JSON.stringify([newJob, ...jobs]));
      return newJob;
    }
  },
  candidates: {
    getAll: async () => {
      await delay(500);
      return JSON.parse(localStorage.getItem('mockCandidates'));
    },
    getById: async (id) => {
      await delay(300);
      const cands = JSON.parse(localStorage.getItem('mockCandidates'));
      return cands.find(c => c.id === id);
    },
    create: async (candData) => {
      await delay(600);
      const cands = JSON.parse(localStorage.getItem('mockCandidates'));
      const newCand = { ...candData, id: 'cand_' + Date.now() };
      localStorage.setItem('mockCandidates', JSON.stringify([newCand, ...cands]));
      return newCand;
    }
  },
  interviews: {
    getAll: async () => {
      await delay(500);
      return JSON.parse(localStorage.getItem('mockInterviews'));
    },
    getById: async (id) => {
      await delay(300);
      const ints = JSON.parse(localStorage.getItem('mockInterviews'));
      return ints.find(i => i.id === id);
    },
    create: async (intData) => {
      await delay(600);
      const ints = JSON.parse(localStorage.getItem('mockInterviews'));
      const newInt = { 
        ...intData, 
        id: 'int_' + Date.now(),
        status: 'Scheduled',
        link: 'https://interview.ai/room/r-' + Math.random().toString(36).substr(2,6)
      };
      localStorage.setItem('mockInterviews', JSON.stringify([newInt, ...ints]));
      return newInt;
    }
  }
};
