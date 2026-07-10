import { apiService } from './api.service';

export const templateService = {
  getTemplates: (companyId) => {
    return apiService.get(`/companies/${companyId}/templates`);
  },

  createTemplate: (companyId, payload) => {
    return apiService.post(`/companies/${companyId}/templates`, payload);
  },

  updateTemplate: (companyId, templateId, payload) => {
    return apiService.put(`/companies/${companyId}/templates/${templateId}`, payload);
  },

  deleteTemplate: (companyId, templateId) => {
    return apiService.delete(`/companies/${companyId}/templates/${templateId}`);
  }
};
