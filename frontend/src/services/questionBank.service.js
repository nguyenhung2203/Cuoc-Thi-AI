import { apiService } from './api.service';

/**
 * Question Bank Service
 * CRUD kho câu hỏi của Company — API_SPEC §K-S4-04
 * Endpoint: /companies/:company_id/question-bank
 */
export const questionBankService = {
  /**
   * Lấy danh sách câu hỏi trong kho
   * @param {String} companyId
   * @param {Object} params { job_id, type, level, keyword, page, page_size }
   */
  getQuestions: (companyId, params = {}) => {
    const query = new URLSearchParams(params).toString();
    return apiService.get(`/companies/${companyId}/question-bank${query ? `?${query}` : ''}`);
  },

  /**
   * Recruiter thêm câu hỏi thủ công
   * @param {String} companyId
   * @param {Object} data { question_text, question_type, skill_tags, level, expected_signals }
   */
  createQuestion: (companyId, data) => {
    return apiService.post(`/companies/${companyId}/question-bank`, data);
  },

  /**
   * Cập nhật câu hỏi
   * @param {String} companyId
   * @param {String} questionId
   * @param {Object} data
   */
  updateQuestion: (companyId, questionId, data) => {
    return apiService.put(`/companies/${companyId}/question-bank/${questionId}`, data);
  },

  /**
   * Xóa câu hỏi
   * @param {String} companyId
   * @param {String} questionId
   */
  deleteQuestion: (companyId, questionId) => {
    return apiService.delete(`/companies/${companyId}/question-bank/${questionId}`);
  },

  /**
   * AI tạo câu hỏi dựa trên JD và lưu vào question_bank
   * API_SPEC §9.1 — POST /companies/:company_id/jobs/:job_id/ai/generate-questions
   * @param {String} companyId
   * @param {String} jobId
   * @param {Object} data { question_types, level, count }
   */
  generateWithAI: (companyId, jobId, data = { count: 10, level: 'middle' }) => {
    return apiService.post(`/companies/${companyId}/jobs/${jobId}/ai/generate-questions`, data);
  }
};
