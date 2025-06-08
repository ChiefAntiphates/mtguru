export interface Match {
  id: string;
  score: number;
  metadata: {
    colors: [string];
    image_url: string;
    name: string;
    rarity: string;
    release_date: string;
    set_name: string;
    scryfall_uri: string;
  }
}

export interface SearchResponse {
  count: number;
  matches: [Match]
} 