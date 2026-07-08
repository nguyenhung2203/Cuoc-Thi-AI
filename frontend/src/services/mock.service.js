import { apiService } from './api.service';

/**
 * Mock Interview Service
 * Quản lý luyện tập phỏng vấn với AI — API_SPEC.md §11
 */
export const mockService = {
  /**
   * Tạo Mock Interview mới (draft)
   * API_SPEC §11.1 — POST /mock-interviews
   * @param {Object} data { target_role, target_level, cv_file_id }
   */
  createMockInterview: (data) => {
    return apiService.post('/mock-interviews', data);
  },

  /**
   * Bắt đầu Mock Interview (lấy câu hỏi đầu tiên)
   * API_SPEC §11.2 — POST /mock-interviews/:mock_interview_id/start
   * @param {String} mockId
   */
  startMockInterview: (mockId) => {
    return apiService.post(`/mock-interviews/${mockId}/start`);
  },

  /**
   * Gửi câu trả lời và nhận feedback + câu hỏi tiếp theo từ AI
   * API_SPEC §11.3 — POST /mock-interviews/:mock_interview_id/answer
   * @param {String} mockId
   * @param {Object} data { question_id, answer_text }
   */
  submitAnswer: (mockId, data) => {
    // API_SPEC §11.3 — POST /mock-interviews/:mock_interview_id/messages
    return apiService.post(`/mock-interviews/${mockId}/messages`, { content: data.answer_text });
  },

  /**
   * Kết thúc Mock Interview và lấy kết quả tổng hợp
   * API_SPEC §11.4 — POST /mock-interviews/:mock_interview_id/end
   * @param {String} mockId
   */
  endMockInterview: (mockId) => {
    return apiService.post(`/mock-interviews/${mockId}/end`);
  },

  /**
   * Lấy lịch sử các Mock Interview của Candidate hiện tại
   * API_SPEC §11.5 — GET /mock-interviews/my
   * @param {Object} params { page, page_size }
   */
  listMyMockInterviews: (params = {}) => {
    const query = new URLSearchParams(params).toString();
    // Backend API uses /mock-interviews/me instead of /my
    return apiService.get(`/mock-interviews/me${query ? `?${query}` : ''}`);
  },

  /**
   * Lấy chi tiết một Mock Interview (kết quả, messages)
   * @param {String} mockId
   */
  getMockInterview: (mockId) => {
    return apiService.get(`/mock-interviews/${mockId}`);
  },

  /**
   * Lấy báo cáo kết quả Mock Interview
   * API_SPEC §11.7 — GET /mock-interviews/:mock_interview_id/report
   * @param {String} mockId
   */
  getReport: (mockId) => {
    return apiService.get(`/mock-interviews/${mockId}/report`);
  },

  /**
   * Lấy lịch sử messages của một Mock Interview
   * @param {String} mockId
   */
  getMessages: (mockId) => {
    return apiService.get(`/mock-interviews/${mockId}/messages`);
  },

  /**
   * [Legacy] Gửi tin nhắn - alias cho submitAnswer với format cũ
   * Dùng sendMessage nếu backend vẫn hỗ trợ endpoint /messages
   * @param {String} mockId
   * @param {String} content
   */
  sendMessage: (mockId, content) => {
    return apiService.post(`/mock-interviews/${mockId}/messages`, { content });
  },

  /**
   * [Legacy] Kết thúc interview - alias cho endMockInterview
   * @param {String} mockId
   */
  finishMockInterview: (mockId) => {
    return apiService.post(`/mock-interviews/${mockId}/finish`);
  }
};
