import React, { useState } from 'react'
import SearchBar from './components/SearchBar'
import Filters from './components/Filters'
import mtguruLogo from './assets/mtguru-logo.png'
import './App.css'

interface FilterOptions {
  type: string
  color: string
  rarity: string
  [key: string]: string
}

function App() {
  const [filters, setFilters] = useState<FilterOptions>({
    type: '',
    color: '',
    rarity: '',
  })

  const handleSearch = async (query: string) => {
    try {
      const response = await fetch('http://localhost:8080/api/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          "query": query,
          "filters": filters
        }),
      })

      if (!response.ok) {
        throw new Error('Search failed')
      }

      const data = await response.json()
      // TODO: Handle the search results
      console.log('Search results:', data)
    } catch (error) {
      console.error('Error during search:', error)
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
        <SearchBar onSearch={handleSearch} />
        <Filters onFilterChange={setFilters} />
        {/* TODO: Add CardResults component here */}
      </main>
    </div>
  )
}

export default App
