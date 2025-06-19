import React from 'react';
import Header from '../common/Header';
import Footer from '../common/Footer';

/**
 * Main layout component that wraps all pages
 */
const Layout = ({ children, showFooter = true }) => {
  return (
    <div className="min-h-screen flex flex-col bg-white">
      <Header />
      
      <main className="flex-1">
        {children}
      </main>
      
      {showFooter && <Footer />}
    </div>
  );
};

export default Layout;