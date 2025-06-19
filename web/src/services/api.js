import { API_BASE_URL, API_ENDPOINTS } from '../utils/constants';

/**
 * Base API class for handling HTTP requests
 */
class ApiService {
  constructor(baseURL = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  /**
   * Generic HTTP request method
   */
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new ApiError(
          errorData.message || `HTTP error! status: ${response.status}`,
          response.status,
          errorData
        );
      }

      return await response.json();
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      throw new ApiError(
        'Network error occurred. Please check your connection.',
        0,
        error
      );
    }
  }

  /**
   * GET request
   */
  async get(endpoint, params = {}) {
    const searchParams = new URLSearchParams();
    
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        searchParams.append(key, value);
      }
    });

    const queryString = searchParams.toString();
    const fullEndpoint = queryString ? `${endpoint}?${queryString}` : endpoint;

    return this.request(fullEndpoint, { method: 'GET' });
  }

  /**
   * POST request
   */
  async post(endpoint, data) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  /**
   * PUT request
   */
  async put(endpoint, data) {
    return this.request(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  /**
   * DELETE request
   */
  async delete(endpoint) {
    return this.request(endpoint, { method: 'DELETE' });
  }
}

/**
 * Custom error class for API errors
 */
export class ApiError extends Error {
  constructor(message, status, details = null) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.details = details;
  }

  /**
   * Check if error is a network error
   */
  isNetworkError() {
    return this.status === 0;
  }

  /**
   * Check if error is a client error (4xx)
   */
  isClientError() {
    return this.status >= 400 && this.status < 500;
  }

  /**
   * Check if error is a server error (5xx)
   */
  isServerError() {
    return this.status >= 500;
  }
}

/**
 * Jobs API service
 */
export class JobsApi extends ApiService {
  /**
   * Search jobs with filters and pagination
   */
  async searchJobs(searchParams = {}) {
    const {
      q = '',
      experience_level = '',
      employment_type = '',
      location = '',
      work_mode = '',
      company = '',
      date_from = '',
      date_to = '',
      limit = 20,
      offset = 0,
    } = searchParams;

    const params = {
      q,
      experience_level,
      employment_type,
      location,
      work_mode,
      company,
      date_from,
      date_to,
      limit,
      offset,
    };

    return this.get(API_ENDPOINTS.JOBS, params);
  }

  /**
   * Get a single job by ID
   * Note: This would be a separate endpoint in a real API
   */
  async getJobById(jobId) {
    // For MVP, we'll simulate this by searching and filtering
    // In a real API, this would be: return this.get(`${API_ENDPOINTS.JOBS}/${jobId}`);
    
    try {
      const response = await this.searchJobs({ q: 'developer', limit: 100 });
      const job = response.data.find(job => job.job_id === parseInt(jobId));
      
      if (!job) {
        throw new ApiError('Job not found', 404);
      }
      
      return { data: job };
    } catch (error) {
      throw error;
    }
  }
}

// Create and export a singleton instance
export const jobsApi = new JobsApi();

// Export the base API service for other potential APIs
export default ApiService;