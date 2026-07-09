DELETE FROM ai_prompt_templates WHERE name IN ('analyze_jd', 'analyze_cv', 'generate_questions') AND version = 1 AND company_id IS NULL;
