import React from 'react';
import { Filter, ChevronDown } from 'lucide-react';
import { useApp } from '../../context/AppContext';
import JobCard from './JobCard';
import FilterSidebar from './FilterSidebar';
import { JobCardSkeleton } from '../common/LoadingSpinner';
import { formatJobCount } from '../../utils/formatters';

/**
 * Job grid component with filtering and pagination
 */
const JobGrid = () => {
  const { 
    jobs, 
    loading, 
    error, 
    pagination, 
    filterSidebarOpen, 
    setFilterSidebarOpen,
    loadMoreJobs,
    hasActiveFilters,
    searchQuery
  } = useApp();

  const handleLoadMore = () => {
    if (!loading && pagination.has_more) {
      loadMoreJobs();
    }
  };

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <ErrorState error={error} />
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8" id="job-listings">
      <div className="lg:flex lg:space-x-8">
        {/* Desktop Filter Sidebar */}
        <div className="hidden lg:block w-80 flex-shrink-0">
          <div className="sticky top-24">
            <FilterSidebar />
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1 min-w-0">
          {/* Header with job count and filter button */}
          <JobGridHeader 
            totalJobs={pagination.total}
            currentCount={jobs.length}
            loading={loading}
            searchQuery={searchQuery}
            hasActiveFilters={hasActiveFilters}
            onToggleFilters={() => setFilterSidebarOpen(!filterSidebarOpen)}
          />

          {/* Job Grid */}
          <div className="space-y-6">
            {/* Loading skeletons */}
            {loading && jobs.length === 0 && (
              <div className="grid gap-6 md:grid-cols-2">
                {Array.from({ length: 6 }, (_, i) => (
                  <JobCardSkeleton key={i} />
                ))}
              </div>
            )}

            {/* No results state */}
            {!loading && jobs.length === 0 && (
              <EmptyState searchQuery={searchQuery} hasActiveFilters={hasActiveFilters} />
            )}

            {/* Job cards */}
            {jobs.length > 0 && (
              <>
                <div className="grid gap-6 md:grid-cols-2">
                  {jobs.map((job) => (
                    <JobCard key={job.job_id} job={job} />
                  ))}
                </div>

                {/* Load more button */}
                {pagination.has_more && (
                  <div className="flex justify-center pt-8">
                    <button
                      onClick={handleLoadMore}
                      disabled={loading}
                      className="btn btn-secondary btn-lg flex items-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {loading ? (
                        <>
                          <div className="w-4 h-4 border-2 border-gray-400 border-t-transparent rounded-full animate-spin" />
                          <span>Loading...</span>
                        </>
                      ) : (
                        <>
                          <span>Load More Jobs</span>
                          <ChevronDown className="w-4 h-4" />
                        </>
                      )}
                    </button>
                  </div>
                )}

                {/* End of results indicator */}
                {!pagination.has_more && jobs.length > 6 && (
                  <div className="text-center pt-8 pb-4">
                    <p className="text-gray-500">
                      You've seen all {formatJobCount(pagination.total)} available jobs
                    </p>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>

      {/* Mobile Filter Sidebar */}
      <FilterSidebar />
    </div>
  );
};

/**
 * Job grid header with counts and filter toggle
 */
const JobGridHeader = ({ 
  totalJobs, 
  currentCount, 
  loading, 
  searchQuery, 
  hasActiveFilters, 
  onToggleFilters 
}) => {
  return (
    <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-8 pb-6 border-b border-gray-200">
      <div>
        <h2 className="text-2xl font-semibold text-gray-900 mb-2">
          {searchQuery ? `Search Results` : 'Latest Tech Jobs'}
        </h2>
        
        <div className="flex flex-wrap items-center gap-2 text-gray-600">
          {loading && currentCount === 0 ? (
            <span>Loading jobs...</span>
          ) : (
            <>
              <span>
                Showing {currentCount} of {formatJobCount(totalJobs)}
              </span>
              {searchQuery && (
                <>
                  <span>•</span>
                  <span>for "{searchQuery}"</span>
                </>
              )}
              {hasActiveFilters && (
                <>
                  <span>•</span>
                  <span className="text-costa-blue font-medium">Filtered</span>
                </>
              )}
            </>
          )}
        </div>
      </div>

      <div className="flex items-center space-x-4 mt-4 sm:mt-0">
        {/* Future: Sort dropdown */}
        <select 
          className="input text-sm opacity-50 cursor-not-allowed"
          disabled
          title="Coming soon"
        >
          <option>Most Recent</option>
          <option>Most Relevant</option>
          <option>Salary: High to Low</option>
        </select>

        {/* Filter toggle button */}
        <button
          onClick={onToggleFilters}
          className={`btn btn-secondary flex items-center space-x-2 lg:hidden ${
            hasActiveFilters ? 'ring-2 ring-costa-blue' : ''
          }`}
        >
          <Filter className="w-4 h-4" />
          <span>Filters</span>
          {hasActiveFilters && (
            <span className="w-2 h-2 bg-costa-blue rounded-full" />
          )}
        </button>
      </div>
    </div>
  );
};

/**
 * Empty state when no jobs are found
 */
const EmptyState = ({ searchQuery, hasActiveFilters }) => {
  const { clearFilters, updateSearchQuery } = useApp();

  return (
    <div className="text-center py-16">
      <div className="max-w-md mx-auto">
        <div className="w-24 h-24 mx-auto mb-6 bg-gray-100 rounded-full flex items-center justify-center">
          <Filter className="w-12 h-12 text-gray-400" />
        </div>
        
        <h3 className="text-xl font-semibold text-gray-900 mb-2">
          No jobs found
        </h3>
        
        <p className="text-gray-600 mb-6">
          {searchQuery || hasActiveFilters 
            ? "We couldn't find any jobs matching your criteria. Try adjusting your search or filters."
            : "There are no jobs available at the moment. Check back later for new opportunities."
          }
        </p>

        <div className="flex flex-col sm:flex-row gap-3 justify-center">
          {hasActiveFilters && (
            <button 
              onClick={clearFilters}
              className="btn btn-secondary"
            >
              Clear Filters
            </button>
          )}
          
          {searchQuery && (
            <button 
              onClick={() => updateSearchQuery('')}
              className="btn btn-primary"
            >
              Show All Jobs
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

/**
 * Error state component
 */
const ErrorState = ({ error }) => {
  const handleRetry = () => {
    window.location.reload();
  };

  return (
    <div className="text-center py-16">
      <div className="max-w-md mx-auto">
        <div className="w-24 h-24 mx-auto mb-6 bg-red-100 rounded-full flex items-center justify-center">
          <svg className="w-12 h-12 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        
        <h3 className="text-xl font-semibold text-gray-900 mb-2">
          Something went wrong
        </h3>
        
        <p className="text-gray-600 mb-6">
          {error || 'We encountered an error while loading jobs. Please try again.'}
        </p>

        <button 
          onClick={handleRetry}
          className="btn btn-primary"
        >
          Try Again
        </button>
      </div>
    </div>
  );
};

export default JobGrid;