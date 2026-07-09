import { apiService } from './api.service';

/**
 * Rubric Service
 * CRUD Rubric & Rubric Criteria — API_SPEC §K-S5-01
 * Endpoint: /companies/:company_id/rubrics
 */
export const rubricService = {
  /**
   * Lấy danh sách Rubric của company
   * @param {String} companyId
   */
  getRubrics: (companyId) => {
    return apiService.get(`/companies/${companyId}/rubrics`);
  },

  /**
   * Lấy chi tiết Rubric kèm các Criteria
   * @param {String} companyId
   * @param {String} rubricId
   */
  getRubric: (companyId, rubricId) => {
    return apiService.get(`/companies/${companyId}/rubrics/${rubricId}`);
  },

  /**
   * Tạo Rubric mới kèm criteria
   * @param {String} companyId
   * @param {Object} data { name, description, criteria: [{ name, weight, min_score, max_score, scoring_guide }] }
   */
  createRubric: (companyId, data) => {
    return apiService.post(`/companies/${companyId}/rubrics`, data);
  },

  /**
   * Cập nhật Rubric
   * @param {String} companyId
   * @param {String} rubricId
   * @param {Object} data
   */
  updateRubric: (companyId, rubricId, data) => {
    return apiService.put(`/companies/${companyId}/rubrics/${rubricId}`, data);
  },

  /**
   * Xóa Rubric
   * @param {String} companyId
   * @param {String} rubricId
   */
  deleteRubric: (companyId, rubricId) => {
    return apiService.delete(`/companies/${companyId}/rubrics/${rubricId}`);
  }
};
