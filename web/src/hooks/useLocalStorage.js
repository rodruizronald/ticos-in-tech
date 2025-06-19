import { useState, useEffect } from 'react';

/**
 * Hook for managing localStorage with React state
 * @param {string} key - localStorage key
 * @param {any} initialValue - Initial value if key doesn't exist
 * @returns {[value, setValue, removeValue]} - State value, setter, and remover
 */
export const useLocalStorage = (key, initialValue) => {
  // Get value from localStorage or use initial value
  const [storedValue, setStoredValue] = useState(() => {
    try {
      const item = window.localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      console.warn(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  });

  // Return a wrapped version of useState's setter function that persists the new value to localStorage
  const setValue = (value) => {
    try {
      // Allow value to be a function so we have the same API as useState
      const valueToStore = value instanceof Function ? value(storedValue) : value;
      
      // Save state
      setStoredValue(valueToStore);
      
      // Save to localStorage
      window.localStorage.setItem(key, JSON.stringify(valueToStore));
    } catch (error) {
      console.warn(`Error setting localStorage key "${key}":`, error);
    }
  };

  // Function to remove the item from localStorage
  const removeValue = () => {
    try {
      window.localStorage.removeItem(key);
      setStoredValue(initialValue);
    } catch (error) {
      console.warn(`Error removing localStorage key "${key}":`, error);
    }
  };

  // Listen for changes to localStorage from other tabs/windows
  useEffect(() => {
    const handleStorageChange = (e) => {
      if (e.key === key && e.newValue !== null) {
        try {
          setStoredValue(JSON.parse(e.newValue));
        } catch (error) {
          console.warn(`Error parsing localStorage value for key "${key}":`, error);
        }
      }
    };

    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, [key]);

  return [storedValue, setValue, removeValue];
};

/**
 * Hook for managing a Set in localStorage (useful for saved jobs)
 * @param {string} key - localStorage key
 * @param {Set} initialValue - Initial Set value
 * @returns {[Set, addItem, removeItem, hasItem, clearAll]} - Set operations
 */
export const useLocalStorageSet = (key, initialValue = new Set()) => {
  const [storedSet, setStoredSet] = useLocalStorage(key, Array.from(initialValue));
  
  const currentSet = new Set(storedSet);

  const addItem = (item) => {
    const newSet = new Set(currentSet);
    newSet.add(item);
    setStoredSet(Array.from(newSet));
  };

  const removeItem = (item) => {
    const newSet = new Set(currentSet);
    newSet.delete(item);
    setStoredSet(Array.from(newSet));
  };

  const hasItem = (item) => {
    return currentSet.has(item);
  };

  const clearAll = () => {
    setStoredSet([]);
  };

  const toggleItem = (item) => {
    if (hasItem(item)) {
      removeItem(item);
    } else {
      addItem(item);
    }
  };

  return [currentSet, addItem, removeItem, hasItem, clearAll, toggleItem];
};