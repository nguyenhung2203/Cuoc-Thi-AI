import { apiService } from './api.service';

/**
 * Room Service
 * Các API liên quan đến phòng phỏng vấn (Realtime Room)
 */
export const roomService = {
  /**
   * Lấy chi tiết phòng (Metadata)
   * @param {String} companyId 
   * @param {String} interviewId 
   */
  getRoomDetail: (companyId, interviewId) => {
    return apiService.get(`/companies/${companyId}/interviews/${interviewId}/room`);
  },

  /**
   * Xin cấp Token LiveKit cho phòng (dành cho Recruiter)
   * Candidate sẽ nhận token từ hàm joinByToken trong interview.service
   * @param {String} companyId 
   * @param {String} interviewId 
   * @param {Object} data { participant_type } (mặc định "recruiter")
   */
  getRoomToken: (companyId, interviewId, data = { participant_type: 'recruiter' }) => {
    return apiService.get(`/companies/${companyId}/interviews/${interviewId}/room/access-token`);
  }
};
