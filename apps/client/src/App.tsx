import { useState, useEffect } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

const SERVER_URL =  "http://localhost:8888"

function App() {
  const [count, setCount] = useState(0)
  const [time, setTime] = useState<string>('')

  useEffect(() => {
    fetch(`${SERVER_URL}/clicks`, {
      method: "POST",
      body: JSON.stringify({
        clicks: count
      }),
      headers: {
        "Content-type": "application/json; charset=UTF-8"
      }
    })
  }, [count])

  return (
    <>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => fetch(`${SERVER_URL}/time`).then(response => response.text()).then(data => setTime(data))}>
          Get time from Go
        </button>
        {time != '' && <p>Go thinks the time is {time}</p>}
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
