import React, { useState } from "react";
import "./App.css";
import cover from "./assets/cover.jpg";

function App() {
    const [game, setGame] = useState(null);
    const [player, setPlayer] = useState(null);
    return (
        <div className="App">
            <Debug game={game} player={player} />
            <img src={cover} className="cover" />
            <div className="buttons">
                <NewGame setGame={setGame} setPlayer={setPlayer} />
                <JoinGame setGame={setGame} setPlayer={setPlayer} />
            </div>
            <div></div>
        </div>
    );
}

function Debug({ game, player }) {
    return (
        <div className="debug">
            <div>Current game: {game}</div>
            <div>Current player: {player}</div>
            <a
                href={`http://localhost:8000/game/${game}/${player}`}
                target="_blank"
                rel="noopener noreferrer"
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
            .then(res => {
                if (!res.ok) {
                    throw res.statusText; // TODO: manage the error
                }
                return res.json();
            })
            .then(res => {
                setGame(res.game);
                setPlayer(res.player);
            })
            .catch(err => console.error(err));
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
        <button className="button" onClick={() => setShowForm(true)}>
            Join game
        </button>
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
        <button className="button" onClick={() => setShowForm(true)}>
            New game
        </button>
    );
}

export default App;
