import React, { createContext, useContext, useState, useCallback } from 'react';
import { jobsApi } from '../services/api';
import { useLocalStorageSet } from '../hooks/useLocalStorage';
import { STORAGE_KEYS } from '../utils/constants';

// Create the context
const AppContext = createContext();

/**
 * App Context Provider - manages global application state
 */
export const AppProvider = ({ children }) => {
  // Jobs state
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  
  // Search and filter state
  const [searchQuery, setSearchQuery] = useState('');
  const [filters, setFilters] = useState({
    experience_level: '',
    employment_type: '',
    location: '',
    work_mode: '',
    company: '',
    date_from: '',
    date_to: ''
  });
  
  // UI state
  const [filterSidebarOpen, setFilterSidebarOpen] = useState(false);
  
  // Pagination state
  const [pagination, setPagination] = useState({
    total: 0,
    limit: 20,
    offset: 0,
    has_more: false
  });
  
  // Saved jobs using localStorage
  const [savedJobs, addSavedJob, removeSavedJob, hasSavedJob, clearSavedJobs, toggleSavedJob] = 
    useLocalStorageSet(STORAGE_KEYS.SAVED_JOBS);

  /**
   * Search jobs with current query and filters
   */
  const searchJobs = useCallback(async (
    query = searchQuery, 
    currentFilters = filters, 
    offset = 0,
    shouldAppend = false
  ) => {
    try {
      setLoading(true);
      setError(null);
      
      const searchParams = {
        q: query.trim(),
        ...currentFilters,
        limit: pagination.limit,
        offset
      };
      
      const response = await jobsApi.searchJobs(searchParams);
      
      if (shouldAppend) {
        setJobs(prevJobs => [...prevJobs, ...(response.data || [])]);
      } else {
        setJobs(response.data || []);
      }
      
      setPagination(response.pagination || {
        total: 0,
        limit: 20,
        offset: 0,
        has_more: false
      });
      
      return response;
    } catch (err) {
      setError(err.message || 'Failed to fetch jobs');
      if (!shouldAppend) {
        setJobs([]);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [searchQuery, filters, pagination.limit]);

  /**
   * Update search query and trigger search
   */
  const updateSearchQuery = useCallback((query) => {
    setSearchQuery(query);
    searchJobs(query, filters, 0, false);
  }, [filters, searchJobs]);

  /**
   * Update filters and trigger search
   */
  const updateFilters = useCallback((newFilters) => {
    setFilters(newFilters);
    searchJobs(searchQuery, newFilters, 0, false);
  }, [searchQuery, searchJobs]);

  /**
   * Update a single filter
   */
  const updateFilter = useCallback((filterKey, value) => {
    const newFilters = { ...filters, [filterKey]: value };
    updateFilters(newFilters);
  }, [filters, updateFilters]);

  /**
   * Clear all filters
   */
  const clearFilters = useCallback(() => {
    const clearedFilters = {
      experience_level: '',
      employment_type: '',
      location: '',
      work_mode: '',
      company: '',
      date_from: '',
      date_to: ''
    };
    updateFilters(clearedFilters);
  }, [updateFilters]);

  /**
   * Load more jobs (pagination)
   */
  const loadMoreJobs = useCallback(async () => {
    if (!pagination.has_more || loading) return;
    
    const nextOffset = pagination.offset + pagination.limit;
    return searchJobs(searchQuery, filters, nextOffset, true);
  }, [pagination, loading, searchJobs, searchQuery, filters]);

  /**
   * Get a single job by ID
   */
  const getJobById = useCallback(async (jobId) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await jobsApi.getJobById(jobId);
      return response.data;
    } catch (err) {
      setError(err.message || 'Failed to fetch job details');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  /**
   * Reset all state
   */
  const resetState = useCallback(() => {
    setJobs([]);
    setSearchQuery('');
    setFilters({
      experience_level: '',
      employment_type: '',
      location: '',
      work_mode: '',
      company: '',
      date_from: '',
      date_to: ''
    });
    setError(null);
    setLoading(false);
    setPagination({
      total: 0,
      limit: 20,
      offset: 0,
      has_more: false
    });
    setFilterSidebarOpen(false);
  }, []);

  /**
   * Check if any filters are active
   */
  const hasActiveFilters = Object.values(filters).some(value => value !== '');

  // Context value
  const value = {
    // Jobs data
    jobs,
    loading,
    error,
    
    // Search state
    searchQuery,
    setSearchQuery,
    updateSearchQuery,
    
    // Filter state
    filters,
    updateFilters,
    updateFilter,
    clearFilters,
    hasActiveFilters,
    
    // UI state
    filterSidebarOpen,
    setFilterSidebarOpen,
    
    // Pagination
    pagination,
    loadMoreJobs,
    
    // Saved jobs
    savedJobs,
    addSavedJob,
    removeSavedJob,
    hasSavedJob,
    clearSavedJobs,
    toggleSavedJob,
    
    // Actions
    searchJobs,
    getJobById,
    resetState
  };

  return (
    <AppContext.Provider value={value}>
      {children}
    </AppContext.Provider>
  );
};

/**
 * Hook to use the App context
 */
export const useApp = () => {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error('useApp must be used within an AppProvider');
  }
  return context;
};