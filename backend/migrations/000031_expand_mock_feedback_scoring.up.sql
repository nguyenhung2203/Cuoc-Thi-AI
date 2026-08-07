-- Expand mock feedback with communication, tone, and professionalism signals.
UPDATE ai_prompt_templates
SET template = replace(
  replace(
    replace(template,
      '  "improvements": ["string"],',
      '  "improvements": ["string"],\n  "communication_score": 0,\n  "tone_score": 0,\n  "personality_score": 0,'),
    '"confidence": 0.0',
    '"confidence": 0.0'),
  'Nêu điểm tốt trước, sau đó góp ý cải thiện.',
  'Nêu điểm tốt trước, sau đó góp ý cải thiện.\n- Chấm riêng communication_score, tone_score và personality_score theo thang 0-10; chỉ dựa trên câu trả lời và không suy đoán đặc điểm nhạy cảm.')
WHERE name = 'mock_feedback' AND company_id IS NULL;
