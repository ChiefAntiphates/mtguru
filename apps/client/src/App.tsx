import { useState, useEffect } from 'react'
import SearchBar from './components/SearchBar'
import Filters from './components/Filters'
import CardGrid from './components/CardGrid'
import mtguruLogo from './assets/mtguru-logo.png'
import { Card, SearchResponse } from './types/card'
import './App.css'

interface FilterOptions {
  set_type: string
  colors: string
  rarity: string
  [key: string]: string
}

function App() {
  const [filters, setFilters] = useState<FilterOptions>({
    set_type: '',
    colors: '',
    rarity: '',
  })
  const [cards, setCards] = useState<Card[]>([])
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
      console.log('Sending request to:', 'https://internal.mtguru.com/api/search')
      const response = await fetch('https://internal.mtguru.com/api/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
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
      
      // Check for different possible response formats
      if (data.data?.Get?.Mtguru) {
        const searchResults = data.data.Get.Mtguru
        console.log('Number of cards found:', searchResults.length)
        
        if (searchResults.length === 0) {
          console.log('No cards found for query:', query)
        }
        
        setCards(searchResults)
      } else if (data.matches?.data?.Get?.Mtguru) {
        // Handle old format
        const searchResults = data.matches.data.Get.Mtguru
        console.log('Number of cards found (old format):', searchResults.length)
        setCards(searchResults)
      } else {
        console.error('Unexpected response format:', data)
        throw new Error('Unexpected response format from server')
      }
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
