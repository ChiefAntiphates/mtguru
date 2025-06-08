import React from 'react';
import { Match } from '../types/responses';
import './CardGrid.css';

interface CardGridProps {
  cards: Match[];
}

const CardGrid: React.FC<CardGridProps> = ({ cards }) => {
  return (
    <div className="card-grid">
      {cards.map((card, index) => (
        <a 
          key={index} 
          href={card.metadata.scryfall_uri} 
          target="_blank" 
          rel="noopener noreferrer"
          className="card-link"
        >
          <div className="card-item">
            
            <img 
              src={card.metadata.image_url} 
              alt={card.metadata.name}
              className="card-image"
              loading="lazy"
            />
           
            <div className="match-bar">
              <div 
                className="match-progress" 
                style={{ width: `${card.score * 100}%` }}
              >
                <span className="match-value">
                  {Math.round(card.score * 100)}%
                </span>
              </div>
            </div>
            <div className="card-info">
              <h3 className="card-name">{card.metadata.name}</h3>
              <div className="card-details">
                <p className="card-set">{card.metadata.set_name}</p>
              </div>
            </div>
          </div>
        </a>
      ))}
    </div>
  );
};

export default CardGrid; 