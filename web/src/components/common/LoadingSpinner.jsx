import React from 'react';
import { Loader2 } from 'lucide-react';

/**
 * Loading spinner component with different sizes and styles
 */
const LoadingSpinner = ({ 
  size = 'md', 
  color = 'primary', 
  text = '',
  className = '' 
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8',
    xl: 'w-12 h-12'
  };

  const colorClasses = {
    primary: 'text-costa-blue',
    secondary: 'text-gray-500',
    white: 'text-white',
    danger: 'text-costa-red'
  };

  return (
    <div className={`flex items-center justify-center ${className}`}>
      <div className="flex flex-col items-center space-y-2">
        <Loader2 
          className={`animate-spin ${sizeClasses[size]} ${colorClasses[color]}`}
          aria-hidden="true"
        />
        {text && (
          <p className={`text-sm ${colorClasses[color]} animate-pulse`}>
            {text}
          </p>
        )}
      </div>
    </div>
  );
};

/**
 * Full page loading spinner
 */
export const FullPageLoader = ({ text = 'Loading...' }) => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <LoadingSpinner 
        size="xl" 
        color="primary" 
        text={text}
        className="flex-col"
      />
    </div>
  );
};

/**
 * Inline loading spinner for buttons
 */
export const ButtonSpinner = ({ size = 'sm', color = 'white' }) => {
  return (
    <LoadingSpinner 
      size={size} 
      color={color}
      className="mr-2"
    />
  );
};

/**
 * Content placeholder with skeleton loading
 */
export const SkeletonLoader = ({ 
  lines = 3, 
  className = '',
  showAvatar = false 
}) => {
  return (
    <div className={`animate-pulse ${className}`}>
      {showAvatar && (
        <div className="flex items-center mb-4">
          <div className="w-12 h-12 bg-gray-200 rounded-lg mr-4"></div>
          <div className="flex-1">
            <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div className="h-3 bg-gray-200 rounded w-1/2"></div>
          </div>
        </div>
      )}
      
      <div className="space-y-3">
        {Array.from({ length: lines }, (_, i) => (
          <div
            key={i}
            className={`h-4 bg-gray-200 rounded ${
              i === lines - 1 ? 'w-2/3' : 'w-full'
            }`}
          ></div>
        ))}
      </div>
    </div>
  );
};

/**
 * Job card skeleton loader
 */
export const JobCardSkeleton = () => {
  return (
    <div className="card p-6 animate-pulse">
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center flex-1">
          <div className="w-12 h-12 bg-gray-200 rounded-lg mr-4"></div>
          <div className="flex-1">
            <div className="h-5 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div className="h-4 bg-gray-200 rounded w-1/2"></div>
          </div>
        </div>
        <div className="w-6 h-6 bg-gray-200 rounded-full"></div>
      </div>
      
      <div className="flex items-center space-x-4 mb-3">
        <div className="h-3 bg-gray-200 rounded w-20"></div>
        <div className="h-3 bg-gray-200 rounded w-16"></div>
        <div className="h-3 bg-gray-200 rounded w-24"></div>
      </div>
      
      <div className="h-6 bg-gray-200 rounded w-20 mb-3"></div>
      
      <div className="flex flex-wrap gap-1 mb-4">
        {Array.from({ length: 4 }, (_, i) => (
          <div key={i} className="h-6 bg-gray-200 rounded-full w-16"></div>
        ))}
      </div>
      
      <div className="flex items-center justify-between">
        <div className="h-4 bg-gray-200 rounded w-24"></div>
        <div className="h-8 bg-gray-200 rounded w-20"></div>
      </div>
    </div>
  );
};

export default LoadingSpinner;