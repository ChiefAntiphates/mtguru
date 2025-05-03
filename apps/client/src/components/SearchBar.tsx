import React, { useState } from 'react';
import '../styles/SearchBar.css';

interface SearchBarProps {
  onSearch: (query: string) => void;
  onFilterToggle: () => void;
  isFiltersExpanded: boolean;
}

const SearchBar: React.FC<SearchBarProps> = ({ onSearch, onFilterToggle, isFiltersExpanded }) => {
  const [searchQuery, setSearchQuery] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(searchQuery);
  };

  return (
    <form className="search-container" onSubmit={handleSubmit}>
      <div className="search-box">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search for cards (e.g., 'make dragons stronger')"
          className="search-input"
        />
        <button 
          type="button" 
          className={`filter-button ${isFiltersExpanded ? 'active' : ''}`}
          onClick={onFilterToggle}
          title="Toggle filters"
        >
          <svg 
            xmlns="http://www.w3.org/2000/svg" 
            width="24" 
            height="24" 
            viewBox="0 0 24 24" 
            fill="none" 
            stroke="currentColor" 
            strokeWidth="2" 
            strokeLinecap="round" 
            strokeLinejoin="round"
          >
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
          </svg>
        </button>
        <button type="submit" className="search-button">
          Search
        </button>
      </div>
    </form>
  );
};

export default SearchBar; 