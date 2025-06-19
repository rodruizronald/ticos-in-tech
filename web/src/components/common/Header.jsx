import React from 'react';
import { useRouter } from '../../utils/router';

/**
 * Header component with navigation and branding
 */
const Header = () => {
  const { navigate, currentPath } = useRouter();

  const handleLogoClick = () => {
    navigate('/');
  };

  return (
    <header className="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-18">
          {/* Logo */}
          <div className="flex items-center">
            <button
              onClick={handleLogoClick}
              className="text-2xl font-bold text-costa-blue hover:text-blue-700 transition-colors focus:outline-none focus:ring-2 focus:ring-costa-blue focus:ring-offset-2 rounded-lg px-2 py-1"
              aria-label="TicosInTech Home"
            >
              TicosInTech
            </button>
          </div>

          {/* Navigation */}
          <nav className="hidden md:flex items-center space-x-8">
            <button
              onClick={() => navigate('/')}
              className={`text-sm font-medium transition-colors px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-costa-blue focus:ring-offset-2 ${
                currentPath === '/' 
                  ? 'text-costa-blue bg-blue-50' 
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Jobs
            </button>
            
            {/* Future navigation items */}
            <button
              className="text-gray-400 text-sm font-medium cursor-not-allowed"
              disabled
              title="Coming soon"
            >
              Companies
            </button>
            
            <button
              className="text-gray-400 text-sm font-medium cursor-not-allowed"
              disabled
              title="Coming soon"
            >
              Resources
            </button>
          </nav>

          {/* User actions */}
          <div className="flex items-center space-x-4">
            {/* Login/Register buttons - Phase 2 */}
            <button
              className="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-lg text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-costa-blue focus:ring-offset-2"
              disabled
              title="Coming soon"
            >
              Login
            </button>
            
            <button
              className="btn btn-primary btn-sm opacity-50 cursor-not-allowed"
              disabled
              title="Coming soon"
            >
              Register
            </button>

            {/* Mobile menu button - for future mobile navigation */}
            <button
              className="md:hidden p-2 rounded-lg text-gray-600 hover:text-gray-900 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-costa-blue focus:ring-offset-2"
              aria-label="Open menu"
              disabled
              title="Coming soon"
            >
              <svg
                className="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 6h16M4 12h16M4 18h16"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;