import { apiService } from './api.service';

/**
 * Room Service
 * Các API liên quan đến phòng phỏng vấn (Realtime Room) — API_SPEC.md §7
 */
export const roomService = {
  /**
   * Lấy metadata chi tiết phòng
   * API_SPEC §7.1 — GET /companies/:company_id/interviews/:interview_id/room
   */
  getRoomDetail: (companyId, interviewId) => {
    return apiService.get(`/companies/${companyId}/interviews/${interviewId}/room`);
  },

  /**
   * Xin cấp Room Access Token cho Recruiter (LiveKit JWT)
   * API_SPEC §7.2 — POST /companies/:company_id/interviews/:interview_id/room/token
   * ⚠️ Đây phải là POST, không phải GET
   * @param {String} companyId
   * @param {String} interviewId
   * @param {Object} data { participant_type: 'recruiter' }
   */
  getRoomToken: (companyId, interviewId, data = { participant_type: 'recruiter' }) => {
    return apiService.post(`/companies/${companyId}/interviews/${interviewId}/room/token`, data);
  },

  /**
   * Lấy lịch sử chat của phòng (REST fallback — chủ yếu dùng WebSocket)
   * API_SPEC §7 (endpoint GET /api/v1/rooms/:room_id/chat từ Gateway của Hùng)
   * @param {String} roomId
   */
  getChatHistory: (roomId) => {
    return apiService.get(`/rooms/${roomId}/chat`);
  }
};
