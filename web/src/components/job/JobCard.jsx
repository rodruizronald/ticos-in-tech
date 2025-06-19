import React from 'react';
import { Heart, MapPin, Globe, Users, Briefcase, Clock, ExternalLink } from 'lucide-react';
import { useApp } from '../../context/AppContext';
import { useRouter } from '../../utils/router';
import { formatRelativeDate, getCompanyInitials, truncateText } from '../../utils/formatters';

/**
 * Job card component for displaying job listings
 */
const JobCard = ({ job, className = '' }) => {
  const { hasSavedJob, toggleSavedJob } = useApp();
  const { navigate } = useRouter();
  
  const isJobSaved = hasSavedJob(job.job_id);

  const handleCardClick = (e) => {
    // Don't navigate if clicking on interactive elements
    if (e.target.closest('button') || e.target.closest('a')) {
      return;
    }
    navigate(`/job/${job.job_id}`);
  };

  const handleSaveJob = (e) => {
    e.stopPropagation();
    toggleSavedJob(job.job_id);
  };

  const handleApplyClick = (e) => {
    e.stopPropagation();
    window.open(job.application_url, '_blank', 'noopener,noreferrer');
  };

  return (
    <div 
      className={`card card-hover p-6 cursor-pointer ${className}`}
      onClick={handleCardClick}
      role="article"
      aria-label={`${job.title} at ${job.company_name}`}
    >
      {/* Header with company logo and save button */}
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center flex-1 min-w-0">
          {/* Company logo */}
          <div className="flex-shrink-0 mr-4">
            {job.company_logo_url ? (
              <img 
                src={job.company_logo_url} 
                alt={`${job.company_name} logo`}
                className="w-12 h-12 rounded-lg object-cover border border-gray-200"
                onError={(e) => {
                  e.target.style.display = 'none';
                  e.target.nextElementSibling.style.display = 'flex';
                }}
              />
            ) : null}
            <div 
              className={`w-12 h-12 rounded-lg bg-costa-blue text-white flex items-center justify-center font-semibold text-sm ${
                job.company_logo_url ? 'hidden' : 'flex'
              }`}
            >
              {getCompanyInitials(job.company_name)}
            </div>
          </div>

          {/* Job title and company */}
          <div className="min-w-0 flex-1">
            <h3 className="text-lg font-semibold text-costa-blue hover:text-blue-700 transition-colors mb-1 truncate">
              {job.title}
            </h3>
            <p className="text-gray-600 font-medium truncate">
              {job.company_name}
            </p>
          </div>
        </div>

        {/* Save button */}
        <button
          onClick={handleSaveJob}
          className={`p-2 rounded-full transition-colors hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-costa-blue focus:ring-offset-2 ${
            isJobSaved ? 'text-costa-red' : 'text-gray-400 hover:text-costa-red'
          }`}
          aria-label={isJobSaved ? 'Remove from saved jobs' : 'Save job'}
        >
          <Heart className={`w-5 h-5 ${isJobSaved ? 'fill-current' : ''}`} />
        </button>
      </div>

      {/* Job metadata */}
      <div className="flex flex-wrap items-center gap-4 text-sm text-gray-500 mb-3">
        <span className="flex items-center">
          <MapPin className="w-4 h-4 mr-1" />
          {job.location}
        </span>
        
        {job.work_mode && (
          <span className="flex items-center">
            <Globe className="w-4 h-4 mr-1" />
            {job.work_mode}
          </span>
        )}
        
        {job.experience_level && (
          <span className="flex items-center">
            <Users className="w-4 h-4 mr-1" />
            {job.experience_level}
          </span>
        )}
      </div>

      {/* Employment type badge */}
      {job.employment_type && (
        <div className="mb-4">
          <span className="badge badge-secondary">
            <Briefcase className="w-3 h-3 mr-1" />
            {job.employment_type}
          </span>
        </div>
      )}

      {/* Job description preview */}
      {job.description && (
        <p className="text-gray-700 text-sm mb-4 truncate-2-lines">
          {truncateText(job.description, 120)}
        </p>
      )}

      {/* Technology tags */}
      {job.technologies && job.technologies.length > 0 && (
        <div className="flex flex-wrap gap-1 mb-4">
          {job.technologies.slice(0, 6).map((tech, index) => (
            <span
              key={index}
              className={`text-xs px-2 py-1 rounded-full border transition-colors ${
                tech.required 
                  ? 'bg-costa-blue text-white border-costa-blue font-medium' 
                  : 'bg-gray-50 text-gray-600 border-gray-200 hover:bg-gray-100'
              }`}
              title={`${tech.name} - ${tech.required ? 'Required' : 'Optional'}`}
            >
              {tech.name}
            </span>
          ))}
          {job.technologies.length > 6 && (
            <span 
              className="text-xs px-2 py-1 rounded-full bg-gray-50 text-gray-500 border border-gray-200"
              title={`${job.technologies.length - 6} more technologies`}
            >
              +{job.technologies.length - 6}
            </span>
          )}
        </div>
      )}

      {/* Footer with posted date and apply button */}
      <div className="flex items-center justify-between pt-2 border-t border-gray-100">
        <span className="text-sm text-gray-500 flex items-center">
          <Clock className="w-4 h-4 mr-1" />
          {formatRelativeDate(job.posted_at)}
        </span>
        
        <button
          onClick={handleApplyClick}
          className="btn btn-primary btn-sm flex items-center space-x-1 shadow-sm hover:shadow-md transition-shadow"
        >
          <span>Apply</span>
          <ExternalLink className="w-3 h-3" />
        </button>
      </div>
    </div>
  );
};

export default JobCard;