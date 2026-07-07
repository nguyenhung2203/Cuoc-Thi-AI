import { apiService } from './api.service';

/**
 * AI Service
 * Các API liên quan đến AI Orchestration — API_SPEC.md §9
 * Lưu ý: Frontend gọi vào Backend Go, Backend Go mới gọi tiếp sang Python AI Orchestrator.
 */
export const aiService = {
  /**
   * AI tạo danh sách câu hỏi dựa trên JD
   * API_SPEC §9.1 — POST /companies/:company_id/jobs/:job_id/ai/generate-questions
   * @param {String} companyId
   * @param {String} jobId
   * @param {Object} data { question_types, level, count }
   */
  generateQuestions: (companyId, jobId, data = {}) => {
    return apiService.post(`/companies/${companyId}/jobs/${jobId}/ai/generate-questions`, data);
  },

  /**
   * AI gợi ý câu hỏi follow-up dựa trên transcript gần nhất
   * API_SPEC §9.2 — POST /companies/:company_id/interviews/:interview_id/ai/suggest-follow-up
   * @param {String} companyId
   * @param {String} interviewId
   * @param {Object} data { last_transcript_id, focus }
   */
  suggestFollowUp: (companyId, interviewId, data = {}) => {
    return apiService.post(`/companies/${companyId}/interviews/${interviewId}/ai/suggest-follow-up`, data);
  },

  /**
   * AI chấm điểm câu trả lời của ứng viên
   * API_SPEC §9.3 — POST /companies/:company_id/interviews/:interview_id/ai/score-answer
   * @param {String} companyId
   * @param {String} interviewId
   * @param {Object} data { transcript_ids, criterion_ids }
   */
  scoreAnswer: (companyId, interviewId, data = {}) => {
    return apiService.post(`/companies/${companyId}/interviews/${interviewId}/ai/score-answer`, data);
  },

  /**
   * AI tạo báo cáo tổng hợp sau phỏng vấn (bất đồng bộ)
   * API_SPEC §9.4 — POST /companies/:company_id/interviews/:interview_id/ai/generate-report
   * Response có thể trả status: "generating" — frontend nhận event report:ready qua WebSocket khi xong
   * @param {String} companyId
   * @param {String} interviewId
   * @param {Boolean} forceRefresh
   */
  generateReport: (companyId, interviewId, forceRefresh = false) => {
    return apiService.post(`/companies/${companyId}/interviews/${interviewId}/ai/generate-report`, {
      force_refresh: forceRefresh
    });
  }
};
