import React, { useState, useEffect } from 'react';
import { 
  ArrowLeft, 
  Heart, 
  Share2, 
  ExternalLink, 
  MapPin, 
  Globe, 
  Users, 
  Briefcase, 
  Building, 
  Clock,
  Copy,
  Check
} from 'lucide-react';
import { useApp } from '../../context/AppContext';
import { useRouter } from '../../utils/router';
import { formatRelativeDate, getCompanyInitials } from '../../utils/formatters';
import { FullPageLoader } from '../common/LoadingSpinner';

/**
 * Job detail page component
 */
const JobDetail = ({ jobId }) => {
  const [job, setJob] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [copied, setCopied] = useState(false);
  
  const { hasSavedJob, toggleSavedJob, getJobById } = useApp();
  const { navigate } = useRouter();

  useEffect(() => {
    const fetchJob = async () => {
      try {
        setLoading(true);
        setError(null);
        const jobData = await getJobById(jobId);
        setJob(jobData);
      } catch (err) {
        setError(err.message || 'Failed to load job details');
      } finally {
        setLoading(false);
      }
    };

    if (jobId) {
      fetchJob();
    }
  }, [jobId, getJobById]);

  const handleShare = async () => {
    const url = window.location.href;
    
    if (navigator.share) {
      try {
        await navigator.share({
          title: `${job.title} at ${job.company_name}`,
          text: `Check out this job opportunity: ${job.title} at ${job.company_name}`,
          url: url
        });
      } catch (err) {
        console.log('Share cancelled');
      }
    } else {
      // Fallback to copying URL
      try {
        await navigator.clipboard.writeText(url);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      } catch (err) {
        console.error('Failed to copy URL');
      }
    }
  };

  const handleApply = () => {
    window.open(job.application_url, '_blank', 'noopener,noreferrer');
  };

  if (loading) {
    return <FullPageLoader text="Loading job details..." />;
  }

  if (error || !job) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center">
          <div className="w-24 h-24 mx-auto mb-6 bg-red-100 rounded-full flex items-center justify-center">
            <svg className="w-12 h-12 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Job Not Found</h1>
          <p className="text-gray-600 mb-6">
            {error || 'The job you\'re looking for could not be found.'}
          </p>
          
          <button 
            onClick={() => navigate('/')}
            className="btn btn-primary"
          >
            Back to Jobs
          </button>
        </div>
      </div>
    );
  }

  const isJobSaved = hasSavedJob(job.job_id);

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Back button */}
      <button 
        onClick={() => navigate('/')}
        className="flex items-center text-costa-blue hover:text-blue-700 mb-8 group"
      >
        <ArrowLeft className="w-4 h-4 mr-2 group-hover:-translate-x-1 transition-transform" />
        Back to Jobs
      </button>

      {/* Job header */}
      <div className="bg-white rounded-xl border border-gray-200 shadow-subtle p-8 mb-8">
        <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between mb-8">
          <div className="flex items-start space-x-6 flex-1">
            {/* Company logo */}
            <div className="flex-shrink-0">
              {job.company_logo_url ? (
                <img 
                  src={job.company_logo_url} 
                  alt={`${job.company_name} logo`}
                  className="w-20 h-20 rounded-xl object-cover border border-gray-200"
                  onError={(e) => {
                    e.target.style.display = 'none';
                    e.target.nextElementSibling.style.display = 'flex';
                  }}
                />
              ) : null}
              <div 
                className={`w-20 h-20 rounded-xl bg-costa-blue text-white flex items-center justify-center font-bold text-2xl ${
                  job.company_logo_url ? 'hidden' : 'flex'
                }`}
              >
                {getCompanyInitials(job.company_name)}
              </div>
            </div>

            {/* Job info */}
            <div className="flex-1 min-w-0">
              <h1 className="text-3xl font-bold text-gray-900 mb-3">
                {job.title}
              </h1>
              
              <div className="flex items-center space-x-2 mb-4">
                <h2 className="text-xl text-gray-700 font-medium">
                  {job.company_name}
                </h2>
                <button className="text-costa-blue hover:text-blue-700 text-sm">
                  View company →
                </button>
              </div>

              {/* Job metadata */}
              <div className="flex flex-wrap items-center gap-6 text-gray-600">
                <div className="flex items-center">
                  <MapPin className="w-5 h-5 mr-2" />
                  <span>{job.location}</span>
                </div>
                
                {job.work_mode && (
                  <div className="flex items-center">
                    <Globe className="w-5 h-5 mr-2" />
                    <span>{job.work_mode}</span>
                  </div>
                )}
                
                {job.employment_type && (
                  <div className="flex items-center">
                    <Briefcase className="w-5 h-5 mr-2" />
                    <span>{job.employment_type}</span>
                  </div>
                )}
                
                {job.experience_level && (
                  <div className="flex items-center">
                    <Users className="w-5 h-5 mr-2" />
                    <span>{job.experience_level}</span>
                  </div>
                )}
              </div>

              {/* Posted date */}
              <div className="flex items-center text-gray-500 mt-3">
                <Clock className="w-4 h-4 mr-2" />
                <span className="text-sm">
                  {formatRelativeDate(job.posted_at)}
                </span>
              </div>
            </div>
          </div>

          {/* Action buttons */}
          <div className="flex items-center space-x-3 mt-6 lg:mt-0">
            <button
              onClick={() => toggleSavedJob(job.job_id)}
              className={`p-3 rounded-xl border transition-colors ${
                isJobSaved 
                  ? 'text-costa-red border-red-300 bg-red-50' 
                  : 'text-gray-400 border-gray-300 hover:text-costa-red hover:border-red-300'
              }`}
              aria-label={isJobSaved ? 'Remove from saved jobs' : 'Save job'}
            >
              <Heart className={`w-6 h-6 ${isJobSaved ? 'fill-current' : ''}`} />
            </button>
            
            <button 
              onClick={handleShare}
              className="p-3 rounded-xl border border-gray-300 text-gray-400 hover:text-gray-600 hover:border-gray-400 transition-colors"
              aria-label="Share job"
            >
              {copied ? <Check className="w-6 h-6 text-green-500" /> : <Share2 className="w-6 h-6" />}
            </button>
          </div>
        </div>

        {/* Apply button */}
        <div className="flex flex-col sm:flex-row gap-4">
          <button
            onClick={handleApply}
            className="btn btn-danger btn-lg flex items-center justify-center space-x-2 flex-1 sm:flex-none"
          >
            <span>Apply for this position</span>
            <ExternalLink className="w-5 h-5" />
          </button>
          
          <div className="text-sm text-gray-500 sm:flex sm:items-center">
            <span>You'll be redirected to the company's application page</span>
          </div>
        </div>
      </div>

      {/* Content grid */}
      <div className="grid lg:grid-cols-3 gap-8">
        {/* Main content */}
        <div className="lg:col-span-2 space-y-8">
          {/* Job description */}
          <section className="bg-white rounded-xl border border-gray-200 shadow-subtle p-8">
            <h3 className="text-xl font-semibold text-gray-900 mb-6">
              Job Description
            </h3>
            
            <div className="prose max-w-none text-gray-700 leading-relaxed">
              {job.description.split('\n').map((paragraph, index) => (
                paragraph.trim() && (
                  <p key={index} className="mb-4 last:mb-0">
                    {paragraph}
                  </p>
                )
              ))}
            </div>
          </section>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Technologies */}
          <section className="bg-white rounded-xl border border-gray-200 shadow-subtle p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Technologies & Skills
            </h3>
            
            {job.technologies && job.technologies.length > 0 ? (
              <div className="space-y-3">
                {job.technologies.map((tech, index) => (
                  <div
                    key={index}
                    className={`flex items-center justify-between p-3 rounded-lg border ${
                      tech.required 
                        ? 'bg-costa-blue/5 border-costa-blue/20' 
                        : 'bg-gray-50 border-gray-200'
                    }`}
                  >
                    <div>
                      <span className={`font-medium ${
                        tech.required ? 'text-costa-blue' : 'text-gray-700'
                      }`}>
                        {tech.name}
                      </span>
                      {tech.category && (
                        <div className="text-xs text-gray-500 mt-1">
                          {tech.category}
                        </div>
                      )}
                    </div>
                    
                    <span className={`text-xs px-2 py-1 rounded-full font-medium ${
                      tech.required 
                        ? 'bg-costa-blue text-white' 
                        : 'bg-gray-200 text-gray-600'
                    }`}>
                      {tech.required ? 'Required' : 'Nice to have'}
                    </span>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-gray-500 text-sm">
                No specific technologies listed
              </p>
            )}
          </section>

          {/* Company info */}
          <section className="bg-white rounded-xl border border-gray-200 shadow-subtle p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              About {job.company_name}
            </h3>
            
            <div className="flex items-center mb-4">
              <Building className="w-5 h-5 text-gray-400 mr-3" />
              <span className="font-medium text-gray-900">{job.company_name}</span>
            </div>
            
            <p className="text-gray-600 text-sm mb-4">
              Join one of Costa Rica's innovative technology companies and be part of the growing tech ecosystem.
            </p>
            
            <button 
              className="text-costa-blue hover:text-blue-700 text-sm font-medium"
              disabled
              title="Coming soon"
            >
              View all jobs from {job.company_name} →
            </button>
          </section>
        </div>
      </div>
    </div>
  );
};

export default JobDetail;