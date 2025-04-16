import React, { useState, useRef, useEffect } from 'react';
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
  const [activeFilter, setActiveFilter] = useState<string | null>(null);
  const [filters, setFilters] = useState<FilterOptions>({
    type: '',
    color: '',
    rarity: '',
  });
  const filterRef = useRef<HTMLDivElement>(null);
  const timeoutRef = useRef<number | null>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (filterRef.current && !filterRef.current.contains(event.target as Node)) {
        setIsExpanded(false);
        setActiveFilter(null);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      if (timeoutRef.current !== null) {
        window.clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const handleFilterChange = (key: string, value: string) => {
    const newFilters = { ...filters, [key]: value };
    setFilters(newFilters);
    onFilterChange(newFilters);
    setActiveFilter(null);
  };

  const handleMouseEnter = (key: string) => {
    if (timeoutRef.current !== null) {
      window.clearTimeout(timeoutRef.current);
    }
    setActiveFilter(key);
  };

  const handleMouseLeave = (key: string) => {
    timeoutRef.current = window.setTimeout(() => {
      if (activeFilter === key) {
        setActiveFilter(null);
      }
    }, 100);
  };

  const filterOptions = {
    type: [
      { value: '', label: 'Any' },
      { value: 'creature', label: 'Creature' },
      { value: 'instant', label: 'Instant' },
      { value: 'sorcery', label: 'Sorcery' },
      { value: 'enchantment', label: 'Enchantment' },
      { value: 'artifact', label: 'Artifact' },
      { value: 'planeswalker', label: 'Planeswalker' },
      { value: 'land', label: 'Land' }
    ],
    color: [
      { value: '', label: 'Any' },
      { value: 'white', label: 'White' },
      { value: 'blue', label: 'Blue' },
      { value: 'black', label: 'Black' },
      { value: 'red', label: 'Red' },
      { value: 'green', label: 'Green' },
      { value: 'colorless', label: 'Colorless' }
    ],
    rarity: [
      { value: '', label: 'Any' },
      { value: 'common', label: 'Common' },
      { value: 'uncommon', label: 'Uncommon' },
      { value: 'rare', label: 'Rare' },
      { value: 'mythic', label: 'Mythic Rare' }
    ]
  };

  return (
    <div className="filters-container" ref={filterRef}>
      <button 
        className={`filters-toggle ${isExpanded ? 'active' : ''}`}
        onClick={() => setIsExpanded(!isExpanded)}
        aria-expanded={isExpanded}
      >
        {Object.values(filters).some(value => value !== '') ? (
          <span className="filters-active-indicator">●</span>
        ) : null}
        {isExpanded ? 'Hide Filters' : 'Show Filters'}
      </button>
      
      {isExpanded && (
        <div className="filters-bubble">
          <div className="filters-content">
            {Object.entries(filterOptions).map(([key, options]) => (
              <div 
                className="filter-group" 
                key={key}
                onMouseEnter={() => handleMouseEnter(key)}
                onMouseLeave={() => handleMouseLeave(key)}
              >
                <button
                  className={`filter-button ${activeFilter === key ? 'active' : ''} ${filters[key] ? 'has-value' : ''}`}
                  onClick={() => setActiveFilter(activeFilter === key ? null : key)}
                >
                  {key.charAt(0).toUpperCase() + key.slice(1)}: {filters[key] ? options.find(opt => opt.value === filters[key])?.label : 'Any'}
                </button>
                {activeFilter === key && (
                  <div 
                    className="filter-options-bubble"
                    onMouseEnter={() => handleMouseEnter(key)}
                    onMouseLeave={() => handleMouseLeave(key)}
                  >
                    <div className="filter-options-content">
                      {options.map((option) => (
                        <button
                          key={option.value}
                          className={`filter-option ${filters[key] === option.value ? 'selected' : ''}`}
                          onClick={() => handleFilterChange(key, option.value)}
                        >
                          {option.label}
                        </button>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default Filters; 