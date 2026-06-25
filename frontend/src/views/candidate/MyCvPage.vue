<script setup>
import Button from '../../components/common/AppButton.vue'
import { Upload } from 'lucide-vue-next'
</script>

<template>
  <div class="page-header">
    <div>
      <h1 class="text-h1">CV của tôi</h1>
      <p class="text-helper" style="margin-top: 4px">Tải lên và quản lý hồ sơ xin việc (CV/Resume) của bạn.</p>
    </div>
  </div>
  
  <div style="padding: 48px; text-align: center; border: 1px dashed var(--border); border-radius: var(--radius-lg); background-color: var(--surface-soft)">
    <div style="margin-bottom: 16px; display: inline-flex; justify-content: center; align-items: center; width: 64px; height: 64px; border-radius: 50%; background-color: rgba(37, 99, 235, 0.1); color: var(--primary)">
      <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" x2="8" y1="13" y2="13"/><line x1="16" x2="8" y1="17" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
    </div>
    <h2 class="text-h2" style="margin-bottom: 8px">Chưa có CV nào</h2>
    <p class="text-body" style="color: var(--text-secondary); max-width: 400px; margin: 0 auto 24px">
      Tải lên CV của bạn (định dạng PDF) để các nhà tuyển dụng dễ dàng đánh giá kỹ năng của bạn hơn.
    </p>
    <Button><Upload size="16" style="margin-right: 8px" /> Tải lên CV mới</Button>
  </div>
</template>
