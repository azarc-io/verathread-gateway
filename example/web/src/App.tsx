import './App.scss'
import {useState} from "react";

function App() {
  const [curCount, setCurCount] = useState(0)

  return (
      <>
        <button onClick={() => setCurCount((count) => count + 1)}>
          remote count is {curCount}
        </button>
      </>
  )
}

export default App
