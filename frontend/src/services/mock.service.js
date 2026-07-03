import { apiService } from './api.service';

/**
 * Mock Service
 * Quản lý luyện tập phỏng vấn với AI
 */
export const mockService = {
  /**
   * Tạo Mock Interview mới
   */
  startMockInterview: async (data) => {
    // data: { target_role, target_level, cv_file_id }
    return apiService.post('/mock-interviews', data);
  },

  /**
   * Lấy lịch sử chat
   */
  getMessages: async (mockId) => {
    return apiService.get(`/mock-interviews/${mockId}/messages`);
  },

  /**
   * Gửi tin nhắn trả lời và nhận feedback/câu hỏi tiếp theo từ AI
   */
  sendMessage: async (mockId, content) => {
    return apiService.post(`/mock-interviews/${mockId}/messages`, { content });
  },

  /**
   * Kết thúc phỏng vấn và lấy kết quả cuối cùng
   */
  finishMockInterview: async (mockId) => {
    return apiService.post(`/mock-interviews/${mockId}/finish`);
  }
};
