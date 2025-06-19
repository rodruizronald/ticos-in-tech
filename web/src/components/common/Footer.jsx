import React from 'react';
import { Heart, MapPin, Mail } from 'lucide-react';

/**
 * Footer component with links and Costa Rica branding
 */
const Footer = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-gray-50 border-t border-gray-200 mt-16">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Brand section */}
          <div className="md:col-span-2">
            <h3 className="text-2xl font-bold text-costa-blue mb-4">
              TicosInTech
            </h3>
            <p className="text-gray-600 mb-4 max-w-md">
              Connecting talented developers with innovative companies across Costa Rica's 
              thriving technology ecosystem.
            </p>
            <div className="flex items-center text-sm text-gray-500">
              <span>Made with</span>
              <Heart className="w-4 h-4 mx-1 text-costa-red" />
              <span>in</span>
              <MapPin className="w-4 h-4 mx-1 text-costa-blue" />
              <span>Costa Rica</span>
            </div>
          </div>

          {/* Job seekers */}
          <div>
            <h4 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-4">
              For Job Seekers
            </h4>
            <ul className="space-y-2">
              <li>
                <button className="text-gray-600 hover:text-costa-blue transition-colors text-sm">
                  Browse Jobs
                </button>
              </li>
              <li>
                <button 
                  className="text-gray-400 text-sm cursor-not-allowed"
                  disabled
                  title="Coming soon"
                >
                  Career Guide
                </button>
              </li>
              <li>
                <button 
                  className="text-gray-400 text-sm cursor-not-allowed"
                  disabled
                  title="Coming soon"
                >
                  Salary Guide
                </button>
              </li>
              <li>
                <button 
                  className="text-gray-400 text-sm cursor-not-allowed"
                  disabled
                  title="Coming soon"
                >
                  Resume Tips
                </button>
              </li>
            </ul>
          </div>

          {/* Employers */}
          <div>
            <h4 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-4">
              For Employers
            </h4>
            <ul className="space-y-2">
              <li>
                <button 
                  className="text-gray-400 text-sm cursor-not-allowed"
                  disabled
                  title="Coming soon"
                >
                  Post a Job
                </button>
              </li>
              <li>
                <button 
                  className="text-gray-400 text-sm cursor-not-allowed"
                  disabled
                  title="Coming soon"
                >
                  Browse Talent
                </button>
              </li>
              <li>
                <button 
                  className="text-gray-400 text-sm cursor-not-allowed"
                  disabled
                  title="Coming soon"
                >
                  Pricing
                </button>
              </li>
              <li>
                <button 
                  className="text-gray-400 text-sm cursor-not-allowed"
                  disabled
                  title="Coming soon"
                >
                  Hiring Guide
                </button>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom section */}
        <div className="border-t border-gray-200 pt-8 mt-8">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <div className="flex items-center space-x-6 mb-4 md:mb-0">
              <p className="text-gray-500 text-sm">
                © {currentYear} TicosInTech. All rights reserved.
              </p>
            </div>
            
            <div className="flex items-center space-x-6">
              <button 
                className="text-gray-400 text-sm cursor-not-allowed"
                disabled
                title="Coming soon"
              >
                Privacy Policy
              </button>
              <button 
                className="text-gray-400 text-sm cursor-not-allowed"
                disabled
                title="Coming soon"
              >
                Terms of Service
              </button>
              <a
                href="mailto:hello@ticosintech.com"
                className="flex items-center text-gray-500 hover:text-costa-blue transition-colors text-sm"
              >
                <Mail className="w-4 h-4 mr-1" />
                Contact
              </a>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;