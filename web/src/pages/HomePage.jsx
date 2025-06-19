import React, { useEffect } from 'react';
import HeroSection from '../components/search/HeroSection';
import JobGrid from '../components/job/JobGrid';
import { useApp } from '../context/AppContext';

/**
 * Home page component with hero section and job listings
 */
const HomePage = () => {
  const { searchJobs, filters, jobs } = useApp();

  // Initial load of jobs when component mounts
  useEffect(() => {
    // Only search if we don't have jobs loaded yet
    if (jobs.length === 0) {
      searchJobs('developer', filters, 0);
    }
  }, []); // Empty dependency array for initial load only

  return (
    <div className="bg-white">
      {/* Hero section with search */}
      <HeroSection />
      
      {/* Job listings */}
      <JobGrid />
    </div>
  );
};

export default HomePage;