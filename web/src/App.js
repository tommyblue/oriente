import React, { useState } from "react";
import Intro from "./intro";
import Game from "./game";
import "./App.css";
import "./fonts/font.ttf";

function App() {
    const [game, setGame] = useState(null);
    const [player, setPlayer] = useState(null);
    return (
        <div className="App">
            <Debug game={game} player={player} />
            {game && player ? (
                <Game />
            ) : (
                <Intro setGame={setGame} setPlayer={setPlayer} />
            )}
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

export default App;
