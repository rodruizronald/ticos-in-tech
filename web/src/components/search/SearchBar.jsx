import React, { useState, useCallback } from 'react';
import { Search, MapPin } from 'lucide-react';
import { useApp } from '../../context/AppContext';
import { debounce } from '../../utils/formatters';

/**
 * Search bar component for job search
 */
const SearchBar = ({ 
  size = 'lg', 
  showLocation = true, 
  placeholder = 'Search jobs...', 
  className = '' 
}) => {
  const { searchQuery, updateSearchQuery, filters, updateFilter } = useApp();
  const [localQuery, setLocalQuery] = useState(searchQuery);
  const [localLocation, setLocalLocation] = useState(filters.location || 'Costa Rica');

  // Debounced search to avoid too many API calls
  const debouncedSearch = useCallback(
    debounce((query) => {
      updateSearchQuery(query);
    }, 300),
    [updateSearchQuery]
  );

  const handleQueryChange = (e) => {
    const value = e.target.value;
    setLocalQuery(value);
    debouncedSearch(value);
  };

  const handleLocationChange = (e) => {
    const value = e.target.value;
    setLocalLocation(value);
    updateFilter('location', value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    updateSearchQuery(localQuery);
    if (showLocation) {
      updateFilter('location', localLocation);
    }
  };

  const sizeClasses = {
    sm: 'h-10',
    md: 'h-12',
    lg: 'h-14'
  };

  const inputSizeClasses = {
    sm: 'text-sm pl-10 pr-4',
    md: 'text-base pl-12 pr-4',
    lg: 'text-lg pl-14 pr-6'
  };

  const iconSizeClasses = {
    sm: 'w-4 h-4 left-3',
    md: 'w-5 h-5 left-4',
    lg: 'w-6 h-6 left-4'
  };

  return (
    <form onSubmit={handleSubmit} className={`w-full ${className}`}>
      <div className={`flex flex-col sm:flex-row bg-white rounded-xl shadow-medium overflow-hidden border border-gray-200 focus-within:ring-2 focus-within:ring-costa-blue focus-within:border-transparent ${sizeClasses[size]}`}>
        {/* Search input */}
        <div className="flex-1 relative">
          <Search className={`absolute top-1/2 transform -translate-y-1/2 text-gray-400 ${iconSizeClasses[size]}`} />
          <input
            type="text"
            placeholder={placeholder}
            value={localQuery}
            onChange={handleQueryChange}
            className={`w-full ${inputSizeClasses[size]} py-0 h-full text-gray-900 placeholder-gray-500 focus:outline-none border-0`}
            aria-label="Search for jobs"
          />
        </div>

        {/* Location input */}
        {showLocation && (
          <div className="flex-1 relative border-l border-gray-200">
            <MapPin className={`absolute top-1/2 transform -translate-y-1/2 text-gray-400 ${iconSizeClasses[size]}`} />
            <input
              type="text"
              placeholder="Location in Costa Rica"
              value={localLocation}
              onChange={handleLocationChange}
              className={`w-full ${inputSizeClasses[size]} py-0 h-full text-gray-900 placeholder-gray-500 focus:outline-none border-0`}
              aria-label="Job location"
            />
          </div>
        )}

        {/* Search button */}
        <button
          type="submit"
          className={`bg-costa-blue hover:bg-blue-700 text-white font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-costa-blue focus:ring-offset-2 px-6 sm:px-8 ${
            size === 'sm' ? 'text-sm' : size === 'lg' ? 'text-lg' : 'text-base'
          }`}
          aria-label="Search jobs"
        >
          <span className="hidden sm:inline">Search Jobs</span>
          <Search className="sm:hidden w-5 h-5" />
        </button>
      </div>
    </form>
  );
};

/**
 * Compact search bar for use in headers or secondary locations
 */
export const CompactSearchBar = ({ className = '' }) => {
  return (
    <SearchBar 
      size="sm"
      showLocation={false}
      placeholder="Search..."
      className={className}
    />
  );
};

export default SearchBar;