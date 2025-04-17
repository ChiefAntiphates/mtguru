import React, { useState } from 'react';
import '../styles/Filters.css';

interface FilterOptions {
  type: string;
  color: string;
  rarity: string;
  [key: string]: string;
}

interface FiltersProps {
  onFilterChange: (filters: FilterOptions) => void;
}

const Filters: React.FC<FiltersProps> = ({ onFilterChange }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [filters, setFilters] = useState<FilterOptions>({
    type: '',
    color: '',
    rarity: '',
  });

  const handleFilterChange = (key: string, value: string) => {
    const newFilters = { ...filters, [key]: value };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  return (
    <div className="filters-container">
      <button 
        className="filters-toggle"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        {isExpanded ? 'Hide Filters' : 'Show Filters'}
      </button>
      
      {isExpanded && (
        <div className="filters-content">
          <div className="filter-group">
            <label>Type</label>
            <select
              value={filters.type}
              onChange={(e) => handleFilterChange('type', e.target.value)}
            >
              <option value="">Any</option>
              <option value="creature">Creature</option>
              <option value="instant">Instant</option>
              <option value="sorcery">Sorcery</option>
              <option value="enchantment">Enchantment</option>
              <option value="artifact">Artifact</option>
              <option value="planeswalker">Planeswalker</option>
              <option value="land">Land</option>
            </select>
          </div>

          <div className="filter-group">
            <label>Color</label>
            <select
              value={filters.color}
              onChange={(e) => handleFilterChange('color', e.target.value)}
            >
              <option value="">Any</option>
              <option value="white">White</option>
              <option value="blue">Blue</option>
              <option value="black">Black</option>
              <option value="red">Red</option>
              <option value="green">Green</option>
              <option value="colorless">Colorless</option>
            </select>
          </div>

          <div className="filter-group">
            <label>Rarity</label>
            <select
              value={filters.rarity}
              onChange={(e) => handleFilterChange('rarity', e.target.value)}
            >
              <option value="">Any</option>
              <option value="common">Common</option>
              <option value="uncommon">Uncommon</option>
              <option value="rare">Rare</option>
              <option value="mythic">Mythic Rare</option>
            </select>
          </div>
        </div>
      )}
    </div>
  );
};

export default Filters; 