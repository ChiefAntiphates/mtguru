import { useState, useEffect } from 'react'
import SearchBar from './components/SearchBar'
import Filters from './components/Filters'
import CardGrid from './components/CardGrid'
import mtguruLogo from './assets/mtguru-logo.png'
import { Match, SearchResponse } from './types/responses'
import './App.css'

interface FilterOptions {
  set_type: string
  colors: string
  rarity: string
  [key: string]: string
}

function App() {
  const API_URL = import.meta.env.VITE_API_URL
  const [filters, setFilters] = useState<FilterOptions>({
    set_type: '',
    colors: '',
    rarity: '',
  })
  const [cards, setCards] = useState<Match[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [searchAttempted, setSearchAttempted] = useState(false)
  const [renderError, setRenderError] = useState<string | null>(null)
  const [isFiltersExpanded, setIsFiltersExpanded] = useState(false)

  // Error boundary for rendering
  useEffect(() => {
    const handleError = (event: ErrorEvent) => {
      console.error('Unhandled error:', event.error)
      setRenderError('An error occurred while rendering the page. Please try refreshing.')
    }

    window.addEventListener('error', handleError)
    return () => window.removeEventListener('error', handleError)
  }, [])

  const handleSearch = async (query: string) => {
    if (!query.trim()) {
      setError('Please enter a search query')
      return
    }

    console.log('Starting search for:', query)
    setIsLoading(true)
    setError(null)
    setCards([])
    setSearchAttempted(true)
    setRenderError(null)

    try {
      console.log('Sending request to:', API_URL + '/api/search')
      const response = await fetch(API_URL + '/api/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          "query": query,
          "filters": filters
        }),
      })

      console.log('Response status:', response.status)
      
      if (!response.ok) {
        const errorText = await response.text()
        console.error('Server error response:', errorText)
        throw new Error(`Search failed with status: ${response.status}. ${errorText}`)
      }

      const rawData = await response.text()
      console.log('Raw response:', rawData)

      let data: SearchResponse
      try {
        data = JSON.parse(rawData)
      } catch (e) {
        console.error('Failed to parse JSON:', e)
        throw new Error('Invalid JSON response from server')
      }

      console.log('Parsed response data:', data)

      console.log('Number of cards found:', data.count)
      
      if (data.count === 0) {
        console.log('No cards found for query:', query)
      }
      
      setCards(data.matches)

    } catch (error) {
      console.error('Error during search:', error)
      setError(error instanceof Error ? error.message : 'An error occurred during search')
      setCards([])
    } finally {
      setIsLoading(false)
    }
  }

  const renderContent = () => {
    try {
      if (renderError) {
        return <div className="error-message">{renderError}</div>
      }

      if (error) {
        return <div className="error-message">{error}</div>
      }

      if (isLoading) {
        return <div className="loading">Searching for cards...</div>
      }

      if (!searchAttempted) {
        return <div className="welcome-message">Enter a search query to find Magic: The Gathering cards</div>
      }

      if (cards.length === 0) {
        return <div className="no-results">No cards found. Try a different search.</div>
      }

      return <CardGrid cards={cards} />
    } catch (error) {
      console.error('Error in renderContent:', error)
      return <div className="error-message">An error occurred while rendering the content</div>
    }
  }

  return (
    <div className="app">
      <header className="app-header">
        <img src={mtguruLogo} alt="MTGuru Logo" className="app-logo" />
        <h1>MTGuru</h1>
        <p>Search for Magic: The Gathering cards using natural language</p>
      </header>
      
      <main className="app-main">
        <SearchBar 
          onSearch={handleSearch} 
          onFilterToggle={() => setIsFiltersExpanded(!isFiltersExpanded)}
          isFiltersExpanded={isFiltersExpanded}
        />
        {isFiltersExpanded && <Filters onFilterChange={(filters) => setFilters(filters)} />}
        {renderContent()}
      </main>
    </div>
  )
}

export default App
