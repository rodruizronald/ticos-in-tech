import React from 'react';
import { Router, Route, useRouter } from './utils/router';
import { AppProvider } from './context/AppContext';
import Layout from './components/layout/Layout';
import HomePage from './pages/HomePage';
import JobDetailPage from './pages/JobDetailPage';
import NotFoundPage from './pages/NotFoundPage';
import './styles/index.css';

/**
 * Main App component with routing and global state
 */
function App() {
  return (
    <AppProvider>
      <Router>
        <Layout>
          <AppRoutes />
        </Layout>
      </Router>
    </AppProvider>
  );
}

/**
 * App routing component
 */
const AppRoutes = () => {
  const { currentPath } = useRouter();
  
  // Home page
  if (currentPath === '/') {
    return <HomePage />;
  }
  
  // Job detail page
  if (currentPath.match(/^\/job\/\d+$/)) {
    return <JobDetailPage />;
  }
  
  // 404 page
  return <NotFoundPage />;
};

export default App;