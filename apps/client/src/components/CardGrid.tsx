import React from 'react';
import { Card } from '../types/card';
import './CardGrid.css';

interface CardGridProps {
  cards: Card[];
}

const CardGrid: React.FC<CardGridProps> = ({ cards }) => {
  return (
    <div className="card-grid">
      {cards.map((card, index) => (
        <a 
          key={index} 
          href={card.scryfall_uri} 
          target="_blank" 
          rel="noopener noreferrer"
          className="card-link"
        >
          <div className="card-item">
            {card.image_uris?.large ? (
              <img 
                src={card.image_uris.large} 
                alt={card.name}
                className="card-image"
                loading="lazy"
              />
            ) : (
              <div className="card-image-placeholder">
                <span className="placeholder-text">No Image Available</span>
              </div>
            )}
            <div className="match-bar">
              <div 
                className="match-progress" 
                style={{ width: `${(1 - card._additional.distance) * 100}%` }}
              >
                <span className="match-value">
                  {Math.round((1 - card._additional.distance) * 100)}%
                </span>
              </div>
            </div>
            <div className="card-info">
              <h3 className="card-name">{card.name}</h3>
              <div className="card-details">
                <p className="card-set">{card.set_name}</p>
              </div>
            </div>
          </div>
        </a>
      ))}
    </div>
  );
};

export default CardGrid; 