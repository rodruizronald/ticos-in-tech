// API Configuration
export const API_BASE_URL = 'http://localhost:8080/api/v1';
export const API_ENDPOINTS = {
  JOBS: '/jobs',
};

// Filter Options
export const EXPERIENCE_LEVELS = [
  'Entry-level',
  'Junior', 
  'Mid-level',
  'Senior',
  'Lead',
  'Principal',
  'Executive'
];

export const EMPLOYMENT_TYPES = [
  'Full-time',
  'Part-time', 
  'Contract',
  'Freelance',
  'Temporary',
  'Internship'
];

export const WORK_MODES = [
  'Remote',
  'Hybrid',
  'Onsite'
];

export const LOCATIONS = [
  'Costa Rica',
  'LATAM'
];

// Pagination
export const DEFAULT_PAGE_SIZE = 20;
export const MAX_PAGE_SIZE = 100;

// UI Constants
export const DEBOUNCE_DELAY = 300;
export const ANIMATION_DURATION = {
  FAST: 150,
  NORMAL: 200,
  SLOW: 300
};

// Routes
export const ROUTES = {
  HOME: '/',
  JOB_DETAIL: '/job/:id',
  NOT_FOUND: '/404'
};

// Local Storage Keys
export const STORAGE_KEYS = {
  SAVED_JOBS: 'ticos_saved_jobs',
  USER_PREFERENCES: 'ticos_user_preferences'
};