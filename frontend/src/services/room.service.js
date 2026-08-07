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
   * @param {String} companyId
   * @param {String} interviewId
   */
  getRoomToken: (companyId, interviewId) => {
    return apiService.get(`/companies/${companyId}/interviews/${interviewId}/room/access-token`);
  },

  /**
   * Lấy lịch sử chat của phòng (REST fallback — chủ yếu dùng WebSocket)
   * API_SPEC §7 — GET /api/v1/rooms/:room_id/chat
   * ⚠️ Endpoint này nằm trên Realtime Gateway (VITE_WS_URL, cổng 8081),
   *    KHÔNG phải REST API server (cổng 8080), nên không dùng apiService.
   * @param {String} roomId
   */
  getChatHistory: async (roomId) => {
    const wsBase = import.meta.env.VITE_WS_URL || `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}`;
    // VITE_WS_URL có thể kèm path (vd ws://host:8081/ws/interview-room);
    // chỉ giữ lại origin (scheme+host+port) và đổi ws(s):// -> http(s)://
    const httpScheme = wsBase.replace(/^ws(s)?:\/\//, (_, s) => (s ? 'https://' : 'http://'));
    const httpBase = new URL(httpScheme).origin;
    const token = localStorage.getItem('access_token');
    const headers = {};
    if (token) headers['Authorization'] = `Bearer ${token}`;

    const response = await fetch(`${httpBase}/api/v1/rooms/${roomId}/chat`, { headers });
    const data = await response.json();
    if (!response.ok || data.success === false) {
      throw data.error || { message: 'Không lấy được lịch sử chat.' };
    }
    return data.data;
  }
};
