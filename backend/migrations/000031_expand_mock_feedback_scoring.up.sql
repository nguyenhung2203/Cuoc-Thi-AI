-- Expand mock feedback with dimensions and level-aware scoring.
UPDATE ai_prompt_templates
SET template = replace(
  replace(
    replace(
      replace(template,
        'Nhiệm vụ: Nhận xét câu trả lời của ứng viên.',
        'Nhiệm vụ: Nhận xét câu trả lời của ứng viên theo đúng cấp độ target_level.\n\nCẤP ĐỘ ĐÁNH GIÁ:\n- fresher: nền tảng, tư duy học hỏi, trình bày cơ bản; không yêu cầu chuyên sâu.\n- junior: nền tảng chắc, áp dụng thực tế có hướng dẫn, xử lý tình huống phổ biến.\n- mid: tự chủ, phân tích trade-off, kinh nghiệm thực tế và trách nhiệm.\nChỉ dùng tiêu chuẩn của cấp độ được cung cấp.'),
        '  "improvements": ["string"],',
        '  "improvements": ["string"],\n  "communication_score": 0,\n  "tone_score": 0,\n  "personality_score": 0,'),
      'TARGET ROLE:\n{{target_role}}',
      'TARGET ROLE:\n{{target_role}}\n\nTARGET LEVEL:\n{{target_level}}'),
    '"target_role":"string"}',
    '"target_role":"string","target_level":"string"}')
WHERE name = 'mock_feedback' AND company_id IS NULL;
