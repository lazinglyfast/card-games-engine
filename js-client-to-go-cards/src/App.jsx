import axios from "axios"
import { useState } from "react"
import "./App.css"

const App = () => {
    const [createRes, setCreateRes] = useState("")
    const [openRes, setOpenRes] = useState("")
    const [drawRes, setDrawRes] = useState("")
    const [cards, setCards] = useState("")
    const [shuffled, setShuffled] = useState(false)
    const [guid, setGuid] = useState("")
    const [count, setCount] = useState(0)

    const baseUrl = "http://localhost:8000"

    const handleCreate = () => {
        var url = addQueryParams(`${baseUrl}/create`)

        const promise = axios.post(url)
        promise
            .then((res) => {
                setCreateRes(JSON.stringify(res.data))
            })
            .catch((err) => {
                setCreateRes(JSON.stringify(err))
            })
    }

    const addQueryParams = (url) => {
        var params = []
        if (cards.length > 0) {
            params.push("cards=" + cards)
        }

        if (shuffled === true) {
            params.push("shuffled=true")
        }

        if (count > 0) {
            params.push("count=" + count)
        }

        if (params.length > 0) {
            url += "?" + params.join("&")
        }

        return url
    }

    const handleOpen = () => {
        var url = `${baseUrl}/open/${guid}`

        const promise = axios.get(url)
        promise
            .then((res) => {
                setOpenRes(JSON.stringify(res.data))
            })
            .catch((err) => {
                setOpenRes(JSON.stringify(err))
            })
    }

    const handleDraw = () => {
        var url = addQueryParams(`${baseUrl}/draw/${guid}`)

        const promise = axios.patch(url)
        promise
            .then((res) => {
                setDrawRes(JSON.stringify(res.data))
            })
            .catch((err) => {
                setDrawRes(JSON.stringify(err))
            })
    }

    return (
        <div>
            <h2>Little utility to quickly test the Card Games Engine backend API.</h2>

            <div>
                <span>
                    <button onClick={handleCreate}>create deck</button>
                    cards:{" "}
                    <input type="text" onChange={(e) => setCards(e.target.value)} />
                    shuffle?
                    <input
                        type="checkbox"
                        onChange={(e) => setShuffled(e.target.checked)}
                    />
                    <div>Create response: {createRes}</div>
                </span>
            </div>

            <div>
                deck guid to open/draw:{" "}
                <input
                    type="text"
                    style={{ width: 400 }}
                    onChange={(e) => setGuid(e.target.value)}
                />
            </div>

            <div>
                <span>
                    <button onClick={handleOpen}>open deck</button>
                    <div>Open response: {openRes}</div>
                </span>
            </div>

            <div>
                <span>
                    <button onClick={handleDraw}>draw</button>
                    number of cards to drawn:{" "}
                    <input type="number" onChange={(e) => setCount(e.target.value)} />
                </span>
                <div>Draw response: {drawRes}</div>
            </div>
        </div>
    )
}

export default App
