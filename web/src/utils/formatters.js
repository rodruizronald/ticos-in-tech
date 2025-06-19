/**
 * Format a date string to show relative time (e.g., "Posted 3 days ago")
 * @param {string} dateString - ISO date string
 * @returns {string} Formatted relative date
 */
export const formatRelativeDate = (dateString) => {
  const date = new Date(dateString);
  const now = new Date();
  const diffTime = Math.abs(now - date);
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  
  if (diffDays === 0) return 'Posted today';
  if (diffDays === 1) return 'Posted 1 day ago';
  if (diffDays < 7) return `Posted ${diffDays} days ago`;
  if (diffDays < 30) {
    const weeks = Math.floor(diffDays / 7);
    return weeks === 1 ? 'Posted 1 week ago' : `Posted ${weeks} weeks ago`;
  }
  
  const months = Math.floor(diffDays / 30);
  return months === 1 ? 'Posted 1 month ago' : `Posted ${months} months ago`;
};

/**
 * Get company initials for logo fallback
 * @param {string} companyName - Company name
 * @returns {string} Company initials (max 2 characters)
 */
export const getCompanyInitials = (companyName) => {
  if (!companyName) return 'CO';
  
  const words = companyName.trim().split(' ');
  if (words.length === 1) {
    return words[0].substring(0, 2).toUpperCase();
  }
  
  return words
    .slice(0, 2)
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase();
};

/**
 * Truncate text to specified length with ellipsis
 * @param {string} text - Text to truncate
 * @param {number} maxLength - Maximum length
 * @returns {string} Truncated text
 */
export const truncateText = (text, maxLength = 100) => {
  if (!text || text.length <= maxLength) return text;
  return text.substring(0, maxLength).trim() + '...';
};

/**
 * Format job count for display
 * @param {number} count - Number of jobs
 * @returns {string} Formatted count
 */
export const formatJobCount = (count) => {
  if (count === 0) return 'No jobs';
  if (count === 1) return '1 job';
  if (count < 1000) return `${count} jobs`;
  if (count < 1000000) return `${(count / 1000).toFixed(1)}k jobs`;
  return `${(count / 1000000).toFixed(1)}M jobs`;
};

/**
 * Create a URL-friendly slug from a string
 * @param {string} text - Text to slugify
 * @returns {string} URL-friendly slug
 */
export const slugify = (text) => {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '');
};

/**
 * Debounce function to limit API calls
 * @param {Function} func - Function to debounce
 * @param {number} delay - Delay in milliseconds
 * @returns {Function} Debounced function
 */
export const debounce = (func, delay) => {
  let timeoutId;
  return (...args) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => func.apply(null, args), delay);
  };
};