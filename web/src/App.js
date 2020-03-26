import React, { useState } from "react";
import "./App.css";

function App() {
    const [game, setGame] = useState(null);
    const [player, setPlayer] = useState(null);
    return (
        <div className="App">
            Current game: {game}
            <br />
            Current player: {player}
            <br />
            <NewGame setGame={setGame} setPlayer={setPlayer} />
            <JoinGame setGame={setGame} setPlayer={setPlayer} />
            <a
                href={`http://localhost:8000/game/${game}/${player}`}
                target="_blank"
            >
                Debug
            </a>
        </div>
    );
}

function JoinGame({ setGame, setPlayer }) {
    const [showForm, setShowForm] = useState(false);
    const [gameID, setGameID] = useState(null);
    const joinGame = () => {
        fetch(`http://localhost:8000/game/${gameID}`)
            .then(res => res.json())
            .then(res => {
                setGame(res.game);
                setPlayer(res.player);
            });
    };
    if (showForm) {
        return (
            <form
                onSubmit={e => {
                    e.preventDefault();
                    joinGame();
                }}
            >
                Game ID:{" "}
                <input type="text" onChange={e => setGameID(e.target.value)} />
                <button>Join!</button>
            </form>
        );
    }
    return (
        <div>
            <button onClick={() => setShowForm(true)}>Join game</button>
        </div>
    );
}

function NewGame({ setGame, setPlayer }) {
    const [showForm, setShowForm] = useState(false);
    const [nPlayers, setNPlayers] = useState(null);
    const setGameId = () => {
        fetch(`http://localhost:8000/game/new/${nPlayers}`)
            .then(res => res.json())
            .then(res => {
                setGame(res.game);
                setPlayer(res.player);
            });
    };
    if (showForm) {
        return (
            <form
                onSubmit={e => {
                    e.preventDefault();
                    setGameId();
                }}
            >
                Number of players:{" "}
                <input
                    type="number"
                    min="4"
                    max="12"
                    onChange={e => setNPlayers(e.target.value)}
                />
                <button>Start!</button>
            </form>
        );
    }
    return (
        <div>
            <button onClick={() => setShowForm(true)}>New game</button>
        </div>
    );
}

export default App;
