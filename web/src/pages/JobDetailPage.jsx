import React from 'react';
import { useRouter } from '../utils/router';
import JobDetail from '../components/job/JobDetail';

/**
 * Job detail page component
 */
const JobDetailPage = () => {
  const { params } = useRouter();
  const jobId = params.id;

  return (
    <div className="bg-gray-50 min-h-screen">
      <JobDetail jobId={jobId} />
    </div>
  );
};

export default JobDetailPage;