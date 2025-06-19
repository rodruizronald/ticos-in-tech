import React from 'react';
import { useRouter } from '../utils/router';
import { Home, Search, ArrowLeft } from 'lucide-react';

/**
 * 404 Not Found page component
 */
const NotFoundPage = () => {
  const { navigate } = useRouter();

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <div className="text-center">
          {/* 404 illustration */}
          <div className="mx-auto w-32 h-32 mb-8">
            <div className="relative">
              <div className="text-8xl font-bold text-gray-200 select-none">404</div>
              <div className="absolute inset-0 flex items-center justify-center">
                <Search className="w-12 h-12 text-gray-400" />
              </div>
            </div>
          </div>

          {/* Title and description */}
          <h1 className="text-3xl font-bold text-gray-900 mb-4">
            Page Not Found
          </h1>
          
          <p className="text-gray-600 mb-8 max-w-md mx-auto">
            The page you're looking for doesn't exist. It might have been moved, 
            deleted, or you entered the wrong URL.
          </p>

          {/* Action buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <button
              onClick={() => navigate('/')}
              className="btn btn-primary flex items-center justify-center space-x-2"
            >
              <Home className="w-4 h-4" />
              <span>Go Home</span>
            </button>
            
            <button
              onClick={() => window.history.back()}
              className="btn btn-secondary flex items-center justify-center space-x-2"
            >
              <ArrowLeft className="w-4 h-4" />
              <span>Go Back</span>
            </button>
          </div>

          {/* Helpful links */}
          <div className="mt-12 pt-8 border-t border-gray-200">
            <p className="text-sm text-gray-500 mb-4">
              Looking for something specific?
            </p>
            
            <div className="flex flex-wrap justify-center gap-4 text-sm">
              <button
                onClick={() => navigate('/')}
                className="text-costa-blue hover:text-blue-700 underline"
              >
                Browse Jobs
              </button>
              
              <span className="text-gray-300">•</span>
              
              <button
                className="text-gray-400 cursor-not-allowed"
                disabled
                title="Coming soon"
              >
                View Companies
              </button>
              
              <span className="text-gray-300">•</span>
              
              <button
                className="text-gray-400 cursor-not-allowed"
                disabled
                title="Coming soon"
              >
                Contact Support
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NotFoundPage;