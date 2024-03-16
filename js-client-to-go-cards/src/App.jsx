import axios from "axios"
import { useState } from "react"

const App = () => {
    const [response, setResponse] = useState()

    const handleCreateDeckClick = () => {
        const promise = axios.get("http://localhost:8000/create")
        promise.then((res) => {
            setResponse(JSON.stringify(res.data))
        })
    }

    return (
        <div>
            <button onClick={handleCreateDeckClick}>create deck</button>
            <div>server response: {response}</div>
        </div>
    )
}

export default App
