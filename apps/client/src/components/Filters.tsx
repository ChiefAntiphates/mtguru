import React, { useState } from 'react';
import '../styles/Filters.css';

interface FilterOptions {
  set_type: string;
  colors: string;
  rarity: string;
  [key: string]: string;
}

interface FiltersProps {
  onFilterChange: (filters: FilterOptions) => void;
}

const Filters: React.FC<FiltersProps> = ({ onFilterChange }) => {
  const [filters, setFilters] = useState<FilterOptions>({
    set_type: '',
    colors: '',
    rarity: '',
  });

  const handleFilterChange = (key: string, value: string) => {
    const newFilters = { ...filters, [key]: value };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  return (
    <div className="filters-container">
      <div className="filters-content">
        <div className="filter-group">
          <label>Type</label>
          <select
            value={filters.set_type}
            onChange={(e) => handleFilterChange('set_type', e.target.value)}
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
            value={filters.colors}
            onChange={(e) => handleFilterChange('colors', e.target.value)}
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
    </div>
  );
};

export default Filters; 