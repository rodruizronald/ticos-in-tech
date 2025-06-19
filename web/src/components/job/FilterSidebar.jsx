import React, { useEffect } from 'react';
import { X, Filter } from 'lucide-react';
import { useApp } from '../../context/AppContext';
import { 
  EXPERIENCE_LEVELS, 
  EMPLOYMENT_TYPES, 
  WORK_MODES 
} from '../../utils/constants';

/**
 * Filter sidebar component for job filtering
 */
const FilterSidebar = () => {
  const { 
    filters, 
    updateFilter, 
    clearFilters, 
    filterSidebarOpen, 
    setFilterSidebarOpen,
    hasActiveFilters,
    loading
  } = useApp();

  // Close sidebar on escape key
  useEffect(() => {
    const handleEscape = (e) => {
      if (e.key === 'Escape') {
        setFilterSidebarOpen(false);
      }
    };

    if (filterSidebarOpen) {
      document.addEventListener('keydown', handleEscape);
      document.body.style.overflow = 'hidden';
    }

    return () => {
      document.removeEventListener('keydown', handleEscape);
      document.body.style.overflow = 'unset';
    };
  }, [filterSidebarOpen, setFilterSidebarOpen]);

  const handleFilterChange = (filterType, value) => {
    updateFilter(filterType, value);
  };

  const handleCheckboxChange = (filterType, value, isChecked) => {
    // For single-select behavior (like radio buttons)
    updateFilter(filterType, isChecked ? value : '');
  };

  if (!filterSidebarOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div 
        className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
        onClick={() => setFilterSidebarOpen(false)}
        aria-hidden="true"
      />
      
      {/* Sidebar */}
      <div className="fixed inset-y-0 right-0 w-full max-w-sm bg-white shadow-strong z-50 overflow-y-auto animate-slide-in-right lg:relative lg:shadow-none lg:animate-none">
        <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 lg:px-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Filter className="w-5 h-5 text-gray-600" />
              <h3 className="text-lg font-semibold text-gray-900">Filters</h3>
              {hasActiveFilters && (
                <span className="badge badge-primary text-xs">
                  Active
                </span>
              )}
            </div>
            
            <div className="flex items-center space-x-3">
              {hasActiveFilters && (
                <button 
                  onClick={clearFilters}
                  disabled={loading}
                  className="text-costa-blue hover:text-blue-700 text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Clear All
                </button>
              )}
              
              <button 
                onClick={() => setFilterSidebarOpen(false)}
                className="p-1 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-costa-blue"
                aria-label="Close filters"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>

        <div className="px-6 py-6 space-y-8 lg:px-4">
          {/* Experience Level */}
          <FilterSection
            title="Experience Level"
            value={filters.experience_level}
            options={EXPERIENCE_LEVELS}
            onChange={(value) => handleFilterChange('experience_level', value)}
            type="radio"
            loading={loading}
          />

          {/* Employment Type */}
          <FilterSection
            title="Employment Type"
            value={filters.employment_type}
            options={EMPLOYMENT_TYPES}
            onChange={(value) => handleFilterChange('employment_type', value)}
            type="radio"
            loading={loading}
          />

          {/* Work Mode */}
          <FilterSection
            title="Work Mode"
            value={filters.work_mode}
            options={WORK_MODES}
            onChange={(value) => handleFilterChange('work_mode', value)}
            type="radio"
            loading={loading}
          />

          {/* Company Search */}
          <div>
            <h4 className="font-medium text-gray-900 mb-3">Company</h4>
            <input
              type="text"
              placeholder="Search companies..."
              value={filters.company}
              onChange={(e) => handleFilterChange('company', e.target.value)}
              disabled={loading}
              className="input w-full disabled:opacity-50 disabled:cursor-not-allowed"
              aria-label="Filter by company name"
            />
            {filters.company && (
              <p className="text-xs text-gray-500 mt-1">
                Press Enter to search
              </p>
            )}
          </div>

          {/* Date Range - Future enhancement */}
          <div className="opacity-50">
            <h4 className="font-medium text-gray-900 mb-3">Posted Date</h4>
            <select 
              className="input w-full cursor-not-allowed"
              disabled
              title="Coming soon"
            >
              <option>Any time</option>
              <option>Last 24 hours</option>
              <option>Last 7 days</option>
              <option>Last 30 days</option>
            </select>
          </div>
        </div>

        {/* Mobile apply button */}
        <div className="sticky bottom-0 bg-white border-t border-gray-200 px-6 py-4 lg:hidden">
          <button
            onClick={() => setFilterSidebarOpen(false)}
            className="btn btn-primary w-full"
            disabled={loading}
          >
            Show Results
          </button>
        </div>
      </div>
    </>
  );
};

/**
 * Filter section component for reusable filter groups
 */
const FilterSection = ({ 
  title, 
  value, 
  options, 
  onChange, 
  type = 'radio', 
  loading = false 
}) => {
  return (
    <div>
      <h4 className="font-medium text-gray-900 mb-3">{title}</h4>
      <div className="space-y-3">
        {options.map((option) => (
          <label 
            key={option} 
            className={`flex items-center cursor-pointer group ${
              loading ? 'opacity-50 cursor-not-allowed' : ''
            }`}
          >
            <input
              type={type}
              name={title.toLowerCase().replace(' ', '_')}
              checked={value === option}
              onChange={(e) => onChange(e.target.checked ? option : '')}
              disabled={loading}
              className={`rounded border-gray-300 text-costa-blue focus:ring-costa-blue focus:ring-offset-0 focus:ring-2 ${
                type === 'radio' ? 'rounded-full' : ''
              } disabled:opacity-50 disabled:cursor-not-allowed`}
            />
            <span className="ml-3 text-sm text-gray-700 group-hover:text-gray-900 select-none">
              {option}
            </span>
          </label>
        ))}
      </div>
    </div>
  );
};

export default FilterSidebar;
    