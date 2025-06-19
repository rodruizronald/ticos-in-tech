import { useState, useEffect, useCallback } from 'react';
import { ApiError } from '../services/api';

/**
 * Generic hook for API calls with loading, error, and success states
 */
export const useApi = (apiFunction, initialData = null) => {
  const [data, setData] = useState(initialData);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const execute = useCallback(async (...args) => {
    try {
      setLoading(true);
      setError(null);
      
      const result = await apiFunction(...args);
      setData(result);
      return result;
    } catch (err) {
      const errorMessage = err instanceof ApiError 
        ? err.message 
        : 'An unexpected error occurred';
      
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [apiFunction]);

  const reset = useCallback(() => {
    setData(initialData);
    setError(null);
    setLoading(false);
  }, [initialData]);

  return {
    data,
    loading,
    error,
    execute,
    reset
  };
};

/**
 * Hook for API calls that should execute immediately
 */
export const useApiEffect = (apiFunction, dependencies = [], initialData = null) => {
  const { data, loading, error, execute } = useApi(apiFunction, initialData);

  useEffect(() => {
    execute();
  }, dependencies);

  return { data, loading, error, refetch: execute };
};

/**
 * Hook for debounced API calls (useful for search)
 */
export const useDebouncedApi = (apiFunction, delay = 300, initialData = null) => {
  const [data, setData] = useState(initialData);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [debouncedArgs, setDebouncedArgs] = useState(null);

  const execute = useCallback((...args) => {
    setDebouncedArgs(args);
  }, []);

  useEffect(() => {
    if (!debouncedArgs) return;

    const timeoutId = setTimeout(async () => {
      try {
        setLoading(true);
        setError(null);
        
        const result = await apiFunction(...debouncedArgs);
        setData(result);
      } catch (err) {
        const errorMessage = err instanceof ApiError 
          ? err.message 
          : 'An unexpected error occurred';
        
        setError(errorMessage);
      } finally {
        setLoading(false);
      }
    }, delay);

    return () => clearTimeout(timeoutId);
  }, [debouncedArgs, delay, apiFunction]);

  const reset = useCallback(() => {
    setData(initialData);
    setError(null);
    setLoading(false);
    setDebouncedArgs(null);
  }, [initialData]);

  return {
    data,
    loading,
    error,
    execute,
    reset
  };
};

/**
 * Hook for infinite scroll / pagination
 */
export const usePaginatedApi = (apiFunction, pageSize = 20) => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [hasMore, setHasMore] = useState(true);
  const [pagination, setPagination] = useState({
    total: 0,
    limit: pageSize,
    offset: 0
  });

  const loadMore = useCallback(async (searchParams = {}) => {
    if (loading) return;

    try {
      setLoading(true);
      setError(null);
      
      const params = {
        ...searchParams,
        limit: pageSize,
        offset: pagination.offset
      };
      
      const result = await apiFunction(params);
      
      if (pagination.offset === 0) {
        // First load or new search
        setData(result.data || []);
      } else {
        // Append to existing data
        setData(prev => [...prev, ...(result.data || [])]);
      }
      
      setPagination(result.pagination || {
        total: result.data?.length || 0,
        limit: pageSize,
        offset: pagination.offset + pageSize
      });
      
      setHasMore(result.pagination?.has_more || false);
      
      return result;
    } catch (err) {
      const errorMessage = err instanceof ApiError 
        ? err.message 
        : 'An unexpected error occurred';
      
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [apiFunction, pageSize, pagination.offset, loading]);

  const reset = useCallback(() => {
    setData([]);
    setError(null);
    setLoading(false);
    setHasMore(true);
    setPagination({
      total: 0,
      limit: pageSize,
      offset: 0
    });
  }, [pageSize]);

  const loadNext = useCallback((searchParams = {}) => {
    setPagination(prev => ({ ...prev, offset: prev.offset + pageSize }));
    return loadMore(searchParams);
  }, [loadMore, pageSize]);

  return {
    data,
    loading,
    error,
    hasMore,
    pagination,
    loadMore,
    loadNext,
    reset
  };
};