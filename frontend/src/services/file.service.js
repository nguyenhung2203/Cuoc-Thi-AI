import { apiService } from './api.service';

/**
 * File Service
 * Các API liên quan đến Upload & Download File — API_SPEC.md §12
 */
export const fileService = {
  /**
   * Upload file lên hệ thống (CV, tài liệu, v.v.)
   * API_SPEC §12.1 — POST /files (multipart/form-data)
   * @param {File} file — File object từ input
   * @param {String} fileType — 'cv' | 'document' | 'image'
   * @param {String} companyId — (optional) company_id nếu file thuộc về company
   */
  uploadFile: (file, fileType = 'cv', companyId = null) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('file_type', fileType);
    if (companyId) {
      formData.append('company_id', companyId);
    }
    return apiService.post('/files/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  },

  /**
   * Lấy Signed URL để tải file về (có thời hạn)
   * API_SPEC §12.2 — GET /files/:file_id/signed-url
   * @param {String} fileId
   * @param {String} companyId — bắt buộc khi người gọi là recruiter (BE yêu cầu)
   */
  getDownloadUrl: (fileId, companyId = null) => {
    const query = companyId ? `?company_id=${encodeURIComponent(companyId)}` : '';
    return apiService.get(`/files/${fileId}/signed-url${query}`);
  }
};
