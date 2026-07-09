import { apiService } from './api.service';

/**
 * Transcript Service
 * Các API liên quan đến Transcript cuộc phỏng vấn — API_SPEC.md §8
 * Lưu ý: Transcript realtime nhận qua WebSocket (event transcript:update).
 * Các API dưới đây dùng để lấy lịch sử hoặc chỉnh sửa thủ công.
 */
export const transcriptService = {
  /**
   * Lấy toàn bộ transcript của một buổi phỏng vấn
   * API_SPEC §8.1 — GET /companies/:company_id/interviews/:interview_id/transcripts
   * @param {String} companyId
   * @param {String} interviewId
   */
  getTranscripts: (companyId, interviewId) => {
    return apiService.get(`/companies/${companyId}/interviews/${interviewId}/transcripts`);
  },

  /**
   * Thêm transcript thủ công (note-like)
   * API_SPEC §8.2 — POST /companies/:company_id/interviews/:interview_id/transcripts
   * @param {String} companyId
   * @param {String} interviewId
   * @param {Object} data { speaker_type, speaker_name, content, source }
   */
  createTranscript: (companyId, interviewId, data) => {
    return apiService.post(`/companies/${companyId}/interviews/${interviewId}/transcripts`, data);
  },

  /**
   * Chỉnh sửa nội dung một transcript (giữ nguyên bản gốc, lưu edited_content)
   * API_SPEC §8.3 — PUT /companies/:company_id/interviews/:interview_id/transcripts/:transcript_id
   * @param {String} companyId
   * @param {String} interviewId
   * @param {String} transcriptId
   * @param {Object} data { edited_content }
   */
  editTranscript: (companyId, interviewId, transcriptId, data) => {
    return apiService.put(
      `/companies/${companyId}/interviews/${interviewId}/transcripts/${transcriptId}`,
      data
    );
  }
};
