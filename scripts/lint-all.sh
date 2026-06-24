#!/bin/bash
echo "Linting frontend..."
npm run lint --prefix frontend
echo "Linting backend..."
go fmt ./backend/...
echo "Linting AI service..."
flake8 ai-service
