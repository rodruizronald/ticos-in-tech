import React from 'react';
import SearchBar from './SearchBar';
import { useApp } from '../../context/AppContext';
import { formatJobCount } from '../../utils/formatters';

/**
 * Hero section component for the homepage
 */
const HeroSection = () => {
  const { pagination, jobs } = useApp();

  return (
    <section className="relative bg-hero-gradient text-white overflow-hidden">
      {/* Background overlay */}
      <div className="absolute inset-0 bg-black opacity-15"></div>
      
      {/* Decorative background elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-10 -right-10 w-80 h-80 bg-white opacity-5 rounded-full"></div>
        <div className="absolute -bottom-16 -left-16 w-96 h-96 bg-white opacity-5 rounded-full"></div>
      </div>

      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 sm:py-24 lg:py-32">
        <div className="text-center max-w-4xl mx-auto">
          {/* Main headline */}
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold mb-6 leading-tight">
            Find Your Next{' '}
            <span className="relative">
              Tech Job
              <div className="absolute -bottom-2 left-0 right-0 h-1 bg-costa-red opacity-75 rounded-full"></div>
            </span>
            {' '}in Costa Rica
          </h1>

          {/* Subtitle */}
          <p className="text-lg sm:text-xl lg:text-2xl mb-8 text-blue-100 leading-relaxed max-w-3xl mx-auto">
            Discover opportunities from leading companies across Costa Rica's 
            growing tech ecosystem
          </p>

          {/* Stats */}
          {pagination.total > 0 && (
            <div className="flex flex-wrap justify-center items-center gap-6 mb-12 text-blue-100">
              <div className="flex items-center space-x-2">
                <span className="text-2xl font-bold text-white">
                  {formatJobCount(pagination.total)}
                </span>
                <span>available</span>
              </div>
              <div className="hidden sm:block w-1 h-1 bg-blue-300 rounded-full"></div>
              <div className="flex items-center space-x-2">
                <span className="text-2xl font-bold text-white">50+</span>
                <span>companies</span>
              </div>
              <div className="hidden sm:block w-1 h-1 bg-blue-300 rounded-full"></div>
              <div className="flex items-center space-x-2">
                <span className="text-2xl font-bold text-white">100%</span>
                <span>remote friendly</span>
              </div>
            </div>
          )}

          {/* Search bar */}
          <div className="max-w-3xl mx-auto mb-8">
            <SearchBar 
              size="lg"
              className="animate-fade-in"
            />
          </div>

          {/* Browse all jobs CTA */}
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <button
              onClick={() => {
                document.getElementById('job-listings')?.scrollIntoView({ 
                  behavior: 'smooth' 
                });
              }}
              className="btn btn-secondary btn-lg bg-white text-costa-blue hover:bg-gray-50 border-0 shadow-medium"
            >
              Browse All Jobs
            </button>
            
            <div className="text-blue-100 text-sm">
              or scroll down to see the latest opportunities
            </div>
          </div>
        </div>
      </div>

      {/* Bottom wave decoration */}
      <div className="absolute bottom-0 left-0 right-0">
        <svg
          className="w-full h-12 text-white"
          viewBox="0 0 1200 120"
          preserveAspectRatio="none"
          fill="currentColor"
        >
          <path d="M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z"
            opacity=".25"
          />
          <path d="M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z"
            opacity=".5"
          />
          <path d="M0,0V5.63C149.93,59,314.09,71.32,475.83,42.57c43-7.64,84.23-20.12,127.61-26.46,59-8.63,112.48,12.24,165.56,35.4C827.93,77.22,886,95.24,951.2,90c86.53-7,172.46-45.71,248.8-84.81V0Z" />
        </svg>
      </div>
    </section>
  );
};

export default HeroSection;