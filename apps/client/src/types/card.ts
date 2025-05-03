export interface CardImageUris {
  large: string;
}

export interface CardAdditional {
  distance: number;
}

export interface Card {
  _additional: CardAdditional;
  image_uris?: CardImageUris;
  name: string;
  oracle_text: string;
  set_name: string;
  scryfall_uri: string;
}

export interface SearchResponse {
  data?: {
    Get: {
      Mtguru: Card[];
    };
  };
  matches?: {
    data: {
      Get: {
        Mtguru: Card[];
      };
    };
  };
} 