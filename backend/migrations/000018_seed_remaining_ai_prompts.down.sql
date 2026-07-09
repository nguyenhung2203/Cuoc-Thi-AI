DELETE FROM ai_prompt_templates WHERE name IN ('suggest_follow_up', 'score_answer', 'generate_report', 'mock_question', 'mock_feedback') AND version = 1 AND company_id IS NULL;
