import { createContext, useContext, useState, useEffect } from 'react';

// Router Context
const RouterContext = createContext();

/**
 * Simple router implementation for the MVP
 * Can be easily replaced with React Router later
 */
export const Router = ({ children }) => {
  const [currentPath, setCurrentPath] = useState(window.location.pathname);
  const [params, setParams] = useState({});

  useEffect(() => {
    const handlePopState = () => {
      setCurrentPath(window.location.pathname);
      parseParams(window.location.pathname);
    };

    window.addEventListener('popstate', handlePopState);
    parseParams(currentPath);

    return () => window.removeEventListener('popstate', handlePopState);
  }, [currentPath]);

  const parseParams = (path) => {
    // Simple parameter parsing for routes like /job/:id
    const jobMatch = path.match(/^\/job\/(\d+)$/);
    if (jobMatch) {
      setParams({ id: jobMatch[1] });
    } else {
      setParams({});
    }
  };

  const navigate = (path) => {
    if (path !== currentPath) {
      window.history.pushState({}, '', path);
      setCurrentPath(path);
      parseParams(path);
    }
  };

  const value = {
    currentPath,
    params,
    navigate
  };

  return (
    <RouterContext.Provider value={value}>
      {children}
    </RouterContext.Provider>
  );
};

/**
 * Hook to access router functionality
 */
export const useRouter = () => {
  const context = useContext(RouterContext);
  if (!context) {
    throw new Error('useRouter must be used within a Router component');
  }
  return context;
};

/**
 * Route component for conditional rendering
 */
export const Route = ({ path, component: Component, exact = false }) => {
  const { currentPath } = useRouter();
  
  const isMatch = exact 
    ? currentPath === path
    : currentPath.startsWith(path) || 
      (path.includes(':') && matchParameterizedRoute(currentPath, path));

  return isMatch ? <Component /> : null;
};

/**
 * Helper function to match parameterized routes
 */
const matchParameterizedRoute = (currentPath, routePath) => {
  const currentSegments = currentPath.split('/').filter(Boolean);
  const routeSegments = routePath.split('/').filter(Boolean);

  if (currentSegments.length !== routeSegments.length) {
    return false;
  }

  return routeSegments.every((segment, index) => {
    return segment.startsWith(':') || segment === currentSegments[index];
  });
};