import { apiService } from './api.service';

/**
 * Interview Service
 * Các API liên quan đến quản lý Lịch phỏng vấn
 */
export const interviewService = {
  /**
   * Lấy danh sách lịch phỏng vấn
   * @param {String} companyId 
   * @param {Object} params { status, date, page, page_size }
   */
  getInterviews: (companyId, params = {}) => {
    const query = new URLSearchParams(params).toString();
    const url = `/companies/${companyId}/interviews${query ? `?${query}` : ''}`;
    return apiService.get(url);
  },

  /**
   * Tạo lịch phỏng vấn mới (Lên lịch)
   * @param {String} companyId 
   * @param {Object} data { job_id, candidate_id, recruiter_id, scheduled_at, duration_minutes, mode, send_invite }
   */
  createInterview: (companyId, data) => {
    return apiService.post(`/companies/${companyId}/interviews`, data);
  },

  /**
   * Lấy chi tiết lịch phỏng vấn
   * @param {String} companyId 
   * @param {String} interviewId 
   */
  getInterview: (companyId, interviewId) => {
    return apiService.get(`/companies/${companyId}/interviews/${interviewId}`);
  },

  /**
   * Bắt đầu phỏng vấn (Start)
   * @param {String} companyId 
   * @param {String} interviewId 
   * @param {Object} data { consent_recording, consent_ai }
   */
  startInterview: (companyId, interviewId, data = { consent_recording: true, consent_ai: true }) => {
    return apiService.post(`/companies/${companyId}/interviews/${interviewId}/start`, data);
  },

  /**
   * Kết thúc phỏng vấn (End)
   * @param {String} companyId 
   * @param {String} interviewId 
   * @param {Object} data { generate_report }
   */
  endInterview: (companyId, interviewId, data = { generate_report: true }) => {
    return apiService.post(`/companies/${companyId}/interviews/${interviewId}/end`, data);
  },

  /**
   * Hủy phỏng vấn (Cancel)
   * @param {String} companyId 
   * @param {String} interviewId 
   * @param {Object} data { reason }
   */
  cancelInterview: (companyId, interviewId, data) => {
    return apiService.post(`/companies/${companyId}/interviews/${interviewId}/cancel`, data);
  },

  /**
   * Lưu ghi chú nội bộ của recruiter cho buổi phỏng vấn
   * @param {String} companyId
   * @param {String} interviewId
   * @param {String} notes
   */
  updateNotes: (companyId, interviewId, notes) => {
    return apiService.put(`/companies/${companyId}/interviews/${interviewId}/notes`, { notes });
  },

  /**
   * Gửi email nhắc nhở cho ứng viên
   * @param {String} companyId
   * @param {String} interviewId
   */
  sendReminder: (companyId, interviewId) => {
    return apiService.post(`/companies/${companyId}/interviews/${interviewId}/send-reminder`, {});
  },

  /**
   * API dành cho Candidate tham gia vào phòng (KHÔNG CẦN companyId, chỉ cần invite_token)
   * @param {String} inviteToken
   */
  joinByToken: (inviteToken) => {
    return apiService.get(`/interviews/join/${inviteToken}`);
  }
};
